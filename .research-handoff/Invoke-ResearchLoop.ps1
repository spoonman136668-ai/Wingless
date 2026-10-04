[CmdletBinding()]
param(
  [string]$ConfigPath = "C:\ProgramData\CKBR\research\controller.json",
  [switch]$DryRun,
  [switch]$SelfTest,
  [string]$FixtureRoot = ""
)

Set-StrictMode -Version 3.0
$ErrorActionPreference = "Stop"
Add-Type -AssemblyName System.Security -ErrorAction Stop
$script:CodexUsed = $false
$script:ProductionAuthorityChanged = $false
$script:AcceptedRefMutated = $false
$script:RuntimeLaunched = $false

function Get-UtcNow { (Get-Date).ToUniversalTime().ToString("o") }

function Get-Sha256Bytes([byte[]]$Bytes) {
  $sha = [System.Security.Cryptography.SHA256]::Create()
  try { return ([BitConverter]::ToString($sha.ComputeHash($Bytes))).Replace("-","").ToLowerInvariant() }
  finally { $sha.Dispose() }
}

function Get-Sha256Text([string]$Text) {
  Get-Sha256Bytes ([Text.Encoding]::UTF8.GetBytes($Text))
}

function Write-AtomicUtf8([string]$Path, [string]$Text) {
  $dir = Split-Path -Parent $Path
  if ($dir) { [IO.Directory]::CreateDirectory($dir) | Out-Null }
  $tmp = "$Path.tmp.$PID"
  [IO.File]::WriteAllText($tmp, $Text, (New-Object Text.UTF8Encoding($false)))
  Move-Item -LiteralPath $tmp -Destination $Path -Force
}

function Read-JsonFile([string]$Path) {
  if (!(Test-Path -LiteralPath $Path)) { throw "MISSING_JSON:$Path" }
  return Get-Content -LiteralPath $Path -Raw | ConvertFrom-Json
}

function ConvertTo-CanonicalJson($Value) {
  return ($Value | ConvertTo-Json -Depth 100 -Compress)
}

function Invoke-ProcessStrict([string]$FilePath, [string[]]$Arguments, [int[]]$AllowedExitCodes = @(0), [string]$WorkingDirectory = "") {
  if ($FilePath -match "(?i)(^|[\\/])codex(\.exe)?$") {
    $script:CodexUsed = $true
    throw "CODEX_USAGE_FORBIDDEN"
  }
  $oldLocation = Get-Location
  $oldEap = $ErrorActionPreference
  try {
    if ($WorkingDirectory) { Set-Location -LiteralPath $WorkingDirectory }
    $ErrorActionPreference = "Continue"
    $captured = @(& $FilePath @Arguments 2>&1)
    $exitCode = $LASTEXITCODE
  } finally {
    $ErrorActionPreference = $oldEap
    Set-Location -LiteralPath $oldLocation
  }
  $stdout = ($captured | ForEach-Object { [string]$_ }) -join [Environment]::NewLine
  if ($AllowedExitCodes -notcontains $exitCode) {
    throw "PROCESS_FAILED:${FilePath}:${exitCode}:$stdout"
  }
  [pscustomobject]@{ exit_code=$exitCode; stdout=$stdout.Trim(); stderr="" }
}

function Invoke-Git([string]$Repo, [string[]]$GitArgs, [int[]]$AllowedExitCodes = @(0)) {
  Invoke-ProcessStrict -FilePath "git.exe" -Arguments (@("-C",$Repo) + $GitArgs) -AllowedExitCodes $AllowedExitCodes
}

function Invoke-Gh([string[]]$GhArgs, [int[]]$AllowedExitCodes = @(0)) {
  Invoke-ProcessStrict -FilePath "gh.exe" -Arguments $GhArgs -AllowedExitCodes $AllowedExitCodes
}

function Get-ResearchLoopMutex([string]$Name) {
  $created = $false
  $m = New-Object Threading.Mutex($false, $Name, [ref]$created)
  try {
    if (!$m.WaitOne(0)) {
      $m.Dispose()
      return $null
    }
  } catch [Threading.AbandonedMutexException] {
  }
  return $m
}

function Apply-QualifiedControllerIdentity($Config, [string]$ScriptPath=$PSCommandPath) {
  $dir=Split-Path -Parent $ScriptPath
  $identityPath=Join-Path $dir "controller-qualified-sha.txt"
  if(!(Test-Path -LiteralPath $identityPath)){return $null}
  $sha=([IO.File]::ReadAllText($identityPath,[Text.Encoding]::UTF8)).Trim().ToLowerInvariant()
  if($sha -notmatch "^[0-9a-f]{40}$"){throw "QUALIFIED_CONTROLLER_IDENTITY_INVALID:$sha"}
  $Config.controller_sha=$sha
  return $sha
}

function Assert-GitHubAccess($Config) {
  if ($FixtureRoot) { return }
  Invoke-Gh @("auth","status") | Out-Null
  foreach ($repo in @($Config.repositories.ckb_plane,$Config.repositories.wingless,$Config.repositories.yggdrasil)) {
    $r = Invoke-Gh @("repo","view",$repo,"--json","nameWithOwner")
    $j = $r.stdout | ConvertFrom-Json
    if ($j.nameWithOwner -ne $repo) { throw "GITHUB_REPOSITORY_ACCESS_MISMATCH:$repo" }
    $p = Invoke-Gh @("api","repos/$repo","--jq",".permissions.push")
    if ($p.stdout.Trim().ToLowerInvariant() -ne "true") { throw "GITHUB_REPOSITORY_WRITE_ACCESS_REQUIRED:$repo" }
  }
}

function Sync-ContractSource($Config) {
  if($FixtureRoot){return $null}
  if(!$Config.PSObject.Properties["contract_source_sync"] -or !$Config.contract_source_sync){return $null}
  $repoPath=[string]$Config.contract_source_repo
  $ref=if($Config.PSObject.Properties["contract_source_ref"]){[string]$Config.contract_source_ref}else{"main"}
  if(!(Test-Path -LiteralPath (Join-Path $repoPath ".git"))){throw "CONTRACT_SOURCE_REPO_MISSING:$repoPath"}
  $dirty=(Invoke-Git $repoPath @("status","--porcelain")).stdout
  if($dirty){throw "CONTRACT_SOURCE_REPO_DIRTY:$repoPath"}
  Invoke-Git $repoPath @("fetch","origin",$ref,"--prune")|Out-Null
  $remote=(Invoke-Git $repoPath @("rev-parse","origin/$ref")).stdout
  if($Config.PSObject.Properties["authority_floor_sha"] -and $Config.authority_floor_sha){
    Invoke-Git $repoPath @("merge-base","--is-ancestor",[string]$Config.authority_floor_sha,$remote)|Out-Null
  }
  Invoke-Git $repoPath @("checkout","--detach",$remote)|Out-Null
  $contracts=Join-Path $repoPath "snapshot\source\research-controller\contracts"
  $templates=Join-Path $repoPath "snapshot\source\research-controller\templates"
  if(!(Test-Path -LiteralPath $contracts) -or !(Test-Path -LiteralPath $templates)){throw "CONTRACT_SOURCE_PACKAGE_MISSING:$remote"}
  $Config.contract_root=$contracts
  $Config.template_root=$templates
  return $remote
}


function Get-RemoteBranch([string]$Repo, [string]$Branch) {
  if ($FixtureRoot) {
    $p = Join-Path $FixtureRoot ("branch-" + (($Repo + "-" + $Branch) -replace '[^A-Za-z0-9_.-]','_') + ".json")
    if (!(Test-Path -LiteralPath $p)) { return $null }
    return Read-JsonFile $p
  }
  $encoded = [Uri]::EscapeDataString($Branch)
  $r = Invoke-Gh @("api","repos/$Repo/branches/$encoded") -AllowedExitCodes @(0,1)
  if ($r.exit_code -ne 0 -or [string]::IsNullOrWhiteSpace($r.stdout)) { return $null }
  return ($r.stdout | ConvertFrom-Json)
}

function Get-WorkflowRuns([string]$Repo, [string]$Branch, [int]$Limit = 20) {
  if ($FixtureRoot) {
    $p = Join-Path $FixtureRoot ("runs-" + (($Repo + "-" + $Branch) -replace '[^A-Za-z0-9_.-]','_') + ".json")
    if (!(Test-Path -LiteralPath $p)) { return @() }
    $j = Read-JsonFile $p
    return @($j.workflow_runs)
  }
  $encoded = [Uri]::EscapeDataString($Branch)
  $r = Invoke-Gh @("api","repos/$Repo/actions/runs?branch=$encoded&per_page=$Limit")
  $j = $r.stdout | ConvertFrom-Json
  return @($j.workflow_runs)
}

function Get-RunArtifacts([string]$Repo, [long]$RunId, [string]$Destination) {
  [IO.Directory]::CreateDirectory($Destination) | Out-Null
  if ($FixtureRoot) {
    $src = Join-Path $FixtureRoot "artifacts\$RunId"
    if (!(Test-Path -LiteralPath $src)) { throw "FIXTURE_ARTIFACTS_MISSING:$RunId" }
    Copy-Item -Path (Join-Path $src "*") -Destination $Destination -Recurse -Force
    return
  }
  Invoke-Gh @("run","download","$RunId","--repo",$Repo,"--dir",$Destination) | Out-Null
}

function Classify-WorkflowFailure($Config, [string]$Lane, $Run) {
  if ($Run.PSObject.Properties["failure_class"] -and $Run.failure_class) {
    return [pscustomobject]@{failure_class=[string]$Run.failure_class;failing_step="fixture";error_signature=(Get-Sha256Text ([string]$Run.failure_class))}
  }
  if ($FixtureRoot) {
    return [pscustomobject]@{failure_class="UNCLASSIFIED_WORKFLOW_FAILURE";failing_step="workflow";error_signature=(Get-Sha256Text ([string]$Run.conclusion))}
  }
  $repo = if ($Lane -eq "Wingless") { $Config.repositories.wingless } else { $Config.repositories.yggdrasil }
  $view = Invoke-Gh @("run","view","$($Run.id)","--repo",$repo,"--log-failed") -AllowedExitCodes @(0,1)
  $s = [string]$view.stdout
  $normalized = ($s -replace '\s+',' ').Trim()
  $sigSource = if ($normalized.Length -gt 4000) { $normalized.Substring(0,4000) } else { $normalized }
  $sig = Get-Sha256Text $sigSource
  if ($s -match '(?i)(runner.+(shutdown|lost|offline)|lost communication|runner process.+stopped|operation was canceled)') {
    return [pscustomobject]@{failure_class="RUNNER_INTERRUPTION";failing_step="runner";error_signature=$sig}
  }
  if ($s -match '(?i)(HTTP (500|502|503|504)|temporarily unavailable|service unavailable|secondary rate limit|connection reset)') {
    return [pscustomobject]@{failure_class="TRANSIENT_GITHUB_API";failing_step="github";error_signature=$sig}
  }
  if ($s -match '(?i)(Validate frozen package|Verify preregistration and bounded implementation delta|No files were found with the provided path|artifact.+not found|workflow packaging)') {
    return [pscustomobject]@{failure_class="WORKFLOW_PACKAGING_ERROR";failing_step="package-validation";error_signature=$sig}
  }
  return [pscustomobject]@{failure_class="UNCLASSIFIED_WORKFLOW_FAILURE";failing_step="workflow";error_signature=$sig}
}

function Find-ClassificationArtifact([string]$Root, [string]$ExperimentId) {
  $files = @(Get-ChildItem -LiteralPath $Root -Recurse -File -Filter "classification.json" -ErrorAction SilentlyContinue)
  if ($files.Count -eq 0) { throw "CLASSIFICATION_ARTIFACT_MISSING:$ExperimentId" }
  $matches = New-Object Collections.Generic.List[object]
  foreach ($f in $files) {
    try {
      $j = Read-JsonFile $f.FullName
      if (($j.experiment -eq $ExperimentId) -or ($j.experiment_id -eq $ExperimentId)) { [void]$matches.Add($j) }
    } catch {}
  }
  if ($matches.Count -ne 1) { throw "CLASSIFICATION_ARTIFACT_NOT_UNIQUE:${ExperimentId}:$($matches.Count)" }
  return $matches[0]
}

function Get-RepairFingerprint([string]$Lane,[string]$Experiment,[string]$Stage,[string]$HeadSha,[string]$FailureClass,[string]$FailingStep,[string]$ErrorSignature) {
  $raw = @($Lane,$Experiment,$Stage,$HeadSha,$FailureClass,$FailingStep,$ErrorSignature) -join [Environment]::NewLine
  return Get-Sha256Text $raw
}

function Test-NeverAutoRepair([string]$FailureClass) {
  return @(
    "SCIENTIFIC_NEGATIVE","SCIENTIFIC_MIXED","PROVENANCE_AMBIGUITY",
    "DETERMINISM_DIVERGENCE","SAFETY_TRIPWIRE","AUTHORITY_EXPANSION",
    "CAPACITY_INCREASE","FROZEN_PREMISE_CHANGE","ARCHITECTURAL_BOUNDARY",
    "ACCEPTED_REF_MUTATION_REQUIRED","PRODUCTION_ACTIVATION_REQUIRED"
  ) -contains $FailureClass
}

function Resolve-Successor($Contract, [string]$Classification) {
  $c = $Classification.ToLowerInvariant()
  if ($c -eq "supported") { return $Contract.successors.supported }
  if ($c -eq "mixed") { return $Contract.successors.mixed }
  if ($c -eq "negative") { return $Contract.successors.negative }
  if ($c -eq "invalid") { return $null }
  throw "UNEXPECTED_SCIENTIFIC_CLASSIFICATION:$Classification"
}

function Render-Template([string]$Template, $Variables) {
  $out = $Template
  $names = @($Variables.PSObject.Properties.Name | Sort-Object)
  foreach ($name in $names) {
    $token = "{{" + $name + "}}"
    $out = $out.Replace($token, [string]$Variables.$name)
  }
  if ($out -match '\{\{[A-Z0-9_]+\}\}') { throw "UNRESOLVED_TEMPLATE_TOKEN:$($Matches[0])" }
  return $out
}

function Assert-Contract($Contract) {
  if ($Contract.schema -ne "ckb-plane.research-contract.v1") { throw "CONTRACT_SCHEMA_INVALID" }
  foreach ($name in @("experiment_id","lane","parent_experiment","parent_sha","hypothesis","experiment_family","manifest_sha","classification_rules","successors","required_authority","qualification_profile","implementation_template")) {
    if ($null -eq $Contract.$name) { throw "CONTRACT_FIELD_MISSING:$name" }
  }
  if ($Contract.lane -notin @("Wingless","Yggdrasil")) { throw "CONTRACT_LANE_INVALID" }
  if ([string]$Contract.manifest_sha -notmatch '^[0-9a-f]{64}$') { throw "CONTRACT_MANIFEST_SHA_INVALID" }
  if ($Contract.codex_usage -ne "DISABLED") { throw "CONTRACT_CODEX_NOT_DISABLED" }
  if ($Contract.production_authority -ne $false) { throw "CONTRACT_PRODUCTION_AUTHORITY_FORBIDDEN" }
  if ($Contract.accepted_ref_mutation -ne $false) { throw "CONTRACT_ACCEPTED_REF_MUTATION_FORBIDDEN" }
  if ($Contract.runtime_launch -ne $false) { throw "CONTRACT_RUNTIME_LAUNCH_FORBIDDEN" }
  if ($Contract.capacity_change_authorized -ne $false) { throw "CONTRACT_CAPACITY_CHANGE_FORBIDDEN" }
  if ($Contract.authority_expansion -ne $false) { throw "CONTRACT_AUTHORITY_EXPANSION_FORBIDDEN" }
  if ($Contract.broker_access -ne $false) { throw "CONTRACT_BROKER_ACCESS_FORBIDDEN" }
  if ($Contract.ktrade_access -ne $false) { throw "CONTRACT_KTRADE_ACCESS_FORBIDDEN" }
  if ($Contract.implementation_template.ready) {
    if ([string]::IsNullOrWhiteSpace([string]$Contract.branch)) { throw "CONTRACT_BRANCH_REQUIRED" }
    if ([string]::IsNullOrWhiteSpace([string]$Contract.preregistration_sha)) { throw "CONTRACT_PREREGISTRATION_SHA_REQUIRED" }
    if ([string]::IsNullOrWhiteSpace([string]$Contract.preregistered_at_utc)) { throw "CONTRACT_PREREGISTERED_TIME_REQUIRED" }
  }
}

function New-ScientistRequest([string]$Lane, $Current, [string]$Reason, $Classification, $Contract, $EvidenceRefs) {
  $req = [ordered]@{
    schema="ckb-plane.research-scientist-request.v1"
    lane=$Lane
    reason=$Reason
    experiment_id=$Current.experiment_id
    head_sha=$Current.head_sha
    classification=$Classification
    key_metrics=$Current.key_metrics
    frozen_constraints=$Contract.frozen_constraints
    eliminated_explanations=$Contract.eliminated_explanations
    remaining_uncertainties=$Contract.remaining_uncertainties
    evidence_refs=@($EvidenceRefs)
  }
  return [pscustomobject]$req
}

function Save-ScientistRequest([string]$StateRoot, $Request) {
  $json = ConvertTo-CanonicalJson $Request
  $sha = Get-Sha256Text $json
  $dir = Join-Path $StateRoot "scientist-requests"
  [IO.Directory]::CreateDirectory($dir) | Out-Null
  $path = Join-Path $dir "$sha.json"
  if (!(Test-Path -LiteralPath $path)) { Write-AtomicUtf8 $path ($json + [Environment]::NewLine) }
  return [pscustomobject]@{ sha256=$sha; path=$path }
}

function Publish-ScientistRequest($Config, $Request, $Saved) {
  if($DryRun -or $FixtureRoot){return}
  if(!$Config.PSObject.Properties["publish_scientist_requests"] -or !$Config.publish_scientist_requests){return}
  $repo=[string]$Config.scientist_request_issue_repo
  if([string]::IsNullOrWhiteSpace($repo)){throw "SCIENTIST_REQUEST_ISSUE_REPO_REQUIRED"}
  $prefix=if($Config.PSObject.Properties["scientist_request_issue_prefix"]){[string]$Config.scientist_request_issue_prefix}else{"Research scientist request"}
  $title=$prefix+" "+$Saved.sha256
  $listed=Invoke-Gh @("issue","list","--repo",$repo,"--state","all","--search",($Saved.sha256+" in:title"),"--limit","20","--json","number,title,state")
  $rows=@($listed.stdout|ConvertFrom-Json)
  $exact=$rows|Where-Object{$p=$_.PSObject.Properties["title"]; $p -and [string]$p.Value -eq $title}|Select-Object -First 1
  if($exact){return}
  $tmp=Join-Path $Config.state_root ("scientist-issue-"+$Saved.sha256+".json")
  try{
    Write-AtomicUtf8 $tmp ((ConvertTo-CanonicalJson $Request)+[Environment]::NewLine)
    Invoke-Gh @("issue","create","--repo",$repo,"--title",$title,"--body-file",$tmp)|Out-Null
  }finally{
    Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue
  }
}

function Load-ControllerState([string]$StateRoot) {
  $p = Join-Path $StateRoot "state.json"
  if (!(Test-Path -LiteralPath $p)) {
    return [pscustomobject]@{schema="ckb-plane.research-loop-state.v1";consumed=[pscustomobject]@{};repairs=[pscustomobject]@{};lanes=[pscustomobject]@{}}
  }
  return Read-JsonFile $p
}

function Save-ControllerState([string]$StateRoot, $State) {
  [IO.Directory]::CreateDirectory($StateRoot) | Out-Null
  Write-AtomicUtf8 (Join-Path $StateRoot "state.json") ((ConvertTo-CanonicalJson $State) + [Environment]::NewLine)
}

function Save-ControllerReceipt([string]$StateRoot, $Receipt) {
  $json = ConvertTo-CanonicalJson $Receipt
  $sha = Get-Sha256Text $json
  $hist = Join-Path $StateRoot "receipts"
  [IO.Directory]::CreateDirectory($hist) | Out-Null
  $path = Join-Path $hist "$sha.json"
  if (!(Test-Path -LiteralPath $path)) { Write-AtomicUtf8 $path ($json + [Environment]::NewLine) }
  Write-AtomicUtf8 (Join-Path $StateRoot "latest-receipt.json") ($json + [Environment]::NewLine)
  return $sha
}

function Get-ContractPath($Config, [string]$ContractId) { Join-Path $Config.contract_root ($ContractId + ".json") }

function Load-Contract($Config, [string]$ContractId) {
  $c = Read-JsonFile (Get-ContractPath $Config $ContractId)
  Assert-Contract $c
  $c | Add-Member -NotePropertyName "_contract_id" -NotePropertyValue $ContractId -Force
  return $c
}

function Get-InitialLaneState($LaneConfig) {
  [pscustomobject]@{
    experiment_id=$LaneConfig.seed.experiment_id
    branch=$LaneConfig.seed.branch
    head_sha=$LaneConfig.seed.head_sha
    workflow_run_id=[long]$LaneConfig.seed.workflow_run_id
    contract_id=$LaneConfig.seed.contract_id
    state="RECONCILE"
    last_scientific_classification=$null
    key_metrics=[pscustomobject]@{}
    repair_attempts=0
    scientist_request=$null
    blocker=$null
  }
}

function Update-Consumed($State, [string]$Lane, [long]$RunId, [string]$ArtifactSha) {
  $key = "$Lane/$RunId"
  $existing = $State.consumed.PSObject.Properties[$key]
  if ($existing -and $existing.Value -ne $ArtifactSha) { throw "COMPLETED_RESULT_REINGEST_IDENTITY_MISMATCH:$key" }
  if (!$existing) { $State.consumed | Add-Member -NotePropertyName $key -NotePropertyValue $ArtifactSha }
}

function Invoke-SafeRepair($Config, $State, [string]$Lane, $Current, [string]$FailureClass, [string]$Stage, [string]$FailingStep, [string]$ErrorSignature) {
  if (Test-NeverAutoRepair $FailureClass) { return [pscustomobject]@{status="HARD_BLOCKER";reason=$FailureClass;attempts=0} }
  if ($FailureClass -notin @("TRANSIENT_GITHUB_API","WORKFLOW_DISPATCH_FAILURE","WORKFLOW_PACKAGING_ERROR","MISSING_READY_RESEARCH_RECEIPT","STALE_METADATA","MISSING_GENERATED_WORKFLOW","CONCURRENT_BRANCH_CREATION","DISPOSABLE_CHECKOUT_CORRUPTION","CACHE_CORRUPTION","RUNNER_INTERRUPTION","SAFE_WORKFLOW_RERUN","MANIFEST_ADMISSION_OMISSION","TEMPORARY_RUNNER_UNAVAILABLE")) {
    return [pscustomobject]@{status="HARD_BLOCKER";reason="UNAPPROVED_REPAIR_CLASS:$FailureClass";attempts=0}
  }
  $fp = Get-RepairFingerprint $Lane $Current.experiment_id $Stage $Current.head_sha $FailureClass $FailingStep $ErrorSignature
  $p = $State.repairs.PSObject.Properties[$fp]
  $n = if ($p) { [int]$p.Value } else { 0 }
  if ($n -ge [int]$Config.max_identical_repair_attempts) { return [pscustomobject]@{status="HARD_BLOCKER";reason="REPAIR_BUDGET_EXHAUSTED:$fp";attempts=$n} }
  $n++
  if ($p) { $p.Value=$n } else { $State.repairs | Add-Member -NotePropertyName $fp -NotePropertyValue $n }
  if (!$DryRun -and !$FixtureRoot -and $FailureClass -in @("RUNNER_INTERRUPTION","SAFE_WORKFLOW_RERUN","TEMPORARY_RUNNER_UNAVAILABLE","WORKFLOW_DISPATCH_FAILURE")) {
    $repo = if ($Lane -eq "Wingless") { $Config.repositories.wingless } else { $Config.repositories.yggdrasil }
    Invoke-Gh @("run","rerun","$($Current.workflow_run_id)","--repo",$repo) | Out-Null
  }
  return [pscustomobject]@{status="REPAIRED";reason=$FailureClass;attempts=$n;fingerprint=$fp}
}

function Get-WorkflowDispatchRuns([string]$Repo, [string]$Workflow, [int]$Limit = 50) {
  $r=Invoke-Gh @("api","repos/$Repo/actions/runs?event=workflow_dispatch&per_page=$Limit")
  $j=$r.stdout|ConvertFrom-Json
  $path=".github/workflows/"+$Workflow
  return @($j.workflow_runs|Where-Object{[string]$_.path -eq $path})
}

function Find-WorkflowDispatchRun([string]$Repo,[string]$Workflow,[string]$ExpectedTitle,[int]$Limit=50) {
  return @(Get-WorkflowDispatchRuns $Repo $Workflow $Limit |
    Where-Object{[string]$_.display_title -eq $ExpectedTitle} |
    Sort-Object{[datetime]$_.created_at} -Descending |
    Select-Object -First 1)
}
function Ensure-Authority($Config, $Contract) {
  $a = $Contract.required_authority
  if ($a.disposition -ne "READY_RESEARCH") { throw "AUTHORITY_DISPOSITION_NOT_READY_RESEARCH" }
  if ($a.production_authority -ne $false) { throw "AUTHORITY_PRODUCTION_FORBIDDEN" }
  if ($a.external_model_calls -ne $false) { throw "AUTHORITY_EXTERNAL_MODEL_CALLS_FORBIDDEN" }
  if ($a.manifest_sha256 -ne $Contract.manifest_sha) { throw "AUTHORITY_MANIFEST_BINDING_MISMATCH" }
  $preissued=$a.PSObject.Properties["preissued_receipt_sha256"]
  if ($preissued -and ![string]::IsNullOrWhiteSpace([string]$preissued.Value)) { return [pscustomobject]@{status="PREISSUED";receipt_sha256=[string]$preissued.Value} }
  if ([string]::IsNullOrWhiteSpace([string]$Contract.preregistration_sha)) { throw "PREREGISTRATION_SHA_REQUIRED_FOR_AUTHORITY" }
  if ([string]::IsNullOrWhiteSpace([string]$a.ckb_plane_sha)) { throw "CKB_PLANE_SHA_REQUIRED_FOR_AUTHORITY" }

  $requestRaw=@($Contract.lane,$Contract.experiment_id,$Contract.preregistration_sha,$Contract.manifest_sha,$a.ckb_plane_sha) -join "|"
  $requestId=(Get-Sha256Text $requestRaw).Substring(0,32)
  $authorityDir=Join-Path $Config.state_root "authority"
  [IO.Directory]::CreateDirectory($authorityDir)|Out-Null
  $cache=Join-Path $authorityDir ($requestId+".json")

  if(Test-Path -LiteralPath $cache){
    $raw=[IO.File]::ReadAllText($cache,[Text.Encoding]::UTF8)
    $j=$raw|ConvertFrom-Json
    if($j.request_id -ne $requestId -or $j.project -ne $Contract.lane -or $j.experiment -ne $Contract.experiment_id -or $j.preregistration_sha -ne $Contract.preregistration_sha -or $j.manifest_sha256 -ne $Contract.manifest_sha -or $j.ckb_plane_main_sha -ne $a.ckb_plane_sha -or $j.research_decision.disposition -ne "READY_RESEARCH"){
      throw "CACHED_AUTHORITY_IDENTITY_MISMATCH:$requestId"
    }
    return [pscustomobject]@{status="CACHED";receipt_sha256=(Get-Sha256Text $raw);receipt_path=$cache;request_id=$requestId}
  }

  if($DryRun -or $FixtureRoot){return [pscustomobject]@{status="DRY_RUN_AUTHORITY";receipt_sha256=$null;receipt_path=$null;request_id=$requestId}}

  $repo=$Config.repositories.ckb_plane
  Invoke-Gh @(
    "workflow","run",$Config.authority_workflow,
    "--repo",$repo,
    "--ref",$Config.authority_workflow_ref,
    "-f","request_id=$requestId",
    "-f","project=$($Contract.lane)",
    "-f","experiment_id=$($Contract.experiment_id)",
    "-f","preregistration_sha=$($Contract.preregistration_sha)",
    "-f","manifest_sha256=$($Contract.manifest_sha)",
    "-f","ckb_plane_sha=$($a.ckb_plane_sha)"
  )|Out-Null

  $deadline=(Get-Date).AddSeconds([int]$Config.authority_timeout_seconds)
  $run=$null
  do{
    $match=@(Find-WorkflowDispatchRun $repo $Config.authority_workflow ("Research authority "+$requestId) 50)
    if($match.Count -eq 1){$run=$match[0]}
    if($run -and [string]$run.status -eq "completed"){break}
    Start-Sleep -Seconds 5
  }while((Get-Date)-lt $deadline)
  if(!$run){throw "AUTHORITY_WORKFLOW_RUN_NOT_FOUND:$requestId"}
  if([string]$run.status -ne "completed"){throw "AUTHORITY_WORKFLOW_TIMEOUT:$requestId"}
  if([string]$run.conclusion -ne "success"){throw "AUTHORITY_WORKFLOW_FAILED:$($requestId):$($run.conclusion)"}

  $tmp=Join-Path $Config.state_root ("authority-download\"+$requestId)
  Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
  [IO.Directory]::CreateDirectory($tmp)|Out-Null
  try{
    Invoke-Gh @("run","download","$($run.id)","--repo",$repo,"--name",("research-authority-"+$requestId),"--dir",$tmp)|Out-Null
    $receiptFile=Get-ChildItem -LiteralPath $tmp -Recurse -File -Filter "research-authority.json"|Select-Object -First 1
    if(!$receiptFile){throw "AUTHORITY_RECEIPT_ARTIFACT_MISSING:$requestId"}
    $bytes=[IO.File]::ReadAllBytes($receiptFile.FullName)
    $raw=[Text.Encoding]::UTF8.GetString($bytes)
    $j=$raw|ConvertFrom-Json
    if($j.schema -ne "ckb-plane.external-exposure-ready-research-receipt.v1" -or
       $j.request_id -ne $requestId -or
       $j.project -ne $Contract.lane -or
       $j.experiment -ne $Contract.experiment_id -or
       $j.preregistration_sha -ne $Contract.preregistration_sha -or
       $j.manifest_sha256 -ne $Contract.manifest_sha -or
       $j.ckb_plane_main_sha -ne $a.ckb_plane_sha -or
       $j.plan_decision.disposition -ne "READY_PLAN" -or
       $j.fetch_decision.disposition -ne "READY_FETCH" -or
       $j.research_decision.disposition -ne "READY_RESEARCH"){
      throw "AUTHORITY_RECEIPT_IDENTITY_MISMATCH:$requestId"
    }
    [IO.File]::WriteAllBytes($cache,$bytes)
    return [pscustomobject]@{status="ISSUED";receipt_sha256=(Get-Sha256Bytes $bytes);receipt_path=$cache;request_id=$requestId}
  }finally{
    Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
  }
}

function Ensure-HostedFocusedQualification($Config,$Contract) {
  $contractIdProp=$Contract.PSObject.Properties["_contract_id"]
  if(!$contractIdProp -or [string]::IsNullOrWhiteSpace([string]$contractIdProp.Value)){throw "QUALIFICATION_CONTRACT_ID_REQUIRED"}
  $contractId=[string]$contractIdProp.Value
  $workflow=[string]$Config.qualification_workflow
  $workflowRef=[string]$Config.qualification_workflow_ref
  $repo=[string]$Config.repositories.ckb_plane
  $qualificationSourceSha=if($Config.PSObject.Properties["qualification_source_sha"] -and ![string]::IsNullOrWhiteSpace([string]$Config.qualification_source_sha)){[string]$Config.qualification_source_sha}else{[string]$Config.controller_sha}
  $requestRaw=@([string]$Config.controller_sha,$qualificationSourceSha,[string]$Contract.lane,[string]$Contract.experiment_id,[string]$Contract.preregistration_sha,[string]$Contract.manifest_sha,$contractId) -join "|"
  $requestId=(Get-Sha256Text $requestRaw).Substring(0,32)
  $dir=Join-Path $Config.state_root "focused-qualification"
  [IO.Directory]::CreateDirectory($dir)|Out-Null
  $cache=Join-Path $dir ($requestId+".json")
  if(Test-Path -LiteralPath $cache){
    $raw=[IO.File]::ReadAllText($cache,[Text.Encoding]::UTF8)
    $j=$raw|ConvertFrom-Json
    if($j.request_id -ne $requestId -or $j.controller_sha -ne $Config.controller_sha -or !$j.PSObject.Properties["qualification_source_sha"] -or $j.qualification_source_sha -ne $qualificationSourceSha -or $j.contract_id -ne $contractId -or $j.experiment_id -ne $Contract.experiment_id -or $j.preregistration_sha -ne $Contract.preregistration_sha -or $j.manifest_sha256 -ne $Contract.manifest_sha -or $j.qualification -ne "PASS"){
      throw "CACHED_FOCUSED_QUALIFICATION_IDENTITY_MISMATCH:$requestId"
    }
    return [pscustomobject]@{status="CACHED";request_id=$requestId;receipt_path=$cache;receipt_sha256=(Get-Sha256Text $raw)}
  }
  if($DryRun -or $FixtureRoot){return [pscustomobject]@{status="DRY_RUN";request_id=$requestId}}
  Invoke-Gh @(
    "workflow","run",$workflow,
    "--repo",$repo,
    "--ref",$workflowRef,
    "-f","request_id=$requestId",
    "-f","controller_sha=$($Config.controller_sha)",
    "-f","qualification_source_sha=$qualificationSourceSha",
    "-f","contract_id=$contractId"
  )|Out-Null
  $deadline=(Get-Date).AddSeconds([int]$Config.qualification_timeout_seconds)
  $run=$null
  do{
    $match=@(Find-WorkflowDispatchRun $repo $workflow ("Focused research qualification "+$requestId) 50)
    if($match.Count -eq 1){$run=$match[0]}
    if($run -and [string]$run.status -eq "completed"){break}
    Start-Sleep -Seconds 5
  }while((Get-Date)-lt $deadline)
  if(!$run){throw "FOCUSED_QUALIFICATION_RUN_NOT_FOUND:$requestId"}
  if([string]$run.status -ne "completed"){throw "FOCUSED_QUALIFICATION_TIMEOUT:$requestId"}
  if([string]$run.conclusion -ne "success"){throw "FOCUSED_QUALIFICATION_FAILED:$($requestId):$($run.conclusion)"}
  $tmp=Join-Path $Config.state_root ("focused-qualification-download\"+$requestId)
  Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
  [IO.Directory]::CreateDirectory($tmp)|Out-Null
  try{
    Invoke-Gh @("run","download","$($run.id)","--repo",$repo,"--name",("research-focused-qualification-"+$requestId),"--dir",$tmp)|Out-Null
    $receiptFile=Get-ChildItem -LiteralPath $tmp -Recurse -File -Filter "research-focused-qualification.json"|Select-Object -First 1
    if(!$receiptFile){throw "FOCUSED_QUALIFICATION_RECEIPT_MISSING:$requestId"}
    $bytes=[IO.File]::ReadAllBytes($receiptFile.FullName)
    $raw=[Text.Encoding]::UTF8.GetString($bytes)
    $j=$raw|ConvertFrom-Json
    if($j.schema -ne "ckb-plane.research-focused-qualification.v1" -or $j.request_id -ne $requestId -or $j.controller_sha -ne $Config.controller_sha -or !$j.PSObject.Properties["qualification_source_sha"] -or $j.qualification_source_sha -ne $qualificationSourceSha -or $j.contract_id -ne $contractId -or $j.experiment_id -ne $Contract.experiment_id -or $j.preregistration_sha -ne $Contract.preregistration_sha -or $j.manifest_sha256 -ne $Contract.manifest_sha -or $j.qualification -ne "PASS" -or $j.codex_used -ne $false){
      throw "FOCUSED_QUALIFICATION_RECEIPT_IDENTITY_MISMATCH:$requestId"
    }
    [IO.File]::WriteAllBytes($cache,$bytes)
    return [pscustomobject]@{status="HOSTED";request_id=$requestId;receipt_path=$cache;receipt_sha256=(Get-Sha256Bytes $bytes)}
  }finally{
    Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
  }
}

function Invoke-FocusedQualification($Config,$Contract,[string]$WorkingDirectory) {
  $commands=@($Contract.qualification_profile.focused_commands)
  $python=@($commands|Where-Object{
    $exe=([string]$_.exe).ToLowerInvariant()
    $exe -eq "python.exe" -or $exe -eq "python3.exe"
  })
  if($python.Count -gt 0){
    if($python.Count -ne $commands.Count){throw "HOSTED_QUALIFICATION_MIXED_COMMAND_SET_FORBIDDEN"}
    foreach($cmd in $commands){
      $args=@($cmd.args|ForEach-Object{[string]$_})
      if($args.Count -ne 3 -or $args[0] -ne "-m" -or $args[1] -ne "py_compile"){throw "HOSTED_QUALIFICATION_UNSUPPORTED_PYTHON_COMMAND"}
    }
    return Ensure-HostedFocusedQualification $Config $Contract
  }
  foreach($cmd in $commands){
    Invoke-ProcessStrict -FilePath ([string]$cmd.exe) -Arguments @($cmd.args|ForEach-Object{[string]$_}) -WorkingDirectory $WorkingDirectory|Out-Null
  }
  return [pscustomobject]@{status="LOCAL"}
}
function New-ResearchWorktree($Config, $Contract, [string]$RepoPath) {
  $root = Join-Path $Config.state_root "worktrees"
  [IO.Directory]::CreateDirectory($root) | Out-Null
  $worktreeIdentity=@([string]$Contract.lane,[string]$Contract.experiment_id,[string]$Contract.preregistration_sha,[string]$Contract.parent_sha) -join "|"
  $slug = "wt-" + (Get-Sha256Text $worktreeIdentity).Substring(0,24)
  $wt = Join-Path $root $slug
  if (Test-Path -LiteralPath $wt) {
    try { Invoke-Git $RepoPath @("worktree","remove","--force",$wt) | Out-Null } catch {}
    Remove-Item -LiteralPath $wt -Recurse -Force -ErrorAction SilentlyContinue
  }
  Invoke-Git $RepoPath @("config","core.longpaths","true") | Out-Null
  Invoke-Git $RepoPath @("fetch","origin","--prune") | Out-Null
  Invoke-Git $RepoPath @("cat-file","-e","$($Contract.parent_sha)^{commit}") | Out-Null
  $base=[string]$Contract.parent_sha
  if($Contract.PSObject.Properties["preregistration_sha"] -and ![string]::IsNullOrWhiteSpace([string]$Contract.preregistration_sha)){
    Invoke-Git $RepoPath @("cat-file","-e","$($Contract.preregistration_sha)^{commit}") | Out-Null
    Invoke-Git $RepoPath @("merge-base","--is-ancestor",$Contract.parent_sha,$Contract.preregistration_sha) | Out-Null
    $base=[string]$Contract.preregistration_sha
  }
  Invoke-Git $RepoPath @("worktree","add","--detach",$wt,$base) | Out-Null
  return $wt
}

function Test-SuccessorTemplate($Config, $Contract, [string]$RepoPath) {
  Assert-Contract $Contract
  $wt=New-ResearchWorktree $Config $Contract $RepoPath
  try{
    foreach($file in @($Contract.implementation_template.files)){
      $src=Join-Path $Config.template_root $file.template
      if(!(Test-Path -LiteralPath $src)){throw "IMPLEMENTATION_TEMPLATE_MISSING:$src"}
      $rendered=Render-Template (Get-Content -LiteralPath $src -Raw) $Contract.implementation_template.variables
      $target=Join-Path $wt $file.target
      [IO.Directory]::CreateDirectory((Split-Path -Parent $target))|Out-Null
      Write-AtomicUtf8 $target $rendered
    }
    [void](Invoke-FocusedQualification $Config $Contract $wt)
    return [pscustomobject]@{status="QUALIFIED";base_sha=(Invoke-Git $wt @("rev-parse","HEAD")).stdout}
  } finally {
    try{Invoke-Git $RepoPath @("worktree","remove","--force",$wt)|Out-Null}catch{}
  }
}

function Render-Successor($Config, $Contract, [string]$RepoPath) {
  Assert-Contract $Contract
  foreach ($file in @($Contract.implementation_template.files)) {
    $src = Join-Path $Config.template_root $file.template
    if (!(Test-Path -LiteralPath $src)) { throw "IMPLEMENTATION_TEMPLATE_MISSING:$src" }
    $raw = Get-Content -LiteralPath $src -Raw
    [void](Render-Template $raw $Contract.implementation_template.variables)
  }
  if ($DryRun -or $FixtureRoot) { return [pscustomobject]@{status="DRY_RUN_RENDERED";head_sha=$Contract.preregistered_head_sha;branch=$Contract.branch} }
  $wt = New-ResearchWorktree $Config $Contract $RepoPath
  try {
    foreach ($file in @($Contract.implementation_template.files)) {
      $src = Join-Path $Config.template_root $file.template
      $rendered = Render-Template (Get-Content -LiteralPath $src -Raw) $Contract.implementation_template.variables
      $target = Join-Path $wt $file.target
      [IO.Directory]::CreateDirectory((Split-Path -Parent $target)) | Out-Null
      Write-AtomicUtf8 $target $rendered
    }
    [void](Invoke-FocusedQualification $Config $Contract $wt)
    $status = Invoke-Git $wt @("status","--porcelain")
    if (!$status.stdout) { throw "RENDER_PRODUCED_NO_DELTA" }
    Invoke-Git $wt @("add","--all") | Out-Null
    $oldAuthorDate=$env:GIT_AUTHOR_DATE
    $oldCommitterDate=$env:GIT_COMMITTER_DATE
    $oldAuthorName=$env:GIT_AUTHOR_NAME
    $oldAuthorEmail=$env:GIT_AUTHOR_EMAIL
    $oldCommitterName=$env:GIT_COMMITTER_NAME
    $oldCommitterEmail=$env:GIT_COMMITTER_EMAIL
    try {
      $env:GIT_AUTHOR_DATE=[string]$Contract.preregistered_at_utc
      $env:GIT_COMMITTER_DATE=[string]$Contract.preregistered_at_utc
      $env:GIT_AUTHOR_NAME="CKB Research Controller"
      $env:GIT_AUTHOR_EMAIL="ckb-research-controller@local.invalid"
      $env:GIT_COMMITTER_NAME="CKB Research Controller"
      $env:GIT_COMMITTER_EMAIL="ckb-research-controller@local.invalid"
      Invoke-Git $wt @("commit","-m","research: preregister $($Contract.experiment_id)") | Out-Null
    } finally {
      $env:GIT_AUTHOR_DATE=$oldAuthorDate
      $env:GIT_COMMITTER_DATE=$oldCommitterDate
      $env:GIT_AUTHOR_NAME=$oldAuthorName
      $env:GIT_AUTHOR_EMAIL=$oldAuthorEmail
      $env:GIT_COMMITTER_NAME=$oldCommitterName
      $env:GIT_COMMITTER_EMAIL=$oldCommitterEmail
    }
    $head = (Invoke-Git $wt @("rev-parse","HEAD")).stdout
    $remote = Invoke-Git $RepoPath @("ls-remote","--heads","origin","refs/heads/$($Contract.branch)")
    if ($remote.stdout) {
      $existing = ($remote.stdout -split '\s+')[0]
      if ($existing -eq $head) {
      } elseif ($existing -eq [string]$Contract.preregistration_sha) {
        Invoke-Git $wt @("push","origin","HEAD:refs/heads/$($Contract.branch)") | Out-Null
      } else {
        throw "CONCURRENT_BRANCH_IDENTITY_MISMATCH:$($Contract.branch):${existing}:$head"
      }
    } else {
      Invoke-Git $wt @("push","origin","HEAD:refs/heads/$($Contract.branch)") | Out-Null
    }
    return [pscustomobject]@{status="PUSHED";head_sha=$head;branch=$Contract.branch}
  } finally {
    try { Invoke-Git $RepoPath @("worktree","remove","--force",$wt) | Out-Null } catch {}
  }
}

function Repair-WorkflowPackaging($Config, [string]$Lane, $Contract, [string]$RepoPath, [string]$ExpectedHeadSha) {
  Assert-Contract $Contract
  if($DryRun -or $FixtureRoot){
    return [pscustomobject]@{status="DRY_RUN";head_sha=$ExpectedHeadSha;run=$null}
  }
  $workflowFiles=@($Contract.implementation_template.files|Where-Object{
    (([string]$_.target).Replace("\\","/")).StartsWith(".github/workflows/")
  })
  if($workflowFiles.Count -ne 1){throw "WORKFLOW_PACKAGING_REPAIR_TARGET_COUNT_INVALID:$($workflowFiles.Count)"}

  $authority=Ensure-Authority $Config $Contract
  $vars=$Contract.implementation_template.variables
  if($null -eq $vars){
    $vars=[pscustomobject]@{}
    $Contract.implementation_template.variables=$vars
  }
  if($authority.receipt_sha256){
    $vars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_SHA256" -NotePropertyValue ([string]$authority.receipt_sha256) -Force
  }
  if($authority.PSObject.Properties["receipt_path"] -and $authority.receipt_path){
    $bytes=[IO.File]::ReadAllBytes([string]$authority.receipt_path)
    $vars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_B64" -NotePropertyValue ([Convert]::ToBase64String($bytes)) -Force
  }

  Invoke-Git $RepoPath @("fetch","origin",$Contract.branch,"--prune")|Out-Null
  $remote=(Invoke-Git $RepoPath @("rev-parse","origin/$($Contract.branch)")).stdout
  if($remote -ne $ExpectedHeadSha){throw "WORKFLOW_PACKAGING_REPAIR_HEAD_DRIFT:${ExpectedHeadSha}:$remote"}

  $root=Join-Path $Config.state_root "worktrees"
  [IO.Directory]::CreateDirectory($root)|Out-Null
  $wt=Join-Path $root ("workflow-repair-"+$Lane.ToLowerInvariant()+"-"+$ExpectedHeadSha.Substring(0,12))
  if(Test-Path -LiteralPath $wt){throw "WORKFLOW_PACKAGING_REPAIR_WORKTREE_EXISTS:$wt"}
  Invoke-Git $RepoPath @("worktree","add","--detach",$wt,$ExpectedHeadSha)|Out-Null
  try{
    foreach($file in $workflowFiles){
      $src=Join-Path $Config.template_root $file.template
      if(!(Test-Path -LiteralPath $src)){throw "IMPLEMENTATION_TEMPLATE_MISSING:$src"}
      $rendered=Render-Template (Get-Content -LiteralPath $src -Raw) $vars
      $target=Join-Path $wt $file.target
      [IO.Directory]::CreateDirectory((Split-Path -Parent $target))|Out-Null
      Write-AtomicUtf8 $target $rendered
    }
    $status=(Invoke-Git $wt @("status","--porcelain")).stdout
    if(!$status){throw "WORKFLOW_PACKAGING_REPAIR_NO_DELTA"}
    $changed=@((Invoke-Git $wt @("diff","--name-only")).stdout -split "[\r\n]+"|Where-Object{$_})
    $allowed=@($workflowFiles|ForEach-Object{([string]$_.target).Replace("\\","/")})
    foreach($p in $changed){
      if($allowed -notcontains $p.Replace("\\","/")){throw "WORKFLOW_PACKAGING_REPAIR_SCOPE_WIDENED:$p"}
    }
    [void](Invoke-FocusedQualification $Config $Contract $wt)
    foreach($file in $workflowFiles){Invoke-Git $wt @("add","--",$file.target)|Out-Null}

    $oldAuthorDate=$env:GIT_AUTHOR_DATE
    $oldCommitterDate=$env:GIT_COMMITTER_DATE
    $oldAuthorName=$env:GIT_AUTHOR_NAME
    $oldAuthorEmail=$env:GIT_AUTHOR_EMAIL
    $oldCommitterName=$env:GIT_COMMITTER_NAME
    $oldCommitterEmail=$env:GIT_COMMITTER_EMAIL
    try{
      $env:GIT_AUTHOR_DATE=[string]$Contract.preregistered_at_utc
      $env:GIT_COMMITTER_DATE=[string]$Contract.preregistered_at_utc
      $env:GIT_AUTHOR_NAME="CKB Research Controller"
      $env:GIT_AUTHOR_EMAIL="ckb-research-controller@local.invalid"
      $env:GIT_COMMITTER_NAME="CKB Research Controller"
      $env:GIT_COMMITTER_EMAIL="ckb-research-controller@local.invalid"
      Invoke-Git $wt @("commit","-m","repair: regenerate hosted workflow packaging")|Out-Null
    }finally{
      $env:GIT_AUTHOR_DATE=$oldAuthorDate
      $env:GIT_COMMITTER_DATE=$oldCommitterDate
      $env:GIT_AUTHOR_NAME=$oldAuthorName
      $env:GIT_AUTHOR_EMAIL=$oldAuthorEmail
      $env:GIT_COMMITTER_NAME=$oldCommitterName
      $env:GIT_COMMITTER_EMAIL=$oldCommitterEmail
    }
    $head=(Invoke-Git $wt @("rev-parse","HEAD")).stdout
    $remoteBefore=(Invoke-Git $RepoPath @("ls-remote","--heads","origin","refs/heads/$($Contract.branch)")).stdout
    $remoteSha=if($remoteBefore){($remoteBefore -split '\\s+')[0]}else{""}
    if($remoteSha -ne $ExpectedHeadSha){throw "WORKFLOW_PACKAGING_REPAIR_PUSH_RACE:${ExpectedHeadSha}:$remoteSha"}
    Invoke-Git $wt @("push","origin","HEAD:refs/heads/$($Contract.branch)")|Out-Null
    $run=Verify-HostedRunExists $Config $Lane $Contract $head
    return [pscustomobject]@{status="REPAIRED";head_sha=$head;run=$run}
  }finally{
    try{Invoke-Git $RepoPath @("worktree","remove","--force",$wt)|Out-Null}catch{}
  }
}

function Verify-HostedRunExists($Config, [string]$Lane, $Contract, [string]$HeadSha) {
  $repo = if ($Lane -eq "Wingless") { $Config.repositories.wingless } else { $Config.repositories.yggdrasil }
  $deadline = (Get-Date).AddSeconds([int]$Config.hosted_run_discovery_timeout_seconds)
  do {
    $match = @(Get-WorkflowRuns $repo $Contract.branch 20 | Where-Object { $_.head_sha -eq $HeadSha } | Sort-Object {[datetime]$_.created_at} -Descending | Select-Object -First 1)
    if ($match.Count -eq 1) { return $match[0] }
    if ($DryRun -or $FixtureRoot) { break }
    Start-Sleep -Seconds 5
  } while ((Get-Date) -lt $deadline)
  throw "HOSTED_RUN_NOT_FOUND:${Lane}:$($Contract.experiment_id):$HeadSha"
}

function Advance-Lane($Config, $State, [string]$Lane, $Current, [string]$SuccessorContractId) {
  $contract = Load-Contract $Config $SuccessorContractId
  if ($contract.parent_experiment -ne $Current.experiment_id) { throw "SUCCESSOR_PARENT_EXPERIMENT_MISMATCH" }
  if ($contract.parent_sha -ne $Current.head_sha) { throw "SUCCESSOR_PARENT_SHA_MISMATCH" }
  if (!$contract.implementation_template.ready) {
    $req = New-ScientistRequest $Lane $Current "NOVEL_HYPOTHESIS_REQUIRED" $Current.last_scientific_classification $contract @("successor-contract:$SuccessorContractId","implementation-template:not-ready")
    $Current.state="NOVEL_HYPOTHESIS_REQUIRED"
    $saved=Save-ScientistRequest $Config.state_root $req
    Publish-ScientistRequest $Config $req $saved
    $Current.scientist_request=$saved
    return $Current
  }
  $repoPath = if ($Lane -eq "Wingless") { $Config.local_repositories.wingless } else { $Config.local_repositories.yggdrasil }
  $pre=($contract|ConvertTo-Json -Depth 100|ConvertFrom-Json)
  if($null -eq $pre.implementation_template.variables){
    $pre.implementation_template.variables=[pscustomobject]@{}
  }
  $preVars=$pre.implementation_template.variables
  $preVars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_SHA256" -NotePropertyValue ("0"*64) -Force
  $preVars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_B64" -NotePropertyValue "e30=" -Force
  [void](Test-SuccessorTemplate $Config $pre $repoPath)
  $authority=Ensure-Authority $Config $contract
  $vars=$contract.implementation_template.variables
  if($null -eq $vars){
    $vars=[pscustomobject]@{}
    $contract.implementation_template.variables=$vars
  }
  if($authority.receipt_sha256){
    $vars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_SHA256" -NotePropertyValue ([string]$authority.receipt_sha256) -Force
  }
  if($authority.PSObject.Properties["receipt_path"] -and $authority.receipt_path){
    $bytes=[IO.File]::ReadAllBytes([string]$authority.receipt_path)
    $vars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_B64" -NotePropertyValue ([Convert]::ToBase64String($bytes)) -Force
  }
  $render = Render-Successor $Config $contract $repoPath
  $run = Verify-HostedRunExists $Config $Lane $contract $render.head_sha
  $Current.experiment_id=$contract.experiment_id
  $Current.branch=$contract.branch
  $Current.head_sha=$render.head_sha
  $Current.workflow_run_id=[long]$run.id
  $Current.contract_id=$SuccessorContractId
  $Current.state="RUNNING"
  $Current.repair_attempts=0
  $Current.scientist_request=$null
  $Current.blocker=$null
  return $Current
}

function Try-AdoptVerifiedWorkflowPackagingDrift($Config, [string]$Lane, $Current, [string]$RepoPath, [string]$ExpectedHeadSha, [string]$RemoteHeadSha) {
  if([string]::IsNullOrWhiteSpace($ExpectedHeadSha) -or [string]::IsNullOrWhiteSpace($RemoteHeadSha) -or $ExpectedHeadSha -eq $RemoteHeadSha){return $null}
  $contract=Load-Contract $Config ([string]$Current.contract_id)
  if(!$contract.implementation_template.ready){return $null}
  $workflowFiles=@($contract.implementation_template.files|Where-Object{([string]$_.target).Replace("\","/").StartsWith(".github/workflows/")})
  if($workflowFiles.Count -eq 0){return $null}

  try{
    Invoke-Git $RepoPath @("fetch","origin",$contract.branch,"--prune")|Out-Null
    Invoke-Git $RepoPath @("merge-base","--is-ancestor",$ExpectedHeadSha,$RemoteHeadSha)|Out-Null
  }catch{return $null}

  $changed=@((Invoke-Git $RepoPath @("diff","--name-only",($ExpectedHeadSha+".."+$RemoteHeadSha))).stdout -split "[\r\n]+"|Where-Object{$_}|ForEach-Object{$_.Replace("\","/")})
  if($changed.Count -eq 0){return $null}
  $allowed=@($workflowFiles|ForEach-Object{([string]$_.target).Replace("\","/")})
  foreach($p in $changed){if($allowed -notcontains $p){return $null}}

  $a=$contract.required_authority
  $requestRaw=@($contract.lane,$contract.experiment_id,$contract.preregistration_sha,$contract.manifest_sha,$a.ckb_plane_sha) -join "|"
  $requestId=(Get-Sha256Text $requestRaw).Substring(0,32)
  $cache=Join-Path $Config.state_root ("authority\"+$requestId+".json")
  if(!(Test-Path -LiteralPath $cache)){return $null}
  try{$authority=Ensure-Authority $Config $contract}catch{return $null}
  if(!$authority.receipt_sha256 -or !$authority.receipt_path){return $null}

  $vars=($contract.implementation_template.variables|ConvertTo-Json -Depth 100|ConvertFrom-Json)
  if($null -eq $vars){$vars=[pscustomobject]@{}}
  $vars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_SHA256" -NotePropertyValue ([string]$authority.receipt_sha256) -Force
  $bytes=[IO.File]::ReadAllBytes([string]$authority.receipt_path)
  $vars|Add-Member -NotePropertyName "AUTHORITY_RECEIPT_B64" -NotePropertyValue ([Convert]::ToBase64String($bytes)) -Force

  $verifyRoot=Join-Path $Config.state_root ("workflow-drift-verify\"+$Lane.ToLowerInvariant()+"-"+$RemoteHeadSha.Substring(0,12))
  Remove-Item -LiteralPath $verifyRoot -Recurse -Force -ErrorAction SilentlyContinue
  [IO.Directory]::CreateDirectory($verifyRoot)|Out-Null
  try{
    foreach($file in $workflowFiles){
      $src=Join-Path $Config.template_root $file.template
      if(!(Test-Path -LiteralPath $src)){return $null}
      $rendered=Render-Template (Get-Content -LiteralPath $src -Raw) $vars
      $tmp=Join-Path $verifyRoot ([IO.Path]::GetRandomFileName())
      Write-AtomicUtf8 $tmp $rendered
      $expectedBlob=(Invoke-Git $RepoPath @("hash-object",$tmp)).stdout
      $target=([string]$file.target).Replace("\","/")
      $remoteBlob=(Invoke-Git $RepoPath @("rev-parse",($RemoteHeadSha+":"+$target))).stdout
      if($expectedBlob -ne $remoteBlob){return $null}
    }
  }catch{return $null}
  finally{Remove-Item -LiteralPath $verifyRoot -Recurse -Force -ErrorAction SilentlyContinue}

  $run=Get-WorkflowRuns (if($Lane -eq "Wingless"){$Config.repositories.wingless}else{$Config.repositories.yggdrasil}) $contract.branch 20 |
    Where-Object{$_.head_sha -eq $RemoteHeadSha -and $_.status -eq "completed" -and $_.conclusion -eq "success"} |
    Sort-Object {[datetime]$_.created_at} -Descending | Select-Object -First 1
  if($null -eq $run){return $null}
  return [pscustomobject]@{head_sha=$RemoteHeadSha;workflow_run_id=[long]$run.id;status="ADOPTED_VERIFIED_WORKFLOW_PACKAGING_DRIFT"}
}

function Reconcile-Lane($Config, $State, [string]$Lane, $LaneConfig) {
  $prop = $State.lanes.PSObject.Properties[$Lane]
  $current = if ($prop) { $prop.Value } else { Get-InitialLaneState $LaneConfig }
  if (!$prop) { $State.lanes | Add-Member -NotePropertyName $Lane -NotePropertyValue $current }

  $repo = if ($Lane -eq "Wingless") { $Config.repositories.wingless } else { $Config.repositories.yggdrasil }
  $branch = Get-RemoteBranch $repo $current.branch
  if ($null -eq $branch) { $current.state="HARD_BLOCKER";$current.blocker="AUTHORITATIVE_BRANCH_MISSING:$($current.branch)";return $current }
  $remoteSha=[string]$branch.commit.sha
  if ($current.head_sha -and $remoteSha -ne $current.head_sha) {
    $repoPath=if($Lane -eq "Wingless"){$Config.local_repositories.wingless}else{$Config.local_repositories.yggdrasil}
    $adopt=Try-AdoptVerifiedWorkflowPackagingDrift $Config $Lane $current $repoPath ([string]$current.head_sha) $remoteSha
    if($null -eq $adopt){
      $current.state="HARD_BLOCKER";$current.blocker="PROVENANCE_AMBIGUITY:HEAD_DRIFT:$($current.head_sha):$remoteSha";return $current
    }
    $current.head_sha=[string]$adopt.head_sha
    $current.workflow_run_id=[long]$adopt.workflow_run_id
    $current.blocker=$null
  }
  $current.head_sha=$remoteSha

  $runs=@(Get-WorkflowRuns $repo $current.branch 20)
  $run=$null
  if ($current.workflow_run_id) { $run=$runs | Where-Object { [long]$_.id -eq [long]$current.workflow_run_id } | Select-Object -First 1 }
  if ($null -eq $run) { $run=$runs | Where-Object { $_.head_sha -eq $remoteSha } | Sort-Object {[datetime]$_.created_at} -Descending | Select-Object -First 1 }
  if ($null -eq $run) { $current.state="HARD_BLOCKER";$current.blocker="AUTHORITATIVE_RUN_MISSING:$($current.branch):$remoteSha";return $current }
  $current.workflow_run_id=[long]$run.id

  if ($run.status -ne "completed") { $current.state="RUNNING";$current.blocker=$null;return $current }
  if ($run.conclusion -ne "success") {
    $failure=Classify-WorkflowFailure $Config $Lane $run
    if ($failure.failure_class -eq "UNCLASSIFIED_WORKFLOW_FAILURE") {
      $current.state="HARD_BLOCKER"
      $current.blocker="UNCLASSIFIED_WORKFLOW_FAILURE:$($failure.error_signature)"
      return $current
    }
    $repair=Invoke-SafeRepair $Config $State $Lane $current ([string]$failure.failure_class) "reconcile" ([string]$failure.failing_step) ([string]$failure.error_signature)
    if ($repair.status -eq "REPAIRED") {
      if([string]$failure.failure_class -eq "WORKFLOW_PACKAGING_ERROR"){
        $contract=Load-Contract $Config $current.contract_id
        $repoPath=if($Lane -eq "Wingless"){$Config.local_repositories.wingless}else{$Config.local_repositories.yggdrasil}
        $pack=Repair-WorkflowPackaging $Config $Lane $contract $repoPath ([string]$current.head_sha)
        $current.head_sha=[string]$pack.head_sha
        $current.workflow_run_id=[long]$pack.run.id
      }
      $current.state="RUNNING";$current.repair_attempts=$repair.attempts;$current.blocker=$null
    }
    else { $current.state="HARD_BLOCKER";$current.blocker=$repair.reason }
    return $current
  }

  $consumedKey="$Lane/$($run.id)"
  $consumedProp=$State.consumed.PSObject.Properties[$consumedKey]
  if($consumedProp -and $current.last_scientific_classification){
    $contract=Load-Contract $Config $current.contract_id
    $next=Resolve-Successor $contract ([string]$current.last_scientific_classification)
    if($next -and $next.contract_id -and (Test-Path -LiteralPath (Get-ContractPath $Config ([string]$next.contract_id)))){
      return Advance-Lane $Config $State $Lane $current ([string]$next.contract_id)
    }
    if($current.state -eq "NOVEL_HYPOTHESIS_REQUIRED" -and $current.scientist_request){return $current}
  }

  $tmp=Join-Path $Config.state_root ("ingest\"+$Lane+"-"+$run.id)
  Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
  Get-RunArtifacts $repo ([long]$run.id) $tmp
  $cls=Find-ClassificationArtifact $tmp $current.experiment_id
  $clsJson=ConvertTo-CanonicalJson $cls
  $clsSha=Get-Sha256Text $clsJson
  Update-Consumed $State $Lane ([long]$run.id) $clsSha

  $classification=[string]$cls.classification
  if ($classification -notin @("supported","mixed","negative","invalid")) { throw "CLASSIFICATION_VALUE_INVALID:$classification" }
  $current.last_scientific_classification=$classification
  $current.key_metrics=$cls.metrics
  $contract=Load-Contract $Config $current.contract_id

  if ($classification -eq "invalid") {
    if ($contract.invalid_behavior -eq "INFRA_INVALID" -and $cls.PSObject.Properties["failure_class"] -and $cls.failure_class) {
      $repair=Invoke-SafeRepair $Config $State $Lane $current ([string]$cls.failure_class) "classification" ([string]$cls.failing_step) ([string]$cls.error_signature)
      if($repair.status -eq "REPAIRED"){$current.state="RUNNING";$current.repair_attempts=$repair.attempts;$current.blocker=$null}
      else{$current.state="HARD_BLOCKER";$current.blocker=$repair.reason}
    } elseif ($contract.invalid_behavior -eq "INFRA_INVALID") {
      $current.state="HARD_BLOCKER"
      $current.blocker="INFRA_INVALID_REQUIRES_EXACT_FAILURE_CLASSIFICATION"
    } else {
      $current.state="HARD_BLOCKER"
      $current.blocker="INVALID_RESULT_NOT_AUTHORIZED_FOR_REPAIR"
    }
    return $current
  }

  $next=Resolve-Successor $contract $classification
  if ($null -eq $next -or [string]::IsNullOrWhiteSpace([string]$next.contract_id) -or !(Test-Path -LiteralPath (Get-ContractPath $Config ([string]$next.contract_id)))) {
    $req=New-ScientistRequest $Lane $current "NOVEL_HYPOTHESIS_REQUIRED" $classification $contract @("github-run:$($run.id)","classification-sha256:$clsSha")
    $current.state="NOVEL_HYPOTHESIS_REQUIRED"
    $saved=Save-ScientistRequest $Config.state_root $req
    Publish-ScientistRequest $Config $req $saved
    $current.scientist_request=$saved
    $current.blocker=$null
    return $current
  }
  return Advance-Lane $Config $State $Lane $current ([string]$next.contract_id)
}

function Invoke-LaneSafely($Config, $State, [string]$Lane, $LaneConfig) {
  try {
    return Reconcile-Lane $Config $State $Lane $LaneConfig
  } catch {
    $p=$State.lanes.PSObject.Properties[$Lane]
    $current=if($p){$p.Value}else{Get-InitialLaneState $LaneConfig}
    if(!$p){$State.lanes|Add-Member -NotePropertyName $Lane -NotePropertyValue $current}
    $current.state="HARD_BLOCKER"
    $current.blocker="CONTROLLER_EXCEPTION:"+$_.Exception.Message
    return $current
  }
}

function Invoke-SelfTest {
  $failed=New-Object Collections.Generic.List[string]
  function Check([bool]$Ok,[string]$Name){if(!$Ok){[void]$failed.Add($Name)}}
  Check ((Get-Sha256Text "abc") -eq "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad") "sha256"
  $vars=[pscustomobject]@{A="x";B="y"}
  Check ((Render-Template "{{A}}-{{B}}" $vars) -eq "x-y") "template"
  $fp1=Get-RepairFingerprint "Wingless" "E1" "run" "abc" "RUNNER_INTERRUPTION" "job" "sig"
  $fp2=Get-RepairFingerprint "Wingless" "E1" "run" "abc" "RUNNER_INTERRUPTION" "job" "sig"
  Check ($fp1 -eq $fp2) "repair-fingerprint"
  Check (Test-NeverAutoRepair "SCIENTIFIC_NEGATIVE") "negative-never-repair"
  Check (Test-NeverAutoRepair "PROVENANCE_AMBIGUITY") "provenance-never-repair"
  Check (!(Test-NeverAutoRepair "RUNNER_INTERRUPTION")) "runner-repairable"
  $c=[pscustomobject]@{successors=[pscustomobject]@{supported=[pscustomobject]@{contract_id="A"};mixed=[pscustomobject]@{contract_id="B"};negative=[pscustomobject]@{contract_id="C"}}}
  Check ((Resolve-Successor $c "supported").contract_id -eq "A") "supported-map"
  Check ((Resolve-Successor $c "mixed").contract_id -eq "B") "mixed-map"
  Check ((Resolve-Successor $c "negative").contract_id -eq "C") "negative-map"
  Check ($null -eq (Resolve-Successor $c "invalid")) "invalid-not-scientific-successor"

  $oldFixture=$script:FixtureRoot
  try {
    $script:FixtureRoot="selftest"
    $rcfg=[pscustomobject]@{max_identical_repair_attempts=3;repositories=[pscustomobject]@{wingless="x";yggdrasil="y"}}
    $rs=[pscustomobject]@{repairs=[pscustomobject]@{}}
    $cur=[pscustomobject]@{experiment_id="E1";head_sha="abc";workflow_run_id=1}
    $d1=Invoke-SafeRepair $rcfg $rs "Wingless" $cur "RUNNER_INTERRUPTION" "run" "runner" "sig"
    $d2=Invoke-SafeRepair $rcfg $rs "Wingless" $cur "RUNNER_INTERRUPTION" "run" "runner" "sig"
    $d3=Invoke-SafeRepair $rcfg $rs "Wingless" $cur "RUNNER_INTERRUPTION" "run" "runner" "sig"
    $d4=Invoke-SafeRepair $rcfg $rs "Wingless" $cur "RUNNER_INTERRUPTION" "run" "runner" "sig"
    Check ($d1.status -eq "REPAIRED" -and $d2.status -eq "REPAIRED" -and $d3.status -eq "REPAIRED" -and $d4.status -eq "HARD_BLOCKER" -and $d4.attempts -eq 3) "repair-budget-three"
  } finally {
    $script:FixtureRoot=$oldFixture
  }

  $root=Join-Path $env:TEMP ("ckb-research-render-selftest-"+[guid]::NewGuid().ToString("N"))
  try {
    $remote=Join-Path $root "remote.git"
    $seed=Join-Path $root "seed"
    $templates=Join-Path $root "templates"
    $stateRoot=Join-Path $root "state"
    [IO.Directory]::CreateDirectory($root)|Out-Null
    [IO.Directory]::CreateDirectory($templates)|Out-Null
    Invoke-ProcessStrict -FilePath "git.exe" -Arguments @("init","--bare",$remote)|Out-Null
    Invoke-ProcessStrict -FilePath "git.exe" -Arguments @("clone",$remote,$seed)|Out-Null
    [IO.File]::WriteAllText((Join-Path $seed "base.txt"),"base",(New-Object Text.UTF8Encoding($false)))
    $sourceContracts=Join-Path $seed "snapshot\source\research-controller\contracts"
    $sourceTemplates=Join-Path $seed "snapshot\source\research-controller\templates"
    [IO.Directory]::CreateDirectory($sourceContracts)|Out-Null
    [IO.Directory]::CreateDirectory($sourceTemplates)|Out-Null
    [IO.File]::WriteAllText((Join-Path $sourceContracts "selftest.json"),"{}",(New-Object Text.UTF8Encoding($false)))
    [IO.File]::WriteAllText((Join-Path $sourceTemplates "selftest.tmpl"),"selftest",(New-Object Text.UTF8Encoding($false)))
    Invoke-Git $seed @("add","--all")|Out-Null
    Invoke-Git $seed @("-c","user.name=SelfTest","-c","user.email=selftest@local.invalid","commit","-m","base")|Out-Null
    Invoke-Git $seed @("push","origin","HEAD:refs/heads/main")|Out-Null
    $parent=(Invoke-Git $seed @("rev-parse","HEAD")).stdout
    $syncCfg=[pscustomobject]@{contract_source_sync=$true;contract_source_repo=$seed;contract_source_ref="main";authority_floor_sha=$parent;contract_root="";template_root=""}
    $synced=Sync-ContractSource $syncCfg
    Check ($synced -eq $parent) "contract-source-sync-sha"
    Check ((Test-Path -LiteralPath $syncCfg.contract_root) -and (Test-Path -LiteralPath $syncCfg.template_root)) "contract-source-sync-paths"
    [IO.File]::WriteAllText((Join-Path $templates "simple.tmpl"),"value={{VALUE}}",(New-Object Text.UTF8Encoding($false)))
    $renderCfg=[pscustomobject]@{state_root=$stateRoot;template_root=$templates}
    $renderContract=[pscustomobject]@{
      schema="ckb-plane.research-contract.v1"
      experiment_id="SELFTEST-DETERMINISTIC-R1"
      lane="Wingless"
      parent_experiment="SELFTEST-PARENT"
      parent_sha=$parent
      preregistration_sha=$parent
      preregistered_at_utc="2026-10-04T00:00:00Z"
      branch="research/selftest-deterministic-r1"
      hypothesis="selftest"
      experiment_family="selftest"
      manifest_sha="7077d2b72d8954afc24f36b2971393acc11e3b345f7c3521102c20347b9c6171"
      classification_rules=[pscustomobject]@{supported="x";mixed="x";negative="x";invalid="x"}
      successors=[pscustomobject]@{supported=[pscustomobject]@{contract_id=""};mixed=[pscustomobject]@{contract_id=""};negative=[pscustomobject]@{contract_id=""}}
      invalid_behavior="INFRA_INVALID"
      required_authority=[pscustomobject]@{disposition="READY_RESEARCH";manifest_sha256="7077d2b72d8954afc24f36b2971393acc11e3b345f7c3521102c20347b9c6171";production_authority=$false;external_model_calls=$false}
      qualification_profile=[pscustomobject]@{focused_commands=@()}
      implementation_template=[pscustomobject]@{ready=$true;files=@([pscustomobject]@{template="simple.tmpl";target="generated.txt"});variables=[pscustomobject]@{VALUE="alpha"}}
      codex_usage="DISABLED"
      production_authority=$false
      accepted_ref_mutation=$false
      runtime_launch=$false
      capacity_change_authorized=$false
      authority_expansion=$false
      broker_access=$false
      ktrade_access=$false
    }
    Invoke-Git $seed @("push","origin","HEAD:refs/heads/research/selftest-deterministic-r1")|Out-Null
    $preq=Test-SuccessorTemplate $renderCfg $renderContract $seed
    Check ($preq.status -eq "QUALIFIED" -and $preq.base_sha -eq $parent) "template-preauthority-qualification"
    $one=Render-Successor $renderCfg $renderContract $seed
    $two=Render-Successor $renderCfg $renderContract $seed
    Check ($one.head_sha -eq $two.head_sha) "deterministic-commit-identity"
    $remoteHead=(Invoke-Git $seed @("ls-remote","--heads","origin","refs/heads/research/selftest-deterministic-r1")).stdout.Split()[0]
    Check ($remoteHead -eq $one.head_sha) "concurrent-branch-adoption"
  } finally {
    Remove-Item -LiteralPath $root -Recurse -Force -ErrorAction SilentlyContinue
  }

  $source=Get-Content -LiteralPath $PSCommandPath -Raw
  $pushLines=@($source -split [Environment]::NewLine | Where-Object { $_ -match 'Invoke-Git.+@\("push"' })
  $badPush=@($pushLines | Where-Object { $_ -match '--force' })
  Check ($badPush.Count -eq 0) "no-force-push"
  Check ($source -match 'function Repair-WorkflowPackaging') "workflow-packaging-repair-present"
  Check ($source -match 'function Try-AdoptVerifiedWorkflowPackagingDrift') "verified-workflow-drift-adoption-present"
  $idRoot=Join-Path $env:TEMP ("ckb-research-id-selftest-"+[guid]::NewGuid().ToString("N"))
  try{
    [IO.Directory]::CreateDirectory($idRoot)|Out-Null
    $idScript=Join-Path $idRoot "Invoke-ResearchLoop.ps1"
    [IO.File]::WriteAllText($idScript,"# identity selftest",(New-Object Text.UTF8Encoding($false)))
    [IO.File]::WriteAllText((Join-Path $idRoot "controller-qualified-sha.txt"),("a"*40),(New-Object Text.UTF8Encoding($false)))
    $idCfg=[pscustomobject]@{controller_sha="old"}
    $applied=Apply-QualifiedControllerIdentity $idCfg $idScript
    Check ($applied -eq ("a"*40) -and $idCfg.controller_sha -eq ("a"*40)) "qualified-controller-identity-sidecar"
  }finally{
    Remove-Item -LiteralPath $idRoot -Recurse -Force -ErrorAction SilentlyContinue
  }
  Check ($source -match 'Validate frozen package') "workflow-packaging-classifier-present"
  if($failed.Count -gt 0){throw ("SELF_TEST_FAILED:"+($failed -join ","))}
  [pscustomobject]@{schema="ckb-plane.research-loop-self-test.v1";passed=$true;tests=21;codex_used=$false}
}

if($SelfTest){Invoke-SelfTest | ConvertTo-Json -Depth 10;exit 0}

$Config=Read-JsonFile $ConfigPath
[void](Apply-QualifiedControllerIdentity $Config $PSCommandPath)
if($Config.schema -ne "ckb-plane.research-loop-controller-config.v1"){throw "CONFIG_SCHEMA_INVALID"}
if([int]$Config.max_identical_repair_attempts -ne 3){throw "REPAIR_BUDGET_MUST_BE_THREE"}
if($Config.codex_usage -ne "DISABLED"){throw "CODEX_USAGE_MUST_BE_DISABLED"}
if($Config.fifth_manifest_sha256 -ne "7077d2b72d8954afc24f36b2971393acc11e3b345f7c3521102c20347b9c6171"){throw "FIFTH_MANIFEST_IDENTITY_INVALID"}
if($Config.PSObject.Properties["github_token_dpapi_file"] -and ![string]::IsNullOrWhiteSpace([string]$Config.github_token_dpapi_file)){
  $tokenPath=[string]$Config.github_token_dpapi_file
  if(!(Test-Path -LiteralPath $tokenPath)){throw "GITHUB_TOKEN_DPAPI_FILE_MISSING:$tokenPath"}
  $cipher=[IO.File]::ReadAllBytes($tokenPath)
  $plain=[System.Security.Cryptography.ProtectedData]::Unprotect($cipher,$null,[System.Security.Cryptography.DataProtectionScope]::LocalMachine)
  try{
    $token=[Text.Encoding]::UTF8.GetString($plain).Trim()
    if([string]::IsNullOrWhiteSpace($token)){throw "GITHUB_TOKEN_DPAPI_EMPTY"}
    $env:GH_TOKEN=$token
  }finally{
    if($plain){[Array]::Clear($plain,0,$plain.Length)}
  }
}
if($Config.PSObject.Properties["github_cli_config_dir"] -and ![string]::IsNullOrWhiteSpace([string]$Config.github_cli_config_dir)){
  if(!(Test-Path -LiteralPath ([string]$Config.github_cli_config_dir))){throw "GITHUB_CLI_CONFIG_DIR_MISSING:$($Config.github_cli_config_dir)"}
  $env:GH_CONFIG_DIR=[string]$Config.github_cli_config_dir
}

$mutex=Get-ResearchLoopMutex $Config.mutex_name
if($null -eq $mutex){
  [pscustomobject]@{schema="ckb-plane.research-loop-controller.v1";status="ALREADY_RUNNING";observed_at_utc=(Get-UtcNow);codex_used=$false}|ConvertTo-Json -Depth 10
  exit 0
}
try{
  Assert-GitHubAccess $Config
  $contractSourceSha=Sync-ContractSource $Config
  if($contractSourceSha){
    $Config|Add-Member -NotePropertyName "qualification_source_sha" -NotePropertyValue ([string]$contractSourceSha) -Force
  }
  $state=Load-ControllerState $Config.state_root
  $wing=Invoke-LaneSafely $Config $state "Wingless" $Config.lanes.Wingless
  $ygg=Invoke-LaneSafely $Config $state "Yggdrasil" $Config.lanes.Yggdrasil
  $overall="RUNNING"
  foreach($x in @($wing,$ygg)){
    if($x.state -eq "HARD_BLOCKER"){$overall="HARD_BLOCKER";break}
    if($x.state -eq "NOVEL_HYPOTHESIS_REQUIRED" -and $overall -ne "HARD_BLOCKER"){$overall="NOVEL_HYPOTHESIS_REQUIRED"}
    if($x.state -eq "PROGRAM_COMPLETE" -and $overall -eq "RUNNING"){$overall="PROGRAM_COMPLETE"}
  }
  $receipt=[ordered]@{
    schema="ckb-plane.research-loop-controller.v1"
    status=$overall
    controller_sha=$Config.controller_sha
    contract_source_sha=$contractSourceSha
    observed_at_utc=(Get-UtcNow)
    wingless=$wing
    yggdrasil=$ygg
    codex_used=$script:CodexUsed
    production_authority_changed=$script:ProductionAuthorityChanged
    accepted_ref_mutated=$script:AcceptedRefMutated
    runtime_launched=$script:RuntimeLaunched
    dry_run=[bool]$DryRun
  }
  Save-ControllerState $Config.state_root $state
  $sha=Save-ControllerReceipt $Config.state_root ([pscustomobject]$receipt)
  $receipt.receipt_sha256=$sha
  [pscustomobject]$receipt|ConvertTo-Json -Depth 100
  if($script:CodexUsed -or $script:ProductionAuthorityChanged -or $script:AcceptedRefMutated -or $script:RuntimeLaunched){exit 3}
  if($overall -eq "HARD_BLOCKER"){exit 2}
  exit 0
} finally {
  try{$mutex.ReleaseMutex()}catch{}
  $mutex.Dispose()
}
