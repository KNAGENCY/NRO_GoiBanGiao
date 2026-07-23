@ECHO OFF
TITLE Ngoc Rong Kame - Buoc 1 - Cai Moi Truong Cho VPS
SETLOCAL EnableExtensions DisableDelayedExpansion
cd /d "%~dp0"

REM ===========================================================================
REM  NGOC RONG KAME - BUOC 1: CAI MOI TRUONG CHO MAY CHU VPS
REM  ---------------------------------------------------------------------
REM  File nay chay TREN VPS (Windows Server). Bam dup chuot la chay.
REM
REM  No lo phan NEN MONG, chua dung gi den game:
REM     1. Kiem tra quyen Administrator
REM     2. Kiem tra o C: con du cho trong
REM     3. Cai JDK 21 (Eclipse Temurin) - de chay may chu game
REM     4. Cai XAMPP 8.2 (Apache + PHP 8.2 + MariaDB) - de chay web
REM     5. Don rac mac dinh cua XAMPP trong thu muc htdocs
REM     6. Cai Apache va MySQL thanh DICH VU WINDOWS (tu len sau khi reboot)
REM     7. Dat bien moi truong JAVA_HOME va them vao PATH
REM     8. Mo tuong lua cong 80 (web) va 14445 (game). KHONG mo 3306.
REM     9. Dat mui gio ve gio Viet Nam
REM
REM  CHAY LAI NHIEU LAN DUOC. Cai gi co roi thi tu dong bo qua.
REM
REM  File nay KHONG chua mat khau nao het, va cung khong hoi mat khau.
REM  Viec dat mat khau CSDL la cua buoc sau (file 3-CAI-GAME-LEN-VPS.ps1).
REM
REM  THU TU DUNG: chay file NAY TRUOC, roi moi chep du lieu game va web len.
REM  Lam nguoc lai thi bo cai XAMPP se de len thu muc htdocs va an mat
REM  trang index.php that cua ban.
REM ===========================================================================

REM --- Bien theo doi ket qua tung buoc, in ra bang tong ket o cuoi ---
set "KQ_QUYEN=CHUA CHAY"
set "KQ_ODIA=CHUA CHAY"
set "KQ_JDK=CHUA CHAY"
set "KQ_XAMPP=CHUA CHAY"
set "KQ_DONRAC=CHUA CHAY"
set "KQ_DICHVU=CHUA CHAY"
set "KQ_BIENMT=CHUA CHAY"
set "KQ_TUONGLUA=CHUA CHAY"
set "KQ_MUIGIO=CHUA CHAY"
set "CO_LOI=0"

set "THU_MUC_TAI=%~dp0tai-ve"
set "FILE_LOG=%~dp0ket-qua-buoc-1.txt"

cls
echo.
echo ===============================================================
echo    NGOC RONG KAME - BUOC 1: CAI MOI TRUONG CHO VPS
echo ===============================================================
echo.
echo  Script nay se cai: JDK 21, XAMPP 8.2, dat Apache + MySQL thanh
echo  dich vu Windows, mo cong tuong lua, dat mui gio Viet Nam.
echo.
echo  Can tai ve khoang 340 MB. Tuy toc do mang, co the mat 5-20 phut.
echo  KHONG dong cua so nay giua chung.
echo.
echo  CHAY FILE NAY TRUOC khi chep du lieu game va web len VPS.
echo.
echo ===============================================================
echo.


REM ===========================================================================
REM  BUOC 1/9 - KIEM TRA QUYEN ADMINISTRATOR
REM ===========================================================================
echo [1/9] Dang kiem tra quyen Administrator...

net session >nul 2>&1
if not errorlevel 1 goto CO_QUYEN

echo.
echo   ***********************************************************
echo   *  DUNG LAI: KHONG DU QUYEN                               *
echo   ***********************************************************
echo.
echo   Script nay phai chay bang quyen Administrator thi moi cai
echo   duoc phan mem va sua duoc tuong lua.
echo.
echo   CACH LAM DUNG:
echo     1. Dong cua so nay lai.
echo     2. Bam CHUOT PHAI vao file 1-CAI-MOI-TRUONG.bat
echo     3. Chon dong "Run as administrator"
echo        (tren ban Windows tieng Viet la "Chay bang quyen quan tri")
echo.
echo   ***********************************************************
echo.
pause
exit /b 1

:CO_QUYEN
set "KQ_QUYEN=OK - dang chay quyen Administrator"
echo       OK - dang chay bang quyen Administrator.
echo.


REM ===========================================================================
REM  BUOC 2/9 - KIEM TRA DUNG LUONG O C:
REM ===========================================================================
echo [2/9] Dang kiem tra dung luong trong cua o C:...

set "CON_TRONG_GB="
for /f "usebackq delims=" %%A in (`powershell -NoProfile -Command "$d = Get-PSDrive C -ErrorAction SilentlyContinue; if ($d) { [math]::Floor($d.Free/1GB) } else { -1 }"`) do set "CON_TRONG_GB=%%A"

if not defined CON_TRONG_GB goto ODIA_KHONG_DOC_DUOC
if "%CON_TRONG_GB%"=="-1" goto ODIA_KHONG_DOC_DUOC

echo       O C: con trong %CON_TRONG_GB% GB.

if %CON_TRONG_GB% GEQ 10 goto ODIA_DU_CHO

echo.
echo   ***********************************************************
echo   *  DUNG LAI: O C: SAP DAY                                 *
echo   ***********************************************************
echo.
echo   Con trong: %CON_TRONG_GB% GB. Can it nhat 10 GB de cai tiep.
echo.
echo   Rieng buoc nay da can khoang 2.5 GB (JDK + XAMPP + file tai ve).
echo   Sau do con phai chep len 1 GB du lieu game va kho anh 822 MB.
echo.
echo   CACH XU LY:
echo     - Xoa bot file trong C:\Users\Administrator\Downloads
echo     - Xoa thu muc rac: go %%TEMP%% vao thanh dia chi cua This PC
echo     - Hoac lien he noi ban VPS de nang dung luong o cung
echo.
echo   ***********************************************************
echo.
set "KQ_ODIA=THAT BAI - chi con %CON_TRONG_GB% GB, can it nhat 10 GB"
pause
exit /b 1

:ODIA_KHONG_DOC_DUOC
echo.
echo   [CANH BAO] Khong doc duoc dung luong o C:.
echo              Script van chay tiep, nhung ban hay tu kiem tra
echo              o C: con it nhat 10 GB trong truoc khi cai.
echo.
set "KQ_ODIA=CANH BAO - khong doc duoc dung luong, hay tu kiem tra"
set "CO_LOI=1"
goto XONG_ODIA

:ODIA_DU_CHO
set "KQ_ODIA=OK - con trong %CON_TRONG_GB% GB"
echo       OK - du cho de cai.

:XONG_ODIA
echo.

REM --- Tao thu muc chua file tai ve ---
if not exist "%THU_MUC_TAI%" mkdir "%THU_MUC_TAI%" >nul 2>&1
if exist "%THU_MUC_TAI%" goto CO_THU_MUC_TAI
echo       [CANH BAO] Khong tao duoc thu muc tai ve canh script.
echo                  Chuyen sang dung thu muc tam cua Windows.
set "THU_MUC_TAI=%TEMP%\kame-tai-ve"
if not exist "%THU_MUC_TAI%" mkdir "%THU_MUC_TAI%" >nul 2>&1
:CO_THU_MUC_TAI


REM ===========================================================================
REM  BUOC 3/9 - CAI JDK 21 (ECLIPSE TEMURIN)
REM ===========================================================================
echo [3/9] Dang kiem tra Java 21...

call :TIM_JAVA
if not defined JAVA_EXE goto CAN_CAI_JDK

echo       DA CO SAN Java 21 tai:
echo         %JAVA_EXE%
echo       Bo qua buoc cai dat.
set "KQ_JDK=DA CO SAN - bo qua"
goto XONG_JDK

:CAN_CAI_JDK
echo       Chua co Java 21 tren may nay. Bat dau cai...
echo.

REM --- Cach 1: dung winget neu may co san ---
REM  Luu y: Windows Server 2019 thuong KHONG co winget. Khong sao,
REM  script se tu chuyen sang cach 2 (tai file .msi ve cai).
where winget >nul 2>&1
if errorlevel 1 goto CAI_JDK_BANG_MSI

echo       Tim thay winget. Dang thu cai bang winget...
winget install --id EclipseAdoptium.Temurin.21.JDK --silent --accept-source-agreements --accept-package-agreements --disable-interactivity
echo.
echo       Kiem tra lai sau khi winget chay xong...
call :TIM_JAVA
if defined JAVA_EXE goto JDK_CAI_XONG
echo       winget khong cai duoc. Chuyen sang cach tai file .msi.
echo.

:CAI_JDK_BANG_MSI
set "FILE_MSI=%THU_MUC_TAI%\OpenJDK21-x64.msi"

if not exist "%FILE_MSI%" goto TAI_MSI
REM --- File da tai lan truoc, kiem tra xem co nguyen ven khong ---
call :KIEM_TRA_KICH_THUOC "%FILE_MSI%" 100000000
if "%FILE_DU_LON%"=="1" goto CAI_MSI
echo       File .msi tai lan truoc bi loi (qua nho). Xoa va tai lai.
del /f /q "%FILE_MSI%" >nul 2>&1

:TAI_MSI
echo       Dang tai JDK 21 tu trang chinh thuc adoptium.net (~180 MB)...
echo       Vui long cho, khong dong cua so.
call :TAI_FILE "https://api.adoptium.net/v3/installer/latest/21/ga/windows/x64/jdk/hotspot/normal/eclipse" "%FILE_MSI%"

call :KIEM_TRA_KICH_THUOC "%FILE_MSI%" 100000000
if "%FILE_DU_LON%"=="1" goto CAI_MSI

echo       Tai tu adoptium.net khong duoc. Thu duong dan du phong tren github...
del /f /q "%FILE_MSI%" >nul 2>&1
REM  CHU Y ky tu %%%%2B o duoi day KHONG phai go nham.
REM  Dia chi that chua "%%2B" (chinh la dau + da ma hoa). Nhung lenh CALL
REM  bat cac ky tu %% HAI LAN, nen phai viet gap doi len thi khi toi noi
REM  moi con dung "%%2B". Viet "%%2B" nhu binh thuong se thanh "B" va tai loi.
call :TAI_FILE "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.11%%%%2B10/OpenJDK21U-jdk_x64_windows_hotspot_21.0.11_10.msi" "%FILE_MSI%"

call :KIEM_TRA_KICH_THUOC "%FILE_MSI%" 100000000
if "%FILE_DU_LON%"=="1" goto CAI_MSI

echo.
echo       [THAT BAI] Khong tai duoc JDK 21.
echo                  Thuong la do VPS chua ra duoc internet, hoac
echo                  trang tai dang bao tri.
echo.
echo       CACH LAM TAY:
echo         1. Mo trinh duyet tren VPS, vao: https://adoptium.net/temurin/releases
echo         2. Chon: Version 21 - Operating System Windows - Architecture x64
echo            - Package Type JDK - dinh dang .msi
echo         3. Tai ve va bam dup de cai binh thuong.
echo         4. Chay lai file 1-CAI-MOI-TRUONG.bat nay.
echo.
set "KQ_JDK=THAT BAI - khong tai duoc file cai dat"
set "CO_LOI=1"
goto XONG_JDK

:CAI_MSI
echo       Dang cai JDK 21 (co the mat 1-3 phut, man hinh se dung yen)...
msiexec /i "%FILE_MSI%" /qn /norestart INSTALLDIR="C:\Program Files\Java\jdk-21" ADDLOCAL=FeatureMain,FeatureEnvironment,FeatureJavaHome /L*v "%THU_MUC_TAI%\nhat-ky-cai-jdk.log"
set "MA_MSI=%ERRORLEVEL%"

if "%MA_MSI%"=="0" goto MSI_XONG
if "%MA_MSI%"=="3010" goto MSI_XONG
echo       [CANH BAO] Bo cai JDK bao ma loi %MA_MSI%.
echo                  Chi tiet o: %THU_MUC_TAI%\nhat-ky-cai-jdk.log
echo                  Van kiem tra tiep xem Java co len khong...

:MSI_XONG
call :TIM_JAVA
if defined JAVA_EXE goto JDK_CAI_XONG

echo.
echo       [THAT BAI] Cai xong nhung van khong tim thay Java 21.
echo                  Xem nhat ky: %THU_MUC_TAI%\nhat-ky-cai-jdk.log
echo.
set "KQ_JDK=THAT BAI - cai xong nhung khong tim thay java.exe"
set "CO_LOI=1"
goto XONG_JDK

:JDK_CAI_XONG
echo       OK - da cai xong Java 21 tai:
echo         %JAVA_EXE%
set "KQ_JDK=OK - da cai moi"

:XONG_JDK
echo.


REM ===========================================================================
REM  BUOC 4/9 - CAI XAMPP 8.2 (APACHE + PHP 8.2 + MARIADB)
REM ===========================================================================
echo [4/9] Dang kiem tra XAMPP...

REM  QUAN TRONG: bat buoc cai vao dung C:\xampp.
REM  Ly do: file .htaccess cua web nap lop chan spam bang duong dan
REM  tuyet doi "C:/xampp/htdocs/ddos_guard.php". Cai sang o khac thi
REM  lop bao ve do IM LANG khong chay, khong bao loi gi het.

if not exist "C:\xampp\apache\bin\httpd.exe" goto CAN_CAI_XAMPP
echo       DA CO SAN XAMPP tai C:\xampp. Bo qua buoc cai dat.
set "KQ_XAMPP=DA CO SAN - bo qua"
goto XONG_XAMPP

:CAN_CAI_XAMPP
if not exist "C:\xampp" goto TAI_XAMPP
echo       [CANH BAO] Da co thu muc C:\xampp nhung thieu Apache.
echo                  Co the lan cai truoc bi dut giua chung.
echo                  Script se cai de len. Neu van loi, hay xoa han
echo                  thu muc C:\xampp roi chay lai file nay.
echo.

:TAI_XAMPP
set "FILE_XAMPP=%THU_MUC_TAI%\xampp-8.2.12-installer.exe"

if not exist "%FILE_XAMPP%" goto TAI_XAMPP_MOI
call :KIEM_TRA_KICH_THUOC "%FILE_XAMPP%" 100000000
if "%FILE_DU_LON%"=="1" goto CAI_XAMPP
echo       File cai XAMPP tai lan truoc bi loi (qua nho). Xoa va tai lai.
del /f /q "%FILE_XAMPP%" >nul 2>&1

:TAI_XAMPP_MOI
echo       Dang tai XAMPP 8.2.12 (PHP 8.2) tu sourceforge (~158 MB)...
echo       Vui long cho, khong dong cua so.
REM  CHU Y ky tu %%%%20 o duoi day KHONG phai go nham.
REM  Dia chi that chua "%%20" (chinh la dau cach da ma hoa, vi ten thu muc
REM  tren sourceforge la "XAMPP Windows" co khoang trang). Lenh CALL bat
REM  cac ky tu %% HAI LAN nen phai viet gap doi len. Neu viet "%%20" nhu
REM  binh thuong thi dia chi bien thanh "XAMPP0Windows" va tai that bai.
call :TAI_FILE "https://downloads.sourceforge.net/project/xampp/XAMPP%%%%20Windows/8.2.12/xampp-windows-x64-8.2.12-0-VS16-installer.exe" "%FILE_XAMPP%"

call :KIEM_TRA_KICH_THUOC "%FILE_XAMPP%" 100000000
if "%FILE_DU_LON%"=="1" goto CAI_XAMPP

echo.
echo       [THAT BAI] Khong tai duoc XAMPP.
echo.
echo       CACH LAM TAY:
echo         1. Mo trinh duyet tren VPS, vao: https://www.apachefriends.org/download.html
echo         2. Tai ban XAMPP for Windows co PHP 8.2 (64-bit)
echo         3. Bam dup de cai, va PHAI chon thu muc C:\xampp
echo            (dung doi sang o D: hay Program Files - se hong lop bao ve web)
echo         4. Chay lai file 1-CAI-MOI-TRUONG.bat nay.
echo.
set "KQ_XAMPP=THAT BAI - khong tai duoc file cai dat"
set "CO_LOI=1"
goto XONG_XAMPP

:CAI_XAMPP
echo       Dang cai XAMPP vao C:\xampp (mat 2-5 phut, man hinh se dung yen)...
"%FILE_XAMPP%" --mode unattended --unattendedmodeui minimal --prefix "C:\xampp"
set "MA_XAMPP=%ERRORLEVEL%"

if "%MA_XAMPP%"=="0" goto XAMPP_KIEM_TRA
echo       [CANH BAO] Bo cai XAMPP bao ma loi %MA_XAMPP%.
echo                  Van kiem tra tiep xem cai duoc den dau...

:XAMPP_KIEM_TRA
set "THIEU="
if not exist "C:\xampp\apache\bin\httpd.exe" set "THIEU=%THIEU% Apache"
if not exist "C:\xampp\php\php.exe" set "THIEU=%THIEU% PHP"
if not exist "C:\xampp\mysql\bin\mysqld.exe" set "THIEU=%THIEU% MariaDB"

if not defined THIEU goto XAMPP_DU

echo.
echo       [THAT BAI] XAMPP cai thieu:%THIEU%
echo                  Hay xoa thu muc C:\xampp roi cai tay tu
echo                  https://www.apachefriends.org/download.html
echo.
set "KQ_XAMPP=THAT BAI - cai thieu:%THIEU%"
set "CO_LOI=1"
goto XONG_XAMPP

:XAMPP_DU
echo       OK - da cai xong XAMPP tai C:\xampp
set "PHIEN_BAN_PHP="
REM  Co y KHONG dat dau nhay kep quanh duong dan php.exe: duong dan nay
REM  khong co khoang trang, ma them dau nhay vao trong 'for /f' se bi
REM  cmd cat mat va sinh loi. (Cung loi da gap o ham :THU_MOT_JAVA.)
for /f "tokens=2" %%P in ('C:\xampp\php\php.exe -v 2^>nul ^| findstr /i /b "PHP"') do if not defined PHIEN_BAN_PHP set "PHIEN_BAN_PHP=%%P"
if defined PHIEN_BAN_PHP echo         Phien ban PHP: %PHIEN_BAN_PHP%
set "KQ_XAMPP=OK - da cai moi vao C:\xampp"

:XONG_XAMPP
echo.


REM ===========================================================================
REM  BUOC 5/9 - DON RAC MAC DINH CUA XAMPP TRONG THU MUC HTDOCS
REM ===========================================================================
echo [5/9] Dang don rac mac dinh cua XAMPP trong C:\xampp\htdocs...

REM  VI SAO PHAI DON:
REM  Bo cai XAMPP tu tao san mot trang gioi thieu cua no trong htdocs:
REM  index.php (chuyen huong sang /dashboard), dashboard\, applications.html,
REM  bitnami.css, favicon.ico, va mot thu muc img\.
REM
REM  Hai cai hai:
REM   1) Thu muc img\ cua XAMPP se TRON LAN voi thu muc img\ cua web game.
REM   2) Trang http://nrokame.com/dashboard/ van mo duoc va khoe ra dung
REM      phien ban XAMPP + PHP dang chay - dung thu ke xau di do lo hong can.
REM
REM  Phai don TRUOC khi chep web that len. Neu chep web truoc roi moi cai
REM  XAMPP thi index.php nguy trang cua ban bi ghi de mat.

if not exist "C:\xampp\htdocs" goto DONRAC_KHONG_CO_THU_MUC

REM  AN TOAN - PHAI DUNG CA HAI DAU HIEU MOI DAM XOA index.php:
REM    1) Co chuoi "dashboard" (ban cua XAMPP chi la 1 dong chuyen huong
REM       sang /dashboard/). Da kiem chung: index.php that cua ban KHONG
REM       he chua chu nay - 0 lan xuat hien.
REM    2) Nho hon 4 KB. index.php that cua ban nang 44 KB.
REM  Xoa nham file nay la mat trang nguy trang, nen tha bo sot con hon.
set "LA_RAC_XAMPP=0"
if not exist "C:\xampp\htdocs\index.php" goto DONRAC_KHONG_CO_INDEX
set "CO_CHU_DASHBOARD=0"
findstr /i /c:"dashboard" "C:\xampp\htdocs\index.php" >nul 2>&1
if not errorlevel 1 set "CO_CHU_DASHBOARD=1"
call :KIEM_TRA_KICH_THUOC "C:\xampp\htdocs\index.php" 4096
REM  FILE_DU_LON=1 nghia la tu 4 KB tro len -> chac chan KHONG phai rac XAMPP.
if "%CO_CHU_DASHBOARD%"=="1" if "%FILE_DU_LON%"=="0" set "LA_RAC_XAMPP=1"

if "%LA_RAC_XAMPP%"=="1" goto DONRAC_XOA_INDEX
echo       - index.php hien tai KHONG phai cua XAMPP (co ve la trang that
echo         cua ban). Giu nguyen, khong dong den.
goto DONRAC_KHONG_CO_INDEX

:DONRAC_XOA_INDEX
del /f /q "C:\xampp\htdocs\index.php" >nul 2>&1
echo       - Da xoa index.php mac dinh cua XAMPP

:DONRAC_KHONG_CO_INDEX
if exist "C:\xampp\htdocs\dashboard" rd /s /q "C:\xampp\htdocs\dashboard" >nul 2>&1
if exist "C:\xampp\htdocs\dashboard" echo       - [CANH BAO] Khong xoa duoc thu muc dashboard
if not exist "C:\xampp\htdocs\dashboard" echo       - Thu muc dashboard: da don sach

if exist "C:\xampp\htdocs\applications.html" del /f /q "C:\xampp\htdocs\applications.html" >nul 2>&1
if exist "C:\xampp\htdocs\bitnami.css" del /f /q "C:\xampp\htdocs\bitnami.css" >nul 2>&1
if exist "C:\xampp\htdocs\favicon.ico" del /f /q "C:\xampp\htdocs\favicon.ico" >nul 2>&1
if exist "C:\xampp\htdocs\webalizer" rd /s /q "C:\xampp\htdocs\webalizer" >nul 2>&1
echo       - Da don applications.html, bitnami.css, favicon.ico, webalizer

REM  Thu muc img\ cua XAMPP: chi xoa khi no la ban cua XAMPP (co file
REM  xampp-logo.svg). Neu ban da chep img\ cua web game vao roi thi GIU.
if not exist "C:\xampp\htdocs\img" goto DONRAC_XONG
if not exist "C:\xampp\htdocs\img\xampp-logo.svg" goto DONRAC_GIU_IMG
rd /s /q "C:\xampp\htdocs\img" >nul 2>&1
echo       - Da xoa thu muc img\ mac dinh cua XAMPP
goto DONRAC_XONG
:DONRAC_GIU_IMG
echo       - Thu muc img\ hien co ve la cua web game. GIU NGUYEN.

:DONRAC_XONG
echo       OK - htdocs da san sang de chep web that vao.
set "KQ_DONRAC=OK - da don rac mac dinh cua XAMPP"
goto XONG_DONRAC

:DONRAC_KHONG_CO_THU_MUC
echo       Bo qua - chua co thu muc C:\xampp\htdocs (XAMPP chua cai duoc).
set "KQ_DONRAC=BO QUA - chua co htdocs"

:XONG_DONRAC
echo.


REM ===========================================================================
REM  BUOC 6/9 - CAI APACHE VA MYSQL THANH DICH VU WINDOWS
REM ===========================================================================
echo [6/9] Dang cai Apache va MySQL thanh dich vu Windows...

REM  VI SAO BUOC NAY QUAN TRONG NHAT TRONG CA FILE:
REM  Bo cai XAMPP chay che do "unattended" KHONG tu dang ky dich vu Windows.
REM  Neu khong lam buoc nay thi moi lan VPS khoi dong lai, ban phai tu tay
REM  mo XAMPP Control Panel bam Start - khong bam thi WEB CHET va CSDL CHET.
REM
REM  Va vi may chu game can CSDL de khoi dong, no se: bat len -> khong noi
REM  duoc CSDL -> thoat loi -> run.bat chay lai sau 10 giay -> lap vo tan.
REM  Ma web cung chet nen khong co gi bao cho ban biet.
REM
REM  Apache va MySQL KHAC may chu game o mot diem: chung khong can man hinh,
REM  nen chay duoc dang dich vu (session 0). Chi rieng may chu game moi
REM  bat buoc phai dung cach tu dong dang nhap - xem file 4-TU-DONG-BAT.bat.

set "LOI_DICHVU=0"

if exist "C:\xampp\apache\bin\httpd.exe" goto DICHVU_APACHE
echo       Bo qua - chua cai duoc XAMPP o buoc 4.
set "KQ_DICHVU=BO QUA - chua co XAMPP"
set "CO_LOI=1"
goto XONG_DICHVU

:DICHVU_APACHE
sc query "Apache2.4" >nul 2>&1
if not errorlevel 1 goto DICHVU_APACHE_DA_CO
echo       - Dang dang ky dich vu Apache...
call "C:\xampp\apache\apache_installservice.bat" >nul 2>&1
sc query "Apache2.4" >nul 2>&1
if errorlevel 1 goto DICHVU_APACHE_LOI
echo         Da dang ky dich vu Apache2.4
goto DICHVU_APACHE_TUCHAY
:DICHVU_APACHE_DA_CO
echo       - Dich vu Apache2.4 da co san, bo qua dang ky.
:DICHVU_APACHE_TUCHAY
sc config "Apache2.4" start= auto >nul 2>&1
sc start "Apache2.4" >nul 2>&1
goto DICHVU_MYSQL
:DICHVU_APACHE_LOI
echo         [LOI] Khong dang ky duoc dich vu Apache.
echo               Lam tay: mo cua so nay bang quyen Administrator va chay:
echo               C:\xampp\apache\apache_installservice.bat
set "LOI_DICHVU=1"

:DICHVU_MYSQL
sc query "mysql" >nul 2>&1
if not errorlevel 1 goto DICHVU_MYSQL_DA_CO
echo       - Dang dang ky dich vu MySQL...
call "C:\xampp\mysql\mysql_installservice.bat" >nul 2>&1
sc query "mysql" >nul 2>&1
if errorlevel 1 goto DICHVU_MYSQL_LOI
echo         Da dang ky dich vu mysql
goto DICHVU_MYSQL_TUCHAY
:DICHVU_MYSQL_DA_CO
echo       - Dich vu mysql da co san, bo qua dang ky.
:DICHVU_MYSQL_TUCHAY
sc config "mysql" start= auto >nul 2>&1
sc start "mysql" >nul 2>&1
goto DICHVU_KIEM_TRA
:DICHVU_MYSQL_LOI
echo         [LOI] Khong dang ky duoc dich vu MySQL.
echo               Lam tay: C:\xampp\mysql\mysql_installservice.bat
echo               Neu van loi, xem file nhat ky loi cua MariaDB tai:
echo               C:\xampp\mysql\data\%COMPUTERNAME%.err
set "LOI_DICHVU=1"

:DICHVU_KIEM_TRA
REM  Cho vai giay cho hai dich vu kip len roi moi kiem tra.
ping -n 6 127.0.0.1 >nul

set "APACHE_DANG_CHAY=0"
set "MYSQL_DANG_CHAY=0"
sc query "Apache2.4" 2>nul | findstr /i "RUNNING" >nul 2>&1
if not errorlevel 1 set "APACHE_DANG_CHAY=1"
sc query "mysql" 2>nul | findstr /i "RUNNING" >nul 2>&1
if not errorlevel 1 set "MYSQL_DANG_CHAY=1"

if "%APACHE_DANG_CHAY%"=="1" echo       - Apache: DANG CHAY
if "%APACHE_DANG_CHAY%"=="0" echo       - Apache: CHUA CHAY (xem cach xu ly trong HUONG-DAN.md)
if "%MYSQL_DANG_CHAY%"=="1" echo       - MySQL : DANG CHAY
if "%MYSQL_DANG_CHAY%"=="0" echo       - MySQL : CHUA CHAY (xem cach xu ly trong HUONG-DAN.md)

if "%APACHE_DANG_CHAY%"=="0" set "LOI_DICHVU=1"
if "%MYSQL_DANG_CHAY%"=="0" set "LOI_DICHVU=1"

if "%LOI_DICHVU%"=="1" goto DICHVU_CO_LOI
echo       OK - ca hai dich vu da chay va se TU LEN sau moi lan khoi dong lai.
set "KQ_DICHVU=OK - Apache + MySQL da thanh dich vu, tu chay"
goto XONG_DICHVU

:DICHVU_CO_LOI
echo.
echo       [CANH BAO] Dich vu chua tron. Xem chi tiet o tren.
echo                  Nguyen nhan hay gap nhat: cong 80 hoac 3306 dang bi
echo                  phan mem khac chiem. Kiem tra bang lenh:
echo                    netstat -ano ^| findstr ":80 "
echo                    netstat -ano ^| findstr ":3306"
echo.
set "KQ_DICHVU=CANH BAO - co dich vu chua chay, xem o tren"
set "CO_LOI=1"

:XONG_DICHVU
echo.


REM ===========================================================================
REM  BUOC 7/9 - DAT BIEN MOI TRUONG JAVA_HOME VA THEM VAO PATH
REM ===========================================================================
echo [7/9] Dang dat bien moi truong JAVA_HOME va PATH...

if defined JAVA_EXE goto DAT_BIEN_MT
echo       Bo qua - vi chua cai duoc Java 21 o buoc 3.
set "KQ_BIENMT=BO QUA - chua co Java"
set "CO_LOI=1"
goto XONG_BIENMT

:DAT_BIEN_MT
REM --- Tu duong dan java.exe suy nguoc ra thu muc goc cua JDK ---
REM     ...\jdk-21\bin\java.exe  ->  ...\jdk-21
set "TM_BIN="
for %%I in ("%JAVA_EXE%") do set "TM_BIN=%%~dpI"
set "TM_BIN=%TM_BIN:~0,-1%"
set "JAVA_GOC="
for %%I in ("%TM_BIN%") do set "JAVA_GOC=%%~dpI"
set "JAVA_GOC=%JAVA_GOC:~0,-1%"

echo       JAVA_HOME se duoc dat thanh:
echo         %JAVA_GOC%

setx /M JAVA_HOME "%JAVA_GOC%" >nul 2>&1
if errorlevel 1 goto BIENMT_LOI

REM  Them %JAVA_HOME%\bin vao PATH cua HE THONG.
REM  Dung PowerShell doc/ghi thang vao registry thay vi "setx /M PATH"
REM  vi setx cat cut PATH o ky tu thu 1024 - lam hong ca he thong.
set "KQ_PATH="
for /f "usebackq delims=" %%R in (`powershell -NoProfile -ExecutionPolicy Bypass -Command "$k = [Microsoft.Win32.Registry]::LocalMachine.OpenSubKey('SYSTEM\CurrentControlSet\Control\Session Manager\Environment', $true); $p = $k.GetValue('Path', '', [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames); if ($p -like '*JAVA_HOME*bin*') { 'DACO' } else { $k.SetValue('Path', $p.TrimEnd(';') + ';%%JAVA_HOME%%\bin', [Microsoft.Win32.RegistryValueKind]::ExpandString); 'DATHEM' }"`) do set "KQ_PATH=%%R"

if "%KQ_PATH%"=="DACO" goto BIENMT_XONG_DACO
if "%KQ_PATH%"=="DATHEM" goto BIENMT_XONG_DATHEM
goto BIENMT_LOI

:BIENMT_XONG_DACO
echo       OK - JAVA_HOME da dat. PATH da co san %%JAVA_HOME%%\bin.
set "KQ_BIENMT=OK - JAVA_HOME da dat, PATH da co san"
goto XONG_BIENMT

:BIENMT_XONG_DATHEM
echo       OK - JAVA_HOME da dat va da them %%JAVA_HOME%%\bin vao PATH.
set "KQ_BIENMT=OK - da dat JAVA_HOME va them vao PATH"
goto XONG_BIENMT

:BIENMT_LOI
echo.
echo       [THAT BAI] Khong ghi duoc bien moi truong he thong.
echo.
echo       CACH LAM TAY:
echo         Bam nut Start, go: bien moi truong  (hoac "environment variables")
echo         Chon "Edit the system environment variables"
echo           - Bam Environment Variables
echo           - O khung DUOI (System variables) bam New:
echo               Variable name : JAVA_HOME
echo               Variable value: %JAVA_GOC%
echo           - Van khung duoi, chon dong Path, bam Edit, bam New, go:
echo               %%JAVA_HOME%%\bin
echo.
set "KQ_BIENMT=THAT BAI - phai dat tay, xem huong dan o tren"
set "CO_LOI=1"

:XONG_BIENMT
echo.


REM ===========================================================================
REM  BUOC 8/9 - TUONG LUA: MO 80 VA 14445, CHAN 3306
REM ===========================================================================
echo [8/9] Dang cau hinh tuong lua Windows...

set "LOI_TUONGLUA=0"

REM ***************************************************************************
REM  CHU Y - DAY LA CHO NGUY HIEM NHAT TRONG CA FILE
REM  ------------------------------------------------------------------------
REM  Ban dang ngoi tren VPS qua Remote Desktop (cong 3389). Neu bat tuong lua
REM  len ma cong 3389 chua duoc cho phep, ban se BI DAY RA KHOI MAY NGAY LAP
REM  TUC va KHONG VAO LAI DUOC - vi muon sua thi phai vao duoc may.
REM
REM  Nhieu ban Windows Server ban giao may voi tuong lua DANG TAT, va bat RDP
REM  bang cach sua thang registry, de nguyen cac luat "Remote Desktop" o trang
REM  thai disabled. Neu VPS cua ban nhu vay thi bat tuong lua = mat VPS.
REM
REM  Nen: PHAI mo duong cho 3389 TRUOC, roi moi bat tuong lua.
REM  Lam ca hai cach cho chac: bat lai nhom luat "Remote Desktop" co san,
REM  VA them mot luat rieng cua minh.
REM ***************************************************************************

echo       - Dang mo duong cho Remote Desktop TRUOC khi bat tuong lua...
netsh advfirewall firewall set rule group="remote desktop" new enable=Yes >nul 2>&1
netsh advfirewall firewall set rule group="Remote Desktop" new enable=Yes >nul 2>&1
call :MO_CONG "Kame - Remote Desktop 3389" 3389 "Remote Desktop (giu duong vao VPS)"

REM  Kiem tra chac chan da co luat cho 3389 truoc khi dam bat tuong lua.
netsh advfirewall firewall show rule name="Kame - Remote Desktop 3389" >nul 2>&1
if not errorlevel 1 goto TUONGLUA_AN_TOAN_BAT

echo.
echo   ***********************************************************
echo   *  DUNG LAI: CHUA CHAC AN TOAN                            *
echo   ***********************************************************
echo.
echo   Khong tao duoc luat cho phep Remote Desktop (cong 3389).
echo   Script se KHONG bat tuong lua, vi bat luc nay co the day ban
echo   ra khoi VPS va khong vao lai duoc nua.
echo.
echo   Cac cong 80 va 14445 van se duoc mo binh thuong.
echo   Hay bat tuong lua bang tay SAU KHI da tu kiem tra 3389 duoc phep.
echo.
echo   ***********************************************************
echo.
set "LOI_TUONGLUA=1"
goto TUONGLUA_MO_CONG

:TUONGLUA_AN_TOAN_BAT
echo       - Da co luat cho 3389. Bay gio bat tuong lua an toan.
netsh advfirewall set allprofiles state on >nul 2>&1
if errorlevel 1 set "LOI_TUONGLUA=1"

:TUONGLUA_MO_CONG
REM --- Cong 80: web (Apache) ---
call :MO_CONG "Kame - Web HTTP 80" 80 "web (Apache)"

REM --- Cong 14445: may chu game ---
REM     So cong nay doc tu zz1\data\config\config.properties dong 7.
REM     Neu sau nay ban doi so cong o do thi phai sua ca luat nay.
call :MO_CONG "Kame - Game 14445" 14445 "may chu game"

REM --- Cong 3306: CSDL MariaDB - CHAN HAN ---
REM     TUYET DOI khong mo cong nay ra internet.
REM     Web va may chu game deu chay ngay tren chinh con VPS nay,
REM     ca hai noi vao CSDL qua duong noi bo (localhost) nen khong
REM     ai can 3306 mo ra ngoai ca.
REM     Trong Windows Firewall, luat CHAN luon thang luat CHO PHEP,
REM     nen day la lop khoa chac chan.
call :CHAN_CONG "Kame - Chan MySQL 3306" 3306 "CSDL MariaDB"

REM --- KHONG dung den cong 443 (HTTPS) o buoc nay ---
REM     Apache cua XAMPP dung chung chi tu ky. Mo 443 luc nay thi
REM     trinh duyet nguoi choi hien canh bao do "Ket noi khong rieng tu".
REM     Chi mo 443 SAU KHI da co chung chi that (Let's Encrypt hoac
REM     bat che do proxy cua Cloudflare).

if "%LOI_TUONGLUA%"=="0" goto TUONGLUA_OK
echo.
echo       [CANH BAO] Co luat tuong lua khong dat duoc.
echo                  Hay mo cua so nay bang quyen Administrator va chay lai.
echo.
set "KQ_TUONGLUA=CANH BAO - co luat khong dat duoc, xem o tren"
set "CO_LOI=1"
goto XONG_TUONGLUA

:TUONGLUA_OK
echo       OK - giu 3389, mo 80 va 14445, da chan 3306.
set "KQ_TUONGLUA=OK - giu 3389, mo 80 + 14445, chan 3306"

:XONG_TUONGLUA
echo.


REM ===========================================================================
REM  BUOC 9/9 - DAT MUI GIO VE GIO VIET NAM
REM ===========================================================================
echo [9/9] Dang kiem tra mui gio...

set "MUI_GIO_HIEN_TAI="
for /f "usebackq delims=" %%T in (`tzutil /g 2^>nul`) do set "MUI_GIO_HIEN_TAI=%%T"

if /i "%MUI_GIO_HIEN_TAI%"=="SE Asia Standard Time" goto MUIGIO_DA_DUNG

echo       Mui gio hien tai: %MUI_GIO_HIEN_TAI%
echo       Dang doi ve SE Asia Standard Time (gio Ha Noi, UTC+7)...
tzutil /s "SE Asia Standard Time" >nul 2>&1
if errorlevel 1 goto MUIGIO_LOI

echo       OK - da doi mui gio. Gio he thong bay gio la: %TIME%
set "KQ_MUIGIO=OK - da doi sang gio Viet Nam"
goto XONG_MUIGIO

:MUIGIO_DA_DUNG
echo       DA DUNG SAN - dang la SE Asia Standard Time (gio Viet Nam).
set "KQ_MUIGIO=DA DUNG SAN - gio Viet Nam"
goto XONG_MUIGIO

:MUIGIO_LOI
echo.
echo       [CANH BAO] Khong doi duoc mui gio.
echo                  Lam tay: bam chuot phai vao dong ho goc man hinh
echo                  chon "Adjust date/time", roi chon mui gio
echo                  (UTC+07:00) Bangkok, Hanoi, Jakarta
echo.
echo       Vi sao quan trong: moi su kien trong game, gio mo server,
echo       va nhat ky nap the deu tinh theo gio cua may nay. Sai mui gio
echo       thi su kien no sai gio, nguoi choi keu ngay.
echo.
set "KQ_MUIGIO=CANH BAO - khong doi duoc, hay dat tay"
set "CO_LOI=1"

:XONG_MUIGIO
echo.


REM ===========================================================================
REM  BANG TONG KET
REM ===========================================================================
echo.
echo ===============================================================
echo    BANG TONG KET - BUOC 1
echo ===============================================================
echo.
echo   1. Quyen Administrator ... %KQ_QUYEN%
echo   2. Dung luong o C: ...... %KQ_ODIA%
echo   3. Java JDK 21 .......... %KQ_JDK%
echo   4. XAMPP 8.2 ............ %KQ_XAMPP%
echo   5. Don rac htdocs ....... %KQ_DONRAC%
echo   6. Dich vu Apache+MySQL . %KQ_DICHVU%
echo   7. JAVA_HOME va PATH .... %KQ_BIENMT%
echo   8. Tuong lua ............ %KQ_TUONGLUA%
echo   9. Mui gio .............. %KQ_MUIGIO%
echo.
echo ===============================================================

REM --- Ghi ket qua ra file de xem lai sau ---
echo NGOC RONG KAME - KET QUA BUOC 1 - %DATE% %TIME%> "%FILE_LOG%"
echo.>> "%FILE_LOG%"
echo 1. Quyen Administrator : %KQ_QUYEN%>> "%FILE_LOG%"
echo 2. Dung luong o C:     : %KQ_ODIA%>> "%FILE_LOG%"
echo 3. Java JDK 21         : %KQ_JDK%>> "%FILE_LOG%"
echo 4. XAMPP 8.2           : %KQ_XAMPP%>> "%FILE_LOG%"
echo 5. Don rac htdocs      : %KQ_DONRAC%>> "%FILE_LOG%"
echo 6. Dich vu Apache+MySQL: %KQ_DICHVU%>> "%FILE_LOG%"
echo 7. JAVA_HOME va PATH   : %KQ_BIENMT%>> "%FILE_LOG%"
echo 8. Tuong lua           : %KQ_TUONGLUA%>> "%FILE_LOG%"
echo 9. Mui gio             : %KQ_MUIGIO%>> "%FILE_LOG%"
if defined JAVA_EXE echo.>> "%FILE_LOG%"
if defined JAVA_EXE echo Duong dan java.exe: %JAVA_EXE%>> "%FILE_LOG%"

echo.
echo   Ket qua nay da duoc luu vao file:
echo     %FILE_LOG%
echo.

if "%CO_LOI%"=="1" goto KET_THUC_CO_LOI


REM ===========================================================================
REM  KET THUC - MOI VIEC TRON TRU
REM ===========================================================================
echo ===============================================================
echo    XONG BUOC 1 - MOI THU DEU OK
echo ===============================================================
echo.
echo   VIEC CAN LAM NGAY BAY GIO:
echo.
echo   1. DONG cua so nay lai.
echo.
echo   2. Bien moi truong JAVA_HOME chi co hieu luc voi cua so
echo      MOI mo. Cac cua so cmd dang mo san se khong thay no.
echo      Neu muon chac chan: khoi dong lai VPS mot lan.
echo.
echo   3. BAY GIO moi chep du lieu game va web len VPS.
echo      Xem muc "Buoc 4" trong file HUONG-DAN.md.
echo.
echo   4. Sau khi chep xong, chay file 3-CAI-GAME-LEN-VPS.ps1
echo      (bam chuot phai vao no, chon "Run with PowerShell").
echo      No se hoi mat khau CSDL va tu dien vao ca 3 noi.
echo.
echo   NEN LAM NGAY LUC NAY: chup mot ban snapshot cua VPS tren trang
echo   quan ly cua noi ban VPS. Day la thoi diem dep nhat de chup -
echo   nen mong da xong, du lieu game chua chep. Hong buoc nao sau nay
echo   cung chi mat 5 phut quay lai, thay vi cai lai tu dau.
echo.
echo   NHUNG DIEU CAN NHO (chua lam o buoc nay):
echo.
echo   - CSDL MySQL van dang KHONG CO MAT KHAU. File 3 se dat.
echo     Nho: mat khau do phai sua DONG THOI o 3 file, thieu 1 la hong:
echo       zz1\data\config\config.properties
echo       C:\xampp\htdocs\sv1\config.php
echo       C:\xampp\htdocs\sv1\boss_status.php
echo     Va phai dat cho ca 4 tai khoan cua MySQL, khong phai 1.
echo     File 3 lo het viec nay - dung lam tay neu tranh duoc.
echo.
echo   - TUYET DOI khong bat FileZilla (FTP) va Mercury (mail) trong
echo     XAMPP Control Panel - ban khong dung den chung, bat len chi
echo     tang cua cho ke xau.
echo.
echo   - Cong 443 (HTTPS) co y KHONG mo o buoc nay, vi chung chi
echo     cua XAMPP la chung chi tu ky, trinh duyet se bao do.
echo.
echo   - May chu game VAN CHUA tu bat sau khi khoi dong lai. Viec do
echo     lam o file 4-TU-DONG-BAT.bat, va chi lam SAU KHI game da
echo     chay ngon lanh bang tay.
echo.
echo ===============================================================
echo.
pause
exit /b 0

:KET_THUC_CO_LOI
echo ===============================================================
echo    XONG BUOC 1 - NHUNG CO VIEC CHUA TRON
echo ===============================================================
echo.
echo   Hay doc lai bang tong ket o tren, tim dong nao ghi
echo   THAT BAI hoac CANH BAO.
echo.
echo   Voi moi dong do, script da in san huong dan lam tay ngay
echo   luc no xay ra - hay keo cua so nay len tren de doc lai.
echo.
echo   Sua xong thi CHAY LAI chinh file nay. Script biet bo qua
echo   nhung phan da lam duoc, chi lam lai phan con thieu.
echo.
echo   DUNG chep du lieu game len khi con dong THAT BAI.
echo   Rieng dong CANH BAO thi doc ky roi tu quyet dinh.
echo.
echo ===============================================================
echo.
pause
exit /b 1


REM ===========================================================================
REM  ==========================  CAC HAM PHU  ================================
REM ===========================================================================

REM ---------------------------------------------------------------------------
REM  :TIM_JAVA - tim java.exe phien ban 21 tren may
REM  Ket qua tra ve trong bien JAVA_EXE (rong = khong tim thay)
REM ---------------------------------------------------------------------------
:TIM_JAVA
set "JAVA_EXE="
REM  Uu tien 1: dung cho ma run.bat cua may chu game tim toi dau tien
call :THU_MOT_JAVA "C:\Program Files\Java\jdk-21\bin\java.exe"
if defined JAVA_EXE goto :eof
REM  Uu tien 2: cho cai mac dinh cua Eclipse Temurin / winget
for /d %%D in ("C:\Program Files\Eclipse Adoptium\jdk-21*") do call :THU_MOT_JAVA "%%~fD\bin\java.exe"
if defined JAVA_EXE goto :eof
REM  Uu tien 3: theo bien JAVA_HOME neu ai do da dat san
if defined JAVA_HOME call :THU_MOT_JAVA "%JAVA_HOME%\bin\java.exe"
goto :eof

REM ---------------------------------------------------------------------------
REM  :THU_MOT_JAVA <duong-dan-java.exe>
REM  Kiem tra mot duong dan cu the: co ton tai khong, va co PHAI ban 21 khong.
REM  Phai kiem tra so 21 that su, khong duoc tin vao ten thu muc.
REM  (Tren may cu dang co san JRE 8 tren PATH - chay zz1.jar bang no se
REM   loi UnsupportedClassVersionError vi zz1.jar la bytecode major 65.)
REM ---------------------------------------------------------------------------
:THU_MOT_JAVA
if defined JAVA_EXE goto :eof
if not exist %1 goto :eof
REM  Phai ghi ket qua ra file tam roi moi doc.
REM  Ly do: duong dan chua khoang trang ("C:\Program Files\...") ma dat
REM  thang vao trong 'for /f' thi cmd cat mat dau nhay kep, bao loi
REM  "'C:\Program' is not recognized". Da bi dinh loi nay khi thu that.
set "FILE_TAM_JAVA=%TEMP%\kame-kiem-tra-java.txt"
del /f /q "%FILE_TAM_JAVA%" >nul 2>&1
%1 -version > "%FILE_TAM_JAVA%" 2>&1
if not exist "%FILE_TAM_JAVA%" goto :eof
set "PHIEN_BAN="
for /f "tokens=3" %%V in ('findstr /i "version" "%FILE_TAM_JAVA%"') do if not defined PHIEN_BAN set "PHIEN_BAN=%%V"
del /f /q "%FILE_TAM_JAVA%" >nul 2>&1
if not defined PHIEN_BAN goto :eof
REM  Dong dau tien co dang:  openjdk version "21.0.11" 2026-04-15 LTS
REM  Nen token thu 3 la:     "21.0.11"   (ke ca dau nhay kep)
REM  Lay 2 ky tu tu vi tri 1 (bo qua dau nhay dau) duoc:  21
REM  Voi JRE 8 thi chuoi la "1.8.0_xxx" -^> ra "1." -^> bi loai dung y muon.
set "PHIEN_BAN=%PHIEN_BAN:~1,2%"
if not "%PHIEN_BAN%"=="21" goto :eof
set "JAVA_EXE=%~1"
goto :eof

REM ---------------------------------------------------------------------------
REM  :TAI_FILE <dia-chi> <luu-vao>
REM  Tai mot file tu internet. Uu tien curl.exe (nhanh, co san tu
REM  Windows Server 2019), khong co thi dung PowerShell.
REM ---------------------------------------------------------------------------
:TAI_FILE
where curl.exe >nul 2>&1
if errorlevel 1 goto TAI_BANG_POWERSHELL
curl.exe -L --fail --retry 2 --retry-delay 5 --connect-timeout 30 -o %2 %1
goto :eof

:TAI_BANG_POWERSHELL
REM  Bat TLS 1.2: mot so ban Windows Server cu mac dinh dung TLS 1.0,
REM  se bi cac trang tai tu choi ket noi.
powershell -NoProfile -ExecutionPolicy Bypass -Command "$ProgressPreference = 'SilentlyContinue'; [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; try { Invoke-WebRequest -Uri %1 -OutFile %2 -UseBasicParsing } catch { Write-Host ('   Loi tai file: ' + $_.Exception.Message) }"
goto :eof

REM ---------------------------------------------------------------------------
REM  :KIEM_TRA_KICH_THUOC <file> <so-byte-toi-thieu>
REM  Dat FILE_DU_LON=1 neu file ton tai va du lon.
REM  Muc dich: bat truong hop tai ve nham trang bao loi HTML vai KB
REM  thay vi bo cai that. Neu khong kiem tra, msiexec se bao loi
REM  kho hieu va rat kho doan nguyen nhan.
REM ---------------------------------------------------------------------------
:KIEM_TRA_KICH_THUOC
set "FILE_DU_LON=0"
if not exist %1 goto :eof
set "KICH_THUOC=0"
for %%S in (%1) do set "KICH_THUOC=%%~zS"
if %KICH_THUOC% GEQ %2 set "FILE_DU_LON=1"
goto :eof

REM ---------------------------------------------------------------------------
REM  :MO_CONG <ten-luat> <so-cong> <mo-ta>
REM  Mo mot cong TCP di vao. Neu luat da co thi bo qua (chay lai duoc).
REM  Chi can TCP - may chu game dung java.net.ServerSocket, luon la TCP,
REM  nen khong phai mo luat UDP.
REM ---------------------------------------------------------------------------
:MO_CONG
netsh advfirewall firewall show rule name=%1 >nul 2>&1
if not errorlevel 1 goto MO_CONG_DA_CO
netsh advfirewall firewall add rule name=%1 dir=in action=allow protocol=TCP localport=%2 >nul 2>&1
if errorlevel 1 goto MO_CONG_LOI
echo       - Da MO cong %2 cho %~3
goto :eof
:MO_CONG_DA_CO
echo       - Cong %2 cho %~3: luat da co san, bo qua.
goto :eof
:MO_CONG_LOI
echo       - [LOI] Khong mo duoc cong %2 cho %~3
set "LOI_TUONGLUA=1"
goto :eof

REM ---------------------------------------------------------------------------
REM  :CHAN_CONG <ten-luat> <so-cong> <mo-ta>
REM ---------------------------------------------------------------------------
:CHAN_CONG
netsh advfirewall firewall show rule name=%1 >nul 2>&1
if not errorlevel 1 goto CHAN_CONG_DA_CO
netsh advfirewall firewall add rule name=%1 dir=in action=block protocol=TCP localport=%2 >nul 2>&1
if errorlevel 1 goto CHAN_CONG_LOI
echo       - Da CHAN cong %2 (%~3) khong cho vao tu internet
goto :eof
:CHAN_CONG_DA_CO
echo       - Cong %2 (%~3): luat chan da co san, bo qua.
goto :eof
:CHAN_CONG_LOI
echo       - [LOI] Khong dat duoc luat chan cong %2 (%~3)
set "LOI_TUONGLUA=1"
goto :eof
