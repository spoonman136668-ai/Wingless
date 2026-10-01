Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$Sidecars = @(
    [pscustomobject]@{
        Name='wing'
        Root='C:\ProgramData\CKBR\research-sidecar'
        Source='C:\ProgramData\CKBR\codex\work\rs\source\Wingless'
        Task='CKBPlane Research Sidecar'
        SearchRoots=@('C:\ProgramData\CKBR\codex\work\rs','C:\ProgramData\CKBR\research-sidecar')
    },
    [pscustomobject]@{
        Name='ygg'
        Root='C:\ProgramData\CKBR\research-sidecar-yggdrasil'
        Source='C:\ProgramData\CKBR\research-sidecar-yggdrasil\source\Yggdrasil'
        Task='CKBPlane Research Sidecar Yggdrasil'
        SearchRoots=@('C:\ProgramData\CKBR\research-sidecar-yggdrasil')
    }
)

function Get-GitInfo {
    param([string]$Repo,[string]$Label)
    $HeadRaw=(& git -C $Repo rev-parse HEAD 2>$null | Select-Object -First 1)
    if($LASTEXITCODE -ne 0 -or $null -eq $HeadRaw){return}
    $Head=([string]$HeadRaw).Trim()
    $TreeRaw=(& git -C $Repo rev-parse 'HEAD^{tree}' 2>$null | Select-Object -First 1)
    $Tree=if($null-eq$TreeRaw){''}else{([string]$TreeRaw).Trim()}
    $BranchRaw=(& git -C $Repo branch --show-current 2>$null | Select-Object -First 1)
    $Branch=if($null-eq$BranchRaw){''}else{([string]$BranchRaw).Trim()}
    $RemoteRaw=(& git -C $Repo remote get-url origin 2>$null | Select-Object -First 1)
    $Remote=if($null-eq$RemoteRaw){''}else{([string]$RemoteRaw).Trim()}
    $Dirty=((& git -C $Repo status --porcelain 2>$null)-join';')
    $LogRaw=(& git -C $Repo log -1 --pretty='format:%ct|%H|%s' 2>$null | Select-Object -First 1)
    $Log=if($null-eq$LogRaw){''}else{[string]$LogRaw}
    Write-Host "GITROOT label=$Label path=$Repo head=$Head tree=$Tree branch=$Branch dirty=$([bool](-not[string]::IsNullOrWhiteSpace($Dirty))) remote=$Remote latest=$Log"
}

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
        }
    }

    if(Test-Path -LiteralPath $S.Source -PathType Container){
        Get-GitInfo -Repo $S.Source -Label "$($S.Name)-source"
        & git -C $S.Source log -8 --pretty='format:SOURCELOG %ct %H %s'
    }

    $Repos=New-Object Collections.Generic.HashSet[string] ([StringComparer]::OrdinalIgnoreCase)
    foreach($Root in $S.SearchRoots){
        if(-not(Test-Path -LiteralPath $Root -PathType Container)){continue}
        foreach($Mark in @(Get-ChildItem -LiteralPath $Root -Force -Recurse -ErrorAction SilentlyContinue | Where-Object {$_.Name -eq '.git'})){
            $Repo=if($Mark.PSIsContainer){$Mark.Parent.FullName}else{$Mark.DirectoryName}
            [void]$Repos.Add($Repo)
        }
    }
    foreach($Repo in ($Repos | Sort-Object)){
        Get-GitInfo -Repo $Repo -Label "$($S.Name)-scan"
    }
}

Write-Host '=== CODEX ==='
foreach($P in @(
    'C:\ProgramData\CKBR\codex\bin\codex.exe',
    'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe'
)){
    if(Test-Path -LiteralPath $P -PathType Leaf){
        Write-Host "CODEX path=$P version=$((& $P --version 2>$null)-join' ')"
    }
}

Write-Host 'REPO_LOCAL_CUTOVER_PROBE=PASS'
