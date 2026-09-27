$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up169c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up169c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP169C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP169C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP169C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP169C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up169c-correction-delay -count=1;if($LASTEXITCODE-ne 0){throw 'UP169C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP169C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up169c-correction-delay)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP169C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up169c-correction-delay)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP169C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP169C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up169c-correction-delay.v1'){throw 'UP169C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 32 -or [int]$R.max_actions_per_arm-ne 1 -or [int]$R.evaluation_writes-ne 64){throw 'UP169C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.repeated_maintenance_used -or $R.adaptive_trigger_used -or $R.warning_threshold_changed){throw 'UP169C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) confirm=$($M.confirmation) intervention=$($M.intervention) actions=$($M.actions_taken) mean_delay=$($M.mean_delay_writes) delayed=$($M.arms_delayed) advanced=$($M.arms_advanced) max_delay=$($M.max_delay_writes) no_positive=$($M.no_positive_delay_actions)"}
 Write-Host 'WINGLESS_UP169_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
