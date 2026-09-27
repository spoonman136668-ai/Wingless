$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3d-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3d-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3D_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3D_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3D_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3D_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3d-compressed-window -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3D_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3D_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3d-compressed-window)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3D_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3d-compressed-window)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3D_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3D_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3d-compressed-window.v1'){throw 'UPLM3D_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 144 -or [int]$R.metrics.Count-ne 12 -or [int]$R.total_budget-ne 6 -or -not $R.distinct_arm_per_round){throw 'UPLM3D_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_onset_used -or $R.adaptive_budget_used -or $R.adaptive_throughput_used -or $R.future_schedule_used){throw 'UPLM3D_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "profile=$($M.profile) start=$($M.action_start_round) throughput=$($M.throughput) earliest_failed=$($M.earliest_failed) fixed_failed=$($M.fixed_failed) earliest_actions=$($M.earliest_actions) advantage=$($M.earliest_advantage)"}
 Write-Host 'WINGLESS_UP192_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
