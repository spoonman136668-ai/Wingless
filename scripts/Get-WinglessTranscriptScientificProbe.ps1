function Get-WinglessTranscriptScientificProbe {
    param(
        [Parameter(Mandatory=$true)][string]$Transcript,
        [Parameter(Mandatory=$true)][string]$ExpectedExperiment
    )
    $Candidates = @(
        foreach ($Line in ($Transcript -split '[\r\n]+')) {
            $Clean = $Line.Trim()
            if ($Clean -notmatch '^\{\s*"schema"\s*:') { continue }
            try { $Candidate = ConvertFrom-Json -InputObject $Clean -ErrorAction Stop }
            catch { continue }
            if ($null -eq $Candidate -or $null -eq $Candidate.PSObject.Properties['schema'] -or
                $null -eq $Candidate.PSObject.Properties['experiment']) { continue }
            if (-not ([string]$Candidate.schema).Equals('wingless.research-scientific-result.v1', [StringComparison]::Ordinal)) { continue }
            if (-not ([string]$Candidate.experiment).Equals($ExpectedExperiment, [StringComparison]::Ordinal)) { continue }
            $Candidate
        }
    )
    if ($Candidates.Count -ne 1) {
        throw "WINGLESS_SCIENTIFIC_JSON_LINE_MISSING_OR_AMBIGUOUS count=$($Candidates.Count)"
    }
    return $Candidates[0]
}
