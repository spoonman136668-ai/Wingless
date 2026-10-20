package main

import (
 "bytes"
 "encoding/json"
 "strings"
 "testing"
)
func TestOfflineRawExperienceContract(t *testing.T){
 a:=[]float64{-.8,-.9,0,0}
 cmds:=[]map[string]any{}
 for i:=0;i<20;i++ {
   cmds=append(cmds,map[string]any{"op":"experience","before":a,"after":a,"action":1,"outcome":1})
 }
 cmds=append(cmds,map[string]any{"op":"infer_experience","sensory":a,"action":1})
 cmds=append(cmds,map[string]any{"op":"experience_snapshot"})
 cmds=append(cmds,map[string]any{"op":"predict","action":0,"prior":0})
 var input strings.Builder
 for _,c:=range cmds{b,e:=json.Marshal(c);if e!=nil{t.Fatal(e)};input.Write(b);input.WriteByte('\n')}
 var output bytes.Buffer
 if e:=runCognition(strings.NewReader(input.String()),&output);e!=nil{t.Fatal(e)}
 var rows []map[string]any
 dec:=json.NewDecoder(&output)
 for dec.More(){var x map[string]any;if e:=dec.Decode(&x);e!=nil{t.Fatal(e)};rows=append(rows,x)}
 if len(rows)!=3||rows[0]["known"]!=true||rows[1]["schema"]!=float64(1){
   t.Fatalf("invalid observed native interface results: %v",rows)
 }
 if e:=runCognition(strings.NewReader("{\"op\":\"experience\",\"action\":0,\"outcome\":1}\n"),&output);e==nil{t.Fatal("missing sensory accepted")}
}
