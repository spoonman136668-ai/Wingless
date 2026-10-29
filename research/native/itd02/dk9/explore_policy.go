// Package dk9 adds an opt-in native Go uncertainty chooser to the original
// unchanged DK8 local cells. It is NOT original Wingless R8 or mission planning.
package dk9

import (
 "crypto/sha256"
 "fmt"
 old "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk8"
)
type Public=old.Public
type Policy struct{*old.Policy}
func New()*Policy{return &Policy{Policy:old.New()}}
func Restore(b []byte)(*Policy,error){p,e:=old.Restore(b);if e!=nil{return nil,e};return &Policy{Policy:p},nil}
func (p *Policy) Choose(before Public, objective string)(*int,string,error){
 choice,status,err:=p.Policy.Choose(before,objective)
 if err!=nil||choice!=nil{return choice,status,err}
 if status!="ABSTAIN_NO_LEARNED_POSITIVE_ACTION"{return nil,"",fmt.Errorf("DK9_UNRECOGNIZED_ORIGINAL_RESPONSE")}
 key,err:=old.Context(before);if err!=nil{return nil,"",err}
 row:=p.Counts[key]
 h:=sha256.Sum256([]byte(fmt.Sprintf("%s:%d",key,before.Direction)))
 start:=int(h[0])%7
 least:=row[0].N
 for _,cell:=range row{if cell.N<least{least=cell.N}}
 for k:=0;k<7;k++{
   a:=(start+k)%7
   if row[a].N==least{
     result:=a
     return &result,"SELECT_NATIVE_LEAST_OBSERVED_ACTION",nil
   }
 }
 return nil,"",fmt.Errorf("DK9_NATIVE_UNREACHABLE")
}
