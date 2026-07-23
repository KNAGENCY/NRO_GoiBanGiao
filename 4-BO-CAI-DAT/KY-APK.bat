@echo off
REM ===========================================================================
REM  BAM DUP VAO FILE NAY DE KY BAN APK
REM
REM  File nay chi lam mot viec: goi KY-VA-DANG-APK.ps1 nam cung thu muc.
REM  Co no de khoi phai mo PowerShell va go duong dan bang tay.
REM ===========================================================================

REM %~dp0 = thu muc chua chinh file .bat nay, ke ca khi bam dup tu cho khac.
set "THUMUC=%~dp0"
set "SCRIPT=%THUMUC%KY-VA-DANG-APK.ps1"

if not exist "%SCRIPT%" (
    echo.
    echo   KHONG TIM THAY: %SCRIPT%
    echo   File .bat nay phai nam cung thu muc voi KY-VA-DANG-APK.ps1
    echo.
    pause
    exit /b 1
)

REM -ExecutionPolicy Bypass: Windows mac dinh chan chay file .ps1. Tham so nay
REM chi noi long CHO RIENG lan chay nay, khong doi cai dat cua may.
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%SCRIPT%"

if errorlevel 1 (
    echo.
    echo   Chua ky xong. Doc dong mau do o tren de biet vi sao.
    echo.
    pause
)
