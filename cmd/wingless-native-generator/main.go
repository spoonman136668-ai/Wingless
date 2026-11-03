package main

import("bufio";"bytes";"encoding/json";"fmt";"io";"os";"github.com/spoonman136668-ai/Wingless/unitary")
type request struct {Op string `json:"op"`;Examples []unitary.NativeTextExample `json:"examples"`;Input string `json:"input"`;Checkpoint json.RawMessage `json:"checkpoint"`}
func handle(raw []byte,state **unitary.NativeParallelGenerator)(any,bool,error){
 if unitary.ValidateNativeGeneratorJSON(raw)!=nil{return nil,false,fmt.Errorf("GENERATOR_JSON_INVALID")}
 var fields map[string]json.RawMessage;if json.Unmarshal(raw,&fields)!=nil{return nil,false,fmt.Errorf("GENERATOR_REQUEST_INVALID")}
 var r request;d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 if d.Decode(&r)!=nil||d.Decode(new(any))!=io.EOF{return nil,false,fmt.Errorf("GENERATOR_REQUEST_INVALID")}
 allowed:=map[string][]string{"train":{"op","examples"},"generate":{"op","input"},"checkpoint":{"op"},"restore":{"op","checkpoint"}}
 keys,ok:=allowed[r.Op];if !ok||len(fields)!=len(keys){return nil,false,fmt.Errorf("GENERATOR_REQUEST_FIELDS_INVALID")}
 for _,k:=range keys{if _,ok:=fields[k];!ok{return nil,false,fmt.Errorf("GENERATOR_REQUEST_FIELDS_INVALID")}}
 if r.Op=="train"||r.Op=="restore"{
  if *state!=nil{return nil,false,fmt.Errorf("SILENT_STATE_REPLACEMENT_DENIED")}
  var m *unitary.NativeParallelGenerator;var err error
  if r.Op=="train"{m,err=unitary.TrainNativeParallelGenerator(r.Examples)}else{m,err=unitary.RestoreNativeParallelGenerator(r.Checkpoint)}
  if err!=nil{return nil,false,err};*state=m;return map[string]any{"ready":true},true,nil
 }
 if *state==nil{return nil,false,fmt.Errorf("TRAIN_OR_RESTORE_REQUIRED")}
 if r.Op=="generate"{p,err:=(*state).Generate(r.Input);return p,p.ValidOutput,err}
 snapshot,err:=(*state).Checkpoint();if err!=nil{return nil,false,err};return map[string]any{"checkpoint":json.RawMessage(snapshot)},true,nil
}
func main(){
 var state *unitary.NativeParallelGenerator
 scanner:=bufio.NewScanner(os.Stdin);scanner.Buffer(make([]byte,4096),unitary.NativeGeneratorJSONCap)
 out:=json.NewEncoder(os.Stdout)
 for scanner.Scan(){
  value,valid,err:=handle(scanner.Bytes(),&state)
  if err!=nil{out.Encode(map[string]any{"ok":false,"error":err.Error(),"language_proof":false});continue}
  out.Encode(map[string]any{"ok":valid,"result":value,"language_proof":false})
 }
 if scanner.Err()!=nil{fmt.Fprintln(os.Stderr,"GENERATOR_REQUEST_SIZE_OR_STREAM_ERROR");os.Exit(1)}
}
