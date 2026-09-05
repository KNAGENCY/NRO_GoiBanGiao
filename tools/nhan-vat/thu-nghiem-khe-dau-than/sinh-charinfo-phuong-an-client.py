# -*- coding: utf-8 -*-
"""sinh-charinfo.py — sinh bảng xương CharInfo (33 cf) TỪ HÌNH HỌC ẢNH NGUỒN, hai chế độ:

  A  "chi-charinfo"  : giữ nguyên sprite + part.bin đang có, chỉ sinh CharInfo mới
                       (đầu đặt theo CẰM của khung thân; thân/chân đặt đúng vị trí khung nguồn).
  B  "tron-ven"      : thêm (1) tỉ lệ theo dải để đầu các dải bằng đầu idle,
                       (2) cắt thân tại CẰM − OV (không còn mặt trong mảnh thân),
                       (3) sinh sprite x1 + part.bin vào scratchpad, rồi CharInfo theo từng cf.

Cả hai chế độ: cf 0 và cf 1 chép NGUYÊN VĂN bảng cũ (NPC Iruka / màn chọn-tạo nhân vật dùng cf 0/1
không đổi), các slot idle (đầu 0-2, thân 0-2, chân 0-2) giữ nguyên sprite + dx/dy trong part.bin.

Phép vẽ client (Char.cs 6057-6084): x = cx + CI.dx + part.dx ; y = cy - CI.dy + part.dy (mép trên).
Suy ra CI.dx = L - part.dx ; CI.dy = part.dy - T, với (L, T) là góc trên-trái mong muốn so với (cx, cy).

Chạy:  python sinh-charinfo.py            (ghi charinfo-A.json, charinfo-B.json, sprites-B/, part-B.bin,
                                           charinfo-B.cs, do-khung.json rồi tự gọi ghep-khung.py cho 3 bảng)
"""
from __future__ import annotations

import json
import math
import shutil
import subprocess
import sys
from pathlib import Path

import numpy as np
from PIL import Image
from scipy import ndimage as ndi

GOC = Path(r"D:\NRO_GoiBanGiao")
if str(GOC) not in sys.path:
    sys.path.insert(0, str(GOC))

from assetpack.part import Frame, Part, encode_part, parse_part  # noqa: E402
from tools.artqa.task6b_closure import (  # noqa: E402
    PLAYER_ACTIONS,
    PLAYER_FEET_OFFSET_X1,
    _frame_rect,
    _load_reviewed_alpha,
    _review_by_name,
    zoom_variants,
)

ART_ROOT = GOC / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = GOC / "docs/my-thuat/review/vertical-slice-v1.json"
MAPPING = GOC / "docs/my-thuat/mapping-id.json"
CHARINFO_CU = GOC / "tools/nhan-vat/charinfo.json"
PART_BIN_CU = GOC / "server/internal/game/assets/part.bin"
SPRITES_CU_X1 = Path(r"D:\NRO_Project\Assets\Resources\res\x1\smallimage")  # chỉ ĐỌC
GHEP_KHUNG = GOC / "tools/nhan-vat/ghep-khung.py"
OUT = Path(__file__).resolve().parent

PART_IDS = {"head": 30, "body": 57, "leg": 58}
FIRST_ID = {"head": 27041, "body": 27044, "leg": 27061}
FEET = PLAYER_FEET_OFFSET_X1  # -11
S_IDLE = 128 / 340  # player_uniform_scale(idle0 width 340)
OV_NECK_X1 = 2      # thân chồng lên cằm 2 px x1 (che mép dưới cằm, không che miệng)
CUT_LEG_TOP, CUT_BODY_BOTTOM = 0.66, 0.76  # giữ như pipeline (thân-chân cùng khung → khớp tự nhiên)
BG = np.array([246, 4, 242])

# ---- gán slot theo ngữ nghĩa cf (đọc Char.cs): chạy = cf 2..6 ; 20/31 = khung BẮT ĐẦU skill (đất/bay)
THAN_DOI = {2: 3, 20: 11, 31: 11}      # cf 2 → run0 (thay idle2 trong vòng chạy) ; 20/31 → cast0 (kết ấn)
CHAN_THEO_THAN = True                   # chân lấy cùng khung với thân khi khung đó có chân (slot ≤ 13)


# ------------------------------------------------------------------ đo ảnh nguồn
def _hsv(a):
    a = a.astype(np.float32) / 255.0
    r, g, b = a[..., 0], a[..., 1], a[..., 2]
    mx, mn = a.max(-1), a.min(-1)
    d = mx - mn
    dd = np.where(d == 0, 1, d)
    h = np.zeros_like(mx)
    m = d > 1e-6
    rm = m & (mx == r)
    gm = m & (mx == g) & ~rm
    bm = m & ~rm & ~gm
    h[rm] = (60 * ((g - b) / dd) % 360)[rm]
    h[gm] = (60 * ((b - r) / dd) + 120)[gm]
    h[bm] = (60 * ((r - g) / dd) + 240)[bm]
    s = np.where(mx > 0, d / np.where(mx == 0, 1, mx), 0)
    return h, s, mx


def _bbox(m):
    ys, xs = np.where(m)
    if len(ys) == 0:
        return None
    return int(xs.min()), int(ys.min()), int(xs.max()) + 1, int(ys.max()) + 1


def do_mat(rgb_strip: np.ndarray, rect):
    """Đo MẶT (khối da lớn nhất có tâm ở 55% trên hộp nội dung) của một khung → toạ độ STRIP (px nguồn)."""
    x, y, w, h = rect
    fr = rgb_strip[y:y + h, x:x + w]
    fg = np.abs(fr.astype(int) - BG).sum(-1) > 90
    fg = ndi.binary_opening(fg, iterations=1)
    hh, ss, vv = _hsv(fr)
    skin = fg & (hh > 20) & (hh < 48) & (ss > 0.35) & (ss < 0.8) & (vv > 0.55)
    skin = ndi.binary_opening(skin, iterations=1)
    lab, n = ndi.label(skin)
    if n == 0:
        return None
    sizes = ndi.sum(skin, lab, range(1, n + 1))
    for k in np.argsort(-sizes):
        comp = lab == (k + 1)
        ys, xs = np.where(comp)
        if ys.mean() < 0.55 * h:
            fx0, fy0, fx1, fy1 = _bbox(comp)
            holes = ndi.binary_fill_holes(comp) & ~comp
            eye = None
            lab2, n2 = ndi.label(holes)
            if n2:
                s2 = ndi.sum(holes, lab2, range(1, n2 + 1))
                eb = _bbox(lab2 == (int(np.argmax(s2)) + 1))
                eye = (eb[1] + eb[3]) / 2
            return {
                "chin_y": y + fy1, "face_top": y + fy0, "face_cx": x + (fx0 + fx1) / 2,
                "face_h": fy1 - fy0, "face_w": fx1 - fx0,
                "eye_to_chin": (fy1 - eye) if eye is not None else None,
            }
    return None


# ------------------------------------------------------------------ khung nguồn
def khung_nguon(records):
    """Danh sách 17 khung theo thứ tự slot thân (idle0..hit2) với rect + đo mặt + strip RGB/alpha."""
    frames = []
    cache = {}
    for action, source, count in PLAYER_ACTIONS:
        rec = records[source]
        rgb = np.array(Image.open(ART_ROOT / source).convert("RGB"))
        alpha = _load_reviewed_alpha(ART_ROOT, records, source, cache)  # RGBA đã tách nền
        for i in range(count):
            rect = _frame_rect(rec, i)
            frames.append({
                "action": action, "source": source, "i": i, "rect": rect,
                "mat": do_mat(rgb, rect), "alpha": alpha,
            })
    return frames


def ti_le_dai(frames):
    """k_dải = median( idle / khung ) trên {face_h, face_w, eye_to_chin} → nhân để đầu bằng đầu idle."""
    idle = [f["mat"] for f in frames if f["action"] == "idle"]
    ref = {k: np.mean([m[k] for m in idle]) for k in ("face_h", "face_w", "eye_to_chin")}
    k = {"idle": 1.0}
    for action in {f["action"] for f in frames} - {"idle"}:
        ratios = []
        for f in frames:
            if f["action"] != action or f["mat"] is None:
                continue
            for key in ref:
                v = f["mat"].get(key)
                if v:
                    ratios.append(ref[key] / v)
        k[action] = round(float(np.median(ratios)), 3) if ratios else 1.0
    return k


def khung_scaled(f, scale):
    """Crop alpha khung, scale NEAREST như pipeline; trả (ảnh, center_x4, feet_x4)."""
    x, y, w, h = f["rect"]
    im = f["alpha"].crop((x, y, x + w, y + h)).resize(
        (max(1, round(w * scale)), max(1, round(h * scale))), Image.Resampling.NEAREST)
    bb = im.getchannel("A").getbbox()
    return im, (bb[0] + bb[2]) / 2.0, float(bb[3])


def cat_manh(im, top_x4, bottom_x4):
    """Cắt hàng [top, bottom) rồi trim alpha; trả (mảnh, left_x4, top_x4 thật)."""
    top = max(0, int(round(top_x4)))
    bottom = min(im.height, int(round(bottom_x4)))
    piece = im.crop((0, top, im.width, bottom))
    bb = piece.getchannel("A").getbbox()
    piece = piece.crop(bb)
    return piece, bb[0], top + bb[1]


def dat(left_x4, top_x4, center_x4, feet_x4):
    """(L, T) x1 so với (cx, cy) — đúng công thức component_placement của pipeline."""
    return round((left_x4 - center_x4) / 4), round((top_x4 - feet_x4) / 4) + FEET


# ------------------------------------------------------------------ hình học từng slot
def hinh_hoc(frames, k_dai, mode):
    """Trả dict slot → {'L','T','w','h','img_x4'|None,'chin_rel','face_rel','face_off'} cho head/body/leg."""
    geo = {"head": {}, "body": {}, "leg": {}}
    for j, f in enumerate(frames):
        scale = S_IDLE * (k_dai[f["action"]] if mode == "B" else 1.0)
        im, cx4, feet4 = khung_scaled(f, scale)
        x, y, w, h = f["rect"]
        H = im.height
        mat = f["mat"]
        chin_x4 = (mat["chin_y"] - y) * scale if mat else None
        face_x4 = (mat["face_cx"] - x) * scale if mat else None
        # làm tròn y như task6b_closure._component_rect + player_component_images:
        # lát cắt tính theo px NGUỒN (round(h*pct)) rồi mới nhân scale và round lần nữa
        def x4(v_src):
            return round(v_src * scale)
        # thân
        if mode == "B" and f["action"] != "idle" and mat is not None:
            body_top = round(chin_x4) - OV_NECK_X1 * 4      # cắt tại cằm − OV (không còn mặt)
        else:
            body_top = x4(round(h * 0.28))
        body, bl, bt = cat_manh(im, body_top, x4(round(h * CUT_BODY_BOTTOM)))
        L, T = dat(bl, bt, cx4, feet4)
        geo["body"][j] = {
            "L": L, "T": T, "w": body.width, "h": body.height, "img": body, "frame": j,
            "chin_rel": (round((chin_x4 - feet4) / 4) + FEET) if chin_x4 is not None else None,
            "face_rel": round((face_x4 - cx4) / 4) if face_x4 is not None else None,
            "khung_h_x1": H / 4,
        }
        # chân (14 khung đầu)
        if j < 14:
            leg, ll, lt = cat_manh(im, x4(round(h * CUT_LEG_TOP)), H)
            L, T = dat(ll, lt, cx4, feet4)
            geo["leg"][j] = {"L": L, "T": T, "w": leg.width, "h": leg.height, "img": leg, "frame": j}
        # đầu (3 khung idle) — luôn cắt 0..42% như hiện hành để giữ nguyên sprite đang có
        if j < 3:
            head, hl, ht = cat_manh(im, 0, x4(round(h * 0.42)))
            L, T = dat(hl, ht, cx4, feet4)
            geo["head"][j] = {
                "L": L, "T": T, "w": head.width, "h": head.height, "img": head, "frame": j,
                # kích cỡ ảnh x1 mà client vẽ (zoom_variants: round(x4/4)) — MỌI phép đặt dùng x1
                "h_x1": max(1, round(head.height / 4)), "w_x1": max(1, round(head.width / 4)),
                # cằm nằm cách mép DƯỚI ảnh đầu bao nhiêu px x1 (lát cắt 42% hơi thấp hơn cằm)
                "chin_off": round((ht + head.height - chin_x4) / 4),
                "face_off": round((face_x4 - hl) / 4),
            }
    return geo


# ------------------------------------------------------------------ CharInfo
def gan_slot(ci_cu):
    """idx (đầu, chân, thân) cho 33 cf theo ngữ nghĩa mới; cf 0/1 giữ nguyên."""
    idx = []
    for cf, row in enumerate(ci_cu):
        h, l, b = row[0][0], row[1][0], row[2][0]
        if cf >= 2:
            b = THAN_DOI.get(cf, b)
            if CHAN_THEO_THAN and b <= 13:
                l = b
        idx.append((h, l, b))
    return idx


def sinh_charinfo(ci_cu, idx, geo, part_cu, mode):
    """Trả (charinfo mới, part frames mới {comp: [(sid,dx,dy)]}, ghi chú)."""
    part_moi = {c: [list(fr) for fr in part_cu[PART_IDS[c]][1]] for c in PART_IDS}
    K = {"head": 0, "leg": 1, "body": 2}
    # slot nào được giữ nguyên dx/dy trong part.bin: idle 0-2 của cả 3 loại + (mode A) mọi slot
    for comp in ("body", "leg"):
        for j, g in geo[comp].items():
            if mode == "A" or j < 3:
                continue  # A: part.bin không đổi ; idle: giữ nguyên
            anchors = [cf for cf in range(33) if idx[cf][K[comp]] == j and cf >= 2]
            if not anchors:
                continue  # slot không cf nào dùng (vd thân 15 khung nằm) → giữ như cũ
            a = anchors[0]
            row = ci_cu[a][K[comp]]
            part_moi[comp][j][1] = g["L"] - row[1]   # dx sao cho CI của cf neo giữ số cũ (L = CI.dx + part.dx)
            part_moi[comp][j][2] = row[2] + g["T"]   # dy (T = part.dy - CI.dy)
    # kiểm tra hình học mode A khớp part.bin hiện hành (slot bất kỳ): L = CI.dx + part.dx của cf neo cũ
    ci = []
    ghi_chu = []
    for cf in range(33):
        if cf < 2:
            ci.append([list(r) for r in ci_cu[cf]])
            continue
        h, l, b = idx[cf]
        gb, gl, gh = geo["body"][b], geo["leg"][l], geo["head"][h]
        pb, pl, ph = part_moi["body"][b], part_moi["leg"][l], part_moi["head"][h]
        # thân, chân: đặt đúng vị trí khung nguồn của chính slot
        row_b = [b, gb["L"] - pb[1], pb[2] - gb["T"]]
        row_l = [l, gl["L"] - pl[1], pl[2] - gl["T"]]
        # đầu: cằm đầu = cằm khung thân ; tâm mặt đầu = tâm mặt khung thân
        if gb["chin_rel"] is None:
            ghi_chu.append(f"cf {cf}: khung thân slot {b} không đo được mặt → đặt đầu theo mép thân +{OV_NECK_X1}")
            head_bottom = gb["T"] + 8
            face_rel = gb["L"] + round(gb["w"] / 8)
        else:
            head_bottom = gb["chin_rel"] + gh["chin_off"]
            face_rel = gb["face_rel"]
        T_h = head_bottom - gh["h_x1"]
        L_h = face_rel - gh["face_off"]
        row_h = [h, L_h - ph[1], ph[2] - T_h]
        rows = [row_h, row_l, row_b]
        if len(ci_cu[cf]) > 3:
            rows.append(list(ci_cu[cf][3]))
        ci.append(rows)
    return ci, part_moi, ghi_chu


# ------------------------------------------------------------------ xuất
def xuat_cs(ci):
    out = ["\t\tpublic static int[][][] CharInfo = new int[33][][]", "\t\t{"]
    for cf, rows in enumerate(ci):
        out.append("\t\t\tnew int[4][]")
        out.append("\t\t\t{")
        r4 = rows[3] if len(rows) > 3 else None
        body = [f"\t\t\t\tnew int[3] {{ {r[0]}, {r[1]}, {r[2]} }}" for r in rows[:3]]
        body.append(f"\t\t\t\tnew int[3] {{ {r4[0]}, {r4[1]}, {r4[2]} }}" if r4 else "\t\t\t\tnew int[3]")
        out.append(",\n".join(body))
        out.append("\t\t\t}" + ("," if cf < 32 else ""))
    out.append("\t\t};")
    return "\n".join(out) + "\n"


def ghi_part(part_cu, part_moi, duong):
    parts = list(part_cu)
    for comp, pid in PART_IDS.items():
        t = part_cu[pid].part_type
        parts[pid] = Part(t, tuple(Frame(int(s), int(dx), int(dy)) for s, dx, dy in part_moi[comp]))
    duong.write_bytes(encode_part(parts))


def chay_ghep(charinfo, out_png, part_bin=None, sprites=None):
    cmd = [sys.executable, str(GHEP_KHUNG), "--charinfo", str(charinfo), "--out", str(out_png)]
    if part_bin:
        cmd += ["--part-bin", str(part_bin)]
    if sprites:
        cmd += ["--sprites", str(sprites)]
    r = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8", errors="replace")
    if r.returncode:
        print(r.stderr)
    return r.stdout


def main():
    review = json.loads(REVIEW.read_text(encoding="utf-8"))
    records = _review_by_name(review)
    ci_cu = json.loads(CHARINFO_CU.read_text(encoding="utf-8"))
    raw = PART_BIN_CU.read_bytes()
    part_cu_parsed = parse_part(raw)
    part_cu = [(p.part_type, [(f.sprite_id, f.dx, f.dy) for f in p.frames]) for p in part_cu_parsed]
    mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
    geo_map = {c["id"]: c for c in mapping["build_inputs"]["player_component_geometry"]["components"]}

    frames = khung_nguon(records)
    k_dai = ti_le_dai(frames)
    print("== tỉ lệ nhân theo dải để đầu bằng idle:", k_dai)
    print("== đo mặt (px nguồn, toạ độ strip):")
    for j, f in enumerate(frames):
        m = f["mat"]
        print("  slot thân %2d %-6s#%d rect=%s mặt=%s" % (j, f["action"], f["i"], f["rect"],
              None if m is None else {k: (round(v, 1) if isinstance(v, float) else v) for k, v in m.items()}))

    idx = gan_slot(ci_cu)
    ket_qua = {"k_dai": k_dai, "OV_NECK_X1": OV_NECK_X1, "slot_idx_moi": idx}

    # ---------- A: chỉ CharInfo
    geoA = hinh_hoc(frames, k_dai, "A")
    # tự kiểm: hình học A tái tạo đúng placement_x1 trong mapping-id.json (34/34)
    sai = []
    for comp, first in FIRST_ID.items():
        for j, g in geoA[comp].items():
            pl = geo_map[first + j]["placement_x1"]
            if (g["L"], g["T"]) != (pl["left_rel"], pl["top_rel"]):
                sai.append((first + j, (g["L"], g["T"]), (pl["left_rel"], pl["top_rel"])))
    print("== tự kiểm hình học A so mapping-id.json: %s" % ("khớp 34/34" if not sai else f"LỆCH {sai}"))
    ciA, partA, ghiA = sinh_charinfo(ci_cu, idx, geoA, part_cu, "A")
    (OUT / "charinfo-A.json").write_text(json.dumps(ciA), encoding="utf-8")
    for g in ghiA:
        print("  A:", g)

    # ---------- B: trọn vẹn
    geoB = hinh_hoc(frames, k_dai, "B")
    ciB, partB, ghiB = sinh_charinfo(ci_cu, idx, geoB, part_cu, "B")
    for g in ghiB:
        print("  B:", g)
    (OUT / "charinfo-B.json").write_text(json.dumps(ciB), encoding="utf-8")
    (OUT / "charinfo-B.cs").write_text(xuat_cs(ciB), encoding="utf-8", newline="\n")
    (OUT / "charinfo-A.cs").write_text(xuat_cs(ciA), encoding="utf-8", newline="\n")
    ghi_part(part_cu_parsed, partB, OUT / "part-B.bin")
    spr = OUT / "sprites-B"
    spr.mkdir(exist_ok=True)
    doi_sprite = []
    for comp, first in FIRST_ID.items():
        for j, g in geoB[comp].items():
            sid = first + j
            if j < 3 or (comp == "body" and j == 15):
                shutil.copy(SPRITES_CU_X1 / f"Small{sid}.png", spr / f"Small{sid}.png")  # giữ nguyên
                continue
            x1 = zoom_variants(g["img"])["x1"]
            x1.save(spr / f"Small{sid}.png")
            doi_sprite.append((sid, comp, j, g["img"].size, x1.size))
    print("== B: sinh lại %d sprite (x4 → x1 NEAREST); giữ nguyên 9 slot idle + thân 15" % len(doi_sprite))
    for d in doi_sprite:
        print("   ", d)
    print("== B: part.bin dx/dy đổi ở:")
    for comp, pid in PART_IDS.items():
        for j, (s, dx, dy) in enumerate(part_cu[pid][1]):
            n = partB[comp][j]
            if (dx, dy) != (n[1], n[2]):
                print("    %s slot %2d sprite %d: (%d,%d) → (%d,%d)" % (comp, j, s, dx, dy, n[1], n[2]))

    # ---------- đo bằng công cụ của dự án
    print("\n##### HIỆN TRẠNG (charinfo.json + part.bin + sprite client)")
    print(chay_ghep(CHARINFO_CU, OUT / "ghep-hien-trang.png"))
    print("\n##### A: CharInfo mới, sprite + part.bin GIỮ NGUYÊN")
    print(chay_ghep(OUT / "charinfo-A.json", OUT / "ghep-A.png"))
    print("\n##### B: CharInfo mới + tỉ lệ dải + cắt cằm + part-B.bin + sprites-B")
    print(chay_ghep(OUT / "charinfo-B.json", OUT / "ghep-B.png", OUT / "part-B.bin", spr))

    # khe ngang (tâm mặt đầu so tâm mặt khung thân) — ghep-khung không đo, tự tính
    def khe_ngang(ci, geo, part):
        out = {}
        for cf in range(2, 33):
            h, l, b = ci[cf][0][0], ci[cf][1][0], ci[cf][2][0]
            gb, gh = geo["body"][b], geo["head"][h]
            if gb["face_rel"] is None:
                continue
            ph = part["head"][h]
            L_h = ci[cf][0][1] + ph[1]
            out[cf] = (L_h + gh["face_off"]) - gb["face_rel"]
        return out
    ket_qua["khe_ngang_A"] = khe_ngang(ciA, geoA, partA)
    ket_qua["khe_ngang_B"] = khe_ngang(ciB, geoB, partB)
    print("== khe ngang tâm mặt (đầu − thân, px x1) A:", ket_qua["khe_ngang_A"])
    print("== khe ngang tâm mặt (đầu − thân, px x1) B:", ket_qua["khe_ngang_B"])
    (OUT / "do-khung.json").write_text(json.dumps(ket_qua, indent=1, default=str), encoding="utf-8")


if __name__ == "__main__":
    main()
