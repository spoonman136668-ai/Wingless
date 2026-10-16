# Archived compatibility entry point. Automatic research continuation is not
# a Codex client; no subscription-backed CLI, alternative API, or unattended
# model invocation is authorized here. Preserve the original CLI parameters
# for callers and return without mutating evidence, tags, branches or science.
param(
    [Parameter(Mandatory=$true)][string]$RepoPath,
    [Parameter(Mandatory=$true)][string]$Repository,
    [Parameter(Mandatory=$true)][long]$RunId,
    [Parameter(Mandatory=$true)][string]$HeadBranch,
    [Parameter(Mandatory=$true)][string]$HeadSha,
    [Parameter(Mandatory=$true)][string]$Conclusion
)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
if($Repository -cne 'spoonman136668-ai/Wingless'){throw 'RESEARCH_REPOSITORY_MISMATCH'}
if($HeadSha -cnotmatch '^[0-9a-f]{40}$'){throw 'RESEARCH_SOURCE_SHA_UNVERIFIED'}
Write-Host 'RESEARCH_AGENT_POLICY=GPT_CHAT_GITHUB_SOURCE_ONLY_NO_CODEX_NO_EXTERNAL_APIS'
Write-Host "RESEARCH_CONTINUATION_NOT_DISPATCHED=NO_UNATTENDED_MODEL run_id=$RunId"
return
