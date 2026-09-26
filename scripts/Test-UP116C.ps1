$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up116c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up116c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP116C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP116C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP116C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP116C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up116c-exact3-generality -count=1;if($LASTEXITCODE-ne 0){throw 'UP116C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP116C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up116c-exact3-generality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP116C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up116c-exact3-generality)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP116C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP116C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up116c-exact3-generality.v1'){throw 'UP116C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 12 -or [int]$R.exact_recall_cap-ne 16){throw 'UP116C_CAPACITY'}
 if([int]$R.generation_interval-ne 32 -or [int]$R.admission_sightings-ne 3){throw 'UP116C_RECURRENCE'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.episodes_per_seed-ne 32){throw 'UP116C_CONFIG'}
 if($R.query_evidence_used_for_admission -or $R.semantic_priority_used -or $R.future_oracle_used -or $R.phase_label_used){throw 'UP116C_LABEL_LEAK'}
 if([int]$R.points.Count-ne 2){throw "UP116C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recall_entries_used-gt 16){throw 'UP116C_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "scenario=$($M.scenario) weak_admit=$($M.weak_admission_rate) strong_admit=$($M.strong_admission_rate) distractor_admit=$($M.distractor_admission_rate) hot=$($M.hot_accuracy) strong=$($M.strong_accuracy) distractor=$($M.distractor_accuracy) target=$($M.target16_accuracy) exact=$($M.target16_exact_accuracy) candidates=$($M.all_candidate_accuracy) fp=$($M.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP176_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
