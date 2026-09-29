$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UPLM5Y_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-uplm5y-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5y-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5Y_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5Y_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5y-eligible-order-length-spectrum -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5Y_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5Y_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5y-eligible-order-length-spectrum)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5Y_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5y-eligible-order-length-spectrum)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5Y_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5y-eligible-order-length-spectrum.v1' -or [int]$R.pooled_condition_cells-ne 72 -or [int]$R.matched_conditions_per_cell-ne 384 -or [int]$R.total_conditions-ne 27648 -or [int]$R.summaries.Count-ne 72 -or [int]$R.total_decision_points-le 0 -or [int]$R.observed_lengths.Count-le 0){throw 'UPLM5Y_DESIGN'}
 if([int]$R.parent_prehazard_only_candidates-ne 0 -or [int]$R.parent_predicate_difference_points-ne 0){throw 'UPLM5Y_PARENT_ANCHOR_DRIFT'}
 if(-not $R.canonical_hazard_selector_used -or $R.policy_changed -or $R.future_schedule_used -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5Y_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "decision_points=$($R.total_decision_points) eligible=$($R.total_eligible_candidates) hazard=$($R.total_hazard_candidates) subhazard=$($R.total_subhazard_candidates) exposed_subhazard_points=$($R.total_exposed_subhazard_decision_points) observed_lengths=$($R.observed_lengths -join ',')"
 Write-Host 'WINGLESS_UP373_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
