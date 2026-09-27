$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm5h-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm5h-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM5H_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM5H_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm5h-residual-collision-trace -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5H_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM5H_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm5h-residual-collision-trace)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM5H_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm5h-residual-collision-trace)|Out-String).Trim();if($P1-cne $P2){throw 'UPLM5H_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm5h-residual-collision-trace.v1' -or [int]$R.residual_condition_cells-ne 24 -or [int]$R.summaries.Count-ne 24 -or [int]$R.points_per_condition_cell-ne 384){throw 'UPLM5H_DESIGN'}
 if($R.coordinate_modified -or $R.adaptive_feature_selection_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UPLM5H_BOUNDARY_LEAK'}
 if([int]$R.collisions.Count-le 0){throw 'UPLM5H_EXPECTED_RESIDUAL_COLLISIONS_MISSING'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 $TotalGroups=0;$TotalMembers=0;$Max=0
 foreach($S in $R.summaries){$TotalGroups+=[int]$S.collision_groups;$TotalMembers+=[int]$S.collision_members;if([int]$S.max_failure_spread-gt $Max){$Max=[int]$S.max_failure_spread}}
 Write-Host "collision_groups=$TotalGroups collision_members=$TotalMembers max_spread=$Max"
 foreach($C in $R.collisions){
   $Desc=@($C.members|ForEach-Object{"b=$($_.budget),start=$($_.action_start),tp=$($_.throughput),out=$($_.outcome)"}) -join ';'
   Write-Host "profile=$($C.deadline_profile) tr=$($C.throughput_reduction) rot=$($C.rotation) perm=$($C.permutation) members=$Desc"
 }
 Write-Host 'WINGLESS_UP327_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
