$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false
 go clean -cache -testcache
 go test ./unitary ./cmd/unitary-up228b-parity-geometry-determinism -count=1;if($LASTEXITCODE-ne 0){throw 'UP228B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP228B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up228b-parity-geometry-determinism)|Out-String).Trim()
 $P2=((go run ./cmd/unitary-up228b-parity-geometry-determinism)|Out-String).Trim()
 if($P1-cne $P2){throw 'UP228B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up228b-parity-geometry-determinism.v1' -or [int]$R.training_phases.Count-ne 15 -or [int]$R.families.Count-ne 4 -or [int]$R.summaries.Count-ne 4){throw 'UP228B_DESIGN'}
 if($R.parity_used_for_modeling -or $R.gate_changed -or $R.evaluation_states_used -or $R.adaptive_feature_selection_used -or $R.new_native_feature_used -or $R.live_activation){throw 'UP228B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "family=$($S.family) even=$($S.distinct_even_geometries) odd=$($S.distinct_odd_geometries) shared=$($S.shared_parity_geometries)"}
 Write-Host 'WINGLESS_UP352_HARNESS_PASS'
}finally{Pop-Location}
