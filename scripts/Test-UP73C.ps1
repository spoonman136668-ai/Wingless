$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up73c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up73c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP73C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP73C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP73C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP73C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up73c-irregular-tag-family-replication -count=1;if($LASTEXITCODE-ne 0){throw 'UP73C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP73C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up73c-irregular-tag-family-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP73C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up73c-irregular-tag-family-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP73C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP73C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up73c-irregular-tag-family-replication.v1'){throw 'UP73C_SCHEMA_MISMATCH'}
 if($R.selection_performed){throw 'UP73C_SELECTION_PRESENT'}
 if([int]$R.points.Count-ne 12){throw "UP73C_POINT_COUNT $($R.points.Count)"}
 $Families=@($R.points|ForEach-Object{$_.tag_family}|Sort-Object -Unique)
 if(($Families -join ',') -cne 'golden_rotation,sqrt2_rotation,sqrt3_rotation'){throw "UP73C_FAMILIES $($Families -join ',')"}
 foreach($P in $R.points){if([int]$P.banks-ne 10){throw 'UP73C_BANK_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "family=$($P.tag_family) schedule=$($P.schedule_base) noise=$($P.memory_noise) gate=$($P.gate) value=$($P.value_accuracy) exact=$($P.exact_scenario_accuracy) error=$($P.evaluation_error)"}
 Write-Host 'WINGLESS_UP73_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
