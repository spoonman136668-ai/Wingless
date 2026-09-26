$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up149c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up149c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP149C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP149C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP149C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP149C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up149c-hand-rotation-causality -count=1;if($LASTEXITCODE-ne 0){throw 'UP149C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP149C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up149c-hand-rotation-causality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP149C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up149c-hand-rotation-causality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP149C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP149C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up149c-hand-rotation-causality.v1'){throw 'UP149C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 16 -or [int]$R.initial_hands.Count-ne 4){throw 'UP149C_DESIGN'}
 foreach($P in $R.points){if(([int[]]$P.target_slots -join ',') -ne '0,1,2,3' -or [int]$P.recall_entries_used-ne 16){throw 'UP149C_POINT'}}
 if($R.memory_policy_changed -or $R.adaptive_refresh_used -or $R.capacity_increased){throw 'UP149C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "hand=$($P.initial_hand) refresh=$($P.refresh_before_admission) present=$($P.present_before_refresh) retained=$($P.target_final_retained) final_hand=$($P.final_hand)"}
 Write-Host 'WINGLESS_UP149_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
