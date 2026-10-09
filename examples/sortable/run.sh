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
import base64, datetime, fcntl, hashlib, json, os, pathlib, pwd, re, shutil, subprocess, sys, time, urllib.request, uuid

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
        fail(json.dumps(dict(error='command_failed',argv=[str(x) for x in argv],exit_code=result.returncode,stdout=result.stdout[-4000:],stderr=result.stderr[-4000:])))
    return result

def kernel(*args, check=True):
    return command([LAUNCHER, *args], check=check)

def write(path, value):
    path.write_text(json.dumps(value, indent=2)+'\n')

def read(path):
    def unique(pairs):
        value={}
        for key,item in pairs:
            if key in value:
                fail('Duplicate lifecycle artifact key: '+key)
            value[key]=item
        return value
    return json.loads(path.read_text(),object_pairs_hook=unique,parse_constant=lambda x: fail('Invalid JSON constant: '+x))

def sha(data):
    return hashlib.sha256(data).hexdigest()

def exact_object(value, fields):
    if type(value) is not dict or set(value)!=set(fields.split()):
        fail('Malformed lifecycle artifact object.')

def string_list(value):
    if type(value) is not list or any(type(x) is not str for x in value):
        fail('Malformed lifecycle artifact list.')

def validate_inventory(value, source):
    exact_object(value,'schema_version kind resources queried absent source evidence')
    if type(value['schema_version']) is not int or value['schema_version']!=1 or value['kind']!='docker-inventory' or value['source']!=source:
        fail('Unsupported Docker inventory schema/source.')
    for field in ['resources','queried','absent','evidence']:
        if type(value[field]) is not list:
            fail('Malformed Docker inventory collection.')
    for field in ['queried','absent']:
        for row in value[field]:
            string_list(row)
            if len(row)!=2 or row[0] not in ['container','network','volume'] or not row[1]:
                fail('Malformed Docker inventory identity.')
    identities=set()
    for resource in value['resources']:
        exact_object(resource,'resource_id kind labels created_at running in_use system inspect_ok name use_known')
        identity=(resource['kind'],resource['resource_id'])
        if resource['kind'] not in ['container','network','volume'] or type(resource['resource_id']) is not str or not resource['resource_id'] or identity in identities:
            fail('Malformed or repeated Docker resource identity.')
        identities.add(identity)
        if type(resource['labels']) is not dict or any(type(k) is not str or type(v) is not str for k,v in resource['labels'].items()):
            fail('Malformed Docker resource labels.')
        if any(type(resource[k]) is not bool for k in ['running','in_use','system','inspect_ok','use_known']) or type(resource['name']) is not str:
            fail('Malformed Docker resource state.')
        if type(resource['created_at']) is not str or datetime.datetime.fromisoformat(resource['created_at'].replace('Z','+00:00')).tzinfo is None:
            fail('Malformed Docker resource timestamp.')
    for result in value['evidence']:
        exact_object(result,'schema_version argv exit_code stdout stderr')
        string_list(result['argv'])
        if type(result['schema_version']) is not int or result['schema_version']!=1 or type(result['exit_code']) is not int or any(type(result[k]) is not str for k in ['stdout','stderr']):
            fail('Malformed Docker observation evidence.')

def validate_cleanup_artifact(value, run_id, node_id, scope_id, stale):
    exact_object(value,'schema_version kind plan inventory')
    if type(value['schema_version']) is not int or value['schema_version']!=1 or value['kind']!='cleanup-plan-artifact':
        fail('Unsupported cleanup artifact schema.')
    plan=value['plan'];exact_object(plan,'schema_version scope before actions dispositions')
    expected=dict(run_id='stale' if stale else run_id,terminal=not stale,stale_sweep=stale,repository_scope_id=scope_id)
    if type(plan['schema_version']) is not int or plan['schema_version']!=1 or plan['scope']!=expected or any(type(plan['scope'][k]) is not bool for k in ['terminal','stale_sweep']):
        fail('Cleanup plan is not bound to the selected run/repository scope.')
    string_list(plan['before'])
    for field in ['actions','dispositions']:
        if type(plan[field]) is not list:
            fail('Malformed cleanup plan collection.')
        for row in plan[field]:
            if field=='actions':
                exact_object(row,'resource_id kind action argv requires_success_of owner lifecycle proof_digest preconditions environment predecessor_result_id evidence_digest')
                string_list(row['argv']);string_list(row['preconditions'])
                if row['action'] not in ['stop','remove'] or not row['argv'] or not row['preconditions'] or len(set(row['preconditions']))!=len(row['preconditions']) or type(row['environment']) is not dict or any(type(k) is not str or type(v) is not str for k,v in row['environment'].items()):
                    fail('Malformed cleanup action.')
                for key in ['proof_digest','evidence_digest','predecessor_result_id']:
                    if key=='predecessor_result_id' and row[key] is None:
                        continue
                    if type(row[key]) is not str or re.fullmatch(r'sha256:[0-9a-f]{64}',row[key]) is None:
                        fail('Malformed cleanup proof digest.')
                dependency=row['requires_success_of']
                if dependency is not None and (type(dependency) is not int or dependency<0):
                    fail('Malformed cleanup dependency.')
            else:
                fields='resource_id kind owner lifecycle disposition action reason command_evidence evidence'
                exact_object(row,fields+(' follow_up' if 'follow_up' in row else ''))
                if row['disposition'] not in ['removed','retained_for_dependency','blocked','foreign','missing'] or row['action'] not in ['none','remove_exact_id'] or type(row['reason']) is not str or any(type(row[k]) is not list for k in ['command_evidence','evidence']):
                    fail('Malformed cleanup disposition.')
            exact_object(row['owner'],'run_id node_id')
            if row['kind'] not in ['container','network','volume'] or type(row['resource_id']) is not str or not row['resource_id'] or row['lifecycle'] not in ['run','chunk'] or any(type(x) is not str or not x for x in row['owner'].values()):
                fail('Malformed cleanup ownership.')
            if not stale and (row['owner']!=dict(run_id=run_id,node_id=node_id) or row['lifecycle']!='run'):
                fail('Cleanup row is outside the selected lifecycle owner.')
    validate_inventory(value['inventory'],'managed_orphan_sweep' if stale else 'registered_exact')
    if not stale:
        for resource in value['inventory']['resources']:
            labels=resource['labels']
            if any(labels.get('com.designmachines.depot.'+key)!=expected for key,expected in [('managed','true'),('run-id',run_id),('node-id',node_id),('repository-scope-id',scope_id),('lifecycle','run')]):
                fail('Current-run inventory ownership mismatch.')
    return any(row['disposition'] in ['blocked','retained_for_dependency'] for row in plan['dispositions'])

def load_reconciliation(descriptor, status, run_id, node_id, scope_id):
    # Status is the maximum of two independent plans, not current-run authority.
    if status not in [0,3]:
        fail('Unexpected plan-reconcile status '+str(status)+'; evidence: '+str(descriptor))
    paths=[descriptor,descriptor.with_name(descriptor.stem+'.current-run.json'),descriptor.with_name(descriptor.stem+'.stale-sweep.json')]
    for path in paths:
        if not path.is_file() or path.resolve()!=path or path.stat().st_uid!=os.getuid():
            fail('Missing or unsafe reconciliation artifact: '+str(path))
    document=read(descriptor);exact_object(document,'schema_version kind current_run_plan stale_sweep_plan ttl_hours')
    if type(document['schema_version']) is not int or document['schema_version']!=1 or document['kind']!='cleanup-plan-set' or document['ttl_hours']!=876000 or document['current_run_plan']!=str(paths[1]) or document['stale_sweep_plan']!=str(paths[2]):
        fail('Reconciliation descriptor schema/path binding mismatch.')
    current=read(paths[1]);sweep=read(paths[2])
    blocked=validate_cleanup_artifact(current,run_id,node_id,scope_id,False)
    stale_blocked=validate_cleanup_artifact(sweep,run_id,node_id,scope_id,True)
    if status!=(3 if blocked or stale_blocked else 0):
        fail('Reconciliation status disagrees with emitted dispositions.')
    return current,sweep,blocked

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

    compose_env=evidence/'compose.env'
    def write_compose_env(receipt):
        if compose_env.is_symlink():
            fail('Unsafe Compose env file.')
        values=dict(SORTABLE_RUN_DIR=str(run),SORTABLE_PROJECT=m['project'],SORTABLE_BUILD_RECEIPT=receipt)
        # Single-quoted dotenv values do not expand dollar signs. Kernel's
        # fixed Docker environment need not forward any application variables.
        compose_env.write_text(''.join(key+"='"+value.replace("'", "\\'")+"'\n" for key,value in values.items()))
        for key in values:
            os.environ.pop(key,None)
    write_compose_env('not-building')
    compose=['docker','compose','--env-file',str(compose_env),'--project-name',m['project'],'--project-directory',str(checkout),'-f',str(checkout/'examples/sortable/compose.yaml')]
    state=pathlib.Path(m['kernelStateDir'])
    if state!=checkout/'.workflow-kernel'/'runs'/m['runId']:
        fail('Kernel state directory is outside the selected run.')

    def docker_preflight():
        info=json.loads(command(['docker','info','--format','{{json .}}']).stdout)
        if info.get('DockerRootDir')!='/home/ned/.local/share/docker' or 'name=rootless' not in info.get('SecurityOptions',[]):
            fail('Verified rootless development Docker daemon required.')

    def initialize_kernel():
        if not (state/'run-state.json').exists():
            kernel('init',state,'--run-id',m['runId'],'--mode','shadow','--occurred-at',datetime.datetime.now(datetime.timezone.utc).isoformat())

    def reconcile(descriptor, allow_blocked=False):
        if descriptor.parent!=evidence or descriptor.exists():
            fail('Reconciliation requires a fresh owned descriptor path.')
        result=kernel('plan-reconcile','--state-dir',state,'--run-id',m['runId'],'--ttl-hours','876000','--output',descriptor,check=False)
        scope_id=read(checkout/'.workflow-kernel/repository-scope.json')['scope_id']
        current,sweep,blocked=load_reconciliation(descriptor,result.returncode,m['runId'],m['nodeId'],scope_id)
        if blocked and not allow_blocked:
            fail('Selected current-run cleanup is blocked; evidence: '+str(descriptor))
        return current,sweep,blocked

    def inventory(prefix, allow_blocked=False):
        # Kernel's inventory projection, not a competing resource registry.
        descriptor=evidence/(prefix+'-inventory-plans.json')
        current,sweep,blocked=reconcile(descriptor,allow_blocked=allow_blocked)
        observation=sweep['inventory']
        # The registered-only projection cannot see partial, not-yet-recorded
        # creation. Use sweep observations, never its cleanup actions.
        observation=dict(observation,resources=[r for r in observation['resources'] if r['labels'].get('com.designmachines.depot.run-id')==m['runId']])
        for resource in observation['resources']:
            if any(resource['labels'].get('com.designmachines.depot.'+key)!=expected for key,expected in [('managed','true'),('node-id',m['nodeId']),('repository-scope-id',current['plan']['scope']['repository_scope_id']),('lifecycle','run')]):
                fail('Creation observation ownership mismatch.')
        path=evidence/(prefix+'-inventory.json')
        write(path,observation)
        return path,blocked

    def create(argv):
        token=uuid.uuid4().hex
        argv_path=evidence/(token+'-argv.json');write(argv_path,argv)
        deps=evidence/(token+'-dependencies.json');write(deps,[])
        plan_path=evidence/(token+'-creation-plan.json')
        before,_=inventory(token+'-before')
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
        after,_=inventory(token+'-after',allow_blocked=True)
        receipt=kernel('record-create','--state-dir',state,'--plan',plan_path,'--result',result_path,'--before-inventory',before,'--after-inventory',after,check=False)
        (evidence/(token+'-creation-receipt.json')).write_text(receipt.stdout)
        if result.returncode or receipt.returncode:
            fail('Creation failed; record-create was attempted. Inspect its receipt before stop. Evidence: '+str(result_path))
        # Registration can resolve a network dependency on the newly created
        # container. The pre-registration plan is observation, not readiness.
        kernel('validate-resource-registry','--state-dir',state,'--run-id',m['runId'],'--node-id',m['nodeId'])
        reconcile(evidence/(token+'-registered-inventory-plans.json'))

    def stop():
        if not (state/'run-state.json').exists():
            return
        descriptor=evidence/('stop-'+uuid.uuid4().hex+'.json')
        reconcile(descriptor)
        plan_path=pathlib.Path(read(descriptor)['current_run_plan'])
        # Never execute the unrelated stale-sweep plan.
        outcomes=evidence/(descriptor.stem+'-outcomes.json');write(outcomes,[])
        step_path=evidence/(descriptor.stem+'-step.json')
        while True:
            kernel('next-cleanup-step','--state-dir',state,'--plan',plan_path,'--outcomes',outcomes,'--output',step_path)
            step=read(step_path)
            if step['complete']:
                break
            fresh=evidence/(descriptor.stem+'-fresh-'+uuid.uuid4().hex+'.json')
            fresh_plan,_,_=reconcile(fresh)
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
        write_compose_env(encoded)
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
