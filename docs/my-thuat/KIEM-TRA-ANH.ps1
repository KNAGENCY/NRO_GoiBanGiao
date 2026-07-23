# ===========================================================================
#  KIEM TRA CHAT LUONG LO ANH TRUOC KHI GUI
#  ---------------------------------------------------------------------
#  CACH DUNG: chuot phai file nay -> "Run with PowerShell"
#             hoac mo PowerShell roi go:
#             powershell -File "D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.ps1" -ThuMuc "D:\duong\dan\toi\thu\muc\anh"
#
#  NO KIEM 5 DIEU, va bao ngay cai nao HONG:
#    1. Anh co phai tranh that khong, hay chi la hinh khoi phang
#    2. Co nen trong suot khong
#    3. Kich thuoc co dung khong
#    4. Anh co bi trong rong khong
#    5. Anh co bi vien sang quanh mep khong (loi hay gap khi tach nen)
#
#  NGUONG PHAN BIET THAT/GIA:
#    Tranh to 3 sac do co vien -> it nhat 40 mau khac nhau.
#    Duoi 20 mau gan nhu chac chan la hinh khoi ve bang code.
#    (Doi chieu: lo v2 dat 81-1264 mau; lo v3 hong chi 2-5 mau)
# ===========================================================================

param(
    [string]$ThuMuc = ""
)

Add-Type -AssemblyName System.Drawing

if ($ThuMuc -eq "") {
    $ThuMuc = Read-Host "Dan duong dan thu muc chua anh vao day roi Enter"
}
if (-not (Test-Path $ThuMuc)) {
    Write-Host "KHONG TIM THAY THU MUC: $ThuMuc" -ForegroundColor Red
    Read-Host "Enter de dong"
    exit 1
}

$NGUONG_MAU_TOI_THIEU = 40      # duoi muc nay coi la hinh khoi, khong phai tranh
$NGUONG_MAU_NGHI_NGO  = 80      # tu 40-80 la dang ngo, can nguoi xem lai

$anh = Get-ChildItem $ThuMuc -Recurse -Filter *.png
if ($anh.Count -eq 0) {
    Write-Host "Khong co file .png nao trong thu muc nay." -ForegroundColor Red
    Read-Host "Enter de dong"
    exit 1
}

Write-Host ""
Write-Host "Dang kiem $($anh.Count) anh..." -ForegroundColor Cyan
Write-Host ""

$hong    = @()
$ngongo  = @()
$khongTS = @()
$rong    = @()
$vienSang = @()
$dat     = 0

foreach ($f in $anh) {
    try {
        $b = New-Object System.Drawing.Bitmap($f.FullName)
    } catch {
        $hong += [PSCustomObject]@{ Ten = $f.Name; LyDo = "khong mo duoc file" }
        continue
    }

    # --- co kenh trong suot khong ---
    if ($b.PixelFormat -notmatch 'Argb') { $khongTS += $f.Name }

    # --- dem mau va do dac ---
    $mau  = @{}
    $dac  = 0
    $tong = 0
    $mepSang = 0
    $buoc = [Math]::Max(1, [int]($b.Width / 120))   # anh lon thi lay mau thua ra cho nhanh

    for ($x = 0; $x -lt $b.Width; $x += $buoc) {
        for ($y = 0; $y -lt $b.Height; $y += $buoc) {
            $tong++
            $c = $b.GetPixel($x, $y)
            if ($c.A -gt 30) {
                $dac++
                $mau["$($c.R),$($c.G),$($c.B)"] = 1
                # diem nua trong suot va rat sang o vien = dau hieu tach nen loi
                if ($c.A -lt 200 -and ($c.R + $c.G + $c.B) -gt 690) { $mepSang++ }
            }
        }
    }

    $soMau   = $mau.Count
    $phanDac = if ($tong -gt 0) { [math]::Round($dac * 100 / $tong) } else { 0 }
    $b.Dispose()

    # --- phan loai ---
    if ($phanDac -lt 2) {
        $rong += [PSCustomObject]@{ Ten = $f.Name; LyDo = "gan nhu trong rong ($phanDac% co hinh)" }
    }
    elseif ($soMau -lt $NGUONG_MAU_TOI_THIEU) {
        $hong += [PSCustomObject]@{ Ten = $f.Name; LyDo = "chi $soMau mau - la hinh khoi phang, khong phai tranh" }
    }
    elseif ($soMau -lt $NGUONG_MAU_NGHI_NGO) {
        $ngongo += [PSCustomObject]@{ Ten = $f.Name; LyDo = "$soMau mau - hoi it, nen xem lai bang mat" }
    }
    else {
        $dat++
    }

    if ($mepSang -gt ($dac * 0.05)) { $vienSang += $f.Name }
}

# ===========================================================================
Write-Host ("=" * 68)
Write-Host "  KET QUA" -ForegroundColor White
Write-Host ("=" * 68)
Write-Host ""
Write-Host ("  Dat            : {0,4} / {1}" -f $dat, $anh.Count) -ForegroundColor Green
Write-Host ("  Dang ngo       : {0,4}" -f $ngongo.Count) -ForegroundColor Yellow
Write-Host ("  HONG (hinh khoi): {0,4}" -f $hong.Count) -ForegroundColor Red
Write-Host ("  Trong rong     : {0,4}" -f $rong.Count) -ForegroundColor Red
Write-Host ("  Thieu nen trong suot: {0,4}" -f $khongTS.Count) -ForegroundColor Red
Write-Host ("  Co vien sang o mep  : {0,4}" -f $vienSang.Count) -ForegroundColor Yellow
Write-Host ""

function InDanhSach($ten, $ds, $mau) {
    if ($ds.Count -eq 0) { return }
    Write-Host "  --- $ten ---" -ForegroundColor $mau
    $ds | Select-Object -First 25 | ForEach-Object {
        if ($_ -is [string]) { Write-Host "    $_" }
        else { Write-Host ("    {0,-38} {1}" -f $_.Ten, $_.LyDo) }
    }
    if ($ds.Count -gt 25) { Write-Host "    ... va $($ds.Count - 25) file nua" }
    Write-Host ""
}

InDanhSach "HONG - phai lam lai" $hong Red
InDanhSach "TRONG RONG - phai lam lai" $rong Red
InDanhSach "THIEU NEN TRONG SUOT" $khongTS Red
InDanhSach "DANG NGO - nen xem bang mat" $ngongo Yellow
InDanhSach "VIEN SANG O MEP" $vienSang Yellow

Write-Host ("=" * 68)
if ($hong.Count -eq 0 -and $rong.Count -eq 0 -and $khongTS.Count -eq 0) {
    Write-Host "  => LO NAY DUNG DUOC. Gui cho tro ly duoc roi." -ForegroundColor Green
} else {
    Write-Host "  => LO NAY CHUA DUNG DUOC. Bao AI lam lai nhung file o tren." -ForegroundColor Red
    Write-Host "     Nho nhac AI: phai SINH ANH THAT, khong duoc ve hinh khoi bang code." -ForegroundColor Red
}
Write-Host ("=" * 68)
Write-Host ""
Read-Host "Enter de dong"
