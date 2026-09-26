$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up165b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up165b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP165B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP165B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP165B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP165B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up165b-report-margin-crossings -count=1;if($LASTEXITCODE-ne 0){throw 'UP165B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP165B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up165b-report-margin-crossings)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP165B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up165b-report-margin-crossings)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP165B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP165B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up165b-report-margin-crossings.v1'){throw 'UP165B_SCHEMA_MISMATCH'}
 if([int]$R.steps.Count-ne 12 -or [int]$R.subjects-ne 6 -or [int]$R.report_steps-ne 2){throw 'UP165B_DESIGN'}
 if($R.decision_threshold-cne 'probability_margin_zero' -or $R.diagnostic_updates_retained -or $R.adaptive_ordering_used){throw 'UP165B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.steps){Write-Host "subject=$($S.subject) report=$($S.report_index) A_down=$($S.path_a_correct_to_incorrect) A_up=$($S.path_a_incorrect_to_correct) B_down=$($S.path_b_correct_to_incorrect) B_up=$($S.path_b_incorrect_to_correct) diff=$($S.different_correctness_after)"}
 Write-Host 'WINGLESS_UP165_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
