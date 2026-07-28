param(
    [Parameter(Mandatory = $true)]
    [string]$UnityProject,

    [Parameter(Mandatory = $true)]
    [string]$SnapshotRoot,

    [string]$ManifestPath
)

$ErrorActionPreference = "Stop"
if ([string]::IsNullOrWhiteSpace($ManifestPath)) {
    $ManifestPath = Join-Path $PSScriptRoot "overlay-manifest.json"
}

function Stop-Overlay {
    param(
        [string]$Code,
        [string]$Message
    )
    throw "$Code`: $Message"
}

function Get-NormalizedPath {
    param([string]$Path)
    return [System.IO.Path]::GetFullPath($Path).TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
}

function Test-SamePath {
    param(
        [string]$Left,
        [string]$Right
    )
    return [string]::Equals(
        (Get-NormalizedPath $Left),
        (Get-NormalizedPath $Right),
        [System.StringComparison]::OrdinalIgnoreCase
    )
}

function Test-PathInside {
    param(
        [string]$Candidate,
        [string]$Root
    )
    $normalizedCandidate = Get-NormalizedPath $Candidate
    $normalizedRoot = Get-NormalizedPath $Root
    if ([string]::Equals(
            $normalizedCandidate,
            $normalizedRoot,
            [System.StringComparison]::OrdinalIgnoreCase)) {
        return $true
    }
    return $normalizedCandidate.StartsWith(
        $normalizedRoot + [System.IO.Path]::DirectorySeparatorChar,
        [System.StringComparison]::OrdinalIgnoreCase
    )
}

function Get-Sha256 {
    param([string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Write-Utf8NoBom {
    param(
        [string]$Path,
        [string]$Content
    )
    [System.IO.File]::WriteAllText(
        $Path,
        $Content,
        [System.Text.UTF8Encoding]::new($false)
    )
}

try {
    $unityRoot = Get-NormalizedPath $UnityProject
    $snapshot = Get-NormalizedPath $SnapshotRoot
    $manifestFile = Get-NormalizedPath $ManifestPath

    if (-not (Test-Path -LiteralPath $manifestFile -PathType Leaf)) {
        Stop-Overlay "OVERLAY-MANIFEST-001" "Manifest not found: $manifestFile"
    }
    $manifest = Get-Content -LiteralPath $manifestFile -Raw -Encoding UTF8 | ConvertFrom-Json
    if ($manifest.schema_version -ne 1) {
        Stop-Overlay "OVERLAY-MANIFEST-001" "Only schema_version=1 is supported"
    }

    if (Test-PathInside $snapshot $unityRoot) {
        Stop-Overlay "OVERLAY-SOURCE-PROTECTED" "Cannot apply overlay to the Unity source"
    }
    foreach ($protectedRoot in @($manifest.protected_source_roots)) {
        if ($protectedRoot -and (Test-PathInside $snapshot $protectedRoot)) {
            Stop-Overlay "OVERLAY-SOURCE-PROTECTED" "Target is a protected Unity source"
        }
    }

    foreach ($requiredDirectory in @("Assets", "Packages", "ProjectSettings")) {
        if (-not (Test-Path -LiteralPath (Join-Path $unityRoot $requiredDirectory) -PathType Container)) {
            Stop-Overlay "OVERLAY-SOURCE-001" "Unity source is missing $requiredDirectory"
        }
    }

    $projectVersionPath = Join-Path $unityRoot "ProjectSettings\ProjectVersion.txt"
    $projectVersionText = Get-Content -LiteralPath $projectVersionPath -Raw -Encoding UTF8
    $expectedVersion = [string]$manifest.unity_version
    if ($expectedVersion -and $projectVersionText -notmatch "m_EditorVersion:\s*$([regex]::Escape($expectedVersion))(\r?\n|$)") {
        Stop-Overlay "OVERLAY-UNITY-VERSION-001" "Unity $expectedVersion is required"
    }

    if (Test-Path -LiteralPath $snapshot) {
        if (Test-Path -LiteralPath (Join-Path $snapshot ".nro-e2e-overlay.json")) {
            Stop-Overlay "OVERLAY-ALREADY-APPLIED" "Snapshot already has an overlay"
        }
        if (@(Get-ChildItem -LiteralPath $snapshot -Force).Count -gt 0) {
            Stop-Overlay "OVERLAY-TARGET-NOT-EMPTY" "Snapshot target must be empty"
        }
    }

    # Validate every patch against the source before creating the snapshot. A failed
    # precondition must never leave an unmarked, partially copied Unity tree behind.
    $manifestDirectory = Split-Path -Parent $manifestFile
    $patchPlans = [System.Collections.Generic.List[object]]::new()
    foreach ($patch in @($manifest.patches)) {
        $sourceFile = Join-Path $unityRoot ([string]$patch.target)
        if (-not (Test-Path -LiteralPath $sourceFile -PathType Leaf)) {
            Stop-Overlay "OVERLAY-SHA-001" "Source is missing: $($patch.target)"
        }
        $actualSha = Get-Sha256 $sourceFile
        if ($actualSha -ne ([string]$patch.source_sha256).ToLowerInvariant()) {
            Stop-Overlay "OVERLAY-SHA-001" "$($patch.target) SHA mismatch"
        }

        $patchFile = Join-Path $manifestDirectory ([string]$patch.patch)
        if (-not (Test-Path -LiteralPath $patchFile -PathType Leaf)) {
            Stop-Overlay "OVERLAY-PATCH-001" "Patch is missing: $($patch.patch)"
        }
        $patchSpec = Get-Content -LiteralPath $patchFile -Raw -Encoding UTF8 | ConvertFrom-Json
        $content = [System.IO.File]::ReadAllText($sourceFile)
        $lineEnding = if ($content.Contains("`r`n")) { "`r`n" } else { "`n" }
        $content = $content.Replace("`r`n", "`n")
        $replacements = @($patchSpec.replacements)
        if ($replacements.Count -eq 0 -or $null -eq $replacements[0]) {
            Stop-Overlay "OVERLAY-PATCH-001" "$($patch.patch) has no replacements"
        }
        foreach ($replacement in $replacements) {
            $before = ([string]$replacement.before).Replace("`r`n", "`n")
            $after = ([string]$replacement.after).Replace("`r`n", "`n")
            $expectedCount = if ($null -ne $replacement.expected_count) {
                [int]$replacement.expected_count
            }
            else {
                1
            }
            $count = [regex]::Matches($content, [regex]::Escape($before)).Count
            if ($count -ne $expectedCount) {
                Stop-Overlay "OVERLAY-PATCH-001" `
                    "$($patch.target): precondition must occur $expectedCount time(s), actual $count"
            }
            $content = $content.Replace($before, $after)
        }
        if ($lineEnding -eq "`r`n") {
            $content = $content.Replace("`n", "`r`n")
        }
        $patchPlans.Add([pscustomobject]@{
            target = ([string]$patch.target)
            content = $content
        })
    }

    [System.IO.Directory]::CreateDirectory($snapshot) | Out-Null
    foreach ($directory in @("Assets", "Packages", "ProjectSettings")) {
        Copy-Item -LiteralPath (Join-Path $unityRoot $directory) `
            -Destination $snapshot -Recurse
    }

    $resultFiles = [System.Collections.Generic.List[object]]::new()

    foreach ($patchPlan in $patchPlans) {
        $relativeTarget = ([string]$patchPlan.target).Replace("/", "\")
        $targetFile = Join-Path $snapshot $relativeTarget
        Write-Utf8NoBom $targetFile $patchPlan.content
        $resultFiles.Add([pscustomobject]@{
            path = ([string]$patchPlan.target).Replace("\", "/")
            sha256 = Get-Sha256 $targetFile
            kind = "patch"
        })
    }

    foreach ($overlayFileValue in @($manifest.overlay_files)) {
        $relativeOverlay = ([string]$overlayFileValue).Replace("/", "\")
        $sourceOverlay = Join-Path $PSScriptRoot $relativeOverlay
        $targetOverlay = Join-Path $snapshot $relativeOverlay
        if (-not (Test-Path -LiteralPath $sourceOverlay -PathType Leaf)) {
            Stop-Overlay "OVERLAY-FILE-001" "Overlay file is missing: $overlayFileValue"
        }
        $targetParent = Split-Path -Parent $targetOverlay
        [System.IO.Directory]::CreateDirectory($targetParent) | Out-Null
        if (Test-Path -LiteralPath $targetOverlay) {
            Stop-Overlay "OVERLAY-FILE-001" "Overlay would replace an existing file: $overlayFileValue"
        }
        Copy-Item -LiteralPath $sourceOverlay -Destination $targetOverlay
        $resultFiles.Add([pscustomobject]@{
            path = ([string]$overlayFileValue).Replace("\", "/")
            sha256 = Get-Sha256 $targetOverlay
            kind = "overlay"
        })
    }

    $marker = [ordered]@{
        schema_version = 1
        source_root = $unityRoot
        source_manifest_sha256 = Get-Sha256 $manifestFile
        unity_version = $expectedVersion
        files = @($resultFiles)
    }
    Write-Utf8NoBom (Join-Path $snapshot ".nro-e2e-overlay.json") `
        ($marker | ConvertTo-Json -Depth 8)
    Write-Host "OVERLAY-OK: $snapshot"
}
catch {
    Write-Error $_.Exception.Message
    exit 1
}
