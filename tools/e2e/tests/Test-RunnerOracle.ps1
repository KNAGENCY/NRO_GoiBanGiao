[CmdletBinding()]
param(
    [string]$FixturePath,
    [switch]$SelfTest
)

$ErrorActionPreference = "Stop"

function Stop-Oracle {
    param([string]$Code, [string]$Message)
    throw "$Code`: $Message"
}

function Resolve-FixtureFile {
    param(
        [Parameter(Mandatory = $true)][string]$Root,
        [Parameter(Mandatory = $true)][string]$Relative
    )
    if ([string]::IsNullOrWhiteSpace($Relative) -or
        [IO.Path]::IsPathRooted($Relative)) {
        Stop-Oracle "E2E-SHA-006" "evidence path khong hop le"
    }
    $rootPath = [IO.Path]::GetFullPath($Root).TrimEnd(
        [IO.Path]::DirectorySeparatorChar,
        [IO.Path]::AltDirectorySeparatorChar)
    $candidate = [IO.Path]::GetFullPath((Join-Path $rootPath $Relative))
    $prefix = $rootPath + [IO.Path]::DirectorySeparatorChar
    if (-not $candidate.StartsWith(
        $prefix,
        [StringComparison]::OrdinalIgnoreCase)) {
        Stop-Oracle "E2E-SHA-006" "evidence vuot khoi fixture root"
    }
    if (-not (Test-Path -LiteralPath $candidate -PathType Leaf) -or
        (Get-Item -LiteralPath $candidate).Length -eq 0) {
        Stop-Oracle "E2E-SHA-006" "evidence thieu hoac rong: $Relative"
    }
    return $candidate
}

function Assert-HashEntries {
    param(
        [Parameter(Mandatory = $true)]$Result,
        [Parameter(Mandatory = $true)][string]$Property,
        [Parameter(Mandatory = $true)][string]$EmptyCode,
        [Parameter(Mandatory = $true)][string]$Root
    )
    $propertyValue = $Result.PSObject.Properties[$Property]
    if ($null -eq $propertyValue) {
        Stop-Oracle $EmptyCode "$Property khong duoc thieu"
    }
    $entries = @($propertyValue.Value)
    if ($entries.Count -eq 0) {
        Stop-Oracle $EmptyCode "$Property khong duoc rong"
    }
    foreach ($entry in $entries) {
        $role = [string]$entry.role
        $relative = [string]$entry.path
        $expected = [string]$entry.sha256
        if ([string]::IsNullOrWhiteSpace($role) -or
            $expected -notmatch '^[0-9a-fA-F]{64}$') {
            Stop-Oracle "E2E-SHA-006" "$Property thieu role/SHA hop le"
        }
        $path = Resolve-FixtureFile -Root $Root -Relative $relative
        $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash
        if ($actual -cne $expected.ToUpperInvariant()) {
            Stop-Oracle "E2E-SHA-006" "SHA khong khop: $relative"
        }
    }
}

function Read-NdjsonEvents {
    param(
        [Parameter(Mandatory = $true)][string]$Root,
        [Parameter(Mandatory = $true)][string]$Relative
    )
    $path = Resolve-FixtureFile -Root $Root -Relative $Relative
    $events = @()
    $previous = [long]0
    foreach ($line in @(Get-Content -LiteralPath $path -Encoding UTF8)) {
        if ([string]::IsNullOrWhiteSpace($line)) {
            continue
        }
        try {
            $event = $line | ConvertFrom-Json
        }
        catch {
            Stop-Oracle "E2E-EVENT-008" "NDJSON khong hop le"
        }
        if ($null -eq $event.PSObject.Properties["sequence"] -or
            $null -eq $event.PSObject.Properties["type"] -or
            $null -eq $event.PSObject.Properties["payload"]) {
            Stop-Oracle "E2E-EVENT-008" "event thieu sequence/type/payload"
        }
        $current = [long]0
        if (-not [long]::TryParse([string]$event.sequence, [ref]$current)) {
            Stop-Oracle "E2E-EVENT-008" "sequence khong phai so nguyen"
        }
        if ($current -le $previous) {
            Stop-Oracle "E2E-SEQUENCE-003" "sequence lui/trung: $previous -> $current"
        }
        if ($current -ne ($previous + 1)) {
            Stop-Oracle "E2E-SEQUENCE-004" "sequence dut: $previous -> $current"
        }
        $previous = $current
        $events += $event
    }
    if ($events.Count -eq 0) {
        Stop-Oracle "E2E-EVIDENCE-002" "NDJSON khong duoc rong"
    }
    return $events
}

function Assert-Identity {
    param(
        [Parameter(Mandatory = $true)]$Fixture,
        [Parameter(Mandatory = $true)][string]$Property,
        [Parameter(Mandatory = $true)][string]$Code
    )
    $values = @(
        [string]$Fixture.aggregate.$Property,
        [string]$Fixture.clean.$Property,
        [string]$Fixture.warm.$Property)
    if (@($values | Where-Object {
        $_ -notmatch '^[0-9a-fA-F]{64}$'
    }).Count -ne 0 -or
        $values[0] -cne $values[1] -or
        $values[0] -cne $values[2]) {
        Stop-Oracle $Code "$Property cua aggregate/clean/warm khong dong nhat"
    }
}

function Assert-RequiredEvidenceRoles {
    param(
        [Parameter(Mandatory = $true)]$Result,
        [Parameter(Mandatory = $true)][string[]]$Roles
    )
    $actual = @($Result.evidence | ForEach-Object { [string]$_.role })
    foreach ($role in $Roles) {
        if (@($actual | Where-Object { $_ -ceq $role }).Count -ne 1) {
            Stop-Oracle "E2E-EVIDENCE-005" "evidence role $role phai co dung mot"
        }
    }
    $eventRecord = @($Result.evidence | Where-Object {
        ([string]$_.role) -ceq "events"
    })[0]
    if ([string]$eventRecord.path -cne [string]$Result.events_path) {
        Stop-Oracle "E2E-EVIDENCE-005" "events_path khong trung evidence"
    }
}

function Assert-CleanDamage {
    param([Parameter(Mandatory = $true)][object[]]$Events)
    $damage = @($Events | Where-Object {
        ([string]$_.type) -ceq "damage_received"
    })
    $expected = @(
        @{ Kind = "mob"; Amount = 5; Before = 150; After = 145 },
        @{ Kind = "player"; Amount = 2; Before = 145; After = 143 },
        @{ Kind = "player"; Amount = 2; Before = 143; After = 141 },
        @{ Kind = "player"; Amount = 2; Before = 141; After = 139 },
        @{ Kind = "player"; Amount = 2; Before = 139; After = 137 }
    )
    if ($damage.Count -ne $expected.Count) {
        Stop-Oracle "E2E-DAMAGE-010" `
            "can dung 5 event damage, thuc te $($damage.Count)"
    }
    for ($index = 0; $index -lt $expected.Count; $index++) {
        $payload = $damage[$index].payload
        $want = $expected[$index]
        if ([string]$payload.source_kind -cne $want.Kind -or
            [int]$payload.amount -ne $want.Amount -or
            [int]$payload.hp_before -ne $want.Before -or
            [int]$payload.hp_after -ne $want.After -or
            ($index -eq 0 -and [int]$payload.source_id -ne 1)) {
            Stop-Oracle "E2E-DAMAGE-010" "damage thu $($index + 1) sai"
        }
    }
}

function Invoke-RunnerOracle {
    param([Parameter(Mandatory = $true)][string]$Path)
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf) -or
        (Get-Item -LiteralPath $Path).Length -eq 0) {
        Stop-Oracle "E2E-FIXTURE-001" "fixture thieu hoac rong"
    }
    try {
        $fixture = Get-Content -LiteralPath $Path -Raw -Encoding UTF8 |
            ConvertFrom-Json
    }
    catch {
        Stop-Oracle "E2E-FIXTURE-001" "fixture JSON khong hop le"
    }
    if ($null -eq $fixture -or $fixture.schema_version -ne 1 -or
        $null -eq $fixture.aggregate -or
        $null -eq $fixture.clean -or
        $null -eq $fixture.warm) {
        Stop-Oracle "E2E-FIXTURE-001" "schema aggregate/clean/warm khong hop le"
    }
    if ([string]$fixture.clean.scope -cne "full") {
        Stop-Oracle "E2E-CLEAN-SCOPE-015" "clean phai co scope=full"
    }
    if ([string]$fixture.warm.scope -cne "warm_smoke") {
        Stop-Oracle "E2E-WARM-SCOPE-014" "warm phai co scope=warm_smoke"
    }
    if ([string]$fixture.aggregate.status -cne "PASS" -or
        [string]$fixture.clean.status -cne "PASS" -or
        [string]$fixture.warm.status -cne "PASS") {
        Stop-Oracle "E2E-CLEAN-WARM-016" `
            "aggregate chi PASS khi clean va warm cung PASS"
    }

    Assert-Identity -Fixture $fixture `
        -Property "executable_sha256" -Code "E2E-EXECUTABLE-012"
    Assert-Identity -Fixture $fixture `
        -Property "server_lock_sha256" -Code "E2E-SERVER-LOCK-017"
    Assert-Identity -Fixture $fixture `
        -Property "asset_lock_sha256" -Code "E2E-ASSET-LOCK-013"

    $root = Split-Path -Parent ([IO.Path]::GetFullPath($Path))
    Assert-HashEntries -Result $fixture.clean -Property "artifacts" `
        -EmptyCode "E2E-ARTIFACT-002" -Root $root
    Assert-HashEntries -Result $fixture.warm -Property "artifacts" `
        -EmptyCode "E2E-ARTIFACT-002" -Root $root
    Assert-HashEntries -Result $fixture.clean -Property "evidence" `
        -EmptyCode "E2E-EVIDENCE-002" -Root $root
    Assert-HashEntries -Result $fixture.warm -Property "evidence" `
        -EmptyCode "E2E-EVIDENCE-002" -Root $root
    Assert-RequiredEvidenceRoles -Result $fixture.clean `
        -Roles @("db_after_shutdown", "db_after_relogin", "events")
    Assert-RequiredEvidenceRoles -Result $fixture.warm `
        -Roles @("cache_after", "events")

    $cleanEvents = @(
        Read-NdjsonEvents -Root $root -Relative ([string]$fixture.clean.events_path))
    $null = @(
        Read-NdjsonEvents -Root $root -Relative ([string]$fixture.warm.events_path))
    Assert-CleanDamage -Events $cleanEvents
    if ([int]$fixture.clean.final_state.hp -ne 137 -or
        [int]$fixture.clean.final_state.hp_max -ne 150) {
        Stop-Oracle "E2E-HP-011" "clean final HP phai bang 137/150"
    }
    return $true
}

function Invoke-OracleSelfTest {
    $fixtureRoot = Join-Path $PSScriptRoot "fixtures"
    $cases = @(
        @{ File = "valid.json"; Code = $null },
        @{ File = "damage-extra.json"; Code = "E2E-DAMAGE-010" },
        @{ File = "hp-final-not-137.json"; Code = "E2E-HP-011" },
        @{ File = "executable-mismatch.json"; Code = "E2E-EXECUTABLE-012" },
        @{ File = "asset-lock-mismatch.json"; Code = "E2E-ASSET-LOCK-013" },
        @{ File = "server-lock-mismatch.json"; Code = "E2E-SERVER-LOCK-017" },
        @{ File = "warm-scope-full.json"; Code = "E2E-WARM-SCOPE-014" },
        @{ File = "db-evidence-missing-sha.json"; Code = "E2E-SHA-006" },
        @{ File = "sequence-backward.json"; Code = "E2E-SEQUENCE-003" },
        @{ File = "sequence-gap.json"; Code = "E2E-SEQUENCE-004" },
        @{ File = "clean-pass-warm-fail.json"; Code = "E2E-CLEAN-WARM-016" },
        @{ File = "artifacts-empty.json"; Code = "E2E-ARTIFACT-002" },
        @{ File = "evidence-empty.json"; Code = "E2E-EVIDENCE-002" },
        @{ File = "evidence-sha-mismatch.json"; Code = "E2E-SHA-006" }
    )
    $failures = 0
    foreach ($case in $cases) {
        $actual = $null
        try {
            $null = Invoke-RunnerOracle -Path (Join-Path $fixtureRoot $case.File)
        }
        catch {
            $actual = $_.Exception.Message.Split(":")[0]
        }
        if ($actual -cne $case.Code) {
            Write-Output (
                "SELFTEST-FAIL {0}: expected={1} actual={2}" -f
                $case.File,
                $(if ($null -eq $case.Code) { "PASS" } else { $case.Code }),
                $(if ($null -eq $actual) { "PASS" } else { $actual }))
            $failures++
        }
        else {
            Write-Output (
                "SELFTEST-PASS {0}: {1}" -f
                $case.File,
                $(if ($null -eq $actual) { "PASS" } else { $actual }))
        }
    }
    if ($failures -ne 0) {
        throw "ORACLE-SELFTEST-FAIL: $failures/$($cases.Count) fixture sai"
    }
    Write-Output "ORACLE-SELFTEST-PASS: $($cases.Count)/$($cases.Count)"
}

if ($SelfTest) {
    Invoke-OracleSelfTest
    exit 0
}

if ([string]::IsNullOrWhiteSpace($FixturePath)) {
    throw "E2E-FIXTURE-001: thieu -FixturePath"
}

$null = Invoke-RunnerOracle -Path $FixturePath
Write-Output "ORACLE-PASS"
