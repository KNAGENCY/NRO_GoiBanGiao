[CmdletBinding()]
param(
    [string]$EvidenceRoot,
    [string]$OracleRunRoot,
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

function Read-E2EJson {
    param(
        [Parameter(Mandatory = $true)][string]$Path,
        [Parameter(Mandatory = $true)][string]$Code
    )
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) {
        throw "$Code`: thieu file $([IO.Path]::GetFileName($Path))"
    }
    try {
        return Get-Content -LiteralPath $Path -Raw -Encoding UTF8 |
            ConvertFrom-Json
    }
    catch {
        throw "$Code`: JSON khong doc duoc: $([IO.Path]::GetFileName($Path))"
    }
}

function Resolve-E2ERelativePath {
    param(
        [Parameter(Mandatory = $true)][string]$Root,
        [Parameter(Mandatory = $true)][string]$Relative,
        [Parameter(Mandatory = $true)][string]$Code
    )
    if ([string]::IsNullOrWhiteSpace($Relative) -or
        [IO.Path]::IsPathRooted($Relative) -or
        $Relative.Split([char[]]@("/", "\")) -contains "..") {
        throw "$Code`: evidence path khong hop le"
    }
    $base = [IO.Path]::GetFullPath($Root).TrimEnd("\") + "\"
    $path = [IO.Path]::GetFullPath((Join-Path $Root $Relative))
    if (-not $path.StartsWith($base, [StringComparison]::OrdinalIgnoreCase)) {
        throw "$Code`: evidence path vuot khoi run root"
    }
    return $path
}

function Assert-E2EArtifactSha {
    param(
        [Parameter(Mandatory = $true)][string]$Root,
        [Parameter(Mandatory = $true)]$Record,
        [Parameter(Mandatory = $true)][string]$Code
    )
    if ($null -eq $Record) {
        throw "$Code`: thieu evidence record"
    }
    $relative = [string]$Record.path
    $expected = [string]$Record.sha256
    if ($expected -notmatch "^[0-9a-fA-F]{64}$") {
        throw "$Code`: thieu SHA-256 cho $relative"
    }
    $path = Resolve-E2ERelativePath -Root $Root -Relative $relative -Code $Code
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
        throw "$Code`: evidence khong ton tai: $relative"
    }
    $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash
    if ($actual.ToLowerInvariant() -cne $expected.ToLowerInvariant()) {
        throw "$Code`: evidence sai SHA-256: $relative"
    }
    return $path
}

function Read-E2ENdjson {
    param(
        [Parameter(Mandatory = $true)][string]$Root,
        [Parameter(Mandatory = $true)]$Record
    )
    $path = Assert-E2EArtifactSha -Root $Root -Record $Record `
        -Code "E2E-EVIDENCE-SHA-016"
    $events = New-Object System.Collections.Generic.List[object]
    $lastByClient = @{}
    $lineNumber = 0
    foreach ($line in Get-Content -LiteralPath $path -Encoding UTF8) {
        $lineNumber++
        if ([string]::IsNullOrWhiteSpace($line)) {
            continue
        }
        try {
            $event = $line | ConvertFrom-Json
        }
        catch {
            throw "E2E-EVENT-008: NDJSON hong o dong $lineNumber"
        }
        $client = [string]$event.client
        $current = [long]0
        if ([string]::IsNullOrWhiteSpace($client) -or
            -not [long]::TryParse([string]$event.sequence, [ref]$current)) {
            throw "E2E-EVENT-008: client/sequence khong hop le o dong $lineNumber"
        }
        $previous = if ($lastByClient.ContainsKey($client)) {
            [long]$lastByClient[$client]
        }
        else {
            [long]0
        }
        if ($current -le $previous) {
            throw "E2E-SEQUENCE-003: sequence lui hoac trung cho $client ($previous -> $current)"
        }
        if ($current -ne ($previous + 1)) {
            throw "E2E-SEQUENCE-004: sequence dut cho $client ($previous -> $current)"
        }
        $lastByClient[$client] = $current
        [void]$events.Add($event)
    }
    if ($events.Count -eq 0) {
        throw "E2E-EVENT-008: NDJSON khong co event"
    }
    return $events.ToArray()
}

function Assert-E2EVisualReview {
    param([Parameter(Mandatory = $true)]$Result)
    if ([string]$Result.visual_review.status -cne "pass" -or
        [string]::IsNullOrWhiteSpace([string]$Result.visual_review.reviewer)) {
        throw "E2E-VISUAL-018: visual review chua PASS"
    }
    foreach ($name in @(
        "blur", "seam", "edge_alignment", "anchor", "loop", "impact_timing")) {
        $property = $Result.visual_review.checklist.PSObject.Properties[$name]
        if ($null -eq $property -or [string]$property.Value -cne "pass") {
            throw "E2E-VISUAL-018: checklist $name chua PASS"
        }
    }
}

function Assert-E2EResultArtifacts {
    param(
        [Parameter(Mandatory = $true)][string]$Root,
        [Parameter(Mandatory = $true)]$Result
    )
    $artifacts = @($Result.artifacts)
    if ($artifacts.Count -eq 0) {
        throw "E2E-EVIDENCE-SHA-016: khong co artifact"
    }
    $hasScreenshot = $false
    foreach ($artifact in $artifacts) {
        Assert-E2EArtifactSha -Root $Root -Record $artifact `
            -Code "E2E-EVIDENCE-SHA-016" | Out-Null
        if ([string]$artifact.kind -ceq "screenshot") {
            $hasScreenshot = $true
        }
    }
    if (-not $hasScreenshot) {
        throw "E2E-EVIDENCE-SHA-016: thieu screenshot evidence"
    }
}

function Test-E2ECleanOracle {
    param([Parameter(Mandatory = $true)][string]$Root)
    $result = Read-E2EJson -Path (Join-Path $Root "result.json") `
        -Code "E2E-SCOPE-014"
    if ([string]$result.scope -cne "full") {
        throw "E2E-SCOPE-014: clean phai co scope=full"
    }
    if ([string]$result.status -cne "pass") {
        throw "E2E-SCOPE-014: clean full chua PASS"
    }
    foreach ($sha in @(
        [string]$result.executable_sha256,
        [string]$result.server_lock_sha256)) {
        if ($sha -notmatch "^[0-9a-f]{64}$") {
            throw "E2E-ASSET-013: clean thieu executable/server lock SHA"
        }
    }

    $events = @(Read-E2ENdjson -Root $Root -Record $result.events)
    $dbStates = @{}
    foreach ($name in @("before_shutdown", "after_shutdown", "after_relogin")) {
        $record = $result.database.PSObject.Properties[$name]
        if ($null -eq $record) {
            throw "E2E-DB-SHA-015: thieu DB snapshot $name"
        }
        $path = Assert-E2EArtifactSha -Root $Root -Record $record.Value `
            -Code "E2E-DB-SHA-015"
        $dbStates[$name] = Read-E2EJson -Path $path -Code "E2E-DB-SHA-015"
    }
    foreach ($field in @("map_id", "zone_id", "x", "y", "hp", "hp_max")) {
        $expected = [string]$dbStates.before_shutdown.$field
        foreach ($name in @("after_shutdown", "after_relogin")) {
            if ([string]$dbStates[$name].$field -cne $expected) {
                throw "E2E-PERSIST-004: DB $field sai sau $name"
            }
        }
    }

    Assert-E2EResultArtifacts -Root $Root -Result $result
    Assert-E2EVisualReview -Result $result

    $target = [string]$result.combat.target_client
    $damage = @($events | Where-Object {
        [string]$_.client -ceq $target -and
        [string]$_.type -ceq "damage_received"
    })
    $expectedDamage = @(5, 2, 2, 2, 2)
    $before = @(150, 145, 143, 141, 139)
    $after = @(145, 143, 141, 139, 137)
    if ($damage.Count -ne $expectedDamage.Count) {
        throw "E2E-COMBAT-010: damage count actual=$($damage.Count), expected=5"
    }
    for ($index = 0; $index -lt $expectedDamage.Count; $index++) {
        $payload = $damage[$index].payload
        $sourceExpected = if ($index -eq 0) { "mob" } else { "player" }
        if ([int]$payload.amount -ne $expectedDamage[$index] -or
            [int]$payload.hp_before -ne $before[$index] -or
            [int]$payload.hp_after -ne $after[$index] -or
            [string]$payload.source_kind -cne $sourceExpected) {
            throw "E2E-COMBAT-010: damage[$index] khong khop oracle"
        }
    }
    if (@($result.combat.expected_damage) -join "," -cne "5,2,2,2,2") {
        throw "E2E-COMBAT-010: result khong ghim expected damage"
    }
    if ([int]$result.combat.expected_final_hp -ne 137 -or
        [int]$result.combat.actual_final_hp -ne 137 -or
        [int]$damage[-1].payload.hp_after -ne 137 -or
        [int]$result.persistence.hp -ne 137 -or
        [int]$dbStates.after_relogin.hp -ne 137) {
        throw "E2E-HP-011: HP cuoi phai bang 137"
    }
    return $result
}

function Test-E2EWarmOracle {
    param([Parameter(Mandatory = $true)][string]$Root)
    $result = Read-E2EJson -Path (Join-Path $Root "result.json") `
        -Code "E2E-SCOPE-014"
    if ([string]$result.scope -cne "warm_smoke") {
        throw "E2E-SCOPE-014: warm phai co scope=warm_smoke"
    }
    if ([string]$result.status -cne "pass") {
        throw "E2E-SCOPE-014: warm smoke chua PASS"
    }
    [void](Read-E2ENdjson -Root $Root -Record $result.events)
    foreach ($name in @("cache_before", "cache_after")) {
        $record = $result.PSObject.Properties[$name]
        if ($null -eq $record) {
            throw "E2E-EVIDENCE-SHA-016: thieu $name"
        }
        Assert-E2EArtifactSha -Root $Root -Record $record.Value `
            -Code "E2E-EVIDENCE-SHA-016" | Out-Null
    }
    Assert-E2EResultArtifacts -Root $Root -Result $result
    if ($result.resources.negotiation_detected -ne $true -or
        [int]$result.resources.entered_map_id -ne 1 -or
        $result.resources.dragon_ball_fallback -ne $false -or
        [string]$result.resources.applied_lock_sha256 -cne
            [string]$result.server_lock_sha256) {
        throw "E2E-WARM-017: warm upgrade negotiation/resource khong hop le"
    }
    return $result
}

function Test-E2ERunnerOracle {
    param([Parameter(Mandatory = $true)][string]$Root)
    $resolved = [IO.Path]::GetFullPath($Root)
    $clean = Test-E2ECleanOracle -Root (Join-Path $resolved "clean")
    $warm = Test-E2EWarmOracle -Root (Join-Path $resolved "warm-upgrade")
    if ([string]$clean.executable_sha256 -cne
        [string]$warm.executable_sha256) {
        throw "E2E-EXECUTABLE-012: clean/warm khac executable"
    }
    if ([string]$clean.server_lock_sha256 -cne
        [string]$warm.server_lock_sha256) {
        throw "E2E-ASSET-013: clean/warm khac server lock"
    }
    $aggregate = [ordered]@{
        schema_version = 1
        scope = "aggregate"
        status = "pass"
        executable_sha256 = [string]$clean.executable_sha256
        server_lock_sha256 = [string]$clean.server_lock_sha256
        clean = @{ scope = "full"; status = "pass"; path = "clean/result.json" }
        warm = @{
            scope = "warm_smoke"
            status = "pass"
            path = "warm-upgrade/result.json"
        }
    }
    $json = $aggregate | ConvertTo-Json -Depth 8
    [IO.File]::WriteAllText(
        (Join-Path $resolved "aggregate-result.json"),
        $json,
        (New-Object Text.UTF8Encoding($false)))
    return $aggregate
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

if (-not [string]::IsNullOrWhiteSpace($OracleRunRoot)) {
    $aggregate = Test-E2ERunnerOracle -Root $OracleRunRoot
    Write-Output (
        "E2E-RUNNER-ORACLE-PASS " +
        (Join-Path ([IO.Path]::GetFullPath($OracleRunRoot)) "aggregate-result.json"))
    exit 0
}

if ([string]::IsNullOrWhiteSpace($EvidenceRoot)) {
    throw "E2E-EVIDENCE-001: thieu -EvidenceRoot"
}
Test-E2EEvidenceContract -Root $EvidenceRoot
Write-Output "EVIDENCE-VALID"
