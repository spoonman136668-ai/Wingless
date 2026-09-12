param()
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

$Repo = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Set-Location $Repo

$OldFlags = $env:GOFLAGS
try {
    $env:GOFLAGS = '-buildvcs=false'
    $Head = (& git rev-parse HEAD).Trim()
    if ($LASTEXITCODE -ne 0 -or $Head -notmatch '^[0-9a-f]{40}$') { throw 'Wingless HEAD unavailable' }

    $Version = (& go version) -join ''
    if ($LASTEXITCODE -ne 0 -or $Version -cne 'go version go1.25.5 windows/amd64') {
        throw "Windows Go 1.25.5 required; got $Version"
    }

    $Main = Get-Content -LiteralPath '.\cmd\wingless\main.go' -Raw
    foreach ($Forbidden in @('NewListener(', 'ListenAndServe(', 'TransportPath', 'integration/ckbplane')) {
        if ($Main.Contains($Forbidden)) { throw "DEFAULT_LISTENER_ACTIVATION_PRESENT: $Forbidden" }
    }

    Write-Host '[WINGLESS-LISTENER] focused inactive listener tests'
    & go test ./integration/ckbplane -run '^TestWinglessListener' -count=1 -v -timeout=180s
    if ($LASTEXITCODE -ne 0) { throw 'Focused Wingless listener tests failed' }

    Write-Host '[WINGLESS-LISTENER] all ckb-plane integration tests'
    & go test ./integration/ckbplane -count=1 -timeout=240s
    if ($LASTEXITCODE -ne 0) { throw 'Wingless ckb-plane integration regression failed' }

    Write-Host '[WINGLESS-LISTENER] full Wingless suite'
    & go test ./... -count=1 -timeout=600s
    if ($LASTEXITCODE -ne 0) { throw 'Full Wingless suite failed' }

    Write-Host '[WINGLESS-LISTENER] build normal Wingless command'
    New-Item -ItemType Directory -Path '.\bin' -Force | Out-Null
    & go build -o .\bin\wingless.exe ./cmd/wingless
    if ($LASTEXITCODE -ne 0) { throw 'Wingless build failed' }

    Write-Host '[WINGLESS-LISTENER] isolated listener request/response replay'
    & go test ./integration/ckbplane -run '^TestWinglessListenerIsolatedProof$' -count=1 -v -timeout=180s
    if ($LASTEXITCODE -ne 0) { throw 'Isolated Wingless listener replay failed' }

    New-Item -ItemType Directory -Path '.\evidence' -Force | Out-Null
    $EvidencePath = Join-Path $Repo 'evidence\stage3-6-inactive-listener.json'
    $Evidence = [ordered]@{
        schema='wingless.stage3_6.inactive-listener.v1'
        status='ISOLATED_LISTENER_REPLAY_PASS_NO_ACTIVATION'
        wingless_commit=$Head
        frozen_plane_source_commit='539d56fec273c8851aea131cf4d31425a62f7250'
        accepted_plane_transport_commit='f3daaf40552a4d4d1bdbd9e4a929e03d61100657'
        observed_at_utc=[DateTimeOffset]::UtcNow.ToString('o')
        go_version=$Version
        focused_listener_tests='PASS'
        existing_ckbplane_tests='PASS'
        full_wingless_suite='PASS'
        wingless_build='PASS'
        isolated_listener_replay='PASS'
        request_body_consumed_before_inference=$true
        inference_parented_to_http_request_context=$true
        cancellation_propagation='PASS'
        broker_backed_mock_replay='PASS'
        candidate_acceptance='external_required'
        work_order_retry_budget=0
        default_listener_start='ABSENT'
        live_listener_launched=$false
        live_plane_contacted=$false
        live_queue_used=$false
        worker_registered=$false
        edits_applied_by_wingless=$false
        accepted_refs_mutated=$false
        ckb_runtime_launched=$false
        coinbase_broker_credentials_production_touched=$false
    }
    $Evidence | ConvertTo-Json -Depth 6 | Set-Content -LiteralPath $EvidencePath -Encoding UTF8

    Write-Host ''
    Write-Host 'WINGLESS STAGE 3.6 INACTIVE LISTENER PROOF PASS'
    Write-Host "evidence=$EvidencePath"
    Write-Host 'default_listener_start=ABSENT'
    Write-Host 'live_listener_launched=False'
    Write-Host 'live_plane_contacted=False'
    Write-Host 'live_queue_used=False'
    Write-Host 'worker_registered=False'
    Write-Host 'candidate_acceptance=external_required'
    Write-Host 'work_order_retry_budget=0'
} finally {
    $env:GOFLAGS = $OldFlags
}
