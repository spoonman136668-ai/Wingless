import json
import os
import sqlite3
import subprocess
import sys

SPECS = [
    {
        "name": "wing",
        "db": r"C:\ProgramData\CKBR\research-sidecar\state\research-sidecar.db",
        "source": r"C:\ProgramData\CKBR\codex\work\rs\source\Wingless",
        "roots": [r"C:\ProgramData\CKBR\codex\work"],
        "inspect": [
            ".wingless/qualification-request.json",
            ".github/workflows/research-qualify-windows.yml",
            "research/bootstrap/wingless-lm-north-star.txt",
        ],
    },
    {
        "name": "ygg",
        "db": r"C:\ProgramData\CKBR\research-sidecar-yggdrasil\state\recovered-r2-20260930_190855\research-sidecar.db",
        "source": r"C:\ProgramData\CKBR\research-sidecar-yggdrasil\source\Yggdrasil",
        "roots": [r"C:\ProgramData\CKBR\research-sidecar-yggdrasil"],
        "inspect": [
            ".yggdrasil/qualification-request.json",
            ".yggdrasil/isolated-run.json",
            ".github/workflows/plane-isolated-research.yml",
            "research/architecture/yggdrasil-north-star.ice",
        ],
    },
]

def run_git(repo, *args, check=True):
    p = subprocess.run(
        ["git", "-C", repo, *args],
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    if check and p.returncode != 0:
        raise RuntimeError(f"git failed repo={repo} args={args} exit={p.returncode}")
    return p.returncode, p.stdout.strip()

def decode(v):
    if isinstance(v, (bytes, bytearray)):
        return v.decode("utf-8")
    return v

def latest_row(db, table):
    cols = [r[1] for r in db.execute(f"pragma table_info({table})")]
    if not cols:
        return None
    row = db.execute(f"select * from {table} order by rowid desc limit 1").fetchone()
    if row is None:
        return None
    return {cols[i]: decode(row[i]) for i in range(len(cols))}

def parse_raw(row):
    if not row:
        return {}
    raw = row.get("raw") or row.get("notice_raw")
    if isinstance(raw, str):
        try:
            return json.loads(raw)
        except Exception:
            pass
    return {}

def repo_roots(roots):
    seen = set()
    for root in roots:
        if not os.path.isdir(root):
            continue
        for current, dirs, files in os.walk(root):
            if ".git" in dirs or ".git" in files:
                seen.add(current)
                if ".git" in dirs:
                    dirs.remove(".git")
    return sorted(seen)

def has_commit(repo, sha):
    rc, _ = run_git(repo, "cat-file", "-e", f"{sha}^{{commit}}", check=False)
    return rc == 0

def show_file(repo, sha, path):
    rc, out = run_git(repo, "show", f"{sha}:{path}", check=False)
    if rc == 0:
        print(f"INSPECT_BEGIN name={path}")
        print(out)
        print(f"INSPECT_END name={path}")

for spec in SPECS:
    print(f"=== {spec['name']} ===")
    if not os.path.isfile(spec["db"]):
        raise RuntimeError(f"DB_MISSING {spec['db']}")
    db = sqlite3.connect(f"file:{spec['db'].replace(os.sep, '/')}?mode=ro", uri=True, timeout=5)
    try:
        rows = {}
        for table in ("research_sidecar_results", "research_implementation_freezes", "research_cycles", "settings"):
            try:
                rows[table] = latest_row(db, table)
            except Exception as exc:
                print(f"TABLE_ERROR table={table} error={exc!r}")
                rows[table] = None
    finally:
        db.close()

    latest = parse_raw(rows.get("research_sidecar_results"))
    freeze = parse_raw(rows.get("research_implementation_freezes"))
    cycle = parse_raw(rows.get("research_cycles"))
    for label, obj in (("result", latest), ("freeze", freeze), ("cycle", cycle)):
        print(
            "ROW"
            f" label={label}"
            f" cycle_id={obj.get('cycle_id','')}"
            f" status={obj.get('status','')}"
            f" experiment_id={obj.get('experiment_id','')}"
            f" baseline_sha={obj.get('baseline_sha','')}"
            f" source_head={obj.get('source_head','')}"
            f" source_commit={obj.get('source_commit','')}"
        )

    targets = []
    for obj in (latest, freeze, cycle):
        for key in ("source_head", "source_commit", "baseline_sha"):
            value = str(obj.get(key, ""))
            if len(value) == 40 and all(ch in "0123456789abcdef" for ch in value.lower()) and value not in targets:
                targets.append(value)

    repos = []
    if os.path.isdir(spec["source"]):
        repos.append(spec["source"])
    for repo in repo_roots(spec["roots"]):
        if repo not in repos:
            repos.append(repo)
    print(f"REPO_COUNT name={spec['name']} count={len(repos)}")

    latest_sha = str(latest.get("source_head", ""))
    latest_repo = None
    for repo in repos:
        rc, head = run_git(repo, "rev-parse", "HEAD", check=False)
        if rc == 0:
            print(f"REPO name={spec['name']} path={repo} head={head}")
        for sha in targets:
            if has_commit(repo, sha):
                print(f"TARGET_FOUND name={spec['name']} sha={sha} repo={repo}")
                if sha == latest_sha and latest_repo is None:
                    latest_repo = repo

    if latest_sha and latest_repo:
        _, tree = run_git(latest_repo, "rev-parse", f"{latest_sha}^{{tree}}")
        _, logline = run_git(latest_repo, "show", "-s", "--format=%H|%ct|%s", latest_sha)
        print(f"LATEST_SOURCE name={spec['name']} sha={latest_sha} tree={tree} repo={latest_repo}")
        print(f"LATEST_LOG name={spec['name']} value={logline}")
        for path in spec["inspect"]:
            show_file(latest_repo, latest_sha, path)
    elif latest_sha:
        print(f"LATEST_SOURCE_NOT_FOUND name={spec['name']} sha={latest_sha}")

for p in (
    r"C:\ProgramData\CKBR\codex\bin\codex.exe",
    r"C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe",
):
    if os.path.isfile(p):
        q = subprocess.run([p, "--version"], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
        print(f"CODEX path={p} version={q.stdout.strip()}")

print("REPO_LOCAL_CUTOVER_PY_PROBE=PASS")
