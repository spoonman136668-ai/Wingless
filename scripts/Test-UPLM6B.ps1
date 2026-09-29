$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR;$GoExe=(Get-Command go -ErrorAction Stop).Source;$GoBin=Split-Path $GoExe -Parent;$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-uplm6b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm6b-gotmp-"+$PID);New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null;$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM6B_GO_ENV_FAILED'};go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM6B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm6b-prepressure-extension-map -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM6B_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM6B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm6b-prepressure-extension-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM6B_PROBE1_FAILED'};$P2=((go run ./cmd/unitary-up-lm6b-prepressure-extension-map)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM6B_NONDETERMINISTIC_OUTPUT'};$R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm6b-prepressure-extension-map.v1' -or [int]$R.rows.Count-ne 108){throw 'UPLM6B_DESIGN'}
 if($R.parent_classification-cne 'PREPRESSURE_CREATES_LENGTH16_FLOOR'){throw 'UPLM6B_PARENT_ANCHOR_DRIFT'}
 if($R.live_activation -or $R.policy_changed){throw 'UPLM6B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';Write-Host "min_writes=$($R.min_writes) max_writes=$($R.max_writes) distinct=$($R.distinct_write_counts) min_post_length=$($R.min_post_length) max_post_length=$($R.max_post_length) classification=$($R.classification)";Write-Host 'WINGLESS_UP382_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$PriorGOROOT;$env:GOCACHE=$PriorGoCache;$env:GOTMPDIR=$PriorGoTmp;Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
