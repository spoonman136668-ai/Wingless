package main

// Offline, project-native source analysis only. ICE remains authoritative.
import (
 "encoding/json"
 "fmt"
 "os"
 "github.com/spoonman136668-ai/Wingless/cognition"
)
func runICESourceCognition(root,snapshot,op,from,to string)error{
 session,err:=cognition.OpenSourceReasoner(root);if err!=nil{return err}
 if _,err=os.Lstat(snapshot);err==nil{
  if err=session.Load(snapshot);err!=nil{return err}
 }else if !os.IsNotExist(err){return err}
 switch op{
 case "observe":
  if err=session.ObserveExistingCall(from,to);err!=nil{return err}
  return session.Save(snapshot)
 case "predict":
  answer,err:=session.PredictCall(from,to);if err!=nil{return err}
  return json.NewEncoder(os.Stdout).Encode(answer)
 default:return fmt.Errorf("unsupported source cognition operation %q",op)
 }
}
