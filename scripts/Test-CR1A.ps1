param(
    [int]$Repetitions = 100
)

$ErrorActionPreference = 'Stop'
Set-Location (Split-Path $PSScriptRoot -Parent)

if ($Repetitions -lt 2 -or $Repetitions -gt 1000) {
    throw 'Repetitions must be between 2 and 1000.'
}

$Evidence = Join-Path (Get-Location) 'evidence\cr1a'
New-Item -ItemType Directory -Force -Path $Evidence | Out-Null

& go version | Tee-Object (Join-Path $Evidence 'go-version.txt')
if ($LASTEXITCODE -ne 0) { throw 'Go unavailable' }

& go run -buildvcs=false ./cmd/ice validate |
    Tee-Object (Join-Path $Evidence 'ice-validate.txt')
if ($LASTEXITCODE -ne 0) { throw 'ICE stale; run go run -buildvcs=false ./cmd/ice build and review the generated diff first' }

& go test -buildvcs=false -race ./... -count=1 -timeout 90s -json |
    Out-File -Encoding utf8 (Join-Path $Evidence 'tests-full.jsonl')
if ($LASTEXITCODE -ne 0) { throw 'Full Wingless test suite failed' }

& go test -buildvcs=false -race ./cognitive ./benchmark -count=1 -timeout 90s -json |
    Out-File -Encoding utf8 (Join-Path $Evidence 'tests-cognitive.jsonl')
if ($LASTEXITCODE -ne 0) { throw 'CR-1A focused tests failed' }

& go run -buildvcs=false ./cmd/cognitive-bench -repetitions $Repetitions |
    Tee-Object (Join-Path $Evidence 'benchmark-cognitive-reuse.json')
if ($LASTEXITCODE -ne 0) { throw 'CR-1A deterministic reuse benchmark failed' }

& go build -buildvcs=false ./cmd/wingless
if ($LASTEXITCODE -ne 0) { throw 'Wingless build failed' }

Write-Host 'CR-1A candidate validation completed. Acceptance remains external_required.'
