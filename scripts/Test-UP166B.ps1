$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up166b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up166b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP166B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP166B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP166B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP166B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up166b-margin-vulnerability-rank -count=1;if($LASTEXITCODE-ne 0){throw 'UP166B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP166B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up166b-margin-vulnerability-rank)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP166B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up166b-margin-vulnerability-rank)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP166B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP166B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up166b-margin-vulnerability-rank.v1'){throw 'UP166B_SCHEMA_MISMATCH'}
 if([int]$R.steps.Count-ne 24 -or [int]$R.subjects-ne 6 -or [int]$R.paths-ne 2 -or [int]$R.report_steps-ne 2){throw 'UP166B_DESIGN'}
 if($R.ranking_signal-cne 'absolute_pre_step_probability_margin' -or $R.post_result_threshold_used -or $R.diagnostic_updates_retained){throw 'UP166B_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.steps){Write-Host "subject=$($S.subject) path=$($S.path) report=$($S.report_index) crossings=$($S.crossing_count) mean_rank=$($S.mean_crossing_rank) max_rank=$($S.max_crossing_rank) frac=$($S.mean_crossing_rank_fraction)"}
 Write-Host 'WINGLESS_UP166_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
