$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2n-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2n-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2N_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2N_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2N_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2N_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2n-recall-loss-lead-time -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2N_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2N_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2n-recall-loss-lead-time)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2N_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2n-recall-loss-lead-time)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2N_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2N_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2n-recall-loss-lead-time.v1'){throw 'UPLM2N_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 16 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM2N_DESIGN'}
 if($R.model_training_used -or $R.intervention_triggered -or $R.future_oracle_used){throw 'UPLM2N_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "rot=$($P.identity_rotation) order=$($P.deferred_report_order) shift=$($P.value_shift) pos=$($P.lost_pending_position) lead_events=$($P.lead_events)"}
 Write-Host 'WINGLESS_UP176_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
