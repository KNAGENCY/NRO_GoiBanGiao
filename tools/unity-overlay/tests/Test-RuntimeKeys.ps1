param(
    [string]$UnityRoot = "D:\NRO_Project",
    [string]$MappingPath = ""
)

$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..\..")).Path
$discoverScript = Join-Path $repoRoot "tools\unity-overlay\discover-runtime-keys.ps1"

function Assert-True {
    param(
        [bool]$Condition,
        [string]$Message
    )

    if (-not $Condition) {
        throw "ASSERT FAIL: $Message"
    }
}

function Get-Sha256 {
    param([string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Invoke-Discovery {
    param(
        [string]$Root,
        [string]$Map,
        [switch]$RequireFutureOutputs
    )

    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $arguments = @(
        "-NoProfile",
        "-ExecutionPolicy", "Bypass",
        "-File", $discoverScript,
        "-UnityRoot", $Root,
        "-MappingPath", $Map
    )
    if ($RequireFutureOutputs) {
        $arguments += "-RequireFutureOutputs"
    }
    $output = & powershell.exe @arguments 2>&1
    $ErrorActionPreference = $oldPreference
    return @{
        ExitCode = $LASTEXITCODE
        Output = ($output -join [Environment]::NewLine)
    }
}

function Copy-JsonObject {
    param([object]$Value)
    return (($Value | ConvertTo-Json -Depth 20) | ConvertFrom-Json)
}

$tempRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("nro-runtime-keys-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempRoot | Out-Null

try {
    $game1Dir = Join-Path $tempRoot "Assets\Scripts\Game1"
    New-Item -ItemType Directory -Path $game1Dir | Out-Null

    $sourceLines = @(
        'logo = GameCanvas.loadImageResize("/Logo/1.png");',
        'normalLeft = GameCanvas.loadImage("/mainImage/btn0left.png");',
        'normalMid = GameCanvas.loadImage("/mainImage/btn0mid.png");',
        'normalRight = GameCanvas.loadImage("/mainImage/btn0right.png");',
        'pressedLeft = GameCanvas.loadImage("/mainImage/btn1left.png");',
        'pressedMid = GameCanvas.loadImage("/mainImage/btn1mid.png");',
        'pressedRight = GameCanvas.loadImage("/mainImage/btn1right.png");',
        'DrawDialogue(g);',
        'DrawAvatar(g);',
        'DrawDisabledButton(g);',
        'DrawItemSlot(g);',
        'DrawSelectedSlot(g);'
    )
    $sourceFile = Join-Path $game1Dir "FixtureUi.cs"
    [System.IO.File]::WriteAllLines($sourceFile, $sourceLines, (New-Object System.Text.UTF8Encoding($false)))
    $sourceSha = Get-Sha256 $sourceFile

    $physicalForLogo = @{
        x1 = @("Assets/Resources/res/x4/Logo/1.png")
        x2 = @("Assets/Resources/res/x4/Logo/1.png")
        x3 = @("Assets/Resources/res/x4/Logo/1.png")
        x4 = @("Assets/Resources/res/x4/Logo/1.png")
    }
    $physicalForNormal = @{
        x1 = @("Assets/Resources/res/x1/mainimage/btn0left.png", "Assets/Resources/res/x1/mainimage/btn0mid.png", "Assets/Resources/res/x1/mainimage/btn0right.png")
        x2 = @("Assets/Resources/res/x2/mainimage/btn0left.png", "Assets/Resources/res/x2/mainimage/btn0mid.png", "Assets/Resources/res/x2/mainimage/btn0right.png")
        x3 = @("Assets/Resources/res/x3/mainimage/btn0left.png", "Assets/Resources/res/x3/mainimage/btn0mid.png", "Assets/Resources/res/x3/mainimage/btn0right.png")
        x4 = @("Assets/Resources/res/x4/mainimage/btn0left.png", "Assets/Resources/res/x4/mainimage/btn0mid.png", "Assets/Resources/res/x4/mainimage/btn0right.png")
    }
    $physicalForPressed = @{
        x1 = @("Assets/Resources/res/x1/mainimage/btn1left.png", "Assets/Resources/res/x1/mainimage/btn1mid.png", "Assets/Resources/res/x1/mainimage/btn1right.png")
        x2 = @("Assets/Resources/res/x2/mainimage/btn1left.png", "Assets/Resources/res/x2/mainimage/btn1mid.png", "Assets/Resources/res/x2/mainimage/btn1right.png")
        x3 = @("Assets/Resources/res/x3/mainimage/btn1left.png", "Assets/Resources/res/x3/mainimage/btn1mid.png", "Assets/Resources/res/x3/mainimage/btn1right.png")
        x4 = @("Assets/Resources/res/x4/mainimage/btn1left.png", "Assets/Resources/res/x4/mainimage/btn1mid.png", "Assets/Resources/res/x4/mainimage/btn1right.png")
    }

    $entries = @(
        @{
            source_asset = "ui_logo_game_x4.png"
            mode = "direct"
            logical_keys = @("/Logo/1.png")
            future_logical_key = $null
            call_sites = @(@{
                file = "Assets/Scripts/Game1/FixtureUi.cs"
                line = 1
                expression = $sourceLines[0]
                file_sha256 = $sourceSha
            })
            consumer_sites = @()
            physical_paths = $physicalForLogo
        },
        @{
            source_asset = "ui_nut_thuong_x4.png"
            mode = "slice_group"
            logical_keys = @("/mainImage/btn0left.png", "/mainImage/btn0mid.png", "/mainImage/btn0right.png")
            future_logical_key = $null
            call_sites = @(
                @{
                    file = "Assets/Scripts/Game1/FixtureUi.cs"
                    line = 2
                    expression = $sourceLines[1]
                    file_sha256 = $sourceSha
                },
                @{
                    file = "Assets/Scripts/Game1/FixtureUi.cs"
                    line = 3
                    expression = $sourceLines[2]
                    file_sha256 = $sourceSha
                },
                @{
                    file = "Assets/Scripts/Game1/FixtureUi.cs"
                    line = 4
                    expression = $sourceLines[3]
                    file_sha256 = $sourceSha
                }
            )
            consumer_sites = @()
            physical_paths = $physicalForNormal
        },
        @{
            source_asset = "ui_nut_nhan_x4.png"
            mode = "slice_group"
            logical_keys = @("/mainImage/btn1left.png", "/mainImage/btn1mid.png", "/mainImage/btn1right.png")
            future_logical_key = $null
            call_sites = @(
                @{
                    file = "Assets/Scripts/Game1/FixtureUi.cs"
                    line = 5
                    expression = $sourceLines[4]
                    file_sha256 = $sourceSha
                },
                @{
                    file = "Assets/Scripts/Game1/FixtureUi.cs"
                    line = 6
                    expression = $sourceLines[5]
                    file_sha256 = $sourceSha
                },
                @{
                    file = "Assets/Scripts/Game1/FixtureUi.cs"
                    line = 7
                    expression = $sourceLines[6]
                    file_sha256 = $sourceSha
                }
            )
            consumer_sites = @()
            physical_paths = $physicalForPressed
        }
    )

    $overlaySpecs = @(
        @("ui_khung_hoithoai_x4.png", "/naruto/ui/dialog-frame.png", 8),
        @("ui_khung_avatar_x4.png", "/naruto/ui/avatar-frame.png", 9),
        @("ui_nut_khoa_x4.png", "/naruto/ui/button-disabled.png", 10),
        @("ui_o_vatpham_x4.png", "/naruto/ui/item-slot.png", 11),
        @("ui_o_dangchon_x4.png", "/naruto/ui/item-slot-selected.png", 12)
    )
    foreach ($spec in $overlaySpecs) {
        $futureKey = [string]$spec[1]
        $line = [int]$spec[2]
        $leaf = [System.IO.Path]::GetFileNameWithoutExtension($futureKey)
        $physical = @{}
        foreach ($zoom in 1..4) {
            $physical["x$zoom"] = @("Assets/Resources/res/x$zoom/naruto/ui/$leaf.png")
        }
        $entries += @{
            source_asset = [string]$spec[0]
            mode = "overlay_required"
            logical_keys = @()
            future_logical_key = $futureKey
            call_sites = @()
            consumer_sites = @(@{
                file = "Assets/Scripts/Game1/FixtureUi.cs"
                line = $line
                expression = $sourceLines[$line - 1]
                file_sha256 = $sourceSha
            })
            physical_paths = $physical
        }
    }

    $mapping = @{
        schema_version = 1
        unity_source_scope = "Assets/Scripts/Game1"
        entries = $entries
    }
    $existingPhysicalPaths = @(
        $physicalForLogo.x4 +
        $physicalForNormal.x4 +
        $physicalForPressed.x4
    ) | Select-Object -Unique
    foreach ($overlayEntry in @($entries | Where-Object { $_.mode -eq "overlay_required" })) {
        foreach ($zoom in 1..4) {
            $existingPhysicalPaths += @($overlayEntry.physical_paths["x$zoom"])
        }
    }
    $existingPhysicalPaths = @($existingPhysicalPaths | Select-Object -Unique)
    foreach ($relativePhysicalPath in $existingPhysicalPaths) {
        $fullPhysicalPath = Join-Path $tempRoot $relativePhysicalPath
        New-Item -ItemType Directory -Path (Split-Path $fullPhysicalPath -Parent) -Force | Out-Null
        [System.IO.File]::WriteAllBytes($fullPhysicalPath, [byte[]](1, 2, 3))
    }
    $validMap = Join-Path $tempRoot "valid.json"
    [System.IO.File]::WriteAllText(
        $validMap,
        ($mapping | ConvertTo-Json -Depth 20),
        (New-Object System.Text.UTF8Encoding($false))
    )

    # Break được bảo vệ: script discovery chưa tồn tại hoặc không chấp nhận contract hợp lệ.
    $validResult = Invoke-Discovery -Root $tempRoot -Map $validMap
    Assert-True ($validResult.ExitCode -eq 0) "contract 8 binding hợp lệ phải qua. Output: $($validResult.Output)"

    # Break được bảo vệ: staging có đủ cả 5 output tương lai nhưng validator không kiểm được.
    $validStagingResult = Invoke-Discovery -Root $tempRoot -Map $validMap -RequireFutureOutputs
    Assert-True ($validStagingResult.ExitCode -eq 0) "staging đủ future output phải qua. Output: $($validStagingResult.Output)"

    # Contract: direct/slice use call-sites; only overlay has a consumer precondition.
    $directWithConsumer = Copy-JsonObject $mapping
    $directWithConsumer.entries[0].consumer_sites = @($mapping.entries[3].consumer_sites[0])
    $directWithConsumerPath = Join-Path $tempRoot "direct-with-consumer.json"
    [System.IO.File]::WriteAllText($directWithConsumerPath, ($directWithConsumer | ConvertTo-Json -Depth 20))
    $directWithConsumerResult = Invoke-Discovery -Root $tempRoot -Map $directWithConsumerPath
    Assert-True ($directWithConsumerResult.ExitCode -ne 0) "direct must not carry an overlay consumer precondition"
    Assert-True ($directWithConsumerResult.Output -match "RUNTIME-BINDING-001") "direct consumer error must state the contract distinction"

    # Mutation: removing an overlay consumer site removes Task 1's required precondition.
    $overlayWithoutConsumer = Copy-JsonObject $mapping
    $overlayWithoutConsumer.entries[3].consumer_sites = @()
    $overlayWithoutConsumerPath = Join-Path $tempRoot "overlay-without-consumer.json"
    [System.IO.File]::WriteAllText($overlayWithoutConsumerPath, ($overlayWithoutConsumer | ConvertTo-Json -Depth 20))
    $overlayWithoutConsumerResult = Invoke-Discovery -Root $tempRoot -Map $overlayWithoutConsumerPath
    Assert-True ($overlayWithoutConsumerResult.ExitCode -ne 0) "overlay without a consumer site must fail"
    Assert-True ($overlayWithoutConsumerResult.Output -match "RUNTIME-OVERLAY-001") "overlay consumer failure must name the precondition"

    # Mutation: an overlay consumer without SHA is not a valid precondition.
    $overlayWithoutSha = Copy-JsonObject $mapping
    $overlayWithoutSha.entries[3].consumer_sites[0].PSObject.Properties.Remove("file_sha256")
    $overlayWithoutShaPath = Join-Path $tempRoot "overlay-without-consumer-sha.json"
    [System.IO.File]::WriteAllText($overlayWithoutShaPath, ($overlayWithoutSha | ConvertTo-Json -Depth 20))
    $overlayWithoutShaResult = Invoke-Discovery -Root $tempRoot -Map $overlayWithoutShaPath
    Assert-True ($overlayWithoutShaResult.ExitCode -ne 0) "overlay without consumer SHA must fail"
    Assert-True ($overlayWithoutShaResult.Output -match "RUNTIME-EVIDENCE-002") "overlay SHA failure must name the consumer precondition"

    # Break được bảo vệ: một source đổi mode làm sai tỷ lệ 1 direct / 2 slice / 5 overlay.
    $wrongPressedMode = Copy-JsonObject $mapping
    $wrongPressedMode.entries[2].mode = "overlay_required"
    $wrongPressedModePath = Join-Path $tempRoot "wrong-pressed-mode.json"
    [System.IO.File]::WriteAllText($wrongPressedModePath, ($wrongPressedMode | ConvertTo-Json -Depth 20))
    $wrongPressedModeResult = Invoke-Discovery -Root $tempRoot -Map $wrongPressedModePath
    Assert-True ($wrongPressedModeResult.ExitCode -ne 0) "nút nhấn đổi từ slice_group sang overlay phải bị chặn"
    Assert-True ($wrongPressedModeResult.Output -match "mode.*slice_group") "lỗi mode phải nêu mode bắt buộc của source"

    # Break được bảo vệ: scope bị nới rộng ngoài đúng thư mục Game1.
    $wrongScope = Copy-JsonObject $mapping
    $wrongScope.unity_source_scope = "Assets/Scripts"
    $wrongScopePath = Join-Path $tempRoot "wrong-scope.json"
    [System.IO.File]::WriteAllText($wrongScopePath, ($wrongScope | ConvertTo-Json -Depth 20))
    $wrongScopeResult = Invoke-Discovery -Root $tempRoot -Map $wrongScopePath
    Assert-True ($wrongScopeResult.ExitCode -ne 0) "unity_source_scope rộng hơn Game1 phải bị chặn"
    Assert-True ($wrongScopeResult.Output -match "Assets/Scripts/Game1") "lỗi scope phải nêu scope bắt buộc"

    # Break được bảo vệ: evidence SHA cũ vẫn được chấp nhận sau khi consumer thay đổi.
    $badSha = Copy-JsonObject $mapping
    $badSha.entries[0].call_sites[0].file_sha256 = ("0" * 64)
    $badShaPath = Join-Path $tempRoot "bad-sha.json"
    [System.IO.File]::WriteAllText($badShaPath, ($badSha | ConvertTo-Json -Depth 20))
    $badShaResult = Invoke-Discovery -Root $tempRoot -Map $badShaPath
    Assert-True ($badShaResult.ExitCode -ne 0) "SHA call-site sai phải bị chặn"
    Assert-True ($badShaResult.Output -match "SHA") "lỗi SHA phải nói rõ nguyên nhân"

    # Break được bảo vệ: đổi key sang tên đoán nhưng call-site thật không dùng key đó.
    $guessedKey = Copy-JsonObject $mapping
    $guessedKey.entries[0].logical_keys[0] = "/Logo/naruto-emblem.png"
    $guessedKeyPath = Join-Path $tempRoot "guessed-key.json"
    [System.IO.File]::WriteAllText($guessedKeyPath, ($guessedKey | ConvertTo-Json -Depth 20))
    $guessedKeyResult = Invoke-Discovery -Root $tempRoot -Map $guessedKeyPath
    Assert-True ($guessedKeyResult.ExitCode -ne 0) "key đoán không có trong call-site phải bị chặn"
    Assert-True ($guessedKeyResult.Output -match "call-site") "lỗi key đoán phải chỉ ra call-site"

    # Break được bảo vệ: placeholder hoặc source stem bị dùng làm key tương lai.
    $placeholder = Copy-JsonObject $mapping
    $placeholder.entries[3].future_logical_key = "TODO"
    $placeholderPath = Join-Path $tempRoot "placeholder.json"
    [System.IO.File]::WriteAllText($placeholderPath, ($placeholder | ConvertTo-Json -Depth 20))
    $placeholderResult = Invoke-Discovery -Root $tempRoot -Map $placeholderPath
    Assert-True ($placeholderResult.ExitCode -ne 0) "placeholder phải bị chặn"

    $inferred = Copy-JsonObject $mapping
    $inferred.entries[3].future_logical_key = "/naruto/ui/ui_khung_hoithoai_x4.png"
    $inferredPath = Join-Path $tempRoot "inferred.json"
    [System.IO.File]::WriteAllText($inferredPath, ($inferred | ConvertTo-Json -Depth 20))
    $inferredResult = Invoke-Discovery -Root $tempRoot -Map $inferredPath
    Assert-True ($inferredResult.ExitCode -ne 0) "key suy thẳng từ tên source phải bị chặn"

    # Break được bảo vệ: trùng key hoặc thiếu một quy tắc zoom lọt qua.
    $duplicate = Copy-JsonObject $mapping
    $duplicate.entries[4].future_logical_key = $duplicate.entries[3].future_logical_key
    $duplicatePath = Join-Path $tempRoot "duplicate.json"
    [System.IO.File]::WriteAllText($duplicatePath, ($duplicate | ConvertTo-Json -Depth 20))
    $duplicateResult = Invoke-Discovery -Root $tempRoot -Map $duplicatePath
    Assert-True ($duplicateResult.ExitCode -ne 0) "logical key trùng phải bị chặn"

    $missingZoom = Copy-JsonObject $mapping
    $missingZoom.entries[0].physical_paths.PSObject.Properties.Remove("x3")
    $missingZoomPath = Join-Path $tempRoot "missing-zoom.json"
    [System.IO.File]::WriteAllText($missingZoomPath, ($missingZoom | ConvertTo-Json -Depth 20))
    $missingZoomResult = Invoke-Discovery -Root $tempRoot -Map $missingZoomPath
    Assert-True ($missingZoomResult.ExitCode -ne 0) "thiếu physical path x3 phải bị chặn"

    # Break được bảo vệ: `..` thoát/đi vòng trong physical path dù chuỗi vẫn có prefix hợp lệ.
    $pathTraversal = Copy-JsonObject $mapping
    $pathTraversal.entries[3].physical_paths.x1[0] = "Assets/Resources/res/x1/naruto/ui/../ui/dialog-frame.png"
    $pathTraversalPath = Join-Path $tempRoot "path-traversal.json"
    [System.IO.File]::WriteAllText($pathTraversalPath, ($pathTraversal | ConvertTo-Json -Depth 20))
    $pathTraversalResult = Invoke-Discovery -Root $tempRoot -Map $pathTraversalPath
    Assert-True ($pathTraversalResult.ExitCode -ne 0) "physical path chứa '..' phải bị chặn"
    Assert-True ($pathTraversalResult.Output -match "\.\.") "lỗi traversal phải nêu rõ '..'"

    # Break được bảo vệ: prefix gần giống x1 nhưng canonical path không nằm trong x1 thật.
    $outsideZoomRoot = Copy-JsonObject $mapping
    $outsideZoomRoot.entries[3].physical_paths.x1[0] = "Assets/Resources/res/x1evil/naruto/ui/dialog-frame.png"
    $outsideZoomRootPath = Join-Path $tempRoot "outside-zoom-root.json"
    [System.IO.File]::WriteAllText($outsideZoomRootPath, ($outsideZoomRoot | ConvertTo-Json -Depth 20))
    $outsideZoomRootResult = Invoke-Discovery -Root $tempRoot -Map $outsideZoomRootPath
    Assert-True ($outsideZoomRootResult.ExitCode -ne 0) "physical path ngoài đúng root res/xN phải bị chặn"
    Assert-True ($outsideZoomRootResult.Output -match "canonicalize.*ngoài") "lỗi containment phải nêu root bị thoát"

    # Break được bảo vệ: map trỏ đúng tên nhưng file runtime direct/slice thật đã mất.
    $missingPhysicalFile = Join-Path $tempRoot "Assets/Resources/res/x4/mainimage/btn0right.png"
    Remove-Item -LiteralPath $missingPhysicalFile
    $missingPhysicalResult = Invoke-Discovery -Root $tempRoot -Map $validMap
    Assert-True ($missingPhysicalResult.ExitCode -ne 0) "file physical direct/slice bị thiếu phải bị chặn"
    Assert-True ($missingPhysicalResult.Output -match "physical file") "lỗi thiếu file phải nêu physical file"
    [System.IO.File]::WriteAllBytes($missingPhysicalFile, [byte[]](1, 2, 3))

    # Break được bảo vệ: preflight staging bỏ lọt một future output chưa được assetpack sinh.
    $missingFutureOutput = Join-Path $tempRoot "Assets/Resources/res/x2/naruto/ui/dialog-frame.png"
    Remove-Item -LiteralPath $missingFutureOutput
    $missingFutureResult = Invoke-Discovery -Root $tempRoot -Map $validMap -RequireFutureOutputs
    Assert-True ($missingFutureResult.ExitCode -ne 0) "staging thiếu future overlay output phải bị chặn"
    Assert-True ($missingFutureResult.Output -match "future.*physical file") "lỗi staging phải nêu future physical file"
    [System.IO.File]::WriteAllBytes($missingFutureOutput, [byte[]](1, 2, 3))

    if ([string]::IsNullOrWhiteSpace($MappingPath)) {
        $MappingPath = Join-Path $repoRoot "docs\my-thuat\unity-runtime-key-map.json"
    }
    if ((Test-Path -LiteralPath $UnityRoot) -and (Test-Path -LiteralPath $MappingPath)) {
        $actualResult = Invoke-Discovery -Root $UnityRoot -Map $MappingPath
        Assert-True ($actualResult.ExitCode -eq 0) "mapping Game1 thật phải qua. Output: $($actualResult.Output)"
    }

    Write-Host "PASS: 8 binding runtime key, evidence SHA/call-site và x1..x4 đều được bảo vệ."
}
finally {
    if (Test-Path -LiteralPath $tempRoot) {
        Remove-Item -LiteralPath $tempRoot -Recurse -Force
    }
}
