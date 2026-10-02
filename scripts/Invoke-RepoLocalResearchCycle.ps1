param(
    [Parameter(Mandatory=$true)][string]$RepoPath,
    [Parameter(Mandatory=$true)][string]$Repository,
    [Parameter(Mandatory=$true)][string]$ActiveBranch,
    [Parameter(Mandatory=$true)][string]$NorthStarPath,
    [Parameter(Mandatory=$true)][string]$NorthStarSha256,
    [string]$Trigger = 'manual',
    [string]$ResultPath = ''
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Get-CanonicalTextSha256([string]$Path) {
    $Text=[IO.File]::ReadAllText($Path)
    $Normalized=$Text.Replace("`r`n","`n").Replace("`r","`n")
    $Utf8NoBom=New-Object Text.UTF8Encoding($false)
    $Bytes=$Utf8NoBom.GetBytes($Normalized)
    $Sha=[Security.Cryptography.SHA256]::Create()
    try {
        return ([BitConverter]::ToString($Sha.ComputeHash($Bytes))).Replace('-','').ToLowerInvariant()
    } finally {
        $Sha.Dispose()
    }
}

function Invoke-Git {
    param([Parameter(ValueFromRemainingArguments=$true)][string[]]$Args)
    & git -C $RepoPath @Args
    if($LASTEXITCODE -ne 0){throw "GIT_FAILED exit=$LASTEXITCODE args=$($Args -join ' ')"}
}

function Get-CommitFiles([string]$Commit){
    $Files=@(& git -C $RepoPath diff-tree --no-commit-id --name-only -r $Commit)
    if($LASTEXITCODE-ne0){throw "DIFF_TREE_FAILED commit=$Commit"}
    return @($Files|Where-Object{$_ -and $_.Trim()})
}

function Write-CycleResult([bool]$Advanced,[string]$Reason,[string]$Head){
    if([string]::IsNullOrWhiteSpace($ResultPath)){return}
    $Obj=[ordered]@{
        schema='research.repo-local-cycle-result.v1'
        program='Wingless'
        advanced=$Advanced
        reason=$Reason
        head_sha=$Head
        generated_at_utc=[DateTime]::UtcNow.ToString('o')
    }
    $Dir=Split-Path -Parent $ResultPath
    if($Dir){New-Item -ItemType Directory -Force -Path $Dir|Out-Null}
    $Obj|ConvertTo-Json -Depth 20|Set-Content -LiteralPath $ResultPath -Encoding UTF8
}

function Test-ProviderAvailabilityFailure([string]$Message){
    if([string]::IsNullOrWhiteSpace($Message)){return $false}
    return [bool]($Message -match '(?i)(status=(?:408|429|500|502|503|504)\b|\b(?:408|429|500|502|503|504)\b|too many requests|rate.?limit|timed?\s*out|timeout|underlying connection was closed|unexpected error occurred on a receive|connection (?:reset|aborted)|internal server error|bad gateway|server unavailable|service unavailable|gateway timeout)')
}

function Resolve-Codex {
    foreach($P in @(
        'C:\ProgramData\CKBR\codex\bin\codex.exe',
        'C:\ProgramData\CKBR\research-sidecar\codex\bin\codex.exe',
        'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe'
    )){
        if(Test-Path -LiteralPath $P -PathType Leaf){return $P}
    }
    $Cmd=Get-Command codex -ErrorAction SilentlyContinue
    if($Cmd){return $Cmd.Source}
    throw 'CODEX_NOT_AVAILABLE'
}

function Get-OptionalMindContext([int]$MaxChars=12000){
    $Candidates=New-Object Collections.Generic.List[string]
    $RepoContext=Join-Path $RepoPath '.research-autonomy\mind-palace-context.txt'
    $Candidates.Add($RepoContext)
    if(-not[string]::IsNullOrWhiteSpace($env:MIND_PALACE_CONTEXT_PATH)){$Candidates.Add($env:MIND_PALACE_CONTEXT_PATH)}
    foreach($P in @(
        'C:\ProgramData\CKBR\mind-palace\contexts\Wingless.txt',
        'C:\ProgramData\CKBR\mind-palace\Wingless-context.txt'
    )){$Candidates.Add($P)}
    foreach($P in $Candidates){
        if(-not(Test-Path -LiteralPath $P -PathType Leaf)){continue}
        $Text=[IO.File]::ReadAllText($P)
        if($Text.Length-gt$MaxChars){$Text=$Text.Substring(0,$MaxChars)}
        Write-Host "MIND_PALACE_CONTEXT_USED path=$P chars=$($Text.Length)"
        return $Text
    }
    Write-Host 'MIND_PALACE_CONTEXT_UNAVAILABLE_FALLBACK=repository-lineage'
    return ''
}

if(-not(Test-Path -LiteralPath $RepoPath -PathType Container)){throw "REPO_PATH_MISSING path=$RepoPath"}
if([string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)){throw 'GITHUB_TOKEN_MISSING'}

Push-Location $RepoPath
try{
    $Branch=(& git branch --show-current 2>$null|Select-Object -First 1)
    if($null-eq$Branch){$Branch=''}
    $Branch=([string]$Branch).Trim()
    if($Branch-cne$ActiveBranch){throw "ACTIVE_BRANCH_MISMATCH expected=$ActiveBranch actual=$Branch"}

    $StartSha=(& git rev-parse HEAD).Trim()
    if($LASTEXITCODE-ne0){throw 'START_HEAD_READ_FAILED'}
    if(((& git status --porcelain)-join'').Trim()){throw 'START_WORKTREE_DIRTY'}

    $NorthStarFull=Join-Path $RepoPath $NorthStarPath
    if(-not(Test-Path -LiteralPath $NorthStarFull -PathType Leaf)){throw "NORTH_STAR_MISSING path=$NorthStarPath"}
    $NorthStarActual=Get-CanonicalTextSha256 $NorthStarFull
    if($NorthStarActual-cne$NorthStarSha256.ToLowerInvariant()){throw "NORTH_STAR_IDENTITY_DRIFT expected=$NorthStarSha256 actual=$NorthStarActual"}

    $StatePath=Join-Path $RepoPath '.research-autonomy\state.json'
    if($Trigger-ceq'schedule' -and (Test-Path -LiteralPath $StatePath -PathType Leaf)){
        try{
            $State=Get-Content -Raw -LiteralPath $StatePath|ConvertFrom-Json
            $Last=[DateTime]::Parse([string]$State.last_completed_utc).ToUniversalTime()
            $Age=([DateTime]::UtcNow-$Last).TotalMinutes
            if($Age-lt45){
                Write-Host "RESEARCH_AUTONOMY_WATCHDOG_SKIP age_minutes=$([math]::Round($Age,1))"
                Write-CycleResult -Advanced $false -Reason 'watchdog-recent-progress' -Head $StartSha
                return
            }
        }catch{
            Write-Host "RESEARCH_AUTONOMY_STATE_PARSE_WARNING $($_.Exception.Message)"
        }
    }

    git config user.name 'repo-local-research-autonomy'
    git config user.email 'repo-local-research-autonomy@users.noreply.github.com'

    $RunRoot=Join-Path $env:RUNNER_TEMP ("wingless-repo-local-"+$env:GITHUB_RUN_ID+"-"+$env:GITHUB_RUN_ATTEMPT)
    if(Test-Path -LiteralPath $RunRoot){Remove-Item -LiteralPath $RunRoot -Recurse -Force}
    New-Item -ItemType Directory -Force -Path $RunRoot|Out-Null

    $MindContext=Get-OptionalMindContext
    $PromptPath=Join-Path $RunRoot 'agent-prompt.txt'
    $LastMessagePath=Join-Path $RunRoot 'agent-last-message.txt'
    $Prompt=@"
You are the bounded repo-local autonomous research worker for Wingless.

AUTHORITY AND SCOPE
- You may modify only this Wingless research checkout.
- CKB-plane is governance-only and is NOT part of this execution cycle.
- Never modify CKB, ckb-plane, KTRADE, accepted refs, production systems, broker/live interfaces, credentials, runner configuration, or authority configuration.
- Do not push to GitHub. The deterministic wrapper validates, qualifies, records evidence, and pushes.
- Do not create another scheduler, queue consumer, daemon, or retry loop.

FROZEN PROGRAM IDENTITY
- Active branch: $ActiveBranch
- Exact parent SHA for this cycle: $StartSha
- North Star path: $NorthStarPath
- North Star SHA-256: $NorthStarSha256

SCIENTIFIC CONTINUATION
1. Read the North Star, the complete current experiment/preregistration lineage, .research-autonomy evidence if present, and the current .wingless/qualification-request.json.
2. Treat supported, mixed, null, and negative scientific outcomes as evidence. Do not weaken an assertion or retune a threshold to make a result pass.
3. Derive exactly ONE bounded next question that materially reduces uncertainty toward the North Star.
4. Before implementation, create exactly ONE new preregistration under docs/experiments/. The preregistration must freeze the question/hypothesis, exact parent SHA, arms/inputs, seeds, budgets, metrics, thresholds, controls, validity criteria, classification rules, stop conditions, and a no-post-result-tuning rule.
5. Commit that preregistration ALONE as the first commit after $StartSha.
6. Then implement only the bounded experiment. Keep changes inside docs/experiments/, unitary/, cmd/, scripts/, and .wingless/qualification-request.json.
7. The final implementation commit must update .wingless/qualification-request.json using schema wingless.research-qualification-request.v1. Its branch MUST be $ActiveBranch and baseline_sha MUST be $StartSha. The test_script must be a repo-relative scripts/*.ps1 path.
8. Preserve deterministic controls, fixed resources, disjoint seeds where required, exact negative-result cardinality, and no post-result threshold/budget/capacity tuning.
9. Make all required commits locally with concise research commit messages. Leave the worktree clean.
10. If implementation becomes invalid after any output has been observed, do NOT repair or rerun it under the same preregistration. Restore every implementation/request change, leave exactly the preregistration commit with a clean worktree, make no scientific claim, and end your final response with a line containing exactly RESEARCH_TERMINAL=PREREG_BLOCKED followed by a concise reason.

If the evidence is insufficient, provenance is incomplete, the next question would widen authority, or safe continuation is ambiguous, make no changes and exit nonzero with a clear blocker.

OPTIONAL ADVISORY MIND-PALACE CONTEXT
$MindContext
"@
    [IO.File]::WriteAllText($PromptPath,$Prompt,(New-Object Text.UTF8Encoding($false)))

    $Succeeded=$false
    $PrimaryAvailabilityFailure=$false
    $SecondaryAvailabilityFailure=$false
    $NemotronAgent=Join-Path $PSScriptRoot 'Invoke-NemotronRepoAgent.ps1'
    if(Test-Path -LiteralPath $NemotronAgent -PathType Leaf){
        try{
            Write-Host "RESEARCH_AGENT_ATTEMPT provider=nvidia model=z-ai/glm-5.3"
            & $NemotronAgent -RepoPath $RepoPath -Program 'Wingless' -ActiveBranch $ActiveBranch -StartSha $StartSha -PromptPath $PromptPath -LastMessagePath $LastMessagePath -Provider nvidia -Model 'z-ai/glm-5.3'
            $Succeeded=$true
            Write-Host 'RESEARCH_AGENT_PRIMARY_PASS provider=nvidia model=z-ai/glm-5.3'
        }catch{
            $PrimaryError=$_.Exception.Message
            $PrimaryAvailabilityFailure=Test-ProviderAvailabilityFailure $PrimaryError
            Write-Host "RESEARCH_AGENT_PRIMARY_FAILURE_CLASS availability=$([string]$PrimaryAvailabilityFailure.ToString().ToLowerInvariant())"
            Write-Host "RESEARCH_AGENT_PRIMARY_FAILED provider=nvidia model=z-ai/glm-5.3 error=$PrimaryError"
            & git reset --hard $StartSha|Out-Null
            & git clean -fd|Out-Null
        }

        if(-not$Succeeded){
            try{
                Write-Host "RESEARCH_AGENT_ATTEMPT provider=openrouter model=nvidia/nemotron-3-ultra-550b-a55b:free"
                & $NemotronAgent -RepoPath $RepoPath -Program 'Wingless' -ActiveBranch $ActiveBranch -StartSha $StartSha -PromptPath $PromptPath -LastMessagePath $LastMessagePath -Provider openrouter -Model 'nvidia/nemotron-3-ultra-550b-a55b:free'
                $Succeeded=$true
                Write-Host 'RESEARCH_AGENT_SECONDARY_PASS provider=openrouter model=nvidia/nemotron-3-ultra-550b-a55b:free'
            }catch{
                $SecondaryError=$_.Exception.Message
                $SecondaryAvailabilityFailure=Test-ProviderAvailabilityFailure $SecondaryError
                Write-Host "RESEARCH_AGENT_SECONDARY_FAILURE_CLASS availability=$([string]$SecondaryAvailabilityFailure.ToString().ToLowerInvariant())"
                Write-Host "RESEARCH_AGENT_SECONDARY_FAILED provider=openrouter model=nvidia/nemotron-3-ultra-550b-a55b:free error=$SecondaryError"
                & git reset --hard $StartSha|Out-Null
                & git clean -fd|Out-Null
            }
        }
    }else{
        Write-Host "RESEARCH_AGENT_PRIMARY_UNAVAILABLE path=$NemotronAgent"
    }

    if(-not$Succeeded -and $PrimaryAvailabilityFailure -and $SecondaryAvailabilityFailure){
        & git reset --hard $StartSha|Out-Null
        & git clean -fd|Out-Null
        $DeferredHead=(& git rev-parse HEAD).Trim()
        if($DeferredHead-cne$StartSha){throw "PROVIDER_DEFER_HEAD_DRIFT expected=$StartSha actual=$DeferredHead"}
        if(((& git status --porcelain)-join'').Trim()){throw 'PROVIDER_DEFER_DIRTY_WORKTREE'}
        Write-CycleResult -Advanced $false -Reason 'provider-deferred' -Head $DeferredHead
        Write-Host 'REPO_LOCAL_PROVIDER_DEFERRED free_providers_unavailable=true codex_invoked=false'
        Write-Host "head_sha=$DeferredHead"
        return
    }

    if(-not$Succeeded){
        $Codex=Resolve-Codex
        $OldCodeHome=$env:CODEX_HOME
        $OldHome=$env:HOME
        $OldProfile=$env:USERPROFILE
        try{
            $CodexHome='C:\ProgramData\CKBR\codex\home'
            if(Test-Path -LiteralPath $CodexHome -PathType Container){
                $env:CODEX_HOME=$CodexHome
                $env:HOME=$CodexHome
                $env:USERPROFILE=$CodexHome
            }

            $Models=@('gpt-5.6-sol','gpt-5.6-luna')
            foreach($Model in $Models){
                Write-Host "RESEARCH_AGENT_ATTEMPT provider=codex model=$Model"
                $Args=@('exec','--approve-for-me','--model',$Model,'--output-last-message',$LastMessagePath,'-')
                $Proc=Start-Process -FilePath $Codex -ArgumentList $Args -WorkingDirectory $RepoPath -RedirectStandardInput $PromptPath -NoNewWindow -Wait -PassThru
                if($Proc.ExitCode-eq0){$Succeeded=$true;break}
                Write-Host "RESEARCH_AGENT_RETRY provider=codex model=$Model exit=$($Proc.ExitCode)"
                & git reset --hard $StartSha|Out-Null
                & git clean -fd|Out-Null
            }
        }finally{
            $env:CODEX_HOME=$OldCodeHome
            $env:HOME=$OldHome
            $env:USERPROFILE=$OldProfile
        }
    }
    if(-not$Succeeded){
        throw 'RESEARCH_AGENT_FAILED_ALL_PROVIDERS'
    }

    if(((& git status --porcelain)-join'').Trim()){throw 'AGENT_LEFT_UNCOMMITTED_CHANGES'}
    $BranchAfter=((& git branch --show-current 2>$null|Select-Object -First 1)|Out-String).Trim()
    if($BranchAfter-cne$ActiveBranch){throw "AGENT_BRANCH_CHANGED expected=$ActiveBranch actual=$BranchAfter"}

    $Commits=@(& git rev-list --reverse "$StartSha..HEAD"|Where-Object{$_ -and $_.Trim()})
    if($LASTEXITCODE-ne0){throw 'REV_LIST_FAILED'}
    if($Commits.Count-lt1){throw "CONTINUATION_COMMIT_COUNT_TOO_SMALL count=$($Commits.Count)"}

    $AllFiles=New-Object Collections.Generic.List[string]
    foreach($Commit in $Commits){foreach($File in (Get-CommitFiles $Commit)){$AllFiles.Add($File)}}
    foreach($File in $AllFiles){
        if($File-match'(?i)(^|/)(ckb|ckb-plane|ktrade)(/|$)' -or $File-match'(?i)(coinbase|broker|credential|accepted-ref)'){throw "FORBIDDEN_PATH_CHANGE path=$File"}
        $Allowed=$File-eq'.wingless/qualification-request.json' -or $File.StartsWith('docs/experiments/') -or $File.StartsWith('unitary/') -or $File.StartsWith('cmd/') -or $File.StartsWith('scripts/')
        if(-not$Allowed){throw "CHANGE_OUTSIDE_RESEARCH_POLICY path=$File"}
    }

    $FirstFiles=@(Get-CommitFiles $Commits[0])
    if($FirstFiles.Count-ne1 -or -not$FirstFiles[0].StartsWith('docs/experiments/')){throw "PREREG_NOT_FIRST_AND_ALONE files=$($FirstFiles -join ',')"}

    if($Commits.Count-eq1){
        if(-not(Test-Path -LiteralPath $LastMessagePath -PathType Leaf)){throw 'PREREG_BLOCKED_AGENT_SUMMARY_MISSING'}
        $AgentSummary=[IO.File]::ReadAllText($LastMessagePath)
        if($AgentSummary -notmatch '(?m)^RESEARCH_TERMINAL=PREREG_BLOCKED\s*$'){throw 'PREREG_ONLY_WITHOUT_BLOCKED_MARKER'}
        $PreregCommit=[string]$Commits[0]
        $PreregPath=[string]$FirstFiles[0]
        $Experiment=([IO.Path]::GetFileNameWithoutExtension($PreregPath)).ToUpperInvariant()
        $EvidenceDir=Join-Path $RepoPath '.research-autonomy\evidence'
        New-Item -ItemType Directory -Force -Path $EvidenceDir|Out-Null
        $EvidencePath=Join-Path $EvidenceDir ("blocked-"+$env:GITHUB_RUN_ID+".json")
        $Record=[ordered]@{
            schema='research.preregistered-blocker.v1'
            program='Wingless'
            parent_sha=$StartSha
            preregistration_commit=$PreregCommit
            preregistration_path=$PreregPath
            experiment=$Experiment
            classification='preregistered-blocked'
            scientific_claim=$false
            agent_summary=$AgentSummary.Trim()
            github_run_id=$env:GITHUB_RUN_ID
            generated_at_utc=[DateTime]::UtcNow.ToString('o')
            authority='research-only'
            ckb_plane_role='governance-only'
        }
        $Record|ConvertTo-Json -Depth 20|Set-Content -LiteralPath $EvidencePath -Encoding UTF8
        $State=[ordered]@{
            schema='research.repo-local-state.v1'
            program='Wingless'
            active_branch=$ActiveBranch
            north_star_path=$NorthStarPath
            north_star_sha256=$NorthStarSha256
            parent_sha=$StartSha
            preregistration_commit=$PreregCommit
            experiment=$Experiment
            classification='preregistered-blocked'
            scientific_claim=$false
            last_completed_utc=[DateTime]::UtcNow.ToString('o')
            github_run_id=$env:GITHUB_RUN_ID
            authority='research-only'
            ckb_plane_role='governance-only'
        }
        $State|ConvertTo-Json -Depth 30|Set-Content -LiteralPath $StatePath -Encoding UTF8
        Invoke-Git add '.research-autonomy'
        Invoke-Git commit -m ("research: seal blocked preregistration "+$Experiment)
        $EvidenceHead=(& git rev-parse HEAD).Trim()
        Invoke-Git push origin ("HEAD:refs/heads/"+$ActiveBranch)
        Write-CycleResult -Advanced $true -Reason 'preregistered-blocked-sealed' -Head $EvidenceHead
        Write-Host 'REPO_LOCAL_PREREG_BLOCKED_SEALED'
        Write-Host "experiment=$Experiment"
        Write-Host "preregistration_commit=$PreregCommit"
        Write-Host "evidence_head=$EvidenceHead"
        return
    }

    $RequestPath=Join-Path $RepoPath '.wingless\qualification-request.json'
    if(-not(Test-Path -LiteralPath $RequestPath -PathType Leaf)){throw 'QUALIFICATION_REQUEST_MISSING'}
    $Request=Get-Content -Raw -LiteralPath $RequestPath|ConvertFrom-Json
    if([string]$Request.schema-cne'wingless.research-qualification-request.v1'){throw "QUALIFICATION_SCHEMA_MISMATCH schema=$($Request.schema)"}
    if([string]$Request.branch-cne$ActiveBranch){throw "QUALIFICATION_BRANCH_MISMATCH expected=$ActiveBranch actual=$($Request.branch)"}
    if([string]$Request.baseline_sha-cne$StartSha){throw "QUALIFICATION_BASELINE_MISMATCH expected=$StartSha actual=$($Request.baseline_sha)"}
    if([string]$Request.test_script-notmatch'^scripts/.+\.ps1$'){throw "QUALIFICATION_TEST_SCRIPT_INVALID path=$($Request.test_script)"}

    $LastFiles=@(Get-CommitFiles $Commits[-1])
    if($LastFiles-notcontains'.wingless/qualification-request.json'){throw 'FINAL_IMPLEMENTATION_COMMIT_MISSING_QUALIFICATION_REQUEST'}

    $QualifiedHead=(& git rev-parse HEAD).Trim()
    $QualDir=Join-Path $RunRoot 'qualification'
    $SavedSha=$env:GITHUB_SHA
    $SavedRef=$env:GITHUB_REF_NAME
    try{
        $env:GITHUB_SHA=''
        $env:GITHUB_REF_NAME=''
        & (Join-Path $RepoPath 'scripts\Invoke-ResearchQualification.ps1') -RequestPath '.wingless/qualification-request.json' -OutputDirectory $QualDir -GuardWaitMinutes 120
        if($LASTEXITCODE-ne0){throw "WINGLESS_QUALIFICATION_FAILED exit=$LASTEXITCODE"}
    }finally{
        $env:GITHUB_SHA=$SavedSha
        $env:GITHUB_REF_NAME=$SavedRef
    }

    foreach($Dirty in @($Request.allowed_dirty_paths)){
        $Path=[string]$Dirty
        if([string]::IsNullOrWhiteSpace($Path)){continue}
        & git checkout -- $Path 2>$null
        if($LASTEXITCODE-ne0){
            $Full=Join-Path $RepoPath $Path
            if(Test-Path -LiteralPath $Full){Remove-Item -LiteralPath $Full -Force -Recurse}
        }
    }
    if(((& git status --porcelain)-join'').Trim()){throw "QUALIFICATION_LEFT_UNEXPECTED_DIRTY_WORKTREE status=$((& git status --porcelain)-join';')"}

    $EvidenceDir=Join-Path $RepoPath '.research-autonomy\evidence'
    New-Item -ItemType Directory -Force -Path $EvidenceDir|Out-Null
    $ExperimentSafe=([string]$Request.experiment)-replace'[^A-Za-z0-9._-]','_'
    $BundleDir=Join-Path $EvidenceDir ("$ExperimentSafe-$env:GITHUB_RUN_ID")
    New-Item -ItemType Directory -Force -Path $BundleDir|Out-Null
    foreach($Name in @('summary.json','probe.json','transcript.txt')){
        $Src=Join-Path $QualDir $Name
        if(Test-Path -LiteralPath $Src -PathType Leaf){Copy-Item -LiteralPath $Src -Destination (Join-Path $BundleDir $Name) -Force}
    }

    $SummaryFile=Join-Path $QualDir 'summary.json'
    if(-not(Test-Path -LiteralPath $SummaryFile -PathType Leaf)){throw 'QUALIFICATION_SUMMARY_MISSING'}
    $Summary=Get-Content -Raw -LiteralPath $SummaryFile|ConvertFrom-Json
    if([string]$Summary.classification-cne'qualified-scientific-result'){throw "QUALIFICATION_NOT_SCIENTIFIC_RESULT classification=$($Summary.classification)"}

    $State=[ordered]@{
        schema='research.repo-local-state.v1'
        program='Wingless'
        active_branch=$ActiveBranch
        north_star_path=$NorthStarPath
        north_star_sha256=$NorthStarSha256
        parent_sha=$StartSha
        qualified_head_sha=$QualifiedHead
        experiment=[string]$Request.experiment
        classification=[string]$Summary.classification
        last_completed_utc=[DateTime]::UtcNow.ToString('o')
        github_run_id=$env:GITHUB_RUN_ID
        authority='research-only'
        ckb_plane_role='governance-only'
    }
    New-Item -ItemType Directory -Force -Path (Split-Path -Parent $StatePath)|Out-Null
    $State|ConvertTo-Json -Depth 30|Set-Content -LiteralPath $StatePath -Encoding UTF8

    Invoke-Git add '.research-autonomy'
    Invoke-Git commit -m ("research: record qualified evidence for "+[string]$Request.experiment)
    $EvidenceHead=(& git rev-parse HEAD).Trim()
    Invoke-Git push origin ("HEAD:refs/heads/"+$ActiveBranch)

    Write-CycleResult -Advanced $true -Reason 'qualified-and-pushed' -Head $EvidenceHead
    Write-Host 'REPO_LOCAL_RESEARCH_CYCLE_PASS'
    Write-Host "experiment=$($Request.experiment)"
    Write-Host "parent_sha=$StartSha"
    Write-Host "qualified_head=$QualifiedHead"
    Write-Host "evidence_head=$EvidenceHead"
} catch {
    try{
        $Head=(& git -C $RepoPath rev-parse HEAD 2>$null|Select-Object -First 1)
        if($null-eq$Head){$Head=''}
        Write-CycleResult -Advanced $false -Reason $_.Exception.Message -Head ([string]$Head).Trim()
    }catch{}
    throw
} finally {
    Pop-Location
}
