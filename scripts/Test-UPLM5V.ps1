$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UPLM5V_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-uplm5v-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5v-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{

 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5V_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5V_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5v-prehazard-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5V_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5V_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5v-prehazard-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5V_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5v-prehazard-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5V_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5v-prehazard-transfer.v1' -or [int]$R.pooled_condition_cells-ne 72 -or [int]$R.matched_conditions_per_policy_cell-ne 384 -or [int]$R.policies.Count-ne 3 -or [int]$R.deadline_profiles.Count-ne 3 -or [int]$R.summaries.Count-ne 216){throw 'UPLM5V_DESIGN'}
 if(-not $R.resource_reductions_used -or -not $R.prepressure_used -or -not $R.prehazard_rule_fixed -or $R.future_schedule_used -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5V_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($Profile in $R.deadline_profiles){foreach($Policy in $R.policies){$Rows=@($R.summaries|Where-Object{$_.deadline_profile-eq $Profile -and $_.policy-eq $Policy});$A=0;$F=0;foreach($S in $Rows){$A+=[int]$S.actions;$F+=[int]$S.failures};Write-Host "profile=$Profile policy=$Policy actions=$A failures=$F"}}
 Write-Host 'WINGLESS_UP365_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
