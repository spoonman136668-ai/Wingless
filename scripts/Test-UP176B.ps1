$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up176b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up176b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP176B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP176B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP176B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP176B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up176b-temporal-predictor-ablation -count=1;if($LASTEXITCODE-ne 0){throw 'UP176B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP176B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up176b-temporal-predictor-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP176B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up176b-temporal-predictor-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP176B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP176B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up176b-temporal-predictor-ablation.v1'){throw 'UP176B_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 2 -or [int]$R.total_slots-ne 17280 -or [int]$R.context_phase-ne 21){throw 'UP176B_DESIGN'}
 if(-not $R.predictor_probabilities_frozen_before_run -or $R.adaptive_threshold_used -or $R.maintenance_triggered){throw 'UP176B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "predictor=$($M.predictor) brier=$($M.brier_score) precision=$($M.warning_precision) recall=$($M.warning_recall)"}
 Write-Host 'WINGLESS_UP176_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
