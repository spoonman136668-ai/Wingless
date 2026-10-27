// ITD-02 R8 public-action IPC bridge. Development engineering only.
// Native Wingless R8 owns action choice and selected-consequence learning.
// This adapter does not execute an environment or grant scientific admission.
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
    r8 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r8"
)

const maxWireRequest = 8 << 20

type Request struct {
    Op string `json:"op"`
    Before r1.PublicGrid `json:"before"`
    After r1.PublicGrid `json:"after"`
    Action int `json:"action"`
    Reward float64 `json:"reward"`
    Terminal bool `json:"terminal"`
    State string `json:"state"`
}
type Reply struct {
    OK bool `json:"ok"`
    Action *int `json:"action,omitempty"`
    State string `json:"state,omitempty"`
    Updates uint64 `json:"updates,omitempty"`
    Steps uint32 `json:"steps,omitempty"`
    Error string `json:"error,omitempty"`
}
func decodePublic(raw []byte)(Request,error){
    var r Request
    dec:=json.NewDecoder(bytes.NewReader(raw))
    dec.DisallowUnknownFields()
    if err:=dec.Decode(&r);err!=nil{return Request{},err}
    if err:=dec.Decode(new(any));err!=io.EOF{return Request{},errors.New("R8_TRAILING_JSON_DATA")}
    return r,nil
}

// process performs one selected public operation. Rejected load must not
// replace a live learner; a pending chosen action cannot be discarded.
func process(p **r8.Policy,r Request)(Reply,bool){
    v:=*p
    out:=Reply{}
    var err error
    switch r.Op{
    case "choose":
        var a int
        a,err=v.Choose(r.Before)
        if err==nil{out.Action=&a}
    case "feedback":
        err=v.Feedback(r.Before,r.Action,r.After)
    case "learn":
        err=v.ObserveFit(r.Before,r.Action,r.After,r.Reward,r.Terminal)
    case "episode_start":
        err=v.BeginEpisode()
    case "snapshot":
        var b []byte
        b,err=v.Snapshot()
        if err==nil{out.State=base64.StdEncoding.EncodeToString(b)}
    case "snapshot_base":
        var b []byte
        b,err=v.Core.Core.Core.Snapshot()
        if err==nil{out.State=base64.StdEncoding.EncodeToString(b)}
    case "load":
        if v.Armed {
            err=errors.New("R8_LOAD_WHILE_SELECTED_ACTION_PENDING")
            break
        }
        var b []byte
        b,err=base64.StdEncoding.DecodeString(r.State)
        if err==nil {
            var candidate *r8.Policy
            candidate,err=r8.Restore(b)
            if err==nil {*p=candidate;v=candidate}
        }
    case "quit":
        out.OK=true
        return out,true
    default:
        err=errors.New("R8_UNKNOWN_PUBLIC_OPERATION")
    }
    if err!=nil{out.Error=err.Error();return out,false}
    out.OK=true
    out.Updates=v.Core.Core.Core.Events
    out.Steps=v.Core.Core.Route.Steps
    return out,false
}

func serve(in io.Reader,out io.Writer)error{
    sc:=bufio.NewScanner(in)
    sc.Buffer(make([]byte,4096),maxWireRequest)
    bw:=bufio.NewWriter(out)
    p:=r8.New()
    for sc.Scan(){
        var resp Reply
        quit:=false
        req,err:=decodePublic(sc.Bytes())
        if err!=nil{resp.Error=err.Error()}else{
            resp,quit=process(&p,req)
        }
        raw,err:=json.Marshal(resp)
        if err!=nil{return err}
        if _,err=bw.Write(append(raw,'\n'));err!=nil{return err}
        if err=bw.Flush();err!=nil{return err}
        if quit{return nil}
    }
    return sc.Err()
}
func main(){
    if err:=serve(os.Stdin,os.Stdout);err!=nil{
        fmt.Fprintln(os.Stderr,err)
        os.Exit(2)
    }
}
