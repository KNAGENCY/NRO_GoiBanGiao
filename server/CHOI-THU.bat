@echo off
chcp 65001 >nul
cd /d "%~dp0"
REM ===========================================================================
REM  CHOI THU — bam dup file nay de bat may chu game Naruto tren may nha
REM ===========================================================================
REM
REM  No lam 3 viec, theo thu tu:
REM    1. Bat PostgreSQL neu chua chay (D:\PostgreSQL, CSDL "nro").
REM    2. Build lai nroserver.exe tu ma nguon moi nhat (mat ~10-20 giay).
REM    3. Chay may chu, lang nghe cong 14445. CUA SO NAY PHAI DE MO trong luc choi.
REM       Tat may chu: bam Ctrl+C trong cua so nay (no se luu nhan vat roi tat).
REM
REM  Sau do: mo Unity (D:\Unity\2022.3.62f2) -> project D:\NRO_Project -> bam Play.
REM  Client da tro san 127.0.0.1:14445. Tai khoan thu: thutest / 1  (hoac naruto / 1).
REM
REM  Neu bao "khong noi duoc CSDL": xem cong PostgreSQL dang nghe (5432 hay 5433)
REM  bang PowerShell:  Get-NetTCPConnection -State Listen | ? { $_.LocalPort -in 5432,5433 }
REM  roi sua so cong trong dong "set NRO_DSN=" ben duoi.
REM ===========================================================================

set PG=D:\PostgreSQL
set GO=D:\Go\bin\go.exe
set NRO_DSN=postgres://postgres:nrodev123@127.0.0.1:5432/nro

echo [1/3] Kiem PostgreSQL...
"%PG%\bin\pg_isready.exe" -h 127.0.0.1 -p 5432 >nul 2>&1
if errorlevel 1 (
  echo       chua chay -> dang bat...
  "%PG%\bin\pg_ctl.exe" -D "%PG%\data" -l "%PG%\data\pg-choi-thu.log" start -w
  if errorlevel 1 (
    echo       KHONG BAT DUOC PostgreSQL. Xem %PG%\data\pg-choi-thu.log
    pause
    exit /b 1
  )
) else (
  echo       dang chay.
)

echo [2/3] Build may chu...
"%GO%" build -o nroserver.exe ./cmd/nroserver
if errorlevel 1 (
  echo       BUILD LOI - ma nguon dang hong, khong chay ban cu de khoi nham.
  pause
  exit /b 1
)

echo [3/3] Chay may chu (cong 14445). Ctrl+C de tat.
echo ---------------------------------------------------------------------------
nroserver.exe
echo ---------------------------------------------------------------------------
echo May chu da tat.
pause
