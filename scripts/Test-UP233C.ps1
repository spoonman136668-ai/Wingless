$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up233c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up233c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP233C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP233C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up233c-nearest-signature-risk -count=1;if($LASTEXITCODE-ne 0){throw 'UP233C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP233C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up233c-nearest-signature-risk)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP233C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up233c-nearest-signature-risk)|Out-String).Trim();if($P1-cne $P2){throw 'UP233C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up233c-nearest-signature-risk.v1' -or [int]$R.arms_total-ne 512 -or [int]$R.arms.Count-ne 512){throw 'UP233C_DESIGN'}
 if($R.heldout_fitting_used -or $R.new_native_feature_used -or $R.adaptive_weighting_used -or $R.threshold_fitting_used -or $R.intervention_changed -or $R.live_activation){throw 'UP233C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "failures=$($R.failures) survivors=$($R.survivors) event_failures=$($R.event_bearing_failures) event_survivors=$($R.event_bearing_survivors)"
 Write-Host "high=$($R.high_risk) low=$($R.low_risk) unknown=$($R.unknown) no_event=$($R.no_event) high_fail=$($R.high_risk_failures) high_survive=$($R.high_risk_survivors) low_fail=$($R.low_risk_failures) recall=$($R.high_risk_failure_recall) precision=$($R.high_risk_precision)"
 Write-Host "failure_dists=$($R.mean_failure_min_failure_distance)/$($R.mean_failure_min_survivor_distance) survivor_dists=$($R.mean_survivor_min_failure_distance)/$($R.mean_survivor_min_survivor_distance)"
 Write-Host 'WINGLESS_UP332_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
