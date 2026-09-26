$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up157b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up157b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP157B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP157B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP157B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP157B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up157b-late-rehearsal-position-map -count=1;if($LASTEXITCODE-ne 0){throw 'UP157B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP157B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up157b-late-rehearsal-position-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP157B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up157b-late-rehearsal-position-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP157B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP157B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up157b-late-rehearsal-position-map.v1'){throw 'UP157B_SCHEMA_MISMATCH'}
 if([int]$R.arms.Count-ne 30 -or [int]$R.rotations-ne 15 -or [int]$R.subjects-ne 2){throw 'UP157B_DESIGN'}
 if([int]$R.rehearsal_after-ne 12 -or [int]$R.rehearsal_before-ne 3 -or [int]$R.old_updates_per_epoch-ne 15 -or [int]$R.new_updates_per_epoch-ne 24){throw 'UP157B_BUDGET'}
 foreach($A in $R.arms){if([int]$A.pre_report_original_indices.Count-ne 3 -or [int]$A.post_report_original_indices.Count-ne 12){throw 'UP157B_POSITION_COUNT'}}
 if($R.extra_updates_used -or $R.adaptive_ordering_used){throw 'UP157B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in $R.arms){Write-Host "subject=$($A.subject) rotation=$($A.rotation) post_recovery=$($A.mean_post_report_rehearsal_recovery) old=$($A.final_mean_old_retention) new=$($A.final_mean_new_accuracy)"}
 Write-Host 'WINGLESS_UP157_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
