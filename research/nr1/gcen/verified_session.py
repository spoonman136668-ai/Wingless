"""Session-scoped verified checkpoint; metadata checks prevent accidental stale reuse.
Not an adversarial filesystem integrity guarantee. Rehash every new session.
"""
import hashlib
import json
from pathlib import Path
from custody import REQUIRED

class VerifiedCheckpoint:
    def __init__(self, root, manifest):
        self.root = Path(root).resolve()
        if not isinstance(manifest, dict) or set(manifest) != set(REQUIRED):
            raise ValueError("incomplete checkpoint manifest")
        self.identities = {}
        for name in REQUIRED:
            path = self.root / name
            if path.is_symlink() or not path.is_file():
                raise ValueError("missing or symlinked checkpoint file: " + name)
            expected = manifest[name]
            if not isinstance(expected, str) or len(expected) != 64 or any(c not in "0123456789abcdef" for c in expected):
                raise ValueError("invalid SHA256: " + name)
            digest = hashlib.sha256()
            with path.open("rb") as stream:
                for chunk in iter(lambda: stream.read(4 * 1024 * 1024), b""):
                    digest.update(chunk)
            if digest.hexdigest() != expected:
                raise ValueError("SHA256 mismatch: " + name)
            self.identities[name] = self._identity(path)

    @staticmethod
    def _identity(path):
        s = path.stat()
        return (s.st_dev, s.st_ino, s.st_size, s.st_mtime_ns, s.st_ctime_ns)

    def _checked(self, name):
        if name not in self.identities:
            raise ValueError("unknown checkpoint file")
        path = self.root / name
        if path.is_symlink() or not path.is_file() or self._identity(path) != self.identities[name]:
            raise ValueError("checkpoint file identity drift: " + name)
        return path

    def extract(self, tensor_name):
        index = json.loads(self._checked("model.safetensors.index.json").read_text(encoding="utf-8"))
        shard = index.get("weight_map", {}).get(tensor_name)
        if shard not in REQUIRED[:16]:
            raise ValueError("unmapped tensor")
        from safetensors import safe_open
        with safe_open(str(self._checked(shard)), framework="np", device="cpu") as f:
            return f.get_tensor(tensor_name)
