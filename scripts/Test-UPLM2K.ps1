$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2k-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2k-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2K_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2K_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2K_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2K_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2k-dependency-closure-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2K_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2K_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2k-dependency-closure-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2K_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2k-dependency-closure-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2K_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2K_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2k-dependency-closure-boundary.v1'){throw 'UPLM2K_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 300 -or [int]$R.deferred_levels.Count-ne 6 -or [int]$R.entity_count-ne 24){throw 'UPLM2K_DESIGN'}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.chunk_sizes.Count-ne 2){throw 'UPLM2K_MEMORY'}
 foreach($M in $R.metrics){if([math]::Abs([double]$M.recall_hit_rate-[double]$M.expected_recall_hit_rate)-gt 1e-12){throw 'UPLM2K_EXPECTED_RECALL_MISMATCH'}}
 if(-not $R.event_multiset_identical -or $R.capacity_changed -or $R.adaptive_chunking_used -or $R.extra_training_used -or $R.router_modified -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2K_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($L in $R.deferred_levels){$X=@($R.metrics|Where-Object{[int]$_.deferred_first_chunk_reports-eq [int]$L});Write-Host "deferred=$L recall=$($X[0].recall_hit_rate) exact_min=$((@($X|ForEach-Object{$_.report_set_exact_accuracy})|Measure-Object -Minimum).Minimum) max_recall=$((@($X|ForEach-Object{$_.max_recall_entries})|Measure-Object -Maximum).Maximum)"}
 Write-Host 'WINGLESS_UP162_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
