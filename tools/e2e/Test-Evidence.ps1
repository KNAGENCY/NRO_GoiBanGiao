[CmdletBinding()]
param(
    [string]$EvidenceRoot,
    [switch]$SelfTest
)

$ErrorActionPreference = "Stop"

function Test-E2EEvidenceContract {
    param([Parameter(Mandatory = $true)][string]$Root)
    $resolved = [System.IO.Path]::GetFullPath($Root)
    $resultPath = Join-Path $resolved "result.json"
    if (-not (Test-Path -LiteralPath $resultPath -PathType Leaf)) {
        throw "E2E-EVIDENCE-001: thieu result.json"
    }
    try {
        $result = Get-Content -LiteralPath $resultPath -Raw -Encoding UTF8 | ConvertFrom-Json
    }
    catch {
        throw "E2E-EVIDENCE-001: result.json khong doc duoc"
    }
    if ($null -eq $result -or $result.schema_version -ne 1) {
        throw "E2E-EVIDENCE-001: schema_version phai bang 1"
    }

    $fixture = $result.fixture
    $prefix = [string]$fixture.prefix
    if ($prefix -notmatch '^e2e_[a-z0-9_]{4,15}$') {
        throw "E2E-FIXTURE-002: prefix fixture khong phai ASCII hop le"
    }
    $accounts = @($fixture.accounts)
    if ($accounts.Count -ne 2 -or
        $accounts[0] -cne ($prefix + "a") -or
        $accounts[1] -cne ($prefix + "b")) {
        throw "E2E-FIXTURE-008: evidence khong ghim dung hai account cua run"
    }
    if ($fixture.cleanup_completed -ne $true) {
        throw "E2E-FIXTURE-008: cleanup chua hoan tat"
    }

    $artifacts = @($result.artifacts)
    if ($artifacts.Count -eq 0) {
        throw "E2E-EVIDENCE-002: artifacts khong duoc rong"
    }
    foreach ($artifact in $artifacts) {
        $relative = [string]$artifact.path
        $expected = [string]$artifact.sha256
        if ([string]::IsNullOrWhiteSpace($relative) -or
            [System.IO.Path]::IsPathRooted($relative) -or
            $relative.Split([char[]]@('/', '\')) -contains "..") {
            throw "E2E-SHA-006: artifact path khong hop le"
        }
        if ($expected -notmatch '^[0-9a-fA-F]{64}$') {
            throw "E2E-SHA-006: artifact thieu SHA-256 hop le: $relative"
        }
        $path = Join-Path $resolved $relative
        if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
            throw "E2E-SHA-006: artifact khong ton tai: $relative"
        }
        $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($actual -cne $expected.ToLowerInvariant()) {
            throw "E2E-SHA-006: artifact sai SHA-256: $relative"
        }
    }

    $events = @($result.events)
    if ($events.Count -eq 0) {
        throw "E2E-EVIDENCE-002: events khong duoc rong"
    }
    $previous = [long]0
    foreach ($event in $events) {
        foreach ($field in @(
            "monotonic_ms", "utc", "client", "sequence", "type", "payload")) {
            if ($null -eq $event.PSObject.Properties[$field]) {
                throw "E2E-EVENT-008: event thieu field $field"
            }
        }
        $monotonic = [long]0
        if (-not [long]::TryParse([string]$event.monotonic_ms, [ref]$monotonic) -or
            $monotonic -le 0) {
            throw "E2E-EVENT-008: monotonic_ms phai lon hon 0"
        }
        $utc = [DateTimeOffset]::MinValue
        if (-not [DateTimeOffset]::TryParse([string]$event.utc, [ref]$utc)) {
            throw "E2E-EVENT-008: utc khong parse duoc"
        }
        if ([string]::IsNullOrWhiteSpace([string]$event.client) -or
            [string]::IsNullOrWhiteSpace([string]$event.type) -or
            $null -eq $event.payload) {
            throw "E2E-EVENT-008: client/type/payload khong hop le"
        }
        $current = [long]0
        if (-not [long]::TryParse([string]$event.sequence, [ref]$current)) {
            throw "E2E-EVENT-008: sequence khong phai so nguyen"
        }
        if ($current -le $previous) {
            throw "E2E-SEQUENCE-003: sequence lui hoac trung ($previous -> $current)"
        }
        if ($current -ne ($previous + 1)) {
            throw "E2E-SEQUENCE-004: sequence dut ($previous -> $current)"
        }
        $previous = $current
    }
    return $true
}

function Invoke-EvidenceSelfTest {
    $root = Join-Path ([System.IO.Path]::GetTempPath()) ("nro-e2e-evidence-" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $root | Out-Null
    try {
        $capture = Join-Path $root "capture.png"
        [System.IO.File]::WriteAllBytes($capture, [byte[]](1, 2, 3))
        $sha = (Get-FileHash -LiteralPath $capture -Algorithm SHA256).Hash.ToLowerInvariant()
        $manifest = @{
            schema_version = 1
            fixture = @{
                prefix = "e2e_valid_123456"
                accounts = @("e2e_valid_123456a", "e2e_valid_123456b")
                cleanup_completed = $true
            }
            artifacts = @(@{ path = "capture.png"; sha256 = $sha })
            events = @(
                @{
                    monotonic_ms = 101
                    utc = "2026-07-31T01:02:03.0000000Z"
                    client = "A"
                    sequence = 1
                    type = "control_ready"
                    payload = @{ ready = $true }
                },
                @{
                    monotonic_ms = 102
                    utc = "2026-07-31T01:02:04.0000000Z"
                    client = "A"
                    sequence = 2
                    type = "shutdown_completed"
                    payload = @{ exit_code = 0 }
                }
            )
        }
        $resultPath = Join-Path $root "result.json"

        function Write-TestResult {
            param($Value)
            $Value | ConvertTo-Json -Depth 8 |
                Set-Content -LiteralPath $resultPath -Encoding UTF8
        }
        function Assert-Rejected {
            param([string]$Code, [string]$Message)
            try {
                Test-E2EEvidenceContract -Root $root | Out-Null
                throw "SELFTEST-FAIL: $Message was accepted"
            }
            catch {
                if (-not $_.Exception.Message.Contains($Code)) {
                    throw "SELFTEST-FAIL: expected $Code, got $($_.Exception.Message)"
                }
            }
        }

        $validArtifact = $manifest.artifacts[0]
        $manifest.artifacts = @()
        Write-TestResult $manifest
        Assert-Rejected "E2E-EVIDENCE-002" "empty artifacts"
        $manifest.artifacts = @($validArtifact)

        $validEvents = @($manifest.events)
        $manifest.events = @()
        Write-TestResult $manifest
        Assert-Rejected "E2E-EVIDENCE-002" "empty events"
        $manifest.events = $validEvents

        foreach ($field in @(
            "monotonic_ms", "utc", "client", "sequence", "type", "payload")) {
            $missing = @{}
            foreach ($key in $validEvents[0].Keys) {
                if ($key -cne $field) {
                    $missing[$key] = $validEvents[0][$key]
                }
            }
            $manifest.events = @($missing)
            Write-TestResult $manifest
            Assert-Rejected "E2E-EVENT-008" "event missing $field"
        }
        $manifest.events = $validEvents

        $manifest.artifacts[0].sha256 = $null
        Write-TestResult $manifest
        Assert-Rejected "E2E-SHA-006" "missing SHA"

        $manifest.artifacts[0].sha256 = ("0" * 64)
        Write-TestResult $manifest
        Assert-Rejected "E2E-SHA-006" "wrong SHA"

        $manifest.artifacts[0].sha256 = $sha
        $manifest.events[1].sequence = 3
        Write-TestResult $manifest
        Assert-Rejected "E2E-SEQUENCE-004" "sequence gap"

        $manifest.events[1].sequence = 2
        $manifest.fixture.prefix = "../prod"
        Write-TestResult $manifest
        Assert-Rejected "E2E-FIXTURE-002" "bad prefix"

        $manifest.fixture.prefix = "e2e_valid_123456"
        $manifest.fixture.cleanup_completed = $false
        Write-TestResult $manifest
        Assert-Rejected "E2E-FIXTURE-008" "missing cleanup"

        $manifest.fixture.cleanup_completed = $true
        Write-TestResult $manifest
        if (-not (Test-E2EEvidenceContract -Root $root)) {
            throw "SELFTEST-FAIL: valid evidence was rejected"
        }
    }
    finally {
        $temp = [System.IO.Path]::GetFullPath([System.IO.Path]::GetTempPath())
        $resolved = [System.IO.Path]::GetFullPath($root)
        if ($resolved.StartsWith($temp, [System.StringComparison]::OrdinalIgnoreCase) -and
            (Test-Path -LiteralPath $resolved)) {
            Remove-Item -LiteralPath $resolved -Recurse -Force
        }
    }
}

if ($SelfTest) {
    Invoke-EvidenceSelfTest
    Write-Output "EVIDENCE-SELFTEST-PASS"
    exit 0
}

if ([string]::IsNullOrWhiteSpace($EvidenceRoot)) {
    throw "E2E-EVIDENCE-001: thieu -EvidenceRoot"
}
Test-E2EEvidenceContract -Root $EvidenceRoot
Write-Output "EVIDENCE-VALID"
