@ECHO OFF
TITLE Ngoc Rong Kame - Buoc 4 - Tu Dong Bat May Chu Sau Khi Khoi Dong Lai
SETLOCAL EnableExtensions DisableDelayedExpansion
cd /d "%~dp0"

REM ===========================================================================
REM  NGOC RONG KAME - BUOC 4: TU DONG BAT MAY CHU SAU KHI VPS KHOI DONG LAI
REM  ---------------------------------------------------------------------
REM  CHI CHAY FILE NAY SAU KHI may chu game da chay ngon lanh bang tay.
REM  Tu dong hoa mot thu chua chay duoc chi lam kho tim loi hon.
REM
REM  ---------------------------------------------------------------------
REM  VI SAO PHAI LAM RIENG MOT BUOC PHUC TAP THE NAY
REM  ---------------------------------------------------------------------
REM  Apache va MySQL da duoc cai thanh dich vu Windows o buoc 1, chung tu
REM  len sau khi khoi dong lai. Xong.
REM
REM  May chu game thi KHONG lam the duoc. Ham main() cua no
REM  (ServerManager.java dong 167) chi mo mot cua so do hoa; toan bo phan
REM  khoi dong nam trong ham dung cua ServerManagerUI. Noi cach khac:
REM  MAY CHU GAME BAT BUOC PHAI CO MAN HINH.
REM
REM  Dich vu Windows chay o "session 0" - noi khong he co man hinh. Nen
REM  khong the bien may chu game thanh dich vu. Day la han che that cua
REM  phan mem, khong phai lua chon.
REM
REM  Hau qua: sau khi VPS khoi dong lai, neu KHONG CO AI DANG NHAP thi
REM  khong co man hinh nao ca -> may chu game khong the len.
REM
REM  Cach duy nhat chay duoc, gom 2 manh:
REM    MANH 1 - Tu dong dang nhap (AutoAdminLogon): bao Windows tu dang
REM             nhap vao tai khoan Administrator ngay khi may khoi dong,
REM             de luon co san mot man hinh.
REM    MANH 2 - Cong viec theo lich: tu chay run.bat khi dang nhap xong.
REM
REM  File nay lam duoc MANH 2 tron ven, va lam gan xong MANH 1 - tru dung
REM  buoc go mat khau Windows, viec do ban phai tu lam. Ly do ghi ro o duoi.
REM ===========================================================================

set "TEN_CV=NgocRongKame-MayChu"
set "GAME_ROOT=C:\Users\Administrator\Music\zz1"

cls
echo.
echo ===============================================================
echo    BUOC 4: TU DONG BAT MAY CHU SAU KHI VPS KHOI DONG LAI
echo ===============================================================
echo.

REM --- Kiem tra quyen ---
net session >nul 2>&1
if not errorlevel 1 goto CO_QUYEN
echo   DUNG LAI: phai chay bang quyen Administrator.
echo   Bam chuot phai vao file nay -^> "Run as administrator".
echo.
pause
exit /b 1
:CO_QUYEN

REM --- Kiem tra thu muc game ---
if exist "%GAME_ROOT%\run.bat" goto CO_GAME
echo   DUNG LAI: khong thay "%GAME_ROOT%\run.bat"
echo.
echo   Neu ban dat thu muc game o cho khac, hay mo file nay bang
echo   Notepad va sua dong  set "GAME_ROOT=..."  o dau file.
echo.
pause
exit /b 1
:CO_GAME

echo   Thu muc game: %GAME_ROOT%
echo.


REM ===========================================================================
REM  MANH 2 - TAO CONG VIEC THEO LICH (phan de, lam truoc)
REM ===========================================================================
echo ---------------------------------------------------------------
echo   MANH 2/2: Tao cong viec tu chay may chu khi dang nhap
echo ---------------------------------------------------------------
echo.

REM  VI SAO DUNG "CONG VIEC THEO LICH" MA KHONG PHAI LOI TAT TRONG STARTUP:
REM    Loi tat trong thu muc Startup KHONG tu chay bang quyen Administrator
REM    duoc - Windows se hien hop thoai hoi "Co cho phep khong?" va dung do
REM    cho ai bam. Khong ai bam thi may chu khong len.
REM
REM    Ma may chu game CAN quyen Administrator: chuc nang tu chan IP pha
REM    hoai goi lenh netsh cua Windows. Khong du quyen thi lenh do THAT BAI
REM    LANG LE - nhat ky van in "[BAN IP VINH VIEN WINDOWS]" nhu that, nhung
REM    tuong lua khong chan ai ca. Ban se tuong minh co lop chong pha hoai
REM    trong khi thuc te khong co.
REM
REM    Cong viec theo lich co o "Chay voi quyen cao nhat" - dung quyen
REM    Administrator, khong hoi han gi.

REM  VI SAO TRO THANG VAO run.bat MA KHONG GOI java.exe:
REM    run.bat co dong  cd /d "%~dp0"  de dat dung thu muc lam viec.
REM    Moi duong dan trong ma Java deu la duong dan TUONG DOI. Sai thu muc
REM    lam viec thi may chu khong doc duoc config.properties (roi ve mat
REM    khau rong), khong nap duoc ban do. Trieu chung luc do RAT GIONG loi
REM    co so du lieu nen cuc ky de chan doan nham.
REM    run.bat cung co san vong canh gac: ma thoat 2 = bao tri chay lai
REM    ngay, 0 = dung han, khac = loi thi tu chay lai sau 10 giay.

schtasks /query /TN "%TEN_CV%" >nul 2>&1
if errorlevel 1 goto TAO_CV
echo   Cong viec "%TEN_CV%" da co tu truoc. Dang tao lai cho moi...
schtasks /delete /TN "%TEN_CV%" /F >nul 2>&1

:TAO_CV
REM Dung %USERNAME% thay vi go cung "Administrator": neu nha cung cap dat ten
REM tai khoan khac thi cong viec se dang ky cho mot tai khoan khong ton tai
REM va khong bao gio chay.
schtasks /create /TN "%TEN_CV%" /TR "\"%GAME_ROOT%\run.bat\"" /SC ONLOGON /RU "%USERNAME%" /RL HIGHEST /F >nul 2>&1
if errorlevel 1 goto CV_LOI

REM ---------------------------------------------------------------------------
REM  GO GIOI HAN 72 GIO - BAT BUOC, KHONG DUOC BO
REM
REM  Cong viec tao bang "schtasks /create" mac dinh bi dat ExecutionTimeLimit
REM  = PT72H. Nghia la sau dung 3 ngay chay lien tuc, Windows TU GIET tien
REM  trinh may chu game. Vong canh gac trong run.bat KHONG cuu duoc, vi day
REM  la Windows giet tu ben ngoai chu khong phai server tu thoat.
REM
REM  Trieu chung neu quen: may chu chay ngon 3 ngay roi dot ngot tat, nguoi
REM  choi rot het, va rat kho doan ra nguyen nhan vi no chi xay ra 3 ngay
REM  mot lan. Dat gioi han ve 0 = chay khong gioi han thoi gian.
REM ---------------------------------------------------------------------------
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "try { Set-ScheduledTask -TaskName '%TEN_CV%' -Settings (New-ScheduledTaskSettingsSet -ExecutionTimeLimit ([TimeSpan]::Zero) -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -RestartCount 3 -RestartInterval ([TimeSpan]::FromMinutes(1))) | Out-Null; 'OK' } catch { 'LOI: ' + $_.Exception.Message }"

REM Kiem chung lai: doc nguoc ra xem gioi han da la 0 chua.
for /f "delims=" %%L in ('powershell -NoProfile -Command "$t=Get-ScheduledTask -TaskName '%TEN_CV%' -EA 0; if ($t) { $t.Settings.ExecutionTimeLimit } else { 'KHONG-THAY' }"') do set "GIOI_HAN=%%L"
if /i "%GIOI_HAN%"=="PT0S" (
    echo       - Da go gioi han 72 gio ^(may chu chay khong gioi han thoi gian^)
) else (
    echo       [CANH BAO] Chua go duoc gioi han thoi gian chay. Hien la: %GIOI_HAN%
    echo                  Neu de nguyen, may chu game se tu tat sau 3 ngay.
    echo                  Mo Task Scheduler ^> tim "%TEN_CV%" ^> tab Settings ^>
    echo                  BO TICK "Stop the task if it runs longer than".
)

echo   [OK] Da tao cong viec "%TEN_CV%".
echo        Tu chay %GAME_ROOT%\run.bat moi khi dang nhap,
echo        voi quyen Administrator day du.
echo.
echo   Muon xem hoac xoa sau nay: bam Start, go "Task Scheduler".
goto XONG_CV

:CV_LOI
echo   [LOI] Khong tao duoc cong viec theo lich.
echo.
echo   CACH LAM TAY:
echo     1. Bam Start, go: Task Scheduler
echo     2. Ben phai bam "Create Task..." (khong phai Create Basic Task)
echo     3. The General : dat ten NgocRongKame-MayChu
echo                      tick "Run with highest privileges"
echo     4. The Triggers: New -^> Begin the task: "At log on"
echo     5. The Actions : New -^> Program: %GAME_ROOT%\run.bat
echo     6. The Settings: BO TICK "Stop the task if it runs longer than..."
echo        (khong bo thi Windows tat may chu sau 3 ngay chay lien tuc)
echo.

:XONG_CV
echo.


REM ===========================================================================
REM  MANH 1 - TU DONG DANG NHAP
REM ===========================================================================
echo ---------------------------------------------------------------
echo   MANH 1/2: Tu dong dang nhap vao Windows
echo ---------------------------------------------------------------
echo.
echo   ***********************************************************
echo   *  DOC KY CHO NAY TRUOC KHI QUYET DINH                    *
echo   ***********************************************************
echo.
echo   De co san mot man hinh cho may chu game, Windows phai TU DANG
echo   NHAP vao tai khoan Administrator sau khi khoi dong.
echo.
echo   CAI GIA PHAI TRA - ban can biet truoc:
echo.
echo   Cach thong thuong (sua registry) luu MAT KHAU WINDOWS cua ban
echo   duoi dang CHU THUONG, doc duoc bang mat, tai khoa:
echo     HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon
echo   Bat ky ai vao duoc may la doc duoc mat khau Administrator.
echo.
echo   May ban DA duoc siet phan nay roi: cong Remote Desktop da doi tu
echo   3389 sang 43389, va chi cho phep dai IP 14.248.0.0/16 ket noi.
echo   Nen nguy co nguoi ngoai vao duoc may la THAP - nhung khong phai
echo   bang khong, vi vay van phai can nhac truoc khi bat tu dang nhap.
echo.
echo   CO CACH AN TOAN HON, va toi khuyen dung cach nay:
echo     Cong cu "Autologon" cua Microsoft (bo Sysinternals). No luu mat
echo     khau vao kho ma hoa cua Windows (LSA secret) chu khong phai chu
echo     thuong trong registry.
echo     Tai tai:  https://learn.microsoft.com/sysinternals/downloads/autologon
echo     Chay len, go ten tai khoan + mat khau, bam Enable. Xong.
echo.
echo   ***********************************************************
echo.
echo   TOI SE KHONG TU GO MAT KHAU WINDOWS CUA BAN VAO DAY.
echo   Mat khau la cua rieng ban - ban tu go, o cong cu ban chon.
echo.
echo   Script chi chuan bi san cac phan KHONG phai mat khau, roi
echo   huong dan ban lam not buoc cuoi.
echo.

set "TRA_LOI="
set /p "TRA_LOI=  Ban co muon chuan bi san phan khong-phai-mat-khau khong? (co/khong): "
if /i "%TRA_LOI%"=="co" goto CHUAN_BI_AUTOLOGON

echo.
echo   Da bo qua phan tu dong dang nhap.
echo.
echo   LUU Y: cong viec theo lich o Manh 2 van da tao xong, nhung no chi
echo   chay KHI CO NGUOI DANG NHAP. Chua bat tu dong dang nhap thi sau moi
echo   lan VPS khoi dong lai, ban phai tu vao Remote Desktop mot lan de
echo   may chu game len.
echo.
goto KET_THUC

:CHUAN_BI_AUTOLOGON
echo.
echo   Dang dat cac gia tri khong phai mat khau...

reg add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v AutoAdminLogon /t REG_SZ /d 1 /f >nul 2>&1
if errorlevel 1 goto AUTOLOGON_LOI
reg add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v DefaultUserName /t REG_SZ /d "Administrator" /f >nul 2>&1
reg add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v DefaultDomainName /t REG_SZ /d "%COMPUTERNAME%" /f >nul 2>&1

echo   [OK] Da dat AutoAdminLogon, DefaultUserName, DefaultDomainName.
echo.
echo   CON DUNG MOT VIEC - BAN TU LAM:
echo.
echo   CACH AN TOAN HON (nen dung):
echo     1. Mo trinh duyet tren VPS, tai cong cu Autologon cua Microsoft:
echo        https://learn.microsoft.com/sysinternals/downloads/autologon
echo     2. Giai nen, chay Autologon64.exe
echo     3. Username: Administrator
echo        Domain  : %COMPUTERNAME%
echo        Password: (mat khau Windows cua ban)
echo     4. Bam Enable.
echo.
echo   CACH DON GIAN HON NHUNG KEM AN TOAN (mat khau nam dang chu thuong):
echo     1. Bam Start, go: regedit
echo     2. Di toi:
echo        HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon
echo     3. Bam chuot phai khoang trong ben phai -^> New -^> String Value
echo        Dat ten: DefaultPassword
echo     4. Bam dup vao no, go mat khau Windows cua ban, bam OK.
echo.
goto KET_THUC

:AUTOLOGON_LOI
echo   [LOI] Khong ghi duoc vao registry. Ban co chac dang chay bang
echo         quyen Administrator khong?
echo.

:KET_THUC
echo ===============================================================
echo    XONG BUOC 4
echo ===============================================================
echo.
echo   CACH THU XEM CO THAT SU CHAY KHONG:
echo.
echo   1. Khoi dong lai VPS (bam Start -^> nut nguon -^> Restart).
echo      Ban se bi mat ket noi Remote Desktop, day la binh thuong.
echo.
echo   2. Doi khoang 3-5 phut cho may len het.
echo.
echo   3. TU DIEN THOAI (tat wifi, dung 4G), mo http://nrokame.com
echo      - Web len duoc = Apache va MySQL da tu chay. Tot.
echo.
echo   4. Vao lai Remote Desktop. Phai thay cua so den cua may chu
echo      game DANG MO SAN. Neu thay tuc la moi thu da tron.
echo.
echo   ***********************************************************
echo   *  MOT THOI QUEN PHAI NHO TU HOM NAY                      *
echo   ***********************************************************
echo.
echo   Khi roi khoi VPS, LUON bam dau X tren thanh xanh o dinh man
echo   hinh (goi la "ngat ket noi" - may chu van chay tiep).
echo.
echo   TUYET DOI KHONG bam nut "Sign out" hay "Dang xuat".
echo   Sign out la huy han phien lam viec -^> cua so may chu game
echo   bi dong theo -^> NGUOI CHOI ROT HET.
echo.
echo   Day la mat trai cua viec may chu bat buoc phai co man hinh.
echo   ***********************************************************
echo.
pause
exit /b 0
