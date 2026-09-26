$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up159c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up159c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP159C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP159C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP159C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP159C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up159c-durability-horizon-replication -count=1;if($LASTEXITCODE-ne 0){throw 'UP159C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP159C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up159c-durability-horizon-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP159C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up159c-durability-horizon-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP159C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP159C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up159c-durability-horizon-replication.v1'){throw 'UP159C_SCHEMA_MISMATCH'}
 if([int]$R.durability_threshold-ne 21 -or $R.threshold_adaptive -or $R.future_schedule_used -or $R.semantic_class_used){throw 'UP159C_BOUNDARY_LEAK'}
 if(-not $R.counterfactual_only -or $R.live_activation){throw 'UP159C_ACTIVATION_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "arms=$($R.arms) permanent=$($R.permanent_saves) delayed=$($R.delayed_losses) tp=$($R.true_positive) fp=$($R.false_positive) fn=$($R.false_negative) tn=$($R.true_negative) auroc=$($R.horizon_auroc)"
 Write-Host 'WINGLESS_UP159_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
