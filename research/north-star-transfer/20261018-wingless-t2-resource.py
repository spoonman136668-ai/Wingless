"""Wingless T2 cost-bound comparison against unchanged actual W-F55 exact public solver.

Uses only negative synthetic fixture ID -900. TRAIN/UNSEEN/REPLICA remain
sealed. Learning the prior on the same fixture is a COST STRESS TEST, not a
transfer evaluation or a valid generalization training protocol.
"""
from __future__ import annotations
import importlib.util,sys,subprocess,time,tracemalloc,json
from pathlib import Path
ROOT=Path(__file__).resolve().parents[2]
HERE=Path(__file__).resolve().parent
F55=ROOT/"research/foundational/20261015-w-f55/source_pilot.py"
KERNEL=HERE/"20261018-wingless-learned-kernel.py"
SOURCE=HERE/"20261018-wingless-t2-source.py"
PIN={"f55":"d02fc933b52fc44d537ed45b061aef56ccd90a7f",
     "kernel":"1c4da8e6ddea0d76f07fa3a47569d3ed18afe417"}

def need(ok,name):
    if not ok:raise ValueError("WING_T2_COST_"+name)
def blob(path):
    return subprocess.run(["git","hash-object",str(path)],
          cwd=ROOT,check=True,capture_output=True,text=True).stdout.strip()
def load(name,p):
    sp=importlib.util.spec_from_file_location(name,p)
    m=importlib.util.module_from_spec(sp);sys.modules[name]=m;sp.loader.exec_module(m)
    return m
def measured(fn):
    tracemalloc.start()
    tic=time.process_time_ns()
    output=fn()
    cpu=time.process_time_ns()-tic
    peak=tracemalloc.get_traced_memory()[1]
    tracemalloc.stop()
    return output,{"CPU_ns":cpu,"Python_peak_bytes":peak}
def attest():
    need(blob(F55)==PIN["f55"] and blob(KERNEL)==PIN["kernel"],"EXACT_F55_LEARNED_OR_CONTROL_SOURCE")
    old=load("t2_w55",F55)
    new=load("t2_wlearned",KERNEL)
    source=load("t2_w_source",SOURCE)
    pub,priv=source.build_fixture(-900,"C7",10)
    role_samples=[list(pub["public_roles"][r]) for r in sorted(pub["public_roles"])]
    (packed,train_meta),training=measured(lambda:new.learn_shared_prior(role_samples))
    def candidate():
        state=new.start()
        for role,direction in pub["public_chain"]:
            state,_=new.compose(state,pub["public_roles"][role],packed,
                   inverse=direction==-1)
        return new.apply(state,pub["start"])
    role_public=[{"role":r,"pairs":list(pub["public_roles"][r])[:2]}
                 for r in sorted(pub["public_roles"])]
    def full_public_solver():
        return old.infer(role_public,pub["start"],pub["public_chain"],
                   "public_full_algebraic_solver",(0,)*12)[0]
    cand,cm=measured(candidate)
    exact,bm=measured(full_public_solver)
    need(exact==priv["expected_answer"],"HONEST_FULL_ALGEBRAIC_CONTROL")
    return {"schema":"wingless.t2.cost-diagnostic.v1",
       "source_case":"negative_synthetic_fixture_ONLY",
       "learned_packed_weights_and_state_bytes":len(packed)+len(new.start()),
       "max_learned_payload_bytes":23,
       "model_cpu_including_amortized_TRAIN_one_fixture_over_128":
           cm["CPU_ns"]+training["CPU_ns"]/128,
       "model_task_CPU_ns":cm["CPU_ns"],
       "model_TRAIN_CPU_ns":training["CPU_ns"],
       "exact_algebraic_control_CPU_ns":bm["CPU_ns"],
       "model_Python_peak_bytes":max(cm["Python_peak_bytes"],training["Python_peak_bytes"]),
       "control_Python_peak_bytes":bm["Python_peak_bytes"],
       "candidate_answer_correct_on_synthetic":cand==priv["expected_answer"],
       "control_answer_correct_on_synthetic":True,
       "noncommuting_S3_representable":False,
       "total_training_and_memory_parity_qualified":False,
       "sealed_original_first_use":False,
       "scientific_transfer_qualified":False}
if __name__=="__main__":print(json.dumps(attest(),sort_keys=True))
