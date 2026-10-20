package main

import (
 "bufio"
 "bytes"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "os"
 "github.com/spoonman136668-ai/Wingless/cognition"
)

// Offline source-free native cognition interface. No queue, tool or broker calls.
func runCognition(in io.Reader,out io.Writer)error {
 scanner:=bufio.NewScanner(io.LimitReader(in,1<<20+1))
 scanner.Buffer(make([]byte,4096),1<<16)
 enc:=json.NewEncoder(out)
 var s cognition.State
 experience:=cognition.NewExperienceMemory()
 count:=0
 for scanner.Scan(){
   count++;if count>4096{return errors.New("too many events")}
   var event struct{
     Op string `json:"op"`
     Action int `json:"action"`
     Prior int `json:"prior"`
     Outcome int `json:"outcome"`
     Rate float64 `json:"rate"`
     NextAction int `json:"next_action"`
     Snapshot json.RawMessage `json:"snapshot"`
     Before *[4]float64 `json:"before"`
     After *[4]float64 `json:"after"`
     Sensory *[4]float64 `json:"sensory"`
   }
   d:=json.NewDecoder(bytes.NewReader(scanner.Bytes()));d.DisallowUnknownFields()
   if err:=d.Decode(&event);err!=nil{return err}
   if err:=d.Decode(new(any));err!=io.EOF{return errors.New("extra JSON data")}
   switch event.Op{
   case "observe":if err:=s.Observe(event.Action,event.Prior,event.Outcome,event.Rate);err!=nil{return err}
   case "predict":
      p,err:=s.Predict(event.Action,event.Prior);if err!=nil{return err}
      if err=enc.Encode(map[string]any{"probability":p});err!=nil{return err}
   case "forecast_two":
      p,err:=s.ForecastTwo(event.Prior,event.Action,event.NextAction);if err!=nil{return err}
      if err=enc.Encode(map[string]any{"probability":p});err!=nil{return err}
   case "snapshot":
      b,err:=s.Snapshot();if err!=nil{return err}
      if _,err=fmt.Fprintln(out,string(b));err!=nil{return err}
   case "restore":if err:=s.Restore(event.Snapshot);err!=nil{return err}
   case "experience":
      if event.Before==nil||event.After==nil{return errors.New("missing temporal sensory pair")}
      if err:=experience.ObserveExperience(*event.Before,*event.After,event.Action,float64(event.Outcome));err!=nil{return err}
   case "infer_experience":
      if event.Sensory==nil{return errors.New("missing sensory observation")}
      value,known,err:=experience.InferExperience(*event.Sensory,event.Action)
      if err!=nil{return err}
      if err:=enc.Encode(map[string]any{"probability":value,"known":known});err!=nil{return err}
   case "experience_snapshot":
      b,err:=experience.ExperienceSnapshot();if err!=nil{return err}
      if _,err=fmt.Fprintln(out,string(b));err!=nil{return err}
   case "experience_restore":
      if err:=experience.ExperienceRestore(event.Snapshot);err!=nil{return err}
   default:return fmt.Errorf("unknown cognition op %q",event.Op)
   }
 }
 return scanner.Err()
}
func cognitionCLI()error{return runCognition(os.Stdin,os.Stdout)}
