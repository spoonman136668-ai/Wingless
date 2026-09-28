$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UP232B_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-up232b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up232b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP232B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP232B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up232b-parity-necessity-control -count=1;if($LASTEXITCODE-ne 0){throw 'UP232B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP232B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up232b-parity-necessity-control)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP232B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up232b-parity-necessity-control)|Out-String).Trim();if($P1-cne $P2){throw 'UP232B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up232b-parity-necessity-control.v1' -or [int]$R.training_phases.Count-ne 15 -or [int]$R.evaluation_phases.Count-ne 32 -or [int]$R.training_points-ne 60 -or [int]$R.evaluation_points-ne 128 -or [int]$R.summaries.Count-ne 4){throw 'UP232B_DESIGN'}
 if($R.evaluation_label_fitting_used -or $R.adaptive_feature_selection_used -or -not $R.parity_input_used -or -not $R.wrong_parity_control_used -or $R.nonlinear_classifier_used -or $R.live_activation){throw 'UP232B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "baseline_accuracy=$($R.baseline_accuracy) true_parity_accuracy=$($R.true_parity_accuracy) wrong_parity_accuracy=$($R.wrong_parity_accuracy)"
 foreach($S in $R.summaries){Write-Host "family=$($S.family) eval=$($S.evaluation_states) baseline=$($S.baseline_correct) true_parity=$($S.true_parity_correct) wrong_parity=$($S.wrong_parity_correct)"}
 Write-Host 'WINGLESS_UP365_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
