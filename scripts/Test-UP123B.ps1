$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up123b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up123b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP123B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP123B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP123B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP123B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up123b-archives-orthogonal -count=1;if($LASTEXITCODE-ne 0){throw 'UP123B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP123B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up123b-archives-orthogonal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP123B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up123b-archives-orthogonal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP123B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP123B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up123b-archives-orthogonal.v1'){throw 'UP123B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP123B_STATE_DIM'}
 if(-not $R.stores_correction_frozen -or -not $R.archives_correction_frozen){throw 'UP123B_CORRECTION_NOT_FROZEN'}
 if($R.adaptive_geometry){throw 'UP123B_ADAPTIVE_GEOMETRY'}
 if([int]$R.points.Count-ne 168){throw "UP123B_POINT_COUNT_$($R.points.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in @('stores_only','stores_plus_archives')){
  $P=@($R.points|Where-Object{$_.arm-ceq $A})
  Write-Host "arm=$A min_archives_original=$(($P|Measure-Object archives_original_accuracy -Minimum).Minimum) min_archives_unseen=$(($P|Measure-Object archives_unseen_accuracy -Minimum).Minimum) min_base_original=$(($P|Measure-Object base_original_accuracy -Minimum).Minimum) min_base_unseen=$(($P|Measure-Object base_unseen_accuracy -Minimum).Minimum) min_acquired=$(($P|Where-Object{[int]$_.stage-gt 0}|Measure-Object acquired_aggregate_accuracy -Minimum).Minimum)"
 }
 Write-Host 'WINGLESS_UP169_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
