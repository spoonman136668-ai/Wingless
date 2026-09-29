$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR;$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-uplm6a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm6a-gotmp-"+$PID);New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null;$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM6A_GO_ENV_FAILED'};go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM6A_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm6a-prepressure-length-support -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM6A_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM6A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm6a-prepressure-length-support)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM6A_PROBE1_FAILED'};$P2=((go run ./cmd/unitary-up-lm6a-prepressure-length-support)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM6A_NONDETERMINISTIC_OUTPUT'};$R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm6a-prepressure-length-support.v1' -or [int]$R.prepressure_on.conditions-ne 27648 -or [int]$R.prepressure_off.conditions-ne 27648){throw 'UPLM6A_DESIGN'}
 if([int]$R.parent_raw_sub16-ne 0 -or [int]$R.parent_eligible_sub16-ne 0 -or [int]$R.prepressure_on.raw_sub16-ne 0 -or [int]$R.prepressure_on.decision_points-ne [int]$R.parent_decision_points){throw 'UPLM6A_PARENT_ANCHOR_DRIFT'}
 if(-not $R.canonical_hazard_selector_used -or $R.policy_changed -or $R.future_schedule_used -or $R.adaptive_policy_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM6A_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';Write-Host "on_raw_sub16=$($R.prepressure_on.raw_sub16) off_raw_sub16=$($R.prepressure_off.raw_sub16) classification=$($R.classification)";Write-Host 'WINGLESS_UP380_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$PriorGOROOT;$env:GOCACHE=$PriorGoCache;$env:GOTMPDIR=$PriorGoTmp;Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
