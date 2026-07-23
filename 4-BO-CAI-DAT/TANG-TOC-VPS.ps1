# ===========================================================================
#  TANG TOC WEB TREN VPS  (mien phi, khong nang cap phan cung)
#
#  Lam 2 viec da chan doan la nguyen nhan chinh:
#    1. Bat OPcache cho PHP  -> PHP nho ma da dich, khong dich lai moi lan
#    2. Loai thu muc web khoi vung Windows Defender quet
#
#  Chay TREN VPS. Co SAO LUU php.ini truoc khi sua, va tu DO LAI toc do
#  truoc/sau de biet co nhanh len that khong.
#
#  Cach chay:
#     powershell -ExecutionPolicy Bypass -File "\\tsclient\C\Users\Admin\Music\bo-cai-vps\TANG-TOC-VPS.ps1"
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

$PHPINI = 'C:\xampp\php\php.ini'
$WEBDIR = 'C:\xampp\htdocs'
$GAMEDIR = 'C:\zz1'
$URL_DO = 'http://localhost/sv1/trangchu'

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  TANG TOC WEB TREN VPS' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan

if (-not (Test-Path $PHPINI)) { Ghi "DUNG LAI: khong thay $PHPINI" 'Red'; Read-Host "`n  Enter de thoat"; exit 1 }

# --- DO TOC DO TRUOC ------------------------------------------------------
function DoTocDo($lan = 4) {
    $ms = @()
    for ($i = 0; $i -lt $lan; $i++) {
        $sw = [Diagnostics.Stopwatch]::StartNew()
        try { Invoke-WebRequest -Uri $URL_DO -UseBasicParsing -TimeoutSec 30 | Out-Null } catch {}
        $sw.Stop(); $ms += $sw.ElapsedMilliseconds
        Start-Sleep -Milliseconds 500
    }
    # bo lan dau (con lanh), lay trung binh phan con lai
    if ($ms.Count -gt 1) { $ms = $ms[1..($ms.Count-1)] }
    return [math]::Round(($ms | Measure-Object -Average).Average)
}

Write-Host ''
Ghi 'Do toc do TRUOC khi sua...' 'Cyan'
$truoc = DoTocDo
Ghi ("Trung binh truoc: {0} ms" -f $truoc) 'Yellow'

# --- 1. BAT OPCACHE -------------------------------------------------------
Write-Host ''
Ghi '--- Bat OPcache ---' 'Cyan'

$sao = "$PHPINI.truoc-tangtoc"
if (-not (Test-Path $sao)) { Copy-Item $PHPINI $sao; Ghi "[OK] Da sao luu php.ini -> php.ini.truoc-tangtoc" 'Green' }
else { Ghi "[OK] Da co ban sao luu tu truoc, giu nguyen" 'Gray' }

$noidung = Get-Content $PHPINI

# Bo dau ; o dong zend_extension=opcache (bat buoc de OPcache nap)
$noidung = $noidung -replace '^\s*;\s*zend_extension\s*=\s*opcache\s*$', 'zend_extension=opcache'

# Cac tham so OPcache: bo comment + dat gia tri tot
$dat = @{
    'opcache.enable'                 = '1'
    'opcache.enable_cli'             = '0'
    'opcache.memory_consumption'     = '256'
    'opcache.interned_strings_buffer'= '16'
    'opcache.max_accelerated_files'  = '20000'
    'opcache.validate_timestamps'    = '1'
    'opcache.revalidate_freq'        = '60'
}
foreach ($k in $dat.Keys) {
    $v = $dat[$k]
    $re = '^\s*;?\s*' + [regex]::Escape($k) + '\s*=.*$'
    if ($noidung -match $re) {
        $noidung = $noidung -replace $re, "$k=$v"
    } else {
        # chua co dong nao -> them vao cuoi khoi [opcache]
        $idx = ($noidung | Select-String -Pattern '^\[opcache\]' | Select-Object -First 1).LineNumber
        if ($idx) { $noidung = $noidung[0..($idx-1)] + "$k=$v" + $noidung[$idx..($noidung.Count-1)] }
        else { $noidung += "$k=$v" }
    }
}
Set-Content $PHPINI $noidung -Encoding UTF8
Ghi '[OK] Da ghi cau hinh OPcache vao php.ini' 'Green'

# --- 2. LOAI TRU DEFENDER -------------------------------------------------
Write-Host ''
Ghi '--- Loai thu muc web + game khoi Windows Defender ---' 'Cyan'
foreach ($d in @($WEBDIR, $GAMEDIR, 'C:\xampp\php', 'C:\xampp\apache')) {
    if (Test-Path $d) {
        try { Add-MpPreference -ExclusionPath $d -ErrorAction Stop; Ghi ("[OK] Da loai tru: " + $d) 'Green' }
        catch { Ghi ("[bo qua] khong loai tru duoc " + $d + " (" + $_.Exception.Message + ")") 'Yellow' }
    }
}
# Loai tru tien trinh PHP/Apache cho chac
foreach ($p in @('C:\xampp\php\php.exe', 'C:\xampp\apache\bin\httpd.exe')) {
    if (Test-Path $p) { try { Add-MpPreference -ExclusionProcess $p -ErrorAction Stop } catch {} }
}

# --- KHOI DONG LAI APACHE -------------------------------------------------
Write-Host ''
Ghi '--- Khoi dong lai Apache de nap cau hinh moi ---' 'Cyan'
$httpd = 'C:\xampp\apache\bin\httpd.exe'
# Kiem cu phap truoc
$cu = $ErrorActionPreference; $ErrorActionPreference = 'Continue'
$kt = & $httpd -t 2>&1 | Out-String
$ErrorActionPreference = $cu
if ($kt -notmatch 'Syntax OK') { Ghi 'CANH BAO: httpd -t khong bao Syntax OK:' 'Red'; Ghi $kt.Trim() 'Red' }

Get-Process httpd -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 2
# XAMPP thuong chay Apache nhu dich vu 'Apache2.4'; thu ca hai cach
$svc = Get-Service -Name 'Apache2.4' -ErrorAction SilentlyContinue
if ($svc) { Start-Service 'Apache2.4'; Ghi '[OK] Da khoi dong lai dich vu Apache2.4' 'Green' }
else {
    Start-Process -FilePath $httpd -WindowStyle Hidden
    Ghi '[OK] Da chay lai httpd.exe' 'Green'
}
Start-Sleep -Seconds 4

# --- KIEM CHUNG OPCACHE DA BAT --------------------------------------------
Write-Host ''
Ghi '--- Kiem chung ---' 'Cyan'
$oc = & C:\xampp\php\php.exe -r "echo ini_get('opcache.enable') ? 'DA BAT' : 'CHUA BAT';" 2>$null
Ghi ("OPcache (dong lenh): " + $oc) 'Gray'

# --- DO TOC DO SAU --------------------------------------------------------
Write-Host ''
Ghi 'Do toc do SAU khi sua...' 'Cyan'
Start-Sleep -Seconds 2
$sau = DoTocDo

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  KET QUA' -ForegroundColor Green
Write-Host '===========================================================' -ForegroundColor Cyan
Ghi ("Truoc : {0} ms" -f $truoc) 'Yellow'
Ghi ("Sau   : {0} ms" -f $sau) 'Green'
if ($truoc -gt 0 -and $sau -gt 0) {
    $lan = [math]::Round($truoc / [math]::Max($sau,1), 1)
    Ghi ("Nhanh len khoang {0} lan" -f $lan) 'Green'
}
Write-Host ''
Ghi 'Neu van cham (tren 1 giay): bao tro ly de xem tiep. Neu nhanh han' 'Gray'
Ghi 'han: xong, khong can nang cap gi.' 'Gray'
Write-Host ''
Ghi 'Muon quay lai nhu cu: chep de php.ini.truoc-tangtoc thanh php.ini' 'Gray'
Write-Host ''
Read-Host '  Bam Enter de dong'
