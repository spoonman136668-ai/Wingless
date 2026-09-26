$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up148b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up148b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP148B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP148B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP148B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP148B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up148b-terminal-tail-crossover -count=1;if($LASTEXITCODE-ne 0){throw 'UP148B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP148B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up148b-terminal-tail-crossover)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP148B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up148b-terminal-tail-crossover)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP148B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP148B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up148b-terminal-tail-crossover.v1'){throw 'UP148B_SCHEMA_MISMATCH'}
 if([int]$R.arm_summaries.Count-ne 4 -or [int]$R.block_metrics.Count-ne 16 -or [int]$R.pair_deltas.Count-ne 2){throw 'UP148B_METRIC_COUNT'}
 if([int]$R.epochs-ne 20 -or [int]$R.distinct_tail_sets-ne 4 -or [int]$R.uses_per_tail_set-ne 5 -or [int]$R.switches-ne 3){throw 'UP148B_DESIGN'}
 if($R.adaptive_ordering_used -or $R.extra_updates_used -or $R.projector_recomputed){throw 'UP148B_BOUNDARY_LEAK'}
 foreach($A in $R.arm_summaries){$B=@($R.block_metrics|Where-Object{$_.arm-ceq $A.arm});$S=@($B|ForEach-Object{[int]$_.shift}|Sort-Object);if(($S -join ',')-cne '0,6,12,18'){throw 'UP148B_COVERAGE'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($D in $R.pair_deltas){Write-Host "pair=$($D.pair) old_delta=$($D.old_retention_delta) new_delta=$($D.new_accuracy_delta)"}
 Write-Host 'WINGLESS_UP148_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
