$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up142b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up142b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP142B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP142B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP142B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP142B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up142b-family-diversity-interaction -count=1;if($LASTEXITCODE-ne 0){throw 'UP142B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP142B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up142b-family-diversity-interaction)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP142B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up142b-family-diversity-interaction)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP142B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP142B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up142b-family-diversity-interaction.v1'){throw 'UP142B_SCHEMA_MISMATCH'}
 if([int]$R.cells.Count-ne 6){throw 'UP142B_CELL_COUNT'}
 if([int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15 -or [int]$R.final_new_updates-ne 4){throw 'UP142B_BUDGET'}
 if($R.extra_updates_used -or $R.adaptive_family_choice -or $R.projector_recomputed){throw 'UP142B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "original_old_delta20v1=$($R.original_old_retention_delta20_vs_1)"
 Write-Host "replication_old_delta20v1=$($R.replication_old_retention_delta20_vs_1)"
 Write-Host "interaction_delta=$($R.old_retention_interaction_delta)"
 Write-Host 'WINGLESS_UP142_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
