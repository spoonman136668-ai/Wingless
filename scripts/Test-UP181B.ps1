$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up181b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up181b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP181B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP181B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP181B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP181B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up181b-margin-conditioned-temporal -count=1;if($LASTEXITCODE-ne 0){throw 'UP181B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP181B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up181b-margin-conditioned-temporal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP181B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up181b-margin-conditioned-temporal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP181B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP181B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up181b-margin-conditioned-temporal.v1' -or [int]$R.context_phase-ne 26 -or [int]$R.metrics.Count-ne 60){throw 'UP181B_DESIGN'}
 if(-not $R.bands_frozen_before_run -or $R.adaptive_bands_used -or $R.maintenance_triggered){throw 'UP181B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){if($M.scope-eq 'ALL'){Write-Host "band=$($M.margin_band) state=$($M.temporal_state) slots=$($M.slots) density=$($M.crossing_density)"}}
 Write-Host 'WINGLESS_UP181_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
