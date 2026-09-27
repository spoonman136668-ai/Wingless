$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UP244C_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-up244c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up244c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP244C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP244C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up244c-endangered-key-failure-partition -count=1;if($LASTEXITCODE-ne 0){throw 'UP244C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP244C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up244c-endangered-key-failure-partition)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP244C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up244c-endangered-key-failure-partition)|Out-String).Trim();if($P1-cne $P2){throw 'UP244C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up244c-endangered-key-failure-partition.v1' -or [int]$R.paired_initial_states-ne 64 -or [int]$R.key_summaries.Count-lt 1){throw 'UP244C_DESIGN'}
 $Total=0;foreach($S in $R.key_summaries){$Total+=[int]$S.states};if($Total-ne 64){throw 'UP244C_STATE_ACCOUNTING'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP244C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.key_summaries){Write-Host "key=$($S.endangered_key) states=$($S.states) a_fail=$($S.a_failures) b_fail=$($S.b_failures) both_fail=$($S.both_fail) neither=$($S.neither_fail) a_events=$($S.a_events) b_events=$($S.b_events)"}
 Write-Host 'WINGLESS_UP358_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
