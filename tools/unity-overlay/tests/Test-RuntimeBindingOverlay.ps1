$ErrorActionPreference = "Stop"
$script:Passed = 0

function Assert-True {
    param([bool]$Condition, [string]$Message)
    if (-not $Condition) {
        throw "ASSERT: $Message"
    }
    $script:Passed++
}

function Write-Json {
    param([string]$Path, [object]$Value)
    [System.IO.File]::WriteAllText(
        $Path,
        ($Value | ConvertTo-Json -Depth 12),
        [System.Text.UTF8Encoding]::new($false)
    )
}

function Invoke-Verify {
    param([string]$Mapping, [string]$Manifest)
    $stdout = Join-Path $fixtureRoot ("stdout-" + [guid]::NewGuid().ToString("N") + ".log")
    $stderr = Join-Path $fixtureRoot ("stderr-" + [guid]::NewGuid().ToString("N") + ".log")
    $process = Start-Process -FilePath "powershell.exe" -ArgumentList @(
        "-NoProfile",
        "-ExecutionPolicy", "Bypass",
        "-File", $verifyScript,
        "-MappingPath", $Mapping,
        "-OverlayManifestPath", $Manifest
    ) -Wait -PassThru -WindowStyle Hidden `
        -RedirectStandardOutput $stdout -RedirectStandardError $stderr
    $output = [System.IO.File]::ReadAllText($stdout) +
        [System.IO.File]::ReadAllText($stderr)
    return [pscustomobject]@{ ExitCode = $process.ExitCode; Output = $output }
}

$overlayRoot = Split-Path -Parent $PSScriptRoot
$verifyScript = Join-Path $overlayRoot "verify-runtime-bindings.ps1"
$fixtureRoot = Join-Path ([System.IO.Path]::GetTempPath()) `
    ("nro-runtime-binding-" + [guid]::NewGuid().ToString("N"))
[System.IO.Directory]::CreateDirectory($fixtureRoot) | Out-Null

try {
    $entries = @()
    for ($index = 1; $index -le 3; $index++) {
        $entries += [ordered]@{
            source_asset = "direct-$index.png"
            mode = if ($index -eq 1) { "direct" } else { "slice_group" }
            logical_keys = @("/legacy/key-$index.png")
            future_logical_key = $null
            call_sites = @()
            consumer_sites = @()
            physical_paths = @{
                x1 = @("Assets/Resources/res/x1/direct-$index.png")
                x2 = @("Assets/Resources/res/x2/direct-$index.png")
                x3 = @("Assets/Resources/res/x3/direct-$index.png")
                x4 = @("Assets/Resources/res/x4/direct-$index.png")
            }
        }
    }

    $patches = @()
    for ($index = 1; $index -le 5; $index++) {
        $file = "Assets/Scripts/Game1/Consumer$index.cs"
        $sha = $index.ToString() * 64
        $futureKey = "/naruto/ui/future-$index.png"
        $entries += [ordered]@{
            source_asset = "overlay-$index.png"
            mode = "overlay_required"
            logical_keys = @()
            future_logical_key = $futureKey
            call_sites = @()
            consumer_sites = @(
                @{
                    file = $file
                    line = 100 + $index
                    expression = "PaintLegacy$index();"
                    file_sha256 = $sha
                }
            )
            physical_paths = @{
                x1 = @("Assets/Resources/res/x1/naruto/ui/future-$index.png")
                x2 = @("Assets/Resources/res/x2/naruto/ui/future-$index.png")
                x3 = @("Assets/Resources/res/x3/naruto/ui/future-$index.png")
                x4 = @("Assets/Resources/res/x4/naruto/ui/future-$index.png")
            }
        }
        $patches += [ordered]@{
            target = $file
            source_sha256 = $sha
            patch = "patches/consumer-$index.patch"
            provides_future_keys = @($futureKey)
        }
    }

    $mapping = [ordered]@{
        schema_version = 1
        unity_source_scope = "Assets/Scripts/Game1"
        entries = $entries
    }
    $manifest = [ordered]@{
        schema_version = 1
        unity_version = "2022.3.62f2"
        patches = $patches
        overlay_files = @()
    }
    $mappingPath = Join-Path $fixtureRoot "mapping.json"
    $manifestPath = Join-Path $fixtureRoot "manifest.json"
    Write-Json $mappingPath $mapping
    Write-Json $manifestPath $manifest

    $green = Invoke-Verify $mappingPath $manifestPath
    Assert-True ($green.ExitCode -eq 0) "Valid bindings must pass: $($green.Output)"

    # Task 0's manifest also has direct/slice patches with no future-key field.
    # That absence must mean zero provided keys, not one empty key.
    $realTaskZeroShapeManifest = $manifest | ConvertTo-Json -Depth 12 | ConvertFrom-Json
    $realTaskZeroShapeManifest.patches = @(
        [pscustomobject]@{
            target = "Assets/Scripts/Game1/Rms.cs"
            source_sha256 = ("a" * 64)
            patch = "patches/rms-profile-root.patch"
        }
    ) + @($realTaskZeroShapeManifest.patches)
    $realTaskZeroShapePath = Join-Path $fixtureRoot "real-task-0-shape.json"
    Write-Json $realTaskZeroShapePath $realTaskZeroShapeManifest
    $realTaskZeroShape = Invoke-Verify $mappingPath $realTaskZeroShapePath
    Assert-True ($realTaskZeroShape.ExitCode -eq 0) `
        "A real Task 0-shaped manifest without provides_future_keys must pass: $($realTaskZeroShape.Output)"

    $nullKeyManifest = $realTaskZeroShapeManifest | ConvertTo-Json -Depth 12 | ConvertFrom-Json
    $nullKeyManifest.patches[0] | Add-Member -NotePropertyName provides_future_keys -NotePropertyValue $null
    $nullKeyPath = Join-Path $fixtureRoot "null-key.json"
    Write-Json $nullKeyPath $nullKeyManifest
    $nullKey = Invoke-Verify $mappingPath $nullKeyPath
    Assert-True ($nullKey.ExitCode -eq 0) `
        "A null provides_future_keys value must mean zero keys: $($nullKey.Output)"

    $emptyKeyManifest = $realTaskZeroShapeManifest | ConvertTo-Json -Depth 12 | ConvertFrom-Json
    $emptyKeyManifest.patches[0] | Add-Member -NotePropertyName provides_future_keys -NotePropertyValue @()
    $emptyKeyPath = Join-Path $fixtureRoot "empty-key.json"
    Write-Json $emptyKeyPath $emptyKeyManifest
    $emptyKey = Invoke-Verify $mappingPath $emptyKeyPath
    Assert-True ($emptyKey.ExitCode -eq 0) `
        "An empty provides_future_keys value must mean zero keys: $($emptyKey.Output)"

    $badShaManifest = $manifest | ConvertTo-Json -Depth 12 | ConvertFrom-Json
    $badShaManifest.patches[0].source_sha256 = "f" * 64
    $badShaPath = Join-Path $fixtureRoot "bad-sha.json"
    Write-Json $badShaPath $badShaManifest
    $badSha = Invoke-Verify $mappingPath $badShaPath
    Assert-True ($badSha.ExitCode -ne 0) "Mismatched consumer SHA must fail."
    Assert-True ($badSha.Output -match "OVERLAY-UI-CONSUMER-001") `
        "Mismatched consumer SHA must return the expected code."

    $missingKeyManifest = $manifest | ConvertTo-Json -Depth 12 | ConvertFrom-Json
    $missingKeyManifest.patches[0].provides_future_keys = @()
    $missingKeyPath = Join-Path $fixtureRoot "missing-key.json"
    Write-Json $missingKeyPath $missingKeyManifest
    $missingKey = Invoke-Verify $mappingPath $missingKeyPath
    Assert-True ($missingKey.ExitCode -ne 0) "Missing future key must fail."
    Assert-True ($missingKey.Output -match "OVERLAY-UI-CONSUMER-001") `
        "Missing future key must return the expected code."

    $duplicateKeyManifest = $manifest | ConvertTo-Json -Depth 12 | ConvertFrom-Json
    $duplicateKeyManifest.patches[1].provides_future_keys = @("/naruto/ui/future-1.png")
    $duplicateKeyPath = Join-Path $fixtureRoot "duplicate-key.json"
    Write-Json $duplicateKeyPath $duplicateKeyManifest
    $duplicateKey = Invoke-Verify $mappingPath $duplicateKeyPath
    Assert-True ($duplicateKey.ExitCode -ne 0) "Duplicate future key must fail."
    Assert-True ($duplicateKey.Output -match "OVERLAY-UI-CONSUMER-001") `
        "Duplicate future key must return the expected code."

    Write-Host "PASS: $script:Passed assertions"
}
finally {
    if (Test-Path -LiteralPath $fixtureRoot) {
        Remove-Item -LiteralPath $fixtureRoot -Recurse -Force
    }
}
