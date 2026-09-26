$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1k-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1k-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1K_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1K_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1K_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1K_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1k-fourth-router-topology -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1K_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1K_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1k-fourth-router-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1K_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1k-fourth-router-topology)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1K_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1K_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1k-fourth-router-topology.v1'){throw 'UPLM1K_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.fourth_grounding_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UPLM1K_PARAMS'}
 if(-not $R.router_reconstructed -or $R.router_modified_after_build -or $R.byte_model_used -or $R.threshold_changed){throw 'UPLM1K_DIAGNOSTIC_BOUNDARY'}
 if([int]$R.family_metrics.Count-ne 6){throw "UPLM1K_FAMILY_COUNT $($R.family_metrics.Count)"}
 if([int]$R.surface_metrics.Count-ne 18){throw "UPLM1K_SURFACE_COUNT $($R.surface_metrics.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.family_metrics){Write-Host "family=$($M.family) split=$($M.split) acc=$($M.accuracy) store=$($M.store_recall) observe=$($M.observe_recall) report=$($M.report_recall)"}
 foreach($M in $R.surface_metrics){if($M.accuracy-lt 1){Write-Host "surface family=$($M.family) split=$($M.split) surface=$($M.surface) target=$($M.target_class) acc=$($M.accuracy) store=$($M.mean_store_probability) observe=$($M.mean_observe_probability) report=$($M.mean_report_probability) margin=$($M.min_target_margin)"}}
 Write-Host 'WINGLESS_UP174_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
