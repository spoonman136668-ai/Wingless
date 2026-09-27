$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3e-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3e-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3E_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3E_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3E_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3E_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3e-deployability-equivalence -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3E_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3E_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3e-deployability-equivalence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3E_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3e-deployability-equivalence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3E_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3E_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3e-deployability-equivalence.v1'){throw 'UPLM3E_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 96 -or [int]$R.action_groups.Count-lt 1 -or -not $R.one_action_per_arm_per_round){throw 'UPLM3E_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_resource_choice_used -or $R.future_schedule_oracle_used){throw 'UPLM3E_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($G in $R.action_groups){Write-Host "profile=$($G.profile) actions=$($G.actual_actions) configurations=$($G.configurations) min_failed=$($G.min_failed) max_failed=$($G.max_failed) spread=$($G.failure_spread)"}
 Write-Host 'WINGLESS_UP193_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
