# kiem-tra-tat-ca.ps1 — CI local: chay HET kiem tra cua du an trong MOT lenh.
# Chay tat ca cac buoc (khong dung o buoc do dau tien) roi bao cao; thoat mã
# khac 0 neu co bat ky buoc nao do. Dung boi pre-push hook va chay tay truoc commit.
#
#   powershell -ExecutionPolicy Bypass -File .\kiem-tra-tat-ca.ps1

$goc = $PSScriptRoot
$loi = New-Object System.Collections.Generic.List[string]

# --- Tim go/gofmt: uu tien PATH, khong co thi D:\Go\bin (theo tai lieu du an) ---
function TimExe($ten, $duPhong) {
    $c = Get-Command $ten -ErrorAction SilentlyContinue
    if ($c) { return $c.Source }
    if (Test-Path $duPhong) { return $duPhong }
    return $null
}
$go = TimExe 'go' 'D:\Go\bin\go.exe'
$gofmt = TimExe 'gofmt' 'D:\Go\bin\gofmt.exe'
if (-not $go) { Write-Host "[DO] khong tim thay 'go' (PATH hoac D:\Go\bin)" -ForegroundColor Red; exit 1 }

# --- Tu do cong PostgreSQL de dat NRO_TEST_DSN (test can DB se chay; neu khong SKIP) ---
$cong = Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue |
    Where-Object { $_.LocalPort -in 5432, 5433 } |
    Select-Object -ExpandProperty LocalPort -Unique | Select-Object -First 1
$khongCoDB = $false
if ($cong) {
    $env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:$cong/nro_test"
    Write-Host "PostgreSQL o cong $cong -> NRO_TEST_DSN da dat (test can DB SE chay)" -ForegroundColor DarkGray
} else {
    $khongCoDB = $true
    Write-Host "PostgreSQL khong chay -> test can DB tu SKIP" -ForegroundColor Yellow
}

function Buoc($ten, [scriptblock]$khoi) {
    Write-Host "`n=== $ten ===" -ForegroundColor Cyan
    & $khoi
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[DO] $ten (exit $LASTEXITCODE)" -ForegroundColor Red
        $loi.Add($ten)
    } else {
        Write-Host "[OK] $ten" -ForegroundColor Green
    }
}

Push-Location "$goc\server"
try {
    Buoc "go build" { & $go build ./... }
    Buoc "go vet"   { & $go vet ./... }

    # gofmt: sau khi chuan hoa LF (.gitattributes `eol=lf` + blob deu LF), gofmt phai
    # SACH tren moi working tree dung. Khong con danh sach nhieu CRLF hardcode nua —
    # gofmt bao file nao la loi that: format sai, HOAC working tree con CRLF cu chua
    # renormalize (sua: `git checkout -- <file>` hoac `git add --renormalize .`).
    Write-Host "`n=== gofmt ===" -ForegroundColor Cyan
    if ($gofmt) {
        $xau = @(& $gofmt -l internal cmd 2>&1 | Where-Object { $_ })
        if ($xau.Count -gt 0) {
            Write-Host "[DO] gofmt: file sai dinh dang (hoac con CRLF cu chua chuan hoa):" -ForegroundColor Red
            $xau | ForEach-Object { Write-Host "     $_" -ForegroundColor Red }
            Write-Host "     Neu chi la CRLF cu: git checkout -- <file>  (hoac: git add --renormalize .)" -ForegroundColor Yellow
            $loi.Add("gofmt")
        } else {
            Write-Host "[OK] gofmt" -ForegroundColor Green
        }
    } else {
        Write-Host "[BO QUA] khong tim thay gofmt" -ForegroundColor Yellow
    }

    Buoc "go test" { & $go test ./... }
} finally { Pop-Location }

# --- assetpack (Python) ---
if (Test-Path "$goc\assetpack") {
    $py = TimExe 'python' 'python'
    if ($py) {
        # Chay TU GOC repo (khong phai trong assetpack/) de "import assetpack" thay goi.
        # KHONG con --ignore nao: test_task6b_canonical_preflight.py tung phai bo vi
        # SHARED_ROOT = ROOT.parents[1] gia dinh repo nam sau 2 tang thu muc (dung
        # trong git worktree, IndexError o checkout thang) -> ca file im lang khong
        # chay. Nay SHARED_ROOT hoi git "--git-common-dir" nen dung o moi vi tri.
        # PYTHONUTF8=1: test_manifest.py chay subprocess roi giai ma output; console
        # khong-UTF-8 (PowerShell gbk/cp936) se UnicodeDecodeError. Ep UTF-8 -> on dinh.
        $env:PYTHONUTF8 = '1'
        Push-Location "$goc"
        try {
            Buoc "pytest assetpack" {
                & $py -m pytest assetpack -q
            }
        } finally { Pop-Location }
    } else {
        Write-Host "`n[BO QUA] khong co python -> bo pytest assetpack" -ForegroundColor Yellow
    }
}

Write-Host "`n================ TONG KET ================" -ForegroundColor Cyan
if ($khongCoDB) {
    # F7: dung de "xanh gia" khi Postgres vang mat — test store/migration (code
    # chay MOI lan khoi dong) da bi SKIP. Bao DO noi bat, khong lang le xanh.
    Write-Host "!! CANH BAO: PostgreSQL KHONG chay -> test CSDL (store/migration) DA BI SKIP, CHUA duoc kiem !!" -ForegroundColor Red
    Write-Host "   Bat DB roi chay lai de kiem day du: D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start" -ForegroundColor Red
}
if ($loi.Count -eq 0) {
    if ($khongCoDB) { Write-Host "Cac buoc chay duoc: XANH (nhung THIEU test CSDL — xem canh bao tren)." -ForegroundColor Yellow }
    else { Write-Host "TAT CA XANH." -ForegroundColor Green }
    exit 0
} else {
    Write-Host "CO $($loi.Count) BUOC DO: $($loi -join ', ')" -ForegroundColor Red
    exit 1
}
