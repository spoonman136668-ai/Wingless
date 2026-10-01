Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$Sidecars = @(
    [pscustomobject]@{
        Name='wing'
        Root='C:\ProgramData\CKBR\research-sidecar'
        Source='C:\ProgramData\CKBR\codex\work\rs\source\Wingless'
        Task='CKBPlane Research Sidecar'
    },
    [pscustomobject]@{
        Name='ygg'
        Root='C:\ProgramData\CKBR\research-sidecar-yggdrasil'
        Source='C:\ProgramData\CKBR\research-sidecar-yggdrasil\source\Yggdrasil'
        Task='CKBPlane Research Sidecar Yggdrasil'
    }
)

foreach($S in $Sidecars){
    Write-Host "=== SIDECAR $($S.Name) ==="
    $Cfg=Join-Path $S.Root 'config\autonomy.yaml'
    $Exe=Join-Path $S.Root 'bin\research-sidecar.exe'
    $Task=Get-ScheduledTask -TaskName $S.Task -ErrorAction SilentlyContinue
    Write-Host "TASK state=$(if($null -eq $Task){'MISSING'}else{[string]$Task.State})"

    if((Test-Path -LiteralPath $Exe -PathType Leaf) -and (Test-Path -LiteralPath $Cfg -PathType Leaf)){
        $Raw=& $Exe -mode status -config $Cfg 2>$null
        if($LASTEXITCODE -eq 0){
            $State=(($Raw|ForEach-Object{[string]$_}) -join [Environment]::NewLine)|ConvertFrom-Json
            Write-Host "STATUS status=$($State.status) cycle_count=$($State.cycle_count) blocked=$($State.blocked_cycle_count) last_cycle_id=$($State.last_cycle_id) last_cycle_status=$($State.last_cycle_status)"
        } else {
            Write-Host "STATUS exit=$LASTEXITCODE"
        }
    }

    if(-not(Test-Path -LiteralPath $S.Source -PathType Container)){
        Write-Host "SOURCE missing=$($S.Source)"
        continue
    }
    Write-Host "SOURCE path=$($S.Source)"
    $Head=(& git -C $S.Source rev-parse HEAD).Trim()
    if($LASTEXITCODE -ne 0){throw "$($S.Name) source is not a git repository"}
    $Tree=(& git -C $S.Source rev-parse 'HEAD^{tree}').Trim()
    $BranchRaw=(& git -C $S.Source branch --show-current 2>$null | Select-Object -First 1)
    $Branch=if($null -eq $BranchRaw){''}else{([string]$BranchRaw).Trim()}
    $RemoteRaw=(& git -C $S.Source remote get-url origin 2>$null | Select-Object -First 1)
    $Remote=if($null -eq $RemoteRaw){''}else{([string]$RemoteRaw).Trim()}
    Write-Host "SOURCE head=$Head tree=$Tree branch=$Branch remote=$Remote"
    $Status=((& git -C $S.Source status --porcelain) -join ';')
    Write-Host "SOURCE dirty=$([bool](-not[string]::IsNullOrWhiteSpace($Status)))"
    & git -C $S.Source log -5 --pretty='format:LOG %H %ct %s'
}

Write-Host '=== CODEX ==='
$Candidates=@(
    'C:\ProgramData\CKBR\codex\bin\codex.cmd',
    'C:\ProgramData\CKBR\codex\bin\codex.exe',
    'C:\ProgramData\CKBR\research-sidecar\codex\bin\codex.cmd',
    'C:\ProgramData\CKBR\research-sidecar\codex\bin\codex.exe',
    'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.cmd',
    'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe'
)
foreach($P in $Candidates){
    if(Test-Path -LiteralPath $P -PathType Leaf){
        Write-Host "CODEX path=$P"
        try { & $P --version } catch { Write-Host "CODEX version_error=$($_.Exception.Message)" }
    }
}
$Home='C:\ProgramData\CKBR\codex\home'
Write-Host "CODEX_HOME exists=$(Test-Path -LiteralPath $Home -PathType Container)"
if(Test-Path -LiteralPath $Home -PathType Container){
    foreach($Name in @('.codex','auth.json','config.toml')){
        $P=Join-Path $Home $Name
        Write-Host "CODEX_HOME_ENTRY name=$Name exists=$(Test-Path -LiteralPath $P)"
    }
}

Write-Host 'REPO_LOCAL_CUTOVER_PROBE=PASS'
