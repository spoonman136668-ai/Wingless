$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up88a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up88a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP88A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP88A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP88A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP88A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up88a-delayed-role-recall -count=1;if($LASTEXITCODE-ne 0){throw 'UP88A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP88A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up88a-delayed-role-recall)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP88A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up88a-delayed-role-recall)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP88A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP88A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up88a-delayed-role-recall.v1'){throw 'UP88A_SCHEMA_MISMATCH'}
 if($R.retraining){throw 'UP88A_RETRAINING_PRESENT'}
 if([int]$R.points.Count-ne 32){throw "UP88A_POINT_COUNT $($R.points.Count)"}
 $Delays=@($R.points|ForEach-Object{[int]$_.delay}|Sort-Object -Unique)
 if(($Delays -join ',') -cne '16,64,256,1024'){throw "UP88A_DELAYS $($Delays -join ',')"}
 foreach($P in $R.points){if([int]$P.scenarios-ne 64){throw 'UP88A_SCENARIO_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "representation=$($P.representation) delay=$($P.delay) noise=$($P.noise_amplitude) target=$($P.target_recall_accuracy) exact=$($P.exact_state_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP88_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
