#!/usr/bin/env python3
import argparse, hashlib, json, math, os, platform, random, sys
from dataclasses import dataclass

EXPERIMENT = "WLM-LM-EXTERNAL-FUTURE-DATA-CAUSAL-STATE-SUFFICIENCY-R31-CONTRIB-CGPT"
PARENT = "ab1d7467edc021d70a96fcc3b749651e4bb46fea"
EQUAL = (582,581,581)
CANDIDATES = [
    (582,436,726),(582,581,581),(582,726,436),(727,291,726),(727,436,581),
    (727,581,436),(727,726,291),(872,291,581),(872,436,436),(872,581,291),
]
BUD_CODE = [0,436,582,727,872]
BUD_OTHER = [0,291,436,581,726]
RIDGE = 1.0
LATENT_DIM = 6
RANDOM_SEED = 31031

MANIFEST_SPECS = {
 "transfer": {
  "code": (41453,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45"),
  "structured": (14365,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932"),
  "technical_prose": (1454,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b"),
 },
 "third": {
  "code": (41453,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff"),
  "structured": (14365,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56"),
  "technical_prose": (1454,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017"),
 },
 "fourth": {
  "code": (41453,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1"),
  "structured": (14365,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933"),
  "technical_prose": (1454,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222"),
 },
 "fifth": {
  "code": (41453,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a"),
  "structured": (14365,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281"),
  "technical_prose": (1454,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49"),
 },
 "sixth": {
  "code": (41453,"ef50a59fbffad91d3ee1e03834dc8d295f5a847464ce8ead4a1366859a766f3a"),
  "structured": (14365,"5bc95ef524c23ec2df24d3d811058abd0c341921d2e1ed4ed29af2abed9b4974"),
  "technical_prose": (1454,"51085ab94eef524982bcb7858c573d1c5cb041decc27ccc137970111efec1592"),
 },
}

@dataclass
class Motif:
    key: tuple
    counts: list
    total: int
    best: int
    best_count: int
    consistency: float

@dataclass
class Arm:
    domain: str
    eval_bytes: bytes
    budgets: list
    bases: dict
    sels: dict
    features: dict
    hits: dict

@dataclass
class Manifest:
    name: str
    arms: list

def sha256(b): return hashlib.sha256(b).hexdigest()

def budget_pair(a,b):
    budget=15819; total=len(a)+len(b)
    if total<budget or total==0: return b'',b''
    na=budget*len(a)//total; nb=budget-na
    if na>len(a): na=len(a); nb=budget-na
    if nb>len(b): nb=len(b); na=budget-nb
    if na<0 or nb<0 or na>len(a) or nb>len(b) or na+nb!=budget: return b'',b''
    return a[:na],b[:nb]

def train_model(files):
    baseline=[[0]*256 for _ in range(256)]
    cand={}
    for data in files:
        for i in range(1,len(data)):
            baseline[data[i-1]][data[i]] += 1
        for i in range(4,len(data)):
            k=(data[i-4],data[i-3],data[i-2],data[i-1]); nxt=data[i]
            if k not in cand: cand[k]=[0,[0]*256]
            cand[k][0]+=1; cand[k][1][nxt]+=1
    rows=[]
    for k,(total,counts) in cand.items():
        if total<4: continue
        best=max(range(256), key=lambda i:(counts[i],-i))
        bc=counts[best]
        rows.append(Motif(k,counts,total,best,bc,bc/total))
    rows.sort(key=lambda r:(-r.best_count,-r.consistency,r.key))
    rows=rows[:256]
    return baseline,{r.key:r for r in rows}

def baseline_pred(base,prev):
    row=base[prev]
    return max(range(256), key=lambda i:(row[i],-i))

def hits(data,base,sel):
    h=0
    for i in range(1,len(data)):
        p=baseline_pred(base,data[i-1])
        if i>=4:
            k=(data[i-4],data[i-3],data[i-2],data[i-1])
            if k in sel: p=sel[k].best
        if p==data[i]: h+=1
    return h

def state_delta(lo_base,hi_base,lo_sel,hi_sel):
    bc=sum(baseline_pred(lo_base,i)!=baseline_pred(hi_base,i) for i in range(256))
    exits=0; best=0
    for k,lo in lo_sel.items():
        if k not in hi_sel: exits+=1
        elif lo.best!=hi_sel[k].best: best+=1
    entries=sum(k not in lo_sel for k in hi_sel)
    return bc,entries,exits,best

def argmax_row(row): return max(range(256),key=lambda i:(row[i],-i))
def top_margin(row):
    vals=sorted(row,reverse=True)
    return vals[0]-vals[1]

def build_manifest(name,srcs,open_eval=True):
    arms=[]
    for ti,(domain,data) in enumerate(srcs):
        others=[srcs[i][1] for i in range(3) if i!=ti]
        a,b=budget_pair(others[0],others[1])
        if len(a)+len(b)!=15819 or len(data)<1454: raise ValueError("budget/input invariant")
        ev=data[:582]; reserve=data[582:1454]
        budgets=BUD_CODE if domain=="code" else BUD_OTHER
        bases={}; sels={}; feats={}; hs={}
        for bud in budgets:
            train=[a,b]
            if bud>0: train.append(reserve[:bud])
            base,sel=train_model(train)
            if len(sel)!=256: raise ValueError(f"motif capacity drift {name}/{domain}/{bud}: {len(sel)}")
            bases[bud]=base; sels[bud]=sel
            if open_eval: hs[bud]=hits(ev,base,sel)
        for i in range(1,len(budgets)):
            lo,hi=budgets[i-1],budgets[i]; w=hi-lo
            lb,hb=bases[lo],bases[hi]; ls,ss=sels[lo],sels[hi]
            bc,en,ex,best=state_delta(lb,hb,ls,ss)
            margin=sum(top_margin(hb[p]) for p in range(256) if argmax_row(lb[p])!=argmax_row(hb[p]))
            active=0
            for k,m in ss.items():
                if k not in ls and m.best!=baseline_pred(hb,k[3]): active+=1
            for k,m in ls.items():
                if k not in ss and m.best!=baseline_pred(lb,k[3]): active+=1
            feats[hi]=[bc/w,margin/w,en/w,ex/w,active/w,best/w]
        arms.append(Arm(domain,ev,budgets,bases,sels,feats,hs))
    return Manifest(name,arms)

def actual_adv(man,alloc):
    return sum(man.arms[i].hits[alloc[i]]-man.arms[i].hits[EQUAL[i]] for i in range(3))

def marginal_examples(man,manifest_idx):
    out=[]
    for arm in man.arms:
        for j in range(1,len(arm.budgets)):
            lo,hi=arm.budgets[j-1],arm.budgets[j]
            out.append((manifest_idx,arm.features[hi],arm.hits[hi]-arm.hits[lo]))
    return out

def trajectory_vector(man,alloc):
    out=[]
    for ai,arm in enumerate(man.arms):
        for j in range(1,len(arm.budgets)):
            hi=arm.budgets[j]; s=(1 if hi<=alloc[ai] else 0)-(1 if hi<=EQUAL[ai] else 0)
            out.extend([s*v for v in arm.features[hi]])
    assert len(out)==72
    return out

def temporal_permute(v):
    out=[]
    for ai in range(3):
        blocks=[v[(ai*24+j*6):(ai*24+(j+1)*6)] for j in range(4)]
        for b in reversed(blocks): out.extend(b)
    return out

def solve_linear(A,b):
    n=len(b); M=[list(A[i])+[b[i]] for i in range(n)]
    for c in range(n):
        p=max(range(c,n),key=lambda r:abs(M[r][c]))
        if abs(M[p][c])<1e-12: return None
        M[c],M[p]=M[p],M[c]; q=M[c][c]
        for j in range(c,n+1): M[c][j]/=q
        for r in range(n):
            if r==c: continue
            f=M[r][c]
            if f==0: continue
            for j in range(c,n+1): M[r][j]-=f*M[c][j]
    return [M[i][n] for i in range(n)]

def ridge_fit(X,y,lam=RIDGE):
    if not X: return None
    d=len(X[0]); n=d+1
    A=[[0.0]*n for _ in range(n)]; b=[0.0]*n
    for x,t in zip(X,y):
        z=[1.0]+list(x)
        for i in range(n):
            b[i]+=z[i]*t
            for j in range(n): A[i][j]+=z[i]*z[j]
    for i in range(1,n): A[i][i]+=lam
    return solve_linear(A,b)

def ridge_pred(beta,x): return beta[0]+sum(beta[i+1]*x[i] for i in range(len(x)))

def control_predict(train_mans,held):
    ex=[]
    for idx,m in enumerate(train_mans): ex += marginal_examples(m,idx)
    X=[e[1] for e in ex]; y=[e[2] for e in ex]
    beta=ridge_fit(X,y)
    preds=[]
    for alloc in CANDIDATES:
        total=0.0
        for ai,arm in enumerate(held.arms):
            for j in range(1,len(arm.budgets)):
                hi=arm.budgets[j]
                if hi<=alloc[ai]: total+=ridge_pred(beta,arm.features[hi])
                if hi<=EQUAL[ai]: total-=ridge_pred(beta,arm.features[hi])
        preds.append(total)
    return preds

def dot(a,b): return sum(x*y for x,y in zip(a,b))
def norm(a): return math.sqrt(max(dot(a,a),0.0))
def normalize(a):
    n=norm(a)
    return [x/n for x in a] if n>1e-15 else None

def center_fit(X):
    d=len(X[0]); mean=[sum(row[j] for row in X)/len(X) for j in range(d)]
    return mean,[[row[j]-mean[j] for j in range(d)] for row in X]

def covariance(Xc):
    d=len(Xc[0]); C=[[0.0]*d for _ in range(d)]
    scale=1.0/max(1,len(Xc)-1)
    for x in Xc:
        for i in range(d):
            xi=x[i]
            if xi==0: continue
            for j in range(d): C[i][j]+=xi*x[j]*scale
    return C

def matvec(C,v): return [dot(row,v) for row in C]

def pca_fit(X,k=LATENT_DIM):
    mean,Xc=center_fit(X); C=covariance(Xc); d=len(mean); comps=[]
    for ci in range(k):
        v=[math.sin((j+1)*(ci+1)*0.731)+math.cos((j+3)*(ci+2)*0.417) for j in range(d)]
        for q in comps:
            a=dot(v,q); v=[v[j]-a*q[j] for j in range(d)]
        v=normalize(v)
        if v is None: v=[1.0 if j==ci%d else 0.0 for j in range(d)]
        for _ in range(256):
            w=matvec(C,v)
            for q in comps:
                a=dot(w,q); w=[w[j]-a*q[j] for j in range(d)]
            nw=normalize(w)
            if nw is None: break
            if norm([nw[j]-v[j] for j in range(d)])<1e-12: v=nw; break
            v=nw
        comps.append(v)
    return mean,comps

def pca_transform(x,mean,comps):
    z=[x[j]-mean[j] for j in range(len(x))]
    return [dot(z,c) for c in comps]

def random_projection(x):
    rnd=random.Random(RANDOM_SEED); d=len(x); out=[]
    for _ in range(LATENT_DIM):
        s=0.0
        for j in range(d): s += (1.0 if rnd.getrandbits(1) else -1.0)*x[j]
        out.append(s/math.sqrt(d))
    return out

def class_metrics(preds,actuals):
    pos=[i for i,a in enumerate(actuals) if a>0]; neg=[i for i,a in enumerate(actuals) if a<=0]
    if not pos or not neg: return {"valid":False}
    tp=sum(preds[i]>0 for i in pos); tn=sum(preds[i]<=0 for i in neg)
    correct=sum((preds[i]>0)==(actuals[i]>0) for i in range(len(actuals)))
    fp=sum(preds[i]>0 for i in neg)
    return {"valid":True,"balanced_accuracy":0.5*(tp/len(pos)+tn/len(neg)),"sign_accuracy":correct/len(actuals),"false_positive_rate":fp/len(neg),"false_positive_count":fp,"positive_count":len(pos),"nonpositive_count":len(neg)}

def fit_policy_model(train_mans,repr_kind,label_permute=False):
    X=[]; y=[]
    for m in train_mans:
        yy=[actual_adv(m,a) for a in CANDIDATES]
        if label_permute: yy=yy[1:]+yy[:1]
        for ai,a in enumerate(CANDIDATES):
            X.append(trajectory_vector(m,a)); y.append(yy[ai])
    if repr_kind=="raw":
        beta=ridge_fit(X,y); return (lambda x:ridge_pred(beta,x)),None
    if repr_kind=="pca":
        mean,comps=pca_fit(X,LATENT_DIM); Z=[pca_transform(x,mean,comps) for x in X]; beta=ridge_fit(Z,y)
        return (lambda x:ridge_pred(beta,pca_transform(x,mean,comps))),(mean,comps)
    if repr_kind=="random":
        Z=[random_projection(x) for x in X]; beta=ridge_fit(Z,y)
        return (lambda x:ridge_pred(beta,random_projection(x))),None
    if repr_kind=="policy":
        Xp=[]
        for m in train_mans:
            for a in CANDIDATES: Xp.append([(a[i]-EQUAL[i])/1744.0 for i in range(3)])
        beta=ridge_fit(Xp,y)
        return (lambda x:ridge_pred(beta,x)),None
    raise ValueError(repr_kind)

def eval_fold(train,held):
    actual=[actual_adv(held,a) for a in CANDIDATES]
    A=control_predict(train,held)
    raw,_=fit_policy_model(train,"raw"); B=[raw(trajectory_vector(held,a)) for a in CANDIDATES]
    pca,_=fit_policy_model(train,"pca"); C=[pca(trajectory_vector(held,a)) for a in CANDIDATES]
    randomp,_=fit_policy_model(train,"random"); R=[randomp(trajectory_vector(held,a)) for a in CANDIDATES]
    pol,_=fit_policy_model(train,"policy"); P=[pol([(a[i]-EQUAL[i])/1744.0 for i in range(3)]) for a in CANDIDATES]
    lp,_=fit_policy_model(train,"pca",label_permute=True); L=[lp(trajectory_vector(held,a)) for a in CANDIDATES]
    Bt=[raw(temporal_permute(trajectory_vector(held,a))) for a in CANDIDATES]
    Ct=[pca(temporal_permute(trajectory_vector(held,a))) for a in CANDIDATES]
    return {"actual":actual,"control":class_metrics(A,actual),"raw":class_metrics(B,actual),"latent":class_metrics(C,actual),"random_projection":class_metrics(R,actual),"policy_only":class_metrics(P,actual),"label_permutation":class_metrics(L,actual),"raw_temporal_permutation":class_metrics(Bt,actual),"latent_temporal_permutation":class_metrics(Ct,actual),"predictions":{"control":A,"raw":B,"latent":C,"random_projection":R,"policy_only":P,"label_permutation":L}}

def mean(vals): return sum(vals)/len(vals) if vals else float('nan')

def main():
    ap=argparse.ArgumentParser(); ap.add_argument("--root",required=True); args=ap.parse_args()
    manifests={}; identity=[]
    for name,spec in MANIFEST_SPECS.items():
        srcs=[]
        for domain in ("code","structured","technical_prose"):
            path=os.path.join(args.root,name,domain+".bin"); b=open(path,"rb").read(); n,h=spec[domain]
            ok=len(b)==n and sha256(b)==h; identity.append({"manifest":name,"domain":domain,"bytes":len(b),"sha256":sha256(b),"identity_ok":ok})
            if not ok: raise SystemExit(f"identity mismatch {name}/{domain}")
            srcs.append((domain,b))
        manifests[name]=build_manifest(name,srcs,True)
    existing=[manifests[n] for n in ("transfer","third","fourth","fifth")]
    folds=[]
    for i,h in enumerate(existing):
        tr=[m for j,m in enumerate(existing) if j!=i]
        f=eval_fold(tr,h); f["holdout_manifest"]=h.name; folds.append(f)
    sixth=eval_fold(existing,manifests["sixth"]); sixth["holdout_manifest"]="sixth"
    invalid=0
    for f in folds+[sixth]:
        for k in ("control","raw","latent","random_projection","policy_only","label_permutation","raw_temporal_permutation","latent_temporal_permutation"):
            if not f[k].get("valid",False): invalid+=1
    macro={}
    for k in ("control","raw","latent","random_projection","policy_only","label_permutation","raw_temporal_permutation","latent_temporal_permutation"):
        macro[k]={"balanced_accuracy":mean([f[k]["balanced_accuracy"] for f in folds if f[k].get("valid")]),"sign_accuracy":mean([f[k]["sign_accuracy"] for f in folds if f[k].get("valid")]),"false_positive_rate":mean([f[k]["false_positive_rate"] for f in folds if f[k].get("valid")])}
    deltas=[f["latent"]["balanced_accuracy"]-f["control"]["balanced_accuracy"] for f in folds]
    stage1_lat=macro["latent"]; stage1_ctl=macro["control"]
    valid=(invalid==0 and all(x["identity_ok"] for x in identity))
    stage1=(valid and stage1_lat["balanced_accuracy"]>=0.65 and stage1_lat["balanced_accuracy"]-stage1_ctl["balanced_accuracy"]>=0.10 and sum(d>0 for d in deltas)>=3 and stage1_lat["false_positive_rate"]<=stage1_ctl["false_positive_rate"]+1e-15)
    s6_lat,s6_ctl=sixth["latent"],sixth["control"]
    stage2=(valid and s6_lat["sign_accuracy"]>=0.80 and s6_lat["sign_accuracy"]-s6_ctl["sign_accuracy"]>=0.20 and s6_lat["false_positive_count"]<=s6_ctl["false_positive_count"])
    support=stage1 and stage2
    raw_delta=macro["raw"]["balanced_accuracy"]-stage1_ctl["balanced_accuracy"]
    lat_delta=stage1_lat["balanced_accuracy"]-stage1_ctl["balanced_accuracy"]
    mixed=(valid and not support and (stage1 or raw_delta>=0.05 or lat_delta>=0.05 or (lat_delta>0 and not stage2)))
    if not valid: cls="INVALID"
    elif support: cls="SCIENTIFIC_SUPPORTED"
    elif mixed: cls="MIXED"
    else: cls="SCIENTIFIC_NEGATIVE"
    result={"schema":"wingless.external-research-result.v1","project":"Wingless","experiment_id":EXPERIMENT,"contributor":"chatgpt-external","contribution_class":"falsification","classification":cls,"validity_pass":valid,"accepted_state_changed":False,"adoption_candidate_only":True,"rsi_success":False,"exact_parent_sha":PARENT,"frozen_parameters":{"ridge_lambda":RIDGE,"latent_dim":LATENT_DIM,"random_projection_seed":RANDOM_SEED,"candidate_count":10,"adaptation_budget":1744},"stage1":{"macro":macro,"latent_minus_control_balanced_accuracy":lat_delta,"raw_minus_control_balanced_accuracy":raw_delta,"positive_latent_delta_holdouts":sum(d>0 for d in deltas),"fold_deltas":deltas,"passes_supported_stage1":stage1,"folds":folds},"stage2":{"sixth_manifest":sixth,"latent_sign_accuracy_minus_control":s6_lat["sign_accuracy"]-s6_ctl["sign_accuracy"],"passes_supported_stage2":stage2},"identity":identity,"invalid_condition_count":invalid,"interpretation_code":("TRANSFERABLE_LATENT_STATE_SUPPORTED" if support else "SEEN_ONLY_OR_PARTIAL" if mixed else "OBSERVABLE_TRAJECTORY_NOT_SUFFICIENT_UNDER_FROZEN_ENCODER"),"runtime":{"python":sys.version.split()[0],"platform":platform.platform()}}
    print(json.dumps(result,sort_keys=True,separators=(",",":")))

if __name__=="__main__": main()
