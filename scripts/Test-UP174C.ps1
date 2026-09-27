$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up174c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up174c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP174C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP174C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP174C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP174C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up174c-warning-delay-tolerance -count=1;if($LASTEXITCODE-ne 0){throw 'UP174C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP174C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up174c-warning-delay-tolerance)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP174C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up174c-warning-delay-tolerance)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP174C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP174C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up174c-warning-delay-tolerance.v1'){throw 'UP174C_SCHEMA_MISMATCH'}
 if([int]$R.delay_intervals.Count-ne 3 -or [int]$R.metrics.Count-ne 3 -or [int]$R.max_actions_per_arm-ne 2){throw 'UP174C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_delay_used -or $R.adaptive_action_budget_used -or $R.warning_threshold_changed){throw 'UP174C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "delay=$($M.delay_intervals) actions=$($M.actions_taken) losses=$($M.losses) prevented=$($M.prevented_losses) efficiency=$($M.prevented_per_action) mean_loss_delay=$($M.mean_loss_delay_writes)"}
 Write-Host 'WINGLESS_UP174_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
