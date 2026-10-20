"""T40 source-free exact learnability audit. No scientific DEV is generated."""
from __future__ import annotations
import itertools,json,math,random,sys,time
from collections import defaultdict
GATES=tuple(range(16))
OPS=("XOR","AND","OR")
PROGRAMS=tuple((g,h,op) for g in GATES for h in GATES for op in range(3))
PRIORS=("uniform_full_256","uniform_composition_768","nonuniform_composition_768")
OBS_K=(2,4,6)
META=64
TRAIN_SEED=40002026
STRIDE=9817
SOURCE_STRIDE=39011
MODELS=("fitted_compositional_factors","fitted_joint_768","fitted_full_256","fitted_independent_Beta","uniform_full_256","omniscient_source_oracle")
STRONG=("fitted_joint_768","fitted_full_256","fitted_independent_Beta","uniform_full_256")
def truth_table(g,h,op):
    result=0
    for q in range(8):
        a,b,c=q&1,(q>>1)&1,(q>>2)&1
        x=(g>>(a|(b<<1)))&1
        y=(h>>(b|(c<<1)))&1
        z=x^y if op==0 else x&y if op==1 else x|y
        result|=z<<q
    return result
LAW=tuple(truth_table(*z) for z in PROGRAMS)
assert len(LAW)==768
def normalized(values):
    total=sum(values)
    if total<=0:raise RuntimeError("T40_ZERO_MASS")
    r=tuple(x/total for x in values)
    if abs(sum(r)-1)>1e-12:raise RuntimeError("T40_NOT_PROBABILITY")
    return r
def true_prior(which):
    if which==PRIORS[0]:return (1/256,)*256,None
    if which==PRIORS[1]:w=(1/768,)*768
    elif which==PRIORS[2]:
        gw=[1+((7*i+3)%11) for i in range(16)]
        hw=[1+((5*i+1)%13) for i in range(16)]
        w=normalized([gw[g]*hw[h]*(5,2,3)[op] for g,h,op in PROGRAMS])
    else:raise RuntimeError("T40_UNKNOWN_PRIOR")
    rules=[0.]*256
    for weight,law in zip(w,LAW):rules[law]+=weight
    return tuple(rules),w
def pick_weighted(rng,weights):
    x=rng.random()*sum(weights)
    for i,w in enumerate(weights):
        x-=w
        if x<0:return i
    return len(weights)-1
def public_train(prior,hidden_weights):
    cases=[];atom_g=set();atom_h=set()
    for i in range(META):
        r=random.Random(TRAIN_SEED+PRIORS.index(prior)*SOURCE_STRIDE+i*STRIDE)
        if prior==PRIORS[0]:rule=r.randrange(256)
        elif i<16:
            g,h,op=i,(i+1)%16,r.randrange(3)
            atom_g.add(g);atom_h.add(h);rule=truth_table(g,h,op)
        elif i<32:
            h,g,op=i-16,(i-15)%16,r.randrange(3)
            atom_g.add(g);atom_h.add(h);rule=truth_table(g,h,op)
        else:
            z=pick_weighted(r,hidden_weights)
            g,h,_=PROGRAMS[z]
            atom_g.add(g);atom_h.add(h);rule=LAW[z]
        mask=sum(1<<q for q in r.sample(range(8),4))
        cases.append((mask,rule&mask))
    if prior!=PRIORS[0] and (len(atom_g)!=16 or len(atom_h)!=16):
        raise RuntimeError("T40_SOURCEFREE_ATOM_ROLE_COVERAGE_MISSING")
    return tuple(cases),{"G":len(atom_g),"H":len(atom_h)}
def train_joint(public,n,laws):
    w=(1/n,)*n
    for _ in range(3):
        acc=[0.05/n]*n
        for mask,signature in public:
            avail=[i for i in range(n) if laws[i]&mask==signature]
            den=sum(w[i] for i in avail)
            if den<=0:raise RuntimeError("T40_JOINT_EM_ZERO")
            for i in avail:acc[i]+=w[i]/den
        w=normalized(acc)
    return w
def train_factors(public):
    g=[1/16]*16;h=[1/16]*16;o=[1/3]*3
    for _ in range(3):
        gc=[0.05]*16;hc=[0.05]*16;oc=[0.05]*3
        for mask,signature in public:
            avail=[i for i,law in enumerate(LAW) if law&mask==signature]
            den=sum(g[PROGRAMS[i][0]]*h[PROGRAMS[i][1]]*o[PROGRAMS[i][2]] for i in avail)
            if den<=0:raise RuntimeError("T40_FACTOR_EM_ZERO")
            for i in avail:
                gg,hh,oo=PROGRAMS[i]
                p=g[gg]*h[hh]*o[oo]/den
                gc[gg]+=p;hc[hh]+=p;oc[oo]+=p
        g=list(normalized(gc));h=list(normalized(hc));o=list(normalized(oc))
    return normalized([g[gg]*h[hh]*o[oo] for gg,hh,oo in PROGRAMS])
def project_programs(weights):
    result=[0.]*256
    for law,w in zip(LAW,weights):result[law]+=w
    return tuple(result)
def independent_beta(public):
    c=[[2.,2.] for _ in range(8)]
    for mask,signature in public:
        for q in range(8):
            if mask&(1<<q):c[q][(signature>>q)&1]+=1
    return tuple(row[1]/sum(row) for row in c)
def models_from_public(public):
    return {"fitted_compositional_factors":project_programs(train_factors(public)),
            "fitted_joint_768":project_programs(train_joint(public,768,LAW)),
            "fitted_full_256":train_joint(public,256,tuple(range(256))),
            "fitted_independent_Beta":independent_beta(public),
            "uniform_full_256":(1/256,)*256}
def groups(mask):
    buckets=defaultdict(list)
    for rule in range(256):buckets[rule&mask].append(rule)
    return buckets
def binary_entropy(p):
    if p<=0 or p>=1:return 0.
    return -p*math.log2(p)-(1-p)*math.log2(1-p)
def oracle_independent_check(law_prior,program_prior,k):
    mask=(1<<k)-1;q=k
    allowed=[law for law in range(256) if (law&mask)==0]
    den=sum(law_prior[law] for law in allowed)
    predicted=sum(law_prior[law] for law in allowed if (law>>q)&1)/den
    if program_prior is None:
        direct_n=sum((1/256) for rule in range(256) if (rule&mask)==0 and (rule>>q)&1)
        direct_d=sum((1/256) for rule in range(256) if (rule&mask)==0)
    else:
        direct_n=sum(w for rule,w in zip(LAW,program_prior) if (rule&mask)==0 and (rule>>q)&1)
        direct_d=sum(w for rule,w in zip(LAW,program_prior) if (rule&mask)==0)
    delta=abs(predicted-direct_n/direct_d)
    if delta>1e-12:raise RuntimeError("T40_ORACLE_PROGRAM_VS_RULE_INCONSISTENT")
    return delta
def scenario(prior,law_prior,program_prior,trained,k):
    risk={model:0. for model in MODELS}
    info=0.;prior_H=0.;posterior_H=0.
    positive_events=0;all_events=0
    marginal=[sum(w for law,w in enumerate(law_prior) if (law>>q)&1) for q in range(8)]
    full_models={name:weights for name,weights in trained.items() if name!="fitted_independent_Beta"}
    for chosen in itertools.combinations(range(8),k):
        mask=sum(1<<q for q in chosen)
        unseen=[q for q in range(8) if not mask&(1<<q)]
        grouped=groups(mask)
        for signature,laws in grouped.items():
            mass=sum(law_prior[law] for law in laws)
            if mass<=0:continue
            positive_events+=1
            p_true={q:sum(law_prior[law] for law in laws if (law>>q)&1)/mass for q in unseen}
            trained_p={}
            for name,w in full_models.items():
                den=sum(w[law] for law in laws)
                if den<=0:raise RuntimeError("T40_TRAINED_CONTROL_HAS_ZERO_SUPPORT")
                trained_p[name]={q:sum(w[law] for law in laws if (law>>q)&1)/den for q in unseen}
            for q in unseen:
                p=p_true[q]
                ph=binary_entropy(marginal[q]);h=binary_entropy(p)
                prior_H+=mass*ph
                posterior_H+=mass*h
                info+=mass*(ph-h)
                for model in MODELS:
                    guess=(p if model=="omniscient_source_oracle" else
                           trained["fitted_independent_Beta"][q] if model=="fitted_independent_Beta" else trained_p[model][q])
                    risk[model]+=mass*(p*(1-guess)**2+(1-p)*guess**2)
        all_events+=len(grouped)
    divisor=math.comb(8,k)*(8-k)
    risk={name:value/divisor for name,value in risk.items()}
    oracle=risk["omniscient_source_oracle"]
    if any(oracle>v+1e-10 for v in risk.values()):
        raise RuntimeError("T40_ORACLE_BAYES_RISK_NOT_OPTIMAL")
    best_loss,best_name=min((risk[m],m) for m in STRONG)
    headroom=max(0.,(best_loss-oracle)/best_loss) if best_loss else 0.
    if prior==PRIORS[0] and (abs(oracle-.25)>1e-12 or abs(info)>1e-12):
        raise RuntimeError("T40_UNRESTRICTED_RULE_INFORMATION_LEAKAGE")
    return {"observed_patterns":k,"unobserved_patterns_per_history":8-k,
            "all_observation_masks":math.comb(8,k),
            "positive_probability_history_events":positive_events,
            "all_truth_table_signature_groups":all_events,
            "mean_mutual_information_bits_per_unseen_target":info/divisor,
            "prior_entropy_bits_per_unseen_target":prior_H/divisor,
            "posterior_entropy_bits_per_unseen_target":posterior_H/divisor,
            "expected_Brier":risk,"strongest_publicly_fitted_competitor":best_name,
            "strongest_publicly_fitted_Brier":best_loss,
            "oracle_Brier_lower_bound":oracle,
            "maximum_possible_relative_gain_beyond_best_publicly_trained_control":headroom,
            "theoretical_at_least_5pct_headroom":headroom>=.05,
            "logical_certified_answers_for_unobserved_input_with_256_rule_consistency":0,
            "independent_oracle_program_vs_rule_difference":oracle_independent_check(law_prior,program_prior,k)}
def audit():
    start=time.process_time_ns()
    multiplicity=[0]*256
    for rule in LAW:multiplicity[rule]+=1
    all_out={}
    for label in PRIORS:
        private_prior,program_prior=true_prior(label)
        public,coverage=public_train(label,program_prior)
        trained=models_from_public(public)
        all_out[label]={"TRAIN_public_cases":len(public),"TRAIN_public_observations_per_case":4,
                "private_fixture_coverage_for_auditor_only":coverage,
                "analytically_exact_scenarios":[scenario(label,private_prior,program_prior,trained,k) for k in OBS_K]}
    uniform=all_out[PRIORS[0]]["analytically_exact_scenarios"]
    if any(abs(x["oracle_Brier_lower_bound"]-.25)>1e-12 for x in uniform):
        raise RuntimeError("T40_UNRESTRICTED_BAYES_BRIER_NOT_QUARTER")
    return {"schema":"wingless.T40.exact-sourcefree-causal-learnability-and-headroom.v1",
        "status":"PASS_EXACT_INFORMATION_MATHEMATICS",
        "all_16_gate_atoms_per_role":True,
        "full_2input_gate_compositions_with_three_operators":len(PROGRAMS),
        "distinct_realizable_3input_rule_tables":len(set(LAW)),
        "composition_program_truth_table_equivalence_multiplicity":{
            "min":min(z for z in multiplicity if z),"max":max(multiplicity),
            "programs_with_nonunique_rule_identity":sum(z for z in multiplicity if z>1)},
        "sourcefree_scenarios":all_out,
        "analytically_evaluated_prior_by_observation_regimes":len(PRIORS)*len(OBS_K),
        "prior_values_used_by_models":False,
        "new_C2_D2_scientific_sources_generated":False,
        "T38_T39_consumed_cases_replayed":False,
        "scientific_DEV_generated_or_opened":0,
        "complete_math_audit_CPU_ns":time.process_time_ns()-start,
        "information_audit_not_native_competitor_regression":True,
        "not_prospective_learning_positive":True,
        "optimizer_promoted":False}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--selftest":
        assert len(set(PROGRAMS))==768 and len(PROGRAMS)==768
        assert all(truth_table(*p)==t for p,t in zip(PROGRAMS,LAW))
        print("T40_SOURCEFREE_SYNTACTIC_SELFTEST_PASS")
    elif len(sys.argv)==1 or (len(sys.argv)==2 and sys.argv[1]=="--audit"):
        print(json.dumps(audit(),sort_keys=True))
    else:raise SystemExit("T40_NO_SCIENTIFIC_DEV_INVOCATION")
