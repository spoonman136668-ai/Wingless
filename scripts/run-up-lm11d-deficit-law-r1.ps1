$ErrorActionPreference = 'Stop'
$env:GOFLAGS = '-buildvcs=false'
& go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_FULL_REGRESSION_FAILED' }
$raw = & go run ./cmd/up-lm11d-deficit-law-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join "`n") | ConvertFrom-Json
$metrics = [ordered]@{}
$v = [double]$probe.metrics.PSObject.Properties['emitted_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['emitted_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['budget_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['budget_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['invalid_permutations'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['invalid_permutations'] = $v
$v = [double]$probe.metrics.PSObject.Properties['unique_states_per_seed'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['unique_states_per_seed'] = $v
$v = [double]$probe.metrics.PSObject.Properties['aggregate_demand_writes'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['aggregate_demand_writes'] = $v
$v = [double]$probe.metrics.PSObject.Properties['aggregate_accepted_writes'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['aggregate_accepted_writes'] = $v
$v = [double]$probe.metrics.PSObject.Properties['aggregate_deficit'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['aggregate_deficit'] = $v
$v = [double]$probe.metrics.PSObject.Properties['law_violation_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['law_violation_cases'] = $v
$v = [double]$probe.metrics.PSObject.Properties['budget_overrun_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['budget_overrun_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['max_conservation_error'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['max_conservation_error'] = $v
$v = [double]$probe.metrics.PSObject.Properties['prefix_digest_comparisons'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['prefix_digest_comparisons'] = $v
$v = [double]$probe.metrics.PSObject.Properties['prefix_digest_mismatch_cases'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['prefix_digest_mismatch_cases'] = $v
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'UP-LM11D-DEFICIT-LAW-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Controller-owned normalized result harness; interpretation follows frozen gates.'
