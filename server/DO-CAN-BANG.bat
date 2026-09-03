@echo off
chcp 65001 >nul
REM ===========================================================================
REM  DO CAN BANG — bam dup file nay de biet game dang de hay kho
REM ===========================================================================
REM
REM  CACH DUNG (khong can biet lap trinh):
REM
REM    1. Mo mot file so trong thu muc  server\kho-noi-dung\  (vi du lang.yaml,
REM       quai.yaml, cua-hang.yaml, bang-roi.yaml, boss.yaml) va sua MOT so.
REM    2. Bam dup chinh file nay.
REM    3. Doc bang vua hien ra — no o  server\bang-can-bang.txt
REM
REM  Bang tra loi bon cau hoi:
REM    D1  Nguoi moi giet mot con Soi Rung mat bao lau?
REM    D2  Nam nguoi ha boss lon pha 1 mat bao lau?
REM    D3  Lang nao thang?
REM    D4  Danh quai bao lau moi mua noi mot mon trong cua hang?
REM
REM  KHONG CAN CO SO DU LIEU. Bon phep do chay hoan toan trong bo nho, dong ho
REM  duoc ghim (khong cho thoi gian that troi), nen ket qua khong phu thuoc may
REM  nhanh hay cham — chay hai lan ra dung mot con so.
REM
REM  Neu bao "go is not recognized": may chua cai Go (https://go.dev/dl/).
REM ===========================================================================

cd /d "%~dp0"

echo.
echo === DANG DO CAN BANG (mat vai giay) ===
echo.

go test ./internal/game/ -run "TestCanBang" -count=1 -v
if errorlevel 1 (
    echo.
    echo *** CO PHEP DO BI DO ***
    echo Doc dong bat dau bang --- FAIL o tren de biet cai gi hong.
    echo Neu la TestCanBang_GoldenSoHienTai thi day la DUNG Y: ban vua doi mot so
    echo trong kho-noi-dung, va test dang bao cho ban biet no lam con so doi bao nhieu.
    echo.
    pause
    exit /b 1
)

echo.
echo === BANG CAN BANG ===
echo.
type "%~dp0bang-can-bang.txt"

echo.
echo Bang tren da duoc luu tai:
echo    %~dp0bang-can-bang.txt
echo.
pause
