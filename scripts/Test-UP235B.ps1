$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-up235b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up235b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP235B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP235B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up235b-centered-parity-feature-subsets -count=1;if($LASTEXITCODE-ne 0){throw 'UP235B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP235B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up235b-centered-parity-feature-subsets)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP235B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up235b-centered-parity-feature-subsets)|Out-String).Trim();if($P1-cne $P2){throw 'UP235B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up235b-centered-parity-feature-subsets.v1' -or [int]$R.subset_count-ne 15 -or [int]$R.feature_names.Count-ne 4){throw 'UP235B_DESIGN'}
 if([double]$R.parent_family_centered_pooled_accuracy-ne 1.0){throw 'UP235B_PARENT_ANCHOR_DRIFT'}
 if(-not $R.training_parity_labels_used -or -not $R.family_centering_used -or $R.phase_input_at_inference_used -or $R.evaluation_label_fitting_used -or $R.adaptive_feature_selection_used -or $R.nonlinear_classifier_used -or $R.live_activation){throw 'UP235B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "minimum_perfect_cardinality=$($R.minimum_perfect_cardinality) minimum_perfect_masks=$($R.minimum_perfect_masks -join ',')"
 foreach($S in $R.subsets){Write-Host "mask=$($S.mask) k=$($S.cardinality) features=$($S.features -join '+') accuracy=$($S.accuracy)"}
 Write-Host 'WINGLESS_UP374_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$PriorGOROOT;$env:GOCACHE=$PriorGoCache;$env:GOTMPDIR=$PriorGoTmp;Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
