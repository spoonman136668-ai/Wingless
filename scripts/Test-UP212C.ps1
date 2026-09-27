$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up212c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up212c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP212C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP212C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up212c-dense-tail-cap-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UP212C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP212C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up212c-dense-tail-cap-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP212C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up212c-dense-tail-cap-boundary)|Out-String).Trim();if($P1-cne $P2){throw 'UP212C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up212c-dense-tail-cap-boundary.v1' -or [int]$R.metrics.Count-ne 80 -or [int]$R.horizons.Count-ne 5){throw 'UP212C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_selection_used){throw 'UP212C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) horizon=$($M.horizon) cap=$($M.opportunity_cap) baseline=$($M.baseline_losses) treated=$($M.treated_losses) prevented=$($M.prevented_losses) actions=$($M.actions_taken)"}
 Write-Host 'WINGLESS_UP291_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
