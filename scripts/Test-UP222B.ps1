$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up222b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up222b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP222B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP222B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up222b-absolute-distance-ood-rejection -count=1;if($LASTEXITCODE-ne 0){throw 'UP222B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP222B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up222b-absolute-distance-ood-rejection)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP222B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up222b-absolute-distance-ood-rejection)|Out-String).Trim();if($P1-cne $P2){throw 'UP222B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up222b-absolute-distance-ood-rejection.v1' -or [int]$R.decisions.Count-ne 128 -or [int]$R.summaries.Count-ne 15){throw 'UP222B_DESIGN'}
 if($R.unseen_family_fitting_used -or $R.new_centroid_used -or $R.new_calibrator_used -or $R.evaluation_derived_threshold_used -or $R.adaptive_gate_selection_used -or $R.phase_input_used -or $R.parity_input_used -or $R.nonlinear_model_used){throw 'UP222B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "train_correct=$($R.training_correct_states) margin_threshold=$($R.training_margin_threshold) distance_threshold=$($R.training_distance_threshold)"
 foreach($S in $R.summaries){Write-Host "gate=$($S.gate) scope=$($S.scope) family=$($S.family) decisions=$($S.decisions) rejected=$($S.rejected) harmful=$($S.harmful_forced_routes) harmful_rejected=$($S.harmful_rejected) harmful_accepted=$($S.harmful_accepted) benign_rejected=$($S.benign_rejected) harmful_rejection_rate=$($S.harmful_route_rejection_rate) benign_rejection_rate=$($S.benign_rejection_rate)"}
 Write-Host 'WINGLESS_UP325_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
