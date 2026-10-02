param(
    [int]$WaitMinutes = 120,
    [int]$PollSeconds = 20,
    [string[]]$PriorityRunnerRoots = @('C:\actions-runner'),
    [string[]]$WinglessRunnerRoots = @(
        'C:\actions-runner-wingless',
        'C:\actions-runner-up-b',
        'C:\actions-runner-up-c'
    )
)

$ErrorActionPreference = 'Stop'

if ($env:RUNNER_ENVIRONMENT -ceq 'github-hosted') {
    Write-Host 'WINGLESS_HOST_GUARD_PASS: GitHub-hosted ephemeral VM is isolated from production runners.'
    Write-Host 'WINGLESS_HOST_GUARD_MODE: ephemeral-github-hosted'
    if ([string]::IsNullOrWhiteSpace($env:GOMAXPROCS)) { $env:GOMAXPROCS = '4' }
    Write-Host "WINGLESS_HOST_GUARD_GOMAXPROCS: $env:GOMAXPROCS"
    return
}

function Test-PathUnderRoot {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path,

        [Parameter(Mandatory = $true)]
        [string]$Root
    )

    $normalizedRoot = $Root.TrimEnd('\')
    if ($Path.Equals($normalizedRoot, [System.StringComparison]::OrdinalIgnoreCase)) {
        return $true
    }

    return $Path.StartsWith(
        $normalizedRoot + '\',
        [System.StringComparison]::OrdinalIgnoreCase
    )
}

function Get-PriorityRunnerWorkers {
    $workers = @(
        Get-CimInstance Win32_Process |
        Where-Object {
            $_.Name -ieq 'Runner.Worker.exe' -and
            -not [string]::IsNullOrWhiteSpace($_.ExecutablePath)
        }
    )

    return @(
        foreach ($worker in $workers) {
            $path = $worker.ExecutablePath

            $isWingless = $false
            foreach ($root in $WinglessRunnerRoots) {
                if (Test-PathUnderRoot -Path $path -Root $root) {
                    $isWingless = $true
                    break
                }
            }
            if ($isWingless) {
                continue
            }

            foreach ($root in $PriorityRunnerRoots) {
                if (Test-PathUnderRoot -Path $path -Root $root) {
                    $worker
                    break
                }
            }
        }
    )
}

$deadline = (Get-Date).AddMinutes($WaitMinutes)

while ($true) {
    $active = @(Get-PriorityRunnerWorkers)

    if ($active.Count -eq 0) {
        Write-Host 'WINGLESS_HOST_GUARD_PASS: no priority production runner worker is active.'
        break
    }

    $descriptions = @(
        foreach ($worker in $active) {
            "PID=$($worker.ProcessId) PATH=$($worker.ExecutablePath)"
        }
    )

    if ((Get-Date) -ge $deadline) {
        throw "WINGLESS_HOST_GUARD_TIMEOUT: priority production runner remained active. $($descriptions -join '; ')"
    }

    Write-Host "WINGLESS_HOST_GUARD_WAIT: production runner active; Wingless yields. $($descriptions -join '; ')"
    Start-Sleep -Seconds $PollSeconds
}

try {
    [System.Diagnostics.Process]::GetCurrentProcess().PriorityClass = 'BelowNormal'
    Write-Host 'WINGLESS_HOST_GUARD_PRIORITY: BelowNormal'
}
catch {
    Write-Warning "Unable to lower qualification shell priority: $($_.Exception.Message)"
}

if ([string]::IsNullOrWhiteSpace($env:GOMAXPROCS)) {
    $env:GOMAXPROCS = '4'
}
Write-Host "WINGLESS_HOST_GUARD_GOMAXPROCS: $env:GOMAXPROCS"
