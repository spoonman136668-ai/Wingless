$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up164b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up164b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP164B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP164B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP164B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP164B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up164b-report-path-sensitivity -count=1;if($LASTEXITCODE-ne 0){throw 'UP164B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP164B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up164b-report-path-sensitivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP164B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up164b-report-path-sensitivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP164B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP164B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up164b-report-path-sensitivity.v1'){throw 'UP164B_SCHEMA_MISMATCH'}
 if([int]$R.subjects_data.Count-ne 6 -or [int]$R.subjects-ne 6 -or [int]$R.report_steps_per_path-ne 2){throw 'UP164B_DESIGN'}
 foreach($S in $R.subjects_data){if([int]$S.path_a_report_steps.Count-ne 2 -or [int]$S.path_b_report_steps.Count-ne 2){throw 'UP164B_STEP_COUNT'}}
 if($R.diagnostic_updates_retained -or $R.adaptive_ordering_used){throw 'UP164B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.subjects_data){Write-Host "subject=$($S.subject) delta_pre=$($S.pre_report_retention_delta_a_minus_b) delta_r13=$($S.after_report13_retention_delta_a_minus_b) delta_r14=$($S.after_report14_retention_delta_a_minus_b)"}
 Write-Host 'WINGLESS_UP164_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
