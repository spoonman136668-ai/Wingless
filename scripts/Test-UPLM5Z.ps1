$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-uplm5z-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5z-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5Z_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5Z_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5z-order-length-origin -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5Z_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5Z_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5z-order-length-origin)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5Z_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5z-order-length-origin)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5Z_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5z-order-length-origin.v1' -or [int]$R.total_conditions-ne 27648 -or [int]$R.total_decision_points-le 0 -or [int]$R.summaries.Count-ne 72){throw 'UPLM5Z_DESIGN'}
 if([int]$R.parent_subhazard_candidates-ne 0 -or [int]$R.parent_observed_lengths.Count-ne 1 -or [int]$R.parent_observed_lengths[0]-ne 16){throw 'UPLM5Z_PARENT_ANCHOR_DRIFT'}
 if([int]$R.parent_decision_points-ne [int]$R.total_decision_points -or [int]$R.parent_eligible_candidates-ne [int]$R.counts.eligible){throw 'UPLM5Z_TRAJECTORY_DRIFT'}
 if([int]$R.eligible_sub16_candidates-ne 0){throw 'UPLM5Z_ELIGIBLE_SUB16_PARENT_DRIFT'}
 if(-not $R.canonical_hazard_selector_used -or $R.policy_changed -or $R.future_schedule_used -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5Z_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "raw_sub16=$($R.raw_sub16_candidates) eligible_sub16=$($R.eligible_sub16_candidates) used=$($R.counts.used_this_decision) empty=$($R.counts.empty) nonoriginal=$($R.counts.non_original_victim) reported=$($R.counts.already_reported) eligible=$($R.counts.eligible)"
 Write-Host 'WINGLESS_UP376_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$PriorGOROOT;$env:GOCACHE=$PriorGoCache;$env:GOTMPDIR=$PriorGoTmp;Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
