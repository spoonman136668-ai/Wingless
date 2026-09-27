$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up227b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up227b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP227B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP227B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up227b-phase-periodicity-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP227B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP227B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up227b-phase-periodicity-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP227B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up227b-phase-periodicity-diagnostic)|Out-String).Trim();if($P1-cne $P2){throw 'UP227B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up227b-phase-periodicity-diagnostic.v1' -or [int]$R.training_phases.Count-ne 15 -or [int]$R.families.Count-ne 4 -or [int]$R.candidate_periods.Count-ne 8 -or [int]$R.summaries.Count-ne 32){throw 'UP227B_DESIGN'}
 if($R.evaluation_states_used -or $R.gate_changed -or $R.evaluation_derived_threshold_used -or $R.adaptive_period_selection_used -or $R.new_native_feature_used -or $R.live_activation){throw 'UP227B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($F in $R.families){$Rows=@($R.summaries|Where-Object{$_.family-eq $F});foreach($S in $Rows){Write-Host "family=$F period=$($S.period) matches=$($S.exact_geometry_matches)/$($S.comparisons)"}}
 Write-Host 'WINGLESS_UP349_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
