#!/usr/bin/env python3
"""Read-only R35 evidence verifier. Output is advisory, never scientific acceptance.

The caller MUST independently retrieve the original GitHub Actions job log from
run 38936916822, job 117094661033, through an authenticated trusted transport.
Checking a supplied text file cannot authenticate that transport by itself.
No network, workflow dispatch, source modification, or experimental replay.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import sys

EXPERIMENT = "WLM-SI-RAW-REP-PHASE-ANCHOR-BOUNDARY44-SURFACE25-RECORD-CONTENT-NECESSITY-R35"
SCHEMA = "wingless.research-scientific-result.v1"
ORIGINAL_RUN = 38936916822
ORIGINAL_JOB = 117094661033


class EvidenceMismatch(ValueError):
    pass


def require(condition, message):
    if not condition:
        raise EvidenceMismatch(message)


def git_blob_sha(data):
    return hashlib.sha1(b"blob " + str(len(data)).encode("ascii") + b"\0" + data).hexdigest()


def git_at(commit, path):
    # Immutable SHA only; neither arguments nor output are passed to a shell.
    require(len(commit) == 40 and all(c in "0123456789abcdef" for c in commit),
            "non-immutable git commit")
    require(not path.startswith("/") and ".." not in Path(path).parts,
            "unsafe repository path")
    try:
        return subprocess.run(
            ["git", "show", commit + ":" + path], check=True,
            capture_output=True
        ).stdout
    except (OSError, subprocess.CalledProcessError) as exc:
        raise EvidenceMismatch("required original git object unavailable") from exc


def verify_original_sources(advisory, prereg_bytes):
    provenance = advisory["provenance"]
    require(advisory["experiment"] == EXPERIMENT, "advisory experiment mismatch")
    require(provenance["original_workflow_run_id"] == ORIGINAL_RUN, "run mismatch")
    require(provenance["original_workflow_job_id"] == ORIGINAL_JOB, "job mismatch")
    require(git_blob_sha(prereg_bytes) == provenance["frozen_prereg_git_blob_sha"],
            "prereg file blob drift")
    qualified_blob = git_at(provenance["qualified_source_commit"],
                            provenance["frozen_prereg_path"])
    require(qualified_blob == prereg_bytes, "qualified-commit prereg drift")
    frozen = json.loads(prereg_bytes.decode("utf-8-sig"))
    require(frozen["experiment"] == EXPERIMENT, "frozen experiment mismatch")
    require(frozen["exact_parent_sha"] == provenance["prereg_exact_parent_commit"],
            "frozen parent mismatch")
    rules = frozen["classification_rules"]
    require("record_content_necessary" in rules, "classification rule absent")
    state_bytes = git_at(provenance["original_evidence_record_commit"],
                         provenance["original_state_path"])
    state = json.loads(state_bytes.decode("utf-8-sig"))
    require(state["experiment"] == EXPERIMENT, "original state experiment mismatch")
    require(int(state["github_run_id"]) == ORIGINAL_RUN, "original state run mismatch")
    require(state["qualified_head_sha"] == provenance["qualified_source_commit"],
            "qualified source not anchored in original state")
    require(state["classification"] == "qualified-scientific-result",
            "original wrapper did not qualify")
    return frozen


def scientific_lines(log_text):
    decoder = json.JSONDecoder()
    for line in log_text.splitlines():
        if '"schema":"wingless.research-scientific-result.v1"' not in line:
            continue
        begin = line.find("{")
        if begin == -1:
            continue
        try:
            result, _ = decoder.raw_decode(line[begin:])
        except json.JSONDecodeError as exc:
            raise EvidenceMismatch("truncated or invalid JSON result in log") from exc
        if isinstance(result, dict) and result.get("schema") == SCHEMA \
                and result.get("experiment") == EXPERIMENT:
            yield result


def verify_result(result, advisory):
    observed = advisory["recovered_from_original_job_log"]
    require(result["schema"] == observed["result_schema"], "result schema drift")
    require(result["experiment"] == EXPERIMENT, "result experiment drift")
    require(result["classification"] == observed["scientific_classification"]
            == "record_content_necessary", "scientific classification mismatch")
    metrics = result["metrics"]
    expected = observed["metrics"]
    require(metrics == expected and set(metrics) == set(expected),
            "metric mismatch or unexpected metric")
    require(all(type(v) is int for v in metrics.values()), "invalid metric type")
    arms = result["arms"]
    require(len(arms) == 44, "missing or added arm")
    authentic = arms[0]
    require(authentic["source_sample_index"] == 43 and authentic["authentic"] is True,
            "authentic arm identity or ordering drift")
    require((authentic["pair_preserving_all4_match_count"],
             authentic["pair_preserving_partial_match_count"],
             authentic["pair_preserving_zero_match_count"]) == (0, 0, 127),
            "authentic scientific values drift")
    require(authentic["full5_alias_class_count"] == 23
            and authentic["full5_cross_entity_collision_pair_count"] == 0,
            "authentic validity controls drift")
    for i, arm in enumerate(arms[1:]):
        require(arm["source_sample_index"] == i and arm["authentic"] is False,
                f"duplicate {i} identity or ordering drift")
        require((arm["pair_preserving_all4_match_count"],
                 arm["pair_preserving_partial_match_count"],
                 arm["pair_preserving_zero_match_count"]) == (120, 0, 7),
                f"duplicate {i} scientific values drift")
        require(arm["full5_alias_class_count"] == 23
                and arm["full5_cross_entity_collision_pair_count"] == 0,
                f"duplicate {i} validity controls drift")
    require(metrics["authentic_sufficient_arm_count"] == 1
            and metrics["duplicate_sufficient_arm_count"] == 0
            and metrics["duplicate_unchanged_arm_count"] == 43,
            "classification counts drift")


def verify_log(log_text, advisory):
    markers = (
        "WINGLESS_AUTOMATED_QUALIFICATION_CLASSIFICATION: qualified-scientific-result",
        "WINGLESS_AUTOMATED_QUALIFICATION_PASS",
        "REPO_LOCAL_RESEARCH_CYCLE_PASS",
    )
    require(all(marker in log_text for marker in markers), "qualification log marker missing")
    require("experiment=" + EXPERIMENT in log_text, "completed experiment marker missing")
    results = list(scientific_lines(log_text))
    require(results, "no original scientific result in job log")
    for result in results:
        verify_result(result, advisory)
    # Contradictory copies must not be silently resolved by choosing one.
    canonical = json.dumps(results[0], sort_keys=True, separators=(",", ":"))
    require(all(json.dumps(x, sort_keys=True, separators=(",", ":")) == canonical
                for x in results[1:]), "conflicting result copies")
    return {"status": "ORIGINAL_LOG_CONTENT_MATCH_ADVISORY_ONLY",
            "original_run": ORIGINAL_RUN, "original_job": ORIGINAL_JOB,
            "scientific_classification": results[0]["classification"],
            "arms": len(results[0]["arms"]), "copies": len(results),
            "transport_authenticated_by_this_script": False,
            "independent_scientific_acceptance": False,
            "experiment_replayed": False}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--job-log", type=Path, required=True)
    parser.add_argument("--advisory", type=Path,
                        default=Path(__file__).with_name(
                            "r35-original-scientific-observation-advisory.json"))
    parser.add_argument("--prereg", type=Path, required=True)
    args = parser.parse_args()
    try:
        advisory = json.loads(args.advisory.read_text(encoding="utf-8-sig"))
        verify_original_sources(advisory, args.prereg.read_bytes())
        result = verify_log(args.job_log.read_text(encoding="utf-8-sig",
                                                  errors="replace"), advisory)
    except (EvidenceMismatch, OSError, KeyError, ValueError) as exc:
        print("R35_EVIDENCE_FAIL_CLOSED: " + str(exc), file=sys.stderr)
        return 1
    print(json.dumps(result, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
