"""Pinned-source manifest admission. No network, downloads, or inference."""
import re
from custody import REQUIRED

MODEL = "Qwen/Qwen3-30B-A3B"
REVISION = "ad44e777bcd18fa416d9da3bd8f70d33ebb85d39"
HEX = re.compile(r"^[0-9a-f]{64}$")


def admit_source_manifest(document):
    """Validate source identity and independent digest provenance before byte verification.

    This is a *preflight*, not proof that the listed digests are authoritative:
    the caller must independently establish and retain the original source record.
    """
    if not isinstance(document, dict) or set(document) != {"model", "revision", "files", "source"}:
        raise ValueError("invalid provenance document shape")
    if document["model"] != MODEL or document["revision"] != REVISION:
        raise ValueError("unpinned checkpoint identity")
    source = document["source"]
    if not isinstance(source, dict) or set(source) != {"uri", "revision"}:
        raise ValueError("missing source custody")
    if not isinstance(source["uri"], str) or not source["uri"].startswith("https://huggingface.co/Qwen/Qwen3-30B-A3B/"):
        raise ValueError("unrecognized source URI")
    if source["revision"] != REVISION:
        raise ValueError("source revision mismatch")
    files = document["files"]
    if not isinstance(files, dict) or set(files) != set(REQUIRED):
        raise ValueError("incomplete checkpoint inventory")
    if any(not isinstance(value, str) or HEX.fullmatch(value) is None for value in files.values()):
        raise ValueError("invalid checkpoint SHA256")
    return dict(files)
