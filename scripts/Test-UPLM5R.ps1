$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false
 go clean -cache -testcache
 go test ./unitary ./cmd/unitary-up-lm5r-action-efficiency-mechanism -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5R_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5R_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5r-action-efficiency-mechanism)|Out-String).Trim()
 $P2=((go run ./cmd/unitary-up-lm5r-action-efficiency-mechanism)|Out-String).Trim()
 if($P1-cne $P2){throw 'UPLM5R_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5r-action-efficiency-mechanism.v1' -or [int]$R.condition_cells-ne 6 -or [int]$R.matched_conditions_per_cell-ne 64 -or [int]$R.summaries.Count-ne 6){throw 'UPLM5R_DESIGN'}
 if($R.resource_reductions_used -or $R.prepressure_used -or -not $R.action_count_measured -or -not $R.only_policy_changed -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5R_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 $EFSame=0;$ELSame=0;$ELF=0;$ELL=0
 foreach($S in $R.summaries){$EFSame+=[int]$S.earliest_fixed_same_actions_differing;$ELSame+=[int]$S.earliest_latest_same_actions_differing;$ELF+=[int]$S.earliest_lower_fixed_same_actions;$ELL+=[int]$S.earliest_lower_latest_same_actions}
 Write-Host "ef_same_actions_differ=$EFSame earliest_lower_fixed_same_actions=$ELF el_same_actions_differ=$ELSame earliest_lower_latest_same_actions=$ELL"
 Write-Host 'WINGLESS_UP354_HARNESS_PASS'
}finally{Pop-Location}
