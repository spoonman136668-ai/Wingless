$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up100b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up100b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP100B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP100B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP100B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP100B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up100b-prefix-timed-routing -count=1;if($LASTEXITCODE-ne 0){throw 'UP100B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP100B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up100b-prefix-timed-routing)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP100B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up100b-prefix-timed-routing)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP100B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP100B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up100b-prefix-timed-routing.v1'){throw 'UP100B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP100B_STATE_DIM'}
 if([int]$R.epochs-ne 20){throw 'UP100B_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UP100B_LR'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP100B_RECALL_CAP'}
 if($R.explicit_class_at_learned_inference){throw 'UP100B_CLASS_LEAK'}
 if([int]$R.class_metrics.Count-ne 5){throw 'UP100B_CLASS_METRIC_COUNT'}
 if([int]$R.verb_metrics.Count-ne 9){throw 'UP100B_VERB_METRIC_COUNT'}
 if([int]$R.routing_points.Count-ne 24){throw 'UP100B_ROUTING_COUNT'}
 foreach($M in $R.routing_points){if([int]$M.recall_entries_used-gt 16){throw 'UP100B_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.class_metrics){Write-Host "split=$($M.split) label=$($M.label) accuracy=$($M.accuracy) precision=$($M.precision) recall=$($M.recall) examples=$($M.examples)"}
 foreach($M in $R.verb_metrics){Write-Host "verb=$($M.verb) accuracy=$($M.accuracy) examples=$($M.examples)"}
 foreach($M in $R.routing_points){Write-Host "arm=$($M.arm) writes=$($M.total_writes) targets=$($M.target_keys) query=$($M.query_accuracy) exact=$($M.exact_target_set_accuracy) entries=$($M.recall_entries_used)"}
 Write-Host 'WINGLESS_UP117_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
