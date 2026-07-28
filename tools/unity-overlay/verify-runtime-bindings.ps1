param(
    [Parameter(Mandatory = $true)]
    [string]$MappingPath,

    [string]$OverlayManifestPath
)

$ErrorActionPreference = "Stop"
if ([string]::IsNullOrWhiteSpace($OverlayManifestPath)) {
    $OverlayManifestPath = Join-Path $PSScriptRoot "overlay-manifest.json"
}

function Stop-BindingCheck {
    param([string]$Message)
    throw "OVERLAY-UI-CONSUMER-001: $Message"
}

function Get-ProvidedFutureKeys {
    param([object]$Patch)

    $property = $Patch.PSObject.Properties["provides_future_keys"]
    if ($null -eq $property -or $null -eq $property.Value) {
        return @()
    }

    return @($property.Value | Where-Object {
        -not [string]::IsNullOrWhiteSpace([string]$_)
    })
}

try {
    $mapping = Get-Content -LiteralPath $MappingPath -Raw -Encoding UTF8 |
        ConvertFrom-Json
    $manifest = Get-Content -LiteralPath $OverlayManifestPath -Raw -Encoding UTF8 |
        ConvertFrom-Json
    if ($mapping.schema_version -ne 1 -or $manifest.schema_version -ne 1) {
        Stop-BindingCheck "Both inputs must use schema_version=1"
    }

    $entries = @($mapping.entries)
    $overlayEntries = @($entries | Where-Object { $_.mode -eq "overlay_required" })
    if ($entries.Count -ne 8 -or $overlayEntries.Count -ne 5) {
        Stop-BindingCheck "Expected eight bindings with five overlay_required entries"
    }

    $providedKeys = [System.Collections.Generic.HashSet[string]]::new(
        [System.StringComparer]::Ordinal
    )
    foreach ($patch in @($manifest.patches)) {
        foreach ($providedKey in @(Get-ProvidedFutureKeys $patch)) {
            if (-not $providedKeys.Add([string]$providedKey)) {
                Stop-BindingCheck "Future key is provided more than once: $providedKey"
            }
        }
    }

    $expectedKeys = [System.Collections.Generic.HashSet[string]]::new(
        [System.StringComparer]::Ordinal
    )
    foreach ($entry in $overlayEntries) {
        $futureKey = [string]$entry.future_logical_key
        if ([string]::IsNullOrWhiteSpace($futureKey) -or
            @($entry.logical_keys).Count -ne 0) {
            Stop-BindingCheck "overlay_required binding has an invalid future key contract"
        }
        if (-not $expectedKeys.Add($futureKey)) {
            Stop-BindingCheck "Duplicate future key in runtime map: $futureKey"
        }

        foreach ($zoomName in @("x1", "x2", "x3", "x4")) {
            if (@($entry.physical_paths.$zoomName).Count -eq 0) {
                Stop-BindingCheck "$futureKey is missing physical path $zoomName"
            }
        }

        $consumerSites = @($entry.consumer_sites)
        if ($consumerSites.Count -eq 0) {
            Stop-BindingCheck "$futureKey has no consumer site"
        }
        foreach ($consumer in $consumerSites) {
            $matches = @(
                $manifest.patches | Where-Object {
                    ([string]$_.target) -eq ([string]$consumer.file) -and
                    ([string]$_.source_sha256).ToLowerInvariant() -eq
                        ([string]$consumer.file_sha256).ToLowerInvariant() -and
                    @(Get-ProvidedFutureKeys $_) -contains $futureKey
                }
            )
            if ($matches.Count -ne 1) {
                Stop-BindingCheck (
                    "$futureKey consumer $($consumer.file) must match exactly one " +
                    "SHA-pinned patch"
                )
            }
        }
    }

    if (-not $providedKeys.SetEquals($expectedKeys)) {
        Stop-BindingCheck "Manifest future keys differ from the runtime key map"
    }

    Write-Host "OVERLAY-UI-BINDINGS-OK: 5/5"
}
catch {
    Write-Error $_.Exception.Message
    exit 1
}
