$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2j-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2j-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2J_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2J_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2J_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2J_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2j-working-set-scheduling -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2J_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2J_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2j-working-set-scheduling)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2J_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2j-working-set-scheduling)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2J_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2J_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2j-working-set-scheduling.v1'){throw 'UPLM2J_SCHEMA'}
 if([int]$R.metrics.Count-ne 200 -or [int]$R.entity_count-ne 24 -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM2J_DESIGN'}
 if([int]$R.store_clauses-ne 24 -or [int]$R.observe_clauses-ne 24 -or [int]$R.report_clauses-ne 24){throw 'UPLM2J_EVENT_COUNT'}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM2J_CAP_BREACH'}}
 if($R.capacity_changed -or $R.adaptive_chunking_used -or $R.extra_training_used -or $R.router_modified -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2J_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "allocation=$($M.allocation) train=$($M.training_schedule) family=$($M.family) working=$($M.working_set_schedule) recall=$($M.recall_hit_rate) expected=$($M.expected_capacity_hit_rate) exact=$($M.report_set_exact_accuracy) route=$($M.event_routing_accuracy)"}
 Write-Host 'WINGLESS_UP172_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
