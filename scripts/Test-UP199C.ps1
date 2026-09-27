$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up199c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up199c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP199C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP199C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up199c-hostile-commitment-depth-gating -count=1;if($LASTEXITCODE-ne 0){throw 'UP199C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP199C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up199c-hostile-commitment-depth-gating)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP199C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up199c-hostile-commitment-depth-gating)|Out-String).Trim();if($P1-cne $P2){throw 'UP199C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up199c-hostile-commitment-depth-gating.v1' -or [int]$R.metrics.Count-ne 24 -or [int]$R.commitment_depths.Count-ne 4 -or [int]$R.action_cap-ne 8){throw 'UP199C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_selection_used){throw 'UP199C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "depth=$($M.commitment_depth) horizon=$($M.horizon) baseline=$($M.baseline_failures) treated=$($M.treated_failures) prevented=$($M.prevented_failures) actions=$($M.actions_taken) skipped=$($M.skipped_opportunities)"}
 Write-Host 'WINGLESS_UP263_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
