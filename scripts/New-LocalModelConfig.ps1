param(
 [Parameter(Mandatory=$true)][string]$Executable,
 [Parameter(Mandatory=$true)][string]$Model,
 [string]$Output = (Join-Path (Split-Path $PSScriptRoot -Parent) 'local-model.json'),
 [ValidateRange(1024,65535)][int]$Port = 18791
)
$ErrorActionPreference = 'Stop'
$ExeFile = Get-Item -LiteralPath $Executable
$ModelFile = Get-Item -LiteralPath $Model
if ($ExeFile.Name -notin @('llama-server.exe','llama-server')) { throw 'Choose a trusted llama-server executable.' }
if ($ExeFile.PSIsContainer -or $ModelFile.PSIsContainer -or $ExeFile.Length -gt 512MB -or $ModelFile.Length -gt 2GB) { throw 'File type/size outside experiment limits.' }
if (Test-Path -LiteralPath $Output) { throw 'Output exists; choose a new configuration path.' }
$Config = [ordered]@{
 executable = $ExeFile.FullName
 executable_sha256 = (Get-FileHash -LiteralPath $ExeFile.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
 model = $ModelFile.FullName
 model_sha256 = (Get-FileHash -LiteralPath $ModelFile.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
 model_id = 'wingless-local-model'
 port = $Port
 threads = 2
 context_tokens = 1024
 startup_seconds = 60
 runtime_seconds = 180
 min_ram = 1073741824
}
$Json = $Config | ConvertTo-Json
$FullOutput = [IO.Path]::GetFullPath($Output)
$Bytes = (New-Object System.Text.UTF8Encoding($false)).GetBytes($Json)
$Stream = [IO.File]::Open($FullOutput, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write)
try { $Stream.Write($Bytes,0,$Bytes.Length) } finally { $Stream.Dispose() }
Write-Host "Created pinned local experiment configuration: $FullOutput"
