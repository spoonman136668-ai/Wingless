$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up150b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up150b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP150B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP150B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP150B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP150B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up150b-terminal-tail-geometry -count=1;if($LASTEXITCODE-ne 0){throw 'UP150B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP150B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up150b-terminal-tail-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP150B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up150b-terminal-tail-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP150B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP150B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up150b-terminal-tail-geometry.v1'){throw 'UP150B_SCHEMA_MISMATCH'}
 if([int]$R.tails.Count-ne 4 -or [int]$R.pairwise.Count-ne 6 -or [int]$R.tail_examples_per_shift-ne 4){throw 'UP150B_COUNT'}
 if([int]$R.common_prefix_epochs-ne 15){throw 'UP150B_PREFIX'}
 if($R.diagnostic_updates_retained -or $R.adaptive_tail_selection -or $R.extra_training_used){throw 'UP150B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($T in $R.tails){Write-Host "shift=$($T.shift) prior=$($T.prior_class) norm=$($T.tail_update_norm) old_cos=$($T.cosine_tail_vs_old_anchor) p=$($T.mean_correct_class_probability)"}
 foreach($P in $R.pairwise){Write-Host "pair=$($P.shift_a)-$($P.shift_b) cos=$($P.cosine)"}
 Write-Host 'WINGLESS_UP150_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
