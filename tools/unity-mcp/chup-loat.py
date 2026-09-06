# -*- coding: utf-8 -*-
"""chup-loat.py — chụp liên tiếp N ảnh Game view qua Unity MCP (dùng để soi vòng chạy cf 2→6).

Vì sao: `cs/chup.cs` chỉ chụp một ảnh vào một đường dẫn cố định. Muốn xem hoạt ảnh
(chạy, đánh) phải chụp nhiều khung cách nhau ~100 ms rồi ghép lại xem từng khung.

Dùng:
    python tools\\unity-mcp\\chup-loat.py --out <thư mục> --ten chay --n 8 --cach 0.12
Ảnh ra: <thư mục>\\<ten>-00.png … (ScreenCapture ghi ở cuối khung, đọc file ở lượt sau).
"""
import argparse
import importlib.util
import os
import sys
import time
from pathlib import Path

DAY = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location("unity_mcp", DAY / "unity_mcp.py")
um = importlib.util.module_from_spec(spec)
spec.loader.exec_module(um)


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--out", required=True)
    ap.add_argument("--ten", default="shot")
    ap.add_argument("--n", type=int, default=6)
    ap.add_argument("--cach", type=float, default=0.12, help="giây giữa hai lần chụp")
    ap.add_argument("--truoc", default=None, help="file .cs chạy TRƯỚC loạt (vd bấm phím chạy)")
    ap.add_argument("--sau", default=None, help="file .cs chạy SAU loạt (vd thả phím)")
    a = ap.parse_args()
    out = Path(a.out)
    out.mkdir(parents=True, exist_ok=True)
    sid = um.phien()

    def code(cs: str) -> str:
        r, _ = um.goi("tools/call", {"name": "execute_code", "arguments": {"action": "execute", "code": cs}}, sid, 3)
        if "error" in r:
            return "LOI " + str(r["error"])
        return " ".join(c.get("text", "") for c in r.get("result", {}).get("content", []))

    if a.truoc:
        print("truoc:", code(open(a.truoc, encoding="utf-8").read())[:200])
    for i in range(a.n):
        p = (out / f"{a.ten}-{i:02d}.png").as_posix()
        # Ép Editor vẽ lại trước: Unity không có focus thì Game view không tự vẽ, ScreenCapture trả khung cũ.
        cs = f'UnityEditorInternal.InternalEditorUtility.RepaintAllViews(); UnityEngine.ScreenCapture.CaptureScreenshot("{p}"); return "{i}";'
        code(cs)
        time.sleep(a.cach)
    if a.sau:
        print("sau:", code(open(a.sau, encoding="utf-8").read())[:200])
    time.sleep(0.5)
    co = sorted(f.name for f in out.glob(f"{a.ten}-*.png"))
    print("da chup", len(co), "anh vao", out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
