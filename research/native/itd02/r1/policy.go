// Package r1 implements an opt-in Wingless native three-action policy.
//
// It is a SOURCE-ONLY public-action learner. Nothing here dispatches a job,
// accesses a MiniGrid hidden state, or proves unfamiliar-world transfer.
package r1

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "math"
 "io"
)

const Schema = "wingless.itd02.r1.public-grid-action-state.v1"
const Features = 9
const Actions = 3

// PublicGrid is built exclusively from the public MiniGrid FullObsWrapper image.
// There is no task family, seed, answer, shortest path, or hidden grid access.
type PublicGrid struct {
 X int `json:"x"`
 Y int `json:"y"`
 GoalX int `json:"goal_x"`
 GoalY int `json:"goal_y"`
 Direction int `json:"direction"`
 FrontBlocked bool `json:"front_blocked"`
}

// Policy has two independently updated local prediction timescales.
type Policy struct {
 Fast [Actions][Features]float64 `json:"fast"`
 Slow [Actions][Features]float64 `json:"slow"`
 Updates uint64 `json:"updates"`
}

func New() *Policy { return &Policy{} }

func allowed(g PublicGrid) bool {
 return g.X>=0 && g.X<=64 && g.Y>=0 && g.Y<=64 &&
  g.GoalX>=0 && g.GoalX<=64 && g.GoalY>=0 && g.GoalY<=64 &&
  g.Direction>=0 && g.Direction<=3
}

func clipped(x float64) float64 {
 if x < -1 {return -1}
 if x > 1 {return 1}
 return x
}

func features(g PublicGrid) ([Features]float64,error) {
 if !allowed(g) { return [Features]float64{},errors.New("R1_INVALID_PUBLIC_GRID") }
 dx,dy:=g.GoalX-g.X,g.GoalY-g.Y
 dist:=dx
 if dist<0 {dist = -dist}
 ddy:=dy
 if ddy<0 {ddy = -ddy}
 z:=[Features]float64{1,clipped(float64(dx)/16),clipped(float64(dy)/16),
    clipped(float64(dist+ddy)/32)}
 z[4+g.Direction]=1
 if g.FrontBlocked {z[8]=1}
 return z,nil
}

func finite(v float64) bool {return !math.IsNaN(v) && !math.IsInf(v,0)}
func (p *Policy) healthy() bool {
 for a:=0;a<Actions;a++ {
  for i:=0;i<Features;i++ {
   if !finite(p.Fast[a][i])||!finite(p.Slow[a][i]) ||
    math.Abs(p.Fast[a][i])>32 || math.Abs(p.Slow[a][i])>32 {return false}
  }
 }
 return true
}
func (p *Policy) score(z [Features]float64,a int) float64 {
 r:=0.0
 for i,v:=range z {r +=v*(.75*p.Fast[a][i]+.25*p.Slow[a][i])}
 return r
}

// Scores returns all three legal, learned action values. Source-only.
//
// A low score is not an abstention: the only lawful actions are left,
// right and forward. Benchmark failure is recorded if this policy loses.
func (p *Policy) Scores(g PublicGrid) ([Actions]float64,error) {
 z,e:=features(g);if e!=nil{return [Actions]float64{},e}
 if p==nil || !p.healthy(){return [Actions]float64{},errors.New("R1_UNHEALTHY_POLICY")}
 var out [Actions]float64
 for a:=0;a<Actions;a++{out[a]=p.score(z,a)}
 return out,nil
}
func (p *Policy) Choose(g PublicGrid) (int,error) {
 q,e:=p.Scores(g);if e!=nil{return 0,e}
 best:=0
 for a:=1;a<Actions;a++ {if q[a]>q[best]{best=a}}
 return best,nil
}
func taxi(g PublicGrid) float64 {
 dx:=g.GoalX-g.X;if dx<0{dx = -dx}
 dy:=g.GoalY-g.Y;if dy<0{dy = -dy}
 return float64(dx+dy)
}
func bounded(x float64) float64 {
 if x>32 {return 32}
 if x< -32 {return -32}
 return x
}

// Observe takes only a *completed* real public transition. Each update is
// causal: it changes the chosen arm, never retrospectively labels other arms.
func (p *Policy) Observe(before PublicGrid,action int,after PublicGrid,reward float64,terminal bool) error {
 if action<0 || action>=Actions || !finite(reward) || math.Abs(reward)>100 ||
   p==nil || !p.healthy() || p.Updates==^uint64(0) {
  return errors.New("R1_ACTION_REWARD_OR_MODEL_INVALID")
 }
 z,e:=features(before);if e!=nil{return e}
 if _,e=features(after);e!=nil{return e}
 current:=p.score(z,action)
 target:=reward+.15*(taxi(before)-taxi(after))
 if !terminal {
  qs,e:=p.Scores(after);if e!=nil{return e}
  best:=qs[0];for a:=1;a<Actions;a++{if qs[a]>best{best=qs[a]}}
  target+=.9*best
 }
 difference:=clipped(target-current)
 next:=*p
 for i,v:=range z {
  next.Fast[action][i]=bounded(next.Fast[action][i]+.12*difference*v)
  next.Slow[action][i]=bounded(next.Slow[action][i]+.015*difference*v)
 }
 next.Updates++
 if !next.healthy(){return errors.New("R1_INVALID_LEARNED_WEIGHT")}
 *p=next
 return nil
}

type envelope struct {
 Schema string `json:"schema"`
 Payload Policy `json:"payload"`
 Digest string `json:"sha256"`
}

func (p *Policy) Snapshot()([]byte,error) {
 if p==nil || !p.healthy(){return nil,errors.New("R1_INVALID_CHECKPOINT_STATE")}
 body,e:=json.Marshal(p);if e!=nil{return nil,e}
 sum:=sha256.Sum256(body)
 return json.Marshal(envelope{Schema:Schema,Payload:*p,Digest:hex.EncodeToString(sum[:])})
}

func Restore(raw []byte)(*Policy,error) {
 if len(raw)==0 || len(raw)>16384{return nil,errors.New("R1_CHECKPOINT_SIZE")}
 var wire envelope
 dec:=json.NewDecoder(bytes.NewReader(raw));dec.DisallowUnknownFields()
 if e:=dec.Decode(&wire);e!=nil{return nil,e}
 if e:=dec.Decode(new(struct{}));e!=io.EOF{return nil,errors.New("R1_CHECKPOINT_TRAILING_DATA")}
 if wire.Schema!=Schema || !wire.Payload.healthy(){return nil,errors.New("R1_CHECKPOINT_SCHEMA_OR_VALUES")}
 body,e:=json.Marshal(wire.Payload);if e!=nil{return nil,e}
 sum:=sha256.Sum256(body)
 if wire.Digest!=hex.EncodeToString(sum[:]){return nil,errors.New("R1_CHECKPOINT_AUTHENTICATION")}
 return &wire.Payload,nil
}
