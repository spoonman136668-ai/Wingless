$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up181c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up181c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP181C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP181C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP181C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP181C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up181c-stage-two-delay-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UP181C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP181C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up181c-stage-two-delay-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP181C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up181c-stage-two-delay-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP181C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP181C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up181c-stage-two-delay-boundary.v1'){throw 'UP181C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 6 -or [int]$R.stage_two_delays.Count-ne 3 -or [int]$R.max_actions-ne 2){throw 'UP181C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_delay_used -or $R.adaptive_target_used){throw 'UP181C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) d2=$($M.stage_two_delay) stage1=$($M.stage_one_actions) stage2=$($M.stage_two_actions) losses=$($M.treated_losses) prevented=$($M.prevented_losses) efficiency=$($M.prevented_per_action) delay=$($M.mean_loss_delay_writes)"}
 Write-Host 'WINGLESS_UP215_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
