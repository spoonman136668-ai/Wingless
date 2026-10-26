package outcome

import (
 "errors"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r8 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r8"
)

// Selected couples the unchanged native R8 selected-action custodian to
// the opt-in outcome ledger. It does not execute source episodes.
type Selected struct {
 Policy *r8.Policy
 Rewards *Ledger
}
func NewSelected() *Selected { return &Selected{Policy:r8.New(),Rewards:New()} }

// Feedback is transactional: validate outcome and apply both changes on
// restored candidates before replacing either live pointer. The caller
// supplies only source-observed selected-action public information.
func(s *Selected) Feedback(before r1.PublicGrid, action int, after r1.PublicGrid, e Event)error {
 if s==nil||s.Policy==nil||s.Rewards==nil||!s.Rewards.Valid()||
  e.Action!=action||e.ID!=s.Rewards.NextID||s.Rewards.Closed {
  return errors.New("OUTCOME_SELECTED_CUSTODY")
 }
 pRaw,err:=s.Policy.Snapshot();if err!=nil{return err}
 p,err:=r8.Restore(pRaw);if err!=nil{return err}
 lRaw,err:=s.Rewards.Snapshot();if err!=nil{return err}
 l,err:=Restore(lRaw);if err!=nil{return err}
 // Prevalidate reward before accepting original selected-action feedback.
 if err=l.Accept(e);err!=nil{return err}
 if err=p.Feedback(before,action,after);err!=nil{return err}
 s.Policy=p
 s.Rewards=l
 return nil
}

func(s *Selected) BeginEpisode()error {
 if s==nil||s.Policy==nil||s.Rewards==nil||!s.Rewards.Valid(){return errors.New("OUTCOME_INVALID_RESET")}
 if s.Rewards.Closed {
  pRaw,err:=s.Policy.Snapshot();if err!=nil{return err}
  p,err:=r8.Restore(pRaw);if err!=nil{return err}
  lRaw,err:=s.Rewards.Snapshot();if err!=nil{return err}
  l,err:=Restore(lRaw);if err!=nil{return err}
  if err=p.BeginEpisode();err!=nil{return err}
  if err=l.Reset();err!=nil{return err}
  s.Policy=p;s.Rewards=l
  return nil
 }
 return errors.New("OUTCOME_RESET_REQUIRES_TERMINAL")
}

// RewardEvidence returns the observed mean reward for a valid action.
// It does not bypass the native selected-action custodian.
func (s *Selected) RewardEvidence(action int) (float64, bool) {
 if s == nil || s.Rewards == nil { return 0, false }
 return s.Rewards.Value(action)
}
