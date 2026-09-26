$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up172b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up172b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP172B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP172B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP172B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP172B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up172b-persistent-risk-shadow -count=1;if($LASTEXITCODE-ne 0){throw 'UP172B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP172B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up172b-persistent-risk-shadow)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP172B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up172b-persistent-risk-shadow)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP172B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP172B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up172b-persistent-risk-shadow.v1'){throw 'UP172B_SCHEMA_MISMATCH'}
 if([int]$R.bins.Count-ne 16 -or [int]$R.warnings.Count-ne 8 -or [int]$R.subjects-ne 6 -or [int]$R.paths-ne 2){throw 'UP172B_DESIGN'}
 if([int]$R.persistence_steps-ne 2 -or [int]$R.steps_per_path-ne 12 -or [int]$R.examples_per_state-ne 120){throw 'UP172B_PERSISTENCE'}
 if(-not $R.bands_frozen_before_run -or $R.adaptive_band_used -or $R.maintenance_triggered -or $R.updates_suppressed){throw 'UP172B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($W in $R.warnings){Write-Host "scope=$($W.scope) warning=$($W.warning) flagged=$($W.flagged_slots) covered=$($W.covered_crossings) density=$($W.crossing_density) coverage=$($W.crossing_coverage)"}
 Write-Host 'WINGLESS_UP172_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
