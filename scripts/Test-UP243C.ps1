$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UP243C_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-up243c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up243c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP243C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP243C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up243c-initial-state-failure-partition -count=1;if($LASTEXITCODE-ne 0){throw 'UP243C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP243C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up243c-initial-state-failure-partition)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP243C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up243c-initial-state-failure-partition)|Out-String).Trim();if($P1-cne $P2){throw 'UP243C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up243c-initial-state-failure-partition.v1' -or [int]$R.paired_initial_states-ne 64 -or [int]$R.cohort_summaries.Count-ne 4 -or [int]$R.hand_summaries.Count-ne 16 -or [int]$R.points.Count-ne 64){throw 'UP243C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP243C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.cohort_summaries){Write-Host "cohort=$($S.cohort) states=$($S.states) a_fail=$($S.a_failures) b_fail=$($S.b_failures) both_fail=$($S.both_fail) neither=$($S.neither_fail)"}
 foreach($S in $R.hand_summaries){if([int]$S.a_failures-gt 0 -or [int]$S.b_failures-gt 0){Write-Host "hand=$($S.initial_hand) states=$($S.states) a_fail=$($S.a_failures) b_fail=$($S.b_failures) both_fail=$($S.both_fail) neither=$($S.neither_fail)"}}
 Write-Host 'WINGLESS_UP355_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
