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
if ($cong) {
    $env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:$cong/nro_test"
    Write-Host "PostgreSQL o cong $cong -> NRO_TEST_DSN da dat (test can DB SE chay)" -ForegroundColor DarkGray
} else {
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

    # gofmt: file cu co dau dong CRLF bi gofmt bao nhieu (§9 BAT-DAU) — KHONG phai
    # loi format that. Loc ra danh sach nhieu da biet; con lai (file MOI sai format)
    # -> DO. Sua goc CRLF (gitattributes) la viec rieng sau.
    Write-Host "`n=== gofmt ===" -ForegroundColor Cyan
    # Nhieu CRLF da biet (§9 BAT-DAU): file cu dau dong CRLF, gofmt bao nhung KHONG
    # phai loi format. So sanh theo DUONG DAN (chuan hoa /), khong theo ten file —
    # de cmd/nroserver/main.go van bi soi (khac cmd/thuclient/main.go trong list).
    # File MOI sai format se khong nam trong list -> DO. Bo list khi chuan hoa gitattributes.
    $nhieuCRLF = @(
        'internal/auth/password_test.go',
        'internal/net/listener.go',
        'internal/protocol/frame.go',
        'internal/protocol/frame_test.go',
        'internal/protocol/integration_test.go',
        'internal/protocol/keystream.go',
        'internal/protocol/keystream_test.go',
        'internal/protocol/message.go',
        'cmd/thuclient/main.go'
    )
    if ($gofmt) {
        $ra = & $gofmt -l internal cmd 2>&1
        $xau = @($ra | Where-Object { $_ } | ForEach-Object { ($_ -replace '\\', '/') } | Where-Object { $_ -notin $nhieuCRLF })
        if ($xau.Count -gt 0) {
            Write-Host "[DO] gofmt: file sai format (ngoai nhieu CRLF cu):" -ForegroundColor Red
            $xau | ForEach-Object { Write-Host "     $_" -ForegroundColor Red }
            $loi.Add("gofmt")
        } else {
            Write-Host "[OK] gofmt (chi con nhieu CRLF cu da biet)" -ForegroundColor Green
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
        # Bo test_task6b_canonical_preflight.py: no gia dinh layout git-worktree
        # (SHARED_ROOT = ROOT.parents[1]) nen vo tren checkout chinh — bug path co san
        # cua test do (khong lien quan code); sua rieng sau.
        # PYTHONUTF8=1: test_manifest.py chay subprocess roi giai ma output; console
        # khong-UTF-8 (PowerShell gbk/cp936) se UnicodeDecodeError. Ep UTF-8 -> on dinh.
        $env:PYTHONUTF8 = '1'
        Push-Location "$goc"
        try {
            Buoc "pytest assetpack" {
                & $py -m pytest assetpack -q --ignore=assetpack/tests/test_task6b_canonical_preflight.py
            }
        } finally { Pop-Location }
    } else {
        Write-Host "`n[BO QUA] khong co python -> bo pytest assetpack" -ForegroundColor Yellow
    }
}

Write-Host "`n================ TONG KET ================" -ForegroundColor Cyan
if ($loi.Count -eq 0) {
    Write-Host "TAT CA XANH." -ForegroundColor Green
    exit 0
} else {
    Write-Host "CO $($loi.Count) BUOC DO: $($loi -join ', ')" -ForegroundColor Red
    exit 1
}
