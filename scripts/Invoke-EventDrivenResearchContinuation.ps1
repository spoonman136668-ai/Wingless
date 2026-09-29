param(
    [Parameter(Mandatory=$true)][string]$RepoPath,
    [Parameter(Mandatory=$true)][string]$Repository,
    [Parameter(Mandatory=$true)][long]$RunId,
    [Parameter(Mandatory=$true)][string]$HeadBranch,
    [Parameter(Mandatory=$true)][string]$HeadSha,
    [Parameter(Mandatory=$true)][string]$Conclusion,
    [Parameter(Mandatory=$true)][ValidateSet('wingless','yggdrasil')][string]$Program,
    [string]$Lane = ''
)

$ErrorActionPreference = 'Stop'

function Invoke-Git {
    param([Parameter(ValueFromRemainingArguments=$true)][string[]]$Args)
    & git -C $RepoPath @Args
    if ($LASTEXITCODE -ne 0) { throw "GIT_FAILED exit=$LASTEXITCODE args=$($Args -join ' ')" }
}

function Get-CommitFiles([string]$Commit) {
    $R = @(& git -C $RepoPath diff-tree --no-commit-id --name-only -r $Commit)
    if ($LASTEXITCODE -ne 0) { throw "DIFF_TREE_FAILED commit=$Commit" }
    @($R | Where-Object { $_ -and $_.Trim() -ne '' })
}

if (-not (Test-Path -LiteralPath $RepoPath -PathType Container)) { throw "REPO_PATH_MISSING path=$RepoPath" }
if ([string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)) { throw 'GITHUB_TOKEN_MISSING' }

Push-Location $RepoPath
try {
    Invoke-Git fetch origin $HeadBranch --tags
    $RemoteSha = (& git rev-parse "origin/$HeadBranch").Trim()
    if ($LASTEXITCODE -ne 0) { throw 'REMOTE_HEAD_RESOLUTION_FAILED' }
    if ($RemoteSha -cne $HeadSha) {
        Write-Host "RESEARCH_CONTINUATION_STALE_EVENT_SKIP expected=$HeadSha remote=$RemoteSha"
        return
    }

    $ReceiptTag = "research-continuation/$RunId"
    & git ls-remote --exit-code --tags origin "refs/tags/$ReceiptTag" *> $null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "RESEARCH_CONTINUATION_ALREADY_PROCESSED run_id=$RunId"
        return
    }

    Invoke-Git switch -C $HeadBranch "origin/$HeadBranch"
    $StartSha = (& git rev-parse HEAD).Trim()
    if ($StartSha -cne $HeadSha) { throw "START_HEAD_MISMATCH expected=$HeadSha actual=$StartSha" }
    if ((& git status --porcelain | Out-String).Trim()) { throw 'START_WORKTREE_DIRTY' }

    $EvidenceRoot = Join-Path $env:RUNNER_TEMP ("research-continuation-" + $RunId)
    if (Test-Path -LiteralPath $EvidenceRoot) { Remove-Item -LiteralPath $EvidenceRoot -Recurse -Force }
    New-Item -ItemType Directory -Force -Path $EvidenceRoot | Out-Null
    $LogsZip = Join-Path $EvidenceRoot 'logs.zip'
    $LogsDir = Join-Path $EvidenceRoot 'logs'
    $Headers = @{ Authorization = "Bearer $env:GITHUB_TOKEN"; Accept = 'application/vnd.github+json'; 'X-GitHub-Api-Version' = '2022-11-28'; 'User-Agent' = 'research-continuation-controller' }
    Invoke-WebRequest -UseBasicParsing -Headers $Headers -Uri "https://api.github.com/repos/$Repository/actions/runs/$RunId/logs" -OutFile $LogsZip
    New-Item -ItemType Directory -Force -Path $LogsDir | Out-Null
    Expand-Archive -LiteralPath $LogsZip -DestinationPath $LogsDir -Force

    $Codex = Get-Command codex -ErrorAction SilentlyContinue
    if (-not $Codex) { $Codex = Get-Command codex.cmd -ErrorAction SilentlyContinue }
    if (-not $Codex) { throw 'CODEX_NOT_AVAILABLE_ON_RESEARCH_CONTROLLER_RUNNER' }

    if ($Program -eq 'wingless') {
        $DispatchPath = '.wingless/qualification-request.json'
        $PreregPrefix = 'docs/experiments/'
        $ProgramRules = @"
PROGRAM-SPECIFIC RULES:
- This is Wingless research only.
- For a completed scientific result, create a NEW successor branch beginning with research/wingless-up.
- The successor must inherit the exact completed HEAD as its baseline.
- First commit on the successor branch must create exactly one new docs/experiments/*.ice preregistration and may change nothing else.
- Only after that first preregistration commit may you implement source/tests.
- The FINAL commit must update .wingless/qualification-request.json so the existing Windows qualification workflow dispatches.
- Preserve zero live activation, no CKB/ckb-plane/KTRADE/broker/credential/production authority, and no post-result tuning.
"@
    }
    else {
        if ($Lane -notin @('a','b','c')) { throw "YGG_LANE_INVALID lane=$Lane" }
        $DispatchPath = "research/control/ygg-$Lane-run.json"
        $PreregPrefix = "research/experiments/ygg-$Lane/"
        $ProgramRules = @"
PROGRAM-SPECIFIC RULES:
- This is Yggdrasil synthetic research only; no wetware or living tissue.
- Stay on the existing lane branch $HeadBranch.
- First successor commit must create exactly one new research/experiments/ygg-$Lane/*.ice preregistration and may change nothing else.
- Only after that preregistration commit may you implement source/workflow changes.
- The FINAL commit must update $DispatchPath so the lane workflow dispatches.
- Preserve fixed resources, frozen scientific controls, disjoint seeds when required, deterministic duplicate checks, and no post-result tuning.
- Never touch another Yggdrasil lane except shared code already explicitly part of the frozen parent substrate.
"@
    }

    $Prompt = @"
You are the bounded event-driven research continuation worker.

Repository: $Repository
Program: $Program
Lane: $Lane
Completed workflow run id: $RunId
Completed branch: $HeadBranch
Completed exact SHA: $HeadSha
Workflow conclusion: $Conclusion
Local extracted workflow logs: $LogsDir

Your job is to continue this one research lane without operator prompting while preserving scientific validity.

Read the repository charter/preregistration lineage, the completed run logs, and the latest exact result. Do not guess missing evidence.

If the completed run contains a valid sealed scientific result:
1. Treat positive and negative scientific outcomes as valid evidence.
2. Derive exactly ONE bounded next question that materially reduces uncertainty toward the program North Star.
3. BEFORE implementation, write a new preregistration freezing hypothesis/question, exact parent identity, arms/inputs, seeds, budgets, thresholds, controls, validity criteria, classifications, stop conditions, and no-post-result-tuning rule.
4. Commit that preregistration ALONE as the first commit after $HeadSha.
5. Then implement the experiment and deterministic qualification.
6. Dispatch it only in the final commit.

If the completed run failed only because of infrastructure or a pre-result implementation/harness defect:
- Do NOT alter the scientific preregistration, frozen thresholds, arms, seeds, budgets, or interpretation.
- Make only the minimum bounded repair and redispatch the exact same experiment.
- Scientific output already observed may not be used to retune the experiment.

If evidence is ambiguous, provenance is incomplete, or safe continuation would widen authority, make NO changes and exit nonzero with a clear blocker.

Never modify accepted refs, production systems, credentials, broker/live interfaces, CKB, ckb-plane, KTRADE, or the other research program.
Never weaken or delete a failing scientific assertion just to pass.
Never increase capacity/budget because of the observed result.
Do not push to GitHub yourself. The deterministic wrapper validates and pushes only after your work is complete.
Make all required git commits locally with concise research commit messages.

$ProgramRules
"@

    $PromptPath = Join-Path $EvidenceRoot 'codex-prompt.txt'
    $LastMessagePath = Join-Path $EvidenceRoot 'codex-last-message.txt'
    $Utf8NoBom = New-Object System.Text.UTF8Encoding($false)
    [IO.File]::WriteAllText($PromptPath, $Prompt, $Utf8NoBom)

    $Proc = Start-Process -FilePath $Codex.Source -ArgumentList @('exec','--full-auto','--output-last-message',$LastMessagePath,'-') -WorkingDirectory $RepoPath -RedirectStandardInput $PromptPath -NoNewWindow -Wait -PassThru
    if ($Proc.ExitCode -ne 0) { throw "CODEX_CONTINUATION_FAILED exit=$($Proc.ExitCode)" }
    if ((& git status --porcelain | Out-String).Trim()) { throw 'CODEX_LEFT_UNCOMMITTED_CHANGES' }

    $CurrentBranch = (& git branch --show-current).Trim()
    if ([string]::IsNullOrWhiteSpace($CurrentBranch)) { throw 'CONTINUATION_BRANCH_MISSING' }

    if ($Program -eq 'wingless') {
        if ($Conclusion -eq 'success' -and $CurrentBranch -ceq $HeadBranch) { throw 'WINGLESS_SUCCESSOR_BRANCH_NOT_CREATED' }
        if ($CurrentBranch -notlike 'research/wingless-up*') { throw "WINGLESS_BRANCH_POLICY_VIOLATION branch=$CurrentBranch" }
    }
    else {
        if ($CurrentBranch -cne $HeadBranch) { throw "YGG_BRANCH_CHANGED expected=$HeadBranch actual=$CurrentBranch" }
    }

    $Commits = @(& git rev-list --reverse "$StartSha..HEAD") | Where-Object { $_ -and $_.Trim() -ne '' }
    if ($LASTEXITCODE -ne 0) { throw 'CONTINUATION_REV_LIST_FAILED' }
    if ($Commits.Count -lt 1) { throw 'CONTINUATION_NO_COMMITS' }

    $AllFiles = New-Object System.Collections.Generic.List[string]
    foreach ($Commit in $Commits) { foreach ($File in (Get-CommitFiles $Commit)) { $AllFiles.Add($File) } }

    foreach ($File in $AllFiles) {
        if ($File -match '(?i)(^|/)(ckb|ckb-plane|ktrade)(/|$)' -or $File -match '(?i)(coinbase|broker|credential|accepted-ref)') { throw "FORBIDDEN_PATH_CHANGE path=$File" }
        if ($Program -eq 'wingless') {
            $Allowed = $File -eq '.wingless/qualification-request.json' -or $File.StartsWith('docs/experiments/') -or $File.StartsWith('unitary/') -or $File.StartsWith('cmd/') -or $File.StartsWith('scripts/') -or $File -eq '.github/workflows/research-qualify-windows.yml'
        }
        else {
            $Allowed = $File.StartsWith("research/experiments/ygg-$Lane/") -or $File -eq ".github/workflows/ygg-$Lane-lane.yml" -or $File -eq $DispatchPath
        }
        if (-not $Allowed) { throw "CHANGE_OUTSIDE_LANE_POLICY path=$File" }
    }

    $FirstFiles = @(Get-CommitFiles $Commits[0])
    $FirstIsPrereg = $FirstFiles.Count -eq 1 -and $FirstFiles[0].StartsWith($PreregPrefix) -and $FirstFiles[0].EndsWith('.ice')
    if ($Conclusion -eq 'success' -and -not $FirstIsPrereg) { throw "PREREG_NOT_FIRST_COMMIT files=$($FirstFiles -join ',')" }

    $LastFiles = @(Get-CommitFiles $Commits[-1])
    if ($LastFiles -notcontains $DispatchPath) { throw "FINAL_COMMIT_DID_NOT_DISPATCH expected=$DispatchPath files=$($LastFiles -join ',')" }

    Invoke-Git push origin "HEAD:refs/heads/$CurrentBranch"
    Invoke-Git tag -f $ReceiptTag HEAD
    Invoke-Git push origin "refs/tags/$ReceiptTag"

    Write-Host 'RESEARCH_CONTINUATION_DISPATCHED'
    Write-Host "run_id=$RunId"
    Write-Host "source_sha=$StartSha"
    Write-Host "successor_branch=$CurrentBranch"
    Write-Host "successor_head=$((& git rev-parse HEAD).Trim())"
    Write-Host "commits=$($Commits.Count)"
}
finally {
    Pop-Location
}
