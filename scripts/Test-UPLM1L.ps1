$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1l-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1l-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1L_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1L_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1L_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1L_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1l-states-store-orthogonal -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1L_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1L_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1l-states-store-orthogonal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1L_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1l-states-store-orthogonal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1L_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1L_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1l-states-store-orthogonal.v1'){throw 'UPLM1L_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.fourth_grounding_epochs-ne 20){throw 'UPLM1L_CONFIG'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM1L_LR'}
 if($R.correction_recomputed -or $R.threshold_changed -or $R.byte_model_used){throw 'UPLM1L_BOUNDARY'}
 if([int]$R.family_metrics.Count-ne 12){throw "UPLM1L_FAMILY_COUNT $($R.family_metrics.Count)"}
 if([int]$R.surface_metrics.Count-ne 36){throw "UPLM1L_SURFACE_COUNT $($R.surface_metrics.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.family_metrics){if($M.family-ceq 'fourth'){Write-Host "arm=$($M.arm) split=$($M.split) acc=$($M.accuracy) store=$($M.store_recall) observe=$($M.observe_recall) report=$($M.report_recall)"}}
 foreach($M in $R.surface_metrics){if($M.surface-ceq 'states'){Write-Host "arm=$($M.arm) split=$($M.split) acc=$($M.accuracy) store_rate=$($M.predicted_store_rate) report_rate=$($M.predicted_report_rate) store=$($M.mean_store_probability) report=$($M.mean_report_probability) margin=$($M.min_target_margin)"}}
 Write-Host 'WINGLESS_UP168_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
