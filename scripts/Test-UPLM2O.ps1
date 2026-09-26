$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2o-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2o-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2O_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2O_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2O_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2O_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2o-unreported-eviction-warning -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2O_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2O_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2o-unreported-eviction-warning)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2O_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2o-unreported-eviction-warning)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2O_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2O_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2o-unreported-eviction-warning.v1'){throw 'UPLM2O_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 288 -or [int]$R.summaries.Count-ne 3 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM2O_DESIGN'}
 if($R.future_oracle_used -or $R.intervention_triggered){throw 'UPLM2O_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "tp=$($R.true_positive) fp=$($R.false_positive) fn=$($R.false_negative) tn=$($R.true_negative) precision=$($R.precision) recall=$($R.recall)"
 foreach($S in $R.summaries){Write-Host "D=$($S.deferred_first_chunk_reports) warnings=$($S.warnings) actual=$($S.actual_unreported_evictions)"}
 Write-Host 'WINGLESS_UP177_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
