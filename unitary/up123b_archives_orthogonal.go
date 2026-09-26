package unitary

import "math"

const UP123BArchivesOrthogonalSchema = "wingless.up123b-archives-orthogonal.v1"

type UP123BPoint struct {
	Arm                    string  `json:"arm"`
	ClassOrder             string  `json:"class_order"`
	ReplayPolicy           string  `json:"replay_policy"`
	Stage                  int     `json:"stage"`
	ArchivesOriginalAcc    float64 `json:"archives_original_accuracy"`
	ArchivesUnseenAcc      float64 `json:"archives_unseen_accuracy"`
	ArchivesOrigMeanMargin float64 `json:"archives_original_mean_margin"`
	ArchivesOrigMinMargin  float64 `json:"archives_original_min_margin"`
	ArchivesUnseenMeanMargin float64 `json:"archives_unseen_mean_margin"`
	ArchivesUnseenMinMargin  float64 `json:"archives_unseen_min_margin"`
	StoresOriginalAcc      float64 `json:"stores_original_accuracy"`
	StoresUnseenAcc        float64 `json:"stores_unseen_accuracy"`
	BaseOriginalAcc        float64 `json:"base_original_accuracy"`
	BaseUnseenAcc          float64 `json:"base_unseen_accuracy"`
	AcquiredAggregateAcc   float64 `json:"acquired_aggregate_accuracy"`
}

type UP123BArchivesOrthogonalResult struct {
	Schema             string        `json:"schema"`
	Experiment         string        `json:"experiment"`
	SourceUP122BSeal   string        `json:"source_up122b_seal"`
	StateDimension     int           `json:"state_dimension"`
	LearningRate       float64       `json:"learning_rate"`
	EpochsPerBatch     int           `json:"epochs_per_batch"`
	GroundingPerVerb   int           `json:"grounding_per_verb"`
	FixedBudget        int           `json:"fixed_budget"`
	StoresCorrectionFrozen bool      `json:"stores_correction_frozen"`
	ArchivesCorrectionFrozen bool    `json:"archives_correction_frozen"`
	AdaptiveGeometry   bool          `json:"adaptive_geometry"`
	Points             []UP123BPoint `json:"points"`
}

func up123bMeanRaw(verb string)[64]float64{
	var out [64]float64
	for _,name:=range up121bOriginalNames {
		h:=up95bEncode(name+" "+verb)
		for i:=0;i<64;i++ { out[i]+=h[i]/6.0 }
	}
	return out
}

func up123bMeanOf(verbs ...string)[64]float64{
	var out [64]float64
	for _,verb:=range verbs {
		v:=up123bMeanRaw(verb)
		for i:=0;i<64;i++ { out[i]+=v[i]/float64(len(verbs)) }
	}
	return out
}

func up123bArchivesDelta()[64]float64{
	s:=up123bMeanRaw("archives")
	o:=up123bMeanOf("notices","spots")
	r:=up123bMeanOf("recounts","retells")
	oo,rr,or:=up120bDot(o,o),up120bDot(r,r),up120bDot(o,r)
	so,sr:=up120bDot(s,o),up120bDot(s,r)
	det:=oo*rr-or*or
	a,b:=0.0,0.0
	if math.Abs(det)>1e-12 {
		a=(so*rr-sr*or)/det
		b=(sr*oo-so*or)/det
	}
	var residual [64]float64
	for i:=0;i<64;i++ { residual[i]=s[i]-a*o[i]-b*r[i] }
	ns,nr:=up120bNorm(s),up120bNorm(residual)
	if nr>1e-12 {
		scale:=ns/nr
		for i:=0;i<64;i++ { residual[i]*=scale }
	}
	var d [64]float64
	for i:=0;i<64;i++ { d[i]=residual[i]-s[i] }
	return d
}

func up123bEncode(arm,name,verb string,storesOrth,archivesDelta [64]float64)[64]float64{
	h:=up121bEncodeName("stores_competitor_orthogonal",name,verb,storesOrth)
	if arm=="stores_plus_archives"&&verb=="archives" {
		for i:=0;i<64;i++ { h[i]+=archivesDelta[i] }
	}
	return h
}

func up123bTrainStep(c *up97bClassifier,arm,name string,spec up106bVerbSpec,storesOrth,archivesDelta [64]float64){
	h:=up123bEncode(arm,name,spec.verb,storesOrth,archivesDelta)
	p:=c.probs(h)
	for k:=0;k<3;k++ {
		g:=p[k]
		if k==spec.class { g-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
		c.b[k]-=0.08*g
	}
}

func up123bTrainBase(arm string,storesOrth,archivesDelta [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++ {
		for _,name:=range up121bOriginalNames {
			for _,spec:=range up106bBase { up123bTrainStep(c,arm,name,spec,storesOrth,archivesDelta) }
		}
	}
	return c
}

func up123bAcquire(c *up97bClassifier,arm string,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,storesOrth,archivesDelta [64]float64){
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ { up123bTrainStep(c,arm,up121bOriginalNames[ni],current,storesOrth,archivesDelta) }
		for _,spec:=range up106bBase { up123bTrainStep(c,arm,up121bOriginalNames[0],spec,storesOrth,archivesDelta) }
		var seq []up109bReplayExample
		if policy=="stage3_anchor2"&&stage==3 {
			seq=up112bExtraReplay(previous,current.class,epoch,2)
		}else{
			seq=up111bCurrentExcluded(previous,current.class,epoch)
		}
		for _,x:=range seq { up123bTrainStep(c,arm,up121bOriginalNames[x.ni],x.spec,storesOrth,archivesDelta) }
	}
}

func up123bVerbAcc(c *up97bClassifier,arm string,spec up106bVerbSpec,names []string,storesOrth,archivesDelta [64]float64)float64{
	hits:=0
	for _,name:=range names {
		if up97bArgmax(c.probs(up123bEncode(arm,name,spec.verb,storesOrth,archivesDelta)))==spec.class { hits++ }
	}
	return float64(hits)/float64(len(names))
}

func up123bBaseAcc(c *up97bClassifier,arm string,names []string,storesOrth,archivesDelta [64]float64)float64{
	s:=0.0
	for _,spec:=range up106bBase { s+=up123bVerbAcc(c,arm,spec,names,storesOrth,archivesDelta) }
	return s/float64(len(up106bBase))
}

func up123bArchivesMargin(c *up97bClassifier,arm string,names []string,storesOrth,archivesDelta [64]float64)(float64,float64){
	sum:=0.0
	min:=math.Inf(1)
	for _,name:=range names {
		z:=up117bLogits(c,up123bEncode(arm,name,"archives",storesOrth,archivesDelta))
		bestOther:=math.Max(z[up97bObserve],z[up97bReport])
		m:=z[up97bStore]-bestOther
		sum+=m
		if m<min { min=m }
	}
	return sum/float64(len(names)),min
}

func up123bAcquiredAcc(c *up97bClassifier,arm string,acquired []up106bVerbSpec,storesOrth,archivesDelta [64]float64)float64{
	if len(acquired)==0{return 0}
	s:=0.0
	for _,spec:=range acquired { s+=up123bVerbAcc(c,arm,spec,up121bUnseenNames[4:6],storesOrth,archivesDelta) }
	return s/float64(len(acquired))
}

func up123bPoint(arm,order,policy string,stage int,c *up97bClassifier,acquired []up106bVerbSpec,storesOrth,archivesDelta [64]float64) UP123BPoint {
	archives:=up106bVerbSpec{verb:"archives",class:up97bStore}
	stores:=up106bVerbSpec{verb:"stores",class:up97bStore}
	om,omin:=up123bArchivesMargin(c,arm,up121bOriginalNames[4:6],storesOrth,archivesDelta)
	um,umin:=up123bArchivesMargin(c,arm,up121bUnseenNames[4:6],storesOrth,archivesDelta)
	return UP123BPoint{
		Arm:arm,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		ArchivesOriginalAcc:up123bVerbAcc(c,arm,archives,up121bOriginalNames[4:6],storesOrth,archivesDelta),
		ArchivesUnseenAcc:up123bVerbAcc(c,arm,archives,up121bUnseenNames[4:6],storesOrth,archivesDelta),
		ArchivesOrigMeanMargin:om,ArchivesOrigMinMargin:omin,
		ArchivesUnseenMeanMargin:um,ArchivesUnseenMinMargin:umin,
		StoresOriginalAcc:up123bVerbAcc(c,arm,stores,up121bOriginalNames,storesOrth,archivesDelta),
		StoresUnseenAcc:up123bVerbAcc(c,arm,stores,up121bUnseenNames,storesOrth,archivesDelta),
		BaseOriginalAcc:up123bBaseAcc(c,arm,up121bOriginalNames,storesOrth,archivesDelta),
		BaseUnseenAcc:up123bBaseAcc(c,arm,up121bUnseenNames,storesOrth,archivesDelta),
		AcquiredAggregateAcc:up123bAcquiredAcc(c,arm,acquired,storesOrth,archivesDelta),
	}
}

func RunUP123B()(UP123BArchivesOrthogonalResult,error){
	storesOrth:=up120bOrthogonalDelta()
	archivesDelta:=up123bArchivesDelta()
	result:=UP123BArchivesOrthogonalResult{
		Schema:UP123BArchivesOrthogonalSchema,Experiment:"UP-123B-archives-orthogonal",
		SourceUP122BSeal:"3381d0138ca59c23d64196124d1d6936523206ec",
		StateDimension:64,LearningRate:0.08,EpochsPerBatch:20,GroundingPerVerb:4,FixedBudget:6,
		StoresCorrectionFrozen:true,ArchivesCorrectionFrozen:true,AdaptiveGeometry:false,
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"stores_only","stores_plus_archives"} {
		for _,order:=range orders {
			name:=up114bOrderName(order)
			seq:=up121bSequence(order,"alternate")
			for _,policy:=range []string{"current_class_excluded","stage3_anchor2"} {
				c:=up123bTrainBase(arm,storesOrth,archivesDelta)
				acquired:=[]up106bVerbSpec{}
				result.Points=append(result.Points,up123bPoint(arm,name,policy,0,c,acquired,storesOrth,archivesDelta))
				for idx,spec:=range seq {
					stage:=idx+1
					up123bAcquire(c,arm,spec,acquired,policy,stage,storesOrth,archivesDelta)
					acquired=append(acquired,spec)
					result.Points=append(result.Points,up123bPoint(arm,name,policy,stage,c,acquired,storesOrth,archivesDelta))
				}
			}
		}
	}
	return result,nil
}
