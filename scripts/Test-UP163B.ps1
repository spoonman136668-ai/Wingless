$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up163b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up163b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP163B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP163B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP163B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP163B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up163b-cleanup-path-curvature -count=1;if($LASTEXITCODE-ne 0){throw 'UP163B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP163B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up163b-cleanup-path-curvature)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP163B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up163b-cleanup-path-curvature)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP163B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP163B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up163b-cleanup-path-curvature.v1'){throw 'UP163B_SCHEMA_MISMATCH'}
 if([int]$R.subjects_data.Count-ne 6 -or [int]$R.subjects-ne 6 -or [int]$R.checkpoints_per_path-ne 4){throw 'UP163B_DESIGN'}
 foreach($S in $R.subjects_data){
  if([int]$S.path_a_checkpoints.Count-ne 4 -or [int]$S.path_b_checkpoints.Count-ne 4 -or [int]$S.matched_gate_distance.Count-ne 3 -or [int]$S.matched_retention_delta_a_minus_b.Count-ne 3){throw 'UP163B_CHECKPOINT_COUNT'}
 }
 if($R.diagnostic_updates_retained -or $R.adaptive_ordering_used){throw 'UP163B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.subjects_data){Write-Host "subject=$($S.subject) deltas=$($S.matched_retention_delta_a_minus_b -join ',') distances=$($S.matched_gate_distance -join ',')"}
 Write-Host 'WINGLESS_UP163_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
