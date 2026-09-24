$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up50b-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up50b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-50B PAIR05 GENERALIZATION ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP50B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP50B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP50B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP50B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up50b-pair05-generalization -count=1;if($LASTEXITCODE-ne 0){throw 'UP50B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP50B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up50b-pair05-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP50B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up50b-pair05-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP50B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP50B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up50b-pair05-generalization.v1'){throw 'UP50B_SCHEMA_MISMATCH'}
 if($R.selection_run){throw 'UP50B_SELECTION_PRESENT'}
 if([int]$R.cases.Count-ne 4){throw 'UP50B_CASE_COUNT'}
 if([int]$R.frozen_pair[0]-ne 0 -or [int]$R.frozen_pair[1]-ne 5){throw 'UP50B_PAIR_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "All cases pass: $($R.all_cases_pass)"
 foreach($C in $R.cases){
  Write-Host "$($C.spec.name) noise=$($C.spec.memory_noise) gate=$($C.hard_gate) held=$($C.selected_hard.static.held_out_accuracy) commit=$($C.selected_hard.mutable_integration.commit_decode_accuracy) final=$($C.selected_hard.mutable_integration.exact_final_table_accuracy) relation=$($C.selected_hard.mutable_integration.relational_query_accuracy)"
 }
 Write-Host 'WINGLESS_UP50_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
