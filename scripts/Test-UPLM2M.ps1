$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2m-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2m-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2M_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2M_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2M_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2M_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2m-parametric-fallback-remap -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2M_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2M_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2m-parametric-fallback-remap)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2M_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2m-parametric-fallback-remap)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2M_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2M_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2m-parametric-fallback-remap.v1'){throw 'UPLM2M_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 160 -or [int]$R.deferred_first_chunk_reports-ne 5 -or [int]$R.value_shifts.Count-ne 4){throw 'UPLM2M_DESIGN'}
 if([int]$R.exact_recall_cap-ne 16 -or $R.capacity_changed -or $R.extra_training_used -or $R.router_modified -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2M_BOUNDARY_LEAK'}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-ne 16){throw 'UPLM2M_CAPACITY'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 $R.metrics|Group-Object value_shift|ForEach-Object{$g=$_.Group;$exact=($g|Where-Object{$_.report_set_exact_accuracy-eq 1}).Count;$missExact=($g|Where-Object{$_.recall_hit_rate-lt 1 -and $_.report_set_exact_accuracy-eq 1}).Count;Write-Host "shift=$($_.Name) cells=$($g.Count) exact=$exact miss_but_exact=$missExact"}
 Write-Host 'WINGLESS_UP175_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
