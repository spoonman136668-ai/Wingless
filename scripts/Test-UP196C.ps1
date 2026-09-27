$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up196c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up196c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP196C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP196C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up196c-hostile-heldout-cohort-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP196C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP196C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up196c-hostile-heldout-cohort-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP196C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up196c-hostile-heldout-cohort-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP196C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up196c-hostile-heldout-cohort-transfer.v1' -or $R.cohort_topology-cne 'contiguous_quartets' -or [int]$R.primary_spacing-ne 4 -or [int]$R.primary_cap-ne 8){throw 'UP196C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.cohort_specific_tuning_used -or $R.warning_threshold_changed){throw 'UP196C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "arms=$($R.arms) baseline=$($R.baseline_losses) primary_prevented=$($R.primary_prevented) primary_actions=$($R.primary_actions) primary_accelerated=$($R.primary_accelerated) primary_mean_extension=$($R.primary_mean_loss_extension_among_failures) primary_max_extension=$($R.primary_max_loss_extension_among_failures) dense_prevented=$($R.dense_prevented) dense_actions=$($R.dense_actions) dense_accelerated=$($R.dense_accelerated)"
 Write-Host 'WINGLESS_UP255_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
