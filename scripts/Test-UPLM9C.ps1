$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGOROOT=$env:GOROOT
$PriorGoCache=$env:GOCACHE
$PriorGoTmp=$env:GOTMPDIR
$GoExe=(Get-Command go -ErrorAction Stop).Source
$GoBin=Split-Path $GoExe -Parent
$GoRoot=Split-Path $GoBin -Parent
$GoCache=Join-Path $env:TEMP ("wingless-uplm9c-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-uplm9c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOROOT=$GoRoot
$env:GOCACHE=$GoCache
$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
    go env -w GOFLAGS=-buildvcs=false
    if($LASTEXITCODE-ne 0){throw 'UPLM9C_GO_ENV_FAILED'}
    go clean -cache -testcache
    if($LASTEXITCODE-ne 0){throw 'UPLM9C_GO_CLEAN_FAILED'}

    go test ./unitary ./cmd/unitary-up-lm9c-profile-family-invariance -count=1
    if($LASTEXITCODE-ne 0){throw 'UPLM9C_FOCUSED_TEST_FAILED'}
    go test ./... -count=1
    if($LASTEXITCODE-ne 0){throw 'UPLM9C_FULL_REGRESSION_FAILED'}

    $P1=((go run ./cmd/unitary-up-lm9c-profile-family-invariance)|Out-String).Trim()
    if($LASTEXITCODE-ne 0){throw 'UPLM9C_PROBE1_FAILED'}
    $P2=((go run ./cmd/unitary-up-lm9c-profile-family-invariance)|Out-String).Trim()
    if($P1-cne $P2){throw 'UPLM9C_NONDETERMINISTIC_OUTPUT'}
    $R=$P1|ConvertFrom-Json

    if($R.schema-cne 'wingless.up-lm9c-profile-family-invariance.v1' -or [int]$R.rows.Count-ne 108){
        throw 'UPLM9C_DESIGN'
    }
    if($R.parent_classification-cne 'PERMUTATION_FAMILY_EXACT_DEFICIT_LAW'){
        throw 'UPLM9C_PARENT_ANCHOR_DRIFT'
    }
    $Profiles=@($R.profiles|ForEach-Object{[string]$_})
    if($Profiles.Count-ne 3 -or $Profiles[0]-cne 'deferred_only' -or $Profiles[1]-cne 'layout_only' -or $Profiles[2]-cne 'hybrid_max'){
        throw 'UPLM9C_PROFILE_FREEZE_DRIFT'
    }
    $Perms=@($R.permutations|ForEach-Object{[string]$_})
    if($Perms.Count-ne 3 -or $Perms[0]-cne 'identity' -or $Perms[1]-cne 'rotate1' -or $Perms[2]-cne 'random_fixed'){
        throw 'UPLM9C_PERMUTATION_DRIFT'
    }
    if((@($R.random_fixed|ForEach-Object{[int]$_}) -join ',') -cne '3,0,5,1,4,2'){
        throw 'UPLM9C_RANDOM_FIXED_DRIFT'
    }
    if($R.live_activation -or $R.policy_changed -or $R.adaptive_selection_used){
        throw 'UPLM9C_BOUNDARY_LEAK'
    }

    Write-Host $P1
    Write-Host ''
    Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
    Write-Host "exact=$($R.exact_law_rows) mismatch=$($R.mismatch_rows) classification=$($R.classification)"
    Write-Host 'WINGLESS_UP392_HARNESS_PASS'
}finally{
    Pop-Location
    if($null -eq $PriorGOROOT){Remove-Item Env:GOROOT -ErrorAction SilentlyContinue}else{$env:GOROOT=$PriorGOROOT}
    if($null -eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
    if($null -eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
    Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
