# ===========================================================================
#  SUA 'localhost' -> '127.0.0.1' CHO KET NOI CSDL  (tang toc web tren VPS)
#
#  Vi sao: tren Windows Server, PHP noi CSDL qua chu "localhost" bi cham vi
#  Windows thu dia chi IPv6 truoc, cho het gio roi moi quay lai IPv4 - mat
#  khoang 1 giay MOI LAN mo trang. Doi thang sang "127.0.0.1" thi bo qua buoc
#  do, nhanh ngay.
#
#  Chi doi CHU HOST. KHONG dung toi mat khau. Co sao luu truoc va do toc do
#  truoc/sau.
#
#  Chay TREN VPS:
#     powershell -ExecutionPolicy Bypass -File "\\tsclient\C\Users\Admin\Music\bo-cai-vps\SUA-LOCALHOST-VPS.ps1"
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

$URL_DO = 'http://localhost/sv1/trangchu'
$files = @(
    'C:\xampp\htdocs\sv1\config.php',
    'C:\xampp\htdocs\sv1\boss_status.php'
)

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  SUA localhost -> 127.0.0.1 (tang toc noi CSDL)' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan

function DoTocDo($lan = 4) {
    $ms = @()
    for ($i = 0; $i -lt $lan; $i++) {
        $sw = [Diagnostics.Stopwatch]::StartNew()
        try { Invoke-WebRequest -Uri $URL_DO -UseBasicParsing -TimeoutSec 30 | Out-Null } catch {}
        $sw.Stop(); $ms += $sw.ElapsedMilliseconds
        Start-Sleep -Milliseconds 500
    }
    if ($ms.Count -gt 1) { $ms = $ms[1..($ms.Count-1)] }
    return [math]::Round(($ms | Measure-Object -Average).Average)
}

Write-Host ''
Ghi 'Do toc do TRUOC...' 'Cyan'
$truoc = DoTocDo
Ghi ("Trung binh truoc: {0} ms" -f $truoc) 'Yellow'

Write-Host ''
Ghi '--- Sua cac file ---' 'Cyan'
$daSua = $false
foreach ($f in $files) {
    if (-not (Test-Path $f)) { Ghi ("[bo qua] khong thay " + $f) 'Yellow'; continue }

    $raw = Get-Content $f -Raw
    $goc = $raw

    # Chi doi chu host trong ngu canh ket noi CSDL, KHONG dong den mat khau:
    #   $DB_HOST = "localhost";      -> "127.0.0.1"
    #   mysql:host=localhost;...     -> host=127.0.0.1
    #   host=localhost / host="localhost" / host='localhost'
    $raw = $raw -replace '(\$DB_HOST\s*=\s*["''])localhost(["''])', '${1}127.0.0.1${2}'
    $raw = $raw -replace '(mysql:host=)localhost', '${1}127.0.0.1'
    $raw = $raw -replace '(host=["''])localhost(["''])', '${1}127.0.0.1${2}'

    if ($raw -ne $goc) {
        $sao = "$f.truoc-127"
        if (-not (Test-Path $sao)) { Copy-Item $f $sao }
        Set-Content $f $raw -Encoding UTF8 -NoNewline
        Ghi ("[OK] Da sua: " + (Split-Path $f -Leaf) + "  (ban luu: " + (Split-Path $sao -Leaf) + ")") 'Green'
        $daSua = $true
    } else {
        Ghi ("[--] Khong thay 'localhost' de doi trong " + (Split-Path $f -Leaf)) 'Gray'
    }
}

if (-not $daSua) {
    Ghi 'Khong file nao thay doi. Co the da sua tu truoc.' 'Yellow'
}

# Kiem lai noi dung (chi in dong host, KHONG in mat khau)
Write-Host ''
Ghi '--- Kiem lai (chi hien dong host) ---' 'Cyan'
foreach ($f in $files) {
    if (Test-Path $f) {
        Select-String -Path $f -Pattern '127\.0\.0\.1|localhost' |
            Where-Object { $_.Line -notmatch 'pass|PASS' } |
            ForEach-Object { Ghi ("  " + (Split-Path $f -Leaf) + ":" + $_.LineNumber + "  " + ($_.Line.Trim() -replace "('root',\s*')[^']*'", "`$1***'")) 'Gray' }
    }
}

Write-Host ''
Ghi 'Do toc do SAU...' 'Cyan'
Start-Sleep -Seconds 1
$sau = DoTocDo

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  KET QUA' -ForegroundColor Green
Write-Host '===========================================================' -ForegroundColor Cyan
Ghi ("Truoc : {0} ms" -f $truoc) 'Yellow'
Ghi ("Sau   : {0} ms" -f $sau) 'Green'
if ($sau -gt 0) { Ghi ("Nhanh len khoang {0} lan" -f [math]::Round($truoc/[math]::Max($sau,1),1)) 'Green' }
Write-Host ''
Ghi 'Muon quay lai: chep de config.php.truoc-127 thanh config.php' 'Gray'
Write-Host ''
Read-Host '  Bam Enter de dong'
