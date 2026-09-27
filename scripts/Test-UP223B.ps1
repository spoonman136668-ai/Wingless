$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up223b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up223b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP223B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP223B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up223b-training-manifold-rejection -count=1;if($LASTEXITCODE-ne 0){throw 'UP223B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP223B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up223b-training-manifold-rejection)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP223B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up223b-training-manifold-rejection)|Out-String).Trim();if($P1-cne $P2){throw 'UP223B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up223b-training-manifold-rejection.v1' -or [int]$R.decisions.Count-ne 128 -or [int]$R.summaries.Count-ne 10){throw 'UP223B_DESIGN'}
 if($R.unseen_family_fitting_used -or $R.new_native_feature_used -or $R.evaluation_derived_threshold_used -or $R.adaptive_gate_selection_used -or $R.nonlinear_model_used -or $R.phase_input_used -or $R.parity_input_used -or $R.live_activation){throw 'UP223B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "train_correct=$($R.training_correct_states) margin=$($R.margin_threshold) centroid=$($R.centroid_distance_threshold) manifold=$($R.manifold_distance_threshold)"
 foreach($S in $R.summaries){Write-Host "gate=$($S.gate) scope=$($S.scope) family=$($S.family) decisions=$($S.decisions) harmful=$($S.harmful) harmful_rejected=$($S.harmful_rejected) harmful_accepted=$($S.harmful_accepted) benign_rejected=$($S.benign_rejected) harmful_rate=$($S.harmful_rejection_rate) benign_rate=$($S.benign_rejection_rate)"}
 Write-Host 'WINGLESS_UP333_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
