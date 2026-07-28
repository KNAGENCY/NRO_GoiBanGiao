param(
    [string]$UnityExe = "D:\Unity\2022.3.62f2\Editor\Unity.exe",
    [switch]$KeepFixture
)

$ErrorActionPreference = "Stop"
$script:Passed = 0

function Assert-True {
    param([bool]$Condition, [string]$Message)
    if (-not $Condition) {
        throw "ASSERT: $Message"
    }
    $script:Passed++
}

function Write-Utf8NoBom {
    param([string]$Path, [string]$Content)
    [System.IO.Directory]::CreateDirectory((Split-Path -Parent $Path)) | Out-Null
    [System.IO.File]::WriteAllText(
        $Path,
        $Content,
        [System.Text.UTF8Encoding]::new($false)
    )
}

function New-Png {
    param([string]$Path, [int]$Width, [int]$Height)
    [System.IO.Directory]::CreateDirectory((Split-Path -Parent $Path)) | Out-Null
    $bitmap = [System.Drawing.Bitmap]::new($Width, $Height)
    try {
        $bitmap.SetPixel(0, 0, [System.Drawing.Color]::Magenta)
        $bitmap.Save($Path, [System.Drawing.Imaging.ImageFormat]::Png)
    }
    finally {
        $bitmap.Dispose()
    }
}

function Invoke-UnityPreflight {
    param([string]$ProjectRoot, [string]$RuntimeMap, [string]$Contract, [string]$Log)
    $unityArguments = @(
        "-batchmode",
        "-nographics",
        "-projectPath", $ProjectRoot,
        "-executeMethod", "NroE2E.NroResourcePreflightTests.RunBatch",
        "-runtime-key-map", $RuntimeMap,
        "-resource-contract", $Contract,
        "-logFile", $Log
    )
    $quotedArguments = $unityArguments | ForEach-Object {
        '"' + ([string]$_).Replace('"', '\"') + '"'
    }
    $startInfo = [System.Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = $UnityExe
    $startInfo.Arguments = $quotedArguments -join " "
    $startInfo.UseShellExecute = $false
    $startInfo.CreateNoWindow = $true
    $process = [System.Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    $process.Start() | Out-Null
    if (-not $process.WaitForExit(180000)) {
        $process.Kill()
        throw "UNITY-PREFLIGHT-TIMEOUT: Unity did not exit within 180 seconds."
    }
    $unityExitCode = $process.ExitCode
    $process.Dispose()
    return [pscustomobject]@{
        ExitCode = $unityExitCode
        Log = if (Test-Path -LiteralPath $Log) {
            [System.IO.File]::ReadAllText($Log)
        }
        else {
            ""
        }
    }
}

Assert-True (Test-Path -LiteralPath $UnityExe -PathType Leaf) "Unity 2022.3.62f2 was not found."
Add-Type -AssemblyName System.Drawing

$overlayRoot = Split-Path -Parent $PSScriptRoot
$fixtureRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("nro-unity-preflight-" + [guid]::NewGuid().ToString("N"))
$assetsRoot = Join-Path $fixtureRoot "Assets"
$resourcesRoot = Join-Path $assetsRoot "Resources"
$runtimeMapPath = Join-Path $fixtureRoot "runtime-key-map.json"
$contractPath = Join-Path $fixtureRoot "resource-contract.json"

try {
    Write-Utf8NoBom (Join-Path $fixtureRoot "Packages\manifest.json") `
        '{"dependencies":{"com.unity.modules.jsonserialize":"1.0.0"}}'
    Write-Utf8NoBom (Join-Path $fixtureRoot "ProjectSettings\ProjectVersion.txt") @'
m_EditorVersion: 2022.3.62f2
m_EditorVersionWithRevision: 2022.3.62f2 (7670c08855a9)
'@
    Write-Utf8NoBom (Join-Path $assetsRoot "Scripts\Game1\Main.cs") @'
namespace Game1
{
    public static class Main
    {
        public static int typeClient;
    }
}
'@
    [System.IO.Directory]::CreateDirectory((Join-Path $assetsRoot "Scripts\NroE2E")) | Out-Null
    [System.IO.Directory]::CreateDirectory((Join-Path $assetsRoot "Tests\EditMode")) | Out-Null
    [System.IO.Directory]::CreateDirectory((Join-Path $assetsRoot "Editor")) | Out-Null
    Copy-Item -LiteralPath (Join-Path $overlayRoot "Assets\Scripts\NroE2E\NroE2EConfig.cs") `
        -Destination (Join-Path $assetsRoot "Scripts\NroE2E\NroE2EConfig.cs")
    Copy-Item -LiteralPath (Join-Path $overlayRoot "Assets\Tests\EditMode\NroResourcePreflightTests.cs") `
        -Destination (Join-Path $assetsRoot "Tests\EditMode\NroResourcePreflightTests.cs")
    Copy-Item -LiteralPath (Join-Path $overlayRoot "Assets\Editor\NroResourceImportSettings.cs") `
        -Destination (Join-Path $assetsRoot "Editor\NroResourceImportSettings.cs")

    $runtimeEntries = @()
    for ($entryIndex = 1; $entryIndex -le 8; $entryIndex++) {
        $paths = [ordered]@{ x1 = @(); x2 = @(); x3 = @(); x4 = @() }
        for ($zoom = 1; $zoom -le 4; $zoom++) {
            $relative = "Assets/Resources/res/x$zoom/ui/runtime-$entryIndex.png"
            $paths["x$zoom"] = @($relative)
            New-Png (Join-Path $fixtureRoot ($relative.Replace("/", "\"))) 8 8
        }
        $runtimeEntries += [ordered]@{
            source_asset = "ui_source_$entryIndex.png"
            mode = "direct"
            logical_keys = @("/ui/runtime-$entryIndex")
            future_logical_key = $null
            call_sites = @()
            consumer_sites = @()
            physical_paths = $paths
        }
    }
    $runtimeMap = [ordered]@{
        schema_version = 1
        unity_source_scope = "Assets/Scripts/Game1"
        entries = $runtimeEntries
    }
    Write-Utf8NoBom $runtimeMapPath ($runtimeMap | ConvertTo-Json -Depth 10)

    $resources = @()
    foreach ($kind in @("small", "background", "hud")) {
        $zooms = @()
        for ($zoom = 1; $zoom -le 4; $zoom++) {
            $relative = "Assets/Resources/res/x$zoom/preflight/$kind.png"
            $width = 7 * $zoom
            $height = 5 * $zoom
            New-Png (Join-Path $fixtureRoot ($relative.Replace("/", "\"))) $width $height
            $zooms += [ordered]@{
                zoom = $zoom
                physical_path = $relative
                expected_width = $width
                expected_height = $height
            }
        }
        $resources += [ordered]@{
            kind = $kind
            name = "fixture-$kind"
            zooms = $zooms
        }
    }
    for ($zoom = 1; $zoom -le 4; $zoom++) {
        for ($tile = 1; $tile -le 2; $tile++) {
            if ($zoom -eq 1) {
                $relative = "Assets/Resources/res/x1/t/1/t_$($tile.ToString('00')).png"
            }
            else {
                $relative = "Assets/Resources/res/x$zoom/t/1`$$tile.png"
            }
            New-Png (Join-Path $fixtureRoot ($relative.Replace("/", "\"))) (24 * $zoom) (24 * $zoom)
        }
    }
    $contract = [ordered]@{
        schema_version = 1
        tile_id = 1
        tile_count = 2
        resources = $resources
    }
    Write-Utf8NoBom $contractPath ($contract | ConvertTo-Json -Depth 10)

    $greenLog = Join-Path $fixtureRoot "preflight-green.log"
    $green = Invoke-UnityPreflight $fixtureRoot $runtimeMapPath $contractPath $greenLog
    Assert-True ($green.ExitCode -eq 0) "Unity preflight fixture must pass. Log: $($green.Log)"
    Assert-True ($green.Log -match "RESOURCE-PREFLIGHT-OK") "Passing log lacks success marker."

    $mutatedTile = Join-Path $resourcesRoot "res\x3\t\1`$1.png"
    Remove-Item -LiteralPath $mutatedTile -Force
    if (Test-Path -LiteralPath ($mutatedTile + ".meta")) {
        Remove-Item -LiteralPath ($mutatedTile + ".meta") -Force
    }
    $redLog = Join-Path $fixtureRoot "preflight-x3-red.log"
    $red = Invoke-UnityPreflight $fixtureRoot $runtimeMapPath $contractPath $redLog
    Assert-True ($red.ExitCode -ne 0) "Removing an x3 tile must fail Unity batch."
    Assert-True ($red.Log -match "ASSET-ZOOM-003") `
        "The x3 mutation must fail with ASSET-ZOOM-003. Log: $($red.Log)"

    Write-Host "PASS: $script:Passed assertions"
}
finally {
    if ($KeepFixture) {
        Write-Host "Fixture: $fixtureRoot"
    }
    elseif (Test-Path -LiteralPath $fixtureRoot) {
        Remove-Item -LiteralPath $fixtureRoot -Recurse -Force
    }
}
