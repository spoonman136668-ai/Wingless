$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UP247C_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-up247c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up247c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP247C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP247C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up247c-pairwise-native-coordinate-purity -count=1;if($LASTEXITCODE-ne 0){throw 'UP247C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP247C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up247c-pairwise-native-coordinate-purity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP247C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up247c-pairwise-native-coordinate-purity)|Out-String).Trim();if($P1-cne $P2){throw 'UP247C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up247c-pairwise-native-coordinate-purity.v1' -or [int]$R.paired_initial_states-ne 64 -or [int]$R.coordinates.Count-ne 9 -or [int]$R.coordinate_pairs-ne 36 -or [int]$R.summaries.Count-ne 36){throw 'UP247C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP247C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "pair=$($S.coordinate_a)+$($S.coordinate_b) groups=$($S.groups) pure=$($S.pure_groups) mixed=$($S.mixed_groups) max_classes=$($S.max_classes_per_group)"}
 Write-Host 'WINGLESS_UP366_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
