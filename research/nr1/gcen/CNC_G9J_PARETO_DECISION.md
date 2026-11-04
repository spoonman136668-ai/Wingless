# CNC-G9J — G9I Pareto evidence and directional decision

**Evidence admission, not an independent new compression trial.**
Authenticated original Qwen3-30B-A3B BF16 weights from G9I source SHA
`f4b60b8de510d654211f5ecac483084e79997b62`, original Actions
run `41802853750` three jobs PASS, artifact `13308383049` ZIP SHA-256
`01688b9f0a999b73fdc4a8dd140f4eced590cfed41856940a5dfbe9b336643c8`,
internal JSON SHA-256
`b5c57f908c8f1d2050be5f77db70ef7990a964806b313f75b4d958ab80f221ad`.
Frozen independent seeds `107,163,229,271,347`; 16 original
128×256 expert windows and synthetic Gaussian activation inputs.

## Exact observed development comparison

| Arm | Relative output squared error | Stored bytes over 16 expert windows | BF16 compression |
|---|---:|---:|---:|
| INT4 control | 0.0123667407 | 295,296 | 3.55× |
| G9F INT2 | 0.1254078417 | 135,552 | 7.74× |
| G9I k=1 | 0.1128922180 | 152,000 | 6.90× |
| G9I k=2 | 0.1055264632 | 160,192 | 6.55× |
| G9I k=3 | 0.1000651796 | 168,384 | 6.23× |

**Judgment: NO_GO for direct full-model admission of uniformly adding
1/2/3 sparse corrections to every expert-weight row.** Authentic
numerical error falls, but even k=3 retains about eight times the INT4
control's squared output error, and its global 12-GiB modeled headroom
under an *unmeasured* 1-GiB workspace reserve is under 100 MiB.

## Next mechanism direction

Do not repeat the same original weight-window tuning. Select an
independent representation family that jointly addresses (1) robust
outlier/error treatment, (2) exact bit and metadata allocation, and
(3) direct compressed execution; verify its *real-model inference*
capability using an authorized native GPU runner before admitting
claims of reasoning preservation or 12-GiB full residency.

G9J adds NO original Qwen activations, new Qwen weight downloads,
GPU execution, protected heldout, scientific acceptance or protected
main changes. A model-level gate requires its own frozen prompts,
reference model, measured resource ceiling, and exact source/model
and GPU runtime custody; not authorized by this document.
