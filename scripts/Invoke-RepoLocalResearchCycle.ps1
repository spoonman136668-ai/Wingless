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
# Keep the original GitHub credential requirement for the normal scheduler,
# but permit the *separately gated*, ephemeral, author-only free endpoint
# to enter its credential-denial path. It must never inherit a GitHub token.
if([string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)){
    $EarlyConfigPath=Join-Path (Split-Path -Parent $PSScriptRoot) '.research-autonomy\config.json'
    if(-not(Test-Path -LiteralPath $EarlyConfigPath -PathType Leaf)){
        throw 'GITHUB_TOKEN_MISSING'
    }
    $EarlyConfig=Get-Content -LiteralPath $EarlyConfigPath -Raw|ConvertFrom-Json
    if([string]$EarlyConfig.schema -cne 'research.autonomy.v2' -or
       [string]$EarlyConfig.program -cne 'Wingless' -or
       [string]$EarlyConfig.authority -cne 'research-only' -or
       [string]$EarlyConfig.research_executor -cne 'free-nemotron-opt-in-v1' -or
       [string]$env:RESEARCH_FREE_EXECUTOR_OWNER_ACK -cne '1' -or
       [string]$env:RUNNER_ENVIRONMENT -cne 'github-hosted'){
        throw 'GITHUB_TOKEN_MISSING'
    }
}

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

    $ExecutorConfig=Join-Path (Split-Path -Parent $PSScriptRoot) '.research-autonomy\config.json'
    if(-not(Test-Path -LiteralPath $ExecutorConfig -PathType Leaf)){throw 'RESEARCH_EXECUTOR_CONFIG_MISSING'}
    $ExecutorDocument=Get-Content -LiteralPath $ExecutorConfig -Raw | ConvertFrom-Json
    if([string]$ExecutorDocument.schema -cne 'research.autonomy.v2' -or
       [string]$ExecutorDocument.program -cne 'Wingless' -or
       [string]$ExecutorDocument.authority -cne 'research-only'){
        throw 'RESEARCH_EXECUTOR_SOURCE_POLICY_DRIFT'
    }
    $ExecutorMode='disabled'
    if($ExecutorDocument.PSObject.Properties.Name -contains 'research_executor'){
        $ExecutorMode=[string]$ExecutorDocument.research_executor
    }
    # External free trials must not receive private machine-local Mind-Palace
    # context even if a public source checkout is configured for the agent.
    if($ExecutorMode -ceq 'free-nemotron-opt-in-v1'){
        $MindContext=''
        Write-Host 'RESEARCH_EXTERNAL_PRIVATE_CONTEXT=EXCLUDED'
    }else{
        $MindContext=Get-OptionalMindContext
    }
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
4. WINDOWS FILE EDIT PREFLIGHT: Some Windows file-edit operations can fail with a sandbox-helper error. Prefer native PowerShell file operations (Set-Content, [IO.File]::WriteAllText with UTF-8 no BOM), git and gofmt; do not depend on a python command or retry a broken apply_patch helper. All paths must be relative to git rev-parse --show-toplevel. Before each commit verify git diff --cached --name-only is within the unchanged allowed research paths, and if required editing remains unavailable, stop with an explicit blocker rather than implying qualification.
4. Before implementation, create exactly ONE new preregistration under docs/experiments/. The preregistration must freeze the question/hypothesis, exact parent SHA, arms/inputs, seeds, budgets, metrics, thresholds, controls, validity criteria, classification rules, stop conditions, and a no-post-result-tuning rule.
5. Commit that preregistration ALONE as the first commit after $StartSha.
6. Then implement only the bounded experiment. Keep changes inside docs/experiments/, unitary/, cmd/, scripts/run-*.ps1 (experiment test harnesses only), and .wingless/qualification-request.json. Never modify controller, agent, workflow, model configuration, or autonomy scripts.
7. The final implementation commit must update .wingless/qualification-request.json using schema wingless.research-qualification-request.v1. Its branch MUST be $ActiveBranch and baseline_sha MUST be $StartSha. The test_script must be a repo-relative scripts/*.ps1 path.
8. Preserve deterministic controls, fixed resources, disjoint seeds where required, exact negative-result cardinality, and no post-result threshold/budget/capacity tuning.
9. Make all required commits locally with concise research commit messages. Leave the worktree clean.
10. If implementation becomes invalid after any output has been observed, do NOT repair or rerun it under the same preregistration. Restore every implementation/request change, leave exactly the preregistration commit with a clean worktree, make no scientific claim, and end your final response with a line containing exactly RESEARCH_TERMINAL=PREREG_BLOCKED followed by a concise reason.

If the evidence is insufficient, provenance is incomplete, the next question would widen authority, or safe continuation is ambiguous, make no changes and exit nonzero with a clear blocker.

OPTIONAL ADVISORY MIND-PALACE CONTEXT
$MindContext
"@
    [IO.File]::WriteAllText($PromptPath,$Prompt,(New-Object Text.UTF8Encoding($false)))

    # The original schedule and qualification wrapper remain authoritative.
    # An explicitly source-configured, second-gated free-only agent may author
    # a frozen research proposal. Without both gates this stays a read-only
    # watchdog; source promotion never silently enables remote inference.
    if($ExecutorMode -ceq 'disabled'){
        Write-Host 'RESEARCH_AGENT_POLICY=GPT_CHAT_GITHUB_SOURCE_ONLY_NO_CODEX_NO_EXTERNAL_APIS'
        Write-Host 'RESEARCH_EXECUTOR_UNATTENDED=NOT_CONFIGURED_NO_MODEL_INVOCATION'
        Write-CycleResult -Advanced $false -Reason 'gpt-authored-research-required-no-unattended-model' -Head $StartSha
        return
    }
    if($ExecutorMode -cne 'free-nemotron-opt-in-v1' -or
       [string]$env:RESEARCH_FREE_EXECUTOR_OWNER_ACK -cne '1'){
        throw 'RESEARCH_FREE_EXECUTOR_EXPLICIT_OWNER_GATE_REQUIRED'
    }
    # No paid provider/model fallback. DPAPI owner credentials and original
    # bound agent tool policy remain unchanged. Neither this script nor CI
    # claims a scientific first-use or automatically merges accepted refs.
    # Model-generated code and tool-initiated native tests must never run on
    # the privileged persistent Windows research machine with host secrets.
    # An ephemeral credential-free execution design must be qualified first.
    if([string]$env:RUNNER_ENVIRONMENT -cne 'github-hosted'){
        throw 'RESEARCH_FREE_EXECUTOR_REQUIRES_EPHEMERAL_RUNNER'
    }
    # This authoring phase must never inherit GitHub write or workstation credentials.
    # Public source metadata can be read without a GitHub bearer token.
    foreach($SecretName in @('GITHUB_TOKEN','GH_TOKEN','NVIDIA_API_KEY','OPENAI_API_KEY','AWS_ACCESS_KEY_ID')){
        if(-not [string]::IsNullOrWhiteSpace([string][Environment]::GetEnvironmentVariable($SecretName))){
            throw 'RESEARCH_FREE_AUTHOR_HAS_UNEXPECTED_PRIVILEGED_CREDENTIAL'
        }
    }
    if([string]::IsNullOrWhiteSpace([string]$env:RESEARCH_FREE_BUNDLE_OUTPUT)){
        throw 'RESEARCH_FREE_AUTHOR_BUNDLE_OUTPUT_REQUIRED'
    }
    if([string]::IsNullOrWhiteSpace([string]$env:RUNNER_TEMP) -or
       -not([IO.Path]::GetFullPath($env:RESEARCH_FREE_BUNDLE_OUTPUT).StartsWith(
           [IO.Path]::GetFullPath($env:RUNNER_TEMP).TrimEnd('\')+'\',
           [StringComparison]::OrdinalIgnoreCase))){
        throw 'RESEARCH_FREE_AUTHOR_BUNDLE_MUST_BE_IN_EPHEMERAL_TEMP'
    }
    if($Repository -cne 'spoonman136668-ai/Wingless'){
        throw 'RESEARCH_FREE_EXECUTOR_REPOSITORY_IDENTITY_DRIFT'
    }
    $Metadata=Invoke-RestMethod -Method Get -TimeoutSec 20 `
        -Uri ('https://api.github.com/repos/'+$Repository) `
        -Headers @{ Accept='application/vnd.github+json' }
    if([string]$Metadata.full_name -cne $Repository -or $Metadata.private -ne $false -or
       $Metadata.archived -ne $false){
        throw 'RESEARCH_FREE_EXECUTOR_PUBLIC_SOURCE_REQUIRED'
    }
    $OnlyFreeModel='nvidia/nemotron-3-ultra-550b-a55b:free'
    Write-Host 'RESEARCH_EXECUTOR_UNATTENDED=EXPLICIT_FREE_ONLY_OPT_IN'
    & (Join-Path $PSScriptRoot 'Invoke-NemotronRepoAgent.ps1') `
        -RepoPath $RepoPath -Program 'Wingless' -ActiveBranch $ActiveBranch `
        -StartSha $StartSha -PromptPath $PromptPath `
        -LastMessagePath $LastMessagePath -Provider 'openrouter' `
        -Model $OnlyFreeModel -MaxRounds 8 -NoProcessTools
    if(-not(Test-Path -LiteralPath $LastMessagePath -PathType Leaf)){
        throw 'RESEARCH_FREE_EXECUTOR_NO_TERMINAL_REPORT'
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
        $Allowed=$File-eq'.wingless/qualification-request.json' -or $File.StartsWith('docs/experiments/') -or $File.StartsWith('unitary/') -or $File.StartsWith('cmd/') -or ($File.StartsWith('scripts/run-') -and $File.EndsWith('.ps1') -and -not $File.Substring('scripts/'.Length).Contains('/'))
        if(-not$Allowed){throw "CHANGE_OUTSIDE_RESEARCH_POLICY path=$File"}
    }

    $FirstFiles=@(Get-CommitFiles $Commits[0])
    if($FirstFiles.Count-ne1 -or -not$FirstFiles[0].StartsWith('docs/experiments/')){throw "PREREG_NOT_FIRST_AND_ALONE files=$($FirstFiles -join ',')"}

    if($ExecutorMode -ceq 'free-nemotron-opt-in-v1'){
        # The remote author may create bounded local commits, never run its
        # scripts, qualify a result, push a branch or claim native acceptance
        # in the provider-credential job. The second independent job receives
        # a Git bundle and must verify base/candidate identity before testing.
        if($Commits.Count-eq1){
            $Summary=[IO.File]::ReadAllText($LastMessagePath)
            if($Summary -notmatch 'RESEARCH_TERMINAL=PREREG_BLOCKED'){
                throw 'FREE_AUTHOR_PREREG_ONLY_WITHOUT_BLOCKED_MARKER'
            }
        }else{
            $LastFiles=@(Get-CommitFiles $Commits[-1])
        if($LastFiles-notcontains'.wingless/qualification-request.json'){throw 'FREE_AUTHOR_REQUEST_MISSING'}
        }
        $BundleOut=[IO.Path]::GetFullPath($env:RESEARCH_FREE_BUNDLE_OUTPUT)
        if(Test-Path -LiteralPath $BundleOut){throw 'FREE_AUTHOR_BUNDLE_ALREADY_EXISTS'}
        & git bundle create $BundleOut HEAD ("^"+$StartSha)
        if($LASTEXITCODE-ne0){throw 'FREE_AUTHOR_BUNDLE_CREATE_FAILED'}
        & git bundle verify $BundleOut
        if($LASTEXITCODE-ne0){throw 'FREE_AUTHOR_BUNDLE_VERIFY_FAILED'}
        Write-CycleResult -Advanced $false -Reason 'credential-free-qualification-required' -Head $StartSha
        Write-Host 'RESEARCH_FREE_AUTHORING_ONLY_BUNDLE_READY'
        Write-Host 'RESEARCH_FREE_AUTHORING_NATIVE_EXECUTION=ZERO'
        Write-Host 'RESEARCH_FREE_AUTHORING_GITHUB_PUSH=ZERO'
        return
    }

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
    $TestScript=[string]$Request.test_script
    if(-not($TestScript.StartsWith('scripts/run-') -and $TestScript.EndsWith('.ps1') -and -not $TestScript.Substring('scripts/'.Length).Contains('/'))){throw "QUALIFICATION_TEST_SCRIPT_INVALID path=$($Request.test_script)"}

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
)
        if(-not$Allowed){throw "CHANGE_OUTSIDE_RESEARCH_POLICY path=$File"}
    }

    $FirstFiles=@(Get-CommitFiles $Commits[0])
    if($FirstFiles.Count-ne1 -or -not$FirstFiles[0].StartsWith('docs/experiments/')){throw "PREREG_NOT_FIRST_AND_ALONE files=$($FirstFiles -join ',')"}

    if($ExecutorMode -ceq 'free-nemotron-opt-in-v1'){
        # The remote author may create bounded local commits, never run its
        # scripts, qualify a result, push a branch or claim native acceptance
        # in the provider-credential job. The second independent job receives
        # a Git bundle and must verify base/candidate identity before testing.
        if($Commits.Count-eq1){
            $Summary=[IO.File]::ReadAllText($LastMessagePath)
            if($Summary -notmatch 'RESEARCH_TERMINAL=PREREG_BLOCKED'){
                throw 'FREE_AUTHOR_PREREG_ONLY_WITHOUT_BLOCKED_MARKER'
            }
        }else{
            $LastFiles=@(Get-CommitFiles $Commits[-1])
        if($LastFiles-notcontains'.wingless/qualification-request.json'){throw 'FREE_AUTHOR_REQUEST_MISSING'}
        }
        $BundleOut=[IO.Path]::GetFullPath($env:RESEARCH_FREE_BUNDLE_OUTPUT)
        if(Test-Path -LiteralPath $BundleOut){throw 'FREE_AUTHOR_BUNDLE_ALREADY_EXISTS'}
        & git bundle create $BundleOut HEAD ("^"+$StartSha)
        if($LASTEXITCODE-ne0){throw 'FREE_AUTHOR_BUNDLE_CREATE_FAILED'}
        & git bundle verify $BundleOut
        if($LASTEXITCODE-ne0){throw 'FREE_AUTHOR_BUNDLE_VERIFY_FAILED'}
        Write-CycleResult -Advanced $false -Reason 'credential-free-qualification-required' -Head $StartSha
        Write-Host 'RESEARCH_FREE_AUTHORING_ONLY_BUNDLE_READY'
        Write-Host 'RESEARCH_FREE_AUTHORING_NATIVE_EXECUTION=ZERO'
        Write-Host 'RESEARCH_FREE_AUTHORING_GITHUB_PUSH=ZERO'
        return
    }

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
    $TestScript=[string]$Request.test_script
    if(-not($TestScript.StartsWith('scripts/run-') -and $TestScript.EndsWith('.ps1') -and -not $TestScript.Substring('scripts/'.Length).Contains('/'))){throw "QUALIFICATION_TEST_SCRIPT_INVALID path=$($Request.test_script)"}

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
)
        if(-not$Allowed){throw "CHANGE_OUTSIDE_RESEARCH_POLICY path=$File"}
    }

    $FirstFiles=@(Get-CommitFiles $Commits[0])
    if($FirstFiles.Count-ne1 -or -not$FirstFiles[0].StartsWith('docs/experiments/')){throw "PREREG_NOT_FIRST_AND_ALONE files=$($FirstFiles -join ',')"}

    if($ExecutorMode -ceq 'free-nemotron-opt-in-v1'){
        # The remote author may create bounded local commits, never run its
        # scripts, qualify a result, push a branch or claim native acceptance
        # in the provider-credential job. The second independent job receives
        # a Git bundle and must verify base/candidate identity before testing.
        if($Commits.Count-eq1){
            $Summary=[IO.File]::ReadAllText($LastMessagePath)
            if($Summary -notmatch 'RESEARCH_TERMINAL=PREREG_BLOCKED'){
                throw 'FREE_AUTHOR_PREREG_ONLY_WITHOUT_BLOCKED_MARKER'
            }
        }else{
            $LastFiles=@(Get-CommitFiles $Commits[-1])
        if($LastFiles-notcontains'.wingless/qualification-request.json'){throw 'FREE_AUTHOR_REQUEST_MISSING'}
        }
        $BundleOut=[IO.Path]::GetFullPath($env:RESEARCH_FREE_BUNDLE_OUTPUT)
        if(Test-Path -LiteralPath $BundleOut){throw 'FREE_AUTHOR_BUNDLE_ALREADY_EXISTS'}
        & git bundle create $BundleOut HEAD ("^"+$StartSha)
        if($LASTEXITCODE-ne0){throw 'FREE_AUTHOR_BUNDLE_CREATE_FAILED'}
        & git bundle verify $BundleOut
        if($LASTEXITCODE-ne0){throw 'FREE_AUTHOR_BUNDLE_VERIFY_FAILED'}
        Write-CycleResult -Advanced $false -Reason 'credential-free-qualification-required' -Head $StartSha
        Write-Host 'RESEARCH_FREE_AUTHORING_ONLY_BUNDLE_READY'
        Write-Host 'RESEARCH_FREE_AUTHORING_NATIVE_EXECUTION=ZERO'
        Write-Host 'RESEARCH_FREE_AUTHORING_GITHUB_PUSH=ZERO'
        return
    }

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
    $TestScript=[string]$Request.test_script
    if(-not($TestScript.StartsWith('scripts/run-') -and $TestScript.EndsWith('.ps1') -and -not $TestScript.Substring('scripts/'.Length).Contains('/'))){throw "QUALIFICATION_TEST_SCRIPT_INVALID path=$($Request.test_script)"}

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
