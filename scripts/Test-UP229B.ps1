$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
if(-not (Test-Path (Join-Path $GoRoot 'src\runtime'))){throw "UP229B_GOROOT_LAYOUT_INVALID:$GoRoot"}
$GoCache=Join-Path $env:TEMP ("wingless-up229b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up229b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP229B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP229B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up229b-parity-holdout-generalization -count=1;if($LASTEXITCODE-ne 0){throw 'UP229B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP229B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up229b-parity-holdout-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP229B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up229b-parity-holdout-generalization)|Out-String).Trim();if($P1-cne $P2){throw 'UP229B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up229b-parity-holdout-generalization.v1' -or [int]$R.training_phases.Count-ne 15 -or [int]$R.holdout_phases.Count-ne 16 -or [int]$R.families.Count-ne 4 -or [int]$R.summaries.Count-ne 4){throw 'UP229B_DESIGN'}
 if($R.heldout_fitting_used -or $R.parity_used_for_modeling -or $R.gate_changed -or $R.adaptive_feature_selection_used -or $R.new_native_feature_used -or $R.live_activation){throw 'UP229B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "family=$($S.family) holdout=$($S.holdout_states) same_parity=$($S.same_parity_exact_matches) opposite_parity=$($S.opposite_parity_exact_matches) unmatched=$($S.unmatched_holdout_states)"}
 Write-Host 'WINGLESS_UP357_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
