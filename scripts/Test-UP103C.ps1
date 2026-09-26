$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up103c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up103c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP103C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP103C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP103C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP103C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up103c-promotion-order -count=1;if($LASTEXITCODE-ne 0){throw 'UP103C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP103C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up103c-promotion-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP103C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up103c-promotion-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP103C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP103C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up103c-promotion-order.v1'){throw 'UP103C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 16){throw 'UP103C_HOT_KEYS'}
 if([int]$R.trigger_threshold-ne 8){throw 'UP103C_TRIGGER'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP103C_RECALL_CAP'}
 if($R.extra_diagnostic_reads_used){throw 'UP103C_EXTRA_READS'}
 if($R.adaptive_ordering_used){throw 'UP103C_ADAPTIVE_ORDER'}
 if($R.future_oracle_used){throw 'UP103C_ORACLE'}
 if($R.query_labels_used_for_admission){throw 'UP103C_QUERY_LABEL'}
 if([int]$R.points.Count-ne 9){throw "UP103C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recall_entries_used-gt 16){throw 'UP103C_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "initial=$($M.initial_order) second=$($M.second_write_order) resident=$($M.mean_resident_hot_before_churn) trigger=$($M.trigger_fired_episodes) trigger_idx=$($M.mean_filtered_write_index_at_trigger) fp=$($M.false_positive_churn_admissions) final=$($M.final_hot_query_accuracy) exact=$($M.final_hot_set_exact_accuracy)"}
 Write-Host 'WINGLESS_UP159_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
