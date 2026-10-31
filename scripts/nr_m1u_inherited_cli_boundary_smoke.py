#!/usr/bin/env python3
"""NR-M1U: execute the inherited Wingless CLI as a subprocess, no network/world.

This is an engineering boundary test, NOT an independent G26 source review.
"""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile


def require(ok, why):
    if not ok:
        raise AssertionError("NRM1U_CLI_BOUNDARY_FAIL: " + why)


def execute(binary, args, data=b""):
    return subprocess.run([str(binary), *args], input=data,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                          timeout=20, check=False)


def main():
    with tempfile.TemporaryDirectory(prefix="nr-m1u-wingless-") as tmp:
        binary = Path(tmp) / ("wingless.exe" if os.name == "nt" else "wingless")
        subprocess.run(["go", "build", "-o", str(binary), "./cmd/wingless"],
                       check=True, timeout=150)
        require(binary.is_file(), "native CLI binary not built")

        # The new CLI is only opt-in. The historical demo path still works.
        demo = execute(binary, ["demo"])
        require(demo.returncode == 0, "legacy offline demo path regressed")
        require(isinstance(json.loads(demo.stdout), dict),
                "legacy demo output not JSON object")

        # Historical unknown-command failure must remain intact.
        unknown = execute(binary, ["this-subcommand-does-not-exist"])
        require(unknown.returncode != 0, "unknown command unexpectedly accepted")
        require(b"usage:" in unknown.stderr and b"cognition" in unknown.stderr,
                "usage lacks old error boundary or new opt-in route")

        # No implicit activation: extra CLI arguments are rejected.
        extra_args = execute(binary, ["cognition", "unexpected"])
        require(extra_args.returncode != 0,
                "cognition with extra arguments unexpectedly accepted")

        safe = execute(binary, ["cognition"], b'{"op":"snapshot"}\n')
        require(safe.returncode == 0, "explicit offline snapshot failed")
        require("weights" in json.loads(safe.stdout),
                "offline native cognition state not serializable")

        for title, data in [
            ("unknown op", b'{"op":"execute"}\n'),
            ("unknown field", b'{"op":"snapshot","shell_command":"whoami"}\n'),
            ("trailing JSON", b'{"op":"snapshot"} {"op":"snapshot"}\n'),
            ("malformed JSON", b'{"op":"snapshot"\n'),
            ("oversized line", b'{"op":"snapshot","blob":"' + b"x" * 70000 + b'"}\n'),
            ("event budget exceeded", b'{"op":"observe","action":1,"prior":0,"outcome":1,"rate":0.1}\n' * 4097),
        ]:
            bad = execute(binary, ["cognition"], data)
            require(bad.returncode != 0,
                    title + " accepted; output=" + repr(bad.stdout[-300:]))

        # Aggregate cap: individually valid JSON lines must not hide truncation.
        line = b'{"op":"snapshot"}' + b" " * 60000 + b"\n"
        under = line * 17
        accepted = execute(binary, ["cognition"], under)
        require(accepted.returncode == 0 and
                len(accepted.stdout.splitlines()) == 17,
                "sub-limit legal records rejected")

        overflow = execute(binary, ["cognition"], line * 18)
        require(overflow.returncode != 0 and
                b"exceeds 1 MiB" in overflow.stderr,
                "multi-line aggregate overflow silently accepted")

        head = b'{"op":"snapshot"}'
        last_spaces = (1 << 20) - len(under) - len(head) - 1
        exact = under + head + b" " * last_spaces + b"\n"
        require(len(exact) == 1 << 20, "boundary test mis-sized")
        boundary = execute(binary, ["cognition"], exact)
        require(boundary.returncode == 0 and
                len(boundary.stdout.splitlines()) == 18,
                "exact 1 MiB stream rejected")

        one_extra = exact[:-1] + b" " + b"\n"
        require(len(one_extra) == (1 << 20) + 1,
                "over-boundary test mis-sized")
        rejected = execute(binary, ["cognition"], one_extra)
        require(rejected.returncode != 0 and
                b"exceeds 1 MiB" in rejected.stderr,
                "1 MiB plus one byte silently accepted")

    print(json.dumps({
        "gate": "NRM1U_INHERITED_CLI_SUBPROCESS_BOUNDARY",
        "result": "PASS",
        "tests": 14,
        "original_world_actions": 0,
        "network_or_controller_actions": 0,
        "independent_source_review": False,
        "review_status": "REVIEW_REQUIRED",
    }, sort_keys=True))


if __name__ == "__main__":
    main()
