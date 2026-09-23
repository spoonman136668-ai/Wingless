# Wingless UP-6: learned observation and query interface

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-5 Windows qualification sealed at `1c61d679b4928e6e4f92fa63236a7872bcce4a15`.

## Question

Can the hybrid unitary-memory substrate operate without exhaustive runtime matching against all 256 possible memory tables and without hard-coded relational arithmetic?

UP-6 replaces both of those UP-5 conveniences with learned shared readout heads.

## Runtime architecture

The memory still uses sixteen complex dimensions for four entities with four values each.

Between commit boundaries:

1. canonical memory is perturbed;
2. it is transported through unitary or matched non-unitary propagation;
3. the known transport is inverted back into the canonical observation frame;
4. a learned shared value decoder reads each entity;
5. the explicit irreversible overwrite/commit boundary updates one entity and re-encodes canonical memory.

No 256-state prototype lookup occurs at runtime.

At the final query:

1. the learned shared value decoder produces four value distributions;
2. a learned shared relation head receives the two queried distributions;
3. it predicts one of four relation classes.

The relation head does not receive entity identity.

## Learned heads

### Shared value decoder

A four-class linear-softmax head maps one entity's four normalized magnitude features to one of four values.

It is trained on deterministic noisy local examples, not on the 48 full memory-program scenarios.

### Shared relation head

A four-class linear-softmax head maps the 16-dimensional outer product of two value distributions to one of four modular relation classes.

It is trained on the sixteen value-pair combinations. Because entity identity is absent, the same head must transfer to ordered entity pairs reserved for the integration test.

## Held-out integration workload

There are 48 deterministic mutable-memory scenarios.

Each scenario has:

- 16 writes/overwrites;
- transport gaps of 32, 128, 512, or 1024 blocks;
- bounded complex perturbation amplitude 0.10 before every transport segment;
- full-table decoding through the learned value head at every commit;
- a final query using only held-out ordered entity pairs.

Held-out query pairs:

- 0 -> 3;
- 3 -> 0;
- 2 -> 0;
- 3 -> 1.

The value and relation heads never train on the full scenario tables or write programs.

## Transport inversion

The unitary path is inverted with the exact reverse unitary program.

The matched non-unitary path is inverted with the exact algebraic inverse of each pair mixer:

```
[a']   [1 -g] [a]
[b'] = [g  1] [b]

M^-1 = 1/(1+g^2) * [1 g; -g 1]
```

This gives the control its strongest reasonable observation path rather than forcing it to decode in a distorted frame.

## Pass conditions

Learned heads:

- shared value decoder training accuracy = 1.0;
- shared relation head training accuracy = 1.0.

Unitary integration:

- all intermediate commit decodes correct;
- exact final-table accuracy = 1.0;
- relational-query accuracy = 1.0;
- maximum forward/inverse round-trip error <= 1e-10;
- maximum forward norm drift <= 1e-12;
- minimum value-classification margin > 0.5;
- minimum relation-classification margin > 0.5.

Matched control metrics must remain finite but are not required to fail.

Complete output must be deterministic and the existing Wingless regression must remain green after advisory ICE rebuild.

## Interpretation boundary

A pass would remove exhaustive full-state prototype lookup and hard-coded relation arithmetic from the UP-5 runtime path. The memory would use learned reusable observation/query heads around an explicit irreversible write boundary while transport remains reversible/unitary.

The inverse transport is still known rather than learned. UP-6 therefore does not yet establish end-to-end learned reasoning.

A positive UP-6 permits UP-7: learn a depth-agnostic observation map directly from transported states, then test unseen transport depths and program lengths without explicitly applying the known inverse before readout.

## Authority boundary

UP-6 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.
