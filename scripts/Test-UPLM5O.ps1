$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm5o-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5o-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5O_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5O_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5o-deadline-order-ablation -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5O_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5O_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5o-deadline-order-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5O_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5o-deadline-order-ablation)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5O_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5o-deadline-order-ablation.v1' -or [int]$R.pooled_condition_cells-ne 72 -or [int]$R.matched_conditions_per_cell-ne 384 -or [int]$R.policies.Count-ne 3 -or [int]$R.summaries.Count-ne 72){throw 'UPLM5O_DESIGN'}
 if(-not $R.only_policy_changed -or $R.adaptive_policy_selection_used -or $R.adaptive_coordinate_search_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5O_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($Profile in @('deferred_only','layout_only','hybrid_min')){$Rows=@($R.summaries|Where-Object{$_.deadline_profile-eq $Profile});$EF=0;$EL=0;$EarliestFixed=0;$Fixed=0;$EarliestLatest=0;$Latest=0;foreach($S in $Rows){$EF+=[int]$S.earliest_fixed_differing;$EL+=[int]$S.earliest_latest_differing;$EarliestFixed+=[int]$S.earliest_lower_than_fixed;$Fixed+=[int]$S.fixed_lower_than_earliest;$EarliestLatest+=[int]$S.earliest_lower_than_latest;$Latest+=[int]$S.latest_lower_than_earliest};Write-Host "profile=$Profile ef_differ=$EF earliest_lt_fixed=$EarliestFixed fixed_lt_earliest=$Fixed el_differ=$EL earliest_lt_latest=$EarliestLatest latest_lt_earliest=$Latest"}
 Write-Host 'WINGLESS_UP345_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
