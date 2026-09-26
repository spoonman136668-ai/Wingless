$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1v-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1v-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1V_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1V_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1V_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1V_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1v-allfamily-integrated-retention -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1V_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1V_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1v-allfamily-integrated-retention)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1V_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1v-allfamily-integrated-retention)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1V_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1V_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1v-allfamily-integrated-retention.v1'){throw 'UPLM1V_SCHEMA_MISMATCH'}
 if([int]$R.byte_metrics.Count-ne 15 -or [int]$R.byte_summaries.Count-ne 3){throw 'UPLM1V_BYTE_COUNT'}
 if([int]$R.router_metrics.Count-ne 5 -or [int]$R.integrated_metrics.Count-ne 30){throw 'UPLM1V_INTEGRATION_COUNT'}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.adaptation_epochs-ne 4 -or [int]$R.router_epochs-ne 20){throw 'UPLM1V_BUDGET'}
 if([math]::Abs([double]$R.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM1V_MASS'}
 if($R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.router_modified -or $R.adaptive_weighting_used -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM1V_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($A in @('equal_mass','fifth_1p5_mass','fifth_2x_mass')){
  foreach($F in @('base','paraphrase','third','fourth','fifth')){
   $M=@($R.integrated_metrics|Where-Object{$_.arm-ceq $A -and $_.family-ceq $F})
   $Mean=($M|Measure-Object -Property { $_.metric.top1_accuracy } -Average).Average
   Write-Host "arm=$A family=$F integrated_cells=$($M.Count)"
  }
 }
 Write-Host 'WINGLESS_UP154_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
