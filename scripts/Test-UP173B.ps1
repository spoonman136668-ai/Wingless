$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up173b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up173b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP173B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP173B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP173B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP173B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up173b-risk-transition-hazard -count=1;if($LASTEXITCODE-ne 0){throw 'UP173B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP173B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up173b-risk-transition-hazard)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP173B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up173b-risk-transition-hazard)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP173B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP173B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up173b-risk-transition-hazard.v1'){throw 'UP173B_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 16 -or [int]$R.summaries.Count-ne 4 -or [int]$R.context_phase-ne 18){throw 'UP173B_DESIGN'}
 if(-not $R.bands_frozen_before_run -or $R.maintenance_triggered -or $R.updates_suppressed){throw 'UP173B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "scope=$($S.scope) transition=$($S.transition_density) established=$($S.established_density) ratio=$($S.transition_vs_established_ratio)"}
 Write-Host 'WINGLESS_UP173_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
