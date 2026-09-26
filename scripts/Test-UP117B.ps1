$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up117b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up117b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP117B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP117B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP117B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP117B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up117b-stores-geometry -count=1;if($LASTEXITCODE-ne 0){throw 'UP117B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP117B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up117b-stores-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP117B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up117b-stores-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP117B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP117B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up117b-stores-geometry.v1'){throw 'UP117B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP117B_STATE_DIM'}
 if($R.training_changed){throw 'UP117B_TRAINING_CHANGED'}
 if([int]$R.static_relations.Count-ne 28){throw "UP117B_STATIC_COUNT $($R.static_relations.Count)"}
 if([int]$R.dynamic_points.Count-ne 168){throw "UP117B_DYNAMIC_COUNT $($R.dynamic_points.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.static_relations){Write-Host "static source=$($M.source_verb) target=$($M.target) cosine=$($M.cosine)"}
 foreach($M in $R.dynamic_points){if($M.verb-ceq 'stores' -or $M.verb-ceq 'keeps'){Write-Host "dynamic order=$($M.class_order) policy=$($M.policy) stage=$($M.stage) verb=$($M.verb) acc=$($M.accuracy) margin=$($M.mean_store_minus_observe_logit_margin) min=$($M.min_store_minus_observe_logit_margin) dot=$($M.mean_store_minus_observe_dot_contribution) bias=$($M.store_minus_observe_bias_contribution)"}}
 Write-Host 'WINGLESS_UP156_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
