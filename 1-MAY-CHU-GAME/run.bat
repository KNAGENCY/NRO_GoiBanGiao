@ECHO OFF
TITLE Ngoc Rong Kame - May Chu
cd /d "%~dp0"

REM ===========================================================================
REM  VONG CANH GAC MAY CHU NGOC RONG KAME
REM
REM  Cua so nay dung canh may chu. Dung dong no khi dang cho nguoi choi vao.
REM
REM  May chu thoat voi cac ma khac nhau, va file nay xu ly tung truong hop:
REM     ma 2  = BAO TRI      -> chay lai ngay sau 3 giay
REM     ma 0  = DUNG HAN     -> thoat, khong chay lai
REM             (nut "Save data" va "Thoat" trong giao dien dung ma nay)
REM     khac  = SERVER LOI   -> tu chay lai sau 10 giay, co ghi log
REM
REM  Moi lan chay lai deu duoc ghi vao logs\restart.log de xem lai sau.
REM ===========================================================================

set "GAME_ROOT=%CD%"
set "WEB_ROOT=C:\xampp\htdocs"
set "BOSS_STATUS_FILE=%WEB_ROOT%\sv1\boss_status.json"

REM zz1.jar bien dich bang Java 21. Lenh "java" tren PATH co the la ban cu hon
REM nen phai goi thang JDK 21, neu khong se loi UnsupportedClassVersionError.
set "JAVA_BIN=C:\Program Files\Java\jdk-21\bin\java.exe"
if not exist "%JAVA_BIN%" set "JAVA_BIN=%JAVA_HOME%\bin\java.exe"
if not exist "%JAVA_BIN%" (
    echo.
    echo [LOI] Khong tim thay Java 21.
    echo       Sua duong dan JAVA_BIN o dau file nay cho dung may cua ban.
    echo.
    pause
    exit /b 1
)

if not exist "zz1.jar" (
    echo.
    echo [LOI] Khong thay file zz1.jar trong thu muc nay.
    echo.
    pause
    exit /b 1
)

if not exist "logs" mkdir "logs"

set "SO_LAN=0"

:VONG_LAP
set /a SO_LAN+=1
echo.
echo ===============================================================
echo   NGOC RONG KAME - dang khoi dong may chu  (lan %SO_LAN%)
echo   %DATE% %TIME%
echo ===============================================================
echo.

"%JAVA_BIN%" -server -Dfile.encoding=UTF-8 -Dgame.root="%GAME_ROOT%" -Dweb.root="%WEB_ROOT%" -Dboss.status.file="%BOSS_STATUS_FILE%" -jar zz1.jar

set "MA=%ERRORLEVEL%"

if "%MA%"=="2" (
    echo.
    echo ---------------------------------------------------------------
    echo   [BAO TRI] May chu yeu cau khoi dong lai.
    echo   Dang chay lai sau 3 giay...
    echo ---------------------------------------------------------------
    echo %DATE% %TIME%  ^| BAO TRI - khoi dong lai>> "logs\restart.log"
    ping -n 4 127.0.0.1 >nul
    goto VONG_LAP
)

if "%MA%"=="0" (
    echo.
    echo ---------------------------------------------------------------
    echo   [DA DUNG] May chu da dung theo yeu cau. Du lieu da luu.
    echo ---------------------------------------------------------------
    echo %DATE% %TIME%  ^| DUNG HAN theo yeu cau>> "logs\restart.log"
    goto KET_THUC
)

echo.
echo ---------------------------------------------------------------
echo   [LOI] May chu thoat bat thuong. Ma loi: %MA%
echo   Tu dong chay lai sau 10 giay de nguoi choi vao duoc.
echo   Xem chi tiet loi o: logs\restart.log
echo ---------------------------------------------------------------
echo %DATE% %TIME%  ^| LOI ma %MA% - tu chay lai>> "logs\restart.log"
ping -n 11 127.0.0.1 >nul
goto VONG_LAP

:KET_THUC
echo.
echo Nhan phim bat ky de dong cua so nay.
pause >nul
