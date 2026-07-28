param(
    [Parameter(Mandatory = $true)]
    [string]$UnityRoot,

    [Parameter(Mandatory = $true)]
    [string]$MappingPath,

    [switch]$RequireFutureOutputs
)

$ErrorActionPreference = "Stop"
$errors = New-Object System.Collections.Generic.List[string]
$allowedModes = @("direct", "slice_group", "overlay_required")
$expectedModes = @{
    "ui_logo_game_x4.png" = "direct"
    "ui_nut_thuong_x4.png" = "slice_group"
    "ui_nut_nhan_x4.png" = "slice_group"
    "ui_khung_avatar_x4.png" = "overlay_required"
    "ui_khung_hoithoai_x4.png" = "overlay_required"
    "ui_nut_khoa_x4.png" = "overlay_required"
    "ui_o_dangchon_x4.png" = "overlay_required"
    "ui_o_vatpham_x4.png" = "overlay_required"
}
$requiredSources = @(
    "ui_logo_game_x4.png",
    "ui_khung_avatar_x4.png",
    "ui_khung_hoithoai_x4.png",
    "ui_nut_thuong_x4.png",
    "ui_nut_nhan_x4.png",
    "ui_nut_khoa_x4.png",
    "ui_o_dangchon_x4.png",
    "ui_o_vatpham_x4.png"
)

function Add-ValidationError {
    param([string]$Message)
    $script:errors.Add($Message)
}

function Get-NormalizedSha256 {
    param([string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Get-SourceStem {
    param([string]$SourceAsset)
    $stem = [System.IO.Path]::GetFileNameWithoutExtension($SourceAsset).ToLowerInvariant()
    return ($stem -replace "_x[1-4]$", "")
}

function Get-KeyStem {
    param([string]$Key)
    return [System.IO.Path]::GetFileNameWithoutExtension($Key).ToLowerInvariant()
}

function Test-KeyShape {
    param(
        [string]$Key,
        [string]$SourceAsset,
        [string]$Label
    )

    if ([string]::IsNullOrWhiteSpace($Key)) {
        Add-ValidationError "$Label bị rỗng."
        return
    }
    if ($Key -notmatch "^/[A-Za-z0-9._/-]+\.png$") {
        Add-ValidationError "$Label '$Key' không phải resource key PNG tuyệt đối."
    }
    if ($Key -match "(?i)(todo|tbd|placeholder|unknown|guess|doan|sample|example)") {
        Add-ValidationError "$Label '$Key' chứa placeholder."
    }
    if ((Get-KeyStem $Key) -eq (Get-SourceStem $SourceAsset)) {
        Add-ValidationError "$Label '$Key' bị suy thẳng từ tên source '$SourceAsset', không có bằng chứng runtime."
    }
}

function Resolve-SafeSourceFile {
    param(
        [string]$RelativePath,
        [string]$Scope,
        [string]$Root
    )

    if ([string]::IsNullOrWhiteSpace($RelativePath)) {
        Add-ValidationError "Evidence thiếu đường dẫn file."
        return $null
    }

    $normalizedRelative = $RelativePath.Replace("\", "/")
    $normalizedScope = $Scope.Trim("/").Replace("\", "/")
    if (-not $normalizedRelative.StartsWith($normalizedScope + "/", [StringComparison]::OrdinalIgnoreCase)) {
        Add-ValidationError "Evidence '$RelativePath' nằm ngoài scope '$Scope'."
        return $null
    }

    $rootFull = [System.IO.Path]::GetFullPath($Root).TrimEnd("\", "/")
    $candidate = [System.IO.Path]::GetFullPath((Join-Path $rootFull $RelativePath))
    if (-not $candidate.StartsWith($rootFull + [System.IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) {
        Add-ValidationError "Evidence '$RelativePath' thoát khỏi Unity root."
        return $null
    }
    if (-not (Test-Path -LiteralPath $candidate -PathType Leaf)) {
        Add-ValidationError "Không tìm thấy evidence file '$RelativePath'."
        return $null
    }
    return $candidate
}

function Test-EvidenceSite {
    param(
        [object]$Site,
        [string]$Kind,
        [string]$Scope,
        [string]$Root
    )

    $file = Resolve-SafeSourceFile -RelativePath ([string]$Site.file) -Scope $Scope -Root $Root
    if ($null -eq $file) {
        return
    }

    $expectedSha = ([string]$Site.file_sha256).ToLowerInvariant()
    $actualSha = Get-NormalizedSha256 $file
    if ($expectedSha -notmatch "^[0-9a-f]{64}$") {
        Add-ValidationError "$Kind '$($Site.file)' có SHA không hợp lệ."
    }
    elseif ($actualSha -ne $expectedSha) {
        Add-ValidationError "$Kind '$($Site.file)' lệch SHA: map=$expectedSha source=$actualSha."
    }

    $lineNumber = 0
    if (-not [int]::TryParse(([string]$Site.line), [ref]$lineNumber) -or $lineNumber -lt 1) {
        Add-ValidationError "$Kind '$($Site.file)' có line không hợp lệ."
        return
    }
    $lines = [System.IO.File]::ReadAllLines($file)
    if ($lineNumber -gt $lines.Length) {
        Add-ValidationError "$Kind '$($Site.file):$lineNumber' vượt quá số dòng $($lines.Length)."
        return
    }
    $actualExpression = $lines[$lineNumber - 1].Trim()
    $expectedExpression = ([string]$Site.expression).Trim()
    if ($actualExpression -ne $expectedExpression) {
        Add-ValidationError "$Kind '$($Site.file):$lineNumber' lệch expression. Map='$expectedExpression' Source='$actualExpression'."
    }
}

function Test-PhysicalPaths {
    param(
        [object]$Entry,
        [string[]]$EffectiveKeys,
        [string]$Root
    )

    $rootFull = [System.IO.Path]::GetFullPath($Root).TrimEnd("\", "/")
    foreach ($zoom in 1..4) {
        $propertyName = "x$zoom"
        $property = $Entry.physical_paths.PSObject.Properties[$propertyName]
        if ($null -eq $property) {
            Add-ValidationError "'$($Entry.source_asset)' thiếu physical path $propertyName."
            continue
        }
        $paths = @($property.Value)
        if ($paths.Count -ne $EffectiveKeys.Count) {
            Add-ValidationError "'$($Entry.source_asset)' $propertyName có $($paths.Count) path nhưng cần $($EffectiveKeys.Count)."
            continue
        }
        for ($index = 0; $index -lt $paths.Count; $index++) {
            $path = ([string]$paths[$index]).Replace("\", "/")
            if ([System.IO.Path]::IsPathRooted($path)) {
                Add-ValidationError "Physical path '$path' phải là đường dẫn tương đối trong Unity root."
                continue
            }
            if (($path -split "/") -contains "..") {
                Add-ValidationError "Physical path '$path' không được chứa '..'."
                continue
            }

            $usesResizeLoader = (
                $Entry.mode -eq "direct" -and
                [string]$Entry.call_sites[0].expression -match "loadImageResize"
            )
            $physicalZoom = $(if ($usesResizeLoader) { 4 } else { $zoom })
            $expectedRelativeRoot = "Assets/Resources/res/x$physicalZoom"
            $expectedFullRoot = [System.IO.Path]::GetFullPath(
                (Join-Path $rootFull $expectedRelativeRoot)
            ).TrimEnd("\", "/")
            $physicalFile = [System.IO.Path]::GetFullPath((Join-Path $rootFull $path))
            if (-not $physicalFile.StartsWith(
                $expectedFullRoot + [System.IO.Path]::DirectorySeparatorChar,
                [StringComparison]::OrdinalIgnoreCase
            )) {
                Add-ValidationError "Physical path '$path' sau canonicalize nằm ngoài '$expectedRelativeRoot'."
                continue
            }

            $keySuffix = $EffectiveKeys[$index].TrimStart("/").Replace("mainImage/", "mainimage/")
            if (-not $path.EndsWith("/" + $keySuffix, [StringComparison]::OrdinalIgnoreCase)) {
                Add-ValidationError "Physical path '$path' không khớp key '$($EffectiveKeys[$index])'."
            }

            if ($Entry.mode -ne "overlay_required") {
                if (-not (Test-Path -LiteralPath $physicalFile -PathType Leaf)) {
                    $callExpressions = (@($Entry.call_sites) | ForEach-Object { [string]$_.expression }) -join "`n"
                    $fallbackPath = "Assets/Resources/res/x4/$keySuffix"
                    $fallbackRoot = [System.IO.Path]::GetFullPath(
                        (Join-Path $rootFull "Assets/Resources/res/x4")
                    ).TrimEnd("\", "/")
                    $fallbackFile = [System.IO.Path]::GetFullPath((Join-Path $rootFull $fallbackPath))
                    $fallbackIsContained = $fallbackFile.StartsWith(
                        $fallbackRoot + [System.IO.Path]::DirectorySeparatorChar,
                        [StringComparison]::OrdinalIgnoreCase
                    )
                    $usesLoadImageFallback = (
                        $zoom -lt 4 -and
                        $callExpressions -match "\bloadImage\(" -and
                        $fallbackIsContained -and
                        (Test-Path -LiteralPath $fallbackFile -PathType Leaf)
                    )
                    if (-not $usesLoadImageFallback) {
                        Add-ValidationError "Không tìm thấy physical file '$path' hoặc fallback x4 cho '$($Entry.source_asset)'."
                    }
                }
            }
            elseif ($RequireFutureOutputs -and -not (Test-Path -LiteralPath $physicalFile -PathType Leaf)) {
                Add-ValidationError "Không tìm thấy future physical file '$path' trong staging."
            }
        }
    }
}

try {
    if (-not (Test-Path -LiteralPath $UnityRoot -PathType Container)) {
        throw "Unity root không tồn tại: $UnityRoot"
    }
    if (-not (Test-Path -LiteralPath $MappingPath -PathType Leaf)) {
        throw "Mapping không tồn tại: $MappingPath"
    }

    $mapping = Get-Content -LiteralPath $MappingPath -Raw -Encoding UTF8 | ConvertFrom-Json
    if ($mapping.schema_version -ne 1) {
        Add-ValidationError "schema_version phải bằng 1."
    }
    if ([string]$mapping.unity_source_scope -cne "Assets/Scripts/Game1") {
        Add-ValidationError "unity_source_scope phải chính xác là 'Assets/Scripts/Game1'."
    }

    $entries = @($mapping.entries)
    if ($entries.Count -ne 8) {
        Add-ValidationError "Mapping phải có đúng 8 binding, hiện có $($entries.Count)."
    }

    $sourceNames = @{}
    $allKeys = @{}
    foreach ($entry in $entries) {
        $source = [string]$entry.source_asset
        if ($requiredSources -notcontains $source) {
            Add-ValidationError "Source asset ngoài danh sách đã duyệt: '$source'."
        }
        if ($sourceNames.ContainsKey($source.ToLowerInvariant())) {
            Add-ValidationError "Source asset bị trùng: '$source'."
        }
        else {
            $sourceNames[$source.ToLowerInvariant()] = $true
        }

        $mode = [string]$entry.mode
        if ($expectedModes.ContainsKey($source) -and $mode -ne $expectedModes[$source]) {
            Add-ValidationError "'$source' phải có mode '$($expectedModes[$source])', không phải '$mode'."
        }
        if ($allowedModes -notcontains $mode) {
            Add-ValidationError "'$source' có mode không hợp lệ '$mode'."
            continue
        }

        $logicalKeys = @($entry.logical_keys)
        $futureKey = [string]$entry.future_logical_key
        $effectiveKeys = @()
        if ($mode -eq "overlay_required") {
            if ($logicalKeys.Count -ne 0) {
                Add-ValidationError "'$source' overlay_required chưa được giả làm logical key đang tồn tại."
            }
            Test-KeyShape -Key $futureKey -SourceAsset $source -Label "future_logical_key"
            $effectiveKeys = @($futureKey)
            if (@($entry.call_sites).Count -ne 0) {
                Add-ValidationError "'$source' overlay_required không được có call-site giả."
            }
            if (@($entry.consumer_sites).Count -lt 1) {
                Add-ValidationError "'$source' overlay_required phải có consumer site + SHA precondition."
            }
        }
        else {
            if (-not [string]::IsNullOrWhiteSpace($futureKey)) {
                Add-ValidationError "'$source' mode $mode không được có future_logical_key."
            }
            if ($mode -eq "direct" -and $logicalKeys.Count -ne 1) {
                Add-ValidationError "'$source' direct phải có đúng một logical key."
            }
            if ($mode -eq "direct" -and @($entry.call_sites).Count -ne 1) {
                Add-ValidationError "'$source' direct phải có đúng một call-site."
            }
            if ($mode -eq "slice_group" -and $logicalKeys.Count -ne 3) {
                Add-ValidationError "'$source' slice_group phải có đúng ba logical key trái-giữa-phải."
            }
            if ($mode -eq "slice_group" -and @($entry.call_sites).Count -ne 3) {
                Add-ValidationError "'$source' slice_group phải có đúng ba call-site trái-giữa-phải."
            }
            if (@($entry.call_sites).Count -lt 1) {
                Add-ValidationError "'$source' mode $mode phải có call-site thật."
            }
            foreach ($key in $logicalKeys) {
                Test-KeyShape -Key ([string]$key) -SourceAsset $source -Label "logical_key"
            }
            $effectiveKeys = @($logicalKeys | ForEach-Object { [string]$_ })
        }

        foreach ($key in $effectiveKeys) {
            $keyLower = $key.ToLowerInvariant()
            if ($allKeys.ContainsKey($keyLower)) {
                Add-ValidationError "Logical key bị trùng: '$key'."
            }
            else {
                $allKeys[$keyLower] = $source
            }
        }

        foreach ($site in @($entry.call_sites)) {
            Test-EvidenceSite -Site $site -Kind "call-site" -Scope ([string]$mapping.unity_source_scope) -Root $UnityRoot
        }
        foreach ($site in @($entry.consumer_sites)) {
            Test-EvidenceSite -Site $site -Kind "consumer-site" -Scope ([string]$mapping.unity_source_scope) -Root $UnityRoot
        }

        if ($mode -ne "overlay_required") {
            $callExpressions = (@($entry.call_sites) | ForEach-Object { [string]$_.expression }) -join "`n"
            foreach ($key in $logicalKeys) {
                if (-not $callExpressions.Contains([string]$key)) {
                    Add-ValidationError "'$source' có logical key '$key' nhưng call-site không tham chiếu key đó."
                }
            }
        }

        if ($null -eq $entry.physical_paths) {
            Add-ValidationError "'$source' thiếu physical_paths."
        }
        elseif ($effectiveKeys.Count -gt 0) {
            Test-PhysicalPaths -Entry $entry -EffectiveKeys $effectiveKeys -Root $UnityRoot
        }
    }

    foreach ($requiredSource in $requiredSources) {
        if (-not $sourceNames.ContainsKey($requiredSource.ToLowerInvariant())) {
            Add-ValidationError "Thiếu source asset đã duyệt '$requiredSource'."
        }
    }

    $directCount = @($entries | Where-Object { $_.mode -eq "direct" }).Count
    $sliceCount = @($entries | Where-Object { $_.mode -eq "slice_group" }).Count
    $overlayCount = @($entries | Where-Object { $_.mode -eq "overlay_required" }).Count
    if ($directCount -ne 1 -or $sliceCount -ne 2 -or $overlayCount -ne 5) {
        Add-ValidationError "Tỷ lệ mode phải là direct=1, slice_group=2, overlay_required=5; hiện là $directCount/$sliceCount/$overlayCount."
    }

    if ($errors.Count -gt 0) {
        foreach ($validationError in $errors) {
            Write-Error $validationError
        }
        exit 1
    }

    Write-Output "OK: 8 runtime bindings hợp lệ; call-site/consumer SHA và physical paths x1..x4 khớp."
    exit 0
}
catch {
    Write-Error $_
    exit 1
}
