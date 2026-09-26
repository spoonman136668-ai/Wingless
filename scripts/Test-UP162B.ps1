$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up162b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up162b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP162B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP162B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP162B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP162B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up162b-cleanup-geometry-sign -count=1;if($LASTEXITCODE-ne 0){throw 'UP162B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP162B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up162b-cleanup-geometry-sign)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP162B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up162b-cleanup-geometry-sign)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP162B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP162B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up162b-cleanup-geometry-sign.v1'){throw 'UP162B_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 6 -or [int]$R.subjects-ne 6){throw 'UP162B_DESIGN'}
 foreach($M in $R.metrics){if([int]$M.training_pair.Count-ne 2){throw 'UP162B_PAIR_COUNT'}}
 if($R.diagnostic_updates_retained -or $R.adaptive_ordering_used -or $R.post_result_classifier_used){throw 'UP162B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "subject=$($M.subject) parent_delta=$($M.parent_retention_delta_sor_minus_osr) align_margin=$($M.store_minus_observe_old_alignment) comm_old_cos=$($M.commutator_cosine_old_reference) comm_norm=$($M.commutator_norm)"}
 Write-Host 'WINGLESS_UP162_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
