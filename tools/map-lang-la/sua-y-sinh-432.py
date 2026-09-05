# -*- coding: utf-8 -*-
"""sua-y-sinh-432.py — đưa điểm sinh người chơi / sói / NPC của MAP 1 từ y=384 (hàng 16, giữa trời)
về mặt cỏ y=432 (hàng 18 của map-1-collision.bin; hàng 18 để RỖNG cho nhân vật đứng, hàng 19 mới đặc).

Vì sao (05/09/2026, chơi thử thật qua Unity MCP): client KHÔNG cho quái rơi (Mob.cs: y luôn kéo về
yFirst) nên hai con sói đứng lơ lửng đúng 2 ô trên mặt đất; NPC client dựng ở y+3 cũng lơ lửng; người
chơi vào game rơi 2 ô rồi mới đứng. SOI-LAN-3 #26/#35 đã nói từ 22/08.

Sửa BỐN nơi cùng lúc (lệch nhau là test tích hợp đỏ: khuvuc_tichhop_test so bd.SpawnY với vị trí thật):
  1. server/kho-noi-dung/quai.yaml       — sinh: map 1, hai con sói y 384 -> 432
  2. server/kho-noi-dung/lang.yaml       — làng Lá vi_tri_y_khoi_dau 384 -> 432
  3. server/internal/game/assets/catalog.json — maps[0].spawn.y 432, npcs[0].y 429 (client +3 = 432),
     mobs[0].spawns[*].y 432   (file này bị ghim SHA trong asset-lock.json)
  4. asset-lock.json — vá đúng mục "server/catalog.json" ở outputs[] và provenance.artifacts[] bằng regex
     (KHÔNG json.dump: file một dòng 444KB, test bám chuỗi nguyên văn — xem sinh-va-cham.py).
Rồi liệt kê các test đang ghim 384 để sửa tay (script không tự sửa test).

Chạy trong worktree/cây có main đã gộp nhánh sprite (asset-lock đụng nhau nếu chạy sớm).
"""
import hashlib, io, json, os, re, sys

GOC = sys.argv[1] if len(sys.argv) > 1 else r"D:\NRO_GoiBanGiao"
QUAI = os.path.join(GOC, "server", "kho-noi-dung", "quai.yaml")
LANG = os.path.join(GOC, "server", "kho-noi-dung", "lang.yaml")
CAT = os.path.join(GOC, "server", "internal", "game", "assets", "catalog.json")
LOCK = os.path.join(GOC, "server", "internal", "game", "assets", "asset-lock.json")

def doc(p): return io.open(p, encoding="utf-8").read()
def ghi(p, s): io.open(p, "w", encoding="utf-8", newline="\n").write(s)

# 1. quai.yaml — chỉ khối map 1
q = doc(QUAI)
khoi = re.search(r"(  - map: 1 # Làng Lá\n    dat:\n)((?:      - \{[^\n]*\}\n)+)", q)
assert khoi, "khong thay khoi sinh map 1 trong quai.yaml"
moi = khoi.group(2).replace("y: 384", "y: 432")
chu = ("  # Map 1: y = 432 = mặt cỏ HÀNG 18 của map-1-collision.bin (hàng 18 để RỖNG cho nhân vật\n"
       "  # đứng, từ hàng 19 mới đặc — tools/map-lang-la/sinh-va-cham.py). Trước 05/09/2026 để 384\n"
       "  # (hàng 16, giữa trời): client KHÔNG cho quái rơi (Mob.cs kéo y về yFirst) nên hai con sói\n"
       "  # lơ lửng đúng 2 ô. Map 2–6 chưa có lưới va chạm riêng nên giữ 384 tạm.\n")
q2 = q.replace(khoi.group(0), chu + khoi.group(1) + moi) if "mặt cỏ HÀNG 18" not in q else q
if q2 != q: ghi(QUAI, q2); print("quai.yaml: map 1 y -> 432")
else: print("quai.yaml: da 432, bo qua")

# 2. lang.yaml — chỉ làng Lá
l = doc(LANG)
m = re.search(r"(    map_khoi_dau: 1\n    vi_tri_x_khoi_dau: 100\n)    vi_tri_y_khoi_dau: 384\n", l)
if m:
    l2 = l.replace(m.group(0), m.group(1) +
        "    # 432 = mặt cỏ hàng 18 của map 1 (map-1-collision.bin). Trước 05/09/2026 để 384 (hàng 16,\n"
        "    # giữa trời): vào game rơi 2 ô. Làng 2–5 chưa có map riêng nên vẫn 384 tạm.\n"
        "    vi_tri_y_khoi_dau: 432\n")
    ghi(LANG, l2); print("lang.yaml: lang La y -> 432")
else: print("lang.yaml: khong thay 384 cua lang La (da sua?)")

# 3. catalog.json — sửa bằng regex GIỮ NGUYÊN định dạng một dòng
c = doc(CAT)
c2 = c
c2 = re.sub(r'("id":\s*1,\s*"name":\s*"Làng Lá",\s*"spawn":\s*\{\s*"x":\s*100,\s*"y":\s*)384', r"\g<1>432", c2, count=1)
# npcs[0]: khoa xep abc trong file (…"portrait":27078,"x":172,"y":384) — "menu" co ']' nen khong dung [^\]]*
c2 = re.sub(r'("portrait":\s*27078,\s*"x":\s*172,\s*"y":\s*)384', r"\g<1>429", c2, count=1)
c2 = re.sub(r'("spawns":\s*\[\s*\{\s*"id":\s*0,\s*"x":\s*160,\s*"y":\s*)384', r"\g<1>432", c2, count=1)
c2 = re.sub(r'("id":\s*1,\s*"x":\s*40,\s*"y":\s*)384', r"\g<1>432", c2, count=1)
if c2 == c: print("catalog.json: khong doi gi (da sua?)")
else:
    json.loads(c2)  # vẫn là JSON hợp lệ
    ghi(CAT, c2); print("catalog.json: maps[0].spawn.y 432, npcs[0].y 429, mobs spawns y 432")
    data = io.open(CAT, "rb").read()
    sha = hashlib.sha256(data).hexdigest()
    # 4. asset-lock.json — vá từng chỗ theo giá trị nó đang có (xem sinh-va-cham.py)
    van = doc(LOCK); kt = json.loads(van); da = set(); n_tong = 0
    for ten, ds in (("outputs", kt["outputs"]), ("provenance.artifacts", kt.get("provenance", {}).get("artifacts", []))):
        x = next((i for i in ds if i.get("path") == "server/catalog.json"), None)
        if x is None: print(f"  ({ten}: khong ghi catalog.json, bo qua)"); continue
        if (x["bytes"], x["sha256"]) in da: continue
        mau = re.compile(r'\{\s*"bytes":\s*' + str(x["bytes"]) + r',\s*"path":\s*"server/catalog\.json",\s*"sha256":\s*"' + x["sha256"] + r'"\s*\}')
        n = len(mau.findall(van))
        if n == 0: raise SystemExit(f"khong tim thay muc {ten} trong van ban lock — DUNG LAI")
        van = mau.sub('{"bytes":%d,"path":"server/catalog.json","sha256":"%s"}' % (len(data), sha), van)
        da.add((x["bytes"], x["sha256"])); n_tong += n
        print(f"  lock {ten}: {x['bytes']}B/{x['sha256'][:12]} -> {len(data)}B/{sha[:12]} ({n} cho)")
    kt2 = json.loads(van)
    assert len(kt2["outputs"]) == len(kt["outputs"])
    io.open(LOCK, "w", encoding="utf-8", newline="").write(van)
    print("asset-lock.json: da va", n_tong, "cho")

print("\nCAC TEST DANG GHIM 384 (sua tay, co chu thich):")
os.system(f'cd /d "{GOC}\\server" && findstr /s /n /c:"384" internal\\game\\*_test.go internal\\noidung\\*_test.go | findstr /v "canbang_test dichuyen_test catalog_test.go:9 catalog_test.go:10 congtrutien store_test"')
