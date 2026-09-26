# Wingless UP-160C — durability geometry

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-159C 3f621ebede6fc7b4aa09bcb0eeb5b06bde295943.

## Question

Can the durable exception among C159's horizon-16 cases be explained by additional *present* replacement geometry, or do identical native states diverge because of later stream interactions?

## Frozen workload

Reuse UP-159C exactly:
- same disjoint 24-step mixed stream;
- same sparse refresh schedule;
- same four cohorts and four starting hands;
- same first harmful-write guided refresh;
- same exact recall cap 16.

No policy or action changes.

## Frozen present-state geometry

Immediately after the guided refresh and before the threatened write, record:
- current replacement hand;
- endangered slot;
- endangered slot offset from hand;
- all 16 two-bit ages, canonicalized by rotating the vector so current hand is index 0;
- unique-write horizon.

The canonical fingerprint contains no key identity, lexical class, or future schedule information.

## Primary test

Among cases with unique-write horizon exactly 16:
- count unique canonical geometry fingerprints;
- for each fingerprint, record whether outcomes include durable save, delayed loss, or both.

If one identical canonical present state produces both outcomes, present replacement geometry is insufficient to determine eventual durability under this stream.

## Bounds

Diagnostic only; counterfactual parent action unchanged. No future schedule input, semantic features, capacity growth, training, adaptive threshold, live activation, or production authority.
