$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3w-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3w-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3W_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3W_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3w-full-deployment-vector -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3W_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3W_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3w-full-deployment-vector)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3W_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3w-full-deployment-vector)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM3W_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3w-full-deployment-vector.v1' -or [int]$R.resource_configurations-ne 80 -or [int]$R.metrics.Count-ne 80 -or [int]$R.summaries.Count-ne 2){throw 'UPLM3W_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_coordinate_search_used){throw 'UPLM3W_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "global_range=$($R.global_outcome_range) distinct=$($R.distinct_outcomes)"
 foreach($S in $R.summaries){Write-Host "coordinate=$($S.coordinate) groups=$($S.groups) nonzero=$($S.nonzero_spread_groups) max_spread=$($S.max_failure_spread) mean_spread=$($S.mean_failure_spread)"}
 Write-Host 'WINGLESS_UP235_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
