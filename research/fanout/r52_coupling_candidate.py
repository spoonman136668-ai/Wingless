#!/usr/bin/env python3
import argparse, hashlib, json, math, pathlib, urllib.request

INTERFACE="state_former(raw_input, history) -> fixed_size_state"
STATE_DIM=6
READOUT_CAPACITY=7
SOURCES=[
("transfer","code","https://raw.githubusercontent.com/python/cpython/ebf955df7a89ed0c7968f79faec1de49f61ed7cb/Lib/statistics.py",41453,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45"),
("transfer","structured","https://raw.githubusercontent.com/json-schema-org/JSON-Schema-Test-Suite/5b0ee1613e45fcc2bddac00e07c19cd49b00d8a8/tests/draft2020-12/uniqueItems.json",14365,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932"),
("transfer","technical_prose","https://raw.githubusercontent.com/golang/go/6f5c275ebdc454197fff5f1496521c8f81e20eef/doc/README.md",1454,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b"),
("third","code","https://raw.githubusercontent.com/torvalds/linux/a74306e2e676f9775457366fc047a660fbf02f26/kernel/sched/core.c",41453,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff"),
("third","structured","https://raw.githubusercontent.com/SchemaStore/schemastore/de76181a2ab215431d3e9314bc14f83cc3b01ad2/src/schemas/json/github-workflow.json",14365,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56"),
("third","technical_prose","https://raw.githubusercontent.com/rust-lang/book/1500248d8f230566e4ec9f27fcbb8fe9e2898ab1/src/ch01-03-hello-cargo.md",1454,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017"),
("fourth","code","https://raw.githubusercontent.com/kubernetes/kubernetes/e7967bb9b76d43a6388abb7a4e90d2f899ef97ff/pkg/kubelet/kubelet.go",41453,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1"),
("fourth","structured","https://raw.githubusercontent.com/microsoft/TypeScript/50d70a3f5f453a79a4323b263165da51f656a4e3/package-lock.json",14365,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933"),
("fourth","technical_prose","https://raw.githubusercontent.com/django/django/a461af8ce48762d7ec602260aaff81014ddccbcb/docs/intro/tutorial01.txt",1454,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222"),
("fifth","code","https://raw.githubusercontent.com/pytorch/pytorch/496340f06ef7bda2800522429ba3f4e3473a92fa/torch/_inductor/scheduler.py",41453,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a"),
("fifth","structured","https://raw.githubusercontent.com/microsoft/vscode/6fad7188e7dbf7db564e5e4a85960eb2bf69bdf4/package-lock.json",14365,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281"),
("fifth","technical_prose","https://raw.githubusercontent.com/numpy/numpy/2f1eca306857b641fb0fef0ab854a2d4137e1581/doc/source/user/absolute_beginners.rst",1454,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49"),
]

def fetch(url,n,expected):
    data=urllib.request.urlopen(urllib.request.Request(url,headers={"User-Agent":"wingless-r52-fanout"}),timeout=45).read(n)
    if len(data)!=n or hashlib.sha256(data).hexdigest()!=expected:
        raise SystemExit("R52_FANOUT_HISTORICAL_IDENTITY_MISMATCH")
    return data

def features(data):
    if not data:
        return (0.0,0.0)
    mean=sum(data)/(255.0*len(data))
    changes=sum(1 for a,b in zip(data,data[1:]) if a!=b)/max(1,len(data)-1)
    return (mean,changes)

def main():
    ap=argparse.ArgumentParser()
    ap.add_argument("--candidate",required=True);ap.add_argument("--candidate-config",required=True);ap.add_argument("--output",required=True)
    a=ap.parse_args()
    cfg=json.loads(pathlib.Path(a.candidate_config).read_text())
    alpha=[float(x) for x in cfg["alpha"]];groups=[int(x) for x in cfg["groups"]]
    if len(alpha)!=6 or len(groups)!=6 or min(groups)<0:
        raise SystemExit("R52_FANOUT_CONFIG_INVALID")
    grouped={}
    context_bytes=0
    by_manifest={}
    for manifest,domain,url,n,h in SOURCES:
        d=fetch(url,n,h);context_bytes+=len(d)
        by_manifest.setdefault(manifest,{})[domain]=features(d)
    order=("code","structured","technical_prose")
    for name,parts in by_manifest.items():
        raw=[]
        for domain in order: raw.extend(parts[domain])
        state=[raw[i]*alpha[i] for i in range(6)]
        gc=max(groups)+1
        sums=[0.0]*gc;counts=[0]*gc
        for i,v in enumerate(state): sums[groups[i]]+=v;counts[groups[i]]+=1
        grouped[name]=[sums[i]/counts[i] for i in range(gc)]
    names=sorted(grouped);dist=[]
    for i in range(len(names)):
        for j in range(i+1,len(names)):
            x,y=grouped[names[i]],grouped[names[j]]
            dist.append(math.sqrt(sum((aa-bb)**2 for aa,bb in zip(x,y))/len(x)))
    metric=sum(dist)/len(dist)
    params=1+len(set(groups))
    result={
      "candidate_id":a.candidate,"interface":INTERFACE,"fixed_state_dimension":STATE_DIM,"fixed_readout_capacity":READOUT_CAPACITY,
      "metric":{"name":"historical_grouped_state_separation","value":round(metric,12)},
      "resource_usage":{"parameters":params,"context_bytes":context_bytes,"model_calls":0},
      "historical_manifest_count":len(names),"sealed_outcome_exposure":False
    }
    pathlib.Path(a.output).write_text(json.dumps(result,indent=2,sort_keys=True)+"\n",encoding="utf-8")
if __name__=="__main__": main()
