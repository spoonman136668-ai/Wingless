"""Opt-in G9U authenticated single-expert CUDA execution entrypoint.

Never fetches the original Qwen checkpoint, chooses a GPU runner, or grants
execution authority. The operator must supply a staged original shard, pinned
manifest, isolated CUDA runtime, and exact clean repository checkout.
"""
import argparse
import json
from pathlib import Path
import subprocess

from g9u_gpu_operator_benchmark import run_authenticated_expert

SCHEMA="wingless.gcen.g9u.opt-in-source-bound-runner.v1"


def check_checkout(repository_root, source_sha):
    root=Path(repository_root).resolve(strict=True)
    if (not isinstance(source_sha,str) or len(source_sha)!=40
            or any(c not in "0123456789abcdef" for c in source_sha)):
        raise ValueError("invalid source SHA")
    def git(*args):
        result=subprocess.run(["git","-C",str(root),*args],
                              capture_output=True,text=True,check=True,timeout=15)
        return result.stdout.strip()
    if Path(git("rev-parse","--show-toplevel")).resolve()!=root:
        raise ValueError("repository root mismatch")
    if git("rev-parse","HEAD") != source_sha:
        raise ValueError("checked-out source SHA mismatch")
    if git("status","--porcelain","--untracked-files=all"):
        raise ValueError("dirty checkout forbidden")
    return root


def execute(*, repository_root, source_sha, shard_path, manifest_path, output):
    root=check_checkout(repository_root,source_sha)
    shard=Path(shard_path).resolve(strict=True)
    manifest=Path(manifest_path).resolve(strict=True)
    dest=Path(output)
    if dest.exists():
        raise FileExistsError("refusing to overwrite original G9U receipt")
    if not dest.parent.is_dir():
        raise ValueError("receipt parent directory does not exist")
    if dest.is_symlink():
        raise ValueError("receipt symlink forbidden")
    if shard == manifest or shard == dest or manifest == dest.resolve():
        raise ValueError("source and evidence paths overlap")
    receipt=run_authenticated_expert(shard_path=shard,manifest_path=manifest,
                                     source_sha=source_sha)
    if receipt.get("exact_source_sha") != source_sha or not receipt.get(
            "original_shard_authenticated_in_process"):
        raise RuntimeError("G9U authenticated source receipt incomplete")
    result={"schema":SCHEMA,"source_repo":"spoonman136668-ai/Wingless",
            "source_sha":source_sha,"checkout_root":str(root),
            "benchmark":receipt,
            "protected_heldout_used":False,"full_model_forward":False,
            "scientific_acceptance":"NOT_GRANTED"}
    # 'x' enforces no overwrite; never publish the actual original weights.
    with dest.open("x",encoding="utf-8") as f:
        json.dump(result,f,sort_keys=True,indent=2)
        f.write("\n")
    return result


def main():
    parser=argparse.ArgumentParser(description="Opt-in authenticated G9U CUDA operator run")
    parser.add_argument("--repository-root",type=Path,required=True)
    parser.add_argument("--source-sha",required=True)
    parser.add_argument("--shard-path",type=Path,required=True)
    parser.add_argument("--manifest",type=Path,required=True)
    parser.add_argument("--output",type=Path,required=True)
    a=parser.parse_args()
    execute(repository_root=a.repository_root,source_sha=a.source_sha,
            shard_path=a.shard_path,manifest_path=a.manifest,output=a.output)


if __name__=="__main__":
    main()
