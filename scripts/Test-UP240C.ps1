$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up240c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up240c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP240C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP240C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up240c-signature-component-overlap -count=1;if($LASTEXITCODE-ne 0){throw 'UP240C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP240C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up240c-signature-component-overlap)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP240C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up240c-signature-component-overlap)|Out-String).Trim();if($P1-cne $P2){throw 'UP240C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up240c-signature-component-overlap.v1' -or [int]$R.arms_total-ne 512 -or [int]$R.schedule_pairs-ne 6 -or [int]$R.current_components.Count-ne 9 -or [int]$R.trajectory_components.Count-ne 18 -or [int]$R.summaries.Count-ne 162){throw 'UP240C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP240C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($Pair in @('0-1','0-2','0-3','1-2','1-3','2-3')){$A=[int]($Pair.Split('-')[0]);$B=[int]($Pair.Split('-')[1]);foreach($Mode in @('current','trajectory')){$Rows=@($R.summaries|Where-Object{[int]$_.schedule_a-eq $A -and [int]$_.schedule_b-eq $B -and $_.signature_mode-eq $Mode});$Fail=0;$Surv=0;$Conflict=0;foreach($S in $Rows){if([int]$S.shared_failure_values-gt 0){$Fail++};if([int]$S.shared_survivor_values-gt 0){$Surv++};if([int]$S.a_failure_b_survivor_conflicts-gt 0 -or [int]$S.a_survivor_b_failure_conflicts-gt 0){$Conflict++}};Write-Host "pair=$Pair mode=$Mode components=$($Rows.Count) fail_overlap_components=$Fail survive_overlap_components=$Surv conflict_components=$Conflict"}}
 Write-Host 'WINGLESS_UP347_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
