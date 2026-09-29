$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UP249C_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-up249c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up249c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP249C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP249C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up249c-quad-native-coordinate-purity -count=1;if($LASTEXITCODE-ne 0){throw 'UP249C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP249C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up249c-quad-native-coordinate-purity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP249C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up249c-quad-native-coordinate-purity)|Out-String).Trim();if($P1-cne $P2){throw 'UP249C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up249c-quad-native-coordinate-purity.v1' -or [int]$R.paired_initial_states-ne 64 -or [int]$R.coordinates.Count-ne 9 -or [int]$R.parent_coordinate_triples-ne 84 -or [int]$R.coordinate_quadruples-ne 126 -or [int]$R.summaries.Count-ne 126){throw 'UP249C_DESIGN'}
 if([int]$R.parent_fully_pure_triples-ne 0){throw 'UP249C_PARENT_ANCHOR_DRIFT'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP249C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "parent_fully_pure_triples=$($R.parent_fully_pure_triples) fully_pure_quadruples=$($R.fully_pure_quadruples)"
 foreach($S in $R.summaries){if([int]$S.mixed_groups-eq 0){Write-Host "pure_quad=$($S.coordinate_a)+$($S.coordinate_b)+$($S.coordinate_c)+$($S.coordinate_d) groups=$($S.groups)"}}
 Write-Host 'WINGLESS_UP372_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
