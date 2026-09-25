$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up76c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up76c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP76C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP76C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP76C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP76C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up76c-partition-granularity -count=1;if($LASTEXITCODE-ne 0){throw 'UP76C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP76C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up76c-partition-granularity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP76C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up76c-partition-granularity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP76C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP76C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up76c-partition-granularity.v1'){throw 'UP76C_SCHEMA_MISMATCH'}
 if($R.exact_recall_used){throw 'UP76C_EXACT_RECALL_PRESENT'}
 if($R.dynamic_routing_used){throw 'UP76C_DYNAMIC_ROUTING_PRESENT'}
 if([int]$R.recurrent_state_bytes-ne 512){throw 'UP76C_STATE_BYTES'}
 if([int]$R.points.Count-ne 27){throw "UP76C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recurrent_state_bytes-ne 512){throw 'UP76C_POINT_STATE_BYTES'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) family=$($M.family) setting=$($M.setting) primary=$($M.primary_accuracy) exact=$($M.exact_episode_accuracy) cross=$($M.cross_bank_corruption_rate)"}
 Write-Host 'WINGLESS_UP98_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
