$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path;$Go=(Get-Command go -ErrorAction Stop).Source;$Root=Split-Path (Split-Path $Go -Parent) -Parent
$OldR=$env:GOROOT;$OldC=$env:GOCACHE;$OldT=$env:GOTMPDIR;$C=Join-Path $env:TEMP ("up237b-cache-"+$PID);$T=Join-Path $env:TEMP ("up237b-tmp-"+$PID);New-Item -ItemType Directory -Force $C,$T|Out-Null;$env:GOROOT=$Root;$env:GOCACHE=$C;$env:GOTMPDIR=$T
Push-Location $Repo
try{go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE){throw'UP237B_GO_ENV'};go clean -cache -testcache;if($LASTEXITCODE){throw'UP237B_CLEAN'};go test ./unitary ./cmd/unitary-up237b-far-centered-nearzero-parity -count=1;if($LASTEXITCODE){throw'UP237B_FOCUSED'};go test ./... -count=1;if($LASTEXITCODE){throw'UP237B_FULL'}
$P1=((go run ./cmd/unitary-up237b-far-centered-nearzero-parity)|Out-String).Trim();$P2=((go run ./cmd/unitary-up237b-far-centered-nearzero-parity)|Out-String).Trim();if($P1-cne$P2){throw'UP237B_NONDET'};$R=$P1|ConvertFrom-Json
if($R.schema-cne'wingless.up237b-far-centered-nearzero-parity.v1'-or [int]$R.training_phases.Count-ne15-or [int]$R.far_evaluation_phases.Count-ne32-or [int]$R.summaries.Count-ne4){throw'UP237B_DESIGN'};if($R.parent_classification-cne'CENTERING_REQUIRED'-or [double]$R.parent_centered_accuracy-ne1){throw'UP237B_ANCHOR'};if($R.evaluation_label_fitting_used-or$R.phase_input_at_inference_used-or$R.adaptive_feature_selection_used-or$R.nonlinear_classifier_used-or$R.live_activation){throw'UP237B_BOUNDARY'}
Write-Host $P1;Write-Host "far_accuracy=$($R.far_accuracy) classification=$($R.classification)";Write-Host 'WINGLESS_UP384_HARNESS_PASS'
}finally{Pop-Location;$env:GOROOT=$OldR;$env:GOCACHE=$OldC;$env:GOTMPDIR=$OldT;Remove-Item -Recurse -Force $C,$T -ErrorAction SilentlyContinue}}
