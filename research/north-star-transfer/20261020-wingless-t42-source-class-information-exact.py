"""T42 exact source-free causal source-class recognition and forecast bounds.

No T41 case is reused and no original scientific case/DEV can be generated.
Two law-source classes: U uniform over all 256 Boolean 3-input rules; C
uniform over 768 syntactic G(a,b) OP H(b,c) compositions (16 gate atoms each,
OP XOR/AND/OR), projected to 256 laws preserving program multiplicity.

Exhaust all mask/signature histories for observed k=0..8, and calculate
1) Bayes-optimal source-class recognition / irreducible ambiguity
2) conditional mutual information between class and public rule evidence
3) exact expected proper Brier risk on truly unobserved truth-table entries
under five true source frequencies for correct-prior Bayes, fixed 50/50,
uniform unrestricted and evaluator-only class-known oracle.

All risk is integrated over actual true distribution; no sampled test cases,
trained new mechanism, private source class labels, or CPU advantage claim.
"""
from __future__ import annotations
import itertools,json,math,sys,time
from collections import defaultdict
from functools import lru_cache

N_BITS=8
OPS=("XOR","AND","OR")
PRIORS=(0.1,0.25,0.5,0.75,0.9)
ALL_RULES=tuple(range(256))
N_COMPOSITIONS=16*16*3
SOURCEFREE_ONLY=True

def composed(g,h,op):
    truth=0
    for pat in range(8):
        a,b,c=(pat&1),(pat>>1)&1,(pat>>2)&1
        gx=(g>>(a+2*b))&1
        hy=(h>>(b+2*c))&1
        output=gx^hy if op==0 else gx&hy if op==1 else gx|hy
        truth|=output<<pat
    return truth

PROGRAM_LAWS=tuple(composed(g,h,op) for g in range(16)
                   for h in range(16) for op in range(3))
assert len(PROGRAM_LAWS)==N_COMPOSITIONS
MULTIPLICITIES=tuple(PROGRAM_LAWS.count(law) for law in ALL_RULES)
C_PRIOR=tuple(m/N_COMPOSITIONS for m in MULTIPLICITIES)
U_PRIOR=(1/256,)*256
assert abs(sum(C_PRIOR)-1)<1e-12

def H(p):
    if p<=0 or p>=1:return 0.
    return -p*math.log2(p)-(1-p)*math.log2(1-p)

def brier(p_true,p_forecast):
    return p_true*(1-p_forecast)**2+(1-p_true)*p_forecast**2

@lru_cache(None)
def histories(mask):
    """For mask return each response group with C/U mass and one masses."""
    group={}
    for law in ALL_RULES:
        signature=law&mask
        value=group.get(signature)
        if value is None:
            value=[0.,0.,[0.]*8,[0.]*8]
            group[signature]=value
        cc=C_PRIOR[law];uu=U_PRIOR[law]
        value[0]+=cc;value[1]+=uu
        for target in range(N_BITS):
            if (law>>target)&1:
                value[2][target]+=cc
                value[3][target]+=uu
    return group

def assert_direct_program_oracle(mask,group):
    """Independent composition-space reconstruction for every signature."""
    prog={}
    for law in PROGRAM_LAWS:
        sig=law&mask
        entry=prog.get(sig)
        if entry is None:
            entry=[0,[0]*N_BITS];prog[sig]=entry
        entry[0]+=1
        for bit in range(N_BITS):
            if law&(1<<bit):entry[1][bit]+=1
    max_diff=0.
    for signature,values in group.items():
        direct=prog.get(signature,[0,[0]*N_BITS])
        dm=direct[0]/N_COMPOSITIONS
        max_diff=max(max_diff,abs(dm-values[0]))
        for q in range(8):
            max_diff=max(max_diff,abs(direct[1][q]/N_COMPOSITIONS-values[2][q]))
    if max_diff>1e-12:raise RuntimeError("T42_768_PROGRAM_INDEPENDENT_ORACLE_MISMATCH")
    return len(group),max_diff

def evaluate_k(k,pi):
    if not 0<=k<=8:raise ValueError("T42_OBSERVATION_BUDGET_INVALID")
    if not 0<pi<1:raise ValueError("T42_PRIOR_BOUND_INVALID")
    mask_count=math.comb(8,k)
    class_accuracy=class_error=conditional_entropy=0.
    mutual_info=0.
    recognition_TV=0.
    fixed_risk=true_risk=uniform_risk=oracle_known_class_risk=0.
    posterior_mass_check=0.
    observation_groups=0
    cases_where_C_possible=0
    # Exact integration over each equally likely public observation-mask set.
    for indices in itertools.combinations(range(8),k):
        mask=sum(1<<bit for bit in indices)
        groups=histories(mask)
        observation_groups+=len(groups)
        mask_mass=0.
        for signature,(massC,massU,oneC,oneU) in groups.items():
            event_C=pi*massC
            event_U=(1-pi)*massU
            event=event_C+event_U
            if event<=0:continue
            observation_groups and None
            mask_mass+=event
            if massC>0:cases_where_C_possible+=1
            cp=event_C/event
            class_accuracy+=max(event_C,event_U)/mask_count
            class_error+=min(event_C,event_U)/mask_count
            conditional_entropy+=event*H(cp)/mask_count
            for target in range(8):
                if mask&(1<<target):continue
                p_U=oneU[target]/massU
                p_C=(oneC[target]/massC) if massC>0 else 0.
                p_true=(pi*oneC[target]+(1-pi)*oneU[target])/event
                denominator=.5*massC+.5*massU
                p_fixed=(.5*oneC[target]+.5*oneU[target])/denominator
                weight=event/(mask_count*(8-k))
                fixed_risk+=weight*brier(p_true,p_fixed)
                true_risk+=weight*brier(p_true,p_true)
                uniform_risk+=weight*brier(p_true,.5)
                oracle_known_class_risk+=(
                    (event_C*p_C*(1-p_C))+
                    (event_U*p_U*(1-p_U)))/(mask_count*(8-k))
        if abs(mask_mass-1)>1e-11:
            raise RuntimeError("T42_PUBLIC_EVENT_MASS_NOT_ONE")
        posterior_mass_check+=mask_mass/mask_count
    mutual_info=H(pi)-conditional_entropy
    if class_accuracy+class_error<1-1e-11 or class_accuracy+class_error>1+1e-11:
        raise RuntimeError("T42_CLASSIFICATION_BAYES_IDENTITY_BROKEN")
    if abs(posterior_mass_check-1)>1e-11:
        raise RuntimeError("T42_MASK_EVENT_NORMALIZATION_BROKEN")
    if mutual_info<-1e-10 or mutual_info>H(pi)+1e-10:
        raise RuntimeError("T42_CLASS_MUTUAL_INFORMATION_OUTSIDE_RANGE")
    if k==0 and abs(mutual_info)>1e-11:
        raise RuntimeError("T42_ZERO_EVIDENCE_INFORMATION_NOT_ZERO")
    if class_error<=0:
        raise RuntimeError("T42_PERFECT_SOURCE_CLASS_IDENTIFICATION_NOT_POSSIBLE")
    if k==8:
        return {
            "observed_inputs":k,
            "observation_input_masks":mask_count,
            "public_observation_signature_groups":observation_groups,
            "source_class_Bayes_MAP_accuracy":class_accuracy,
            "source_class_irreducible_MAP_error":class_error,
            "posterior_class_entropy_bits":conditional_entropy,
            "class_source_mutual_information_bits":mutual_info,
            "max_possible_class_information_bits":H(pi),
            "class_never_perfectly_identifiable":class_error>0,
            "future_unseen_targets":0,
            "expected_Brier":None,
            "learning_class_prior_gain_over_fixed_50pct":None,
            "private_known_class_oracle_headroom_vs_true_prior_Bayes":None,
        }
    if true_risk>fixed_risk+1e-11 or true_risk>uniform_risk+1e-11:
        raise RuntimeError("T42_TRUE_DISTRIBUTION_BAYES_RISK_NOT_MINIMAL")
    if oracle_known_class_risk>true_risk+1e-11:
        raise RuntimeError("T42_CLASS_KNOWN_RISK_NOT_LOWER_BOUND")
    prior_gain=(fixed_risk-true_risk)/fixed_risk
    class_headroom=(true_risk-oracle_known_class_risk)/true_risk
    if abs(pi-.5)<1e-12 and abs(prior_gain)>1e-10:
        raise RuntimeError("T42_BALANCED_MIXTURE_NOT_BAYES_OPTIMAL")
    if k==0 and class_accuracy < max(pi,1-pi)-1e-12:
        raise RuntimeError("T42_BELOW_NO_EVIDENCE_MAJORITY_BASELINE")
    return {
        "observed_inputs":k,
        "observation_input_masks":mask_count,
        "public_observation_signature_groups":observation_groups,
        "source_class_Bayes_MAP_accuracy":class_accuracy,
        "source_class_irreducible_MAP_error":class_error,
        "posterior_class_entropy_bits":conditional_entropy,
        "class_source_mutual_information_bits":mutual_info,
        "max_possible_class_information_bits":H(pi),
        "class_never_perfectly_identifiable":class_error>0,
        "future_unseen_targets":8-k,
        "expected_Brier":{
           "true_class_frequency_sourceblind_Bayes_optimum":true_risk,
           "fixed_50pct_sourceblind_Bayes_mixture":fixed_risk,
           "uniform_unrestricted_Bayes_0p5":uniform_risk,
           "clairvoyant_TRUE_source_class_oracle":oracle_known_class_risk,
        },
        "learning_class_prior_gain_over_fixed_50pct":prior_gain,
        "private_known_class_oracle_headroom_vs_true_prior_Bayes":class_headroom,
        "frozen_5pct_class_frequency_gain_possible":prior_gain>=.05,
        "frozen_5pct_unattainable_private_class_oracle_headroom":class_headroom>=.05
    }

def run():
    start=time.process_time_ns()
    class_occurrences={}
    direct_checks=[]
    for k in range(9):
        patterns=(tuple(range(k)),tuple(sorted(((q+2)%8 for q in range(k)))))
        unique=set()
        for pat in patterns:
            if len(set(pat))!=k:continue
            mask=sum(1<<bit for bit in pat)
            if mask in unique:continue
            unique.add(mask)
            event_count,delta=assert_direct_program_oracle(mask,histories(mask))
            direct_checks.append({"k":k,"mask":mask,"signatures_checked":event_count,
                                  "program_vs_truth_table_max_absolute_mass_error":delta})
    summaries={}
    for pi in PRIORS:
        rows=[evaluate_k(k,pi) for k in range(9)]
        if any(rows[i+1]["class_source_mutual_information_bits"]+1e-10<
               rows[i]["class_source_mutual_information_bits"] for i in range(8)):
            raise RuntimeError("T42_INFORMATION_DECREASED_WITH_MORE_OBSERVATIONS")
        first75=next((r["observed_inputs"] for r in rows if r["source_class_Bayes_MAP_accuracy"]>=.75),None)
        first90=next((r["observed_inputs"] for r in rows if r["source_class_Bayes_MAP_accuracy"]>=.90),None)
        gainmax=max(r["learning_class_prior_gain_over_fixed_50pct"] for r in rows[:-1])
        class_occurrences[str(pi)]={"source_C_prior":pi,"source_U_prior":1-pi,
            "first_k_source_MAP_accuracy_at_least_75pct":first75,
            "first_k_source_MAP_accuracy_at_least_90pct":first90,
            "highest_possible_prior_frequency_gain_over_fixed_mixture":gainmax,
            "exact_observation_budget_rows":rows}
    balanced=class_occurrences["0.5"]["exact_observation_budget_rows"]
    if any(abs(r["learning_class_prior_gain_over_fixed_50pct"])>1e-10 for r in balanced[:-1]):
        raise RuntimeError("T42_BALANCED_TRUE_PRIOR_MORE_LEARNABLE_THAN_BAYES")
    if abs(sum(C_PRIOR)-1)>1e-12 or len(PROGRAM_LAWS)!=768:
        raise RuntimeError("T42_PROGRAM_PRIOR_CUSTODY_FAIL")
    return {"schema":"wingless.T42.sourcefree-exact-causal-structure-recognition-and-Brier-ceiling.v1",
        "status":"EXACT_INFORMATION_BOUND_PASS",
        "sourcefree_rule_programs":768,"distinct_composition_rules":sum(x>0 for x in MULTIPLICITIES),
        "uniform_unrestricted_rules":256,
        "no_scientific_DEV_created":True,
        "T16_T41_previous_case_reuse":False,
        "exact_source_prior_frequencies":list(PRIORS),
        "number_of_observation_budgets":9,
        "sourcefree_information_and_risk":class_occurrences,
        "independent_direct_768_program_vs_projected_256_rule_checks":direct_checks,
        "balanced_50pct_mixture_is_already_Bayes_optimal":True,
        "correct_model_source_class_never_fully_identifiable_even_after_eight_bits":True,
        "not_a_trained_forecasting_or_native_CPU_discovery":True,
        "information_audit_full_CPU_ns":time.process_time_ns()-start,
        "scientific_positive":False,"optimizer_admission":False}

if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--selftest":
        if len(PROGRAM_LAWS)!=768 or sum(MULTIPLICITIES)!=768:
            raise RuntimeError("T42_PROGRAM_SPACE_INCORRECT")
        for i in range(16):
            for j in range(16):
                for op in range(3):
                    if composed(i,j,op) not in ALL_RULES:
                        raise RuntimeError("T42_COMPOSITION_RULE_OUT_OF_RANGE")
        print("T42_EXACT_GRAMMAR_SELFTEST_PASS")
    elif len(sys.argv)==1 or (len(sys.argv)==2 and sys.argv[1]=="--audit"):
        print(json.dumps(run(),sort_keys=True))
    else:raise SystemExit("T42_NO_UNFROZEN_SOURCE_MODE")
