$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm5q-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5q-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5Q_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5Q_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5q-no-prepressure-policy-control -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5Q_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5Q_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5q-no-prepressure-policy-control)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5Q_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5q-no-prepressure-policy-control)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5Q_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5q-no-prepressure-policy-control.v1' -or [int]$R.condition_cells-ne 6 -or [int]$R.matched_conditions_per_cell-ne 64 -or [int]$R.policies.Count-ne 3 -or [int]$R.summaries.Count-ne 6){throw 'UPLM5Q_DESIGN'}
 if($R.resource_reductions_used -or $R.prepressure_used -or -not $R.only_policy_changed -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5Q_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 $EF=0;$EL=0;$ELF=0;$FL=0;$ELL=0;$LL=0;$EFail=0;$FFail=0;$LFail=0
 foreach($S in $R.summaries){$EF+=[int]$S.earliest_fixed_differing;$EL+=[int]$S.earliest_latest_differing;$ELF+=[int]$S.earliest_lower_than_fixed;$FL+=[int]$S.fixed_lower_than_earliest;$ELL+=[int]$S.earliest_lower_than_latest;$LL+=[int]$S.latest_lower_than_earliest;$EFail+=[int]$S.earliest_failures;$FFail+=[int]$S.fixed_failures;$LFail+=[int]$S.latest_failures}
 Write-Host "ef_differ=$EF earliest_lt_fixed=$ELF fixed_lt_earliest=$FL el_differ=$EL earliest_lt_latest=$ELL latest_lt_earliest=$LL earliest_failures=$EFail fixed_failures=$FFail latest_failures=$LFail"
 Write-Host 'WINGLESS_UP351_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
