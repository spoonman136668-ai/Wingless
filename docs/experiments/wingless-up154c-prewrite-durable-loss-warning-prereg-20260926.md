# Wingless UP-154C — prewrite durable-loss warning

Status: preregistered scientific durable-memory shadow-warning experiment.

Scientific parent: sealed UP-153C da9e256704ddafc9cb20c990f347e9bb59cf2181.

## Question

Can Wingless convert its exact hand+age next-slot forecast into an exact warning that the next write will destroy an original durable item, under mixed unique admissions and harmless rewrites?

## Frozen traffic

Use the same 16 durable originals and starting hands 0,4,8,12.

Run 24 writes per arm:
- writes at positions 4,8,12,16,20,24 rewrite the immediately preceding new key;
- all other positions are unique new keys.

Before write positions 3,7,11,15,19,23, issue fixed refresh queries to original keys:
2,6,10,14,1,5 respectively.
A missing refresh simply misses; no substitute is chosen.

## Frozen warning rule

Immediately before each write:
1. if the key already exists, predict no eviction and no durable-loss warning;
2. otherwise compute the next replacement slot from current hand+age state;
3. inspect only the current occupant of that predicted slot;
4. warn iff that occupant is one of original keys 0..15.

The warning does not alter the write.

## Measurements

Across 96 writes:
- TP/FP/FN/TN for original-durable loss;
- precision and recall;
- unique versus rewrite write counts;
- per-write predicted slot, actual slot, predicted loss, actual loss.

## Bounds

Shadow only. No write suppression, no refresh adaptation, no semantic priority, no capacity increase, no policy change, no future oracle, no result-informed retry, no live activation, or production authority.
