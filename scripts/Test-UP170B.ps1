$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up170b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up170b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP170B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP170B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP170B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP170B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up170b-margin-rank-monotonic-risk -count=1;if($LASTEXITCODE-ne 0){throw 'UP170B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP170B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up170b-margin-rank-monotonic-risk)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP170B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up170b-margin-rank-monotonic-risk)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP170B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP170B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up170b-margin-rank-monotonic-risk.v1'){throw 'UP170B_SCHEMA_MISMATCH'}
 if([int]$R.bin_metrics.Count-ne 15 -or [int]$R.class_trends.Count-ne 3 -or [int]$R.context_phase-ne 17 -or [int]$R.new_prefix_rotation-ne 7){throw 'UP170B_DESIGN'}
 if($R.adaptive_rotation_used -or $R.threshold_fitted -or $R.maintenance_triggered){throw 'UP170B_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.bin_metrics){Write-Host "class=$($M.class) ranks=$($M.rank_start)-$($M.rank_end) crossings=$($M.crossings) density=$($M.crossing_density)"}
 foreach($T in $R.class_trends){Write-Host "trend class=$($T.class) nonincreasing=$($T.non_increasing_across_bins)"}
 Write-Host 'WINGLESS_UP170_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
