[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$fixtures = Join-Path $PSScriptRoot "fixtures"
$e2eRoot = Split-Path -Parent $PSScriptRoot
$repositoryRoot = Split-Path -Parent (Split-Path -Parent $e2eRoot)
$runner = Join-Path $e2eRoot "run-unity-e2e.ps1"

function Get-Sha256 {
    param([Parameter(Mandatory = $true)][string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function New-OracleCase {
    param(
        [Parameter(Mandatory = $true)][string]$CleanResult,
        [Parameter(Mandatory = $true)][string]$CleanEvents,
        [string]$WarmResult = "warm-valid-result.json"
    )
    $root = Join-Path ([IO.Path]::GetTempPath()) (
        "nro-task10-oracle-" + [Guid]::NewGuid().ToString("N"))
    $clean = Join-Path $root "clean"
    $warm = Join-Path $root "warm-upgrade"
    foreach ($path in @(
        $clean, $warm,
        (Join-Path $clean "screenshots"),
        (Join-Path $warm "screenshots"))) {
        New-Item -ItemType Directory -Path $path -Force | Out-Null
    }

    Copy-Item -LiteralPath (Join-Path $fixtures $CleanEvents) `
        -Destination (Join-Path $clean "events.ndjson")
    Copy-Item -LiteralPath (Join-Path $fixtures "warm-valid.ndjson") `
        -Destination (Join-Path $warm "events.ndjson")
    foreach ($name in @(
        "db-before.json", "db-after-shutdown.json", "db-after-relogin.json")) {
        Copy-Item -LiteralPath (Join-Path $fixtures "db-state.json") `
            -Destination (Join-Path $clean $name)
    }
    Copy-Item -LiteralPath (Join-Path $fixtures "cache-before.json") `
        -Destination (Join-Path $warm "cache-before.json")
    Copy-Item -LiteralPath (Join-Path $fixtures "cache-after.json") `
        -Destination (Join-Path $warm "cache-after.json")
    Copy-Item -LiteralPath (Join-Path $fixtures "checkpoint.png") `
        -Destination (Join-Path $clean "screenshots\checkpoint.png")
    Copy-Item -LiteralPath (Join-Path $fixtures "checkpoint.png") `
        -Destination (Join-Path $warm "screenshots\checkpoint.png")

    $tokens = @{
        "__DB_SHA__" = Get-Sha256 (Join-Path $clean "db-before.json")
        "__SCREEN_SHA__" = Get-Sha256 (Join-Path $clean "screenshots\checkpoint.png")
        "__CACHE_BEFORE_SHA__" = Get-Sha256 (Join-Path $warm "cache-before.json")
        "__CACHE_AFTER_SHA__" = Get-Sha256 (Join-Path $warm "cache-after.json")
    }
    $cleanText = [IO.File]::ReadAllText(
        (Join-Path $fixtures $CleanResult), [Text.Encoding]::UTF8)
    $warmText = [IO.File]::ReadAllText(
        (Join-Path $fixtures $WarmResult), [Text.Encoding]::UTF8)
    foreach ($key in $tokens.Keys) {
        $cleanText = $cleanText.Replace($key, $tokens[$key])
        $warmText = $warmText.Replace($key, $tokens[$key])
    }
    $cleanText = $cleanText.Replace(
        "__EVENTS_SHA__", (Get-Sha256 (Join-Path $clean "events.ndjson")))
    $warmText = $warmText.Replace(
        "__EVENTS_SHA__", (Get-Sha256 (Join-Path $warm "events.ndjson")))
    [IO.File]::WriteAllText(
        (Join-Path $clean "result.json"), $cleanText,
        (New-Object Text.UTF8Encoding($false)))
    [IO.File]::WriteAllText(
        (Join-Path $warm "result.json"), $warmText,
        (New-Object Text.UTF8Encoding($false)))
    return $root
}

function Invoke-Oracle {
    param([Parameter(Mandatory = $true)][string]$Root)
    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
        $lines = @(
            & powershell -NoProfile -ExecutionPolicy Bypass -File $runner `
                -UnityProject $repositoryRoot `
                -ServerRoot $repositoryRoot `
                -ArtRoot $repositoryRoot `
                -OracleOnly -OracleRunRoot $Root 2>&1 |
                ForEach-Object { [string]$_ })
        $exitCode = $LASTEXITCODE
    }
    finally {
        $ErrorActionPreference = $previousPreference
    }
    return @{
        ExitCode = $exitCode
        Output = ($lines -join "`n")
    }
}

$cases = @(
    @{ Name="valid"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-valid.ndjson"; Expected=$null },
    @{ Name="extra_damage"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-extra-damage.ndjson"; Expected="E2E-COMBAT-010" },
    @{ Name="hp_139"; CleanResult="clean-hp139-result.json";
        CleanEvents="clean-valid.ndjson"; Expected="E2E-HP-011" },
    @{ Name="wrong_executable"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-valid.ndjson"; Warm="warm-wrong-executable-result.json";
        Expected="E2E-EXECUTABLE-012" },
    @{ Name="wrong_lock"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-valid.ndjson"; Warm="warm-wrong-lock-result.json";
        Expected="E2E-ASSET-013" },
    @{ Name="warm_scope_full"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-valid.ndjson"; Warm="warm-wrong-scope-result.json";
        Expected="E2E-SCOPE-014" },
    @{ Name="missing_db_sha"; CleanResult="clean-missing-db-sha-result.json";
        CleanEvents="clean-valid.ndjson"; Expected="E2E-DB-SHA-015" },
    @{ Name="missing_evidence_sha";
        CleanResult="clean-missing-evidence-sha-result.json";
        CleanEvents="clean-valid.ndjson"; Expected="E2E-EVIDENCE-SHA-016" },
    @{ Name="sequence_backward"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-sequence-backward.ndjson"; Expected="E2E-SEQUENCE-003" },
    @{ Name="sequence_gap"; CleanResult="clean-valid-result.json";
        CleanEvents="clean-sequence-gap.ndjson"; Expected="E2E-SEQUENCE-004" }
)

$failures = New-Object System.Collections.Generic.List[string]
foreach ($case in $cases) {
    $caseArgs = @{
        CleanResult = $case.CleanResult
        CleanEvents = $case.CleanEvents
    }
    if (-not [string]::IsNullOrWhiteSpace([string]$case.Warm)) {
        $caseArgs.WarmResult = $case.Warm
    }
    $root = New-OracleCase @caseArgs
    try {
        $actual = Invoke-Oracle -Root $root
        if ($null -eq $case.Expected) {
            if ($actual.ExitCode -ne 0 -or
                -not $actual.Output.Contains("E2E-RUNNER-ORACLE-PASS")) {
                [void]$failures.Add(
                    "$($case.Name): expected PASS, exit=$($actual.ExitCode), output=$($actual.Output)")
            }
            else {
                Write-Output "ORACLE-CASE-PASS $($case.Name)"
            }
        }
        elseif ($actual.ExitCode -eq 0 -or
            -not $actual.Output.Contains([string]$case.Expected)) {
            [void]$failures.Add(
                "$($case.Name): expected $($case.Expected), exit=$($actual.ExitCode), output=$($actual.Output)")
        }
        else {
            Write-Output "ORACLE-CASE-REJECTED $($case.Name) $($case.Expected)"
        }
    }
    finally {
        $temp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
        $resolved = [IO.Path]::GetFullPath($root)
        if ($resolved.StartsWith(
            $temp, [StringComparison]::OrdinalIgnoreCase)) {
            Remove-Item -LiteralPath $resolved -Recurse -Force
        }
    }
}

if ($failures.Count -ne 0) {
    $failures | ForEach-Object { Write-Output "ORACLE-CASE-FAIL $_" }
    throw "RUNNER-ORACLE-TEST-FAIL count=$($failures.Count)"
}
Write-Output "RUNNER-ORACLE-TEST-PASS cases=$($cases.Count)"
