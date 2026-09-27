$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up178c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up178c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP178C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP178C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP178C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP178C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up178c-warning-stage-ablation -count=1;if($LASTEXITCODE-ne 0){throw 'UP178C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP178C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up178c-warning-stage-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP178C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up178c-warning-stage-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP178C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP178C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up178c-warning-stage-ablation.v1'){throw 'UP178C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 8 -or [int]$R.arms.Count-ne 4 -or [int]$R.delay_intervals-ne 1){throw 'UP178C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_stage_used -or $R.adaptive_delay_used){throw 'UP178C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) arm=$($M.arm) actions=$($M.actions_taken) losses=$($M.treated_losses) prevented=$($M.prevented_losses) efficiency=$($M.prevented_per_action) delay=$($M.mean_loss_delay_writes)"}
 Write-Host 'WINGLESS_UP206_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
