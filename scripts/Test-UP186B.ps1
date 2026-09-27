$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up186b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up186b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP186B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP186B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP186B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP186B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up186b-native-context-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP186B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP186B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up186b-native-context-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP186B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up186b-native-context-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP186B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP186B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up186b-native-context-diagnostic.v1'){throw 'UP186B_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 8 -or $R.fitting_used -or $R.maintenance_triggered){throw 'UP186B_DESIGN'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "phase=$($P.phase) crossings=$($P.actual_crossings) mean_abs=$($P.mean_abs_margin) near=$($P.near_boundary_count) correct=$($P.correct_count) norm=$($P.gate_weight_l2)"}
 Write-Host 'WINGLESS_UP186_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
