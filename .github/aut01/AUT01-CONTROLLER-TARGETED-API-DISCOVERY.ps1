$ErrorActionPreference='Stop'
Set-StrictMode -Version 2.0
if($env:RUNNER_NAME -cne 'WINGLESS-UP-B'){throw "WRONG_RUNNER=$env:RUNNER_NAME"}
$Plane='C:\Users\camar\GolandProjects\ckb-plane'
$Roots=@((Join-Path $Plane 'cmd'),(Join-Path $Plane 'internal\plane'))
$Files=@()
foreach($r in $Roots){if(Test-Path -LiteralPath $r){$Files+=Get-ChildItem -LiteralPath $r -Recurse -File -Filter '*.go' -ErrorAction Stop}}
$Regex='api/native|run-pending|workorders|HandleFunc|ServeMux|/v1/|/api/|local operation|RunLocalOperation|ExternalOrchestrator|controller stop|controller run'
$Hits=@($Files|Select-String -Pattern $Regex -CaseSensitive:$false -ErrorAction SilentlyContinue)
foreach($h in $Hits|Select-Object -First 400){Write-Host "$($h.Path):$($h.LineNumber):$($h.Line.Trim())"}
Write-Host "hit_count=$($Hits.Count)"
Write-Host 'result=AUT01_CONTROLLER_TARGETED_API_DISCOVERY_COMPLETE'
