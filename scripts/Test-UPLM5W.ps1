$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UPLM5W_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-uplm5w-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5w-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5W_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5W_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5w-selector-equivalence-trace -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5W_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5W_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5w-selector-equivalence-trace)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5W_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5w-selector-equivalence-trace)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5W_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5w-selector-equivalence-trace.v1' -or [int]$R.pooled_condition_cells-ne 72 -or [int]$R.matched_conditions_per_cell-ne 384 -or [int]$R.total_compared_conditions-ne 27648 -or [int]$R.summaries.Count-ne 72){throw 'UPLM5W_DESIGN'}
 if(-not $R.resource_reductions_used -or -not $R.prepressure_used -or -not $R.prehazard_rule_fixed -or $R.future_schedule_used -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5W_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "selector_differences=$($R.total_selector_differences) action_differences=$($R.total_action_differences) failure_differences=$($R.total_failure_differences)"
 foreach($S in $R.summaries){if([int]$S.selector_differences-ne 0 -or [int]$S.action_differences-ne 0 -or [int]$S.failure_differences-ne 0){Write-Host "cell=$($S.deadline_profile)|br$($S.budget_reduction)|tr$($S.throughput_reduction)|rot$($S.rotation)|$($S.permutation) selector_diff=$($S.selector_differences) action_diff=$($S.action_differences) failure_diff=$($S.failure_differences)"}}
 Write-Host 'WINGLESS_UP367_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
