package cognition

// This opt-in adapter consumes the real immutable ICE source index.  It never
// interprets unverified model output as source truth or an accepted result.
import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "strings"
 "github.com/spoonman136668-ai/Wingless/ice"
)

type SourceReasoner struct{
 Root string
 State LatentState
 IndexSHA string
}
type sourceEnvelope struct{
 Schema string `json:"schema"`
 IndexSHA string `json:"index_sha"`
 State LatentState `json:"state"`
}
func OpenSourceReasoner(root string)(*SourceReasoner,error){
 if err:=ice.Validate(root);err!=nil{return nil,err}
 b,err:=os.ReadFile(filepath.Join(root,".ice","manifest.json"))
 if err!=nil{return nil,err}
 h:=sha256.Sum256(b)
 return &SourceReasoner{Root:root,IndexSHA:hex.EncodeToString(h[:])},nil
}
func (r *SourceReasoner) validate()error{
 if err:=ice.Validate(r.Root);err!=nil{return err}
 b,err:=os.ReadFile(filepath.Join(r.Root,".ice","manifest.json"));if err!=nil{return err}
 h:=sha256.Sum256(b)
 if r.IndexSHA!=hex.EncodeToString(h[:]){return errors.New("ICE source identity drift")}
 return nil
}
func (r *SourceReasoner) symbol(name string)(ice.Symbol,error){
 var zero ice.Symbol
 if name==""||len(name)>160{return zero,errors.New("invalid symbol query")}
 hits,err:=ice.Query(r.Root,name);if err!=nil{return zero,err}
 matched:=[]ice.Symbol{}
 for _,h:=range hits{if h.Name==name && (h.Kind=="function"||h.Kind=="method"){matched=append(matched,h)}}
 if len(matched)!=1{return zero,fmt.Errorf("symbol must resolve uniquely: %q",name)}
 return matched[0],nil
}
func vectorFromSource(s ice.Symbol)[4]float32{
 b:=sha256.Sum256([]byte(s.File+"\x00"+s.Kind+"\x00"+s.Receiver+"\x00"+s.Name))
 var v [4]float32
 for i:=0;i<4;i++{v[i]=float32(int(b[2*i])*256+int(b[2*i+1]))/32767.5-1}
 return v
}
func lexicalCall(from,to ice.Symbol)bool{
 for _,c:=range from.Calls{
   if c==to.Name || strings.HasSuffix(c,"."+to.Name){return true}
 }
 return false
}
// ObserveExistingCall admits only positive syntax-call observations already
// supported by ICE. Missing syntax is unknown, NOT proof of an absent call.
func (r *SourceReasoner) ObserveExistingCall(fromName,toName string)error{
 if e:=r.validate();e!=nil{return e}
 from,e:=r.symbol(fromName);if e!=nil{return e}
 to,e:=r.symbol(toName);if e!=nil{return e}
 if !lexicalCall(from,to){return errors.New("unverified source call; no learned update")}
 return r.State.Observe(vectorFromSource(from),vectorFromSource(to),0,1)
}
type SourcePrediction struct{
 LearnedProbability float64 `json:"learned_probability"`
 VerifiedSyntaxCall bool `json:"verified_syntax_call"`
 EvidenceBasis string `json:"evidence_basis"`
 IndexSHA string `json:"index_sha"`
}
func (r *SourceReasoner) PredictCall(fromName,toName string)(SourcePrediction,error){
 if e:=r.validate();e!=nil{return SourcePrediction{},e}
 from,e:=r.symbol(fromName);if e!=nil{return SourcePrediction{},e}
 to,e:=r.symbol(toName);if e!=nil{return SourcePrediction{},e}
 p,e:=r.State.Predict(vectorFromSource(from),vectorFromSource(to),0);if e!=nil{return SourcePrediction{},e}
 return SourcePrediction{p,lexicalCall(from,to),"syntax-only; exact ICE source authoritative",r.IndexSHA},nil
}
func (r *SourceReasoner) Save(path string)error{
 if e:=r.validate();e!=nil{return e}
 if path==""{return errors.New("empty snapshot path")}
 if st,e:=os.Lstat(path);e==nil&&st.Mode()&os.ModeSymlink!=0{return errors.New("snapshot symlink rejected")}else if e!=nil&&!os.IsNotExist(e){return e}
 data,e:=json.Marshal(sourceEnvelope{"wingless.nci04d2.ice-state.v1",r.IndexSHA,r.State});if e!=nil{return e}
 if len(data)>65536{return errors.New("oversized state")}
 dir:=filepath.Dir(path)
 f,e:=os.CreateTemp(dir,".wingless-ice-*");if e!=nil{return e}
 name:=f.Name();defer os.Remove(name)
 if e=f.Chmod(0600);e!=nil{f.Close();return e}
 if _,e=f.Write(data);e!=nil{f.Close();return e}
 if e=f.Sync();e!=nil{f.Close();return e}
 if e=f.Close();e!=nil{return e}
 return os.Rename(name,path)
}
func (r *SourceReasoner) Load(path string)error{
 if e:=r.validate();e!=nil{return e}
 st,e:=os.Lstat(path);if e!=nil{return e}
 if !st.Mode().IsRegular()||st.Size()>65536{return errors.New("invalid snapshot type or size")}
 b,e:=os.ReadFile(path);if e!=nil{return e}
 var snapshot sourceEnvelope
 if e=json.Unmarshal(b,&snapshot);e!=nil{return e}
 if snapshot.Schema!="wingless.nci04d2.ice-state.v1"||snapshot.IndexSHA!=r.IndexSHA{return errors.New("snapshot wrong schema/source")}
 raw,e:=json.Marshal(snapshot.State);if e!=nil{return e}
 var next LatentState
 if e=next.Restore(raw);e!=nil{return e}
 r.State=next
 return nil
}
