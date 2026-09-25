$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up69c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up69c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP69C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP69C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP69C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP69C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up69c-boundary-depth-scale -count=1;if($LASTEXITCODE-ne 0){throw 'UP69C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP69C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up69c-boundary-depth-scale)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP69C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up69c-boundary-depth-scale)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP69C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP69C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up69c-boundary-depth-scale.v1'){throw 'UP69C_SCHEMA_MISMATCH'}
 if($R.selection_performed){throw 'UP69C_SELECTION_PRESENT'}
 if([int]$R.metrics.Count-ne 18){throw "UP69C_METRIC_COUNT $($R.metrics.Count)"}
 $Depths=@($R.metrics|ForEach-Object{[int]$_.depth}|Sort-Object -Unique)
 if(($Depths -join ',') -cne '1024,2048,4096'){throw "UP69C_DEPTH_LEVELS $($Depths -join ',')"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "family=$($M.family) schedule=$($M.schedule_base) noise=$($M.memory_noise) depth=$($M.depth) gate=$($M.gate) value=$($M.value_accuracy) exact=$($M.exact_scenario_accuracy)"}
 Write-Host 'WINGLESS_UP69_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
