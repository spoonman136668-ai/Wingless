$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up209c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up209c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP209C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP209C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up209c-committed-tail-cap-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP209C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP209C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up209c-committed-tail-cap-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP209C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up209c-committed-tail-cap-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP209C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up209c-committed-tail-cap-transfer.v1' -or [int]$R.metrics.Count-ne 12 -or [int]$R.opportunity_caps.Count-ne 3){throw 'UP209C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_cap_selection_used){throw 'UP209C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) cap=$($M.opportunity_cap) baseline=$($M.baseline_losses) treated=$($M.treated_losses) prevented=$($M.prevented_losses) actions=$($M.actions_taken) accelerated=$($M.accelerated_losses) failure_extension=$($M.mean_loss_step_extension_among_failures)"}
 Write-Host 'WINGLESS_UP284_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
