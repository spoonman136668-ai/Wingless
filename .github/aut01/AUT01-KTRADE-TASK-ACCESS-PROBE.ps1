$ErrorActionPreference='Continue'
Set-StrictMode -Version 2.0
if($env:RUNNER_NAME -cne 'WINGLESS-UP-B'){throw "WRONG_RUNNER=$env:RUNNER_NAME"}
$TaskName='CKB-plane Controller'
Write-Host "runner_identity=$([Security.Principal.WindowsIdentity]::GetCurrent().Name)"
Write-Host "runner_name=$env:RUNNER_NAME"
Write-Host '=== SCHTASKS QUERY ==='
& schtasks.exe /Query /TN $TaskName /V /FO LIST
Write-Host "query_exit=$LASTEXITCODE"
Write-Host '=== POWERSHELL TASK QUERY ==='
try{$t=Get-ScheduledTask -TaskName $TaskName -ErrorAction Stop;Write-Host 'get_scheduled_task=PASS';$t|Select-Object TaskName,State,@{n='MultipleInstances';e={$_.Settings.MultipleInstances}}|Format-List}catch{Write-Host "get_scheduled_task=FAIL:$($_.Exception.GetType().Name):$($_.Exception.Message)"}
Write-Host 'result=AUT01_WINGLESS_TASK_ACCESS_PROBE_COMPLETE'
