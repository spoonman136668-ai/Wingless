Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$Url='https://github.com/openai/codex/releases/download/rust-v0.154.0/codex-code-mode-host-x86_64-pc-windows-msvc.exe'
$Expected='7b4987007702973dfeb49ec9a0c11f737488890e208ccb04f7a147769c4bb1f1'
$CodexBins=@(
  'C:\ProgramData\CKBR\codex\bin\codex.exe',
  'C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\bin\codex.exe'
)

$Temp=Join-Path $env:RUNNER_TEMP 'codex-code-mode-host-0.154.0-x64.exe'
Invoke-WebRequest -UseBasicParsing -Uri $Url -OutFile $Temp
$Hash=(Get-FileHash -LiteralPath $Temp -Algorithm SHA256).Hash.ToLowerInvariant()
if($Hash-cne$Expected){throw "HOST_SHA256_MISMATCH expected=$Expected actual=$Hash"}
$Sig=Get-AuthenticodeSignature -LiteralPath $Temp
Write-Host "HOST_SIGNATURE status=$($Sig.Status) subject=$($Sig.SignerCertificate.Subject)"
if([string]$Sig.Status -cne 'Valid'){throw "HOST_SIGNATURE_INVALID status=$($Sig.Status)"}

foreach($Codex in $CodexBins){
  if(-not(Test-Path -LiteralPath $Codex -PathType Leaf)){throw "CODEX_MISSING path=$Codex"}
  $Version=((& $Codex --version 2>$null)|ForEach-Object{[string]$_}) -join ' '
  if($LASTEXITCODE-ne0 -or $Version-notmatch'0\.154\.0'){throw "CODEX_VERSION_MISMATCH path=$Codex version=$Version"}
  $Dir=Split-Path -Parent $Codex
  $Target=Join-Path $Dir 'codex-code-mode-host.exe'
  if(Test-Path -LiteralPath $Target -PathType Leaf){
    $Existing=(Get-FileHash -LiteralPath $Target -Algorithm SHA256).Hash.ToLowerInvariant()
    if($Existing-cne$Expected){throw "EXISTING_HOST_DRIFT path=$Target sha256=$Existing"}
    Write-Host "HOST_ALREADY_CORRECT path=$Target"
  }else{
    Copy-Item -LiteralPath $Temp -Destination $Target -Force
    $Actual=(Get-FileHash -LiteralPath $Target -Algorithm SHA256).Hash.ToLowerInvariant()
    if($Actual-cne$Expected){throw "HOST_COPY_VERIFY_FAILED path=$Target sha256=$Actual"}
    Write-Host "HOST_INSTALLED path=$Target sha256=$Actual"
  }
}

$SmokeRoot=Join-Path $env:RUNNER_TEMP 'codex-host-smoke'
if(Test-Path -LiteralPath $SmokeRoot){Remove-Item -LiteralPath $SmokeRoot -Recurse -Force}
New-Item -ItemType Directory -Force -Path $SmokeRoot|Out-Null
[IO.File]::WriteAllText((Join-Path $SmokeRoot 'smoke.txt'),'CODEX_HOST_SMOKE_OK',(New-Object Text.UTF8Encoding($false)))

foreach($Codex in $CodexBins){
  $Dir=Split-Path -Parent $Codex
  $Host=Join-Path $Dir 'codex-code-mode-host.exe'
  if(-not(Test-Path -LiteralPath $Host -PathType Leaf)){throw "HOST_MISSING_AFTER_INSTALL path=$Host"}
  Write-Host "SMOKE_BEGIN codex=$Codex"
  $OldCodeHome=$env:CODEX_HOME
  $OldHome=$env:HOME
  $OldProfile=$env:USERPROFILE
  try{
    $Home='C:\ProgramData\CKBR\codex\home'
    if($Codex-like'*research-sidecar-yggdrasil*'){
      $Alt='C:\ProgramData\CKBR\research-sidecar-yggdrasil\codex\home'
      if(Test-Path -LiteralPath $Alt -PathType Container){$Home=$Alt}
    }
    if(Test-Path -LiteralPath $CodexHome -PathType Container){
      $env:CODEX_HOME=$CodexHome
      $env:HOME=$CodexHome
      $env:USERPROFILE=$CodexHome
    }
    $Out=Join-Path $SmokeRoot ((Split-Path -Leaf $Dir)+'-last.txt')
    $Prompt='Read smoke.txt using your repository/file tools. Reply with exactly CODEX_HOST_SMOKE_OK and make no changes.'
    $PromptPath=Join-Path $SmokeRoot ((Split-Path -Leaf $Dir)+'-prompt.txt')
    [IO.File]::WriteAllText($PromptPath,$Prompt,(New-Object Text.UTF8Encoding($false)))
    $Proc=Start-Process -FilePath $Codex -ArgumentList @('exec','--approve-for-me','--model','gpt-5.6-luna','--output-last-message',$Out,'-') -WorkingDirectory $SmokeRoot -RedirectStandardInput $PromptPath -NoNewWindow -Wait -PassThru
    if($Proc.ExitCode-ne0){throw "CODEX_SMOKE_EXEC_FAILED path=$Codex exit=$($Proc.ExitCode)"}
    if(-not(Test-Path -LiteralPath $Out -PathType Leaf)){throw "CODEX_SMOKE_OUTPUT_MISSING path=$Codex"}
    $Text=[IO.File]::ReadAllText($Out).Trim()
    if($Text-cne'CODEX_HOST_SMOKE_OK'){throw "CODEX_SMOKE_OUTPUT_INVALID path=$Codex output=$Text"}
    if([IO.File]::ReadAllText((Join-Path $SmokeRoot 'smoke.txt')).Trim()-cne'CODEX_HOST_SMOKE_OK'){throw 'SMOKE_INPUT_MUTATED'}
    Write-Host "SMOKE_PASS codex=$Codex"
  }finally{
    $env:CODEX_HOME=$OldCodeHome
    $env:HOME=$OldHome
    $env:USERPROFILE=$OldProfile
  }
}
Write-Host 'CODEX_CODE_MODE_HOST_REPAIR=PASS'
