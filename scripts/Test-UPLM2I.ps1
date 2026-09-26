$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2i-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2i-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2I_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2I_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2I_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2I_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2i-entity-concurrency-capacity -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2I_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2I_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2i-entity-concurrency-capacity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2I_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2i-entity-concurrency-capacity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2I_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2I_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2i-entity-concurrency-capacity.v1'){throw 'UPLM2I_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 500 -or [int]$R.allocations-ne 5 -or [int]$R.training_schedules-ne 2 -or [int]$R.families-ne 5 -or [int]$R.report_orders-ne 2){throw 'UPLM2I_COUNT'}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.entity_name_count-ne 24 -or -not $R.entity_names_use_existing_alphabet){throw 'UPLM2I_DESIGN'}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM2I_RECALL_CAP_BREACH'}}
 if($R.recall_cap_changed -or $R.recurrent_parameters_trained -or $R.router_modified -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2I_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "allocation=$($M.allocation) schedule=$($M.training_schedule) family=$($M.family) order=$($M.report_order) entities=$($M.entities) recall_hit=$($M.recall_hit_rate) expected=$($M.expected_capacity_hit_rate) route=$($M.event_routing_accuracy) dep=$($M.dependent_first_byte_accuracy) maxrecall=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP169_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
