param(
    [int]$WaitMinutes = 120,
    [int]$PollSeconds = 20,
    [string[]]$PriorityRunnerRoots = @('C:\actions-runner'),
    [string]$WinglessRunnerRoot = 'C:\actions-runner-wingless'
)

$ErrorActionPreference = 'Stop'

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
            if ($path.StartsWith($WinglessRunnerRoot, [System.StringComparison]::OrdinalIgnoreCase)) {
                continue
            }

            foreach ($root in $PriorityRunnerRoots) {
                if ($path.StartsWith($root, [System.StringComparison]::OrdinalIgnoreCase)) {
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
