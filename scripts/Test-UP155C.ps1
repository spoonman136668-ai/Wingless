$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up155c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up155c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP155C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP155C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP155C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP155C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up155c-protected-cohort-warning -count=1;if($LASTEXITCODE-ne 0){throw 'UP155C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP155C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up155c-protected-cohort-warning)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP155C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up155c-protected-cohort-warning)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP155C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP155C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up155c-protected-cohort-warning.v1'){throw 'UP155C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 4 -or [int]$R.cohorts-ne 4){throw 'UP155C_DESIGN'}
 foreach($M in $R.metrics){if(([int]$M.true_positive+[int]$M.false_positive+[int]$M.false_negative+[int]$M.true_negative)-ne 96){throw 'UP155C_METRIC_COUNT'}}
 if($R.predictor_uses_semantic_class -or $R.future_oracle_used -or $R.intervention_triggered){throw 'UP155C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cohort=$($M.cohort) tp=$($M.true_positive) fp=$($M.false_positive) fn=$($M.false_negative) precision=$($M.precision) recall=$($M.recall)"}
 Write-Host 'WINGLESS_UP155_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
