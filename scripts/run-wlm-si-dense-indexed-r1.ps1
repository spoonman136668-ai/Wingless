$ErrorActionPreference = 'Stop'
$env:GOFLAGS = '-buildvcs=false'
& go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_FULL_REGRESSION_FAILED' }
$raw = & go run ./cmd/wlm-si-dense-indexed-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join "`n") | ConvertFrom-Json
$metrics = [ordered]@{}
$v = [double]$probe.metrics.PSObject.Properties['budget_overrun_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['budget_overrun_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['completed_substrate_case_runs'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['completed_substrate_case_runs'] = $v
$v = [double]$probe.metrics.PSObject.Properties['invalid_state_schedules'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['invalid_state_schedules'] = $v
$v = [double]$probe.metrics.PSObject.Properties['law_violation_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['law_violation_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['max_conservation_error'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['max_conservation_error'] = $v
$v = [double]$probe.metrics.PSObject.Properties['max_deficit_law_error'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['max_deficit_law_error'] = $v
$v = [double]$probe.metrics.PSObject.Properties['paired_case_count'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['paired_case_count'] = $v
$v = [double]$probe.metrics.PSObject.Properties['paired_observable_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['paired_observable_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['row_count_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['row_count_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['state_cardinality_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['state_cardinality_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['total_emitted_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['total_emitted_rows'] = $v
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-SI-DENSE-INDEXED-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Controller-owned normalized result harness; interpretation follows frozen gates.'
