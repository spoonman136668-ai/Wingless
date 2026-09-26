$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up160b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up160b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP160B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP160B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP160B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP160B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up160b-cleanup-noncommutativity -count=1;if($LASTEXITCODE-ne 0){throw 'UP160B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP160B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up160b-cleanup-noncommutativity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP160B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up160b-cleanup-noncommutativity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP160B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP160B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up160b-cleanup-noncommutativity.v1'){throw 'UP160B_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 18 -or [int]$R.subjects-ne 6 -or [int]$R.comparisons_per_subject-ne 3){throw 'UP160B_DESIGN'}
 if([int]$R.cleanup_updates_per_arm-ne 12){throw 'UP160B_CLEANUP_COUNT'}
 foreach($M in $R.metrics){if([int]$M.training_pair.Count-ne 2){throw 'UP160B_PAIR_COUNT'}}
 if($R.diagnostic_updates_retained -or $R.adaptive_ordering_used){throw 'UP160B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "subject=$($M.subject) comparison=$($M.comparison) retention_delta=$($M.retention_delta_a_minus_b) gate_distance=$($M.gate_distance_a_to_b)"}
 Write-Host 'WINGLESS_UP160_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
