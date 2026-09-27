$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up194c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up194c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP194C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP194C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up194c-hostile-spacing-allocation -count=1;if($LASTEXITCODE-ne 0){throw 'UP194C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP194C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up194c-hostile-spacing-allocation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP194C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up194c-hostile-spacing-allocation)|Out-String).Trim();if($P1-cne $P2){throw 'UP194C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up194c-hostile-spacing-allocation.v1' -or [int]$R.metrics.Count-ne 5 -or [int]$R.action_cap-ne 16){throw 'UP194C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.warning_threshold_changed -or $R.adaptive_spacing_used){throw 'UP194C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) spacing=$($M.spacing_intervals) cap=$($M.action_cap) baseline=$($M.baseline_losses) prevented=$($M.prevented_losses) actions=$($M.actions_taken) accelerated=$($M.accelerated_losses) mean_extension=$($M.mean_loss_step_extension_among_failures) max_extension=$($M.max_loss_step_extension_among_failures)"}
 Write-Host 'WINGLESS_UP249_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
