$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-up250c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up250c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP250C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP250C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up250c-five-native-coordinate-purity -count=1;if($LASTEXITCODE-ne 0){throw 'UP250C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP250C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up250c-five-native-coordinate-purity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP250C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up250c-five-native-coordinate-purity)|Out-String).Trim();if($P1-cne $P2){throw 'UP250C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up250c-five-native-coordinate-purity.v1' -or [int]$R.paired_initial_states-ne 64 -or [int]$R.coordinate_count-ne 9 -or [int]$R.parent_coordinate_quadruples-ne 126 -or [int]$R.coordinate_quintuples-ne 126 -or [int]$R.summaries.Count-ne 126){throw 'UP250C_DESIGN'}
 if([int]$R.parent_fully_pure_quadruples-ne 0){throw 'UP250C_PARENT_ANCHOR_DRIFT'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP250C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "parent_fully_pure_quadruples=$($R.parent_fully_pure_quadruples) fully_pure_quintuples=$($R.fully_pure_quintuples)"
 foreach($S in $R.summaries){if([int]$S.mixed_groups-eq 0){Write-Host "pure_five=$($S.coordinates -join '+') groups=$($S.groups)"}}
 Write-Host 'WINGLESS_UP375_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$PriorGOROOT;$env:GOCACHE=$PriorGoCache;$env:GOTMPDIR=$PriorGoTmp;Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
