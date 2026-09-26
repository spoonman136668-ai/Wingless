$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up171b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up171b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP171B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP171B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP171B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP171B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up171b-shadow-risk-trajectory -count=1;if($LASTEXITCODE-ne 0){throw 'UP171B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP171B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up171b-shadow-risk-trajectory)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP171B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up171b-shadow-risk-trajectory)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP171B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP171B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up171b-shadow-risk-trajectory.v1'){throw 'UP171B_SCHEMA'}
 if([int]$R.steps.Count-ne 144 -or [int]$R.subjects-ne 6 -or [int]$R.paths-ne 2 -or [int]$R.steps_per_path-ne 12){throw 'UP171B_DESIGN'}
 if([int]$R.examples_per_state-ne 120 -or -not $R.bands_frozen_before_run){throw 'UP171B_BANDS'}
 if($R.adaptive_band_used -or $R.maintenance_triggered -or $R.updates_suppressed){throw 'UP171B_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.steps){Write-Host "subject=$($S.subject) path=$($S.path) step=$($S.step) class=$($S.class) crossings=$($S.crossings) covered=$($S.covered_crossings) coverage=$($S.coverage) density=$($S.crossing_density)"}
 Write-Host 'WINGLESS_UP171_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
