#!/usr/bin/env bash
set -euo pipefail
# The repository owns the lifecycle. Kernel instruments each exact creating
# Docker argv and owns registration/cleanup; no raw Compose up/down fallback.
if [[ $# != 3 || "$2" != --run-dir ]]; then
  printf 'usage: run.sh assets|start|status|rebuild|stop --run-dir ABSOLUTE_OWNED_ROOT\n' >&2
  exit 2
fi
case "$1" in assets|start|status|rebuild|stop) ;; *) exit 2 ;; esac
export DOCKER_CONTEXT=rootless
unset DOCKER_HOST
exec python3 - "$1" "$3" "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd -P)" <<'PY'
import base64, datetime, fcntl, hashlib, json, os, pathlib, pwd, shutil, subprocess, sys, time, urllib.request, uuid

action, run_arg, checkout_arg = sys.argv[1:]
checkout = pathlib.Path(checkout_arg)
run = pathlib.Path(run_arg)
APP = 'livewires-templ-sortable'
LAUNCHER = pathlib.Path('/home/ned/.codex/plugins/cache/depot/workflow-kernel/0.27.0/skills/workflow-kernel/references/workflow-kernel-launcher.sh')

def fail(message):
    raise SystemExit(message)

def command(argv, check=True):
    result = subprocess.run([str(x) for x in argv], cwd=checkout, capture_output=True, text=True)
    if check and result.returncode:
        fail('Command failed: '+repr(argv)+'\n'+result.stderr[-4000:])
    return result

def kernel(*args, check=True):
    return command([LAUNCHER, *args], check=check)

def write(path, value):
    path.write_text(json.dumps(value, indent=2)+'\n')

def read(path):
    return json.loads(path.read_text())

def sha(data):
    return hashlib.sha256(data).hexdigest()

if not run.is_absolute() or not run.is_dir() or run.resolve()!=run or run.stat().st_uid!=os.getuid():
    fail('Use the fresh, existing, ned-owned absolute run root allocated by the host.')
if run==checkout or run.is_relative_to(checkout) or checkout.is_relative_to(run):
    fail('Run root must be external to the selected physical checkout.')
if pwd.getpwuid(os.getuid()).pw_name!='ned' or command(['hostname','-s']).stdout.strip()!='ned9000':
    fail('This lifecycle is for the verified ned/rootless development boundary on NED.')
if not LAUNCHER.is_file() or not os.access(LAUNCHER,os.X_OK):
    fail('Trusted installed Workflow Kernel launcher is unavailable.')
kernel('kernel-info','--minimum-version','0.27.0')

# Serialize lifecycle invocations for this exact root, without locking the
# application store (which remains entirely owned by the running Go process).
lock_file = run / 'sortable-lifecycle.lock'
if lock_file.is_symlink():
    fail('Unsafe lifecycle lock.')
with lock_file.open('a') as lifecycle_lock:
    fcntl.flock(lifecycle_lock,fcntl.LOCK_EX)
    evidence = run / 'evidence'
    manifest = evidence / 'lifecycle.json'
    if not manifest.exists():
        if action!='assets':
            fail('Run assets first using this exact-owned root.')
        # These helpers verify the root allocation and record only absent paths.
        receipts=[]
        for name, kind in [('evidence','raw-output'),('assets','temporary-directory'),('data','temporary-directory'),('cache','cache'),('build','temporary-directory')]:
            receipts.append(json.loads(kernel('owned-run-create','--run-root',run,'--kind',kind,'--relative-path',name).stdout))
        run_id = 'sortable-'+uuid.uuid4().hex
        state = checkout / '.workflow-kernel' / 'runs' / run_id
        m=dict(application=APP,physicalCheckout=str(checkout),runDirectory=str(run),runId=run_id,nodeId='sortable-example',project=run_id,kernelStateDir=str(state),launcher=str(LAUNCHER),filesystemReceipts=receipts)
        write(manifest,m)
    m=read(manifest)
    if m['application']!=APP or m['physicalCheckout']!=str(checkout) or m['runDirectory']!=str(run) or m['launcher']!=str(LAUNCHER):
        fail('Lifecycle binding does not match this checkout/run/launcher.')
    for receipt in m['filesystemReceipts']:
        path=pathlib.Path(receipt['path'])
        if not path.is_dir() or path.resolve()!=path or path.stat().st_uid!=os.getuid():
            fail('Recorded task directory is missing or unsafe.')
    locked=(checkout/'examples/sortable/assets.lock.json').read_bytes()
    asset_lock=json.loads(locked)

    def assets(acquire=False):
        for asset in asset_lock['assets']:
            relative=pathlib.PurePosixPath(asset['file'])
            if relative.is_absolute() or '..' in relative.parts:
                fail('Unsafe asset lock path.')
            target=run/'assets'/relative
            if target.exists():
                if target.resolve()!=target or not target.is_file() or sha(target.read_bytes())!=asset['sha256']:
                    fail('Unverified existing asset: '+asset['file'])
                continue
            if not acquire:
                fail('Missing locked asset: '+asset['file']+'; run assets.')
            data=urllib.request.urlopen(asset['url'],timeout=30).read()
            if len(data)!=asset['size'] or sha(data)!=asset['sha256']:
                fail('Downloaded bytes disagree with lock: '+asset['file'])
            target.parent.mkdir(parents=True,exist_ok=True)
            temporary=target.with_suffix(target.suffix+'.download')
            with temporary.open('xb') as f:
                f.write(data);f.flush();os.fsync(f.fileno())
            temporary.replace(target)
            target.chmod(0o444)

    def source():
        # Only files contributing to the executable or its serving contract.
        names=command(['git','ls-files','-z','--cached','--others','--exclude-standard']).stdout.split('\0')
        names=sorted(set(n for n in names if n and (
            (n.endswith(('.go','.templ')) and not n.endswith('_test.go'))
            or n in ['go.mod','go.sum','Dockerfile','examples/sortable/bridge.js','examples/sortable/assets.lock.json','examples/sortable/compose.yaml','examples/sortable/run.sh'])))
        hashes={n:sha((checkout/n).read_bytes()) for n in names}
        fingerprint=sha(json.dumps(hashes,sort_keys=True,separators=(',',':')).encode())
        dirt=command(['git','status','--porcelain','--',*names]).stdout
        return dict(application=APP,physicalCheckout=str(checkout),branch=command(['git','branch','--show-current']).stdout.strip(),commit=command(['git','rev-parse','HEAD']).stdout.strip(),sourceFingerprint=fingerprint,dirty=bool(dirt),lockDigest=sha(locked)),hashes

    os.environ.update(SORTABLE_RUN_DIR=str(run),SORTABLE_PROJECT=m['project'],SORTABLE_BUILD_RECEIPT='not-building')
    compose=['docker','compose','--project-name',m['project'],'--project-directory',str(checkout),'-f',str(checkout/'examples/sortable/compose.yaml')]
    state=pathlib.Path(m['kernelStateDir'])

    def docker_preflight():
        info=json.loads(command(['docker','info','--format','{{json .}}']).stdout)
        if info.get('DockerRootDir')!='/home/ned/.local/share/docker' or 'name=rootless' not in info.get('SecurityOptions',[]):
            fail('Verified rootless development Docker daemon required.')

    def initialize_kernel():
        if not (state/'run-state.json').exists():
            kernel('init',state,'--run-id',m['runId'],'--mode','shadow','--occurred-at',datetime.datetime.now(datetime.timezone.utc).isoformat())

    def inventory(prefix):
        # Kernel's inventory projection, not a competing resource registry.
        descriptor=evidence/(prefix+'-inventory-plans.json')
        kernel('plan-reconcile','--state-dir',state,'--run-id',m['runId'],'--ttl-hours','876000','--output',descriptor)
        plans=read(descriptor)
        path=evidence/(prefix+'-inventory.json')
        write(path,read(pathlib.Path(plans['stale_sweep_plan']))['inventory'])
        return path

    def create(argv):
        token=uuid.uuid4().hex
        argv_path=evidence/(token+'-argv.json');write(argv_path,argv)
        deps=evidence/(token+'-dependencies.json');write(deps,[])
        plan_path=evidence/(token+'-creation-plan.json')
        before=inventory(token+'-before')
        kernel('plan-compose','--state-dir',state,'--run-id',m['runId'],'--node-id',m['nodeId'],'--lifecycle','run','--cleanup-policy','stop-remove','--repository-project-name',m['project'],'--argv-json',argv_path,'--dependent-node-ids-json',deps,'--output',plan_path)
        plan=read(plan_path)
        if not plan['managed'] or plan['project_name']!=m['project']:
            fail('Kernel declined this exact Docker boundary.')
        override=checkout/plan['compose_override']
        if not override.is_relative_to(checkout/'.workflow-kernel'):
            fail('Unsafe returned override path.')
        override.parent.mkdir(parents=True,exist_ok=True)
        override.write_text(plan['compose_override_content'])
        os.environ.update(plan['environment'] or {})
        result=command(plan['argv'],check=False)
        result_path=evidence/(token+'-result.json')
        write(result_path,dict(schema_version=1,argv=plan['argv'],exit_code=result.returncode,stdout=result.stdout,stderr=result.stderr))
        after=inventory(token+'-after')
        receipt=kernel('record-create','--state-dir',state,'--plan',plan_path,'--result',result_path,'--before-inventory',before,'--after-inventory',after,check=False)
        (evidence/(token+'-creation-receipt.json')).write_text(receipt.stdout)
        if result.returncode or receipt.returncode:
            fail('Creation failed; partial resources were recorded. Run stop. Evidence: '+str(result_path))

    def stop():
        if not (state/'run-state.json').exists():
            return
        descriptor=evidence/('stop-'+uuid.uuid4().hex+'.json')
        kernel('plan-reconcile','--state-dir',state,'--run-id',m['runId'],'--ttl-hours','876000','--output',descriptor)
        plan_path=pathlib.Path(read(descriptor)['current_run_plan'])
        # Never execute the unrelated stale-sweep plan.
        outcomes=evidence/(descriptor.stem+'-outcomes.json');write(outcomes,[])
        step_path=evidence/(descriptor.stem+'-step.json')
        while True:
            kernel('next-cleanup-step','--state-dir',state,'--plan',plan_path,'--outcomes',outcomes,'--output',step_path)
            step=read(step_path)
            if step['complete']:
                break
            fresh=evidence/(descriptor.stem+'-fresh.json')
            kernel('plan-reconcile','--state-dir',state,'--run-id',m['runId'],'--ttl-hours','876000','--output',fresh)
            fresh_plan=read(pathlib.Path(read(fresh)['current_run_plan']))
            witness=evidence/(descriptor.stem+'-witness.json');write(witness,fresh_plan['inventory'])
            run_state=read(state/'run-state.json')
            statuses=evidence/(descriptor.stem+'-nodes.json')
            write(statuses,dict(schema_version=1,run_id=m['runId'],revision=run_state['revision'],updated_at=run_state['updated_at'],node_statuses={k:v['status'] for k,v in run_state['nodes'].items()}))
            result_path=evidence/(descriptor.stem+'-step-result.json')
            kernel('execute-cleanup-step','--state-dir',state,'--plan',plan_path,'--step-index',str(step['step_index']),'--inventory',witness,'--node-statuses',statuses,'--outcomes',outcomes,'--output',result_path)
            prior=read(outcomes);prior.append(read(result_path));write(outcomes,prior)
        result=kernel('record-cleanup','--state-dir',state,'--plan',plan_path,'--outcomes',outcomes)
        (evidence/(descriptor.stem+'-receipt.json')).write_text(result.stdout)

    def status():
        current,_=source()
        result=dict(**current,actualOrigin=None,consumerUrl=None,referenceUrl=None,health='stopped',producerCommit=asset_lock['producerCommit'],referenceHTMLDigest=next(x['sha256'] for x in asset_lock['assets'] if x['serve']=='/manual/components/sortable-list.html'),loadedAssets=[],compiledSource=None,sourceMatchesBuild=False,project=m['project'],runDirectory=str(run),kernelStateDir=str(state))
        if not (evidence/'compiled-receipt.json').exists():
            return result
        compiled=read(evidence/'compiled-receipt.json');result['compiledSource']=compiled
        mapped=command(compose+['port','app','8080'],check=False)
        origin='http://'+mapped.stdout.strip()
        if mapped.returncode or not mapped.stdout.strip():
            return result
        if not origin.startswith('http://127.0.0.1:') or not origin.removeprefix('http://127.0.0.1:').isdigit():
            fail('Unexpected host publication; loopback dynamic port required.')
        kernel('validate-resource-registry','--state-dir',state,'--run-id',m['runId'],'--node-id',m['nodeId'])
        result.update(actualOrigin=origin,consumerUrl=origin+'/sortable',referenceUrl=origin+'/manual/components/sortable-list.html')
        try:
            health=json.load(urllib.request.urlopen(origin+'/healthz',timeout=3))
            result['compiledSource']=health['source']
            result['loadedAssets']=health['assets']
            result['sourceMatchesBuild']=health['source']==compiled and all(current[k]==compiled[k] for k in ['physicalCheckout','branch','commit','sourceFingerprint','dirty','lockDigest'])
            expected=[{k:x[k] for k in ['file','url','sha256','serve']} for x in asset_lock['assets']]
            result['health']='ready' if health['healthy'] and result['sourceMatchesBuild'] and health['producerCommit']==asset_lock['producerCommit'] and health['assets']==expected else 'source-mismatch'
        except (OSError,ValueError,KeyError):
            result['health']='unreachable'
        return result

    if action=='assets':
        assets(acquire=True)
        print(json.dumps(dict(application=APP,status='verified',runDirectory=str(run),lockDigest=sha(locked),assets=len(asset_lock['assets']))))
    elif action=='status':
        docker_preflight();print(json.dumps(status()))
    elif action=='stop':
        docker_preflight();stop();print(json.dumps(dict(application=APP,status='stopped',store=str(run/'data/orders.json'),filesystemCleanupOwner='host/root')))
    else:
        docker_preflight();assets();initialize_kernel()
        if action=='rebuild':
            stop()
        elif command(compose+['ps','--all','--quiet'],check=False).stdout.strip():
            fail('Application resources already exist. Use status or rebuild; start never adopts an existing instance.')
        receipt,hashes=source()
        snapshot=run/'build/source'
        if snapshot.exists():
            # Only this creation-recorded build child, never data/assets.
            shutil.rmtree(snapshot)
        snapshot.mkdir()
        for name,expected in hashes.items():
            data=(checkout/name).read_bytes()
            if sha(data)!=expected:
                fail('Source changed while snapshotting; retry rebuild.')
            target=snapshot/name;target.parent.mkdir(parents=True,exist_ok=True);target.write_bytes(data)
        if source()[0]!=receipt:
            fail('Source changed during build preparation.')
        for folder in ['modules','go-build']:
            (run/'cache'/folder).mkdir(exist_ok=True)
        encoded=base64.urlsafe_b64encode(json.dumps(receipt,separators=(',',':')).encode()).decode().rstrip('=')
        os.environ['SORTABLE_BUILD_RECEIPT']=encoded
        # Config/assets/source validation all precede application listening.
        command(compose+['--profile','*','config','--quiet'])
        create(compose+['run','--rm','--build','builder'])
        write(evidence/'compiled-receipt.json',receipt)
        binary=run/'build/sortable';write(evidence/'binary-receipt.json',dict(sha256=sha(binary.read_bytes()),source=receipt))
        create(compose+['up','--detach','--wait','--wait-timeout','30','app'])
        image=command(['docker','image','inspect',m['project']+'-toolchain','--format','{{.Id}}']).stdout.strip()
        write(evidence/'image-receipt.json',dict(image=image,tag=m['project']+'-toolchain',cleanupOwner='host/root; verify exact recorded image before removal'))
        final=status();print(json.dumps(final))
        if final['health']!='ready':
            fail('Source-bound application readiness failed; preserve evidence and use stop/rebuild.')
PY
