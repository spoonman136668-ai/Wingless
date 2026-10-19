"""T21-C TRAIN-only preflight: no DEV build, source first-use or private truth."""
from __future__ import annotations
import importlib.util
import json
import sys
from pathlib import Path

ROOT=Path(__file__).resolve().parent

def load(name,file):
    spec=importlib.util.spec_from_file_location(name,ROOT/file)
    m=importlib.util.module_from_spec(spec)
    sys.modules[name]=m
    spec.loader.exec_module(m)
    return m

def main():
    model=load("t21c_model_for_train_only","20261019-wingless-t21c-learned-relations.py")
    basis,cert,source=model.deps()
    if len(source.TRAIN)!=40 or len(source.DEV)!=64 or set(source.TRAIN)&set(source.DEV):
        raise RuntimeError("T21C_SOURCE_IDENTITY_FAILURE")
    public=[source.public_train_case(source.build(label)) for label in source.TRAIN]
    info={}
    for arm in model.ARMS[:3]:
        dictionary,updates,footprint=model.frozen_dictionary(basis,cert,public,arm)
        one={}
        for role in range(3):
            for inverse in (False,True):
                word=((role,inverse),)
                one[word]=basis.construct(word)
                if not basis.independently_verify(one[word],cert)["verified"]:
                    raise RuntimeError("T21C_ATOMIC_PROOF_FAILURE")
        proofs=0
        for case in public:
            compiled,hits,parts=model.compile_word(
                basis,cert,tuple(case["word"]),dictionary,one)
            for evidence in (
                tuple((role[0],) for role in case["roles"]),
                case["roles"],
            ):
                for start in range(3):
                    got,p=basis.infer(compiled,evidence,start)
                    ref,r=cert.shared_role_certificate(1,evidence,case["word"],start,
                                                       budget=4096)
                    if bool(p["certified"])!=bool(r["certified"]) or got!=ref:
                        raise RuntimeError("T21C_TRAIN_ALGEBRA_SOUNDNESS_FAILURE")
                    proofs+=1
        info[arm]={"entries":len(dictionary),"train_updates":updates,
                   "dictionary_estimated_bytes":footprint,
                   "verified_train_public_proofs":proofs}
    if any(x["entries"]!=16 or x["dictionary_estimated_bytes"]>8192
           or x["train_updates"]>2048
           or x["verified_train_public_proofs"]!=240 for x in info.values()):
        raise RuntimeError("T21C_TRAIN_ONLY_RESOURCE_OR_COUNT_INVALID")
    return {"schema":"wingless.t21c.train-only-preflight.v1",
       "passed":True,"arms":info,"source_train_cases":40,
       "source_DEV_cases_generated":0,"T19_T20_originals_opened":0,
       "future_unseen_cases_used":False,
       "external_independent_replication":False,
       "optimizer_admission":False}

if __name__=="__main__":
    print(json.dumps(main(),sort_keys=True))
