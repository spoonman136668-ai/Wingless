// wingless-ice-cognition is an opt-in native ICE learning research consumer.
// It cannot dispatch workers, issue tools or change acceptance.
package main

import (
 "encoding/json"
 "fmt"
 "os"
 "github.com/spoonman136668-ai/Wingless/cognition"
)
func run(args []string)error{
 if len(args)!=5{return fmt.Errorf("usage: wingless-ice-cognition ROOT SNAPSHOT observe|predict FROM TO")}
 root,path,op,source,target:=args[0],args[1],args[2],args[3],args[4]
 r,e:=cognition.OpenSourceReasoner(root);if e!=nil{return e}
 if _,e=os.Lstat(path);e==nil{
   if e=r.Load(path);e!=nil{return e}
 }else if !os.IsNotExist(e){return e}
 switch op{
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
