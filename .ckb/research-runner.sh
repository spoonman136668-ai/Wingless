#!/usr/bin/env bash
set -euo pipefail
export E="$RUNNER_TEMP/ckb-research-evidence"; rm -rf "$E"; mkdir -p "$E"
test "$(git rev-parse HEAD)" = "$PACKAGE_SHA"
git merge-base --is-ancestor '1d9ded192dacb3dcf2230eeb5ee3c580d4e13660' HEAD
git merge-base --is-ancestor '4595c30daf9ba8d264368b745aa7a5ceef5b368d' HEAD
python3 - <<'PY'
import base64,hashlib,json,os,pathlib
raw=base64.b64decode("eyJhdXRob3JpdHlfd29ya2Zsb3dfcmVmIjoibWFpbiIsImF1dGhvcml0eV93b3JrZmxvd19yZWZfc2hhIjoiM2JmZDBkNmJmZGNjZjI2Yzc1ZDc0MDI5NzllYzk2NThmMDVmMzkwMiIsImNrYl9wbGFuZV9tYWluX3NoYSI6ImU1MzI0MDU5MmViMjA5MWIwYjhmMDQ1ZDc1OGMzNDE3NzQ4NDYzZTIiLCJleHBlcmltZW50IjoiV0xNLUxNLUVYVEVSTkFMLUZVVFVSRS1EQVRBLUxFQVJORUQtU1RBVEUtRk9STUVSLVJFQURPVVQtSU5URVJGQUNFLUZSRVNILVJFUExJQ0FUSU9OLVI1MSIsImZldGNoX2RlY2lzaW9uIjp7InNjaGVtYSI6ImNrYi1wbGFuZS5leHRlcm5hbC1leHBvc3VyZS1mZXRjaC52MSIsImRpc3Bvc2l0aW9uIjoiUkVBRFlfRkVUQ0giLCJzY29wZSI6ImZldGNoLXZlcmlmeS1vbmx5IiwicmVhc29ucyI6bnVsbH0sIm1hbmlmZXN0X3NoYTI1NiI6IjgyYjRmZGMwMzc2NmZiMTZlYTIwMTE5MzU2YTA0NjRlZDI1NmNiNTU0ZWMxMzg3ZWQ4Nzg5ZjUxZGI4MDYwNzIiLCJwbGFuX2RlY2lzaW9uIjp7InNjaGVtYSI6ImNrYi1wbGFuZS5leHRlcm5hbC1leHBvc3VyZS52MSIsImRpc3Bvc2l0aW9uIjoiUkVBRFlfUExBTiIsInNjb3BlIjoicGxhbi1vbmx5IiwicmVhc29ucyI6bnVsbH0sInByZXJlZ2lzdHJhdGlvbl9zaGEiOiI0NTk1YzMwZGFmOWJhOGQyNjQzNjhiNzQ1YWE3YTVjZWVmNWIzNjhkIiwicHJvamVjdCI6IldpbmdsZXNzIiwicmVxdWVzdF9pZCI6Ijk5YTJhODg0MzE1Y2FkYWJlZTU1M2RiZDk5Yzc4MjdiIiwicmVzZWFyY2hfZGVjaXNpb24iOnsic2NoZW1hIjoiY2tiLXBsYW5lLmV4dGVybmFsLWV4cG9zdXJlLXJlc2VhcmNoLnYxIiwiZGlzcG9zaXRpb24iOiJSRUFEWV9SRVNFQVJDSCIsInNjb3BlIjoicmVzZWFyY2gtY29uc3VtZS1vbmNlIiwicmVhc29ucyI6bnVsbH0sInJlc2VhcmNoX2hvc3QiOiJDS0ItUExBTkUtUkVNT1RFIiwicnVubmVyX2lkIjoiQ0tCLVBMQU5FLVJFTU9URSIsInNjaGVtYSI6ImNrYi1wbGFuZS5leHRlcm5hbC1leHBvc3VyZS1yZWFkeS1yZXNlYXJjaC1yZWNlaXB0LnYxIn0=",validate=True)
assert hashlib.sha256(raw).hexdigest()=="01a27addf3b9cae05a872d2a734bf6299693a02733c4c4a6913c2430ba1a56af"
r=json.loads(raw)
assert r["schema"]=="ckb-plane.external-exposure-ready-research-receipt.v1"
assert r["project"]=="Wingless"
assert r["experiment"]=="WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-READOUT-INTERFACE-FRESH-REPLICATION-R51"
assert r["preregistration_sha"]=="4595c30daf9ba8d264368b745aa7a5ceef5b368d"
assert r["manifest_sha256"]=="82b4fdc03766fb16ea20119356a0464ed256cb554ec1387ed8789f51db806072"
assert r["ckb_plane_main_sha"]=="e53240592eb2091b0b8f045d758c3417748463e2"
assert r["research_decision"]["disposition"]=="READY_RESEARCH"
(pathlib.Path(os.environ["E"])/"authority.json").write_bytes(raw)
PY
python3 - <<'PY'
import hashlib,json,os,pathlib,urllib.parse,urllib.request
root=pathlib.Path(os.environ["RUNNER_TEMP"]); names=("transfer","third","fourth","fifth","twenty-fourth","twenty-fifth","twenty-sixth")
dirs={n:root/("r51-"+n) for n in names}
for p in dirs.values(): p.mkdir(parents=True,exist_ok=True)
hist={
"transfer":[("code","https://raw.githubusercontent.com/python/cpython/ebf955df7a89ed0c7968f79faec1de49f61ed7cb/Lib/statistics.py",61896,"3023ec949802c2e880a000775bfd1a6b70c21b5423881118141e3ec69d8e3336",41453,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45"),("structured","https://raw.githubusercontent.com/json-schema-org/JSON-Schema-Test-Suite/5b0ee1613e45fcc2bddac00e07c19cd49b00d8a8/tests/draft2020-12/uniqueItems.json",14490,"ed84ef6ddc827659e257ba64cbe9679f06932b6320bea5c07622bf3d2cc4daa5",14365,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932"),("technical_prose","https://raw.githubusercontent.com/golang/go/6f5c275ebdc454197fff5f1496521c8f81e20eef/doc/README.md",3124,"c1bf000e7a873b3afed329c5b9f9c07f7692b1d7fb2bfca2d7600cec8b1d1407",1454,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b")],
"third":[("code","https://raw.githubusercontent.com/torvalds/linux/a74306e2e676f9775457366fc047a660fbf02f26/kernel/sched/core.c",303863,"9297982652b9810b3508993c06d3c6e097ebca5cbe1ef3a6246bccfa0c1128b1",41453,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff"),("structured","https://raw.githubusercontent.com/SchemaStore/schemastore/de76181a2ab215431d3e9314bc14f83cc3b01ad2/src/schemas/json/github-workflow.json",118601,"d10c9f4656e1bd5bc6727e9b35080e017dc167154726fca93da33c7a6bd1c4f3",14365,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56"),("technical_prose","https://raw.githubusercontent.com/rust-lang/book/1500248d8f230566e4ec9f27fcbb8fe9e2898ab1/src/ch01-03-hello-cargo.md",11025,"61369f359b84b646fc3773eb569a26bc18ba6edb4cf2be06a84472c7054c0e39",1454,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017")],
"fourth":[("code","https://raw.githubusercontent.com/kubernetes/kubernetes/e7967bb9b76d43a6388abb7a4e90d2f899ef97ff/pkg/kubelet/kubelet.go",156626,"e0f247ef45cb281166b58ba1a013dc6925023388193bd74c1f126c695b770d34",41453,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1"),("structured","https://raw.githubusercontent.com/microsoft/TypeScript/50d70a3f5f453a79a4323b263165da51f656a4e3/package-lock.json",143844,"52b6c9bd0a26ef1b0fc27dfbfaa00df6a047cdddd587d07c4f6c1c3345faf5cf",14365,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933"),("technical_prose","https://raw.githubusercontent.com/django/django/a461af8ce48762d7ec602260aaff81014ddccbcb/docs/intro/tutorial01.txt",10844,"1cf1817e6211b7c0fa9e8044ee808fde1e8f16ca5f9d5ab5358510a5c9591b96",1454,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222")],
"fifth":[("code","https://raw.githubusercontent.com/pytorch/pytorch/496340f06ef7bda2800522429ba3f4e3473a92fa/torch/_inductor/scheduler.py",539420,"b13a89d0fd77430e2a637dae41980086909c4ad4a855b8418dc17286c39953b8",41453,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a"),("structured","https://raw.githubusercontent.com/microsoft/vscode/6fad7188e7dbf7db564e5e4a85960eb2bf69bdf4/package-lock.json",784171,"73224989048e176e7eeafa8aa729df7e4d510f4bd453018902d95345db0fe886",14365,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281"),("technical_prose","https://raw.githubusercontent.com/numpy/numpy/2f1eca306857b641fb0fef0ab854a2d4137e1581/doc/source/user/absolute_beginners.rst",52653,"945a0a598d080ee7526c9bd5459ce615a5b49fe0331adfdf8bbd8812e908e7e3",1454,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49")]}
for group,sources in hist.items():
  for domain,url,full_n,full_h,prefix_n,prefix_h in sources:
    data=urllib.request.urlopen(urllib.request.Request(url,headers={"User-Agent":"wingless-r51"}),timeout=90).read(full_n+1)
    assert len(data)==full_n and hashlib.sha256(data).hexdigest()==full_h
    p=data[:prefix_n]; assert hashlib.sha256(p).hexdigest()==prefix_h
    (dirs[group]/{"code":"code.bin","structured":"structured.bin","technical_prose":"technical-prose.bin"}[domain]).write_bytes(p)
m=json.load(open("docs/experiments/wlm-lm-external-future-data-r51-manifest.json"))
assert m["manifest_sha256"]=="82b4fdc03766fb16ea20119356a0464ed256cb554ec1387ed8789f51db806072"
for group in m["unseen"]:
  for src in group["sources"]:
    url=f"https://raw.githubusercontent.com/{src['repo']}/{src['commit']}/{urllib.parse.quote(src['path'])}"
    data=urllib.request.urlopen(urllib.request.Request(url,headers={"User-Agent":"wingless-r51"}),timeout=90).read()
    assert hashlib.sha1((f"blob {len(data)}\0").encode()+data).hexdigest()==src["git_blob"]
    p=data[:src["prefix_bytes"]];assert hashlib.sha256(p).hexdigest()==src["prefix_sha256"]
    (dirs[group["name"]]/{"code":"code.bin","structured":"structured.bin","technical_prose":"technical-prose.bin"}[src["domain"]]).write_bytes(p)
PY
go test ./unitary -count=1
go test ./cmd/wlm-lm-external-future-data-learned-state-former-readout-interface-fresh-replication-r51 -count=1
go test ./... -count=1
args=()
for n in transfer third fourth fifth twenty-fourth twenty-fifth twenty-sixth; do key="${n//-/_}"; args+=("--${n}-root" "$RUNNER_TEMP/r51-$n"); done
go run ./cmd/wlm-lm-external-future-data-learned-state-former-readout-interface-fresh-replication-r51 "${args[@]}" --resource-out "$E/resource.json" > "$E/result1.json"
go run ./cmd/wlm-lm-external-future-data-learned-state-former-readout-interface-fresh-replication-r51 "${args[@]}" > "$E/result2.json"
cmp "$E/result1.json" "$E/result2.json"
python3 - <<'PY'
import json,math,os,pathlib
E=pathlib.Path(os.environ["RUNNER_TEMP"])/"ckb-research-evidence";r=json.loads((E/"result1.json").read_text());m=r["metrics"];a={x["id"]:x for x in r["arms"]};full=a["full-6"];sel=a["two-group-2"]
valid=(m["invalid_row_count"]==0 and m["source_identity_mismatch_count"]==0 and m["readout_arm_count"]==2 and m["state_dimension"]==6 and m["total_adaptation_budget"]==1744 and m["heldout_outcome_use_before_arm_freeze"]==0 and m["post_result_arm_or_threshold_choice_count"]==0 and all(math.isfinite(float(v)) for v in m.values()))
support=valid and sel["accuracy_delta_vs_full"]>=0.05 and sel["positive_manifest_count_vs_full"]>=2 and sel["mean_decision_regret"]<=full["mean_decision_regret"]
mixed=valid and not support and sel["accuracy_delta_vs_full"]>=0.015 and sel["positive_manifest_count_vs_full"]>=2 and sel["mean_decision_regret"]<=full["mean_decision_regret"]
cls="invalid" if not valid else "supported" if support else "mixed" if mixed else "negative"
out={"schema":"wingless.research-classification.v1","experiment":"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-READOUT-INTERFACE-FRESH-REPLICATION-R51","classification":cls,"validity_pass":valid,"selected_readout":"two-group-2","fresh_replication":True,"rsi_success":bool(support),"next_successor":None,"package_sha":os.environ["PACKAGE_SHA"],"metrics":m,"arms":r["arms"]}
(E/"classification.json").write_text(json.dumps(out,separators=(",",":")),encoding="utf-8")
print(json.dumps(out,separators=(",",":")))
PY
rm -rf "$RUNNER_TEMP"/r51-*
