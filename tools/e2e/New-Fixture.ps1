[CmdletBinding()]
param(
    [string]$Dsn,
    [string]$RunId,
    [ValidateSet("Create", "Validate", "Cleanup")]
    [string]$Action = "Create",
    [string]$ManifestPath,
    [switch]$SelfTest
)

$ErrorActionPreference = "Stop"

function Assert-E2ETestDsn {
    param([string]$Value)
    try {
        $uri = [Uri]$Value
    }
    catch {
        throw "E2E-FIXTURE-001: DSN khong hop le"
    }
    $database = $uri.AbsolutePath.Trim('/')
    if (($uri.Scheme -ne "postgres" -and $uri.Scheme -ne "postgresql") -or
        $database -cne "nro_test") {
        throw "E2E-FIXTURE-001: fixture chi duoc dung database nro_test"
    }
    $hostName = $uri.Host
    $loopback = $false
    if ($hostName -eq "localhost") {
        $loopback = $true
    }
    else {
        $ip = $null
        if ([System.Net.IPAddress]::TryParse($hostName, [ref]$ip)) {
            $loopback = [System.Net.IPAddress]::IsLoopback($ip)
        }
    }
    if (-not $loopback) {
        throw "E2E-FIXTURE-001: fixture chi duoc noi PostgreSQL loopback"
    }
    return $true
}

function ConvertTo-E2EPrefix {
    param([string]$Value)
    if ([string]::IsNullOrWhiteSpace($Value)) {
        throw "E2E-FIXTURE-002: run ID rong"
    }
    $normalized = $Value.Normalize([Text.NormalizationForm]::FormD)
    $ascii = New-Object Text.StringBuilder
    foreach ($character in $normalized.ToCharArray()) {
        if ([Globalization.CharUnicodeInfo]::GetUnicodeCategory($character) -eq
            [Globalization.UnicodeCategory]::NonSpacingMark) {
            continue
        }
        $code = [int][char]$character
        if ($code -lt 128) {
            [void]$ascii.Append([char]$code)
        }
    }
    $slug = $ascii.ToString().ToLowerInvariant()
    $slug = [regex]::Replace($slug, '[^a-z0-9]+', '_').Trim('_')
    if ([string]::IsNullOrWhiteSpace($slug)) {
        $slug = "run"
    }
    if ($slug.Length -gt 8) {
        $slug = $slug.Substring(0, 8)
    }
    $bytes = [Text.Encoding]::UTF8.GetBytes($Value)
    $hasher = [Security.Cryptography.SHA256]::Create()
    try {
        $hash = ([BitConverter]::ToString($hasher.ComputeHash($bytes))).Replace("-", "").ToLowerInvariant()
    }
    finally {
        $hasher.Dispose()
    }
    $prefix = "e2e_${slug}_$($hash.Substring(0, 6))"
    if ($prefix -notmatch '^e2e_[a-z0-9_]{4,15}$') {
        throw "E2E-FIXTURE-002: khong tao duoc prefix ASCII hop le"
    }
    return $prefix
}

function New-E2EPasswordHash {
    param([Parameter(Mandatory = $true)][string]$Password)
    if ([string]::IsNullOrWhiteSpace($Password)) {
        throw "E2E-FIXTURE-009: mat khau fixture rong"
    }
    $salt = New-Object byte[] 16
    $random = [Security.Cryptography.RandomNumberGenerator]::Create()
    try {
        $random.GetBytes($salt)
    }
    finally {
        $random.Dispose()
    }
    $derive = [Security.Cryptography.Rfc2898DeriveBytes]::new(
        $Password,
        $salt,
        100000,
        [Security.Cryptography.HashAlgorithmName]::SHA256)
    try {
        $digest = $derive.GetBytes(32)
    }
    finally {
        $derive.Dispose()
    }
    return (
        "pbkdf2" + '$' + "100000" + '$' +
        [Convert]::ToBase64String($salt) + '$' +
        [Convert]::ToBase64String($digest))
}

function New-E2ECleanupSql {
    param(
        [string]$Prefix,
        [string[]]$Accounts
    )
    if ($Prefix -notmatch '^e2e_[a-z0-9_]{4,15}$' -or
        $Accounts.Count -ne 2 -or
        $Accounts[0] -cne ($Prefix + "a") -or
        $Accounts[1] -cne ($Prefix + "b")) {
        throw "E2E-FIXTURE-008: cleanup target khong khop dung hai account run ID"
    }
    return @"
BEGIN;
DELETE FROM nhan_vat
 WHERE tai_khoan_id IN (
   SELECT id FROM tai_khoan
    WHERE ten_dang_nhap IN ('$($Accounts[0])', '$($Accounts[1])')
 );
DELETE FROM tai_khoan
 WHERE ten_dang_nhap IN ('$($Accounts[0])', '$($Accounts[1])');
COMMIT;
"@
}

function Assert-FixtureManifest {
    param($Manifest)
    if ($null -eq $Manifest -or $Manifest.schema_version -ne 1) {
        throw "E2E-FIXTURE-003: manifest schema sai"
    }
    $prefix = [string]$Manifest.prefix
    $accounts = @($Manifest.accounts)
    if ($prefix -notmatch '^e2e_[a-z0-9_]{4,15}$' -or
        $accounts.Count -ne 2 -or
        $accounts[0] -cne ($prefix + "a") -or
        $accounts[1] -cne ($prefix + "b")) {
        throw "E2E-FIXTURE-002: manifest prefix/account sai"
    }
    $characters = @($Manifest.characters)
    if ($characters.Count -ne 2) {
        throw "E2E-FIXTURE-003: manifest phai co dung hai character"
    }
    $expectedX = @(100, 148)
    for ($index = 0; $index -lt 2; $index++) {
        $character = $characters[$index]
        if ($character.lang -ne 0 -or $character.map_id -ne 1 -or
            $character.zone_id -ne 0 -or $character.kieu_toc -ne 30 -or
            $character.head -ne 30 -or $character.body -ne 57 -or
            $character.leg -ne 58 -or $character.x -ne $expectedX[$index] -or
            $character.y -ne 384 -or $character.hp -ne 150 -or
            $character.hp_max -ne 150 -or $character.chakra -ne 150 -or
            $character.chakra_max -ne 150) {
            throw "E2E-FIXTURE-003: character $index lech contract Lang La"
        }
        $skills = @($character.skills)
        if ($skills.Count -ne 4) {
            throw "E2E-FIXTURE-003: character $index thieu bon skill"
        }
        for ($skill = 0; $skill -lt 4; $skill++) {
            if ($skills[$skill].template_id -ne ($skill + 1) -or
                $skills[$skill].level -ne 1) {
                throw "E2E-FIXTURE-003: skill manifest sai"
            }
        }
    }
    return $true
}

function Write-Utf8Json {
    param([string]$Path, $Value)
    $parent = Split-Path -Parent $Path
    if (-not [string]::IsNullOrWhiteSpace($parent)) {
        New-Item -ItemType Directory -Path $parent -Force | Out-Null
    }
    $json = $Value | ConvertTo-Json -Depth 10
    [IO.File]::WriteAllText(
        [IO.Path]::GetFullPath($Path),
        $json,
        (New-Object Text.UTF8Encoding($false)))
}

function Invoke-E2EPsql {
    param([string]$Connection, [string]$Sql, [switch]$ReturnRows)
    $psql = if ($env:NRO_PSQL) { $env:NRO_PSQL } else { "D:\PostgreSQL\bin\psql.exe" }
    if (-not (Test-Path -LiteralPath $psql -PathType Leaf)) {
        throw "E2E-FIXTURE-004: khong tim thay psql"
    }
    $arguments = @(
        "--no-psqlrc",
        "--set=ON_ERROR_STOP=1",
        "--quiet",
        "--dbname=$Connection"
    )
    if ($ReturnRows) {
        $arguments += @("--tuples-only", "--no-align")
    }
    $arguments += @("--command=$Sql")
    $output = & $psql @arguments 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "E2E-FIXTURE-005: psql that bai"
    }
    if ($ReturnRows) {
        return @($output)
    }
}

function New-E2EManifest {
    param([string]$Id)
    $prefix = ConvertTo-E2EPrefix $Id
    $hashPart = $prefix.Substring([Math]::Max(0, $prefix.Length - 6)).ToUpperInvariant()
    $skills = @(
        @{ template_id = 1; level = 1 },
        @{ template_id = 2; level = 1 },
        @{ template_id = 3; level = 1 },
        @{ template_id = 4; level = 1 }
    )
    return @{
        schema_version = 1
        run_id = $Id
        prefix = $prefix
        accounts = @(($prefix + "a"), ($prefix + "b"))
        characters = @(
            @{
                name = ("E2A" + $hashPart + "0").Substring(0, 10)
                lang = 0; map_id = 1; zone_id = 0; kieu_toc = 30
                head = 30; body = 57; leg = 58
                x = 100; y = 384
                hp = 150; hp_max = 150; chakra = 150; chakra_max = 150
                skills = $skills
            },
            @{
                name = ("E2B" + $hashPart + "0").Substring(0, 10)
                lang = 0; map_id = 1; zone_id = 0; kieu_toc = 30
                head = 30; body = 57; leg = 58
                x = 148; y = 384
                hp = 150; hp_max = 150; chakra = 150; chakra_max = 150
                skills = $skills
            }
        )
        cleanup_completed = $false
    }
}

function Invoke-FixtureSelfTest {
    try {
        Assert-E2ETestDsn "postgres://localhost/nro"
        throw "SELFTEST-FAIL: production DB was accepted"
    }
    catch {
        if (-not $_.Exception.Message.Contains("E2E-FIXTURE-001")) {
            throw "SELFTEST-FAIL: expected E2E-FIXTURE-001, got $($_.Exception.Message)"
        }
    }

    $prefix = ConvertTo-E2EPrefix "Đợt chạy / ../ 42"
    if ($prefix -notmatch '^e2e_[a-z0-9_]{4,15}$') {
        throw "SELFTEST-FAIL: E2E-FIXTURE-002 prefix sai: $prefix"
    }
    $rawPassword = "fixture-secret-only-in-env-42"
    $passwordHash = New-E2EPasswordHash -Password $rawPassword
    if ($passwordHash -notmatch '^pbkdf2\$100000\$[A-Za-z0-9+/]+=*\$[A-Za-z0-9+/]+=*$') {
        throw "SELFTEST-FAIL: E2E-FIXTURE-009 hash khong dung format production"
    }
    if ($passwordHash -ceq $rawPassword -or $passwordHash.Contains($rawPassword)) {
        throw "SELFTEST-FAIL: E2E-FIXTURE-009 mat khau tho lot vao hash"
    }
    $accounts = @("${prefix}a", "${prefix}b")
    $sql = New-E2ECleanupSql -Prefix $prefix -Accounts $accounts
    foreach ($account in $accounts) {
        if (-not $sql.Contains($account)) {
            throw "SELFTEST-FAIL: E2E-FIXTURE-008 cleanup thieu $account"
        }
    }
    if (-not $sql.Contains("DELETE FROM nhan_vat") -or -not $sql.Contains("DELETE FROM tai_khoan")) {
        throw "SELFTEST-FAIL: E2E-FIXTURE-008 cleanup khong xoa child truoc parent"
    }
    if ($sql.IndexOf("DELETE FROM nhan_vat") -gt $sql.IndexOf("DELETE FROM tai_khoan")) {
        throw "SELFTEST-FAIL: E2E-FIXTURE-008 cleanup sai thu tu child/parent"
    }
    $manifest = New-E2EManifest "self-test-42"
    if (-not (Assert-FixtureManifest $manifest)) {
        throw "SELFTEST-FAIL: manifest dung bi tu choi"
    }
    $manifest.accounts = @($manifest.accounts[0], "production_admin")
    try {
        Assert-FixtureManifest $manifest | Out-Null
        throw "SELFTEST-FAIL: E2E-FIXTURE-008 account ngoai prefix duoc chap nhan"
    }
    catch {
        if (-not $_.Exception.Message.Contains("E2E-FIXTURE-002")) {
            throw
        }
    }
}

if ($SelfTest) {
    Invoke-FixtureSelfTest
    Write-Output "FIXTURE-SELFTEST-PASS"
    exit 0
}

Assert-E2ETestDsn $Dsn | Out-Null
if ([string]::IsNullOrWhiteSpace($ManifestPath)) {
    throw "E2E-FIXTURE-003: thieu -ManifestPath"
}
$fullManifestPath = [IO.Path]::GetFullPath($ManifestPath)

switch ($Action) {
    "Create" {
        if ([string]::IsNullOrWhiteSpace($RunId)) {
            throw "E2E-FIXTURE-002: thieu -RunId"
        }
        $password = $env:NRO_E2E_FIXTURE_PASSWORD
        if ($password -notmatch '^[A-Za-z0-9_-]{16,64}$') {
            throw "E2E-FIXTURE-006: can NRO_E2E_FIXTURE_PASSWORD tam thoi, 16-64 ky tu ASCII"
        }
        $manifest = New-E2EManifest $RunId
        Assert-FixtureManifest $manifest | Out-Null
        $accounts = @($manifest.accounts)
        $characters = @($manifest.characters)
        $passwordHashA = New-E2EPasswordHash -Password $password
        $passwordHashB = New-E2EPasswordHash -Password $password
        $sql = @"
BEGIN;
INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam)
VALUES ('$($accounts[0])', '$passwordHashA'), ('$($accounts[1])', '$passwordHashB');
INSERT INTO nhan_vat
  (tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
   hp_toi_da, hp_hien_tai, chakra_toi_da, chakra_hien_tai)
VALUES
  ((SELECT id FROM tai_khoan WHERE ten_dang_nhap = '$($accounts[0])'),
   '$($characters[0].name)', 0, 30, 1, 100, 384, 150, 150, 150, 150),
  ((SELECT id FROM tai_khoan WHERE ten_dang_nhap = '$($accounts[1])'),
   '$($characters[1].name)', 0, 30, 1, 148, 384, 150, 150, 150, 150);
COMMIT;
"@
        try {
            Invoke-E2EPsql -Connection $Dsn -Sql $sql
            Write-Utf8Json -Path $fullManifestPath -Value $manifest
        }
        catch {
            try {
                $cleanup = New-E2ECleanupSql -Prefix $manifest.prefix -Accounts $accounts
                Invoke-E2EPsql -Connection $Dsn -Sql $cleanup
            }
            catch {
                # Bao loi goc; runner van thu cleanup trong finally.
            }
            throw
        }
        Write-Output $fullManifestPath
    }
    "Validate" {
        if (-not (Test-Path -LiteralPath $fullManifestPath -PathType Leaf)) {
            throw "E2E-FIXTURE-003: manifest khong ton tai"
        }
        $manifest = Get-Content -LiteralPath $fullManifestPath -Raw -Encoding UTF8 | ConvertFrom-Json
        Assert-FixtureManifest $manifest | Out-Null
        $accounts = @($manifest.accounts)
        $rows = Invoke-E2EPsql -Connection $Dsn -ReturnRows -Sql @"
SELECT count(*)
  FROM tai_khoan tk
  JOIN nhan_vat nv ON nv.tai_khoan_id = tk.id
 WHERE tk.ten_dang_nhap IN ('$($accounts[0])', '$($accounts[1])')
   AND nv.lang = 0 AND nv.kieu_toc = 30 AND nv.map_id = 1
   AND nv.hp_toi_da = 150 AND nv.hp_hien_tai = 150
   AND nv.chakra_toi_da = 150 AND nv.chakra_hien_tai = 150;
"@
        $count = [int](($rows | Where-Object { "$_".Trim() -match '^\d+$' } | Select-Object -Last 1).Trim())
        if ($count -ne 2) {
            throw "E2E-FIXTURE-007: DB khong co dung hai fixture"
        }
        $hashRows = Invoke-E2EPsql -Connection $Dsn -ReturnRows -Sql @"
SELECT mat_khau_bam
  FROM tai_khoan
 WHERE ten_dang_nhap IN ('$($accounts[0])', '$($accounts[1])')
 ORDER BY ten_dang_nhap;
"@
        $hashes = @(
            $hashRows |
                ForEach-Object { "$_".Trim() } |
                Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
        if ($hashes.Count -ne 2) {
            throw "E2E-FIXTURE-009: DB khong co dung hai password hash"
        }
        $rawPassword = $env:NRO_E2E_FIXTURE_PASSWORD
        foreach ($hash in $hashes) {
            if ($hash -notmatch '^pbkdf2\$100000\$[A-Za-z0-9+/]+=*\$[A-Za-z0-9+/]+=*$' -or
                (-not [string]::IsNullOrEmpty($rawPassword) -and
                    $hash -ceq $rawPassword)) {
                throw "E2E-FIXTURE-009: mat_khau_bam khong dung format production"
            }
        }
        Write-Output "FIXTURE-VALID"
    }
    "Cleanup" {
        if (-not (Test-Path -LiteralPath $fullManifestPath -PathType Leaf)) {
            throw "E2E-FIXTURE-003: khong cleanup khi thieu manifest target"
        }
        $manifest = Get-Content -LiteralPath $fullManifestPath -Raw -Encoding UTF8 | ConvertFrom-Json
        Assert-FixtureManifest $manifest | Out-Null
        $cleanup = New-E2ECleanupSql -Prefix $manifest.prefix -Accounts @($manifest.accounts)
        Invoke-E2EPsql -Connection $Dsn -Sql $cleanup
        $manifest.cleanup_completed = $true
        Write-Utf8Json -Path $fullManifestPath -Value $manifest
        Write-Output "FIXTURE-CLEANED"
    }
}
