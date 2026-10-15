param(
    [Parameter(Mandatory=$true)][ValidateSet('Wingless','Yggdrasil')][string]$Lane,
    [Parameter(Mandatory=$true)][string]$RepoPath,
    [Parameter(Mandatory=$true)][ValidatePattern('^[0-9a-f]{40}$')][string]$ParentSha,
    [Parameter(Mandatory=$true)][ValidatePattern('^[0-9a-f]{40}$')][string]$PreregistrationCommit,
    [Parameter(Mandatory=$true)][ValidatePattern('^[0-9a-f]{40}$')][string]$QualifiedHead,
    [Parameter(Mandatory=$true)][ValidatePattern('^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$')][string]$Experiment,
    [Parameter(Mandatory=$true)][ValidateSet('PRE_QUALIFICATION','PRE_SCIENTIFIC_RUN')][string]$Stage,
    [Parameter(Mandatory=$true)][ValidateRange(1,9223372036854775807)][long]$RunId,
    [Parameter(Mandatory=$true)][ValidateRange(1,100)][int]$RunAttempt,
    [Parameter(Mandatory=$true)][string]$OutputPath
)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

# Source-only owner observation, NOT a signer, original-claim writer or
# substitute for independently witnessed first-use and randomized selection.
if(-not[IO.Path]::IsPathRooted($RepoPath) -or
   -not(Test-Path -LiteralPath $RepoPath -PathType Container) -or
   -not[IO.Path]::IsPathRooted($OutputPath) -or
   $RunAttempt -ne 1){
    throw 'R31_SOURCE_OR_RUN_IDENTITY_UNQUALIFIED'
}
if(($Lane -ceq 'Wingless' -and $Stage -cne 'PRE_QUALIFICATION') -or
   ($Lane -ceq 'Yggdrasil' -and $Stage -cne 'PRE_SCIENTIFIC_RUN')){
    throw 'R31_LANE_EXECUTION_STAGE_MISMATCH'
}
$Head=(& git -C $RepoPath rev-parse HEAD)
if($LASTEXITCODE-ne0 -or ([string]$Head).Trim() -cne $QualifiedHead){
    throw 'R31_EXACT_RESEARCH_HEAD_CHANGED'
}
& git -C $RepoPath merge-base --is-ancestor $ParentSha $QualifiedHead
if($LASTEXITCODE-ne0){throw 'R31_FROZEN_PARENT_NOT_ANCESTOR'}
$Commits=@(& git -C $RepoPath rev-list --reverse "$ParentSha..$QualifiedHead")
if($LASTEXITCODE-ne0 -or $Commits.Count -lt 2 -or
   [string]$Commits[0] -cne $PreregistrationCommit){
    throw 'R31_PREREGISTRATION_NOT_FIRST'
}
$PrereqFiles=@(& git -C $RepoPath diff-tree --no-commit-id --name-only -r $PreregistrationCommit)
if($LASTEXITCODE-ne0 -or $PrereqFiles.Count-ne1){
    throw 'R31_PREREGISTRATION_SOURCE_NOT_ALONE'
}
$ExpectedPrefix=if($Lane -ceq 'Wingless'){'docs/experiments/'}else{'research/experiments/'}
if(-not([string]$PrereqFiles[0]).StartsWith($ExpectedPrefix,[StringComparison]::Ordinal)){
    throw 'R31_PREREGISTRATION_WRONG_LANE_PATH'
}
$Parent=Split-Path -Parent $OutputPath
if(-not(Test-Path -LiteralPath $Parent -PathType Container)){
    throw 'R31_EXPECTED_EXISTING_EVIDENCE_DIR'
}
$Event=[ordered]@{
    schema='forge.research-owner.r31-genuine-preexecution-observation.v1'
    lane=$Lane
    stage=$Stage
    experiment=$Experiment
    parent_sha=$ParentSha
    preregistration_commit=$PreregistrationCommit
    qualified_source_head=$QualifiedHead
    github_run_id=$RunId
    github_run_attempt=$RunAttempt
    observed_utc=[DateTime]::UtcNow.ToString('o')
    evidence_source='ACTUAL_REPOSITORY_LOCAL_OWNER_BEFORE_EXECUTION'
    observation_before_next_worker_invocation=$true
    predecision_signed_by_independent_owner=$false
    independent_external_first_use_verified=$false
    independent_randomized_selection_verified=$false
    logged_selection_probability_bp=0
    experimental_outcome_available=$false
    scientific_improvement_established=$false
    scientific_acceptance_authorized=$false
    optimizer_dispatch_authorized=$false
    original_first_use_claimed_by_this_probe=$false
}
$Utf8=New-Object Text.UTF8Encoding($false)
$Raw=$Utf8.GetBytes(($Event|ConvertTo-Json -Depth 10 -Compress)+"`n")
$Stream=[IO.FileStream]::new($OutputPath,[IO.FileMode]::CreateNew,[IO.FileAccess]::Write,[IO.FileShare]::None)
try{
    $Stream.Write($Raw,0,$Raw.Length)
    $Stream.Flush($true)
}finally{$Stream.Dispose()}
$Sha=[Security.Cryptography.SHA256]::Create()
try{
    $Fingerprint=([BitConverter]::ToString($Sha.ComputeHash($Raw))).Replace('-','').ToLowerInvariant()
}finally{$Sha.Dispose()}
Write-Host "R31_OWNER_PREEXECUTION_SOURCE_OBSERVATION=PASS lane=$Lane stage=$Stage sha256=$Fingerprint"
Write-Host 'R31_AUTHENTIC_RANDOMIZED_SCIENTIFIC_DECISION=NOT_ESTABLISHED'
