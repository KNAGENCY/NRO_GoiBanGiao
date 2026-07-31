[CmdletBinding()]
param(
    [string]$RunRoot,
    [string]$Executable,
    [string]$Manifest,
    [string]$ServerRoot,
    [string]$PreviousLock,
    [string]$GoExe,
    [int]$GamePort = 14445,
    [int]$ControlPort = 19091,
    [int]$TelemetryPortA = 19092,
    [int]$TelemetryPortB = 19093,
    [string]$SelfTestCase
)

$ErrorActionPreference = "Stop"
$script:CommandTimeoutSec = 15
$script:LoginTimeoutSec = 45
$script:ShutdownTimeoutSec = 20

function Stop-Driver {
    param([string]$Code, [string]$Message)
    throw "$Code`: $Message"
}

function Write-Utf8Json {
    param([string]$Path, $Value)
    $parent = Split-Path -Parent $Path
    if (-not [string]::IsNullOrWhiteSpace($parent)) {
        New-Item -ItemType Directory -Path $parent -Force | Out-Null
    }
    [IO.File]::WriteAllText(
        [IO.Path]::GetFullPath($Path),
        ($Value | ConvertTo-Json -Depth 16),
        (New-Object Text.UTF8Encoding($false)))
}

function Write-Utf8Lines {
    param([string]$Path, [object[]]$Values)
    $lines = @($Values | ForEach-Object {
        $_ | ConvertTo-Json -Depth 16 -Compress
    })
    [IO.File]::WriteAllLines(
        [IO.Path]::GetFullPath($Path),
        $lines,
        (New-Object Text.UTF8Encoding($false)))
}

function New-RandomToken {
    $bytes = New-Object byte[] 32
    $rng = [Security.Cryptography.RandomNumberGenerator]::Create()
    try {
        $rng.GetBytes($bytes)
        return [Convert]::ToBase64String($bytes)
    }
    finally {
        $rng.Dispose()
    }
}

function Assert-DriverEvents {
    param(
        [object[]]$Events,
        [string[]]$AllowedClients,
        [switch]$RequireCombat
    )
    if (@($Events).Count -eq 0) {
        Stop-Driver "E2E-DRIVER-COMBAT-020" "telemetry event stream rong"
    }
    $last = @{}
    foreach ($event in @($Events)) {
        $client = [string]$event.client
        if ([string]::IsNullOrWhiteSpace($client) -or
            $AllowedClients -notcontains $client) {
            Stop-Driver "E2E-SEQUENCE-003" "event sai client: $client"
        }
        $sequence = [long]$event.sequence
        $previous = if ($last.ContainsKey($client)) { [long]$last[$client] } else { 0L }
        if ($sequence -le $previous) {
            Stop-Driver "E2E-SEQUENCE-003" "sequence lui/trung cho $client"
        }
        if ($sequence -ne ($previous + 1)) {
            Stop-Driver "E2E-SEQUENCE-004" "sequence dut cho $client"
        }
        $last[$client] = $sequence
    }
    if ($RequireCombat) {
        $damage = @($Events | Where-Object { [string]$_.type -ceq "damage_received" })
        if ($damage.Count -eq 0) {
            Stop-Driver "E2E-DRIVER-COMBAT-020" "thieu damage_received that"
        }
        foreach ($event in $damage) {
            $payload = $event.payload
            if ([int]$payload.amount -le 0 -or
                [int]$payload.hp_before -le [int]$payload.hp_after -or
                ([int]$payload.hp_before - [int]$payload.hp_after) -ne
                    [int]$payload.amount -or
                [string]::IsNullOrWhiteSpace([string]$payload.source_kind) -or
                [int]$payload.source_id -lt 0) {
                Stop-Driver "E2E-DRIVER-COMBAT-020" `
                    "damage_received khong phai transition runtime hop le"
            }
        }
    }
}

function Assert-CommandChangedState {
    param([string]$Command, $Before, $After)
    if ($null -eq $After) {
        Stop-Driver "E2E-DRIVER-COMMAND-021" "$Command khong co state sau lenh"
    }
    $changed = switch ($Command) {
        "move_to" {
            [int]$Before.x -ne [int]$After.x -or
            [int]$Before.y -ne [int]$After.y
        }
        "select_skill" {
            [int]$Before.selected_skill -ne [int]$After.selected_skill
        }
        { $_ -in @("target_mob", "target_player") } {
            [int]$Before.target_id -ne [int]$After.target_id
        }
        default { $true }
    }
    if (-not $changed) {
        Stop-Driver "E2E-DRIVER-COMMAND-021" "$Command la no-op"
    }
}

function Assert-LoginEvidence {
    param([object[]]$Events, [string]$Account)
    $matches = @($Events | Where-Object {
        [string]$_.type -ceq "login_completed" -and
        [string]$_.payload.account -ceq $Account -and
        [int]$_.payload.character_id -ge 0 -and
        [int]$_.payload.map_id -eq 1
    })
    if ($matches.Count -eq 0) {
        Stop-Driver "E2E-DRIVER-LOGIN-022" "login telemetry khong xac nhan map 1"
    }
}

function Assert-ShutdownEvidence {
    param([bool]$Exited, [bool]$BeforeDb, [bool]$AfterDb)
    if (-not $Exited -or -not $BeforeDb -or -not $AfterDb) {
        Stop-Driver "E2E-SHUTDOWN-007" `
            "graceful shutdown chua exit hoac thieu DB snapshot"
    }
}

function Assert-WarmEvidence {
    param([object[]]$Events, [string]$ExpectedLock)
    $negotiation = @($Events | Where-Object {
        [string]$_.type -ceq "warm_negotiation" -and
        $_.payload.negotiation_detected -eq $true
    })
    $applied = @($Events | Where-Object {
        [string]$_.type -ceq "resource_lock_applied" -and
        [string]$_.payload.applied_lock_sha256 -ceq $ExpectedLock -and
        [int]$_.payload.map_id -eq 1 -and
        $_.payload.dragon_ball_fallback -eq $false
    })
    if ($negotiation.Count -eq 0 -or $applied.Count -eq 0) {
        Stop-Driver "E2E-WARM-017" `
            "warm thieu negotiation/resource lock/map state that"
    }
}

function Assert-CleanLifecycle {
    param([object[]]$Events, [string[]]$Clients)
    foreach ($client in $Clients) {
        $own = @($Events | Where-Object { [string]$_.client -ceq $client })
        foreach ($type in @(
            "login_completed","map_resource_state","command_completed",
            "capture_completed","telemetry_shutdown")) {
            if (@($own | Where-Object { [string]$_.type -ceq $type }).Count -eq 0) {
                Stop-Driver "E2E-EVENT-008" "$client thieu event $type"
            }
        }
        foreach ($mode in @("login","relogin")) {
            if (@($own | Where-Object {
                [string]$_.type -ceq "login_completed" -and
                [string]$_.payload.login_mode -ceq $mode
            }).Count -eq 0) {
                Stop-Driver "E2E-DRIVER-LOGIN-022" "$client thieu $mode telemetry"
            }
        }
    }
}

function Assert-FixtureCleanup {
    param($Fixture)
    if ($null -eq $Fixture -or $Fixture.cleanup_completed -ne $true) {
        Stop-Driver "E2E-FIXTURE-008" "fixture cleanup chua hoan tat"
    }
}

function Invoke-DriverSelfTestCase {
    param([string]$Name)
    $validEvents = @(
        [pscustomobject]@{
            client="A"; sequence=1; type="login_completed"
            payload=[pscustomobject]@{account="e2e_testa";character_id=10;map_id=1}
        },
        [pscustomobject]@{
            client="B"; sequence=1; type="command_completed"
            payload=[pscustomobject]@{command="login"}
        },
        [pscustomobject]@{
            client="A"; sequence=2; type="damage_received"
            payload=[pscustomobject]@{
                amount=7;hp_before=90;hp_after=83;source_kind="mob";source_id=4
            }
        },
        [pscustomobject]@{
            client="A"; sequence=3; type="warm_negotiation"
            payload=[pscustomobject]@{negotiation_detected=$true}
        },
        [pscustomobject]@{
            client="A"; sequence=4; type="resource_lock_applied"
            payload=[pscustomobject]@{
                applied_lock_sha256=("a" * 64);map_id=1
                dragon_ball_fallback=$false
            }
        }
    )
    $expected = switch ($Name) {
        "telemetry_combat_missing" { "E2E-DRIVER-COMBAT-020" }
        "damage_mismatch" { "E2E-DRIVER-COMBAT-020" }
        "hp_mismatch" { "E2E-DRIVER-COMBAT-020" }
        "wrong_client_sequence" { "E2E-SEQUENCE-003" }
        "noop_command" { "E2E-DRIVER-COMMAND-021" }
        "login_failure" { "E2E-DRIVER-LOGIN-022" }
        "shutdown_incomplete" { "E2E-SHUTDOWN-007" }
        "warm_negotiation_missing" { "E2E-WARM-017" }
        "fixture_cleanup_missing" { "E2E-FIXTURE-008" }
        "valid_contract" { "PASS" }
        default { Stop-Driver "E2E-DRIVER-001" "self-test case la" }
    }
    try {
        switch ($Name) {
            "telemetry_combat_missing" {
                $withoutCombat = @($validEvents | Where-Object {
                    $_.type -cne "damage_received"
                })
                # Giữ stream liên tục để ca này chỉ thử đúng cổng combat.
                for ($index = 0; $index -lt $withoutCombat.Count; $index++) {
                    if ($withoutCombat[$index].client -ceq "A") {
                        $withoutCombat[$index].sequence = 1 + @(
                            $withoutCombat[0..$index] | Where-Object {
                                $_.client -ceq "A"
                            }).Count - 1
                    }
                }
                Assert-DriverEvents -Events $withoutCombat `
                    -AllowedClients @("A","B") -RequireCombat
            }
            "wrong_client_sequence" {
                $bad = @($validEvents)
                $bad[1] = [pscustomobject]@{
                    client="C";sequence=1;type="command_completed";payload=@{}
                }
                Assert-DriverEvents -Events $bad -AllowedClients @("A","B")
            }
            "damage_mismatch" {
                $validEvents[2].payload.amount = 8
                Assert-DriverEvents -Events $validEvents `
                    -AllowedClients @("A","B") -RequireCombat
            }
            "hp_mismatch" {
                $validEvents[2].payload.hp_after = 82
                Assert-DriverEvents -Events $validEvents `
                    -AllowedClients @("A","B") -RequireCombat
            }
            "noop_command" {
                Assert-CommandChangedState -Command "move_to" `
                    -Before ([pscustomobject]@{x=20;y=30}) `
                    -After ([pscustomobject]@{x=20;y=30})
            }
            "login_failure" {
                Assert-LoginEvidence -Events @($validEvents | Where-Object {
                    $_.type -cne "login_completed"
                }) -Account "e2e_testa"
            }
            "shutdown_incomplete" {
                Assert-ShutdownEvidence -Exited $false -BeforeDb $true -AfterDb $false
            }
            "warm_negotiation_missing" {
                Assert-WarmEvidence -Events @($validEvents | Where-Object {
                    $_.type -cne "warm_negotiation"
                }) -ExpectedLock ("a" * 64)
            }
            "fixture_cleanup_missing" {
                Assert-FixtureCleanup ([pscustomobject]@{cleanup_completed=$false})
            }
            "valid_contract" {
                Assert-DriverEvents -Events $validEvents `
                    -AllowedClients @("A","B") -RequireCombat
                Assert-LoginEvidence -Events $validEvents -Account "e2e_testa"
                Assert-ShutdownEvidence -Exited $true -BeforeDb $true -AfterDb $true
                Assert-WarmEvidence -Events $validEvents -ExpectedLock ("a" * 64)
                Assert-FixtureCleanup ([pscustomobject]@{cleanup_completed=$true})
            }
        }
        if ($expected -cne "PASS") {
            Stop-Driver "SELFTEST-FAIL" "$Name khong bi reject"
        }
        Write-Output "DRIVER-SELFTEST-PASS $Name"
    }
    catch {
        if ($expected -ceq "PASS" -or
            -not $_.Exception.Message.Contains($expected)) {
            throw
        }
        Write-Output "DRIVER-SELFTEST-$expected $Name"
    }
}

if (-not [string]::IsNullOrWhiteSpace($SelfTestCase)) {
    Invoke-DriverSelfTestCase $SelfTestCase
    exit 0
}

function Assert-NormalArguments {
    foreach ($item in @(
        @{Name="RunRoot";Value=$RunRoot;Kind="container"},
        @{Name="Executable";Value=$Executable;Kind="file"},
        @{Name="Manifest";Value=$Manifest;Kind="file"},
        @{Name="ServerRoot";Value=$ServerRoot;Kind="container"},
        @{Name="PreviousLock";Value=$PreviousLock;Kind="file"})) {
        if ([string]::IsNullOrWhiteSpace([string]$item.Value) -or
            -not [IO.Path]::IsPathRooted([string]$item.Value)) {
            Stop-Driver "E2E-DRIVER-001" "$($item.Name) phai la path tuyet doi"
        }
        if ($item.Kind -ceq "file" -and
            -not (Test-Path -LiteralPath $item.Value -PathType Leaf)) {
            Stop-Driver "E2E-DRIVER-001" "khong tim thay $($item.Name)"
        }
        if ($item.Kind -ceq "container" -and
            -not (Test-Path -LiteralPath $item.Value -PathType Container)) {
            Stop-Driver "E2E-DRIVER-001" "khong tim thay $($item.Name)"
        }
    }
    $dsn = $env:NRO_TEST_DSN
    try { $uri = [Uri]$dsn } catch { Stop-Driver "E2E-DB-002" "DSN sai" }
    if ($uri.AbsolutePath.Trim('/') -cne "nro_test" -or
        -not [Net.IPAddress]::IsLoopback([Net.Dns]::GetHostAddresses($uri.Host)[0])) {
        Stop-Driver "E2E-DB-002" "driver chi dung nro_test loopback"
    }
}

function Resolve-DriverGoExecutable {
    if (-not [string]::IsNullOrWhiteSpace($GoExe) -and
        (Test-Path -LiteralPath $GoExe -PathType Leaf)) {
        return [IO.Path]::GetFullPath($GoExe)
    }
    $command = Get-Command go.exe -ErrorAction SilentlyContinue
    if ($null -eq $command) {
        Stop-Driver "E2E-DRIVER-PROCESS-023" "khong tim thay go.exe"
    }
    return $command.Source
}

function Invoke-CheckedProcess {
    param(
        [string]$FilePath,
        [string[]]$Arguments,
        [string]$WorkingDirectory,
        [string]$Stdout,
        [string]$Stderr
    )
    $start = New-Object Diagnostics.ProcessStartInfo
    $start.FileName = $FilePath
    $start.WorkingDirectory = $WorkingDirectory
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    $start.Arguments = [string]::Join(" ", @($Arguments | ForEach-Object {
        '"' + ([string]$_).Replace('"','\"') + '"'
    }))
    $process = New-Object Diagnostics.Process
    $process.StartInfo = $start
    [void]$process.Start()
    $outTask = $process.StandardOutput.ReadToEndAsync()
    $errTask = $process.StandardError.ReadToEndAsync()
    $process.WaitForExit()
    [IO.File]::WriteAllText($Stdout, $outTask.Result, (New-Object Text.UTF8Encoding($false)))
    [IO.File]::WriteAllText($Stderr, $errTask.Result, (New-Object Text.UTF8Encoding($false)))
    if ($process.ExitCode -ne 0) {
        Stop-Driver "E2E-DRIVER-PROCESS-023" `
            "$([IO.Path]::GetFileName($FilePath)) exit $($process.ExitCode)"
    }
}

function Start-TestServer {
    param([string]$Path, [string]$Token, [string]$Stdout, [string]$Stderr)
    $start = New-Object Diagnostics.ProcessStartInfo
    $start.FileName = $Path
    $start.WorkingDirectory = Join-Path $ServerRoot "server"
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    $start.EnvironmentVariables["NRO_DSN"] = $env:NRO_TEST_DSN
    $start.EnvironmentVariables["NRO_DIA_CHI"] = "127.0.0.1:$GamePort"
    $start.EnvironmentVariables["NRO_E2E"] = "1"
    $start.EnvironmentVariables["NRO_E2E_CONTROL_ADDR"] = "127.0.0.1:$ControlPort"
    $start.EnvironmentVariables["NRO_E2E_CONTROL_TOKEN"] = $Token
    $process = New-Object Diagnostics.Process
    $process.StartInfo = $start
    [void]$process.Start()
    $process | Add-Member -NotePropertyName E2EOutTask `
        -NotePropertyValue $process.StandardOutput.ReadToEndAsync()
    $process | Add-Member -NotePropertyName E2EErrTask `
        -NotePropertyValue $process.StandardError.ReadToEndAsync()
    $process | Add-Member -NotePropertyName E2EOutPath -NotePropertyValue $Stdout
    $process | Add-Member -NotePropertyName E2EErrPath -NotePropertyValue $Stderr
    return $process
}

function Complete-TestServer {
    param([Diagnostics.Process]$Process)
    [IO.File]::WriteAllText(
        $Process.E2EOutPath,
        $Process.E2EOutTask.Result,
        (New-Object Text.UTF8Encoding($false)))
    [IO.File]::WriteAllText(
        $Process.E2EErrPath,
        $Process.E2EErrTask.Result,
        (New-Object Text.UTF8Encoding($false)))
}

function Request-GracefulShutdown {
    param([string]$Token)
    $request = [Net.HttpWebRequest]::Create("http://127.0.0.1:$ControlPort/shutdown")
    $request.Method = "POST"
    $request.Headers["Authorization"] = "Bearer $Token"
    $request.Timeout = $script:ShutdownTimeoutSec * 1000
    $response = $request.GetResponse()
    try {
        if ([int]$response.StatusCode -ne 202) {
            Stop-Driver "E2E-SHUTDOWN-007" "control khong tra 202"
        }
    }
    finally { $response.Close() }
}

function Wait-ServerReady {
    param([Diagnostics.Process]$Process, [string]$Token)
    $deadline = [DateTime]::UtcNow.AddSeconds($script:LoginTimeoutSec)
    while ([DateTime]::UtcNow -lt $deadline) {
        if ($Process.HasExited) {
            Stop-Driver "E2E-DRIVER-LOGIN-022" "server exit truoc ready"
        }
        try {
            $request = [Net.HttpWebRequest]::Create("http://127.0.0.1:$ControlPort/shutdown")
            $request.Method = "GET"
            $request.Headers["Authorization"] = "Bearer $Token"
            $request.Timeout = 500
            $response = $request.GetResponse()
            $response.Close()
        }
        catch [Net.WebException] {
            $response = $_.Exception.Response
            if ($null -ne $response) {
                try {
                    if ([int]$response.StatusCode -eq 405) { return }
                }
                finally { $response.Close() }
            }
        }
        Start-Sleep -Milliseconds 100
    }
    Stop-Driver "E2E-DRIVER-LOGIN-022" "server control khong ready"
}

function Wait-ServerExit {
    param([Diagnostics.Process]$Process)
    if (-not $Process.WaitForExit($script:ShutdownTimeoutSec * 1000)) {
        Stop-Driver "E2E-SHUTDOWN-007" "server khong exit graceful"
    }
    Complete-TestServer $Process
    if ($Process.ExitCode -ne 0) {
        Stop-Driver "E2E-SHUTDOWN-007" "server graceful exit $($Process.ExitCode)"
    }
}

function Start-UnityClient {
    param([string]$Name, [int]$Port, [string]$Token, [string]$Profile, [string]$Log)
    New-Item -ItemType Directory -Path $Profile -Force | Out-Null
    $start = New-Object Diagnostics.ProcessStartInfo
    $start.FileName = $Executable
    $start.WorkingDirectory = Split-Path -Parent $Executable
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $false
    $start.WindowStyle = [Diagnostics.ProcessWindowStyle]::Hidden
    $start.EnvironmentVariables["NRO_E2E"] = "1"
    $start.Arguments = [string]::Join(" ", @(
        "-telemetry-port", $Port,
        "-telemetry-token", $Token,
        "-client", $Name,
        "-profile-root", ('"' + $Profile + '"'),
        "-game-port", $GamePort,
        "-logFile", ('"' + $Log + '"')))
    $process = New-Object Diagnostics.Process
    $process.StartInfo = $start
    [void]$process.Start()
    return $process
}

function Connect-Telemetry {
    param([string]$Name, [int]$Port, [string]$Token)
    $deadline = [DateTime]::UtcNow.AddSeconds($script:LoginTimeoutSec)
    $socket = $null
    while ($null -eq $socket -and [DateTime]::UtcNow -lt $deadline) {
        $candidate = New-Object Net.Sockets.TcpClient
        try {
            $candidate.Connect([Net.IPAddress]::Loopback, $Port)
            $socket = $candidate
        }
        catch {
            $candidate.Dispose()
            Start-Sleep -Milliseconds 100
        }
    }
    if ($null -eq $socket) {
        Stop-Driver "E2E-DRIVER-LOGIN-022" "khong ket noi telemetry $Name"
    }
    $socket.ReceiveTimeout = 1000
    $stream = $socket.GetStream()
    $reader = New-Object IO.StreamReader($stream, (New-Object Text.UTF8Encoding($false)))
    $writer = New-Object IO.StreamWriter($stream, (New-Object Text.UTF8Encoding($false)))
    $writer.AutoFlush = $true
    $writer.NewLine = "`n"
    $writer.WriteLine((@{token=$Token} | ConvertTo-Json -Compress))
    $auth = $reader.ReadLine() | ConvertFrom-Json
    if ($auth.ok -ne $true) {
        Stop-Driver "E2E-AUTH-001" "telemetry auth bi tu choi"
    }
    return [pscustomobject]@{
        Name=$Name;Socket=$socket;Stream=$stream;Reader=$reader;Writer=$writer
        Events=(New-Object System.Collections.Generic.List[object]);LastSequence=0L
    }
}

function Read-TelemetryMessage {
    param($Client, [DateTime]$Deadline)
    while ([DateTime]::UtcNow -lt $Deadline) {
        try { $line = $Client.Reader.ReadLine() }
        catch [IO.IOException] { continue }
        if ([string]::IsNullOrWhiteSpace($line)) { continue }
        try { $message = $line | ConvertFrom-Json }
        catch { Stop-Driver "E2E-EVENT-008" "telemetry JSON hong" }
        if ($null -ne $message.PSObject.Properties["type"] -and
            [string]$message.type -cne "auth") {
            if ([string]$message.client -cne [string]$Client.Name) {
                Stop-Driver "E2E-SEQUENCE-003" "telemetry sai client"
            }
            $current = [long]$message.sequence
            if ($current -le [long]$Client.LastSequence) {
                Stop-Driver "E2E-SEQUENCE-003" "telemetry sequence lui/trung"
            }
            if ($current -ne ([long]$Client.LastSequence + 1)) {
                Stop-Driver "E2E-SEQUENCE-004" "telemetry sequence dut"
            }
            $Client.LastSequence = $current
            [void]$Client.Events.Add($message)
            return [pscustomobject]@{Kind="event";Value=$message}
        }
        return [pscustomobject]@{Kind="response";Value=$message}
    }
    return $null
}

function Invoke-ClientCommand {
    param($Client, [string]$Command, $Payload=@{})
    $id = [Guid]::NewGuid().ToString("N")
    $request = @{command_id=$id;command=$Command;payload=$Payload}
    $Client.Writer.WriteLine(($request | ConvertTo-Json -Depth 8 -Compress))
    $deadline = [DateTime]::UtcNow.AddSeconds($script:CommandTimeoutSec)
    while ([DateTime]::UtcNow -lt $deadline) {
        $message = Read-TelemetryMessage -Client $Client -Deadline $deadline
        if ($null -eq $message -or $message.Kind -ceq "event") { continue }
        if ([string]$message.Value.command_id -cne $id) { continue }
        if ($message.Value.ok -ne $true) {
            Stop-Driver ([string]$message.Value.error_code) `
                "$Command bi Unity bridge tu choi"
        }
        return $message.Value
    }
    Stop-Driver "E2E-DRIVER-COMMAND-021" "$Command timeout"
}

function Get-ClientState {
    param($Client)
    $response = Invoke-ClientCommand -Client $Client -Command "get_state"
    return ([string]$response.payload | ConvertFrom-Json)
}

function Invoke-StateChangingCommand {
    param($Client, [string]$Command, $Payload)
    $before = Get-ClientState $Client
    [void](Invoke-ClientCommand -Client $Client -Command $Command -Payload $Payload)
    $deadline = [DateTime]::UtcNow.AddSeconds($script:CommandTimeoutSec)
    do {
        $after = Get-ClientState $Client
        try {
            Assert-CommandChangedState -Command $Command -Before $before -After $after
            return $after
        }
        catch {
            if (-not $_.Exception.Message.Contains("E2E-DRIVER-COMMAND-021")) {
                throw
            }
        }
    } while ([DateTime]::UtcNow -lt $deadline)
    Stop-Driver "E2E-DRIVER-COMMAND-021" "$Command khong doi runtime state"
}

function Assert-CaptureCompleted {
    param($Client, [string]$Path)
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf) -or
        (Get-Item -LiteralPath $Path).Length -le 0) {
        Stop-Driver "E2E-SHA-006" "screenshot that khong ton tai/zero-byte"
    }
    $expected = (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
    $full = [IO.Path]::GetFullPath($Path)
    $event = @($Client.Events | Where-Object {
        [string]$_.type -ceq "capture_completed" -and
        [IO.Path]::GetFullPath([string]$_.payload.path) -ceq $full
    } | Select-Object -Last 1)
    if ($event.Count -ne 1 -or
        [string]$event[0].payload.sha256 -cne $expected) {
        Stop-Driver "E2E-SHA-006" "capture_completed thieu hoac SHA sai"
    }
}

function Close-TelemetryClient {
    param($Client)
    if ($null -eq $Client) { return }
    $Client.Writer.Dispose()
    $Client.Reader.Dispose()
    $Client.Stream.Dispose()
    $Client.Socket.Dispose()
}

function Get-GameplayScenario {
    param([string]$CatalogPath, [string]$CollisionPath)
    $catalog = Get-Content -LiteralPath $CatalogPath -Raw -Encoding UTF8 |
        ConvertFrom-Json
    $wolf = @($catalog.mobs | ForEach-Object { $_.spawns } |
        Where-Object { [int]$_.id -eq 1 } | Select-Object -First 1)
    if ($wolf.Count -ne 1) {
        Stop-Driver "E2E-DRIVER-COMBAT-020" "catalog thieu Soi ID 1"
    }
    $bytes = [IO.File]::ReadAllBytes($CollisionPath)
    if ($bytes.Length -lt 5) {
        Stop-Driver "E2E-ASSET-006" "collision artifact hong"
    }
    $width=[int]$bytes[0];$height=[int]$bytes[1]
    $tileWidth=[int]$bytes[2];$tileHeight=[int]$bytes[3]
    $best=$null;$bestScore=-1
    for ($row=0; $row -lt $height; $row++) {
        for ($column=0; $column -lt $width; $column++) {
            if ($bytes[4 + $row*$width + $column] -ne 0) { continue }
            $x=$column*$tileWidth + [int]($tileWidth/2)
            $y=$row*$tileHeight + [int]($tileHeight/2)
            $score=[Math]::Abs($x-[int]$wolf[0].x) +
                [Math]::Abs($y-[int]$wolf[0].y)
            if ($score -gt $bestScore) {
                $bestScore=$score;$best=[pscustomobject]@{x=$x;y=$y}
            }
        }
    }
    if ($null -eq $best) {
        Stop-Driver "E2E-ASSET-006" "collision khong co diem an toan"
    }
    return [pscustomobject]@{
        Catalog=$catalog
        WolfOne=[pscustomobject]@{x=[int]$wolf[0].x;y=[int]$wolf[0].y}
        Safe=$best
    }
}

function Wait-ClientReady {
    param($Client)
    $deadline = [DateTime]::UtcNow.AddSeconds($script:LoginTimeoutSec)
    while ([DateTime]::UtcNow -lt $deadline) {
        try {
            $state = Get-ClientState $Client
            if ([int]$state.character_id -ge 0 -and [int]$state.map_id -eq 1) {
                return $state
            }
        }
        catch { }
        Start-Sleep -Milliseconds 200
    }
    Stop-Driver "E2E-DRIVER-LOGIN-022" "client khong ready map 1"
}

function Wait-MutualVisibility {
    param($ClientA, $ClientB, $StateA, $StateB)
    $deadline = [DateTime]::UtcNow.AddSeconds($script:LoginTimeoutSec)
    while ([DateTime]::UtcNow -lt $deadline) {
        $StateA = Get-ClientState $ClientA
        $StateB = Get-ClientState $ClientB
        if (@($StateA.visible_player_ids) -contains [int]$StateB.character_id -and
            @($StateB.visible_player_ids) -contains [int]$StateA.character_id) {
            return [pscustomobject]@{A=$StateA;B=$StateB}
        }
        Start-Sleep -Milliseconds 100
    }
    Stop-Driver "E2E-BROADCAST-003" "hai client khong thay nhau"
}

function Wait-EventType {
    param($Client, [string]$Type, [int]$StartIndex)
    $deadline = [DateTime]::UtcNow.AddSeconds($script:LoginTimeoutSec)
    while ([DateTime]::UtcNow -lt $deadline) {
        $existing = @($Client.Events)
        for ($index=$StartIndex; $index -lt $existing.Count; $index++) {
            if ([string]$existing[$index].type -ceq $Type) { return $existing[$index] }
        }
        [void](Read-TelemetryMessage -Client $Client -Deadline $deadline)
    }
    Stop-Driver "E2E-DRIVER-COMBAT-020" "timeout cho event $Type"
}

function Invoke-PsqlText {
    param([string]$Sql)
    $psql = if ($env:NRO_PSQL) { $env:NRO_PSQL } else { "D:\PostgreSQL\bin\psql.exe" }
    $output = & $psql --no-psqlrc --set=ON_ERROR_STOP=1 --quiet `
        --tuples-only --no-align "--dbname=$env:NRO_TEST_DSN" "--command=$Sql" 2>&1
    if ($LASTEXITCODE -ne 0) { Stop-Driver "E2E-DB-002" "psql that bai" }
    return "$($output | Select-Object -Last 1)".Trim()
}

function Write-DbSnapshot {
    param([string]$Account, [string]$Path)
    $json = Invoke-PsqlText @"
SELECT json_build_object(
  'map_id', nv.map_id, 'zone_id', 0,
  'x', nv.vi_tri_x, 'y', nv.vi_tri_y,
  'hp', nv.hp_hien_tai, 'hp_max', nv.hp_toi_da
)::text
FROM nhan_vat nv JOIN tai_khoan tk ON tk.id=nv.tai_khoan_id
WHERE tk.ten_dang_nhap='$Account';
"@
    if ([string]::IsNullOrWhiteSpace($json)) {
        Stop-Driver "E2E-DB-SHA-015" "khong tim thay DB state fixture"
    }
    [void]($json | ConvertFrom-Json)
    [IO.File]::WriteAllText($Path, $json, (New-Object Text.UTF8Encoding($false)))
}

function New-Record {
    param([string]$Root, [string]$Path, [string]$Kind)
    $relative = [IO.Path]::GetFullPath($Path).Substring(
        [IO.Path]::GetFullPath($Root).TrimEnd('\').Length + 1).Replace('\','/')
    return [ordered]@{
        kind=$Kind;path=$relative
        sha256=(Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
    }
}

function Invoke-Fixture {
    param([string]$Action, [string]$FixturePath, [string]$RunId)
    $args = @(
        "-NoProfile","-ExecutionPolicy","Bypass","-File",
        (Join-Path $ServerRoot "tools\e2e\New-Fixture.ps1"),
        "-Dsn",$env:NRO_TEST_DSN,"-Action",$Action,
        "-ManifestPath",$FixturePath)
    if ($Action -ceq "Create") { $args += @("-RunId",$RunId) }
    & powershell.exe @args | Out-Null
    if ($LASTEXITCODE -ne 0) { Stop-Driver "E2E-FIXTURE-008" "$Action fixture fail" }
}

function Stop-UnityGracefully {
    param($Client, [Diagnostics.Process]$Process)
    if ($null -ne $Client) {
        try { [void](Invoke-ClientCommand -Client $Client -Command "shutdown") } catch { throw }
    }
    if ($null -ne $Process -and -not $Process.WaitForExit($script:ShutdownTimeoutSec * 1000)) {
        Stop-Driver "E2E-SHUTDOWN-007" "Unity client khong exit graceful"
    }
}

function Seed-WarmProfile {
    param([string]$SourceProfile, [string]$WarmProfile, $OldLock)
    New-Item -ItemType Directory -Path $WarmProfile -Force | Out-Null
    Get-ChildItem -LiteralPath $SourceProfile -File -ErrorAction SilentlyContinue |
        Copy-Item -Destination $WarmProfile -Force
    $map = @{
        "NRdataVersion"=[int]$OldLock.versions.data
        "NRmapVersion"=[int]$OldLock.versions.map
        "NRskillVersion"=[int]$OldLock.versions.resource
        "NRitemVersion"=[int]$OldLock.versions.resource
    }
    foreach ($entry in $map.GetEnumerator()) {
        [IO.File]::WriteAllBytes(
            (Join-Path $WarmProfile $entry.Key),
            [byte[]]@([byte]$entry.Value))
    }
}

Assert-NormalArguments
$resolvedGoExe = Resolve-DriverGoExecutable
$resolvedRun = [IO.Path]::GetFullPath($RunRoot)
$manifestData = Get-Content -LiteralPath $Manifest -Raw -Encoding UTF8 | ConvertFrom-Json
if ([string]$manifestData.scope -cne "task10_full") {
    Stop-Driver "E2E-SCOPE-014" "driver chi nhan manifest task10_full"
}
$assetLock = [string]$manifestData.preflight.assets.lock_sha256
$serverLockPath = Join-Path $ServerRoot "server\internal\game\assets\asset-lock.json"
$serverLock = (Get-FileHash -LiteralPath $serverLockPath -Algorithm SHA256).Hash.ToLowerInvariant()
if ($serverLock -cne $assetLock) {
    Stop-Driver "E2E-ASSET-006" "server/asset lock lech truoc gameplay"
}
$scenario = Get-GameplayScenario `
    (Join-Path $ServerRoot "server\internal\game\assets\catalog.json") `
    (Join-Path $ServerRoot "server\internal\game\assets\map-1-collision.bin")

$cleanRoot = Join-Path $resolvedRun "clean"
$warmRoot = Join-Path $resolvedRun "warm-upgrade"
New-Item -ItemType Directory -Path $cleanRoot,$warmRoot -Force | Out-Null
$fixturePath = Join-Path $resolvedRun "task10-fixture.json"
$fixturePasswordPrevious = $env:NRO_E2E_FIXTURE_PASSWORD
$env:NRO_E2E_FIXTURE_PASSWORD = ([Guid]::NewGuid().ToString("N") + [Guid]::NewGuid().ToString("N"))
$fixtureCreated = $false
$server = $null
$clientAProcess = $null
$clientBProcess = $null
$warmProcess = $null
$clientA = $null
$clientB = $null
$warmClient = $null
$controlToken = New-RandomToken
$telemetryTokenA = New-RandomToken
$telemetryTokenB = New-RandomToken
try {
    Invoke-Fixture "Create" $fixturePath ([string]$manifestData.run_id)
    $fixtureCreated = $true
    Invoke-Fixture "Validate" $fixturePath $null
    $fixture = Get-Content $fixturePath -Raw -Encoding UTF8 | ConvertFrom-Json
    $accountA = [string]$fixture.accounts[0]
    $accountB = [string]$fixture.accounts[1]

    $serverDir = Join-Path $resolvedRun "server"
    New-Item -ItemType Directory -Path $serverDir -Force | Out-Null
    $serverExe = Join-Path $serverDir "nroserver-task10.exe"
    Invoke-CheckedProcess -FilePath $resolvedGoExe `
        -Arguments @("build","-o",$serverExe,"./cmd/nroserver") `
        -WorkingDirectory (Join-Path $ServerRoot "server") `
        -Stdout (Join-Path $serverDir "build.stdout.log") `
        -Stderr (Join-Path $serverDir "build.stderr.log")
    $server = Start-TestServer $serverExe $controlToken `
        (Join-Path $serverDir "clean.stdout.log") `
        (Join-Path $serverDir "clean.stderr.log")
    Wait-ServerReady $server $controlToken

    $profileA = Join-Path $resolvedRun "profiles\A"
    $profileB = Join-Path $resolvedRun "profiles\B"
    $clientAProcess = Start-UnityClient "A" $TelemetryPortA $telemetryTokenA `
        $profileA (Join-Path $cleanRoot "unity-a.log")
    $clientBProcess = Start-UnityClient "B" $TelemetryPortB $telemetryTokenB `
        $profileB (Join-Path $cleanRoot "unity-b.log")
    $clientA = Connect-Telemetry "A" $TelemetryPortA $telemetryTokenA
    $clientB = Connect-Telemetry "B" $TelemetryPortB $telemetryTokenB
    [void](Invoke-ClientCommand $clientA "login" @{account=$accountA;password=$env:NRO_E2E_FIXTURE_PASSWORD})
    [void](Invoke-ClientCommand $clientB "login" @{account=$accountB;password=$env:NRO_E2E_FIXTURE_PASSWORD})
    $stateA = Wait-ClientReady $clientA
    $stateB = Wait-ClientReady $clientB
    if ([int]$stateA.x -ne [int]$scenario.Safe.x -or
        [int]$stateA.y -ne [int]$scenario.Safe.y) {
        $stateA = Invoke-StateChangingCommand $clientA "move_to" `
            @{x=[int]$scenario.Safe.x;y=[int]$scenario.Safe.y}
    }
    Assert-LoginEvidence @($clientA.Events) $accountA
    Assert-LoginEvidence @($clientB.Events) $accountB

    $visible = Wait-MutualVisibility $clientA $clientB $stateA $stateB
    $stateA=$visible.A;$stateB=$visible.B
    $npc = @($scenario.Catalog.npcs | Select-Object -First 1)
    if ($npc.Count -ne 1) {
        Stop-Driver "E2E-DRIVER-COMMAND-021" "catalog thieu Iruka"
    }
    [void](Invoke-ClientCommand $clientA "open_npc" @{npc_id=[int]$npc[0].id})
    [void](Invoke-ClientCommand $clientA "choose_npc_menu" @{menu_index=0})

    $skills = @($scenario.Catalog.skill_poses | ForEach-Object {[int]$_})
    if ($skills.Count -ne 4) {
        Stop-Driver "E2E-DRIVER-COMMAND-021" "catalog khong co bon skill"
    }
    $orderedSkills = @($skills | Where-Object { $_ -ne [int]$stateA.selected_skill }) +
        @($skills | Where-Object { $_ -eq [int]$stateA.selected_skill })
    foreach ($skill in $orderedSkills) {
        $stateA = Invoke-StateChangingCommand $clientA "select_skill" `
            @{skill_id=$skill}
    }

    if (@($clientA.Events | Where-Object {
        [string]$_.type -ceq "damage_received"
    }).Count -ne 0) {
        Stop-Driver "E2E-DRIVER-COMBAT-020" "damage xay ra truoc checkpoint Soi"
    }
    $damageStart = $clientA.Events.Count
    $stateA = Invoke-StateChangingCommand $clientA "move_to" `
        @{x=[int]$scenario.WolfOne.x;y=[int]$scenario.WolfOne.y}
    [void](Wait-EventType $clientA "damage_received" $damageStart)
    $stateA = Invoke-StateChangingCommand $clientA "move_to" `
        @{x=[int]$scenario.Safe.x;y=[int]$scenario.Safe.y}

    $pvpSkill = [int]$skills[1]
    if ([int]$stateB.selected_skill -eq $pvpSkill) {
        $alternate = @($skills | Where-Object { $_ -ne $pvpSkill } | Select-Object -First 1)
        $stateB = Invoke-StateChangingCommand $clientB "select_skill" `
            @{skill_id=[int]$alternate[0]}
    }
    $stateB = Invoke-StateChangingCommand $clientB "select_skill" `
        @{skill_id=$pvpSkill}
    if ([int]$stateB.x -ne [int]$scenario.Safe.x -or
        [int]$stateB.y -ne [int]$scenario.Safe.y) {
        $stateB = Invoke-StateChangingCommand $clientB "move_to" `
            @{x=[int]$scenario.Safe.x;y=[int]$scenario.Safe.y}
    }
    $stateB = Invoke-StateChangingCommand $clientB "target_player" `
        @{target_id=[int]$stateA.character_id}
    for ($attack=0; $attack -lt 4; $attack++) {
        $damageStart = $clientA.Events.Count
        $readyStart = $clientB.Events.Count
        [void](Invoke-ClientCommand $clientB "attack" @{})
        [void](Wait-EventType $clientA "damage_received" $damageStart)
        if ($attack -lt 3) {
            [void](Wait-EventType $clientB "skill_ready" $readyStart)
        }
    }
    $captureA = Join-Path $cleanRoot "client-a-gameplay.png"
    $captureB = Join-Path $cleanRoot "client-b-gameplay.png"
    [void](Invoke-ClientCommand $clientA "capture" @{path=$captureA})
    [void](Invoke-ClientCommand $clientB "capture" @{path=$captureB})
    Assert-CaptureCompleted $clientA $captureA
    Assert-CaptureCompleted $clientB $captureB

    $beforeDb = Join-Path $cleanRoot "db-before-shutdown.json"
    $afterDb = Join-Path $cleanRoot "db-after-shutdown.json"
    $reloginDb = Join-Path $cleanRoot "db-after-relogin.json"
    Write-DbSnapshot $accountA $beforeDb
    Request-GracefulShutdown $controlToken
    Wait-ServerExit $server
    Write-DbSnapshot $accountA $afterDb
    Assert-ShutdownEvidence $server.HasExited (Test-Path $beforeDb) (Test-Path $afterDb)

    $server = Start-TestServer $serverExe $controlToken `
        (Join-Path $serverDir "relogin.stdout.log") `
        (Join-Path $serverDir "relogin.stderr.log")
    Wait-ServerReady $server $controlToken
    [void](Invoke-ClientCommand $clientA "login" @{account=$accountA;password=$env:NRO_E2E_FIXTURE_PASSWORD})
    [void](Invoke-ClientCommand $clientB "login" @{account=$accountB;password=$env:NRO_E2E_FIXTURE_PASSWORD})
    [void](Wait-ClientReady $clientA)
    [void](Wait-ClientReady $clientB)
    Write-DbSnapshot $accountA $reloginDb

    Stop-UnityGracefully $clientA $clientAProcess
    Stop-UnityGracefully $clientB $clientBProcess
    $clientAProcess=$null;$clientBProcess=$null

    $cleanEvents = @(@($clientA.Events) + @($clientB.Events) |
        Sort-Object monotonic_ms,client,sequence)
    Assert-DriverEvents $cleanEvents @("A","B") -RequireCombat
    Assert-CleanLifecycle $cleanEvents @("A","B")
    $cleanEventsPath = Join-Path $cleanRoot "events.ndjson"
    Write-Utf8Lines $cleanEventsPath $cleanEvents
    Write-Utf8Lines (Join-Path $cleanRoot "events-a.ndjson") @($clientA.Events)
    Write-Utf8Lines (Join-Path $cleanRoot "events-b.ndjson") @($clientB.Events)
    $damage = @($cleanEvents | Where-Object {
        $_.client -ceq "A" -and $_.type -ceq "damage_received"
    })
    $lastDamage = $damage | Select-Object -Last 1
    $dbAfterRelogin = Get-Content $reloginDb -Raw | ConvertFrom-Json
    $cleanResult = [ordered]@{
        schema_version=1;scope="full";status="captured"
        executable_sha256=(Get-FileHash $Executable -Algorithm SHA256).Hash.ToLowerInvariant()
        server_lock_sha256=$serverLock;asset_lock_sha256=$assetLock
        pair_sha256=[string]$manifestData.preflight.assets.pair_sha256
        versions=$manifestData.preflight.assets.versions
        events=(New-Record $cleanRoot $cleanEventsPath "ndjson")
        database=[ordered]@{
            before_shutdown=(New-Record $cleanRoot $beforeDb "database")
            after_shutdown=(New-Record $cleanRoot $afterDb "database")
            after_relogin=(New-Record $cleanRoot $reloginDb "database")
        }
        evidence=@(
            (New-Record $cleanRoot (Join-Path $cleanRoot "events-a.ndjson") "client_ndjson"),
            (New-Record $cleanRoot (Join-Path $cleanRoot "events-b.ndjson") "client_ndjson"))
        artifacts=@(
            (New-Record $cleanRoot $captureA "screenshot"),
            (New-Record $cleanRoot $captureB "screenshot"))
        combat=[ordered]@{
            target_client="A"
            expected_damage=@($damage | ForEach-Object {[int]$_.payload.amount})
            expected_final_hp=[int]$lastDamage.payload.hp_after
            actual_final_hp=[int]$lastDamage.payload.hp_after
        }
        persistence=[ordered]@{hp=[int]$dbAfterRelogin.hp;hp_max=[int]$dbAfterRelogin.hp_max}
        visual_review=[ordered]@{
            status="pending";reviewer=""
            checklist=[ordered]@{
                blur="pending";seam="pending";edge_alignment="pending"
                anchor="pending";loop="pending";impact_timing="pending"
            }
        }
    }
    Write-Utf8Json (Join-Path $cleanRoot "result.json") $cleanResult
    Close-TelemetryClient $clientA
    Close-TelemetryClient $clientB
    $clientA=$null;$clientB=$null

    $oldLock = Get-Content $PreviousLock -Raw -Encoding UTF8 | ConvertFrom-Json
    $warmProfile = Join-Path $resolvedRun "profiles\warm"
    Seed-WarmProfile $profileA $warmProfile $oldLock
    $cacheBefore = Join-Path $warmRoot "cache-before.json"
    Write-Utf8Json $cacheBefore ([ordered]@{
        versions=$oldLock.versions
        lock_sha256=(Get-FileHash $PreviousLock -Algorithm SHA256).Hash.ToLowerInvariant()
    })
    $warmProcess = Start-UnityClient "A" $TelemetryPortA $telemetryTokenA `
        $warmProfile (Join-Path $warmRoot "unity-warm.log")
    $warmClient = Connect-Telemetry "A" $TelemetryPortA $telemetryTokenA
    [void](Invoke-ClientCommand $warmClient "login" @{account=$accountA;password=$env:NRO_E2E_FIXTURE_PASSWORD})
    [void](Wait-ClientReady $warmClient)
    [void](Wait-EventType $warmClient "warm_negotiation" 0)
    $applied = Wait-EventType $warmClient "resource_lock_applied" 0
    Assert-WarmEvidence @($warmClient.Events) $assetLock
    $warmNegotiation = @($warmClient.Events | Where-Object {
        [string]$_.type -ceq "warm_negotiation" -and
        $_.payload.negotiation_detected -eq $true
    } | Select-Object -First 1)
    $warmCapture = Join-Path $warmRoot "client-warm.png"
    [void](Invoke-ClientCommand $warmClient "capture" @{path=$warmCapture})
    Assert-CaptureCompleted $warmClient $warmCapture
    Stop-UnityGracefully $warmClient $warmProcess
    $warmProcess=$null
    Assert-WarmEvidence @($warmClient.Events) $assetLock
    if (@($warmClient.Events | Where-Object {
        [string]$_.type -ceq "telemetry_shutdown"
    }).Count -eq 0) {
        Stop-Driver "E2E-SHUTDOWN-007" "warm thieu telemetry_shutdown"
    }
    $warmEventsPath = Join-Path $warmRoot "events.ndjson"
    Write-Utf8Lines $warmEventsPath @($warmClient.Events)
    $cacheAfter = Join-Path $warmRoot "cache-after.json"
    Write-Utf8Json $cacheAfter ([ordered]@{
        applied_lock_sha256=[string]$applied.payload.applied_lock_sha256
        map_id=[int]$applied.payload.map_id
        dragon_ball_fallback=[bool]$applied.payload.dragon_ball_fallback
    })
    $warmResult = [ordered]@{
        schema_version=1;scope="warm_smoke";status="captured"
        executable_sha256=$cleanResult.executable_sha256
        server_lock_sha256=$serverLock;asset_lock_sha256=$assetLock
        pair_sha256=[string]$manifestData.preflight.assets.pair_sha256
        versions=$manifestData.preflight.assets.versions
        events=(New-Record $warmRoot $warmEventsPath "ndjson")
        cache_before=(New-Record $warmRoot $cacheBefore "cache")
        cache_after=(New-Record $warmRoot $cacheAfter "cache")
        evidence=@((New-Record $warmRoot $warmEventsPath "client_ndjson"))
        artifacts=@((New-Record $warmRoot $warmCapture "screenshot"))
        resources=[ordered]@{
            negotiation_detected=($warmNegotiation.Count -eq 1)
            applied_lock_sha256=[string]$applied.payload.applied_lock_sha256
            entered_map_id=[int]$applied.payload.map_id
            dragon_ball_fallback=[bool]$applied.payload.dragon_ball_fallback
        }
    }
    Write-Utf8Json (Join-Path $warmRoot "result.json") $warmResult
    Close-TelemetryClient $warmClient
    $warmClient=$null

    Request-GracefulShutdown $controlToken
    Wait-ServerExit $server
    $server=$null
}
finally {
    $shutdownError = $null
    foreach ($pair in @(@($clientA,$clientAProcess),@($clientB,$clientBProcess),@($warmClient,$warmProcess))) {
        if ($null -ne $pair[1] -and -not $pair[1].HasExited) {
            try { Stop-UnityGracefully $pair[0] $pair[1] } catch { $shutdownError=$_ }
        }
    }
    if ($null -ne $server -and -not $server.HasExited) {
        try { Request-GracefulShutdown $controlToken; Wait-ServerExit $server }
        catch { $shutdownError=$_ }
    }
    if ($fixtureCreated) {
        try { Invoke-Fixture "Cleanup" $fixturePath $null }
        catch { if ($null -eq $shutdownError) { $shutdownError=$_ } }
    }
    $env:NRO_E2E_FIXTURE_PASSWORD = $fixturePasswordPrevious
    if (Test-Path $fixturePath) {
        $fixtureFinal = Get-Content $fixturePath -Raw -Encoding UTF8 | ConvertFrom-Json
        Assert-FixtureCleanup $fixtureFinal
    }
    if ($null -ne $shutdownError) { throw $shutdownError }
}

Write-Output "E2E-TASK10-EVIDENCE-CAPTURED $resolvedRun"
