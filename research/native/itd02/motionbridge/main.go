// ITD-02 opt-in source-public motion IPC adapter. Development-only; it never executes a
// world or authenticates a source episode. Native R8 selected custody is intact.
package main

import (
 "bufio"
 "bytes"
 "encoding/base64"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "os"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 outcome "github.com/spoonman136668-ai/Wingless/research/native/itd02/outcome"
)

const maxWireRequest=8<<20

type Request struct{
 Op string `json:"op"`
 Before r1.PublicGrid `json:"before"`
 After r1.PublicGrid `json:"after"`
 Action int `json:"action"`
 Reward float64 `json:"reward"`
 Terminal bool `json:"terminal"`
 Truncated bool `json:"truncated"`
 EventID uint64 `json:"event_id"`
 MinCount uint32 `json:"min_count"`
 State string `json:"state"`
}
type Reply struct{
 OK bool `json:"ok"`
 Action *int `json:"action,omitempty"`
 State string `json:"state,omitempty"`
 Updates uint64 `json:"updates,omitempty"`
 Error string `json:"error,omitempty"`
}
func decode(raw []byte)(Request,error){
 var r Request
 d:=json.NewDecoder(bytes.NewReader(raw))
 d.DisallowUnknownFields()
 if err:=d.Decode(&r);err!=nil{return Request{},err}
 if err:=d.Decode(new(any));err!=io.EOF{return Request{},errors.New("MOTION_TRAILING_JSON")}
 return r,nil
}
func process(p **outcome.MotionSelected,r Request)(Reply,bool){
 v:=*p
 var response Reply
 var err error
 switch r.Op{
 case "choose":
  var a int
  a,err=v.Choose(r.Before,r.MinCount)
  if err==nil{response.Action=&a}
 case "peek_frozen":
  // Counterfactual action query against an isolated checkpoint clone.
  // Never arms the live route, consumes feedback, or learns during probes.
  if v.Base.Base.Policy.Armed {
   err=errors.New("MOTION_FROZEN_PEEK_PENDING_ACTION")
   break
  }
  var snap []byte
  snap,err=v.Snapshot()
  if err!=nil{break}
  var clone *outcome.MotionSelected
  clone,err=outcome.RestoreMotionSelected(snap)
  if err!=nil{break}
  var a int
  a,err=clone.Choose(r.Before,r.MinCount)
  if err==nil{response.Action=&a}
 case "feedback":
  err=v.Feedback(r.Before,r.Action,r.After,outcome.Event{
   ID:r.EventID,Action:r.Action,Reward:r.Reward,
   Terminal:r.Terminal,Truncated:r.Truncated,
  })
 case "episode_start":
  err=v.BeginEpisode()
 case "snapshot":
  var raw []byte
  raw,err=v.Snapshot()
  if err==nil{response.State=base64.StdEncoding.EncodeToString(raw)}
 case "load":
  if v.Base!=nil&&v.Base.Base!=nil&&v.Base.Base.Policy!=nil&&v.Base.Base.Policy.Armed{
   err=errors.New("MOTION_LOAD_PENDING_SELECTED_ACTION")
   break
  }
  var raw []byte
  raw,err=base64.StdEncoding.DecodeString(r.State)
  if err==nil{
   var restored *outcome.MotionSelected
   restored,err=outcome.RestoreMotionSelected(raw)
   if err==nil{*p=restored;v=restored}
  }
 case "quit":
  if v.Base!=nil&&v.Base.Base!=nil&&v.Base.Base.Policy!=nil&&v.Base.Base.Policy.Armed{
   err=errors.New("MOTION_QUIT_PENDING_SELECTED_ACTION")
   break
  }
  response.OK=true
  return response,true
 default:
  err=errors.New("MOTION_UNKNOWN_OPERATION")
 }
 if err!=nil{response.Error=err.Error();return response,false}
 response.OK=true
 response.Updates=v.Base.Base.Policy.LiveSeen
 return response,false
}
func serve(in io.Reader,out io.Writer)error{
 scan:=bufio.NewScanner(in)
 scan.Buffer(make([]byte,4096),maxWireRequest)
 bw:=bufio.NewWriter(out)
 p:=outcome.NewMotionSelected()
 for scan.Scan(){
  response:=Reply{}
  quit:=false
  r,err:=decode(scan.Bytes())
  if err!=nil{response.Error=err.Error()}else{
   response,quit=process(&p,r)
  }
  encoded,err:=json.Marshal(response)
  if err!=nil{return err}
  if _,err=bw.Write(append(encoded,'\n'));err!=nil{return err}
  if err=bw.Flush();err!=nil{return err}
  if quit{return nil}
 }
 return scan.Err()
}
func main(){
 if err:=serve(os.Stdin,os.Stdout);err!=nil{
  fmt.Fprintln(os.Stderr,err)
  os.Exit(2)
 }
}
