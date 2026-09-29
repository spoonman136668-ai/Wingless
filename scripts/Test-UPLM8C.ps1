$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT;$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source
$GoRoot=Split-Path (Split-Path $GoExe -Parent) -Parent
$GoCache=Join-Path $env:TEMP ("wingless-uplm8c-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-uplm8c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache,$GoTmp|Out-Null
$env:GOROOT=$GoRoot;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
    go env -w GOFLAGS=-buildvcs=false
    if($LASTEXITCODE-ne0){throw 'UPLM8C_GO_ENV_FAILED'}
    go clean -cache -testcache
    if($LASTEXITCODE-ne0){throw 'UPLM8C_GO_CLEAN_FAILED'}
    go test ./unitary ./cmd/unitary-up-lm8c-heldout-permutation-deficit-law -count=1
    if($LASTEXITCODE-ne0){throw 'UPLM8C_FOCUSED_TEST_FAILED'}
    go test ./... -count=1
    if($LASTEXITCODE-ne0){throw 'UPLM8C_FULL_REGRESSION_FAILED'}
    $P1=((go run ./cmd/unitary-up-lm8c-heldout-permutation-deficit-law)|Out-String).Trim()
    if($LASTEXITCODE-ne0){throw 'UPLM8C_PROBE1_FAILED'}
    $P2=((go run ./cmd/unitary-up-lm8c-heldout-permutation-deficit-law)|Out-String).Trim()
    if($P1-cne $P2){throw 'UPLM8C_NONDETERMINISTIC_OUTPUT'}
    $R=$P1|ConvertFrom-Json
    if($R.schema-cne 'wingless.up-lm8c-heldout-permutation-deficit-law.v1' -or [int]$R.rows.Count-ne108){throw 'UPLM8C_DESIGN'}
    if($R.parent_classification-cne 'HELDOUT_ROTATION_EXACT_DEFICIT_LAW'){throw 'UPLM8C_PARENT_ANCHOR_DRIFT'}
    $Perms=@($R.permutations)
    if($Perms.Count-ne3 -or $Perms[0]-cne'identity' -or $Perms[1]-cne'rotate1' -or $Perms[2]-cne'random_fixed'){throw 'UPLM8C_PERMUTATION_FREEZE_DRIFT'}
    $Fixed=@($R.random_fixed|ForEach-Object{[int]$_})
    if(($Fixed -join ',')-cne'3,0,5,1,4,2'){throw 'UPLM8C_FIXED_PERMUTATION_DRIFT'}
    if($R.live_activation -or $R.policy_changed){throw 'UPLM8C_BOUNDARY_LEAK'}
    Write-Host $P1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "exact=$($R.exact_law_rows) mismatch=$($R.mismatch_rows) classification=$($R.classification)"
    Write-Host 'WINGLESS_UP391_HARNESS_PASS'
}finally{
    Pop-Location
    if($null-eq$PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
    if($null-eq$PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
    if($null-eq$PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
    Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
