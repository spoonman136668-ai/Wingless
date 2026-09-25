$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up63b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up63b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP63B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP63B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP63B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP63B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up63b-margin-depth-stratification -count=1;if($LASTEXITCODE-ne 0){throw 'UP63B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP63B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up63b-margin-depth-stratification)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP63B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up63b-margin-depth-stratification)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP63B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP63B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up63b-margin-depth-stratification.v1'){throw 'UP63B_SCHEMA_MISMATCH'}
 if($R.oracle_used){throw 'UP63B_ORACLE_PRESENT'}
 if($R.policy_changed){throw 'UP63B_POLICY_CHANGED'}
 if($R.training_changed){throw 'UP63B_TRAINING_CHANGED'}
 if($R.threshold_selected){throw 'UP63B_THRESHOLD_SELECTION_PRESENT'}
 if([int]$R.metrics.Count-ne 16){throw "UP63B_METRIC_COUNT $($R.metrics.Count)"}
 $Depths=@($R.metrics|ForEach-Object{[int]$_.depth}|Sort-Object -Unique)
 if(($Depths -join ',') -cne '32,128,512,1024'){throw "UP63B_DEPTH_LEVELS $($Depths -join ',')"}
 foreach($M in $R.metrics){
   if([int]$M.buckets.Count-ne 4){throw 'UP63B_BUCKET_COUNT'}
   $Events=0;foreach($B in $M.buckets){$Events += [int]$B.events}
   $Expected=48*[int]$M.writes_per_scenario
   if($Events-ne $Expected){throw "UP63B_EVENT_COUNT schedule=$($M.schedule_base) noise=$($M.memory_noise) depth=$($M.depth) events=$Events expected=$Expected"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){foreach($B in $M.buckets){Write-Host "schedule=$($M.schedule_base) noise=$($M.memory_noise) depth=$($M.depth) writes=$($M.writes_per_scenario) lower=$($B.lower_bound) upper=$($B.upper_bound) events=$($B.events) correct=$($B.correct) accuracy=$($B.accuracy)"}}
 Write-Host 'WINGLESS_UP63_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
