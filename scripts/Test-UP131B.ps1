$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up131b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up131b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP131B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP131B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP131B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP131B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up131b-fewshot-lexical-grounding -count=1;if($LASTEXITCODE-ne 0){throw 'UP131B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP131B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up131b-fewshot-lexical-grounding)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP131B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up131b-fewshot-lexical-grounding)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP131B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP131B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up131b-fewshot-lexical-grounding.v1'){throw 'UP131B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128){throw 'UP131B_DIM'}
 if([int]$R.original_gate_epochs-ne 20 -or [int]$R.grounding_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP131B_TRAINING'}
 if($R.projector_recomputed -or $R.threshold_changed){throw 'UP131B_BOUNDARY'}
 if([int]$R.arm_metrics.Count-ne 4){throw 'UP131B_ARM_COUNT'}
 if([int]$R.surface_metrics.Count-ne 48){throw "UP131B_SURFACE_COUNT $($R.surface_metrics.Count)"}
 $Expected=@(0,1,2,4)
 for($i=0;$i-lt 4;$i++){if([int]$R.arm_metrics[$i].grounding_subjects-ne $Expected[$i]){throw 'UP131B_ARM_ORDER'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.arm_metrics){Write-Host "arm=$($M.arm) n=$($M.grounding_subjects) primary=$($M.primary_new.overall_accuracy) secondary=$($M.secondary_new.overall_accuracy) worst=$($M.primary_new.worst_surface_accuracy) oldH=$($M.old_heldout_retained.class_accuracy) oldU=$($M.old_unseen_retained.class_accuracy)"}
 Write-Host 'WINGLESS_UP172_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
