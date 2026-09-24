$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up61b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up61b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP61B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP61B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP61B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP61B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up61b-margin-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP61B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP61B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up61b-margin-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP61B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up61b-margin-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP61B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP61B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up61b-margin-calibration.v1'){throw 'UP61B_SCHEMA_MISMATCH'}
 if($R.oracle_used){throw 'UP61B_ORACLE_PRESENT'}
 if($R.policy_changed){throw 'UP61B_POLICY_CHANGED'}
 if($R.training_changed){throw 'UP61B_TRAINING_CHANGED'}
 if($R.threshold_selected){throw 'UP61B_THRESHOLD_SELECTION_PRESENT'}
 if([int]$R.metrics.Count-ne 12){throw 'UP61B_METRIC_COUNT'}
 foreach($M in $R.metrics){
   if([int]$M.buckets.Count-ne 4){throw 'UP61B_BUCKET_COUNT'}
   $Events=0;foreach($B in $M.buckets){$Events += [int]$B.events}
   if($Events-ne 3072){throw "UP61B_EVENT_COUNT schedule=$($M.schedule_base) noise=$($M.memory_noise) events=$Events"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){
   foreach($B in $M.buckets){
     Write-Host "schedule=$($M.schedule_base) noise=$($M.memory_noise) lower=$($B.lower_bound) upper=$($B.upper_bound) events=$($B.events) correct=$($B.correct) accuracy=$($B.accuracy)"
   }
 }
 Write-Host 'WINGLESS_UP61_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
