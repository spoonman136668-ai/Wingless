$ErrorActionPreference = 'Stop'
$env:GOFLAGS = '-buildvcs=false'
& go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_FULL_REGRESSION_FAILED' }
$raw = & go run ./cmd/up-lm9d-substrate-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join "`n") | ConvertFrom-Json
$metrics = [ordered]@{}
$v = [double]$probe.metrics.PSObject.Properties['mismatch_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['mismatch_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['max_post_target_error'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['max_post_target_error'] = $v
$v = [double]$probe.metrics.PSObject.Properties['max_write_error'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['max_write_error'] = $v
$v = [double]$probe.metrics.PSObject.Properties['paired_comparisons'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['paired_comparisons'] = $v
$v = [double]$probe.metrics.PSObject.Properties['emitted_rows'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['emitted_rows'] = $v
$v = [double]$probe.metrics.PSObject.Properties['unique_states_per_seed_per_substrate'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['unique_states_per_seed_per_substrate'] = $v
$v = [double]$probe.metrics.PSObject.Properties['invalid_permutations'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['invalid_permutations'] = $v
$v = [double]$probe.metrics.PSObject.Properties['packed_high_bit_violations'].Value
if ([double]::IsNaN($v) -or [double]::IsInfinity($v)) { throw 'WINGLESS_METRIC_NONFINITE' }
$metrics['packed_high_bit_violations'] = $v
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'UP-LM9D-SUBSTRATE-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Controller-owned normalized result harness; interpretation follows frozen gates.'
