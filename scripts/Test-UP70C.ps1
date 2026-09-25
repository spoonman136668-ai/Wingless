$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up70c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up70c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP70C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP70C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP70C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP70C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up70c-golden-bank-capacity -count=1;if($LASTEXITCODE-ne 0){throw 'UP70C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP70C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up70c-golden-bank-capacity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP70C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up70c-golden-bank-capacity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP70C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP70C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up70c-golden-bank-capacity.v1'){throw 'UP70C_SCHEMA_MISMATCH'}
 if($R.selection_performed){throw 'UP70C_SELECTION_PRESENT'}
 if([int]$R.points.Count-ne 24){throw "UP70C_POINT_COUNT $($R.points.Count)"}
 $Banks=@($R.points|ForEach-Object{[int]$_.banks}|Sort-Object -Unique)
 if(($Banks -join ',') -cne '6,7,8,10'){throw "UP70C_BANK_LEVELS $($Banks -join ',')"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "schedule=$($P.schedule_base) banks=$($P.banks) noise=$($P.memory_noise) gate=$($P.gate) value=$($P.value_accuracy) exact=$($P.exact_scenario_accuracy) magnitude=$($P.magnitude_value_accuracy)"}
 Write-Host 'WINGLESS_UP70_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
