$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up51c-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up51c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-51C POWERED DEPTH LADDER ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP51C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP51C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP51C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP51C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up51c-depth-ladder -count=1;if($LASTEXITCODE-ne 0){throw 'UP51C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP51C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up51c-depth-ladder)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP51C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up51c-depth-ladder)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP51C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP51C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up51c-powered-depth-ladder.v1'){throw 'UP51C_SCHEMA_MISMATCH'}
 if([int]$R.dimension-ne 128){throw 'UP51C_DIMENSION'}
 if([int]$R.metrics.Count-ne 5){throw 'UP51C_METRIC_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Last passing depth:  $($R.last_passing_depth)"
 Write-Host "First failing depth: $($R.first_failing_depth)"
 foreach($M in $R.metrics){Write-Host "depth=$($M.depth) norm=$($M.max_norm_drift) perturb=$($M.max_perturbation_gain) gate=$($M.gate)"}
 Write-Host 'WINGLESS_UP51_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
