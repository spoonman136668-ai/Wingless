$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot '..\Get-WinglessTranscriptScientificProbe.ps1')
$Exp='WLM-SI-RAW-REP-PHASE-ANCHOR-BOUNDARY44-RECORD-LOCALITY-R34'
$Good='{"schema":"wingless.research-scientific-result.v1","experiment":"'+$Exp+'","classification":"novel_support_not_necessary","metrics":{"replacement_arm_count":27}}'
function MustReject([string]$Transcript,[string]$Expected) {
    $Rejected=$false
    try { $null = Get-WinglessTranscriptScientificProbe -Transcript $Transcript -ExpectedExperiment $Expected }
    catch { $Rejected=$true }
    if (-not $Rejected) { throw 'WINGLESS_PARSER_FALSE_POSITIVE' }
}
$Parsed=Get-WinglessTranscriptScientificProbe -Transcript ("Control transcript\n"+$Good+"\nOther text") -ExpectedExperiment $Exp
if($Parsed.classification -cne 'novel_support_not_necessary' -or $Parsed.metrics.replacement_arm_count -ne 27) { throw 'WINGLESS_PARSER_REAL_RESULT_MISMATCH' }
MustReject -Transcript ("\n"+$Good+"\n"+$Good) -Expected $Exp
MustReject -Transcript ($Good -replace 'BOUNDARY44','BOUNDARY43') -Expected $Exp
MustReject -Transcript '{"schema":"wingless.research-scientific-result.v1","experiment":"FAKE"}' -Expected $Exp
MustReject -Transcript '{"schema":"wingless.research-scientific-result.v1","experiment":"'+$Exp+'"' -Expected $Exp
MustReject -Transcript ('warning: '+$Good) -Expected $Exp
Write-Host 'WINGLESS_R34_TRANSCRIPT_JSONLINE_VALIDATION=PASS'
Write-Host 'WINGLESS_DUPLICATE_MISMATCH_INCOMPLETE_LOG_NEGATIVES=PASS'
