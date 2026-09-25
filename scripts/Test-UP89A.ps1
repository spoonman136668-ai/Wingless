$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up89a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up89a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP89A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP89A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP89A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP89A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up89a-multiquery-latest-value -count=1;if($LASTEXITCODE-ne 0){throw 'UP89A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP89A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up89a-multiquery-latest-value)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP89A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up89a-multiquery-latest-value)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP89A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP89A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up89a-multiquery-latest-value.v1'){throw 'UP89A_SCHEMA_MISMATCH'}
 if($R.retraining){throw 'UP89A_RETRAINING_PRESENT'}
 if($R.exact_recall_side_channel){throw 'UP89A_EXACT_RECALL_SIDE_CHANNEL_PRESENT'}
 if([int]$R.points.Count-ne 18){throw "UP89A_POINT_COUNT $($R.points.Count)"}
 $Lengths=@($R.points|ForEach-Object{[int]$_.length}|Sort-Object -Unique)
 if(($Lengths -join ',') -cne '64,256,1024'){throw "UP89A_LENGTHS $($Lengths -join ',')"}
 foreach($P in $R.points){if([int]$P.scenarios-ne 64){throw 'UP89A_SCENARIO_COUNT'};if([int]$P.queries_per_scenario-ne 4){throw 'UP89A_QUERY_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "representation=$($P.representation) length=$($P.length) noise=$($P.noise_amplitude) query=$($P.query_accuracy) exact=$($P.exact_query_set_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP89_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
