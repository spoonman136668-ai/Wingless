$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2l-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2l-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2L_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2L_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2L_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2L_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2l-five-dependency-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2L_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2L_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2l-five-dependency-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2L_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2l-five-dependency-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2L_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2L_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2l-five-dependency-boundary.v1'){throw 'UPLM2L_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 600 -or [int]$R.deferred_levels.Count-ne 3 -or [int]$R.identity_rotations.Count-ne 2 -or [int]$R.deferred_report_orders.Count-ne 2){throw 'UPLM2L_DESIGN'}
 if([int]$R.exact_recall_cap-ne 16 -or -not $R.event_multiset_identical -or $R.capacity_changed -or $R.adaptive_chunking_used -or $R.extra_training_used){throw 'UPLM2L_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "rot=$($M.identity_rotation) order=$($M.deferred_report_order) defer=$($M.deferred_first_chunk_reports) recall=$($M.recall_hit_rate) expected=$($M.expected_recall_hit_rate) exact=$($M.report_set_exact_accuracy)"}
 Write-Host 'WINGLESS_UP174_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
