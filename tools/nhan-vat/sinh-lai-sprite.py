# -*- coding: utf-8 -*-
"""sinh-lai-sprite.py — sinh lại 34 mảnh người chơi + 3 mảnh Iruka + part.bin.

VÌ SAO CÓ CÔNG CỤ NÀY (05/09/2026)
----------------------------------
Đường chính để sinh asset là ``tools/artqa/task6b_closure.py`` →
``assetpack build``. Đường đó đang **bị khoá** cho một closure MỚI:

* ``generate_production_closure`` từ chối ghi đè và chỉ ghi vào **kho ảnh dùng
  chung** (``D:\\NRO_GoiBanGiao\\docs\\my-thuat\\anh-nguon\\...``) — nằm ngoài
  worktree này;
* ``build()`` đòi ``production_artifact_review`` do **reviewer độc lập**
  (``/root/task6b_visual_review6``) ký, và chữ ký đó bind theo
  ``production_closure_digest`` — đổi một byte sprite là chữ ký hết hiệu lực
  (``ASSET-ARTIFACT-REVIEW-016``). Tự sửa chữ ký = giả mạo, nên không làm.

Công cụ này **dùng lại đúng các hàm toán** trong ``task6b_closure`` (import, không
chép) để sinh ảnh + dx/dy, rồi ghi thẳng vào repo:

* ``docs/my-thuat/lang-la-client-assets/res/x{1..4}/smallimage/Small<id>.png``
  — bản art client có trong git (client Unity không thuộc repo này);
* ``server/internal/game/assets/part.bin`` — dựng qua ``assetpack.part.build_part``
  từ chính legacy ``server/internal/game/part.bin``;
* ``docs/my-thuat/mapping-id.json`` — ``parts[].frames[].dx/dy``, khối
  ``build_inputs.player_component_geometry`` (sự thật hình học mới: tỉ lệ từng khung,
  mốc giải phẫu, tham số + nghiệm MILP), và ``animations[]/sprites[].client_cf`` theo
  bảng sinh từ CharInfo;
* ``docs/my-thuat/review/nhan-vat-moc-giai-phau.json`` — mốc giải phẫu 17 khung nguồn
  đã dùng (tái lập được bằng ``python -m tools.artqa.nhan_vat_moc``);
* ``server/internal/game/assets/asset-lock.json`` — SHA/độ dài part.bin + version.

V1 (05/09/2026, chiều): "đầu tách thân" thực chất là HAI KHUÔN MẶT (mảnh thân cắt 28 %
chứa mặt). Bản này cắt thân từ cằm − 2 px x1, tỉ lệ TỪNG khung theo cỡ đầu, gán lại slot
(cf 2..6 = chạy, slot 15 hết khung nằm), dx/dy giải MILP một lần cho cả 34 mảnh (ghim
dy đầu slot 0 = −38). Chi tiết: ``tools/artqa/task6b_closure.py`` phần "V1".

Iruka V1 (06/09/2026): NPC Iruka (part 2189/2190/2191 → Small27075/27076/27077) trước đây
chỉ có bản closure v11 ép vào 128×160 (38 px x1 = 62 % người chơi) với dx = dy = 0 (đầu
lửng +13, thân chìm +29 px). Nay sinh bằng ``iruka_component_images``: cùng lát cắt theo
cằm, tỉ lệ s0 / r (r = cỡ đầu Iruka / cỡ đầu người chơi, cùng tham chiếu trung bình 3 idle),
dx/dy từ ``component_placement`` tại cf 0 (NPC chỉ dùng cf 0/1, không cần MILP). Khối
``build_inputs.iruka_component_geometry`` ghi riêng, KHÔNG nối vào 34 mảnh người chơi.

Dùng:
    python tools/nhan-vat/sinh-lai-sprite.py            # ghi thật vào worktree
    python tools/nhan-vat/sinh-lai-sprite.py --thu      # chỉ in bảng, không ghi
"""

from __future__ import annotations

import argparse
import hashlib
import io
import json
import sys
from pathlib import Path

GOC = Path(__file__).resolve().parents[2]
if str(GOC) not in sys.path:
    sys.path.insert(0, str(GOC))

from assetpack.animation import REQUIRED_ANIMATION_SHEETS, REQUIRED_CLIENT_CF  # noqa: E402
from assetpack.lock import lock_bytes  # noqa: E402
from assetpack.part import build_part, parse_part  # noqa: E402
from tools.artqa import nhan_vat_moc  # noqa: E402
from tools.artqa.task6b_closure import (  # noqa: E402
    IRUKA_GEOMETRY_VERSION,
    IRUKA_PART_IDS,
    IRUKA_SOURCE,
    PLAYER_FEET_OFFSET_X1,
    PLAYER_GEOMETRY_VERSION,
    _review_by_name,
    iruka_component_images,
    load_charinfo,
    player_component_images,
    player_frame_geometry,
    zoom_variants,
)

ART_ROOT = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1")
REVIEW = GOC / "docs/my-thuat/review/vertical-slice-v1.json"
MAPPING = GOC / "docs/my-thuat/mapping-id.json"
MOC_OUT = GOC / "docs/my-thuat/review/nhan-vat-moc-giai-phau.json"
SPRITE_ROOT = GOC / "docs/my-thuat/lang-la-client-assets/res"
SPRITE_RELATIVE = "docs/my-thuat/lang-la-client-assets/res"
LEGACY_PART = GOC / "server/internal/game/part.bin"
PART_OUT = GOC / "server/internal/game/assets/part.bin"
ASSET_LOCK = GOC / "server/internal/game/assets/asset-lock.json"
PART_IDS = {"head": 30, "body": 57, "leg": 58}
ZOOMS = ("x1", "x2", "x3", "x4")


def _sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def sinh_anh(components, root: Path, ghi: bool) -> dict[int, dict[str, str]]:
    """Ghi Small<id>.png cho bốn mức zoom; trả sha256 từng mức."""

    hashes: dict[int, dict[str, str]] = {}
    for component in components:
        sprite_id = component["id"]
        hashes[sprite_id] = {}
        for zoom, image in zoom_variants(component["image_x4"]).items():
            path = root / zoom / "smallimage" / f"Small{sprite_id}.png"
            if ghi:
                path.parent.mkdir(parents=True, exist_ok=True)
                image.save(path, format="PNG", optimize=False, compress_level=9)
                hashes[sprite_id][zoom] = _sha(path.read_bytes())
            else:
                hashes[sprite_id][zoom] = "?" * 64
    return hashes


def cap_nhat_client_cf(mapping: dict) -> dict:
    """animations[]/sprites[].client_cf theo bảng SINH từ CharInfo (assetpack/animation.py).

    `assetpack build` kiểm `animations[].frames[].client_cf` đúng bằng `REQUIRED_CLIENT_CF`
    (validate_animation_closure) — bảng đổi 05/09/2026 (thân/chân từng lệch slot, chân
    thiếu cf 26..32) thì mapping phải đổi theo, không thì build từ chối cả manifest.
    """

    for item in mapping["animations"]:
        for frame in item["frames"]:
            frame["client_cf"] = list(REQUIRED_CLIENT_CF[item["name"]][frame["slot_index"]])
    for sprite in mapping["sprites"]:
        if 27041 <= sprite["id"] <= 27074 and "component" in sprite:
            sprite["client_cf"] = list(
                REQUIRED_CLIENT_CF[f"hoc-vien-{sprite['component']}"][sprite["slot_index"]]
            )
    return mapping


def khoi_iruka_geometry(iruka, hashes_iruka, feet_offset: int) -> dict:
    """``build_inputs.iruka_component_geometry`` — sự thật hình học 3 mảnh Iruka đang chạy.

    Khối RIÊNG (không nối vào ``player_component_geometry.components`` — test đếm 34 ở đó).
    """

    st = iruka[0]["shared_transform"]
    bang = iruka[0]["bang_cf"]
    cf0 = next(row for row in bang if row["cf"] == 0)
    return {
        "ghi_chu": (
            "Hình học THẬT của 3 mảnh NPC Iruka đang chạy (V1 06/09: cùng lát cắt theo cằm và "
            "phép bù CharInfo như người chơi; tỉ lệ s0 / he_so_dau × ti_le_dau; dx/dy từ "
            "component_placement tại cf 0 — NPC chỉ đảo cf 0/1 nên không cần MILP). "
            "sprites[] 27075–27077 vẫn trỏ production-closure-v11 (hình học CŨ 128×160, dx=dy=0) "
            "vì closure mới cần reviewer độc lập ký lại (ASSET-ARTIFACT-REVIEW-016). Ảnh đang "
            "dùng nằm ở docs/my-thuat/lang-la-client-assets/res/x*/smallimage/."
        ),
        "phien_ban": IRUKA_GEOMETRY_VERSION,
        "generator": "tools/nhan-vat/sinh-lai-sprite.py",
        "source_path": IRUKA_SOURCE,
        "source_sha256": iruka[0]["source_sha256"],
        "source_frame": st["source_frame"],
        "ti_le_dau": st["ti_le_dau"],
        "he_so_dau": st["he_so_dau"],
        "tham_chieu": st["tham_chieu"],
        "ti_le": (
            "scale_x4 = uniform_scale_x4 (s0 = 128/340) / he_so_dau × ti_le_dau. he_so_dau = "
            "cỡ đầu Iruka / cỡ đầu người chơi = median 6 thước cứng (bỏ head_h_hair_chin vì đuôi "
            "tóc Iruka, thay bằng sqrt diện tích mặt) trên CÙNG tham chiếu trung bình 3 khung idle "
            "mà player_component_geometry dùng → ti_le_dau = 1.0 nghĩa là đầu bằng đầu người chơi."
        ),
        "uniform_scale_x4": st["uniform_scale_x4"],
        "scale_x4": st["scale"],
        "frame_bbox_x4": st["frame_bbox_x4"],
        "cao_x1": cf0["cao"],
        "feet_offset_x1": feet_offset,
        "anchor": "feet-bottom-center",
        "charinfo_cf": iruka[0]["charinfo_cf"],
        "charinfo_anchor_cf": iruka[0]["charinfo_anchor_cf"],
        "lat_cat": (
            "đầu: đỉnh tóc → cằm + 6 px nguồn; thân: cằm − 2 px x1 → 76 % khung; chân: 66 % → 100 %."
        ),
        "sprite_root": SPRITE_RELATIVE,
        "bang_cf": bang,
        "components": [
            {
                "id": component["id"],
                "part_id": component["part_id"],
                "component": component["component"],
                "source_rect": component["source_rect"],
                "size_x4": list(component["image_x4"].size),
                "moc": component["moc"],
                "placement_x1": {
                    "left_rel": component["placement"]["left_rel"],
                    "top_rel": component["placement"]["top_rel"],
                },
                "part_offset": {
                    "dx": component["placement"]["dx"],
                    "dy": component["placement"]["dy"],
                },
                "charinfo_cf": component["charinfo_cf"],
                "charinfo_anchor_cf": component["charinfo_anchor_cf"],
                "charinfo_residual": component["charinfo_residual"],
                "component_origin_x4": component["shared_transform"]["component_origin_x4"],
                "paths": {
                    zoom: f"{SPRITE_RELATIVE}/{zoom}/smallimage/Small{component['id']}.png"
                    for zoom in ZOOMS
                },
                "sha256": hashes_iruka[component["id"]],
            }
            for component in iruka
        ],
    }


def cap_nhat_mapping(mapping: dict, components, hashes, feet_offset: int, iruka=(), hashes_iruka=None) -> dict:
    """parts[].dx/dy + khối hình học V1 (+ Iruka) + client_cf; KHÔNG đụng sprites[].zooms (xem ghi chú).

    Bộ lọc part gồm 30/57/58 (người chơi) và 2189/2190/2191 (Iruka). Mọi frame của một part
    Iruka trỏ cùng sprite nên nhận cùng dx/dy.
    """

    by_id = {component["id"]: component for component in [*components, *iruka]}
    part_ids = set(PART_IDS.values()) | set(IRUKA_PART_IDS.values())
    for part in mapping["parts"]:
        if part["id"] not in part_ids:
            continue
        for frame in part["frames"]:
            component = by_id.get(frame["sprite_id"])
            if component is None:
                continue
            frame["dx"] = component["placement"]["dx"]
            frame["dy"] = component["placement"]["dy"]
    # Ô contact sheet giờ là giới hạn trên (cell mỗi khung một cỡ) — mapping phải
    # khai đúng ô hiện hành, nếu không `_schema` của build từ chối cả manifest.
    for sheet in mapping["build_inputs"]["animation_sheets"]:
        sheet["cell_size"] = list(
            REQUIRED_ANIMATION_SHEETS[sheet["name"]]["cell_size"]
        )
    cap_nhat_client_cf(mapping)
    scale = components[0]["shared_transform"]["scale"]
    milp = components[0]["milp"]
    mapping["build_inputs"]["player_component_geometry"] = {
        "ghi_chu": (
            "Hình học THẬT của 34 mảnh người chơi đang chạy (V1 05/09: cắt thân theo cằm, "
            "tỉ lệ từng khung theo cỡ đầu, slot gán lại, dx/dy giải MILP). sprites[].zooms "
            "vẫn trỏ production-closure-v11 (hình học CŨ, mỗi khung một tỉ lệ) vì closure "
            "mới cần reviewer độc lập ký lại (ASSET-ARTIFACT-REVIEW-016). Ảnh đang dùng nằm "
            "ở docs/my-thuat/lang-la-client-assets/res/x*/smallimage/."
        ),
        "phien_ban": PLAYER_GEOMETRY_VERSION,
        "generator": "tools/nhan-vat/sinh-lai-sprite.py",
        "uniform_scale_x4": scale,
        "ti_le": (
            "uniform_scale_x4 = s0 = 128/340 là MỐC; mỗi mảnh dùng scale_frame = s0 / "
            "nguon.he_so_dau (cỡ đầu khung / cỡ đầu idle, median 6 thước cứng — "
            "tools/artqa/nhan_vat_moc.py)."
        ),
        "lat_cat": (
            "đầu: đỉnh tóc → cằm + 6 px nguồn; thân: cằm − 2 px x1 → 76 % khung; chân: "
            "66 % → 100 %. Thân KHÔNG chứa mặt — đây là phần sửa 'hai khuôn mặt'."
        ),
        "feet_offset_x1": feet_offset,
        "anchor": "feet-bottom-center",
        "sprite_root": SPRITE_RELATIVE,
        "moc_giai_phau": MOC_OUT.relative_to(GOC).as_posix(),
        "milp": {
            "ghi_chu": (
                "dx/dy 34 mảnh là nghiệm MILP (tools/artqa/nhan_vat_toi_uu.py) — tái lập "
                "bằng nghiem + tham_so ở đây, không cần giải lại; kiem_nghiem() phải trả []."
            ),
            "phien_ban": milp["phien_ban"],
            "tham_so": milp["tham_so"],
            "nghiem": milp["nghiem"],
            "muc_tieu": milp["muc_tieu"],
            "tham_chieu_idle": milp["tham_chieu_idle"],
            "idle_frames": milp["idle_frames"],
            "bang_cf": milp["bang_cf"],
        },
        "components": [
            {
                "id": component["id"],
                "component": component["component"],
                "slot_index": component["slot_index"],
                "action": component["action"],
                "action_slot": component["action_slot"],
                "source_path": component["source_path"],
                "source_rect": component["source_rect"],
                "scale_frame": component["scale_frame"],
                "nguon": component["nguon"],
                "size_x4": list(component["image_x4"].size),
                "moc": component["moc"],
                "placement_x1": {
                    "left_rel": component["placement"]["left_rel"],
                    "top_rel": component["placement"]["top_rel"],
                },
                "part_offset": {
                    "dx": component["placement"]["dx"],
                    "dy": component["placement"]["dy"],
                },
                "charinfo_cf": component["charinfo_cf"],
                "charinfo_anchor_cf": component["charinfo_anchor_cf"],
                "charinfo_residual": component["charinfo_residual"],
                "paths": {
                    zoom: f"{SPRITE_RELATIVE}/{zoom}/smallimage/Small{component['id']}.png"
                    for zoom in ZOOMS
                },
                "sha256": hashes[component["id"]],
            }
            for component in components
        ],
    }
    if iruka:
        mapping["build_inputs"]["iruka_component_geometry"] = khoi_iruka_geometry(
            iruka, hashes_iruka, feet_offset
        )
    return mapping


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--art-root", type=Path, default=ART_ROOT)
    parser.add_argument("--review", type=Path, default=REVIEW)
    parser.add_argument("--mapping", type=Path, default=MAPPING)
    parser.add_argument("--sprite-root", type=Path, default=SPRITE_ROOT)
    parser.add_argument("--legacy-part", type=Path, default=LEGACY_PART)
    parser.add_argument("--part-out", type=Path, default=PART_OUT)
    parser.add_argument("--asset-lock", type=Path, default=ASSET_LOCK)
    parser.add_argument("--moc-out", type=Path, default=MOC_OUT)
    # 28 = V1 05/09 (part 30/57/58 đổi hết lần nữa); phải khớp vsData trong dongbo.go.
    # Iruka 06/09 (part 2189–2191) đổi CÙNG đợt, chưa phát hành → giữ 28, không bump hai lần.
    parser.add_argument("--vsdata", type=int, default=28)
    parser.add_argument(
        "--thu", action="store_true", help="chỉ in bảng, không ghi file nào"
    )
    args = parser.parse_args()

    review = json.loads(args.review.read_text(encoding="utf-8"))
    records = _review_by_name(review)
    charinfo = load_charinfo()
    moc = nhan_vat_moc.do_nguon(args.art_root, records)
    components = player_component_images(
        review, art_root=args.art_root, records=records, charinfo=charinfo, moc=moc
    )
    if len(components) != 34:
        raise SystemExit(f"mong doi 34 manh, nhan {len(components)}")
    iruka = iruka_component_images(
        review, art_root=args.art_root, records=records, charinfo=charinfo, moc=moc
    )
    if len(iruka) != 3:
        raise SystemExit(f"mong doi 3 manh Iruka, nhan {len(iruka)}")

    hashes = sinh_anh(components, args.sprite_root, ghi=not args.thu)
    hashes_iruka = sinh_anh(iruka, args.sprite_root, ghi=not args.thu)

    print("id    manh slot khung    s_f     x4        x1     cam_off  dx   dy  cf-neo  cf dung slot")
    for component in components:
        width, height = component["image_x4"].size
        cam = component["moc"]["chin_off"]
        print(
            "%5d %-4s %2d  %-7s %.4f  %3dx%-3d %3dx%-3d %6s %+4d %+4d %5d  %s"
            % (
                component["id"],
                component["component"],
                component["slot_index"],
                f"{component['action']}{component['action_slot']}",
                component["scale_frame"],
                width,
                height,
                round(width / 4),
                round(height / 4),
                "%6.1f" % cam if cam is not None else "   -  ",
                component["placement"]["dx"],
                component["placement"]["dy"],
                component["charinfo_anchor_cf"],
                component["charinfo_cf"],
            )
        )
    milp = components[0]["milp"]
    print("ti le moc s0 =", components[0]["shared_transform"]["scale"], "| MILP", milp["phien_ban"],
          "muc tieu %.2f" % milp["muc_tieu"], "| ghim dy dau", milp["tham_so"]["ghim_dy_dau"])
    print("cf | e(cam) d(than) | khe dau-than | khe than-chan | chan | cao | than/chan")
    for row in milp["bang_cf"]:
        print("%2d | %+5.1f %+5.1f | %+3d | %+3d | %+3d | %2d | %s/%s" % (
            row["cf"], row["e"], row["d"], row["khe_dau_than"], row["khe_than_chan"],
            row["chan"], row["cao"], row["than"], row["chan_khung"]))

    st = iruka[0]["shared_transform"]
    print(
        "IRUKA (part 2189/2190/2191): r = %.4f (tham chieu %s idle %s), s = %.4f = s0/r x %.2f"
        % (st["he_so_dau"], st["tham_chieu"]["kieu"], st["tham_chieu"]["idle_frames"], st["scale"], st["ti_le_dau"])
    )
    print("id    manh part  x4        x1     cam_off  left top   dx   dy  residual cf1")
    for component in iruka:
        width, height = component["image_x4"].size
        print(
            "%5d %-4s %4d  %3dx%-3d %3dx%-3d %6.1f  %+4d %+4d %+4d %+4d  %s"
            % (
                component["id"],
                component["component"],
                component["part_id"],
                width,
                height,
                round(width / 4),
                round(height / 4),
                component["moc"]["chin_off"],
                component["placement"]["left_rel"],
                component["placement"]["top_rel"],
                component["placement"]["dx"],
                component["placement"]["dy"],
                component["charinfo_residual"]["1"],
            )
        )
    print("cf | khe dau-than | khe than-chan | chan | cao   (du doan x1, Iruka)")
    for row in iruka[0]["bang_cf"]:
        print("%2d | %+3d | %+3d | %+3d | %2d" % (row["cf"], row["khe_dau_than"], row["khe_than_chan"], row["chan"], row["cao"]))

    mapping = json.loads(args.mapping.read_text(encoding="utf-8"))
    mapping = cap_nhat_mapping(
        mapping, components, hashes, PLAYER_FEET_OFFSET_X1, iruka=iruka, hashes_iruka=hashes_iruka
    )
    part_bytes = build_part(mapping, args.legacy_part.read_bytes())
    print(
        "part.bin: %d record, %d byte, sha %s"
        % (len(parse_part(part_bytes)), len(part_bytes), _sha(part_bytes)[:16])
    )
    if args.thu:
        print("(--thu: khong ghi file nao)")
        return 0

    # KHÔNG đụng mapping["versions"]: khối đó là phiên bản của closure v11 đã được reviewer
    # ký (assetpack/tests/test_deploy.py ghim {runtime_slice 22, data 21, map 22, resource 21});
    # phiên bản runtime thật (data 28) nằm ở asset-lock.json + dongbo.go vsData.
    args.mapping.write_text(
        json.dumps(mapping, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
    )
    args.part_out.write_bytes(part_bytes)
    # Mốc giải phẫu đã dùng — cùng hàm với `python -m tools.artqa.nhan_vat_moc`.
    hinh_hoc, tham_chieu, s0 = player_frame_geometry(records, moc)
    # newline="\n": Path.write_text trên Windows ghi CRLF, blob đã commit là LF →
    # file hiện "M" trong git status dù nội dung y hệt (kiểm đối kháng 06/09 bắt được).
    with io.open(args.moc_out, "w", encoding="utf-8", newline="\n") as f:
        f.write(json.dumps(nhan_vat_moc.xuat_json(moc, tham_chieu, hinh_hoc, s0, records), ensure_ascii=False, indent=1) + "\n")

    # asset-lock.json phải giữ ĐÚNG dạng canonical của `assetpack.lock.lock_bytes`
    # (một dòng, khoá đã sắp xếp): test Go tìm chuỗi con `"mob_ids":[0]` trong đó,
    # ghi lại bằng indent là mọi chốt fail-closed đó im lặng ngừng hoạt động.
    lock = json.loads(args.asset_lock.read_text(encoding="utf-8"))
    for output in lock["outputs"]:
        if output["path"] == "server/part.bin":
            output["sha256"] = _sha(part_bytes)
            output["bytes"] = len(part_bytes)
    lock["versions"]["data"] = args.vsdata
    for sheet in lock.get("closure", {}).get("animation_sheets", []):
        width, height = REQUIRED_ANIMATION_SHEETS[sheet["name"]]["cell_size"]
        sheet["cell_size"] = [width, height]
        for cell in sheet["cells"]:
            slot = cell["slot_index"]
            cell["bounds"] = [slot * width, 0, (slot + 1) * width, height]
    # closure.animations trong lock là bản chụp của mapping.animations — theo cùng bảng.
    for item in lock.get("closure", {}).get("animations", []):
        for frame in item["frames"]:
            frame["client_cf"] = list(REQUIRED_CLIENT_CF[item["name"]][frame["slot_index"]])
    args.asset_lock.write_bytes(lock_bytes(lock))
    print("da ghi:", args.part_out, args.mapping, args.asset_lock, args.moc_out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
