$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source
$GoRoot=Split-Path (Split-Path $GoExe -Parent) -Parent
$GoCache=Join-Path $env:TEMP ("wingless-up239b-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up239b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache,$GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne0){throw 'UP239B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne0){throw 'UP239B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up239b-heldout-composition-parity -count=1
 if($LASTEXITCODE-ne0){throw 'UP239B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1
 if($LASTEXITCODE-ne0){throw 'UP239B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up239b-heldout-composition-parity)|Out-String).Trim()
 if($LASTEXITCODE-ne0){throw 'UP239B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up239b-heldout-composition-parity)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP239B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up239b-heldout-composition-parity.v1' -or [int]$R.heldout_families.Count-ne 4 -or [int]$R.summaries.Count-ne 4){throw 'UP239B_DESIGN'}
 if($R.parent_classification-cne 'EXTREME_TRANSFER_PERFECT' -or [double]$R.parent_extreme_accuracy-ne 1.0){throw 'UP239B_PARENT_ANCHOR_DRIFT'}
 if(-not $R.unlabeled_family_centering_used -or $R.evaluation_label_fitting_used -or $R.phase_input_at_inference_used -or $R.adaptive_feature_selection_used -or $R.nonlinear_classifier_used -or $R.live_activation){throw 'UP239B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "far_accuracy=$($R.far_accuracy) extreme_accuracy=$($R.extreme_accuracy) classification=$($R.classification)"
 Write-Host 'WINGLESS_UP239B_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq$PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq$PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq$PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
