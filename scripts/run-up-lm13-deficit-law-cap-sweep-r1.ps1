$ErrorActionPreference = 'Stop'
$env:GOFLAGS = '-buildvcs=false'
& go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_FULL_REGRESSION_FAILED' }
$raw = & go run ./cmd/up-lm13-deficit-law-cap-sweep-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join "`n") | ConvertFrom-Json
$metrics = [ordered]@{}
$v = [double]$probe.metrics.PSObject.Properties['substrate_case_runs'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['substrate_case_runs'] = $v
$v = [double]$probe.metrics.PSObject.Properties['cap_level_count'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['cap_level_count'] = $v
$v = [double]$probe.metrics.PSObject.Properties['total_emitted_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['total_emitted_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['aggregate_demand_writes'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['aggregate_demand_writes'] = $v
$v = [double]$probe.metrics.PSObject.Properties['aggregate_accepted_writes'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['aggregate_accepted_writes'] = $v
$v = [double]$probe.metrics.PSObject.Properties['aggregate_deficit'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['aggregate_deficit'] = $v
$v = [double]$probe.metrics.PSObject.Properties['accepted_write_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['accepted_write_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['deficit_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['deficit_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['law_violation_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['law_violation_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['budget_overrun_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['budget_overrun_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['invalid_state_schedules'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['invalid_state_schedules'] = $v
$v = [double]$probe.metrics.PSObject.Properties['row_count_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['row_count_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['state_cardinality_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['state_cardinality_mismatch_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['max_conservation_error'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['max_conservation_error'] = $v
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'UP-LM13-DEFICIT-LAW-CAP-SWEEP-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Controller-owned normalized result harness; interpretation follows frozen gates.'
