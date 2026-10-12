#!/usr/bin/env python3
"""Repository-owned gallery lifecycle. Kernel owns Docker registration/cleanup.

No review state or approval is consulted. A host-selected installed Kernel is
required only for exact Docker ownership, as in the existing sortable runner.
"""
import argparse
import base64
import datetime
import fcntl
import hashlib
import json
import os
import pathlib
import pwd
import shutil
import subprocess
import urllib.request
import uuid

CHECKOUT = pathlib.Path(__file__).resolve().parents[2]


def fail(message):
    raise SystemExit(message)


def read(path):
    if path.is_symlink():
        fail("Refusing symlink: " + str(path))
    return json.loads(path.read_text())


def write(path, value):
    if path.is_symlink():
        fail("Refusing symlink: " + str(path))
    path.write_text(json.dumps(value, indent=2) + "\n")


def sha(data):
    return hashlib.sha256(data).hexdigest()


def command(argv, check=True):
    result = subprocess.run([str(x) for x in argv], cwd=CHECKOUT, capture_output=True, text=True)
    if check and result.returncode:
        fail(json.dumps(dict(argv=[str(x) for x in argv], exit=result.returncode,
                             stdout=result.stdout[-4000:], stderr=result.stderr[-4000:])))
    return result


def resolve_kernel():
    # The installed host manifest selects the trusted bundle. Never glob caches
    # or select an executable by mtime. Other hosts pass their selected launcher.
    explicit = os.environ.get("WORKFLOW_KERNEL")
    if explicit:
        launcher = pathlib.Path(explicit)
    else:
        account = pathlib.Path(pwd.getpwuid(os.getuid()).pw_dir)
        manifest = account / ".claude/plugins/installed_plugins.json"
        if not manifest.is_file():
            fail("Set WORKFLOW_KERNEL to the launcher in your host-selected installed workflow-kernel bundle; see README.md.")
        entries = read(manifest).get("plugins", {}).get("workflow-kernel@depot", [])
        selected = [row for row in entries if row.get("scope") == "user"]
        if len(selected) != 1:
            fail("No unique user-installed Workflow Kernel; pass host-selected WORKFLOW_KERNEL.")
        root = pathlib.Path(selected[0]["installPath"])
        plugin = read(root / ".claude-plugin/plugin.json")
        if plugin.get("name") != "workflow-kernel" or plugin.get("version") != selected[0]["version"]:
            fail("Installed Workflow Kernel manifest mismatch.")
        launcher = root / "skills/workflow-kernel/references/workflow-kernel-launcher.sh"
    if not launcher.is_absolute() or not launcher.is_file() or not os.access(launcher, os.X_OK):
        fail("Host-selected Workflow Kernel launcher unavailable.")
    command([launcher, "kernel-info", "--minimum-version", "0.25.0"])
    return launcher


def source():
    lock = (CHECKOUT / "examples/gallery/assets.lock.json").read_bytes()
    names = command(["git", "ls-files", "-z", "--cached", "--others", "--exclude-standard"]).stdout.split("\0")
    # Freeze all build inputs, including generated Go and repository-owned JS.
    names = sorted(set(n for n in names if n and not n.startswith(("plans/", ".workflow-kernel/", ".gallery/"))
                       and (n.endswith((".go", ".templ", ".js", ".json", ".yaml", ".py", ".sh", ".html"))
                            or n in ["go.mod", "go.sum", "Dockerfile"])))
    hashes = {n: sha((CHECKOUT / n).read_bytes()) for n in names if (CHECKOUT / n).is_file()}
    identity = dict(checkout=str(CHECKOUT), branch=command(["git", "branch", "--show-current"]).stdout.strip(),
                    commit=command(["git", "rev-parse", "HEAD"]).stdout.strip(),
                    dirty=bool(command(["git", "status", "--porcelain"]).stdout),
                    fingerprint=sha(json.dumps(hashes, sort_keys=True, separators=(",", ":")).encode()),
                    assetRevision=json.loads(lock)["producerCommit"], lockDigest=sha(lock))
    return identity, hashes


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("action", choices=["assets", "start", "status", "rebuild", "stop", "verify"])
    args = parser.parse_args()
    os.environ["DOCKER_CONTEXT"] = "rootless"
    os.environ.pop("DOCKER_HOST", None)
    info = json.loads(command(["docker", "info", "--format", "{{json .}}"]).stdout)
    if "name=rootless" not in info.get("SecurityOptions", []):
        fail("A rootless development Docker daemon is required.")
    launcher = resolve_kernel()

    def kernel(*argv, check=True):
        return command([launcher, *argv], check=check)

    local = CHECKOUT / ".gallery"
    if local.is_symlink():
        fail("Unsafe gallery state directory.")
    local.mkdir(exist_ok=True)
    with (local / "lifecycle.lock").open("a") as serial:
        fcntl.flock(serial, fcntl.LOCK_EX)
        pointer = local / "instance.json"
        if not pointer.exists():
            if args.action in ["status", "stop"]:
                print(json.dumps(dict(health="stopped", checkout=str(CHECKOUT))))
                return
            run_id = "gallery-" + uuid.uuid4().hex
            run = pathlib.Path(json.loads(kernel("owned-run-start", "--workflow", "gallery", "--run-id", run_id).stdout)["path"])
            for name, kind in [("evidence", "raw-output"), ("assets", "temporary-directory"), ("build", "temporary-directory"), ("cache", "cache")]:
                kernel("owned-run-create", "--run-root", run, "--kind", kind, "--relative-path", name)
            state = CHECKOUT / ".workflow-kernel/runs" / run_id
            kernel("init", state, "--run-id", run_id, "--mode", "shadow", "--occurred-at", datetime.datetime.now(datetime.timezone.utc).isoformat())
            write(pointer, dict(checkout=str(CHECKOUT), run=str(run), runId=run_id, state=str(state), project=run_id, node="gallery"))
        instance = read(pointer)
        if instance["checkout"] != str(CHECKOUT):
            fail("Gallery instance belongs to another checkout.")
        run = pathlib.Path(instance["run"])
        if not run.is_dir() or run.resolve() != run or run.stat().st_uid != os.getuid():
            fail("Missing or unsafe owned run root.")
        state = pathlib.Path(instance["state"])
        if state != CHECKOUT / ".workflow-kernel/runs" / instance["runId"]:
            fail("Gallery Kernel state binding mismatch.")
        evidence = run / "evidence"
        env = evidence / "compose.env"
        compose = ["docker", "compose", "--env-file", env, "--project-name", instance["project"],
                   "--project-directory", CHECKOUT, "-f", CHECKOUT / "examples/gallery/compose.yaml"]

        def compose_env(receipt):
            values = dict(GALLERY_RUN=str(run), GALLERY_PROJECT=instance["project"], GALLERY_RECEIPT=receipt)
            env.write_text("".join(key + "='" + value.replace("'", "\\'") + "'\n" for key, value in values.items()))
            for key in values:
                os.environ.pop(key, None)

        compose_env("not-building")

        def reconcile():
            descriptor = evidence / (uuid.uuid4().hex + "-cleanup.json")
            result = kernel("plan-reconcile", "--state-dir", state, "--run-id", instance["runId"], "--ttl-hours", "876000", "--output", descriptor, check=False)
            if result.returncode not in [0, 3]:
                fail("Kernel reconciliation failed: " + str(descriptor))
            paths = read(descriptor)
            return paths, read(pathlib.Path(paths["current_run_plan"])), read(pathlib.Path(paths["stale_sweep_plan"]))

        def inventory():
            _, _, sweep = reconcile()
            result = sweep["inventory"]
            result["resources"] = [r for r in result["resources"] if r["labels"].get("com.designmachines.depot.run-id") == instance["runId"]]
            path = evidence / (uuid.uuid4().hex + "-inventory.json")
            write(path, result)
            return path

        def create(argv):
            token = uuid.uuid4().hex
            argv_path = evidence / (token + "-argv.json"); write(argv_path, [str(x) for x in argv])
            deps = evidence / (token + "-deps.json"); write(deps, [])
            before = inventory()
            plan_path = evidence / (token + "-creation.json")
            kernel("plan-compose", "--state-dir", state, "--run-id", instance["runId"], "--node-id", instance["node"],
                   "--lifecycle", "run", "--cleanup-policy", "stop-remove", "--repository-project-name", instance["project"],
                   "--argv-json", argv_path, "--dependent-node-ids-json", deps, "--output", plan_path)
            plan = read(plan_path)
            if not plan["managed"] or plan["project_name"] != instance["project"]:
                fail("Kernel declined this exact gallery Docker boundary.")
            override = CHECKOUT / plan["compose_override"]
            if not override.is_relative_to(CHECKOUT / ".workflow-kernel"):
                fail("Unsafe Kernel override path.")
            override.parent.mkdir(parents=True, exist_ok=True); override.write_text(plan["compose_override_content"])
            os.environ.update(plan["environment"] or {})
            result = command(plan["argv"], check=False)
            result_path = evidence / (token + "-result.json")
            write(result_path, dict(schema_version=1, argv=plan["argv"], exit_code=result.returncode, stdout=result.stdout, stderr=result.stderr))
            after = inventory()
            registered = kernel("record-create", "--state-dir", state, "--plan", plan_path, "--result", result_path,
                                "--before-inventory", before, "--after-inventory", after, check=False)
            (evidence / (token + "-registered.json")).write_text(registered.stdout)
            if result.returncode or registered.returncode:
                fail("Docker creation failed; exact result/registration retained at " + str(result_path))
            kernel("validate-resource-registry", "--state-dir", state, "--run-id", instance["runId"], "--node-id", instance["node"])

        def stop():
            # Execute only the current run plan; never execute the stale sweep.
            for _ in range(12):
                paths, current, _ = reconcile()
                if not current["inventory"]["resources"]:
                    return
                plan_path = pathlib.Path(paths["current_run_plan"])
                outcomes = evidence / (uuid.uuid4().hex + "-outcomes.json"); write(outcomes, [])
                while True:
                    step_path = evidence / (uuid.uuid4().hex + "-step.json")
                    kernel("next-cleanup-step", "--state-dir", state, "--plan", plan_path, "--outcomes", outcomes, "--output", step_path)
                    step = read(step_path)
                    if step["complete"]:
                        break
                    _, fresh, _ = reconcile()
                    witness = evidence / (uuid.uuid4().hex + "-witness.json"); write(witness, fresh["inventory"])
                    run_state = read(state / "run-state.json")
                    statuses = evidence / (uuid.uuid4().hex + "-nodes.json")
                    write(statuses, dict(schema_version=1, run_id=instance["runId"], revision=run_state["revision"], updated_at=run_state["updated_at"], node_statuses={k:v["status"] for k,v in run_state["nodes"].items()}))
                    result_path = evidence / (uuid.uuid4().hex + "-removed.json")
                    kernel("execute-cleanup-step", "--state-dir", state, "--plan", plan_path, "--step-index", str(step["step_index"]),
                           "--inventory", witness, "--node-statuses", statuses, "--outcomes", outcomes, "--output", result_path)
                    rows = read(outcomes); rows.append(read(result_path)); write(outcomes, rows)
                result = kernel("record-cleanup", "--state-dir", state, "--plan", plan_path, "--outcomes", outcomes, check=False)
                (evidence / (uuid.uuid4().hex + "-cleanup-receipt.json")).write_text(result.stdout)
                if result.returncode not in [0, 3] or not read(outcomes):
                    fail("Cleanup blocked; preserve gallery state and inspect " + str(plan_path))
            fail("Cleanup did not converge; exact resources retained.")

        def status():
            identity, _ = source()
            result = dict(instance, health="stopped", url=None, currentSource=identity)
            mapped = command(compose + ["port", "app", "8080"], check=False)
            if mapped.returncode or not mapped.stdout.strip():
                return result
            address = mapped.stdout.strip()
            if not address.startswith("127.0.0.1:") or not address.split(":")[-1].isdigit():
                fail("Gallery must publish only to loopback.")
            kernel("validate-resource-registry", "--state-dir", state, "--run-id", instance["runId"], "--node-id", instance["node"])
            result["url"] = "http://" + address
            try:
                health = json.load(urllib.request.urlopen(result["url"] + "/healthz", timeout=3))
                result["builtSource"] = health["source"]
                result["health"] = "ready" if health["source"] == identity else "source-mismatch"
            except (OSError, ValueError):
                result["health"] = "unreachable"
            return result

        if args.action == "status":
            print(json.dumps(status())); return
        if args.action == "stop":
            stop(); print(json.dumps(dict(health="stopped", run=str(run), retained="verified assets, build cache and lifecycle evidence"))); return
        lock = read(CHECKOUT / "examples/gallery/assets.lock.json")
        for asset in lock["assets"]:
            relative = pathlib.PurePosixPath(asset["file"])
            if relative.is_absolute() or ".." in relative.parts:
                fail("Unsafe asset lock path.")
            target = run / "assets" / relative
            if target.exists():
                data = target.read_bytes()
                if target.resolve() != target or len(data) != asset["size"] or sha(data) != asset["sha256"]:
                    fail("Corrupt asset: " + str(target) + "; remove only this cache file, then retry assets.")
            else:
                data = urllib.request.urlopen(asset["url"], timeout=30).read()
                if len(data) != asset["size"] or sha(data) != asset["sha256"]:
                    fail("Downloaded asset does not match lock: " + asset["file"])
                target.parent.mkdir(parents=True, exist_ok=True)
                temporary = target.with_suffix(target.suffix + ".download")
                with temporary.open("xb") as out:
                    out.write(data)
                temporary.replace(target); target.chmod(0o444)
        if args.action == "assets":
            print(json.dumps(dict(assets="verified", revision=lock["producerCommit"], run=str(run)))); return
        if args.action == "rebuild":
            stop()
        elif args.action == "start" and command(compose + ["ps", "--all", "--quiet"], check=False).stdout.strip():
            fail("Gallery resources exist. Use status or rebuild.")
        identity, hashes = source()
        snapshot = run / "build/source"
        if snapshot.exists():
            shutil.rmtree(snapshot)
        snapshot.mkdir()
        for name, expected in hashes.items():
            data = (CHECKOUT / name).read_bytes()
            if sha(data) != expected:
                fail("Source changed while freezing build; retry.")
            target = snapshot / name; target.parent.mkdir(parents=True, exist_ok=True); target.write_bytes(data)
        for name in ["modules", "go-build"]:
            (run / "cache" / name).mkdir(exist_ok=True)
        receipt = base64.urlsafe_b64encode(json.dumps(identity, separators=(",", ":")).encode()).decode().rstrip("=")
        compose_env(receipt)
        if args.action == "verify":
            create(compose + ["run", "--rm", "--build", "builder", "sh", "-c", "templ generate && go test ./... && go vet ./..."])
        else:
            create(compose + ["run", "--rm", "--build", "builder"])
        changed = [n for n,h in hashes.items() if n.endswith("_templ.go") and sha((snapshot / n).read_bytes()) != h]
        if changed:
            fail("Generated files are stale. Run make generate in Docker and commit them: " + ", ".join(changed))
        if source()[0] != identity:
            fail("Checkout changed during build; retry rebuild.")
        write(evidence / "compiled-source.json", identity)
        if args.action == "verify":
            print(json.dumps(dict(verification="passed", source=identity))); return
        create(compose + ["up", "--detach", "--wait", "--wait-timeout", "30", "app"])
        result = status(); print(json.dumps(result))
        if result["health"] != "ready":
            fail("Gallery readiness failed; preserve evidence and use status/stop.")


if __name__ == "__main__":
    main()
