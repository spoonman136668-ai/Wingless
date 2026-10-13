"""W-R6 frozen synthetic context formula replay; no model, training or network."""
from collections import defaultdict
import json

def context(entity, index):
    return (192+(7*entity+3*index)%16,
            208+(5*entity+5*index)%16,
            224+(11*entity+7*index)%16,
            240+(13*entity+11*index)%16)

def outcome():
    signatures={}
    for arm in ("A0","A1","A2","A3"):
        classes=defaultdict(list)
        for entity in range(9):
            for family in range(3):
                rows=tuple(context(entity,t) for t in range(16))
                phase=(rows[0] if (arm=="A1" or arm=="A2" and family==0 or arm=="A3" and family==1) else None)
                classes[(tuple(sorted(rows)),phase)].append((entity,family))
        complete=sum(len(v)==3 and len({e for e,f in v})==1 for v in classes.values())
        collisions=sum(len({e for e,f in v})>1 for v in classes.values())
        signatures[arm]={"signature_classes":len(classes),"complete_alias_triplets":complete,"cross_entity_collision_classes":collisions}
    a=signatures
    valid=a["A0"]=={"signature_classes":8,"complete_alias_triplets":7,"cross_entity_collision_classes":1} and a["A1"]=={"signature_classes":9,"complete_alias_triplets":9,"cross_entity_collision_classes":0}
    classification=("INVALID" if not valid else "SUPPORTED" if all(a[z]["complete_alias_triplets"]<9 for z in ("A2","A3")) else "NEGATIVE")
    return {"schema":"wingless.w-r6-replay.v1","classification":classification,"arms":a,"bridge_records":432,"external_model_calls":0,"limits":"Synthetic formula replay; no unseen corpus or native Go scientific certification"}

if __name__=="__main__":
    one=outcome();two=outcome()
    assert one==two
    print(json.dumps(one,sort_keys=True,indent=2))
