# Wingless UP-157B — late rehearsal cyclic position map

Status: preregistered scientific lexical-stability attribution experiment.

Scientific parent: sealed UP-156B c1a5ffb93fcb4664012eeddbf891b157b69b7fbf.

## Question

UP-156B showed that post-REPORT repair depends strongly on rehearsal order as well as dose. At a fixed 12-after / 3-before split, how does rotating the same 15 old-rehearsal items through every possible position change repair?

## Frozen design

Reconstruct the accepted 15-epoch common prefix.

For every terminal epoch:
- 20 non-REPORT new-family updates;
- exactly 3 old-rehearsal updates before REPORT;
- fixed four-example REPORT block;
- exactly 12 old-rehearsal updates after REPORT;
- all 15 old-rehearsal items used exactly once;
- same 0.08 learning rate and all other accepted settings.

Two new-family subjects are tested independently: Mia and Noah.

## Cyclic order arms

Use the exact accepted 15 old-rehearsal items for that epoch. Create 15 fixed cyclic rotations, offsets 0 through 14. Each rotation preserves coverage and total update count and changes only which old items occupy the three pre-REPORT and twelve post-REPORT positions.

No rotation is selected adaptively.

## Measurements

Per subject × rotation:
- ordered original rehearsal indices;
- three pre-REPORT indices;
- twelve post-REPORT indices;
- mean pre-REPORT recovery;
- mean REPORT damage;
- mean post-REPORT recovery;
- mean net epoch change;
- final mean old retention;
- final mean new accuracy.

## Interpretation

The exact 15-position response is the result. Variation across cyclic rotations demonstrates position/identity sensitivity at fixed dose and coverage. A flat response would reject that explanation.

## Bounds

No extra updates, no rehearsal coverage change, no memory change, no projector recomputation, no architecture change, no adaptive ordering, no result-informed retry, no live activation, no production authority.
