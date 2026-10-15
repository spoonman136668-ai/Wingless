"""Wingless F49 TRAIN-only hidden-world identifiability probe, NOT an original.

Public evidence is identical for two hidden bijection completions. A public
policy never receives either hidden world; evaluator-only comparisons determine
whether one paid action can distinguish otherwise indistinguishable outcomes.
"""
import importlib.util,json,random,sys
from pathlib import Path
P=Path(__file__).resolve().parent.parent/"20261015-w-f48"/"source_pilot.py"
s=importlib.util.spec_from_file_location("w_f49_parent",P)
old=importlib.util.module_from_spec(s);sys.modules[s.name]=old;s.loader.exec_module(old)
p=old.old.f45.prior.prior
SEEDS=range(461100,461104)
FIT=range(461000,461004)
FAMILIES=(700,701,702,703)
RESERVED=range(462000,462016)
OLD_RESERVED=range(412000,412016)
BEAM=4
def source(seed,family,length):
 if seed not in SEEDS or family not in FAMILIES or length not in (3,5):
  raise ValueError("not a preregistered TRAIN-only diagnostic source")
 rng=random.Random(seed*17713+family*3037+length*521+49001)
 private=[];public=[]
 for role in range(4):
  perm=list(range(6));rng.shuffle(perm);private.append(tuple(perm))
  for x in rng.sample(range(6),3):public.append((role,x,perm[x]))
 rng.shuffle(public)
 templates={
 700:((0,1),(2,-1),(3,1),(1,-1),(0,-1)),
 701:((1,-1),(0,1),(2,1),(3,-1),(1,1)),
 702:((3,1),(1,-1),(0,1),(2,-1),(3,-1)),
 703:((2,-1),(3,1),(1,-1),(0,1),(2,1))}
 return tuple(private),tuple(public),rng.randrange(6),templates[family][:length]
def outcome(world,public,start,chain,action):
 if action is None:return p.predict(public,start,chain)
 role,subject=action
 value=world[role][subject]
 return p.predict((*public,(role,subject,value)),start,chain)
def run():
 assert not (set(SEEDS)&set(FIT) or set(SEEDS)&set(RESERVED)
             or set(SEEDS)&set(OLD_RESERVED))
 rows=[];all_beneficial=0;public_correct=0;oracle_correct=0;distinct=0
 for seed in SEEDS:
  for family in FAMILIES:
   for length in (3,5):
    _,public,start,chain=source(seed,family,length)
    worlds=p.worlds(public)
    assert len(worlds)==1296
    choices=sorted(p.candidate_actions(public))[:BEAM]
    assert len(choices)==BEAM
    public_action,cost=old.old.f45.one_step(public,start,chain,[0])
    evals=[]
    for world in (worlds[0],worlds[-1]):
     truth=p.follow(world,start,chain)
     prediction=outcome(world,public,start,chain,public_action)
     public_correct+=int(prediction==truth)
     outcomes={str(a):outcome(world,public,start,chain,a) for a in choices}
     beneficial=[i for i,a in enumerate(choices)
                 if outcomes[str(a)]==truth and p.predict(public,start,chain)!=truth]
     best=max((None,*choices),key=lambda a:(
       int(outcome(world,public,start,chain,a)==truth),int(a is None)))
     oracle_correct+=int(outcome(world,public,start,chain,best)==truth)
     evals.append(dict(truth=truth,public_policy_prediction=prediction,
       oracle_one_purchase_prediction=outcome(world,public,start,chain,best),
       beneficial_action_indices=beneficial,oracle_action=None if best is None else list(best)))
    # Evaluator-only: the two hypothetical worlds are hidden from any public selector.
    sets=[set(e["beneficial_action_indices"]) for e in evals]
    conflict=bool(sets[0] and sets[1] and sets[0].isdisjoint(sets[1]))
    distinct+=int(conflict)
    all_beneficial+=int(any(sets))
    rows.append(dict(seed=seed,family=family,length=length,public=list(map(list,public)),
      start=start,chain=list(map(list,chain)),public_one_step_action=
      None if public_action is None else list(public_action),counterfactual_worlds=evals,
      incompatible_oracle_actions=conflict,posterior_world_count=len(worlds),
      prereg_beam=[list(a) for a in choices]))
 n=len(rows)
 conflict_rate=distinct/n
 gain=(oracle_correct-public_correct)/(2*n)
 accepted=conflict_rate>=.25 and gain>=.10
 return dict(schema="wingless.f49.public-identifiability.trainonly.v1",
  stage="FRESH_PUBLIC_INDISCERNIBILITY_DIAGNOSTIC_NOT_PRIMARY_ORIGINAL",
  questions=n,counterfactual_world_evaluations=2*n,
  fit_seeds_not_used=[FIT.start,FIT.stop],
  diagnostic_seeds=[SEEDS.start,SEEDS.stop],families=list(FAMILIES),
  reserved_originals_not_consumed=[RESERVED.start,RESERVED.stop],
  old_reserved_originals_not_consumed=[OLD_RESERVED.start,OLD_RESERVED.stop],
  candidate_action_beam=BEAM,
  public_one_step_correct=public_correct,unfair_hindsight_oracle_correct=oracle_correct,
  distinct_beneficial_action_cases=distinct,any_beneficial_action_cases=all_beneficial,
  incompatible_hidden_optimal_action_fraction=conflict_rate,
  unfair_oracle_headroom_fraction=gain,criterion_met=accepted,
  verdict="TRAIN_ONLY_IDENTIFIABILITY_BARRIER_SUPPORTED_NOT_F3" if accepted else
    "TRAIN_ONLY_IDENTIFIABILITY_BARRIER_NOT_SHOWN",
  rows=rows,original_executed=False,independent_F3_custody=False,
  scientific_accepted=False,accepted_ref_mutation=False,
  native_576_byte_resident_qualified=False,paid_api_calls=0)
if __name__=="__main__":print(json.dumps(run(),sort_keys=True))
