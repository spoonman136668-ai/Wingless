$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-up236b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up236b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP236B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP236B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up236b-nearzero-centering-necessity -count=1;if($LASTEXITCODE-ne 0){throw 'UP236B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP236B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up236b-nearzero-centering-necessity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP236B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up236b-nearzero-centering-necessity)|Out-String).Trim();if($P1-cne $P2){throw 'UP236B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up236b-nearzero-centering-necessity.v1' -or [int]$R.parent_minimum_perfect_cardinality-ne 1 -or [int]$R.parent_minimum_perfect_masks.Count-ne 1 -or [int]$R.parent_minimum_perfect_masks[0]-ne 4 -or $R.feature-cne 'near_zero_margin_count' -or [int]$R.summaries.Count-ne 4){throw 'UP236B_DESIGN'}
 if([double]$R.centered_accuracy-ne 1.0){throw 'UP236B_PARENT_ANCHOR_DRIFT'}
 if($R.phase_input_at_inference_used -or $R.evaluation_label_fitting_used -or $R.adaptive_feature_selection_used -or $R.nonlinear_classifier_used -or $R.live_activation){throw 'UP236B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "raw_accuracy=$($R.raw_accuracy) centered_accuracy=$($R.centered_accuracy) classification=$($R.classification)"
 foreach($S in $R.summaries){Write-Host "family=$($S.family) raw=$($S.raw_correct)/$($S.states) centered=$($S.centered_correct)/$($S.states)"}
 Write-Host 'WINGLESS_UP378_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$PriorGOROOT;$env:GOCACHE=$PriorGoCache;$env:GOTMPDIR=$PriorGoTmp;Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
