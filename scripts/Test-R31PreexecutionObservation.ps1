param([Parameter(Mandatory=$true)][ValidateSet('Wingless','Yggdrasil')][string]$Lane)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$HomeRepo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
foreach($ScriptName in @('scripts/Write-ResearchPreexecutionObservation.ps1','scripts/Invoke-RepoLocalResearchCycle.ps1')){
    $Tokens=$null;$ParseErrors=$null
    [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $HomeRepo $ScriptName),[ref]$Tokens,[ref]$ParseErrors)
    if(@($ParseErrors).Count-ne0){throw "R31_POWERSHELL_PARSE_FAILED $ScriptName"}
}
$Root=Join-Path $env:RUNNER_TEMP ("r31-qual-"+$Lane+"-"+[guid]::NewGuid().ToString('N'))
$Repo=Join-Path $Root 'fresh-local-source'
New-Item -ItemType Directory -Force -Path $Repo | Out-Null
try{
    & git -C $Repo init -b r31-test
    if($LASTEXITCODE-ne0){throw 'INIT_FAILED'}
    & git -C $Repo config user.email 'r31-test@invalid.local'
    & git -C $Repo config user.name 'r31-test'
    [IO.File]::WriteAllText((Join-Path $Repo 'baseline.txt'),"frozen parent")
    & git -C $Repo add baseline.txt
    & git -C $Repo commit -m baseline
    if($LASTEXITCODE-ne0){throw 'PARENT_COMMIT_FAILED'}
    $Parent=([string](& git -C $Repo rev-parse HEAD)).Trim()
    $PreregPath=if($Lane -ceq 'Wingless'){'docs/experiments/r31-test.md'}else{'research/experiments/r31-test.ice'}
    New-Item -ItemType Directory -Force -Path (Split-Path -Parent (Join-Path $Repo $PreregPath))|Out-Null
    [IO.File]::WriteAllText((Join-Path $Repo $PreregPath),"new pre-registered hypothesis")
    & git -C $Repo add -- $PreregPath
    & git -C $Repo commit -m "prereg only"
    if($LASTEXITCODE-ne0){throw 'PREREG_COMMIT_FAILED'}
    $Prereg=([string](& git -C $Repo rev-parse HEAD)).Trim()
    [IO.File]::WriteAllText((Join-Path $Repo 'implementation.txt'),"source only")
    & git -C $Repo add implementation.txt
    & git -C $Repo commit -m "implementation after prereg"
    if($LASTEXITCODE-ne0){throw 'IMPLEMENTATION_COMMIT_FAILED'}
    $Qualified=([string](& git -C $Repo rev-parse HEAD)).Trim()
    $Evidence=Join-Path $Root 'owner-event.json'
    $Stage=if($Lane -ceq 'Wingless'){'PRE_QUALIFICATION'}else{'PRE_SCIENTIFIC_RUN'}
    $Helper=Join-Path $HomeRepo 'scripts/Write-ResearchPreexecutionObservation.ps1'
    & $Helper -Lane $Lane -RepoPath $Repo -ParentSha $Parent -PreregistrationCommit $Prereg -QualifiedHead $Qualified -Experiment 'FUTURE_R31_TEST' -Stage $Stage -RunId 90000001 -RunAttempt 1 -OutputPath $Evidence
    if(-not$?){throw 'TRUE_PATH_FAILED'}
    $Event=Get-Content -Raw -LiteralPath $Evidence|ConvertFrom-Json
    if($Event.schema -cne 'forge.research-owner.r31-genuine-preexecution-observation.v1' -or
       $Event.lane -cne $Lane -or $Event.stage -cne $Stage -or
       $Event.parent_sha -cne $Parent -or $Event.preregistration_commit -cne $Prereg -or
       $Event.qualified_source_head -cne $Qualified -or
       $Event.github_run_id -ne 90000001 -or
       $Event.logged_selection_probability_bp -ne 0 -or
       $Event.optimizer_dispatch_authorized -ne $false -or
       $Event.predecision_signed_by_independent_owner -ne $false -or
       $Event.scientific_improvement_established -ne $false -or
       $Event.original_first_use_claimed_by_this_probe -ne $false){
       throw 'UNTRUSTED_SCIENTIFIC_AUTHORITY_OR_SOURCE_DRIFT'
    }
    $Bytes=[IO.File]::ReadAllBytes($Evidence)
    $Checksum=([BitConverter]::ToString([Security.Cryptography.SHA256]::HashData($Bytes))).Replace('-','').ToLowerInvariant()
    foreach($Case in @('duplicate','wrong_source','retry','wrong_lane')){
        $Rejected=$false;$ShaArg=$Qualified;$AttemptArg=1
        $OutArg=Join-Path $Root ('extra-'+$Case+'.json')
        $LaneArg=$Lane
        if($Case -ceq 'duplicate'){$OutArg=$Evidence}
        if($Case -ceq 'wrong_source'){$ShaArg='f'*40}
        if($Case -ceq 'retry'){$AttemptArg=2}
        if($Case -ceq 'wrong_lane'){$LaneArg=if($Lane -ceq 'Wingless'){'Yggdrasil'}else{'Wingless'}}
        try{
            & $Helper -Lane $LaneArg -RepoPath $Repo -ParentSha $Parent -PreregistrationCommit $Prereg -QualifiedHead $ShaArg -Experiment 'FUTURE_R31_TEST' -Stage $Stage -RunId 90000001 -RunAttempt $AttemptArg -OutputPath $OutArg
        }catch{$Rejected=$true}
        if(-not$Rejected){throw ('R31_NEGATIVE_'+$Case+'_ADMITTED')}
    }
    $After=([BitConverter]::ToString([Security.Cryptography.SHA256]::HashData([IO.File]::ReadAllBytes($Evidence)))).Replace('-','').ToLowerInvariant()
    if($After -cne $Checksum){throw 'SOURCE_RECORD_MUTATED_ON_NEGATIVE'}
    Write-Host "R31_SOURCE_ONLY_OWNER_EXECUTION_BOUNDARY_PASS lane=$Lane negative=4"
}finally{
    Remove-Item -LiteralPath $Root -Recurse -Force -ErrorAction SilentlyContinue
}
