param(
    [Parameter(Mandatory = $true)]
    [string]$LlamaRoot,

    [switch]$Build,

    [ValidateSet('CUDA','CPU')]
    [string]$Backend = 'CUDA'
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$ExpectedCommit = '5266f24da75dc449bd56cbed7addb9c8e4a6a73e'
$SourceDir = Join-Path $PSScriptRoot 'llama-b10809'
$ExampleDir = Join-Path $LlamaRoot 'examples\trace-moe'
$ExamplesCMake = Join-Path $LlamaRoot 'examples\CMakeLists.txt'

if (-not (Test-Path -LiteralPath (Join-Path $LlamaRoot '.git'))) {
    throw "Not a llama.cpp checkout: $LlamaRoot"
}

$Head = (& git -C $LlamaRoot rev-parse HEAD).Trim()
if ($LASTEXITCODE -ne 0 -or $Head -ne $ExpectedCommit) {
    throw "llama.cpp identity mismatch. Expected $ExpectedCommit, got $Head"
}

$Status = @(& git -C $LlamaRoot status --porcelain --untracked-files=all)
if ($LASTEXITCODE -ne 0) {
    throw 'git status failed'
}
if ($Status.Count -ne 0) {
    throw "llama.cpp checkout must be clean before installing the NR-1A research overlay.`n$($Status -join "`n")"
}

if (-not (Test-Path -LiteralPath $ExamplesCMake)) {
    throw "Missing examples CMake file: $ExamplesCMake"
}
foreach ($SourceName in @('trace-moe.cpp','measure-expert-bytes.cpp','CMakeLists.txt')) {
    if (-not (Test-Path -LiteralPath (Join-Path $SourceDir $SourceName) -PathType Leaf)) {
        throw "NR-1A overlay source missing from Wingless research branch: $SourceName"
    }
}

New-Item -ItemType Directory -Force -Path $ExampleDir | Out-Null
Copy-Item -LiteralPath (Join-Path $SourceDir 'trace-moe.cpp') -Destination (Join-Path $ExampleDir 'trace-moe.cpp')
Copy-Item -LiteralPath (Join-Path $SourceDir 'measure-expert-bytes.cpp') -Destination (Join-Path $ExampleDir 'measure-expert-bytes.cpp')
Copy-Item -LiteralPath (Join-Path $SourceDir 'CMakeLists.txt') -Destination (Join-Path $ExampleDir 'CMakeLists.txt')

$CMakeText = Get-Content -LiteralPath $ExamplesCMake -Raw
if ($CMakeText -notmatch '(?m)^\s*add_subdirectory\(trace-moe\)\s*$') {
    $Needle = '    add_subdirectory(eval-callback)'
    if (-not $CMakeText.Contains($Needle)) {
        throw 'Could not locate exact eval-callback insertion point in b10809 examples/CMakeLists.txt'
    }
    $CMakeText = $CMakeText.Replace($Needle, "$Needle`r`n    add_subdirectory(trace-moe)")
    [IO.File]::WriteAllText($ExamplesCMake, $CMakeText, [Text.UTF8Encoding]::new($false))
}

$Diff = @(& git -C $LlamaRoot diff -- examples/CMakeLists.txt examples/trace-moe)
if ($LASTEXITCODE -ne 0) {
    throw 'Could not inspect research overlay diff'
}
if (-not ($Diff -match 'add_subdirectory\(trace-moe\)')) {
    throw 'Research overlay verification failed: CMake registration missing'
}

Write-Host "NR-1A research overlay installed on exact llama.cpp $ExpectedCommit"
Write-Host 'Routing/model graph source was not modified.'

if ($Build) {
    $BuildDir = Join-Path $LlamaRoot 'build-nr1a'
    $Configure = @('-S', $LlamaRoot, '-B', $BuildDir, '-DLLAMA_BUILD_EXAMPLES=ON')
    if ($Backend -eq 'CUDA') {
        $Configure += '-DGGML_CUDA=ON'
    } else {
        $Configure += '-DGGML_CUDA=OFF'
    }

    & cmake @Configure
    if ($LASTEXITCODE -ne 0) { throw 'NR-1A CMake configure failed' }

    & cmake --build $BuildDir --config Release --target llama-trace-moe llama-measure-experts
    if ($LASTEXITCODE -ne 0) { throw 'NR-1A research utility build failed' }

    function Find-NR1AExe([string]$Name) {
        $Candidates = @(
            (Join-Path $BuildDir "bin\Release\$Name.exe"),
            (Join-Path $BuildDir "bin\$Name.exe"),
            (Join-Path $BuildDir "examples\trace-moe\Release\$Name.exe"),
            (Join-Path $BuildDir "examples\trace-moe\$Name.exe")
        )
        return $Candidates | Where-Object { Test-Path -LiteralPath $_ } | Select-Object -First 1
    }

    $Tracer = Find-NR1AExe 'llama-trace-moe'
    $Measure = Find-NR1AExe 'llama-measure-experts'
    if (-not $Tracer) {
        throw 'Build reported success but llama-trace-moe.exe was not found in expected locations'
    }
    if (-not $Measure) {
        throw 'Build reported success but llama-measure-experts.exe was not found in expected locations'
    }

    $TracerHash = (Get-FileHash -LiteralPath $Tracer -Algorithm SHA256).Hash.ToLowerInvariant()
    $MeasureHash = (Get-FileHash -LiteralPath $Measure -Algorithm SHA256).Hash.ToLowerInvariant()
    Write-Host "TRACE_EXE=$Tracer"
    Write-Host "TRACE_EXE_SHA256=$TracerHash"
    Write-Host "EXPERT_MEASURE_EXE=$Measure"
    Write-Host "EXPERT_MEASURE_EXE_SHA256=$MeasureHash"
}
