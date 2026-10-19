"""T21-B universal proof-carrying S3 relational summaries.

This is a finite-algebra engineering foundation, NOT a trained learner or
an independently admitted scientific source. It never imports prior DEV
source families, starts experiments, or grants optimizer admission.

A summary is a 3x3 matrix of 216-bit relation masks. A bit represents
one globally shared assignment of the three role permutations. Matrix
composition uses intersected SAME-assignment bits at every transition;
therefore repeated role identities cannot silently become independent.
"""
from __future__ import annotations

import hashlib
import itertools
import json
from dataclasses import dataclass

OPS = tuple(itertools.permutations(range(3)))
ASSIGNMENTS = tuple(itertools.product(range(6), repeat=3))
assert len(ASSIGNMENTS) == 216
FULL = (1 << len(ASSIGNMENTS)) - 1
GRAMMAR = "wingless.t21.s3.global-roles.relational-bitmatrix.v1"
CERTIFIER_BLOB = "1041db8ff76365eba5eb8f31cecd0e2172d397d8"
KERNEL_BLOB = "85dbd2673a94a6d72db29e57ee836619405f7f1e"

@dataclass(frozen=True)
class Summary:
    word: tuple
    relation: tuple
    identity: str

def normalize(word):
    word = tuple((r, inv) for r, inv in word)
    if len(word) > 16 or any(type(r) is not int or r not in range(3)
                           or type(inv) is not bool for r, inv in word):
        raise ValueError("T21_INVALID_WORD")
    return word

def direct(word, start, op_indices):
    value = start
    for role, inverse in word:
        op = OPS[op_indices[role]]
        value = op.index(value) if inverse else op[value]
    return value

def digest(word, relation):
    payload = json.dumps({
        "grammar": GRAMMAR,
        "certifier_blob": CERTIFIER_BLOB,
        "kernel_blob": KERNEL_BLOB,
        "word": word,
        "relation_hex": [[format(mask, "x") for mask in row]
                         for row in relation],
    }, sort_keys=True, separators=(",", ":")).encode("ascii")
    return hashlib.sha256(payload).hexdigest()

def _make(word, relation):
    relation = tuple(tuple(int(bitset) for bitset in row) for row in relation)
    if len(relation) != 3 or any(len(row) != 3 for row in relation):
        raise ValueError("T21_BAD_RELATION_SHAPE")
    if any(mask < 0 or mask & ~FULL for row in relation for mask in row):
        raise ValueError("T21_BAD_RELATION_MASK")
    if any(row[0] | row[1] | row[2] != FULL or
           (row[0] & row[1]) or (row[1] & row[2]) or
           (row[0] & row[2]) for row in relation):
        raise ValueError("T21_NOT_FUNCTIONAL_FOR_EVERY_BINDING")
    return Summary(word, relation, digest(word, relation))

def construct(word):
    """Construct an exact, universal input/output relation without source data."""
    word = normalize(word)
    relation = [[0, 0, 0] for _ in range(3)]
    for index, assignment in enumerate(ASSIGNMENTS):
        for start in range(3):
            out = direct(word, start, assignment)
            relation[start][out] |= 1 << index
    return _make(word, relation)

def compose(a, b):
    """Sound concatenation; bitwise AND enforces identical role bindings."""
    if not isinstance(a, Summary) or not isinstance(b, Summary):
        raise TypeError("T21_UNCERTIFIED_OPERAND")
    matrix = [[0, 0, 0] for _ in range(3)]
    for x in range(3):
        for z in range(3):
            for y in range(3):
                matrix[x][z] |= a.relation[x][y] & b.relation[y][z]
    return _make(a.word + b.word, matrix)

def mask_public(rows):
    """Evidence constrains one assignment bit per global role identity."""
    if len(rows) != 3:
        raise ValueError("T21_EXPECTED_THREE_PUBLIC_ROLES")
    mask = 0
    for index, assignment in enumerate(ASSIGNMENTS):
        valid = True
        for role, pairs in enumerate(rows):
            if not 1 <= len(pairs) <= 2:
                raise ValueError("T21_BAD_PUBLIC_CALIBRATION_COUNT")
            op = OPS[assignment[role]]
            if any(type(x) is not int or type(y) is not int or x not in range(3)
                   or y not in range(3) for x, y in pairs):
                raise ValueError("T21_BAD_PUBLIC_PAIR")
            if any(op[x] != y for x, y in pairs):
                valid = False
                break
        if valid:
            mask |= 1 << index
    return mask

def infer(summary, rows, start):
    if not isinstance(summary, Summary):
        raise TypeError("T21_NOT_SUMMARY")
    if start not in range(3):
        raise ValueError("T21_BAD_INPUT_STATE")
    admissible = mask_public(rows)
    if not admissible:
        return None, {"certified": False, "reason": "INCONSISTENT_PUBLIC"}
    answers = tuple(y for y in range(3)
                    if summary.relation[start][y] & admissible)
    if len(answers) != 1:
        return None, {"certified": False, "reason": "AMBIGUOUS_PUBLIC"}
    return answers[0], {"certified": True, "reason": "RELATION_COVERS_ALL_BINDINGS"}

def independently_verify(summary, public_cert):
    """Recompute each of the 648 cells through an independent exact oracle.

    A matching content digest alone does NOT constitute this verification.
    The unchanged T16B exact operator implementation is the second oracle.
    """
    if not isinstance(summary, Summary):
        raise TypeError("T21_NOT_SUMMARY")
    if summary.identity != digest(summary.word, summary.relation):
        raise ValueError("T21_RELATION_DIGEST_MISMATCH")
    if normalize(summary.word) != summary.word:
        raise ValueError("T21_WORD_DRIFT")
    for index, assignment in enumerate(ASSIGNMENTS):
        ops = tuple(OPS[k] for k in assignment)
        for start in range(3):
            current = start
            for role, inv in summary.word:
                current = public_cert.apply_op(1, ops[role], current, inv)
            for output in range(3):
                found = bool(summary.relation[start][output] & (1 << index))
                if found != (output == current):
                    raise ValueError("T21_UNIVERSAL_RELATION_NOT_EQUIVALENT")
    return {
        "full_assignment_comparisons": len(ASSIGNMENTS) * 3,
        "independent_oracle": "UNCHANGED_T16B_S3_APPLY_OP",
        "verified": True,
        "not_independent_science": True,
    }
