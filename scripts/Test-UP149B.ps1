$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up149b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up149b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP149B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP149B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP149B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP149B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up149b-four-tail-causal-ranking -count=1;if($LASTEXITCODE-ne 0){throw 'UP149B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP149B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up149b-four-tail-causal-ranking)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP149B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up149b-four-tail-causal-ranking)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP149B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP149B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up149b-four-tail-causal-ranking.v1'){throw 'UP149B_SCHEMA_MISMATCH'}
 if([int]$R.arms.Count-ne 4 -or [int]$R.common_prefix_epochs-ne 15 -or [int]$R.terminal_epochs-ne 5){throw 'UP149B_DESIGN'}
 $Expected=@(0,6,12,18);for($i=0;$i-lt 4;$i++){if([int]$R.arms[$i].terminal_shift-ne $Expected[$i]){throw 'UP149B_TAIL_ORDER'}}
 if($R.adaptive_tail_selection -or $R.extra_updates_used){throw 'UP149B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in $R.arms){Write-Host "shift=$($A.terminal_shift) prefix=$($A.prefix_state_old_retention) tail=$($A.mean_terminal_tail_damage) old=$($A.final_mean_old_retention) new=$($A.final_mean_new_accuracy)"}
 Write-Host 'WINGLESS_UP149_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
