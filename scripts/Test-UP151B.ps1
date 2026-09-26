$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up151b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up151b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP151B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP151B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP151B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP151B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up151b-tail-class-composition -count=1;if($LASTEXITCODE-ne 0){throw 'UP151B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP151B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up151b-tail-class-composition)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP151B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up151b-tail-class-composition)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP151B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP151B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up151b-tail-class-composition.v1'){throw 'UP151B_SCHEMA_MISMATCH'}
 if([int]$R.arms.Count-ne 10 -or [int]$R.subjects-ne 2 -or [int]$R.compositions-ne 5){throw 'UP151B_DESIGN'}
 if([int]$R.common_prefix_epochs-ne 15 -or [int]$R.terminal_epochs-ne 5 -or [int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15 -or [int]$R.final_new_updates-ne 4){throw 'UP151B_BUDGET'}
 foreach($A in $R.arms){if([int]$A.tail_verbs.Count-ne 4){throw 'UP151B_TAIL_COUNT'}}
 if($R.adaptive_tail_selection -or $R.extra_updates_used){throw 'UP151B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in $R.arms){Write-Host "subject=$($A.subject) composition=$($A.composition) tail=$($A.mean_terminal_tail_damage) old=$($A.final_mean_old_retention) new=$($A.final_mean_new_accuracy)"}
 Write-Host 'WINGLESS_UP151_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
