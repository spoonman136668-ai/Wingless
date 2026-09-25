$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up87a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up87a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP87A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP87A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP87A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP87A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up87a-sequential-role-overwrite -count=1;if($LASTEXITCODE-ne 0){throw 'UP87A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP87A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up87a-sequential-role-overwrite)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP87A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up87a-sequential-role-overwrite)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP87A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP87A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up87a-sequential-role-overwrite.v1'){throw 'UP87A_SCHEMA_MISMATCH'}
 if($R.retraining){throw 'UP87A_RETRAINING_PRESENT'}
 if([int]$R.points.Count-ne 24){throw "UP87A_POINT_COUNT $($R.points.Count)"}
 $Lengths=@($R.points|ForEach-Object{[int]$_.length}|Sort-Object -Unique)
 if(($Lengths -join ',') -cne '16,64,256'){throw "UP87A_LENGTHS $($Lengths -join ',')"}
 foreach($P in $R.points){if([int]$P.scenarios-ne 64){throw 'UP87A_SCENARIO_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "representation=$($P.representation) length=$($P.length) noise=$($P.noise_amplitude) value=$($P.value_accuracy) exact=$($P.exact_state_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP87_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
