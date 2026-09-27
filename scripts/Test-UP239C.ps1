$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up239c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up239c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP239C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP239C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up239c-cross-schedule-signature-overlap -count=1;if($LASTEXITCODE-ne 0){throw 'UP239C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP239C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up239c-cross-schedule-signature-overlap)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP239C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up239c-cross-schedule-signature-overlap)|Out-String).Trim();if($P1-cne $P2){throw 'UP239C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up239c-cross-schedule-signature-overlap.v1' -or [int]$R.arms_total-ne 512 -or [int]$R.schedule_pairs-ne 6 -or [int]$R.signature_modes.Count-ne 2 -or [int]$R.summaries.Count-ne 12){throw 'UP239C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP239C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "pair=$($S.schedule_a)-$($S.schedule_b) mode=$($S.signature_mode) fail_shared=$($S.shared_failure_signatures) survive_shared=$($S.shared_survivor_signatures) af_bs_conflict=$($S.a_failure_b_survivor_conflicts) as_bf_conflict=$($S.a_survivor_b_failure_conflicts)"}
 Write-Host 'WINGLESS_UP344_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
