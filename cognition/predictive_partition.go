package cognition
// NCI04-C experimental predictive equivalence states. State splits follow
// observed predictive residual and bounded description cost, never caller-supplied
// entity identifiers. No claim of causal identifiability or generalized cognition.
import("bytes";"encoding/json";"errors";"io";"math")
const PartitionCapacity=8
const partitionSchema=1
type partitionCell struct {
 Used bool
 Prototype [4]float64
 Drift [4]float64
 Effect [2]float64
 Delta [2][4]float64
 Seen [2]uint32
 Last uint32
}
type PredictivePartition struct {
 Schema int
 Cells [PartitionCapacity]partitionCell
 Stability [4]float64
 Updates uint32
}
func NewPredictivePartition()*PredictivePartition{
 p:=&PredictivePartition{Schema:partitionSchema};for i:=range p.Stability{p.Stability[i]=.2};return p
}
func partitionValid(v [4]float64)bool{
 for _,x:=range v{if math.IsInf(x,0)||math.IsNaN(x)||x < -1||x>1{return false}};return true
}
func(p *PredictivePartition)check()error{
 if p==nil||p.Schema!=partitionSchema{return errors.New("invalid predictive partition schema")}
 for _,s:=range p.Stability{if math.IsNaN(s)||math.IsInf(s,0)||s<0||s>2{return errors.New("invalid learned stability")}}
 for _,c:=range p.Cells{if !c.Used{continue};if !partitionValid(c.Prototype){return errors.New("invalid prototype")}
  for _,v:=range c.Drift{if math.IsNaN(v)||math.IsInf(v,0)||v<0||v>2{return errors.New("invalid drift")}}
  for _,v:=range c.Effect{if math.IsNaN(v)||math.IsInf(v,0)||v<0||v>1{return errors.New("invalid effect")}}
  for _,d:=range c.Delta{for _,v:=range d{if math.IsNaN(v)||math.IsInf(v,0)||math.Abs(v)>2{return errors.New("invalid transition")}}}
 }
 return nil
}
func(p *PredictivePartition)used()int{n:=0;for _,c:=range p.Cells{if c.Used{n++}};return n}
func(p *PredictivePartition)nearest(view [4]float64)(int,float64,float64){
 best:=-1;distance:=math.Inf(1);runner:=math.Inf(1)
 for i,c:=range p.Cells{
  if !c.Used{continue}
  d:=0.0;for j,x:=range view{z:=x-c.Prototype[j];d+=z*z/(.03+p.Stability[j]+c.Drift[j])}
  if d<distance{best,runner,distance=i,distance,d}else if d<runner{runner=d}
 }
 return best,distance,runner
}
func predictiveError(c partitionCell,before,after [4]float64,action int,outcome float64)float64{
 d:=outcome-c.Effect[action];loss:=d*d
 for j:=0;j<4;j++{v:=after[j]-before[j]-c.Delta[action][j];loss+=.5*v*v}
 return loss
}
func coldError(before,after [4]float64,outcome float64)float64{
 t:=outcome-.5;loss:=t*t
 for j:=0;j<4;j++{v:=after[j]-before[j];loss+=.5*v*v}
 return loss
}
func newCell(before [4]float64)partitionCell{
 return partitionCell{Used:true,Prototype:before,Effect:[2]float64{.5,.5}}
}
func(p *PredictivePartition)Learn(before,after [4]float64,action int,outcome float64)error{
 if err:=p.check();err!=nil{return err}
 if !partitionValid(before)||!partitionValid(after)||action<0||action>1||outcome<0||outcome>1||math.IsNaN(outcome)||math.IsInf(outcome,0){return errors.New("invalid experience")}
 if p.Updates==math.MaxUint32{return errors.New("event budget exhausted")}
 next:=*p;selected:=-1;bestLoss:=math.Inf(1)
 for i,c:=range next.Cells{if !c.Used{continue}
  loss:=predictiveError(c,before,after,action,outcome)
  proximity:=0.0
  for j:=0;j<4;j++{d:=before[j]-c.Prototype[j];proximity+=d*d/(.03+next.Stability[j]+c.Drift[j])}
  loss+=.02*proximity
  if loss<bestLoss{bestLoss,selected=loss,i}
 }
 fresh:=coldError(before,after,outcome)+.04
 if selected<0||(fresh<bestLoss&&next.used()<PartitionCapacity){
  for i,c:=range next.Cells{if !c.Used{selected=i;next.Cells[i]=newCell(before);break}}
 }
 if selected<0{return errors.New("model allocation failure")}
 cell:=&next.Cells[selected]
 for j:=0;j<4;j++{
  difference:=math.Abs(after[j]-before[j])
  next.Stability[j]=.96*next.Stability[j]+.04*difference
  old:=cell.Prototype[j]
  cell.Prototype[j]=math.Max(-1,math.Min(1,.92*old+.08*before[j]))
  cell.Drift[j]=.9*cell.Drift[j]+.1*math.Abs(before[j]-old)
  cell.Delta[action][j]=.85*cell.Delta[action][j]+.15*(after[j]-before[j])
 }
 cell.Effect[action]=.85*cell.Effect[action]+.15*outcome
 if cell.Seen[action]<math.MaxUint32{cell.Seen[action]++}
 next.Updates++;cell.Last=next.Updates;*p=next
 return nil
}
func(p *PredictivePartition)Infer(view [4]float64,action int)(float64,bool,error){
 if err:=p.check();err!=nil{return 0,false,err}
 if !partitionValid(view)||action<0||action>1{return 0,false,errors.New("invalid query")}
 i,best,second:=p.nearest(view)
 if i<0||p.Cells[i].Seen[action]==0{return 0,false,nil}
 if second-best<=1e-6{
  for j,c:=range p.Cells{if j!=i&&c.Used&&c.Seen[action]>0&&math.Abs(c.Effect[action]-p.Cells[i].Effect[action])>.2{return 0,false,nil}}
 }
 return p.Cells[i].Effect[action],true,nil
}
func(p *PredictivePartition)Snapshot()([]byte,error){if err:=p.check();err!=nil{return nil,err};return json.Marshal(p)}
func(p *PredictivePartition)Restore(blob []byte)error{
 if p==nil||len(blob)==0||len(blob)>20000{return errors.New("invalid checkpoint bytes")}
 dec:=json.NewDecoder(bytes.NewReader(blob));dec.DisallowUnknownFields()
 var next PredictivePartition
 if err:=dec.Decode(&next);err!=nil{return err}
 var extra any;if err:=dec.Decode(&extra);err!=io.EOF{return errors.New("trailing checkpoint data")}
 if err:=next.check();err!=nil{return err}
 *p=next;return nil
}
