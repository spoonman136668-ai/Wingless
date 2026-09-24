$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up68a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up68a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP68A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP68A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP68A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP68A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up68a-split-robustness -count=1;if($LASTEXITCODE-ne 0){throw 'UP68A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP68A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up68a-split-robustness)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP68A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up68a-split-robustness)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP68A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP68A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up68a-split-robustness.v1'){throw 'UP68A_SCHEMA_MISMATCH'}
 if($R.observer_class_changed){throw 'UP68A_OBSERVER_CHANGED'}
 if([int]$R.points.Count-ne 6){throw 'UP68A_POINT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "residue=$($P.residue) representation=$($P.representation) train=$($P.mean_train_accuracy) held=$($P.mean_heldout_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP68_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
