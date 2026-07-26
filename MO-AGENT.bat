@echo off
REM ===========================================================================
REM  MO MOT AGENT LAM VIEC TREN DU AN NARUTO
REM
REM  Nhay dup vao file nay de mo Claude Code DUNG trong thu muc du an.
REM  Quan trong: chi khi mo o day thi agent moi:
REM    - tu dong doc duoc boi canh (git log + NHAT-KY-AGENT.md) qua hook
REM    - co san lenh /viec
REM  Mo Claude Code o thu muc khac (vi du C:\Users\Admin) thi MAT ca hai thu do.
REM
REM  Muon nhieu agent lam song song: nhay dup file nay nhieu lan, moi cua so
REM  la mot agent. Nho cho moi agent lam mot NHANH GIT RIENG (xem NHAT-KY-AGENT.md).
REM ===========================================================================

cd /d "%~dp0"
echo.
echo   Dang mo agent trong: %CD%
echo   Go  /viec  de nhan viec. Xong viec go  /viec  lan nua de ghi so.
echo.

claude

REM Neu cua so tu dong tat ma khong thay gi, bo dau REM o dong duoi de doc loi:
REM pause
