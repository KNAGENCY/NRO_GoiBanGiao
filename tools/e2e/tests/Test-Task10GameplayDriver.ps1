[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"

$driver = [IO.Path]::GetFullPath(
    (Join-Path $PSScriptRoot "..\Task10-GameplayDriver.ps1"))
if (-not (Test-Path -LiteralPath $driver -PathType Leaf)) {
    throw "E2E-DRIVER-001: thieu Task10 gameplay driver"
}

$cases = @(
    @{ Name = "telemetry_combat_missing"; Code = "E2E-DRIVER-COMBAT-020" },
    @{ Name = "damage_mismatch"; Code = "E2E-DRIVER-COMBAT-020" },
    @{ Name = "hp_mismatch"; Code = "E2E-DRIVER-COMBAT-020" },
    @{ Name = "wrong_client_sequence"; Code = "E2E-SEQUENCE-003" },
    @{ Name = "noop_command"; Code = "E2E-DRIVER-COMMAND-021" },
    @{ Name = "login_failure"; Code = "E2E-DRIVER-LOGIN-022" },
    @{ Name = "shutdown_incomplete"; Code = "E2E-SHUTDOWN-007" },
    @{ Name = "warm_negotiation_missing"; Code = "E2E-WARM-017" },
    @{ Name = "fixture_cleanup_missing"; Code = "E2E-FIXTURE-008" },
    @{ Name = "valid_contract"; Code = "PASS" }
)

$failures = New-Object System.Collections.Generic.List[string]
foreach ($case in $cases) {
    $stdoutPath = Join-Path $env:TEMP (
        "nro-driver-test-" + [Guid]::NewGuid().ToString("N") + ".out")
    $stderrPath = Join-Path $env:TEMP (
        "nro-driver-test-" + [Guid]::NewGuid().ToString("N") + ".err")
    try {
        $process = Start-Process -FilePath "powershell.exe" -Wait -PassThru `
            -NoNewWindow -RedirectStandardOutput $stdoutPath `
            -RedirectStandardError $stderrPath -ArgumentList @(
                "-NoProfile", "-ExecutionPolicy", "Bypass",
                "-File", $driver,
                "-SelfTestCase", $case.Name)
        $output = ((Get-Content -LiteralPath $stdoutPath -Raw -ErrorAction SilentlyContinue) +
            (Get-Content -LiteralPath $stderrPath -Raw -ErrorAction SilentlyContinue))
        if ($process.ExitCode -ne 0 -or
            -not $output.Contains("DRIVER-SELFTEST-$($case.Code) $($case.Name)")) {
            [void]$failures.Add(
                "$($case.Name): expected DRIVER-SELFTEST-$($case.Code), " +
                "exit=$($process.ExitCode), output=$($output.Trim())")
        }
        else {
            Write-Output "DRIVER-CASE-PASS $($case.Name) $($case.Code)"
        }
    }
    finally {
        Remove-Item -LiteralPath $stdoutPath, $stderrPath -Force `
            -ErrorAction SilentlyContinue
    }
}

if ($failures.Count -ne 0) {
    throw "TASK10-DRIVER-SELFTEST-FAIL count=$($failures.Count)`n$($failures -join "`n")"
}

Write-Output "TASK10-DRIVER-SELFTEST-PASS cases=$($cases.Count)"
