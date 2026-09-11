param([string]$NvidiaSMI = "")
$ErrorActionPreference = "Stop"
Set-Location (Split-Path $PSScriptRoot -Parent)
New-Item -ItemType Directory -Force evidence | Out-Null
& go version | Tee-Object evidence/go-version.txt
if ($LASTEXITCODE -ne 0) { throw "Go unavailable" }
& go test -buildvcs=false -race ./... -count=1 -timeout 90s -json | Out-File -Encoding utf8 evidence/tests.jsonl
if ($LASTEXITCODE -ne 0) { throw "Native tests failed; see evidence/tests.jsonl" }
& go run -buildvcs=false ./cmd/ice validate
if ($LASTEXITCODE -ne 0) { throw "ICE stale" }
& go run -buildvcs=false ./cmd/ice build
if ($LASTEXITCODE -ne 0) { throw "ICE build failed" }
& go run -buildvcs=false ./cmd/ice validate
if ($LASTEXITCODE -ne 0) { throw "Rebuilt ICE validation failed" }
if ($NvidiaSMI) {
 & go run -buildvcs=false ./cmd/wingless gpu --nvidia-smi $NvidiaSMI | Out-File -Encoding utf8 evidence/gpu.json
 if ($LASTEXITCODE -ne 0) { throw "Requested NVIDIA collection failed" }
}
& go run -buildvcs=false ./cmd/wingless host | Out-File -Encoding utf8 evidence/host.json
if ($LASTEXITCODE -ne 0) { throw "Host collection failed" }
& go build -buildvcs=false ./cmd/wingless
if ($LASTEXITCODE -ne 0) { throw "Native build failed" }
