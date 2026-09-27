$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm5n-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5n-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5N_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5N_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5n-policy-causal-swap -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5N_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5N_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5n-policy-causal-swap)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5N_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5n-policy-causal-swap)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5N_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5n-policy-causal-swap.v1' -or [int]$R.pooled_condition_cells-ne 72 -or [int]$R.matched_conditions_per_cell-ne 384 -or [int]$R.summaries.Count-ne 72 -or [int]$R.policies.Count-ne 2){throw 'UPLM5N_DESIGN'}
 if(-not $R.only_policy_changed -or $R.adaptive_policy_selection_used -or $R.adaptive_coordinate_search_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5N_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($Profile in @('deferred_only','layout_only','hybrid_min')){$Rows=@($R.summaries|Where-Object{$_.deadline_profile-eq $Profile});$Diff=0;$Earliest=0;$Fixed=0;$Equal=0;$Max=0;foreach($S in $Rows){$Diff+=[int]$S.differing_pairs;$Earliest+=[int]$S.earliest_lower_failures;$Fixed+=[int]$S.fixed_lower_failures;$Equal+=[int]$S.equal_failures;if([int]$S.max_abs_failure_delta-gt $Max){$Max=[int]$S.max_abs_failure_delta}};Write-Host "profile=$Profile cells=$($Rows.Count) differing_pairs=$Diff earliest_lower=$Earliest fixed_lower=$Fixed equal=$Equal max_abs_delta=$Max"}
 Write-Host 'WINGLESS_UP342_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
