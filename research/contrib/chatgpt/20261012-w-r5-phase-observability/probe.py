"""Read-only Wingless W-R5 signature mechanism. No dataset, training or launch path."""
from collections import defaultdict

def _context_rows(rows):
    if not isinstance(rows, (tuple, list)) or not rows:
        raise ValueError("context rows required")
    out = []
    for row in rows:
        if not isinstance(row, (tuple, list)) or len(row) != 4:
            raise ValueError("context requires exactly four bytes")
        if any(type(b) is not int or b < 0 or b > 255 for b in row):
            raise ValueError("invalid context byte")
        out.append(tuple(row))
    if len(set(out)) != len(out):
        raise ValueError("repeated contexts are not allowed for frozen exhaustive arm")
    return tuple(out)

def phase_signature(context_rows, arm, *, family_control_index=None):
    """Return label-free canonical signature; A2's family-dependent phase is supplied by the experimental control only."""
    rows = _context_rows(context_rows)
    if arm not in ("A0", "A1", "A2"):
        raise ValueError("unsupported arm")
    support = tuple(sorted(rows))
    if arm == "A0":
        return (support, None)
    if arm == "A1":
        return (support, rows[0])
    if type(family_control_index) is not int or family_control_index not in (0, 1, 2):
        raise ValueError("A2 requires frozen family rotation control 0..2")
    if len(rows) != 16:
        raise ValueError("A2 requires sixteen frozen contexts")
    phase = (5 * (family_control_index + 1)) % 16
    return (support, rows[phase])

def induce_classes(signatures):
    """Induce classes by equality alone; cannot observe entity/family annotations."""
    if not isinstance(signatures, (list, tuple)) or not signatures:
        raise ValueError("signatures required")
    grouped = defaultdict(list)
    for index, sig in enumerate(signatures):
        grouped[sig].append(index)
    groups = sorted(grouped.values(), key=lambda members: members[0])
    return tuple(tuple(g) for g in groups)
