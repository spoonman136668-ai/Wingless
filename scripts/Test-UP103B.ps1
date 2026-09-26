$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up103b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up103b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP103B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP103B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP103B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP103B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up103b-rehearsal-grounding -count=1;if($LASTEXITCODE-ne 0){throw 'UP103B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP103B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up103b-rehearsal-grounding)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP103B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up103b-rehearsal-grounding)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP103B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP103B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up103b-rehearsal-grounding.v1'){throw 'UP103B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP103B_STATE_DIM'}
 if([int]$R.epochs-ne 20){throw 'UP103B_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UP103B_LR'}
 if([int]$R.grounding_per_verb-ne 4){throw 'UP103B_GROUNDING'}
 if([int]$R.points.Count-ne 2){throw 'UP103B_POINT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) heldout=$($M.heldout_accuracy) store_recall=$($M.store_recall) observe_recall=$($M.observe_recall) report_recall=$($M.report_recall) seen=$($M.seen_verb_accuracy)"}
 Write-Host 'WINGLESS_UP125_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
