$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm4m-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm4m-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM4M_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM4M_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm4m-budget-revocation-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM4M_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM4M_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm4m-budget-revocation-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM4M_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm4m-budget-revocation-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM4M_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm4m-budget-revocation-transfer.v1' -or [int]$R.resource_configurations_per_cell-ne 64 -or [int]$R.topology_condition_cells-ne 108 -or [int]$R.summaries.Count-ne 216){throw 'UPLM4M_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_revocation_used -or $R.adaptive_coordinate_used){throw 'UPLM4M_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($C in $R.candidate_coordinates){
   $Rows=@($R.summaries|Where-Object{$_.coordinate-eq $C});$Bad=@($Rows|Where-Object{[int]$_.nonzero_spread_groups-gt 0})
   $Max=0;$Multi=0;foreach($X in $Rows){if([int]$X.max_failure_spread-gt $Max){$Max=[int]$X.max_failure_spread};$Multi+=[int]$X.multi_member_groups}
   Write-Host "coordinate=$C cells=$($Rows.Count) nonzero_cells=$($Bad.Count) max_spread=$Max multi_member_groups=$Multi"
 }
 Write-Host 'WINGLESS_UP277_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
