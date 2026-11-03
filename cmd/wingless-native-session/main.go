// wingless-native-session exposes native structured-data operations only.
package main

import (
 "bufio"
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "os"
 "github.com/spoonman136668-ai/Wingless/unitary"
)
type request struct {
 Op string `json:"op"`
 Examples []unitary.NativeSessionCase `json:"examples"`
 Leaky bool `json:"leaky"`
 Edges []unitary.NativeSessionEdge `json:"edges"`
 From int `json:"from"`
 To int `json:"to"`
 Checkpoint json.RawMessage `json:"checkpoint"`
}
func handle(raw []byte,state **unitary.NativeSession)(any,error){
 if unitary.ValidateNativeSessionJSON(raw)!=nil{return nil,fmt.Errorf("INVALID_JSON")}
 var fields map[string]json.RawMessage
 if json.Unmarshal(raw,&fields)!=nil{return nil,fmt.Errorf("INVALID_JSON")}
 var r request
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 if d.Decode(&r)!=nil||d.Decode(new(any))!=io.EOF{return nil,fmt.Errorf("INVALID_REQUEST")}
 allowed:=map[string][]string{"train":{"op","examples","leaky"},
  "observe":{"op","edges"},"predict":{"op","from","to"},
  "checkpoint":{"op"},"restore":{"op","checkpoint"}}
 keys,ok:=allowed[r.Op];if !ok||len(fields)!=len(keys){return nil,fmt.Errorf("INVALID_REQUEST_FIELDS")}
 for _,k:=range keys{if _,ok:=fields[k];!ok{return nil,fmt.Errorf("MISSING_REQUEST_FIELD")}}
 switch r.Op {
 case "train":
  if *state!=nil{return nil,fmt.Errorf("SILENT_RETRAIN_DENIED")}
  s,err:=unitary.TrainNativeSession(r.Examples,r.Leaky);if err!=nil{return nil,err}
  *state=s;return map[string]any{"trained":true},nil
 case "restore":
  if *state!=nil{return nil,fmt.Errorf("SILENT_STATE_REPLACEMENT_DENIED")}
  s,err:=unitary.RestoreNativeSession(r.Checkpoint);if err!=nil{return nil,err}
  *state=s;return map[string]any{"restored":true},nil
 }
 if *state==nil{return nil,fmt.Errorf("TRAIN_OR_RESTORE_REQUIRED")}
 switch r.Op {
 case "observe":
  if err:=(*state).ObserveGraph(r.Edges);err!=nil{return nil,err}
  return map[string]any{"observed":true},nil
 case "predict":
  value,logit,err:=(*state).Predict(r.From,r.To);if err!=nil{return nil,err}
  return map[string]any{"prediction":value,"native_logit":logit},nil
 default:
  value,err:=(*state).Checkpoint();if err!=nil{return nil,err}
  return map[string]any{"checkpoint":json.RawMessage(value)},nil
 }
}
func main(){
 var state *unitary.NativeSession
 scanner:=bufio.NewScanner(os.Stdin);scanner.Buffer(make([]byte,4096),65536)
 out:=json.NewEncoder(os.Stdout)
 for scanner.Scan(){
  value,err:=handle(scanner.Bytes(),&state)
  if err!=nil{out.Encode(map[string]any{"ok":false,"error":err.Error(),"language_capability":false});continue}
  out.Encode(map[string]any{"ok":true,"result":value,"language_capability":false})
 }
 if scanner.Err()!=nil{fmt.Fprintln(os.Stderr,"REQUEST_SIZE_OR_STREAM_ERROR");os.Exit(1)}
}
