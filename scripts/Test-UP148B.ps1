$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up148b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up148b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP148B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP148B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP148B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP148B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up148b-terminal-tail-causality -count=1;if($LASTEXITCODE-ne 0){throw 'UP148B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP148B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up148b-terminal-tail-causality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP148B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up148b-terminal-tail-causality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP148B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP148B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up148b-terminal-tail-causality.v1'){throw 'UP148B_SCHEMA_MISMATCH'}
 if([int]$R.arms.Count-ne 2 -or [int]$R.common_prefix_epochs-ne 15 -or [int]$R.terminal_epochs-ne 5){throw 'UP148B_DESIGN'}
 if([int]$R.arms[0].terminal_shift-ne 18 -or [int]$R.arms[1].terminal_shift-ne 12){throw 'UP148B_ARM_ORDER'}
 if([math]::Abs([double]$R.arms[0].prefix_state_old_retention-[double]$R.arms[1].prefix_state_old_retention)-gt 1e-12){throw 'UP148B_PREFIX_MISMATCH'}
 if([int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15 -or [int]$R.final_new_updates-ne 4){throw 'UP148B_BUDGET'}
 if($R.adaptive_tail_selection -or $R.extra_updates_used){throw 'UP148B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in $R.arms){Write-Host "arm=$($A.arm) shift=$($A.terminal_shift) prefix_state=$($A.prefix_state_old_retention) tail=$($A.mean_terminal_tail_damage) net=$($A.mean_terminal_net_epoch_change) final_old=$($A.final_mean_old_retention) final_new=$($A.final_mean_new_accuracy)"}
 Write-Host "old_delta_safe_minus_damaging=$($R.final_old_retention_delta_safe_minus_damaging)"
 Write-Host "new_delta_safe_minus_damaging=$($R.final_new_accuracy_delta_safe_minus_damaging)"
 Write-Host 'WINGLESS_UP148_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
