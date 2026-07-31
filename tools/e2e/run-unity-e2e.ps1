[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$UnityProject,
    [Parameter(Mandatory = $true)]
    [string]$ServerRoot,
    [Parameter(Mandatory = $true)]
    [string]$ArtRoot,
    [string]$Dsn = $env:NRO_TEST_DSN,
    [string]$UnityExe = "D:\Unity\2022.3.62f2\Editor\Unity.exe",
    [string]$GoExe,
    [string]$OutputRoot,
    [string]$PreviousLock,
    [string]$GameplayScript,
    [int]$GamePort = 14445,
    [int]$ControlPort = 19091,
    [int]$TelemetryPortA = 19092,
    [int]$TelemetryPortB = 19093,
    [int]$SourceTimeoutSec = 60,
    [int]$DbTimeoutSec = 30,
    [int]$UnityImportTimeoutSec = 1200,
    [int]$UnityBuildTimeoutSec = 600,
    [int]$ServerStartTimeoutSec = 15,
    [int]$ServerShutdownTimeoutSec = 15,
    [int]$LoginTimeoutSec = 30,
    [int]$CommandTimeoutSec = 10,
    [int]$EvidenceTimeoutSec = 30,
    [switch]$PreflightOnly,
    [switch]$ContractOnly,
    [switch]$ControlSmokeOnly,
    [string]$OracleRunRoot,
    [switch]$OracleOnly,
    [switch]$SelfTest
)

$ErrorActionPreference = "Stop"
$script:PhaseDurations = [ordered]@{}

function Stop-E2E {
    param([string]$Code, [string]$Message)
    throw "$Code`: $Message"
}

function Resolve-FullPath {
    param([string]$Path)
    return [IO.Path]::GetFullPath($Path).TrimEnd(
        [IO.Path]::DirectorySeparatorChar,
        [IO.Path]::AltDirectorySeparatorChar)
}

function Write-Utf8Json {
    param([string]$Path, $Value)
    $json = $Value | ConvertTo-Json -Depth 12
    [IO.File]::WriteAllText(
        (Resolve-FullPath $Path),
        $json,
        (New-Object Text.UTF8Encoding($false)))
}

function Invoke-Phase {
    param([string]$Name, [scriptblock]$Action)
    $clock = [Diagnostics.Stopwatch]::StartNew()
    try {
        return & $Action
    }
    finally {
        $clock.Stop()
        $script:PhaseDurations[$Name] = [Math]::Round($clock.Elapsed.TotalSeconds, 3)
    }
}

function Assert-FiniteTimeouts {
    $values = @(
        $SourceTimeoutSec, $DbTimeoutSec, $UnityImportTimeoutSec,
        $UnityBuildTimeoutSec, $ServerStartTimeoutSec,
        $ServerShutdownTimeoutSec, $LoginTimeoutSec, $CommandTimeoutSec,
        $EvidenceTimeoutSec)
    if (@($values | Where-Object { $_ -le 0 }).Count -ne 0) {
        Stop-E2E "E2E-TIMEOUT-001" "moi phase phai co timeout duong, huu han"
    }
}

function Assert-TestDsn {
    param([string]$Value)
    try {
        $uri = [Uri]$Value
    }
    catch {
        Stop-E2E "E2E-DB-002" "DSN khong hop le"
    }
    if (($uri.Scheme -ne "postgres" -and $uri.Scheme -ne "postgresql") -or
        $uri.AbsolutePath.Trim('/') -cne "nro_test") {
        Stop-E2E "E2E-DB-002" "runner chi nhan database nro_test"
    }
}

function Assert-PortFree {
    param([int]$Port)
    $probe = $null
    try {
        $probe = New-Object Net.Sockets.TcpListener(
            [Net.IPAddress]::Loopback,
            $Port)
        $probe.Start()
    }
    catch {
        Stop-E2E "E2E-PORT-003" "cong loopback $Port dang ban"
    }
    finally {
        if ($null -ne $probe) {
            $probe.Stop()
        }
    }
}

function Get-UnityAssetLockOwners {
    param([Parameter(Mandatory = $true)][string]$Path)
    if ($null -eq ("NroE2ERestartManager" -as [type])) {
        $source = @'
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.Runtime.InteropServices;
using System.Text;

public static class NroE2ERestartManager
{
    const int ErrorSuccess = 0;
    const int ErrorMoreData = 234;

    [StructLayout(LayoutKind.Sequential)]
    struct UniqueProcess
    {
        public int ProcessId;
        public System.Runtime.InteropServices.ComTypes.FILETIME StartTime;
    }

    enum AppType
    {
        Unknown = 0, MainWindow = 1, OtherWindow = 2, Service = 3,
        Explorer = 4, Console = 5, Critical = 1000
    }

    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
    struct ProcessInfo
    {
        public UniqueProcess Process;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 256)]
        public string AppName;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 64)]
        public string ServiceName;
        public AppType ApplicationType;
        public uint AppStatus;
        public uint SessionId;
        [MarshalAs(UnmanagedType.Bool)]
        public bool Restartable;
    }

    [DllImport("rstrtmgr.dll", CharSet = CharSet.Unicode)]
    static extern int RmStartSession(
        out uint handle, int flags, StringBuilder sessionKey);
    [DllImport("rstrtmgr.dll", CharSet = CharSet.Unicode)]
    static extern int RmRegisterResources(
        uint handle, uint fileCount, string[] files,
        uint appCount, IntPtr apps, uint serviceCount, string[] services);
    [DllImport("rstrtmgr.dll")]
    static extern int RmGetList(
        uint handle, out uint needed, ref uint count,
        [In, Out] ProcessInfo[] info, ref uint reasons);
    [DllImport("rstrtmgr.dll")]
    static extern int RmEndSession(uint handle);

    public static string[] Query(string path)
    {
        uint handle;
        var key = new StringBuilder(33);
        int result = RmStartSession(out handle, 0, key);
        if (result != ErrorSuccess) throw new Win32Exception(result);
        try
        {
            result = RmRegisterResources(
                handle, 1, new[] { path }, 0, IntPtr.Zero, 0, null);
            if (result != ErrorSuccess) throw new Win32Exception(result);
            uint needed = 0, count = 0, reasons = 0;
            result = RmGetList(handle, out needed, ref count, null, ref reasons);
            if (result == ErrorSuccess) return new string[0];
            if (result != ErrorMoreData) throw new Win32Exception(result);
            var info = new ProcessInfo[needed];
            count = needed;
            result = RmGetList(
                handle, out needed, ref count, info, ref reasons);
            if (result != ErrorSuccess) throw new Win32Exception(result);
            var owners = new List<string>();
            for (int index = 0; index < count; index++)
            {
                int processId = info[index].Process.ProcessId;
                string processName;
                try
                {
                    processName = Process.GetProcessById(processId).ProcessName;
                }
                catch
                {
                    processName = "exited";
                }
                owners.Add(processId + ":" + processName);
            }
            return owners.ToArray();
        }
        finally
        {
            RmEndSession(handle);
        }
    }
}
'@
        try {
            Add-Type -TypeDefinition $source -Language CSharp -ErrorAction Stop
        }
        catch {
            Stop-E2E "E2E-PREFLIGHT-LOCK-004" `
                "khong nap duoc Restart Manager de xac minh asset lock"
        }
    }
    try {
        return @([NroE2ERestartManager]::Query(
            [IO.Path]::GetFullPath($Path)))
    }
    catch {
        Stop-E2E "E2E-PREFLIGHT-LOCK-004" `
            "khong chung minh duoc asset lock owner=0"
    }
}

function Assert-NoUnityAssetLock {
    param([string]$Root)
    foreach ($relative in @("Library\ArtifactDB-lock", "Library\SourceAssetDB-lock")) {
        $path = Join-Path $Root $relative
        if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
            continue
        }
        $owners = @(Get-UnityAssetLockOwners -Path $path)
        if ($owners.Count -ne 0) {
            Stop-E2E "E2E-PREFLIGHT-LOCK-004" `
                "Unity source dang co asset lock: $relative; owners=$([string]::Join(',', $owners))"
        }
        Write-Host "E2E-UNITY-LOCK-STALE owner_count=0 path=$relative"
    }
}

function Invoke-PsqlScalar {
    param([string]$Sql)
    $psql = if ($env:NRO_PSQL) { $env:NRO_PSQL } else { "D:\PostgreSQL\bin\psql.exe" }
    if (-not (Test-Path -LiteralPath $psql -PathType Leaf)) {
        Stop-E2E "E2E-DB-002" "khong tim thay psql"
    }
    $job = Start-Job -ScriptBlock {
        param($Exe, $Connection, $Query)
        & $Exe --no-psqlrc --set=ON_ERROR_STOP=1 --quiet --tuples-only `
            --no-align "--dbname=$Connection" "--command=$Query" 2>&1
        exit $LASTEXITCODE
    } -ArgumentList $psql, $Dsn, $Sql
    try {
        if ($null -eq (Wait-Job -Job $job -Timeout $DbTimeoutSec)) {
            Stop-E2E "E2E-TIMEOUT-DB" "DB preflight vuot $DbTimeoutSec giay"
        }
        $output = @(Receive-Job -Job $job)
        if ($job.State -ne "Completed") {
            Stop-E2E "E2E-DB-002" "DB preflight that bai"
        }
        return "$($output | Select-Object -Last 1)".Trim()
    }
    finally {
        Remove-Job -Job $job -Force -ErrorAction SilentlyContinue
    }
}

function ConvertTo-NativeArgument {
    param([AllowEmptyString()][string]$Value)
    if ($Value.Length -gt 0 -and $Value -notmatch '[\s"]') {
        return $Value
    }
    $builder = New-Object Text.StringBuilder
    [void]$builder.Append('"')
    $slashes = 0
    foreach ($character in $Value.ToCharArray()) {
        if ($character -eq '\') {
            $slashes++
            continue
        }
        if ($character -eq '"') {
            [void]$builder.Append(('\' * (($slashes * 2) + 1)))
            [void]$builder.Append('"')
            $slashes = 0
            continue
        }
        if ($slashes -gt 0) {
            [void]$builder.Append(('\' * $slashes))
            $slashes = 0
        }
        [void]$builder.Append($character)
    }
    if ($slashes -gt 0) {
        [void]$builder.Append(('\' * ($slashes * 2)))
    }
    [void]$builder.Append('"')
    return $builder.ToString()
}

function Invoke-ExternalWithTimeout {
    param(
        [string]$Name,
        [string]$FilePath,
        [string[]]$Arguments,
        [int]$TimeoutSec,
        [string]$StdoutPath,
        [string]$StderrPath,
        [string]$WorkingDirectory = $ServerRoot
    )
    $start = New-Object Diagnostics.ProcessStartInfo
    $start.FileName = $FilePath
    $start.Arguments = [string]::Join(
        " ",
        @($Arguments | ForEach-Object {
            ConvertTo-NativeArgument ([string]$_)
        }))
    $start.WorkingDirectory = $WorkingDirectory
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    $process = New-Object Diagnostics.Process
    $process.StartInfo = $start
    if (-not $process.Start()) {
        Stop-E2E "E2E-$Name-FAILED" "$Name khong khoi dong duoc"
    }
    $stdoutTask = $process.StandardOutput.ReadToEndAsync()
    $stderrTask = $process.StandardError.ReadToEndAsync()
    if (-not $process.WaitForExit($TimeoutSec * 1000)) {
        # Chi dung cho import/build/tool ngoai. Server production tuyet doi
        # khong di qua ham nay: server shutdown chi qua control endpoint.
        $process.Kill()
        $process.WaitForExit()
        [IO.File]::WriteAllText($StdoutPath, $stdoutTask.Result)
        [IO.File]::WriteAllText($StderrPath, $stderrTask.Result)
        Stop-E2E "E2E-TIMEOUT-$Name" "$Name vuot $TimeoutSec giay"
    }
    $process.WaitForExit()
    [IO.File]::WriteAllText($StdoutPath, $stdoutTask.Result)
    [IO.File]::WriteAllText($StderrPath, $stderrTask.Result)
    $exitCode = $process.ExitCode
    if ($exitCode -ne 0) {
        Stop-E2E "E2E-$Name-FAILED" "$Name exit $exitCode"
    }
}

function Invoke-GitReadWithTimeout {
    param([string]$Arguments)
    $start = New-Object Diagnostics.ProcessStartInfo
    $start.FileName = "git.exe"
    $start.Arguments = $Arguments
    $start.WorkingDirectory = $ServerRoot
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    $process = New-Object Diagnostics.Process
    $process.StartInfo = $start
    if (-not $process.Start()) {
        Stop-E2E "E2E-SOURCE-005" "khong khoi dong duoc git"
    }
    $stdoutTask = $process.StandardOutput.ReadToEndAsync()
    $stderrTask = $process.StandardError.ReadToEndAsync()
    if (-not $process.WaitForExit($SourceTimeoutSec * 1000)) {
        $process.Kill()
        $process.WaitForExit()
        Stop-E2E "E2E-TIMEOUT-SOURCE" "git preflight vuot $SourceTimeoutSec giay"
    }
    $stdout = $stdoutTask.Result
    [void]$stderrTask.Result
    if ($process.ExitCode -ne 0) {
        Stop-E2E "E2E-SOURCE-005" "git preflight exit $($process.ExitCode)"
    }
    return $stdout
}

function Get-SourceSnapshot {
    $head = (Invoke-GitReadWithTimeout "-C `"$ServerRoot`" rev-parse HEAD").Trim()
    if ($head -notmatch '^[0-9a-f]{40}$') {
        Stop-E2E "E2E-SOURCE-005" "ServerRoot khong phai git worktree"
    }
    $dirtyText = Invoke-GitReadWithTimeout `
        "-C `"$ServerRoot`" status --porcelain --untracked-files=all"
    $dirtyLines = @($dirtyText -split "`r?`n" | Where-Object { $_ -ne "" })
    return [ordered]@{
        commit = $head
        dirty = ($dirtyLines.Count -gt 0)
        dirty_count = $dirtyLines.Count
    }
}

function Assert-ArtifactPair {
    param([string]$Root, [string]$Server)
    $lockPath = Join-Path $Root "asset-lock.json"
    if (-not (Test-Path -LiteralPath $lockPath -PathType Leaf) -or
        -not (Test-Path -LiteralPath (Join-Path $Root "server") -PathType Container) -or
        -not (Test-Path -LiteralPath (Join-Path $Root "client") -PathType Container)) {
        Stop-E2E "E2E-ASSET-006" "ArtRoot phai co asset-lock.json + server/client"
    }
    $lock = Get-Content -LiteralPath $lockPath -Raw -Encoding UTF8 | ConvertFrom-Json
    if ($lock.schema_version -ne 1 -or
        ([string]$lock.server_client_pair_sha256) -notmatch '^[0-9a-f]{64}$') {
        Stop-E2E "E2E-ASSET-006" "asset lock/pair SHA khong hop le"
    }
    $embedded = Join-Path $Server "server\internal\game\assets\asset-lock.json"
    if (Test-Path -LiteralPath $embedded -PathType Leaf) {
        $want = (Get-FileHash -LiteralPath $lockPath -Algorithm SHA256).Hash
        $got = (Get-FileHash -LiteralPath $embedded -Algorithm SHA256).Hash
        if ($want -cne $got) {
            Stop-E2E "E2E-ASSET-006" "server source va ArtRoot khac asset lock"
        }
    }
    return [ordered]@{
        lock_sha256 = (Get-FileHash -LiteralPath $lockPath -Algorithm SHA256).Hash.ToLowerInvariant()
        pair_sha256 = ([string]$lock.server_client_pair_sha256).ToLowerInvariant()
        versions = $lock.versions
    }
}

function Assert-ProfileIsolation {
    param([string]$A, [string]$B)
    $aRoot = Resolve-FullPath $A
    $bRoot = Resolve-FullPath $B
    if ([string]::Equals($aRoot, $bRoot, [StringComparison]::OrdinalIgnoreCase)) {
        Stop-E2E "E2E-CACHE-002" "profile A/B trung nhau"
    }
    New-Item -ItemType Directory -Path $aRoot -Force | Out-Null
    New-Item -ItemType Directory -Path $bRoot -Force | Out-Null
    [IO.File]::WriteAllText((Join-Path $aRoot "ONLY_A"), "A")
    [IO.File]::WriteAllText((Join-Path $bRoot "ONLY_B"), "B")
    if ((Test-Path -LiteralPath (Join-Path $aRoot "ONLY_B")) -or
        (Test-Path -LiteralPath (Join-Path $bRoot "ONLY_A"))) {
        Stop-E2E "E2E-CACHE-002" "profile sentinel doc cheo"
    }
}

function Request-GracefulServerShutdown {
    param([int]$Port, [string]$Token)
    $request = [Net.HttpWebRequest]::Create("http://127.0.0.1:$Port/shutdown")
    $request.Method = "POST"
    $request.Headers["Authorization"] = "Bearer $Token"
    $request.Timeout = $ServerShutdownTimeoutSec * 1000
    try {
        $response = $request.GetResponse()
        if ([int]$response.StatusCode -ne 202) {
            Stop-E2E "E2E-SHUTDOWN-007" "control status khong phai 202"
        }
        $response.Close()
    }
    catch {
        Stop-E2E "E2E-SHUTDOWN-007" "graceful control that bai"
    }
}

function Invoke-ControlSmokeLifecycle {
    param(
        [Parameter(Mandatory = $true)][scriptblock]$CreateFixture,
        [Parameter(Mandatory = $true)][scriptblock]$ValidateFixture,
        [Parameter(Mandatory = $true)][scriptblock]$BuildServer,
        [Parameter(Mandatory = $true)][scriptblock]$StartServer,
        [Parameter(Mandatory = $true)][scriptblock]$WaitControlReady,
        [Parameter(Mandatory = $true)][scriptblock]$RequestShutdown,
        [Parameter(Mandatory = $true)][scriptblock]$WaitServerExit,
        [Parameter(Mandatory = $true)][scriptblock]$CleanupFixture,
        [Parameter(Mandatory = $true)][scriptblock]$WriteEvidence
    )
    $failure = $null
    $cleanupFailure = $null
    $process = $null
    try {
        & $CreateFixture
        & $ValidateFixture
        & $BuildServer
        $process = & $StartServer
        & $WaitControlReady $process
        & $RequestShutdown
        & $WaitServerExit $process
    }
    catch {
        $failure = $_
    }
    finally {
        try {
            & $CleanupFixture
        }
        catch {
            $cleanupFailure = $_
        }
    }
    if ($null -ne $failure) {
        throw $failure
    }
    if ($null -ne $cleanupFailure) {
        throw $cleanupFailure
    }
    # Bang chung cleanup_completed chi duoc tao sau khi finally cleanup xanh.
    & $WriteEvidence
}

function Resolve-GoExecutable {
    param([string]$Requested)
    if (-not [string]::IsNullOrWhiteSpace($Requested)) {
        $candidate = Resolve-FullPath $Requested
        if (-not (Test-Path -LiteralPath $candidate -PathType Leaf)) {
            Stop-E2E "E2E-GO-010" "khong tim thay Go executable da chi dinh"
        }
        return $candidate
    }
    $command = Get-Command "go.exe" -ErrorAction SilentlyContinue
    if ($null -ne $command) {
        return Resolve-FullPath $command.Source
    }
    $machineDefault = "D:\Go\bin\go.exe"
    if (Test-Path -LiteralPath $machineDefault -PathType Leaf) {
        return $machineDefault
    }
    Stop-E2E "E2E-GO-010" "khong tim thay go.exe trong PATH hay D:\Go\bin"
}

function New-ControlToken {
    $bytes = New-Object byte[] 32
    $random = [Security.Cryptography.RandomNumberGenerator]::Create()
    try {
        $random.GetBytes($bytes)
        return [Convert]::ToBase64String($bytes)
    }
    finally {
        $random.Dispose()
    }
}

function Invoke-FixtureAction {
    param(
        [string]$Action,
        [string]$FixtureManifest,
        [string]$FixtureRunId,
        [string]$LogPrefix
    )
    $fixtureScript = Join-Path $server "tools\e2e\New-Fixture.ps1"
    $arguments = @(
        "-NoProfile", "-ExecutionPolicy", "Bypass",
        "-File", $fixtureScript,
        "-Dsn", $Dsn,
        "-Action", $Action,
        "-ManifestPath", $FixtureManifest)
    if ($Action -eq "Create") {
        $arguments += @("-RunId", $FixtureRunId)
    }
    Invoke-ExternalWithTimeout -Name ("FIXTURE-" + $Action.ToUpperInvariant()) `
        -FilePath "powershell.exe" -Arguments $arguments `
        -TimeoutSec $DbTimeoutSec `
        -StdoutPath (Join-Path $runRoot "$LogPrefix.stdout.log") `
        -StderrPath (Join-Path $runRoot "$LogPrefix.stderr.log") `
        -WorkingDirectory $server
}

function Start-E2EServer {
    param(
        [string]$Executable,
        [string]$Token,
        [string]$StdoutPath,
        [string]$StderrPath
    )
    $start = New-Object Diagnostics.ProcessStartInfo
    $start.FileName = $Executable
    $start.WorkingDirectory = Join-Path $server "server"
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    $start.EnvironmentVariables["NRO_DSN"] = $Dsn
    $start.EnvironmentVariables["NRO_DIA_CHI"] = "127.0.0.1:$GamePort"
    $start.EnvironmentVariables["NRO_E2E"] = "1"
    $start.EnvironmentVariables["NRO_E2E_CONTROL_ADDR"] =
        "127.0.0.1:$ControlPort"
    $start.EnvironmentVariables["NRO_E2E_CONTROL_TOKEN"] = $Token
    $start.EnvironmentVariables["NRO_MUC_NHAT_KY"] = "info"
    $process = New-Object Diagnostics.Process
    $process.StartInfo = $start
    if (-not $process.Start()) {
        Stop-E2E "E2E-SERVER-START-011" "khong khoi dong duoc server"
    }
    $process | Add-Member -NotePropertyName E2EStdoutTask `
        -NotePropertyValue $process.StandardOutput.ReadToEndAsync()
    $process | Add-Member -NotePropertyName E2EStderrTask `
        -NotePropertyValue $process.StandardError.ReadToEndAsync()
    $process | Add-Member -NotePropertyName E2EStdoutPath `
        -NotePropertyValue $StdoutPath
    $process | Add-Member -NotePropertyName E2EStderrPath `
        -NotePropertyValue $StderrPath
    return $process
}

function Wait-E2EControlReady {
    param([Diagnostics.Process]$Process, [int]$Port, [string]$Token)
    $deadline = [DateTime]::UtcNow.AddSeconds($ServerStartTimeoutSec)
    while ([DateTime]::UtcNow -lt $deadline) {
        if ($Process.HasExited) {
            Stop-E2E "E2E-SERVER-START-011" `
                "server exit $($Process.ExitCode) truoc khi control ready"
        }
        $request = [Net.HttpWebRequest]::Create(
            "http://127.0.0.1:$Port/shutdown")
        $request.Method = "GET"
        $request.Headers["Authorization"] = "Bearer $Token"
        $request.Timeout = 500
        try {
            $response = $request.GetResponse()
            $response.Close()
        }
        catch [Net.WebException] {
            $response = $_.Exception.Response
            if ($null -ne $response) {
                try {
                    if ([int]$response.StatusCode -eq 405) {
                        return
                    }
                }
                finally {
                    $response.Close()
                }
            }
        }
        Start-Sleep -Milliseconds 100
    }
    Stop-E2E "E2E-TIMEOUT-SERVER-START" `
        "control khong ready trong $ServerStartTimeoutSec giay"
}

function Wait-E2EServerExit {
    param([Diagnostics.Process]$Process)
    if (-not $Process.WaitForExit($ServerShutdownTimeoutSec * 1000)) {
        Stop-E2E "E2E-TIMEOUT-SERVER-SHUTDOWN" `
            "server khong graceful exit trong $ServerShutdownTimeoutSec giay"
    }
    $Process.WaitForExit()
    [IO.File]::WriteAllText(
        $Process.E2EStdoutPath,
        $Process.E2EStdoutTask.Result)
    [IO.File]::WriteAllText(
        $Process.E2EStderrPath,
        $Process.E2EStderrTask.Result)
    if ($Process.ExitCode -ne 0) {
        Stop-E2E "E2E-SHUTDOWN-007" `
            "server graceful exit code $($Process.ExitCode)"
    }
}

function Get-E2EFailureScope {
    param(
        [bool]$IsControlSmoke,
        [bool]$IsContract
    )
    if ($IsControlSmoke) {
        return "task9_control_smoke"
    }
    if ($IsContract) {
        return "task9_contract"
    }
    return "task10_full"
}

function Get-E2ERunScope {
    param(
        [bool]$IsControlSmoke,
        [bool]$IsContract
    )
    return Get-E2EFailureScope `
        -IsControlSmoke $IsControlSmoke `
        -IsContract $IsContract
}

function Invoke-RunnerSelfTest {
    if ((Get-E2ERunScope -IsControlSmoke $false -IsContract $false) -cne
        "task10_full" -or
        (Get-E2ERunScope -IsControlSmoke $true -IsContract $false) -cne
        "task9_control_smoke" -or
        (Get-E2ERunScope -IsControlSmoke $false -IsContract $true) -cne
        "task9_contract") {
        throw "SELFTEST-FAIL: top-level manifest scope Task 9/10 sai"
    }
    if ((Get-E2EFailureScope -IsControlSmoke $false -IsContract $false) -cne
        "task10_full" -or
        (Get-E2EFailureScope -IsControlSmoke $true -IsContract $false) -cne
        "task9_control_smoke" -or
        (Get-E2EFailureScope -IsControlSmoke $false -IsContract $true) -cne
        "task9_contract") {
        throw "SELFTEST-FAIL: failure scope Task 9/10 sai"
    }

    $lockRoot = Join-Path ([IO.Path]::GetTempPath()) (
        "nro-e2e-lock-selftest-" + [Guid]::NewGuid().ToString("N"))
    $lockLibrary = Join-Path $lockRoot "Library"
    New-Item -ItemType Directory -Path $lockLibrary -Force | Out-Null
    $artifactLock = Join-Path $lockLibrary "ArtifactDB-lock"
    $sourceLock = Join-Path $lockLibrary "SourceAssetDB-lock"
    [IO.File]::WriteAllBytes($artifactLock, [byte[]](1, 2, 3))
    [IO.File]::WriteAllBytes($sourceLock, [byte[]](4, 5, 6))
    try {
        Assert-NoUnityAssetLock $lockRoot

        $held = [IO.File]::Open(
            $sourceLock,
            [IO.FileMode]::Open,
            [IO.FileAccess]::Read,
            [IO.FileShare]::None)
        try {
            try {
                Assert-NoUnityAssetLock $lockRoot
                throw "SELFTEST-FAIL: lock dang co owner da bi chap nhan"
            }
            catch {
                if (-not $_.Exception.Message.Contains(
                    "E2E-PREFLIGHT-LOCK-004")) {
                    throw
                }
                Write-Host "RUNNER-LOCK-OWNER-REJECTED"
            }
        }
        finally {
            $held.Dispose()
        }
    }
    finally {
        $temp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
        $resolvedLockRoot = [IO.Path]::GetFullPath($lockRoot)
        if ($resolvedLockRoot.StartsWith(
            $temp, [StringComparison]::OrdinalIgnoreCase) -and
            (Test-Path -LiteralPath $resolvedLockRoot)) {
            Remove-Item -LiteralPath $resolvedLockRoot -Recurse -Force
        }
    }

    $steps = New-Object System.Collections.Generic.List[string]
    $fixture = @{ cleanup_completed = $false }
    Invoke-ControlSmokeLifecycle `
        -CreateFixture { [void]$steps.Add("create") } `
        -ValidateFixture { [void]$steps.Add("validate") } `
        -BuildServer { [void]$steps.Add("build") } `
        -StartServer {
            [void]$steps.Add("start")
            return [pscustomobject]@{ HasExited = $false }
        } `
        -WaitControlReady { param($Process) [void]$steps.Add("ready") } `
        -RequestShutdown { [void]$steps.Add("shutdown") } `
        -WaitServerExit { param($Process) [void]$steps.Add("exit") } `
        -CleanupFixture {
            [void]$steps.Add("cleanup")
            $fixture.cleanup_completed = $true
        } `
        -WriteEvidence {
            if ($fixture.cleanup_completed -ne $true) {
                throw "SELFTEST-FAIL: evidence duoc ghi truoc cleanup thanh cong"
            }
            [void]$steps.Add("evidence")
        }
    $actual = [string]::Join(",", $steps)
    $expected = "create,validate,build,start,ready,shutdown,exit,cleanup,evidence"
    if ($actual -cne $expected) {
        throw "SELFTEST-FAIL: E2E-CONTROL-SMOKE-001 order=$actual"
    }

    $failureSteps = New-Object System.Collections.Generic.List[string]
    try {
        Invoke-ControlSmokeLifecycle `
            -CreateFixture { [void]$failureSteps.Add("create") } `
            -ValidateFixture {
                [void]$failureSteps.Add("validate")
                throw "expected-smoke-failure"
            } `
            -BuildServer { [void]$failureSteps.Add("build") } `
            -StartServer { throw "start must not run" } `
            -WaitControlReady { throw "ready must not run" } `
            -RequestShutdown { throw "shutdown must not run" } `
            -WaitServerExit { throw "exit must not run" } `
            -CleanupFixture { [void]$failureSteps.Add("cleanup") } `
            -WriteEvidence { [void]$failureSteps.Add("evidence") }
        throw "SELFTEST-FAIL: control-smoke failure bi nuot"
    }
    catch {
        if (-not $_.Exception.Message.Contains("expected-smoke-failure")) {
            throw
        }
    }
    if ([string]::Join(",", $failureSteps) -cne "create,validate,cleanup") {
        throw "SELFTEST-FAIL: E2E-CONTROL-SMOKE-002 finally/evidence sai"
    }
}

if ($SelfTest) {
    Invoke-RunnerSelfTest
    Write-Output "RUNNER-SELFTEST-PASS"
    exit 0
}

if ($OracleOnly) {
    if ([string]::IsNullOrWhiteSpace($OracleRunRoot)) {
        Stop-E2E "E2E-EVIDENCE-001" "OracleOnly can -OracleRunRoot"
    }
    & (Join-Path $PSScriptRoot "Test-Evidence.ps1") `
        -OracleRunRoot (Resolve-FullPath $OracleRunRoot)
    exit $LASTEXITCODE
}

Assert-FiniteTimeouts
$unityRoot = Resolve-FullPath $UnityProject
$server = Resolve-FullPath $ServerRoot
$art = Resolve-FullPath $ArtRoot
if ([string]::IsNullOrWhiteSpace($OutputRoot)) {
    $OutputRoot = Join-Path $server "e2e-runs"
}
$runId = "run-" + [DateTime]::UtcNow.ToString("yyyyMMddTHHmmssZ") + "-" +
    [Guid]::NewGuid().ToString("N").Substring(0, 8)
$runRoot = Join-Path (Resolve-FullPath $OutputRoot) $runId
New-Item -ItemType Directory -Path $runRoot -Force | Out-Null
$manifestPath = Join-Path $runRoot "manifest.json"

try {
    $preflight = Invoke-Phase "source_preflight" {
        Assert-TestDsn $Dsn
        $ports = if ($ControlSmokeOnly) {
            @($GamePort, $ControlPort)
        }
        else {
            @($GamePort, $ControlPort, $TelemetryPortA, $TelemetryPortB)
        }
        foreach ($port in $ports) {
            Assert-PortFree $port
        }
        $source = Get-SourceSnapshot
        if (-not $ControlSmokeOnly -and $source.dirty) {
            Stop-E2E "E2E-SOURCE-005" `
                "worktree dang dirty ($($source.dirty_count) path); can commit duy nhat truoc full E2E"
        }
        if ($ControlSmokeOnly) {
            [ordered]@{
                mode = "control_smoke"
                source = $source
                go_executable = Resolve-GoExecutable $GoExe
                ports = $ports
            }
            return
        }
        if (-not (Test-Path -LiteralPath $UnityExe -PathType Leaf)) {
            Stop-E2E "E2E-UNITY-008" "khong tim thay Unity 2022.3.62f2"
        }
        $versionText = Get-Content -LiteralPath `
            (Join-Path $unityRoot "ProjectSettings\ProjectVersion.txt") `
            -Raw -Encoding UTF8
        if ($versionText -notmatch 'm_EditorVersion:\s*2022\.3\.62f2(\r?\n|$)') {
            Stop-E2E "E2E-UNITY-008" "project phai dung Unity 2022.3.62f2"
        }
        Assert-NoUnityAssetLock $unityRoot
        $assets = Assert-ArtifactPair -Root $art -Server $server
        [ordered]@{
            mode = "unity_contract"
            source = $source
            assets = $assets
            unity_version = "2022.3.62f2"
            unity_project = $unityRoot
            ports = $ports
        }
    }

    $schemaCount = Invoke-Phase "db_preflight" {
        Invoke-PsqlScalar @"
SELECT count(*)
FROM information_schema.columns
WHERE table_schema='public'
  AND ((table_name='tai_khoan' AND column_name IN ('id','ten_dang_nhap','mat_khau_bam'))
    OR (table_name='nhan_vat' AND column_name IN
      ('id','tai_khoan_id','lang','kieu_toc','map_id','vi_tri_x','vi_tri_y',
       'hp_toi_da','hp_hien_tai','chakra_toi_da','chakra_hien_tai')));
"@
    }
    if ([int]$schemaCount -ne 14) {
        Stop-E2E "E2E-DB-002" "schema nro_test thieu cot fixture (actual=$schemaCount)"
    }

    $profileRoots = @()
    if (-not $ControlSmokeOnly) {
        $profileA = Join-Path $runRoot "profiles\A"
        $profileB = Join-Path $runRoot "profiles\B"
        Invoke-Phase "profile_preflight" {
            Assert-ProfileIsolation -A $profileA -B $profileB
        } | Out-Null
        $profileRoots = @($profileA, $profileB)
    }

    $runScope = Get-E2ERunScope `
        -IsControlSmoke ([bool]$ControlSmokeOnly) `
        -IsContract ([bool]$ContractOnly)
    $manifest = [ordered]@{
        schema_version = 1
        scope = $runScope
        run_id = $runId
        status = "preflight_pass"
        preflight = $preflight
        db_schema_column_count = [int]$schemaCount
        profile_roots = $profileRoots
        timeouts_seconds = [ordered]@{
            source = $SourceTimeoutSec; db = $DbTimeoutSec
            unity_import = $UnityImportTimeoutSec; unity_build = $UnityBuildTimeoutSec
            server_start = $ServerStartTimeoutSec; server_shutdown = $ServerShutdownTimeoutSec
            login = $LoginTimeoutSec; command = $CommandTimeoutSec
            evidence = $EvidenceTimeoutSec
        }
        durations_seconds = $script:PhaseDurations
    }
    Write-Utf8Json -Path $manifestPath -Value $manifest
    if ($PreflightOnly) {
        Write-Output "E2E-PREFLIGHT-PASS $manifestPath"
        exit 0
    }

    if ($ControlSmokeOnly) {
        $fixtureManifest = Join-Path $runRoot "fixture.json"
        $serverDirectory = Join-Path $runRoot "server"
        New-Item -ItemType Directory -Path $serverDirectory -Force | Out-Null
        $serverExecutable = Join-Path $serverDirectory "nroserver-e2e.exe"
        $serverStdout = Join-Path $serverDirectory "server.stdout.log"
        $serverStderr = Join-Path $serverDirectory "server.stderr.log"
        $controlToken = New-ControlToken
        $controlState = @{
            process = $null
            events = New-Object System.Collections.Generic.List[object]
            clock = [Diagnostics.Stopwatch]::StartNew()
        }
        $fixturePasswordPrevious = [Environment]::GetEnvironmentVariable(
            "NRO_E2E_FIXTURE_PASSWORD",
            [EnvironmentVariableTarget]::Process)
        [Environment]::SetEnvironmentVariable(
            "NRO_E2E_FIXTURE_PASSWORD",
            ([Guid]::NewGuid().ToString("N") + [Guid]::NewGuid().ToString("N")),
            [EnvironmentVariableTarget]::Process)
        try {
            Invoke-ControlSmokeLifecycle `
                -CreateFixture {
                    Invoke-Phase "fixture_create" {
                        Invoke-FixtureAction -Action "Create" `
                            -FixtureManifest $fixtureManifest `
                            -FixtureRunId ("control-" + $runId) `
                            -LogPrefix "fixture-create"
                    } | Out-Null
                } `
                -ValidateFixture {
                    Invoke-Phase "fixture_validate" {
                        Invoke-FixtureAction -Action "Validate" `
                            -FixtureManifest $fixtureManifest `
                            -FixtureRunId $null `
                            -LogPrefix "fixture-validate"
                    } | Out-Null
                } `
                -BuildServer {
                    Invoke-Phase "server_build" {
                        Invoke-ExternalWithTimeout -Name "SERVER-BUILD" `
                            -FilePath ([string]$preflight.go_executable) `
                            -Arguments @(
                                "build", "-o", $serverExecutable,
                                "./cmd/nroserver") `
                            -TimeoutSec $UnityBuildTimeoutSec `
                            -StdoutPath (Join-Path $runRoot "server-build.stdout.log") `
                            -StderrPath (Join-Path $runRoot "server-build.stderr.log") `
                            -WorkingDirectory (Join-Path $server "server")
                    } | Out-Null
                } `
                -StartServer {
                    $controlState.process = Start-E2EServer `
                        -Executable $serverExecutable `
                        -Token $controlToken `
                        -StdoutPath $serverStdout `
                        -StderrPath $serverStderr
                    return $controlState.process
                } `
                -WaitControlReady {
                    param($Process)
                    Invoke-Phase "server_start" {
                        Wait-E2EControlReady `
                            -Process $Process `
                            -Port $ControlPort `
                            -Token $controlToken
                    } | Out-Null
                    [void]$controlState.events.Add([ordered]@{
                        monotonic_ms = [Math]::Max(
                            1,
                            [long]$controlState.clock.ElapsedMilliseconds)
                        utc = [DateTime]::UtcNow.ToString("O")
                        client = "runner"
                        sequence = 1
                        type = "control_ready"
                        payload = @{ port = $ControlPort }
                    })
                } `
                -RequestShutdown {
                    Invoke-Phase "server_shutdown_request" {
                        Request-GracefulServerShutdown `
                            -Port $ControlPort `
                            -Token $controlToken
                    } | Out-Null
                } `
                -WaitServerExit {
                    param($Process)
                    Invoke-Phase "server_shutdown_wait" {
                        Wait-E2EServerExit -Process $Process
                    } | Out-Null
                    [void]$controlState.events.Add([ordered]@{
                        monotonic_ms = [Math]::Max(
                            1,
                            [long]$controlState.clock.ElapsedMilliseconds)
                        utc = [DateTime]::UtcNow.ToString("O")
                        client = "runner"
                        sequence = 2
                        type = "shutdown_completed"
                        payload = @{ exit_code = $Process.ExitCode }
                    })
                } `
                -CleanupFixture {
                    $shutdownFailure = $null
                    $running = $controlState.process
                    if ($null -ne $running -and -not $running.HasExited) {
                        try {
                            Request-GracefulServerShutdown `
                                -Port $ControlPort `
                                -Token $controlToken
                            Wait-E2EServerExit -Process $running
                        }
                        catch {
                            # Khong Kill server. Van phai cleanup DB, roi bao loi graceful.
                            $shutdownFailure = $_
                        }
                    }
                    Invoke-Phase "fixture_cleanup" {
                        Invoke-FixtureAction -Action "Cleanup" `
                            -FixtureManifest $fixtureManifest `
                            -FixtureRunId $null `
                            -LogPrefix "fixture-cleanup"
                    } | Out-Null
                    if ($null -ne $shutdownFailure) {
                        throw $shutdownFailure
                    }
                } `
                -WriteEvidence {
                    $fixture = Get-Content -LiteralPath $fixtureManifest `
                        -Raw -Encoding UTF8 | ConvertFrom-Json
                    if ($fixture.cleanup_completed -ne $true) {
                        Stop-E2E "E2E-FIXTURE-008" `
                            "khong ghi evidence truoc cleanup thanh cong"
                    }
                    $relativeLog = "server/server.stdout.log"
                    $result = [ordered]@{
                        schema_version = 1
                        fixture = [ordered]@{
                            prefix = $fixture.prefix
                            accounts = @($fixture.accounts)
                            cleanup_completed = $fixture.cleanup_completed
                        }
                        artifacts = @([ordered]@{
                            path = $relativeLog
                            sha256 = (
                                Get-FileHash -LiteralPath $serverStdout `
                                    -Algorithm SHA256
                            ).Hash.ToLowerInvariant()
                        })
                        events = $controlState.events.ToArray()
                    }
                    $resultPath = Join-Path $runRoot "result.json"
                    Write-Utf8Json -Path $resultPath -Value $result
                    Invoke-ExternalWithTimeout -Name "EVIDENCE" `
                        -FilePath "powershell.exe" `
                        -Arguments @(
                            "-NoProfile", "-ExecutionPolicy", "Bypass",
                            "-File", (Join-Path $server "tools\e2e\Test-Evidence.ps1"),
                            "-EvidenceRoot", $runRoot) `
                        -TimeoutSec $EvidenceTimeoutSec `
                        -StdoutPath (Join-Path $runRoot "evidence.stdout.log") `
                        -StderrPath (Join-Path $runRoot "evidence.stderr.log") `
                        -WorkingDirectory $server
                    $manifest.status = "control_smoke_pass"
                    $manifest.fixture_manifest = $fixtureManifest
                    $manifest.result = $resultPath
                    $manifest.server_executable = [ordered]@{
                        path = $serverExecutable
                        sha256 = (
                            Get-FileHash -LiteralPath $serverExecutable `
                                -Algorithm SHA256
                        ).Hash.ToLowerInvariant()
                    }
                    $manifest.durations_seconds = $script:PhaseDurations
                    Write-Utf8Json -Path $manifestPath -Value $manifest
                }
        }
        finally {
            [Environment]::SetEnvironmentVariable(
                "NRO_E2E_FIXTURE_PASSWORD",
                $fixturePasswordPrevious,
                [EnvironmentVariableTarget]::Process)
        }
        Write-Output "E2E-CONTROL-SMOKE-PASS $(Join-Path $runRoot 'result.json')"
        exit 0
    }

    if ([string]::IsNullOrWhiteSpace($PreviousLock) -or
        -not (Test-Path -LiteralPath $PreviousLock -PathType Leaf)) {
        Stop-E2E "E2E-ASSET-006" "full contract build can -PreviousLock"
    }

    $snapshot = Join-Path $runRoot "unity-snapshot"
    $overlayOut = Join-Path $runRoot "unity-overlay.stdout.log"
    $overlayErr = Join-Path $runRoot "unity-overlay.stderr.log"
    Invoke-Phase "unity_import" {
        Invoke-ExternalWithTimeout -Name "UNITY-IMPORT" `
            -FilePath "powershell.exe" `
            -Arguments @(
                "-NoProfile", "-ExecutionPolicy", "Bypass",
                "-File", (Join-Path $server "tools\e2e\..\unity-overlay\apply-overlay.ps1"),
                "-UnityProject", $unityRoot,
                "-SnapshotRoot", $snapshot) `
            -TimeoutSec $UnityImportTimeoutSec `
            -StdoutPath $overlayOut -StderrPath $overlayErr
    } | Out-Null

    foreach ($relative in @(
        "Assets\Scripts\NroE2E\NroE2EBridge.cs",
        "Assets\Scripts\NroE2E\NroE2EState.cs",
        "Assets\Scripts\NroE2E\NroE2ECommands.cs",
        "Assets\Tests\EditMode\NroE2EContractTests.cs")) {
        $sourceFile = Join-Path (Join-Path $server "tools\unity-overlay") $relative
        $targetFile = Join-Path $snapshot $relative
        New-Item -ItemType Directory -Path (Split-Path -Parent $targetFile) -Force | Out-Null
        Copy-Item -LiteralPath $sourceFile -Destination $targetFile -Force
    }

    $releaseRoot = Join-Path $runRoot "releases"
    Invoke-Phase "asset_deploy" {
        Invoke-ExternalWithTimeout -Name "ASSET-DEPLOY" -FilePath "python.exe" `
            -Arguments @(
                "-m", "assetpack", "deploy",
                "--server-artifact", $art,
                "--client-artifact", $art,
                "--release-root", $releaseRoot,
                "--snapshot", $snapshot,
                "--version-ledger", (Join-Path $server "docs\my-thuat\asset-version-ledger.json"),
                "--previous-lock", (Resolve-FullPath $PreviousLock),
                "--runtime-key-map", (Join-Path $server "docs\my-thuat\unity-runtime-key-map.json"),
                "--unity-exe", $UnityExe,
                "--protected-unity-root", $unityRoot) `
            -TimeoutSec $UnityImportTimeoutSec `
            -StdoutPath (Join-Path $runRoot "asset-deploy.stdout.log") `
            -StderrPath (Join-Path $runRoot "asset-deploy.stderr.log")
    } | Out-Null

    $testResults = Join-Path $runRoot "unity-editmode-results.xml"
    Invoke-Phase "unity_contract_tests" {
        Invoke-ExternalWithTimeout -Name "UNITY-TEST" -FilePath $UnityExe `
            -Arguments @(
                "-batchmode", "-nographics", "-quit", "-projectPath", $snapshot,
                "-executeMethod", "NroE2E.NroE2EContractTests.RunBatch",
                "-logFile", (Join-Path $runRoot "unity-tests.log")) `
            -TimeoutSec $UnityImportTimeoutSec `
            -StdoutPath (Join-Path $runRoot "unity-tests.stdout.log") `
            -StderrPath (Join-Path $runRoot "unity-tests.stderr.log")
    } | Out-Null

    $executable = Join-Path $runRoot "build\NRO-E2E.exe"
    Invoke-Phase "unity_build" {
        Invoke-ExternalWithTimeout -Name "UNITY-BUILD" -FilePath $UnityExe `
            -Arguments @(
                "-batchmode", "-nographics", "-quit", "-projectPath", $snapshot,
                "-executeMethod", "NroE2E.NroE2EBridge.BuildBatch",
                "-e2e-output", $executable,
                "-logFile", (Join-Path $runRoot "unity-build.log")) `
            -TimeoutSec $UnityBuildTimeoutSec `
            -StdoutPath (Join-Path $runRoot "unity-build.stdout.log") `
            -StderrPath (Join-Path $runRoot "unity-build.stderr.log")
    } | Out-Null

    $manifest.status = "contract_pass"
    $manifest.executable = [ordered]@{
        path = $executable
        sha256 = (Get-FileHash -LiteralPath $executable -Algorithm SHA256).Hash.ToLowerInvariant()
    }
    $manifest.durations_seconds = $script:PhaseDurations
    Write-Utf8Json -Path $manifestPath -Value $manifest
    if ($ContractOnly) {
        Write-Output "E2E-CONTRACT-PASS $manifestPath"
        exit 0
    }

    if ([string]::IsNullOrWhiteSpace($GameplayScript) -or
        -not (Test-Path -LiteralPath $GameplayScript -PathType Leaf)) {
        Stop-E2E "E2E-TASK10-NOT-RUN" `
            "full run can -GameplayScript de tao raw clean/warm evidence"
    }

    # Driver chi dieu khien Unity/server va ghi raw evidence. Quyet dinh PASS
    # phai di qua Test-Evidence.ps1 da duoc pressure-test bang bad fixtures.
    $previousTestDsn = [Environment]::GetEnvironmentVariable(
        "NRO_TEST_DSN",
        [EnvironmentVariableTarget]::Process)
    [Environment]::SetEnvironmentVariable(
        "NRO_TEST_DSN",
        $Dsn,
        [EnvironmentVariableTarget]::Process)
    try {
        Invoke-Phase "clean_warm_gameplay" {
            Invoke-ExternalWithTimeout -Name "GAMEPLAY-DRIVER" `
                -FilePath "powershell.exe" `
                -Arguments @(
                    "-NoProfile", "-ExecutionPolicy", "Bypass",
                    "-File", (Resolve-FullPath $GameplayScript),
                    "-RunRoot", $runRoot,
                    "-Executable", $executable,
                    "-Manifest", $manifestPath,
                    "-GamePort", $GamePort,
                    "-ControlPort", $ControlPort,
                    "-TelemetryPortA", $TelemetryPortA,
                    "-TelemetryPortB", $TelemetryPortB) `
                -TimeoutSec (
                    $LoginTimeoutSec +
                    (20 * $CommandTimeoutSec) +
                    (2 * $ServerStartTimeoutSec) +
                    (2 * $ServerShutdownTimeoutSec)) `
                -StdoutPath (Join-Path $runRoot "gameplay.stdout.log") `
                -StderrPath (Join-Path $runRoot "gameplay.stderr.log")
        } | Out-Null
    }
    finally {
        [Environment]::SetEnvironmentVariable(
            "NRO_TEST_DSN",
            $previousTestDsn,
            [EnvironmentVariableTarget]::Process)
    }

    Invoke-Phase "clean_warm_oracle" {
        Invoke-ExternalWithTimeout -Name "RUNNER-ORACLE" `
            -FilePath "powershell.exe" `
            -Arguments @(
                "-NoProfile", "-ExecutionPolicy", "Bypass",
                "-File", (Join-Path $server "tools\e2e\Test-Evidence.ps1"),
                "-OracleRunRoot", $runRoot) `
            -TimeoutSec $EvidenceTimeoutSec `
            -StdoutPath (Join-Path $runRoot "oracle.stdout.log") `
            -StderrPath (Join-Path $runRoot "oracle.stderr.log")
    } | Out-Null

    $manifest.status = "full_pass"
    $manifest.aggregate_result = Join-Path $runRoot "aggregate-result.json"
    $manifest.durations_seconds = $script:PhaseDurations
    Write-Utf8Json -Path $manifestPath -Value $manifest
    Write-Output "E2E-FULL-PASS $($manifest.aggregate_result)"
    exit 0
}
catch {
    $failureScope = Get-E2EFailureScope `
        -IsControlSmoke ([bool]$ControlSmokeOnly) `
        -IsContract ([bool]$ContractOnly)
    $failure = [ordered]@{
        schema_version = 1
        scope = $failureScope
        run_id = $runId
        status = "failed"
        error = $_.Exception.Message
        durations_seconds = $script:PhaseDurations
    }
    Write-Utf8Json -Path (Join-Path $runRoot "failure.json") -Value $failure
    Write-Error $_.Exception.Message
    exit 1
}
