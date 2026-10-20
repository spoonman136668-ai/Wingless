package nativecognition

// NCI04 source-free predictor. The hidden transition operator and task identity
// are never inputs. Two bounded predictive models are selected using only prior
// prediction error. This does not yet demonstrate learned entity identity.
import ("errors";"math")
const PredictiveChannels=4

type PredictiveState struct {
 Weights [2][4][4]float64
 Active int
 Occupied [2]bool
 Updates int
 Rate float64
}
func NewPredictiveState(rate float64)(*PredictiveState,error){
 if math.IsNaN(rate)||math.IsInf(rate,0)||rate<0||rate>.5{return nil,errors.New("invalid rate")}
 p:=&PredictiveState{Rate:rate};p.Occupied[0]=true
 for k:=0;k<2;k++{for i:=0;i<4;i++{p.Weights[k][i][i]=.75}}
 return p,nil
}
func validVector(x [4]float64)bool{for _,v:=range x{if math.IsNaN(v)||math.IsInf(v,0)||math.Abs(v)>2{return false}};return true}
func (p *PredictiveState)predictUsing(k int,x [4]float64)[4]float64{
 var out [4]float64
 for i:=0;i<4;i++{for j:=0;j<4;j++{out[i]+=p.Weights[k][i][j]*x[j]}}
 return out
}
func (p *PredictiveState)Predict(x [4]float64)([4]float64,error){
 if p==nil||!validVector(x)||p.Active<0||p.Active>=2||!p.Occupied[p.Active]{return [4]float64{},errors.New("invalid predictor/input")}
 return p.predictUsing(p.Active,x),nil
}
func squareError(a,b [4]float64)float64{sum:=0.0;for i:=0;i<4;i++{d:=a[i]-b[i];sum+=d*d};return sum/4}
// Learn uses only consecutive observed states. Surprise-guided context selection
// changes the next prediction; never the one just scored.
func (p *PredictiveState)Learn(previous,observed [4]float64)error{
 if p==nil||!validVector(previous)||!validVector(observed){return errors.New("invalid transition")}
 current:=p.Active;errCurrent:=squareError(p.predictUsing(current,previous),observed)
 chosen:=current
 if p.Updates>=12&&!p.Occupied[1]&&errCurrent>.002{chosen=1;p.Occupied[1]=true
 }else if p.Occupied[1]{
  other:=1-current
  if squareError(p.predictUsing(other,previous),observed)+.00001<errCurrent{chosen=other}
 }
 guess:=p.predictUsing(chosen,previous)
 for i:=0;i<4;i++{for j:=0;j<4;j++{
  next:=p.Weights[chosen][i][j]+p.Rate*(observed[i]-guess[i])*previous[j]
  p.Weights[chosen][i][j]=math.Max(-2,math.Min(2,next))
 }}
 p.Active=chosen;p.Updates++;return nil
}
func (p *PredictiveState)Rollout(start [4]float64,horizon int)([4]float64,error){
 if horizon<1||horizon>4{return [4]float64{},errors.New("invalid horizon")}
 x:=start
 for i:=0;i<horizon;i++{v,e:=p.Predict(x);if e!=nil{return [4]float64{},e};x=v}
 return x,nil
}
func (p *PredictiveState)Erase(){if p!=nil{q,_:=NewPredictiveState(p.Rate);*p=*q}}
// 32-weight trained recurrent comparator with equal observational information.
type MatchedRecurrent struct{W [4][4]float64;V [4][4]float64;H [4]float64;Rate float64}
func NewMatchedRecurrent(rate float64)*MatchedRecurrent{
 p:=&MatchedRecurrent{Rate:rate};for i:=0;i<4;i++{p.W[i][i]=.75};return p
}
func (p *MatchedRecurrent)Predict(x [4]float64)[4]float64{
 var y [4]float64
 for i:=0;i<4;i++{for j:=0;j<4;j++{y[i]+=p.W[i][j]*x[j]+p.V[i][j]*p.H[j]}}
 return y
}
func (p *MatchedRecurrent)Learn(x,y [4]float64){
 guess:=p.Predict(x)
 for i:=0;i<4;i++{for j:=0;j<4;j++{
  e:=y[i]-guess[i]
  p.W[i][j]=math.Max(-2,math.Min(2,p.W[i][j]+p.Rate*e*x[j]))
  p.V[i][j]=math.Max(-2,math.Min(2,p.V[i][j]+p.Rate*e*p.H[j]))
 }}
 for j:=0;j<4;j++{p.H[j]=.7*p.H[j]+.3*x[j]}
}
