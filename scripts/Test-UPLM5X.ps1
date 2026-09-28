$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UPLM5X_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-uplm5x-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5x-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5X_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5X_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5x-threat-predicate-equivalence -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5X_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5X_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5x-threat-predicate-equivalence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5X_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5x-threat-predicate-equivalence)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5X_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5x-threat-predicate-equivalence.v1' -or [int]$R.pooled_condition_cells-ne 72 -or [int]$R.matched_conditions_per_cell-ne 384 -or [int]$R.total_conditions-ne 27648 -or [int]$R.summaries.Count-ne 72 -or [int]$R.total_decision_points-le 0){throw 'UPLM5X_DESIGN'}
 if(-not $R.resource_reductions_used -or -not $R.prepressure_used -or $R.future_schedule_used -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5X_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "decision_points=$($R.total_decision_points) predicate_difference_points=$($R.total_predicate_difference_points) exposed_prehazard_only_points=$($R.total_exposed_prehazard_only_points) hazard_candidates=$($R.total_hazard_candidates) prehazard_only_candidates=$($R.total_prehazard_only_candidates)"
 foreach($S in $R.summaries){if([int]$S.predicate_difference_points-ne 0 -or [int]$S.exposed_prehazard_only_points-ne 0){Write-Host "cell=$($S.deadline_profile)|br$($S.budget_reduction)|tr$($S.throughput_reduction)|rot$($S.rotation)|$($S.permutation) decision=$($S.decision_points) predicate_diff=$($S.predicate_difference_points) exposed=$($S.exposed_prehazard_only_points)"}}
 Write-Host 'WINGLESS_UP370_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
