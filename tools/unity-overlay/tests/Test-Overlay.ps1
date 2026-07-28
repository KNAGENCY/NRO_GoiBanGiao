param(
    [switch]$KeepFixture
)

$ErrorActionPreference = "Stop"
$script:Passed = 0

function Assert-True {
    param(
        [bool]$Condition,
        [string]$Message
    )
    if (-not $Condition) {
        throw "ASSERT: $Message"
    }
    $script:Passed++
}

function Write-Utf8NoBom {
    param(
        [string]$Path,
        [string]$Content
    )
    $parent = Split-Path -Parent $Path
    if ($parent) {
        [System.IO.Directory]::CreateDirectory($parent) | Out-Null
    }
    [System.IO.File]::WriteAllText(
        $Path,
        $Content,
        [System.Text.UTF8Encoding]::new($false)
    )
}

function Get-Sha256 {
    param([string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Invoke-Overlay {
    param(
        [string]$Script,
        [string[]]$Arguments
    )
    $stdout = Join-Path $fixtureRoot ("stdout-" + [guid]::NewGuid().ToString("N") + ".log")
    $stderr = Join-Path $fixtureRoot ("stderr-" + [guid]::NewGuid().ToString("N") + ".log")
    $process = Start-Process -FilePath "powershell.exe" `
        -ArgumentList (@("-NoProfile", "-ExecutionPolicy", "Bypass", "-File", $Script) + $Arguments) `
        -Wait -PassThru -WindowStyle Hidden -RedirectStandardOutput $stdout `
        -RedirectStandardError $stderr
    $text = ""
    if (Test-Path -LiteralPath $stdout) {
        $text += [System.IO.File]::ReadAllText($stdout)
    }
    if (Test-Path -LiteralPath $stderr) {
        $text += [System.IO.File]::ReadAllText($stderr)
    }
    return [pscustomobject]@{
        ExitCode = $process.ExitCode
        Output = $text
    }
}

function Assert-FailureCode {
    param(
        [pscustomobject]$Result,
        [string]$Code
    )
    Assert-True ($Result.ExitCode -ne 0) "Command must fail with $Code."
    Assert-True ($Result.Output -match [regex]::Escape($Code)) `
        "Missing error code $Code. Output: $($Result.Output)"
}

$overlayRoot = Split-Path -Parent $PSScriptRoot
$applyScript = Join-Path $overlayRoot "apply-overlay.ps1"
$verifyScript = Join-Path $overlayRoot "verify-overlay.ps1"
$configSource = Join-Path $overlayRoot "Assets\Scripts\NroE2E\NroE2EConfig.cs"

$fixtureRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("nro-overlay-test-" + [guid]::NewGuid().ToString("N"))
$sourceRoot = Join-Path $fixtureRoot "unity-source"
$snapshotRoot = Join-Path $fixtureRoot "unity-snapshot"
$badSnapshotRoot = Join-Path $fixtureRoot "unity-bad-sha"
$badPatchSnapshotRoot = Join-Path $fixtureRoot "unity-bad-patch"
$copyFailureSnapshotRoot = Join-Path $fixtureRoot "unity-copy-failure"
$defaultSnapshotRoot = Join-Path $fixtureRoot "unity-default-manifest"
$patchRoot = Join-Path $fixtureRoot "patches"
$fixtureManifest = Join-Path $fixtureRoot "overlay-manifest.json"

try {
    foreach ($directory in @(
        "Assets\Scripts\Game1",
        "Assets\Scripts\Game2",
        "Packages",
        "ProjectSettings"
    )) {
        [System.IO.Directory]::CreateDirectory((Join-Path $sourceRoot $directory)) | Out-Null
    }

    $rmsGame1 = @'
using System;
using UnityEngine;
namespace Game1 {
public class Rms {
    public static string GetiPhoneDocumentsPath()
    {
        return Application.persistentDataPath;
    }
}}
'@
    $rmsGame2 = $rmsGame1.Replace("namespace Game1", "namespace Game2")
    $gameScrGame1 = @'
namespace Game1 {
public class GameScr {
    public void Read(short levelId) {
        if (sks[levelId].effectHappenOnMob <= 0)
        {
            sks[levelId].effectHappenOnMob = 80;
        }
    }
}}
'@
    $gameScrGame2 = $gameScrGame1.Replace("namespace Game1", "namespace Game2")
    $charGame1 = @'
namespace Game1 {
public class Char {
    public void Attack() {
        effPaints[i] = new EffectPaint();
        effPaints[i].effCharPaint = GameScr.efs[skillPaint.effectHappenOnMob - 1];
        effPaints[i].eMob = mob;
    }
}}
'@
    $charGame2 = $charGame1.Replace("namespace Game1", "namespace Game2")
    $chatPopupGame1 = @'
namespace Game1 {
public class ChatPopup {
    public void Paint() {
        PopUp.paintPopUp(g, num, num2, num3, num4, 16777215, isButton: false);
        SmallImage.drawSmallImage(g, c.avatar, cx + 14, cy + num5, 0, StaticObj.BOTTOM_LEFT);
    }
}}
'@
    $commandGame1 = @'
namespace Game1 {
public class Command {
    public bool isPlaySoundButton = true;
    public void performAction()
	{
        GameCanvas.clearAllPointerEvent();
    }
    public void Paint() {
        if (caption != string.Empty)
        {
            PaintLegacy();
        }
    }
}}
'@
    $panelGame1 = @'
namespace Game1 {
public class Panel {
    public void Paint() {
                Item item = arrItemBox[i];
                g.setColor(11771523);
                g.fillRect(x, y, w, h, 4);
                if (item != null)
                {
                    PaintItem();
                }
                if (i == selected)
                {
                    g.setColor(0xCC6600);
                    g.drawRect(x - 2, y - 2, w + 3, h + 3);

                    g.setColor(0xFFCC66);
                    g.drawRect(x - 1, y - 1, w + 1, h + 1);
                }
    }
}}
'@

    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\Rms.cs") $rmsGame1
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game2\Rms.cs") $rmsGame2
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\GameScr.cs") $gameScrGame1
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game2\GameScr.cs") $gameScrGame2
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\Char.cs") $charGame1
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game2\Char.cs") $charGame2
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\ChatPopup.cs") $chatPopupGame1
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\Command.cs") $commandGame1
    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\Panel.cs") $panelGame1
    Write-Utf8NoBom (Join-Path $sourceRoot "Packages\manifest.json") '{"dependencies":{}}'
    Write-Utf8NoBom (Join-Path $sourceRoot "ProjectSettings\ProjectVersion.txt") "m_EditorVersion: 2022.3.62f2`n"

    [System.IO.Directory]::CreateDirectory($patchRoot) | Out-Null
    $rmsPatch = @{
        replacements = @(
            @{
                before = @'
    public static string GetiPhoneDocumentsPath()
    {
        return Application.persistentDataPath;
    }
'@
                after = @'
    public static string GetiPhoneDocumentsPath()
    {
        return NroE2E.NroE2EConfig.ResolveProfileRoot(
            Environment.GetEnvironmentVariable("NRO_E2E"),
            NroE2E.NroE2EConfig.GetArgumentValue("-profile-root"),
            Application.persistentDataPath);
    }
'@
            }
        )
    }
    $skillPatch = @{
        replacements = @(
            @{
                before = "if (sks[levelId].effectHappenOnMob <= 0)"
                after = "if (NroE2E.NroE2EConfig.ShouldUseLegacyEffect(sks[levelId].effectHappenOnMob))"
            }
        )
    }
    $charPatch = @{
        replacements = @(
            @{
                before = @'
        effPaints[i] = new EffectPaint();
        effPaints[i].effCharPaint = GameScr.efs[skillPaint.effectHappenOnMob - 1];
        effPaints[i].eMob = mob;
'@
                after = @'
        if (NroE2E.NroE2EConfig.HasSkillEffect(skillPaint.effectHappenOnMob))
        {
            effPaints[i] = new EffectPaint();
            effPaints[i].effCharPaint = GameScr.efs[skillPaint.effectHappenOnMob - 1];
            effPaints[i].eMob = mob;
        }
'@
            }
        )
    }
    Write-Utf8NoBom (Join-Path $patchRoot "rms.json") ($rmsPatch | ConvertTo-Json -Depth 8)
    Write-Utf8NoBom (Join-Path $patchRoot "skill.json") ($skillPatch | ConvertTo-Json -Depth 8)
    Write-Utf8NoBom (Join-Path $patchRoot "char.json") ($charPatch | ConvertTo-Json -Depth 8)
    Copy-Item -LiteralPath (Join-Path $overlayRoot "patches\ui-chat-popup.patch") `
        -Destination (Join-Path $patchRoot "ui-chat-popup.patch")
    Copy-Item -LiteralPath (Join-Path $overlayRoot "patches\ui-command-disabled.patch") `
        -Destination (Join-Path $patchRoot "ui-command-disabled.patch")
    Copy-Item -LiteralPath (Join-Path $overlayRoot "patches\ui-panel-slots.patch") `
        -Destination (Join-Path $patchRoot "ui-panel-slots.patch")

    $manifest = @{
        schema_version = 1
        unity_version = "2022.3.62f2"
        protected_source_roots = @($sourceRoot)
        patches = @(
            @{
                target = "Assets/Scripts/Game1/Rms.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game1\Rms.cs")
                patch = "patches/rms.json"
            },
            @{
                target = "Assets/Scripts/Game2/Rms.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game2\Rms.cs")
                patch = "patches/rms.json"
            },
            @{
                target = "Assets/Scripts/Game1/GameScr.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game1\GameScr.cs")
                patch = "patches/skill.json"
            },
            @{
                target = "Assets/Scripts/Game2/GameScr.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game2\GameScr.cs")
                patch = "patches/skill.json"
            },
            @{
                target = "Assets/Scripts/Game1/Char.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game1\Char.cs")
                patch = "patches/char.json"
            },
            @{
                target = "Assets/Scripts/Game2/Char.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game2\Char.cs")
                patch = "patches/char.json"
            },
            @{
                target = "Assets/Scripts/Game1/ChatPopup.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game1\ChatPopup.cs")
                patch = "patches/ui-chat-popup.patch"
            },
            @{
                target = "Assets/Scripts/Game1/Command.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game1\Command.cs")
                patch = "patches/ui-command-disabled.patch"
            },
            @{
                target = "Assets/Scripts/Game1/Panel.cs"
                source_sha256 = Get-Sha256 (Join-Path $sourceRoot "Assets\Scripts\Game1\Panel.cs")
                patch = "patches/ui-panel-slots.patch"
            }
        )
        overlay_files = @(
            "Assets/Scripts/NroE2E/NroE2EConfig.cs",
            "Assets/Scripts/NroE2E/NroUiOverlay.cs"
        )
    }
    Write-Utf8NoBom $fixtureManifest ($manifest | ConvertTo-Json -Depth 8)

    # A replacement count is a precondition. It must be rejected before any
    # snapshot copy, otherwise a failed apply strands an unmarked partial tree.
    $invalidRmsPatch = $rmsPatch | ConvertTo-Json -Depth 8 | ConvertFrom-Json
    $invalidRmsPatch.replacements[0] | Add-Member -NotePropertyName expected_count -NotePropertyValue 2
    Write-Utf8NoBom (Join-Path $patchRoot "rms-invalid.json") `
        ($invalidRmsPatch | ConvertTo-Json -Depth 8)
    $invalidPatchManifest = $manifest | ConvertTo-Json -Depth 8 | ConvertFrom-Json
    $invalidPatchManifest.patches[0].patch = "patches/rms-invalid.json"
    $invalidPatchManifestPath = Join-Path $fixtureRoot "overlay-manifest-invalid-patch.json"
    Write-Utf8NoBom $invalidPatchManifestPath ($invalidPatchManifest | ConvertTo-Json -Depth 8)
    $badPatch = Invoke-Overlay $applyScript @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", $badPatchSnapshotRoot,
        "-ManifestPath", $invalidPatchManifestPath
    )
    Assert-FailureCode $badPatch "OVERLAY-PATCH-001"
    Assert-True (-not (Test-Path -LiteralPath $badPatchSnapshotRoot)) `
        "A rejected patch precondition must not create an unmarked snapshot."

    # Hold a real source file exclusively so the child process fails during its
    # recursive Copy-Item. A failed copy must not strand a partial final snapshot.
    $copyFailureSource = Join-Path $sourceRoot "Assets\copy-blocked.txt"
    Write-Utf8NoBom $copyFailureSource "copy must fail while this handle is held"
    $copyLock = [System.IO.File]::Open(
        $copyFailureSource,
        [System.IO.FileMode]::Open,
        [System.IO.FileAccess]::Read,
        [System.IO.FileShare]::None
    )
    try {
        $copyFailure = Invoke-Overlay $applyScript @(
            "-UnityProject", $sourceRoot,
            "-SnapshotRoot", $copyFailureSnapshotRoot,
            "-ManifestPath", $fixtureManifest
        )
    }
    finally {
        $copyLock.Dispose()
    }
    Assert-True ($copyFailure.ExitCode -ne 0) "An exclusively locked source file must make copy fail."
    Assert-True (-not (Test-Path -LiteralPath $copyFailureSnapshotRoot)) `
        "A failed copy must leave no partial final snapshot."
    Assert-True (-not (Test-Path -LiteralPath (Join-Path $copyFailureSnapshotRoot ".nro-e2e-overlay.json"))) `
        "A failed copy must not leave a false overlay marker."
    $stagingLeftovers = @(Get-ChildItem -LiteralPath $fixtureRoot -Directory -Force -Filter `
        ((Split-Path -Leaf $copyFailureSnapshotRoot) + ".staging-*"))
    Assert-True ($stagingLeftovers.Count -eq 0) `
        "A failed copy must roll back its staging directory."

    $defaultToolRoot = Join-Path $fixtureRoot "default-tool"
    [System.IO.Directory]::CreateDirectory((Join-Path $defaultToolRoot "patches")) | Out-Null
    [System.IO.Directory]::CreateDirectory(
        (Join-Path $defaultToolRoot "Assets\Scripts\NroE2E")
    ) | Out-Null
    Copy-Item -LiteralPath $applyScript -Destination (Join-Path $defaultToolRoot "apply-overlay.ps1")
    Copy-Item -LiteralPath $verifyScript -Destination (Join-Path $defaultToolRoot "verify-overlay.ps1")
    Copy-Item -LiteralPath $fixtureManifest -Destination (Join-Path $defaultToolRoot "overlay-manifest.json")
    Copy-Item -LiteralPath (Join-Path $patchRoot "rms.json") -Destination (Join-Path $defaultToolRoot "patches\rms.json")
    Copy-Item -LiteralPath (Join-Path $patchRoot "skill.json") -Destination (Join-Path $defaultToolRoot "patches\skill.json")
    Copy-Item -LiteralPath (Join-Path $patchRoot "char.json") -Destination (Join-Path $defaultToolRoot "patches\char.json")
    Copy-Item -LiteralPath (Join-Path $patchRoot "ui-chat-popup.patch") -Destination (Join-Path $defaultToolRoot "patches\ui-chat-popup.patch")
    Copy-Item -LiteralPath (Join-Path $patchRoot "ui-command-disabled.patch") -Destination (Join-Path $defaultToolRoot "patches\ui-command-disabled.patch")
    Copy-Item -LiteralPath (Join-Path $patchRoot "ui-panel-slots.patch") -Destination (Join-Path $defaultToolRoot "patches\ui-panel-slots.patch")
    Copy-Item -LiteralPath $configSource `
        -Destination (Join-Path $defaultToolRoot "Assets\Scripts\NroE2E\NroE2EConfig.cs")
    Copy-Item -LiteralPath (Join-Path $overlayRoot "Assets\Scripts\NroE2E\NroUiOverlay.cs") `
        -Destination (Join-Path $defaultToolRoot "Assets\Scripts\NroE2E\NroUiOverlay.cs")
    $defaultManifest = Invoke-Overlay (Join-Path $defaultToolRoot "apply-overlay.ps1") @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", $defaultSnapshotRoot
    )
    Assert-True ($defaultManifest.ExitCode -eq 0) `
        "Default manifest beside script must be used: $($defaultManifest.Output)"

    $samePath = Invoke-Overlay $applyScript @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", $sourceRoot,
        "-ManifestPath", $fixtureManifest
    )
    Assert-FailureCode $samePath "OVERLAY-SOURCE-PROTECTED"

    $childOfSource = Invoke-Overlay $applyScript @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", (Join-Path $sourceRoot "staging-must-not-be-written"),
        "-ManifestPath", $fixtureManifest
    )
    Assert-FailureCode $childOfSource "OVERLAY-SOURCE-PROTECTED"

    $apply = Invoke-Overlay $applyScript @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", $snapshotRoot,
        "-ManifestPath", $fixtureManifest
    )
    Assert-True ($apply.ExitCode -eq 0) "Fixture overlay failed: $($apply.Output)"
    Assert-True (Test-Path -LiteralPath (Join-Path $snapshotRoot ".nro-e2e-overlay.json")) `
        "Snapshot is missing the overlay marker."
    Assert-True (Test-Path -LiteralPath (Join-Path $snapshotRoot "Assets\Scripts\NroE2E\NroE2EConfig.cs")) `
        "Snapshot is missing NroE2EConfig."
    $patchedCommand = [System.IO.File]::ReadAllText(
        (Join-Path $snapshotRoot "Assets\Scripts\Game1\Command.cs")
    )
    Assert-True ($patchedCommand -match "public bool isDisabled;") `
        "Command overlay must expose an explicit disabled state."
    Assert-True ($patchedCommand -match "if \(isDisabled\)\s*\{\s*return;") `
        "Disabled commands must not perform their action."
    Assert-True (
        $patchedCommand -match "if \(isDisabled && NroE2E\.NroUiOverlay\.PaintButtonDisabled"
    ) "Disabled skin must be activated by isDisabled, not a guessed caption rule."

    $verify = Invoke-Overlay $verifyScript @(
        "-SnapshotRoot", $snapshotRoot,
        "-ManifestPath", $fixtureManifest
    )
    Assert-True ($verify.ExitCode -eq 0) "Overlay verification failed: $($verify.Output)"

    $doubleApply = Invoke-Overlay $applyScript @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", $snapshotRoot,
        "-ManifestPath", $fixtureManifest
    )
    Assert-FailureCode $doubleApply "OVERLAY-ALREADY-APPLIED"

    Write-Utf8NoBom (Join-Path $sourceRoot "Assets\Scripts\Game1\Rms.cs") ($rmsGame1 + "// mutation")
    $badSha = Invoke-Overlay $applyScript @(
        "-UnityProject", $sourceRoot,
        "-SnapshotRoot", $badSnapshotRoot,
        "-ManifestPath", $fixtureManifest
    )
    Assert-FailureCode $badSha "OVERLAY-SHA-001"

    Add-Type -Path $configSource
    $productionRoot = Join-Path $fixtureRoot "production-profile"
    $profileA = Join-Path $fixtureRoot "profile-a"
    $profileB = Join-Path $fixtureRoot "profile-b"
    Assert-True (
        [NroE2E.NroE2EConfig]::ResolveProfileRoot($null, $profileA, $productionRoot) -eq $productionRoot
    ) "Without NRO_E2E, production persistentDataPath must be preserved."
    Assert-True (
        [NroE2E.NroE2EConfig]::ResolveProfileRoot("0", $profileA, $productionRoot) -eq $productionRoot
    ) "NRO_E2E=0 must not enable isolated profiles."
    $resolvedA = [NroE2E.NroE2EConfig]::ResolveProfileRoot("1", $profileA, $productionRoot)
    $resolvedB = [NroE2E.NroE2EConfig]::ResolveProfileRoot("1", $profileB, $productionRoot)
    Write-Utf8NoBom (Join-Path $resolvedA "A.sentinel") "A"
    Write-Utf8NoBom (Join-Path $resolvedB "B.sentinel") "B"
    Assert-True (-not (Test-Path -LiteralPath (Join-Path $resolvedA "B.sentinel"))) `
        "Profile A can see B's sentinel."
    Assert-True (-not (Test-Path -LiteralPath (Join-Path $resolvedB "A.sentinel"))) `
        "Profile B can see A's sentinel."

    $relativeRejected = $false
    try {
        [NroE2E.NroE2EConfig]::ResolveProfileRoot("1", "relative-profile", $productionRoot) | Out-Null
    }
    catch {
        $relativeRejected = $true
    }
    Assert-True $relativeRejected "Relative E2E profile must be rejected."
    Assert-True (-not [NroE2E.NroE2EConfig]::ShouldUseLegacyEffect(-1)) `
        "-1 must remain no-effect and must not become effect 80."
    Assert-True ([NroE2E.NroE2EConfig]::ShouldUseLegacyEffect(0)) `
        "Legacy zero must retain the effect-80 fallback."
    Assert-True (-not [NroE2E.NroE2EConfig]::HasSkillEffect(-1)) `
        "Skill -1 must not create EffectPaint or index the effect array."

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
