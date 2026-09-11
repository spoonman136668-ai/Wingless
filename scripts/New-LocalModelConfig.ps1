param(
 [Parameter(Mandatory=$true)][string]$Executable,
 [Parameter(Mandatory=$true)][string]$Model,
 [string]$Output = (Join-Path (Split-Path $PSScriptRoot -Parent) 'local-model.json'),
 [ValidateRange(1024,65535)][int]$Port = 18791,
 [ValidateRange(0,128)][int]$GPULayers = 0,
 [string]$NvidiaSMI = '',
 [ValidateRange(0,31)][int]$GPUIndex = 0,
 [UInt64]$MinVRAM = 1073741824,
 [ValidateRange(128,4096)][int]$ContextTokens = 1024,
 [ValidateRange(1,2048)][int]$Batch = 512,
 [ValidateRange(1,2048)][int]$MicroBatch = 128
)
$ErrorActionPreference = 'Stop'
$ExeFile = Get-Item -LiteralPath $Executable
$ModelFile = Get-Item -LiteralPath $Model
if ($ExeFile.Name -notin @('llama-server.exe','llama-server')) { throw 'Choose a trusted llama-server executable.' }
if ($ExeFile.PSIsContainer -or $ModelFile.PSIsContainer -or $ExeFile.Length -gt 512MB -or $ModelFile.Length -gt 2GB) { throw 'File type/size outside experiment limits.' }
if (Test-Path -LiteralPath $Output) { throw 'Output exists; choose a new configuration path.' }
if ($MicroBatch -gt $Batch) { throw 'MicroBatch must not exceed Batch.' }
if ($GPULayers -gt 0 -and (-not $NvidiaSMI -or $MinVRAM -eq 0)) { throw 'GPU configuration requires NVIDIA utility and VRAM floor.' }
if ($NvidiaSMI) { $NvidiaSMI = (Get-Item -LiteralPath $NvidiaSMI).FullName }
$Config = [ordered]@{
 gpu_layers = $GPULayers
 nvidia_smi = $NvidiaSMI
 gpu_index = $GPUIndex
 min_vram = $MinVRAM
 batch = $Batch
 micro_batch = $MicroBatch
 executable = $ExeFile.FullName
 executable_sha256 = (Get-FileHash -LiteralPath $ExeFile.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
 model = $ModelFile.FullName
 model_sha256 = (Get-FileHash -LiteralPath $ModelFile.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
 model_id = 'wingless-local-model'
 port = $Port
 threads = 2
 context_tokens = $ContextTokens
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
