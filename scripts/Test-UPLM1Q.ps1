$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1q-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1q-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1Q_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1Q_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1Q_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1Q_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1q-fivefamily-cyclic-position -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1Q_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1Q_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1q-fivefamily-cyclic-position)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1Q_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1q-fivefamily-cyclic-position)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1Q_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1Q_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1q-fivefamily-cyclic-position.v1'){throw 'UPLM1Q_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM1Q_CONFIG'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM1Q_LR'}
 if($R.recurrent_parameters_trained -or $R.router_used -or $R.exact_recall_used -or $R.adaptive_ordering_used -or $R.sixth_family_used){throw 'UPLM1Q_BOUNDARY'}
 if([int]$R.metrics.Count-ne 25){throw "UPLM1Q_METRIC_COUNT $($R.metrics.Count)"}
 if([int]$R.summaries.Count-ne 5){throw "UPLM1Q_SUMMARY_COUNT $($R.summaries.Count)"}\n $ExpectedArms=@("rotate_0","rotate_1","rotate_2","rotate_3","rotate_4")\n for($i=0;$i-lt 5;$i++){if($R.summaries[$i].arm-cne $ExpectedArms[$i]){throw "UPLM1Q_ARM_LABEL $($R.summaries[$i].arm)"}}\n if((@($R.summaries.arm | Sort-Object -Unique)).Count-ne 5){throw "UPLM1Q_DUPLICATE_ARM_LABEL"}
 foreach($M in $R.metrics){if([int]$M.update_position-lt 0 -or [int]$M.update_position-gt 4){throw 'UPLM1Q_POSITION'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.summaries){Write-Host "arm=$($M.arm) min=$($M.min_heldout_accuracy) mean=$($M.mean_heldout_accuracy) spread=$($M.heldout_accuracy_spread) min_family=$($M.min_family) max_family=$($M.max_family)"}
 Write-Host 'WINGLESS_UP182_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
