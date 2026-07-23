@ECHO OFF
TITLE Ngoc Rong Kame - Go Ra Lam Lai
SETLOCAL EnableExtensions DisableDelayedExpansion
cd /d "%~dp0"

REM ===========================================================================
REM  NGOC RONG KAME - GO RA LAM LAI (DUONG LUI)
REM  ---------------------------------------------------------------------
REM  Dung file nay khi cai hong va muon lam lai tu dau tren VPS.
REM
REM  No THAO NGUOC nhung gi cac file 1 va 4 da lam:
REM     - go cong viec theo lich tu bat may chu
REM     - tat tu dong dang nhap
REM     - go Apache va MySQL khoi danh sach dich vu Windows
REM     - xoa cac luat tuong lua ten "Kame - ..."
REM
REM  NO KHONG DUNG DEN:
REM     - du lieu game (thu muc zz1)
REM     - thu muc web (htdocs)
REM     - co so du lieu "z"
REM     - JDK va XAMPP (van cai nguyen)
REM  Muon xoa nhung thu do thi tu tay xoa - toi co y khong lam, vi xoa
REM  nham la mat sach du lieu nguoi choi.
REM
REM  LUU Y VE LUAT TUONG LUA 3389 (Remote Desktop):
REM  Script nay CO Y KHONG xoa luat "Kame - Remote Desktop 3389".
REM  Xoa no di co the day ban ra khoi VPS va khong vao lai duoc.
REM ===========================================================================

cls
echo.
echo ===============================================================
echo    NGOC RONG KAME - GO RA LAM LAI
echo ===============================================================
echo.

net session >nul 2>&1
if not errorlevel 1 goto CO_QUYEN
echo   DUNG LAI: phai chay bang quyen Administrator.
echo.
pause
exit /b 1
:CO_QUYEN

echo   Script nay se thao nguoc phan CAI DAT, nhung KHONG xoa
echo   du lieu game, khong xoa web, khong xoa co so du lieu.
echo.
set "XN="
set /p "XN=  Go: DONG Y  de lam tiep (hoac Enter de thoat): "
if /i not "%XN%"=="DONG Y" goto THOAT

echo.
echo   [1/4] Go cong viec tu bat may chu game...
schtasks /delete /TN "NgocRongKame-MayChu" /F >nul 2>&1
if errorlevel 1 echo         (khong co cong viec nao - bo qua)
if not errorlevel 1 echo         Da go.

echo   [2/4] Tat tu dong dang nhap...
reg add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v AutoAdminLogon /t REG_SZ /d 0 /f >nul 2>&1
REM  Xoa han khoa chua mat khau Windows - day la viec NEN LAM du co go
REM  cai dat hay khong, vi no luu mat khau duoi dang chu thuong.
reg delete "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v DefaultPassword /f >nul 2>&1
echo         Da tat, va da xoa khoa DefaultPassword trong registry.

echo   [3/4] Go Apache va MySQL khoi danh sach dich vu Windows...
net stop "Apache2.4" >nul 2>&1
net stop "mysql" >nul 2>&1
if exist "C:\xampp\apache\apache_uninstallservice.bat" call "C:\xampp\apache\apache_uninstallservice.bat" >nul 2>&1
if exist "C:\xampp\mysql\mysql_uninstallservice.bat" call "C:\xampp\mysql\mysql_uninstallservice.bat" >nul 2>&1
echo         Da go. (XAMPP van cai nguyen, chi khong tu chay nua.)

echo   [4/4] Xoa cac luat tuong lua ten "Kame - ..."...
REM  CO Y KHONG xoa luat 3389 - xoa la co the mat duong vao VPS.
netsh advfirewall firewall delete rule name="Kame - Web HTTP 80" >nul 2>&1
netsh advfirewall firewall delete rule name="Kame - Game 14445" >nul 2>&1
netsh advfirewall firewall delete rule name="Kame - Chan MySQL 3306" >nul 2>&1
echo         Da xoa luat cong 80, 14445 va luat chan 3306.
echo         GIU LAI luat 3389 (Remote Desktop) de ban khong mat duong vao VPS.

echo.
echo ===============================================================
echo    DA GO XONG
echo ===============================================================
echo.
echo   Nhung thu VAN CON tren may (co y giu lai):
echo     - JDK 21 va XAMPP
echo     - Thu muc game va thu muc web
echo     - Co so du lieu "z" cung toan bo du lieu nguoi choi
echo     - Luat tuong lua cho Remote Desktop
echo.
echo   MUON LAM LAI TU DAU: chay lai 1-CAI-MOI-TRUONG.bat.
echo.
echo   MUON KHOI PHUC MOT FILE CAU HINH DA BI SUA:
echo   moi file bi sua deu co ban goc nam ngay canh no, ten giong het
echo   nhung them duoi .goc . Vi du:
echo     C:\xampp\htdocs\sv1\config.php.goc
echo   Chi can doi ten bo duoi .goc di la ve nhu cu.
echo.
echo   MUON XOA HAN CO SO DU LIEU de nhap lai tu dau (MAT HET DU LIEU
echo   NGUOI CHOI - can nhac ky):
echo     C:\xampp\mysql\bin\mysql.exe -u root -p -e "DROP DATABASE z;"
echo.
pause
exit /b 0

:THOAT
echo.
echo   Da huy, khong thay doi gi.
echo.
pause
exit /b 0
