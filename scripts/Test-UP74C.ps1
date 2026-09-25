$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up74c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up74c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP74C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP74C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP74C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP74C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up74c-partitioned-carrier -count=1;if($LASTEXITCODE-ne 0){throw 'UP74C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP74C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up74c-partitioned-carrier)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP74C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up74c-partitioned-carrier)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP74C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP74C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up74c-partitioned-carrier.v1'){throw 'UP74C_SCHEMA_MISMATCH'}
 if($R.exact_recall_used){throw 'UP74C_EXACT_RECALL_PRESENT'}
 if([int]$R.recurrent_state_bytes-ne 512){throw 'UP74C_STATE_BYTES'}
 if([int]$R.points.Count-ne 26){throw "UP74C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recurrent_state_bytes-ne 512){throw 'UP74C_POINT_STATE_BYTES'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) family=$($M.family) setting=$($M.setting) primary=$($M.primary_accuracy) exact=$($M.exact_episode_accuracy) cross=$($M.cross_bank_corruption_rate)"}
 Write-Host 'WINGLESS_UP74_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
