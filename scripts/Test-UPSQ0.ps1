$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE
$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-upsq0-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-upsq0-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache
$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_GO_ENV_FAILED'}

 go clean -cache -testcache
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_GO_CLEAN_FAILED'}

 go run ./cmd/ice build
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_ICE_BUILD_FAILED'}

 go run ./cmd/ice validate
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_ICE_VALIDATE_FAILED'}

 go test ./unitary ./cmd/unitary-up-sq0-sequence-qualification -count=1
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_FOCUSED_TEST_FAILED'}

 go test ./... -count=1
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_FULL_REGRESSION_FAILED'}

 $P=((go run ./cmd/unitary-up-sq0-sequence-qualification)|Out-String).Trim()
 if($LASTEXITCODE-ne 0){throw 'UPSQ0_PROBE_FAILED'}

 $R=$P|ConvertFrom-Json
 if($R.Schema-cne 'wingless.up-sq0-sequence-qualification.v1'){throw 'UPSQ0_SCHEMA_MISMATCH'}
 if($R.FullGlobalAttention){throw 'UPSQ0_GLOBAL_ATTENTION_PRESENT'}
 if($R.AggregateScoreSelected){throw 'UPSQ0_AGGREGATE_SELECTION_PRESENT'}
 if([int]$R.RecurrentStateBytes-ne 512){throw "UPSQ0_RECURRENT_BYTES $($R.RecurrentStateBytes)"}
 if([int]$R.ExactRecallEntryCap-ne 16){throw "UPSQ0_RECALL_CAP $($R.ExactRecallEntryCap)"}
 if(($R.GeneratorSeedBases -join ',') -cne '120000000,121000000,122000000'){throw "UPSQ0_SEEDS $($R.GeneratorSeedBases -join ',')"}
 if([int]$R.Metrics.Count-ne 138){throw "UPSQ0_METRIC_COUNT $($R.Metrics.Count)"}

 $Arms=@($R.Metrics|ForEach-Object{$_.Arm}|Sort-Object -Unique)
 if(($Arms -join ',') -cne 'transport_gated_correction,transport_gated_correction_bounded_exact_recall,transport_only'){
   throw "UPSQ0_ARMS $($Arms -join ',')"
 }

 $Families=@($R.Metrics|ForEach-Object{$_.Family}|Sort-Object -Unique)
 if(($Families -join ',') -cne 'delayed_copy,multi_bank_interference,multi_query_associative_recall,overwrite_latest_value_wins,role_filler_recombination,selective_copy,state_machine_composition'){
   throw "UPSQ0_FAMILIES $($Families -join ',')"
 }

 foreach($M in $R.Metrics){
   if([int]$M.Events-le 0){throw "UPSQ0_ZERO_EVENTS $($M.Arm) $($M.Family) $($M.Setting)"}
   if([int]$M.RecurrentStateBytes-ne 512){throw "UPSQ0_METRIC_RECURRENT_BYTES"}
   if([int]$M.ExactRecallBytes-gt 256){throw "UPSQ0_RECALL_BYTES_EXCEEDED $($M.ExactRecallBytes)"}
   if([double]$M.PrimaryAccuracy-lt 0 -or [double]$M.PrimaryAccuracy-gt 1){throw "UPSQ0_PRIMARY_RANGE"}
   if([double]$M.SecondaryAccuracy-lt 0 -or [double]$M.SecondaryAccuracy-gt 1){throw "UPSQ0_SECONDARY_RANGE"}
 }

 Write-Host $P
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.Metrics){
   Write-Host "arm=$($M.Arm) family=$($M.Family) setting=$($M.Setting) primary=$($M.PrimaryAccuracy) secondary=$($M.SecondaryAccuracy) tertiary=$($M.TertiaryMetric):$($M.TertiaryValue) recall_bytes=$($M.ExactRecallBytes) events_per_sec=$($M.EventsPerSecond)"
 }
 Write-Host 'WINGLESS_UP90_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
