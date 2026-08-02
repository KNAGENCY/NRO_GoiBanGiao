param(
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

function Get-Sha256 {
    param([string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

try {
    $snapshot = [System.IO.Path]::GetFullPath($SnapshotRoot)
    $manifestFile = [System.IO.Path]::GetFullPath($ManifestPath)
    $markerPath = Join-Path $snapshot ".nro-e2e-overlay.json"
    if (-not (Test-Path -LiteralPath $markerPath -PathType Leaf)) {
        Stop-Overlay "OVERLAY-MARKER-001" "Snapshot has no overlay marker"
    }
    $marker = Get-Content -LiteralPath $markerPath -Raw -Encoding UTF8 | ConvertFrom-Json
    if ($marker.schema_version -ne 1) {
        Stop-Overlay "OVERLAY-MARKER-001" "Marker schema mismatch"
    }
    if ((Get-Sha256 $manifestFile) -ne ([string]$marker.source_manifest_sha256).ToLowerInvariant()) {
        Stop-Overlay "OVERLAY-MANIFEST-CHANGED" "Manifest changed after overlay application"
    }
    foreach ($file in @($marker.files)) {
        $path = Join-Path $snapshot (([string]$file.path).Replace("/", "\"))
        if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
            Stop-Overlay "OVERLAY-VERIFY-001" "Missing $($file.path)"
        }
        if ((Get-Sha256 $path) -ne ([string]$file.sha256).ToLowerInvariant()) {
            Stop-Overlay "OVERLAY-VERIFY-001" "$($file.path) was modified"
        }
    }
    Write-Host "OVERLAY-VERIFY-OK: $snapshot"
}
catch {
    Write-Error $_.Exception.Message
    exit 1
}
