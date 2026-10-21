// wingless-ice-cognition is an opt-in native ICE learning research consumer.
// It cannot dispatch workers, issue tools or change acceptance.
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "path/filepath"
 "os"
 "strings"
 "github.com/spoonman136668-ai/Wingless/cognition"
)
// The native task outcome is produced by executed project source, never by
// caller-selected truth. Restrict the snapshot outside canonical ICE indexing
// so source records cannot be invalidated by mutable execution state.
func checkedNativeAction(from,to string)(int,error){
 if to!="-"{return 0,fmt.Errorf("native action accepts no target/answer label")}
 switch from{
 case "0":return 0,nil
 case "1":return 1,nil
 case "2":return 2,nil
 default:return 0,fmt.Errorf("native action must be 0, 1 or 2")
 }
}
func checkedNativeSnapshot(root,snapshot string)error{
 absRoot,err:=filepath.EvalSymlinks(root);if err!=nil{return err}
 absRoot,err=filepath.Abs(absRoot);if err!=nil{return err}
 parent,err:=filepath.EvalSymlinks(filepath.Dir(snapshot));if err!=nil{return err}
 parent,err=filepath.Abs(parent);if err!=nil{return err}
 rel,err:=filepath.Rel(absRoot,parent);if err!=nil{return err}
 if rel=="."||rel==".."||(!strings.HasPrefix(rel,".."+string(filepath.Separator))&&!filepath.IsAbs(rel)){
  return fmt.Errorf("native mutable snapshot must be outside canonical source root")
 }
 return nil
}

func run(args []string)error{
 if len(args)!=5{return fmt.Errorf("usage: wingless-ice-cognition ROOT SNAPSHOT observe|predict|native-observe|native-predict|native-intervene FROM TO")}
 root,path,op,source,target:=args[0],args[1],args[2],args[3],args[4]
 r,e:=cognition.OpenSourceReasoner(root);if e!=nil{return e}
 if _,e=os.Lstat(path);e==nil{
   if e=r.Load(path);e!=nil{return e}
 }else if !os.IsNotExist(e){return e}
 switch op{
 case "native-intervene":
   if source!="-"||target!="-" {return fmt.Errorf("native intervention does not accept outcome/identity arguments")}
   if err:=checkedNativeSnapshot(root,path);err!=nil{return err}
   receipt,err:=r.ObserveNativeICEIntervention(context.Background());if err!=nil{return err}
   if err:=r.Save(path);err!=nil{return err}
   return json.NewEncoder(os.Stdout).Encode(receipt)
 case "native-observe":
   action,err:=checkedNativeAction(source,target);if err!=nil{return err}
   if err=checkedNativeSnapshot(root,path);err!=nil{return err}
   result,err:=r.ObserveNativeReasoningTask(context.Background(),action);if err!=nil{return err}
   if err=r.Save(path);err!=nil{return err}
   return json.NewEncoder(os.Stdout).Encode(result)
 case "native-predict":
   action,err:=checkedNativeAction(source,target);if err!=nil{return err}
   p,err:=r.PredictNativeReasoningTask(action);if err!=nil{return err}
   return json.NewEncoder(os.Stdout).Encode(map[string]any{"learned_probability":p,"source_index_sha":r.IndexSHA,"scientific_transfer":"NOT_QUALIFIED"})
 case "observe":
   if e=r.ObserveExistingCall(source,target);e!=nil{return e}
   return r.Save(path)
 case "predict":
   ans,e:=r.PredictCall(source,target);if e!=nil{return e}
   return json.NewEncoder(os.Stdout).Encode(ans)
 }
 return fmt.Errorf("unrecognized task op %q",op)
}
func main(){
 if e:=run(os.Args[1:]);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
}
