$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up86a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up86a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP86A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP86A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP86A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP86A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up86a-factorized-noise-bridge -count=1;if($LASTEXITCODE-ne 0){throw 'UP86A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP86A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up86a-factorized-noise-bridge)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP86A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up86a-factorized-noise-bridge)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP86A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP86A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up86a-factorized-noise-bridge.v1'){throw 'UP86A_SCHEMA_MISMATCH'}
 if($R.retraining){throw 'UP86A_RETRAINING_PRESENT'}
 if([int]$R.points.Count-ne 10){throw "UP86A_POINT_COUNT $($R.points.Count)"}
 $Noises=@($R.points|ForEach-Object{[double]$_.noise_amplitude}|Sort-Object -Unique)
 if(($Noises -join ',') -cne '0,0.01,0.025,0.05,0.1'){throw "UP86A_NOISE_LEVELS $($Noises -join ',')"}
 foreach($P in $R.points){if([int]$P.training_states-ne 27){throw 'UP86A_TRAIN_COUNT'};if([int]$P.heldout_states-ne 486){throw 'UP86A_HELDOUT_COUNT'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "representation=$($P.representation) noise=$($P.noise_amplitude) held=$($P.mean_heldout_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP86_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
