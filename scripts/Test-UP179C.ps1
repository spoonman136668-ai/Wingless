$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up179c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up179c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP179C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP179C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP179C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP179C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up179c-stage-target-ablation -count=1;if($LASTEXITCODE-ne 0){throw 'UP179C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP179C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up179c-stage-target-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP179C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up179c-stage-target-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP179C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP179C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up179c-stage-target-ablation.v1'){throw 'UP179C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 8 -or [int]$R.target_arms.Count-ne 4 -or [int]$R.max_actions-ne 2){throw 'UP179C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_target_used -or $R.adaptive_timing_used){throw 'UP179C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) arm=$($M.arm) stage1=$($M.stage_one_actions) stage2=$($M.stage_two_actions) losses=$($M.treated_losses) prevented=$($M.prevented_losses) delay=$($M.mean_loss_delay_writes)"}
 Write-Host 'WINGLESS_UP209_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
