"""T22 source-safe TRAIN-only admission checks; NEVER generate DEV task IDs."""
import importlib.util,json,sys
from pathlib import Path
HERE=Path(__file__).parent
def load(n,p):
    spec=importlib.util.spec_from_file_location(n,HERE/p)
    m=importlib.util.module_from_spec(spec)
    sys.modules[n]=m
    spec.loader.exec_module(m)
    return m
def main():
    model=load("t22_model_train","20261019-wingless-t22-sparse-learner.py")
    cert,source=model.deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert not set(source.TRAIN)&set(source.DEV)
    public=[source.public_train_case(source.build(k)) for k in source.TRAIN]
    evidence={}
    for arm in model.ARMS[:3]:
        rules,checked,payload,hist=model.fit(cert,public,arm)
        if checked>384 or len(rules)>16 or payload>4096:raise RuntimeError("T22_CAPACITY_VIOLATION")
        if any(not model.universal_identity(cert,w) for w in rules):
            raise RuntimeError("T22_UNSOUND_DICTIONARY")
        checks=0;changed=0
        for p in public:
            word=tuple(p["word"])
            reduced,hits=model.reduce_with_rules(cert,word,rules)
            changed+=hits
            for start in range(3):
                r,proof=cert.shared_role_certificate(1,p["roles"],word,start,budget=4096)
                q,qp=cert.shared_role_certificate(1,p["roles"],reduced,start,budget=4096)
                if not proof["certified"] or not qp["certified"] or r!=q:
                    raise RuntimeError("T22_TRAIN_IDENTITY_CERT_MISMATCH")
                checks+=1
        evidence[arm]={"identity_rules":len(rules),"public_TRAIN_proofs":checks,
                       "TRAIN_rewrite_hits":changed,"candidates_verified":checked,
                       "distinct_training_substrings":hist,"dictionary_bytes":payload}
    return {"schema":"wingless.t22.train-only-frozen-preflight.v1",
       "status":"PASS_SOURCE_FREE_DEV_UNOPENED","train_cases":40,"DEV_cases_generated":0,
       "arms":evidence,"T16_T21_originals_opened":0,
       "independent_external_source":False,"scientific_acceptance":False,
       "optimizer_admission":False}
if __name__=="__main__":
    print(json.dumps(main(),sort_keys=True))
