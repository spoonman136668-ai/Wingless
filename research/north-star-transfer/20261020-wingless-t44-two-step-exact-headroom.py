"""T44 exact contingent two-intervention source-free mathematical preflight.

Compare optimized Bellman two-step causal intervention planning, exact
myopic-first/best-second, best precommitted pair, and fixed pair.
All rewards are reductions of Bayes Brier risk on two disjoint target
rule bits. The observed intervention bits never themselves earn reward.

No learned candidate admitted before the preregistered two-step headroom
threshold passes. No unfamiliar causal graph source / scientific DEV born.
"""
from __future__ import annotations
import functools,itertools,json,math,random,sys,time
from pathlib import Path

TRAIN_COUNT=64
VALID_COUNT=32
PROGRAMS=[]
for g in range(16):
    for h in range(16):
        for op in range(3):
            law=0
            for p in range(8):
                a,b,c=p&1,(p>>1)&1,(p>>2)&1
                x=(g>>(a+2*b))&1
                y=(h>>(b+2*c))&1
                z=x^y if op==0 else x&y if op==1 else x|y
                law|=z<<p
            PROGRAMS.append(law)
assert len(PROGRAMS)==768
multiplicity=[0]*256
for rule in PROGRAMS:multiplicity[rule]+=1
PRIOR=tuple(0.5/256+0.5*m/768 for m in multiplicity)

def fixture(split,index):
    if split not in ("TRAIN","VALID"):
        raise RuntimeError("T44_SOURCEFREE_ONLY")
    if not 0<=index<(TRAIN_COUNT if split=="TRAIN" else VALID_COUNT):
        raise RuntimeError("T44_UNFROZEN_SOURCEFREE_ID")
    rng=random.Random((44002026 if split=="TRAIN" else 44012026)+index*17011)
    regime=rng.randrange(2)
    hidden=rng.choice(PROGRAMS) if regime==0 else rng.randrange(256)
    order=tuple(rng.sample(range(8),8))
    return {"id":f"t44-public-sourcefree-{split.lower()}-{index:03d}",
            "observed":tuple((p,(hidden>>p)&1) for p in order[:2]),
            "targets":tuple(order[2:4]),"actions":tuple(order[4:]),
            "_hidden":hidden,"_source_regime":regime}

def push(mask,signature,pattern,bit):
    if mask&(1<<pattern):
        if ((signature>>pattern)&1)!=bit:
            raise RuntimeError("T44_INCONSISTENT_CAUSAL_OBSERVATION")
        return mask,signature
    return mask|(1<<pattern),signature|(bit<<pattern)
def source_state(public):
    mask=signature=0
    for p,bit in public:mask,signature=push(mask,signature,p,bit)
    return mask,signature

@functools.lru_cache(None)
def posterior(mask,signature):
    rows=tuple((rule,prior) for rule,prior in enumerate(PRIOR)
           if (rule&mask)==signature)
    denom=sum(p for r,p in rows)
    if denom<=0:raise RuntimeError("T44_POSSIBLE_CAUSAL_WORLD_EXCLUDED")
    return tuple((rule,p/denom) for rule,p in rows)
@functools.lru_cache(None)
def risk(mask,signature,targets):
    state=posterior(mask,signature)
    return sum((lambda p:p*(1-p))(
         sum(weight*((rule>>target)&1) for rule,weight in state))
         for target in targets)/len(targets)

@functools.lru_cache(None)
def transitions(mask,signature,action):
    state=posterior(mask,signature)
    result=[]
    for outcome in (0,1):
        prob=sum(weight for rule,weight in state
                 if ((rule>>action)&1)==outcome)
        if prob<=1e-15:continue
        nmask,nsig=push(mask,signature,action,outcome)
        result.append((prob,nmask,nsig,outcome))
    if abs(sum(t[0] for t in result)-1)>1e-10:
        raise RuntimeError("T44_BAYES_OUTCOME_MASS_DRIFT")
    return tuple(result)

def expected_single_loss(mask,signature,action,targets):
    return sum(prob*risk(m,s,targets) for prob,m,s,y
               in transitions(mask,signature,action))
@functools.lru_cache(None)
def one_step(mask,signature,actions,targets):
    choices=tuple((expected_single_loss(mask,signature,a,targets),a)
                  for a in actions)
    return min(choices)
@functools.lru_cache(None)
def exact_two_step(mask,signature,actions,targets):
    choices=[]
    for first in actions:
        other=tuple(a for a in actions if a!=first)
        posterior_expected=0.
        branches=[]
        for prob,m,s,outcome in transitions(mask,signature,first):
            next_loss,second=one_step(m,s,other,targets)
            posterior_expected+=prob*next_loss
            branches.append((outcome,second,prob,next_loss))
        choices.append((posterior_expected,first,tuple(branches)))
    return min(choices,key=lambda x:(x[0],actions.index(x[1])))

def exact_myopic_first(mask,signature,actions,targets):
    _,first=one_step(mask,signature,actions,targets)
    other=tuple(a for a in actions if a!=first)
    terminal=0.
    branches=[]
    for prob,m,s,outcome in transitions(mask,signature,first):
        next_loss,next_action=one_step(m,s,other,targets)
        terminal+=prob*next_loss
        branches.append((outcome,next_action,prob,next_loss))
    return terminal,first,tuple(branches)

def noncontingent_pair(mask,signature,actions,targets):
    records=[]
    for a,b in itertools.permutations(actions,2):
        term=0.
        for prob,m,s,y in transitions(mask,signature,a):
            term+=prob*expected_single_loss(m,s,b,targets)
        records.append((term,a,b))
    return min(records,key=lambda x:(x[0],actions.index(x[1]),actions.index(x[2])))

def independent_path_risk(mask,signature,first,branches,targets):
    # Independent direct summation across 256 hidden worlds, rather than
    # reusing cached one_step terminal expectation.
    state=posterior(mask,signature)
    total=0.
    for y,second,prob,conditional_loss in branches:
        if prob<=0:continue
        for v in (0,1):
            subset=[(law,weight) for law,weight in state
                if ((law>>first)&1)==y and ((law>>second)&1)==v]
            mass=sum(w for law,w in subset)
            if mass<=1e-15:continue
            target_risk=0.
            for t in targets:
                pt=sum(w*((law>>t)&1) for law,w in subset)/mass
                target_risk+=pt*(1-pt)/len(targets)
            total+=mass*target_risk
    return total

def report():
    begin=time.process_time_ns()
    cases=[fixture("VALID",i) for i in range(VALID_COUNT)]
    # No private rule or source class in any policy/DP oracle invocation.
    opt=myopic=nonadapt=fixed=before=0.
    cases_opt_strict=cases_myopic_diff=cases_contingent_second=0
    result=[]
    correct=abstained=incorrect=0
    for item in cases:
        mask,signature=source_state(item["observed"])
        targets=item["targets"];actions=item["actions"]
        first_risk=risk(mask,signature,targets)
        optimum,first,branches=exact_two_step(mask,signature,actions,targets)
        myopic_terminal,myopic_action,myopic_branches=exact_myopic_first(mask,signature,actions,targets)
        no_contingent,a_fixed,b_fixed=noncontingent_pair(mask,signature,actions,targets)
        dumb=sum(prob*expected_single_loss(m,s,actions[1],targets)
                 for prob,m,s,y in transitions(mask,signature,actions[0]))
        independently=independent_path_risk(mask,signature,first,branches,targets)
        if abs(independently-optimum)>1e-10:
            raise RuntimeError("T44_EXACT_DP_NOT_EQUAL_INDEPENDENT_WORLD_ENUMERATION")
        if optimum>myopic_terminal+1e-11 or optimum>no_contingent+1e-11:
            raise RuntimeError("T44_BELLMAN_NOT_OPTIMAL")
        if any(first_risk+1e-11<r for r in (optimum,myopic_terminal,no_contingent,dumb)):
            raise RuntimeError("T44_INTERVENTION_INCREASED_EXPECTED_BAYES_RISK")
        # Only exact logical unanimity is a hard certification.
        actual=item["_hidden"]
        remaining=[law for law,w in posterior(mask,signature)]
        if actual not in remaining:raise RuntimeError("T44_PRIVATE_SOURCE_TRUE_RULE_LOST")
        for target in targets:
            values={(law>>target)&1 for law in remaining}
            if len(values)==1:
                if next(iter(values))!=((actual>>target)&1):incorrect+=1
                else:correct+=1
            else:abstained+=1
        opt+=first_risk-optimum
        myopic+=first_risk-myopic_terminal
        nonadapt+=first_risk-no_contingent
        fixed+=first_risk-dumb
        before+=first_risk
        cases_opt_strict+=int(myopic_terminal-optimum>1e-10)
        cases_myopic_diff+=int(first!=myopic_action)
        branch_actions={z[1] for z in branches}
        cases_contingent_second+=int(len(branch_actions)>1)
        result.append({"id":item["id"],"optimal_first_action":first,
            "myopic_first_action":myopic_action,
            "optimal_first_expected_loss":round(optimum,12),
            "myopic_first_expected_loss":round(myopic_terminal,12),
            "best_precommitted_pair_expected_loss":round(no_contingent,12),
            "first_response_second_policy":[(a,b) for a,b,p,x in branches]})
    rel=(opt-myopic)/myopic if myopic>0 else 0.
    non_rel=(opt-nonadapt)/nonadapt if nonadapt>0 else 0.
    return {"schema":"wingless.T44.original-two-step-exact-adaptive-planning-headroom.v1",
        "status":"EXACT_MATHEMATICAL_PRE_FIRSTUSE_QUALIFICATION_PASS",
        "TRAIN_sourcefree_cases_reserved":TRAIN_COUNT,
        "validation_sourcefree_cases":VALID_COUNT,
        "actions_available_each_case":4,"two_adaptive_interventions_budget":2,
        "causal_targets_disjoint_from_all_actions":True,
        "sourcefree_prior_programs":768,"full_exact_Boolean_rules":256,
        "initial_public_causal_observations":2,
        "total_initial_expected_Brier_risk":before,
        "optimal_two_step_expected_gain":opt,
        "exact_myopic_first_two_step_expected_gain":myopic,
        "best_noncontingent_pair_expected_gain":nonadapt,
        "fixed_first_two_expected_gain":fixed,
        "relative_contingent_DP_advantage_vs_myopic":rel,
        "relative_DP_advantage_vs_best_precommitted_pair":non_rel,
        "cases_where_optimal_first_beats_myopic":cases_opt_strict,
        "cases_optimal_and_myopic_first_differ":cases_myopic_diff,
        "cases_with_branch_contingent_different_second_actions":cases_contingent_second,
        "independent_full_256_world_terminal_Brier_proofs":VALID_COUNT,
        "hard_correct_certifications":correct,"safe_abstentions":abstained,
        "hard_false_certifications":incorrect,
        "DP_two_percent_nonmyopic_headroom_pass":rel>=.02,
        "scientific_DEV_generated":0,"previous_consumed_scientific_cases_replayed":False,
        "sourcefree_case_rows":result,
        "original_whole_analytic_CPU_ns":time.process_time_ns()-begin,
        "exact_solver_cache_statistics":{
            "posterior":posterior.cache_info()._asdict(),
            "risk":risk.cache_info()._asdict(),
            "transitions":transitions.cache_info()._asdict(),
            "one_step":one_step.cache_info()._asdict(),
            "two_step":exact_two_step.cache_info()._asdict()},
        "learned_candidate_admitted":False,
        "optimizer_promotion":False}

if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--headroom":
        print(json.dumps(report(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--syntax":
        print("T44_TWO_STEP_SOURCEFREE_MATH_SYNTAX=PASS")
    else:raise SystemExit("T44_ONLY_FROZEN_SOURCEFREE_EXACT_HEADROOM_NO_DEV")
