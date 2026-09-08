import hashlib
import json
import subprocess
import sys
from pathlib import Path

from PIL import Image


SCRIPT = Path(__file__).with_name("lap_ba_mau_ui.py")
ZOOMS = {"x1": 0.25, "x2": 0.5, "x3": 0.75, "x4": 1.0}


def _png(path: Path, size: tuple[int, int], color: tuple[int, int, int, int]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    Image.new("RGBA", size, color).save(path)


def _sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _fixture(tmp_path: Path) -> dict[str, Path]:
    source = tmp_path / "source"
    zoom = tmp_path / "zoom"
    unity = tmp_path / "unity"
    backup = tmp_path / "backup"
    staging = tmp_path / "staging"
    contract_path = tmp_path / "contract.json"
    state_path = tmp_path / "unity-state.json"

    specs = {
        "ui_hud_khung_trai_x4.png": [12, 8],
        "ui_khung_popup_9o_x4.png": [4, 32],
        "ui_nut_thuong_3manh_x4.png": [12, 8],
    }
    colors = {
        "ui_hud_khung_trai_x4.png": (80, 40, 20, 255),
        "ui_khung_popup_9o_x4.png": (230, 170, 90, 255),
        "ui_nut_thuong_3manh_x4.png": (220, 100, 20, 255),
    }
    expected = {}
    for name, size in specs.items():
        path = source / name
        _png(path, tuple(size), colors[name])
        expected[name] = {"size": size, "sha256": _sha(path)}
        stem = name.removesuffix("_x4.png")
        for level, ratio in ZOOMS.items():
            scaled = (round(size[0] * ratio), round(size[1] * ratio))
            _png(zoom / level / f"{stem}_{level}.png", scaled, colors[name])

    contract_path.write_text(json.dumps({"sources": expected}), encoding="utf-8")
    state_path.write_text(
        json.dumps(
            {
                "data_path": str((unity / "Assets").resolve()),
                "is_playing": False,
                "is_playing_or_will_change": False,
                "is_compiling": False,
                "is_updating": False,
            }
        ),
        encoding="utf-8",
    )

    for level in ("x2", "x3", "x4"):
        ratio = ZOOMS[level]
        target_root = unity / "Assets" / "Resources" / "res" / level / "mainimage"
        _png(target_root / "myTexture2dpanel.png", (round(12 * ratio), round(8 * ratio)), (1, 2, 3, 255))
        _png(target_root / "myTexture2dimgPopup.png", (round(4 * ratio), round(32 * ratio)), (4, 5, 6, 255))
        for name in ("btn0left.png", "btn0mid.png", "btn0right.png"):
            _png(target_root / name, (round(4 * ratio), round(8 * ratio)), (7, 8, 9, 255))

    return {
        "source": source,
        "zoom": zoom,
        "unity": unity,
        "backup": backup,
        "staging": staging,
        "contract": contract_path,
        "state": state_path,
    }


def _run(fx: dict[str, Path], action: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [
            sys.executable,
            str(SCRIPT),
            action,
            "--source-root",
            str(fx["source"]),
            "--zoom-root",
            str(fx["zoom"]),
            "--unity-root",
            str(fx["unity"]),
            "--backup-root",
            str(fx["backup"]),
            "--staging-root",
            str(fx["staging"]),
            "--contract",
            str(fx["contract"]),
            "--unity-state-file",
            str(fx["state"]),
        ],
        text=True,
        encoding="utf-8",
        capture_output=True,
        check=False,
    )


def test_preflight_tu_choi_install_khi_backup_thieu(tmp_path: Path) -> None:
    """Bắt đột biến bỏ cổng backup rồi ghi thẳng vào Unity."""
    fx = _fixture(tmp_path)
    before = {
        p.relative_to(fx["unity"]): _sha(p)
        for p in fx["unity"].rglob("*.png")
    }

    result = _run(fx, "install")

    assert result.returncode == 2
    assert "UI-DOT4-BACKUP-MISSING" in result.stdout
    assert before == {
        p.relative_to(fx["unity"]): _sha(p)
        for p in fx["unity"].rglob("*.png")
    }


def test_preflight_tu_choi_khi_unity_dang_play(tmp_path: Path) -> None:
    """Bắt đột biến cho phép chép asset khi Unity đang Play."""
    fx = _fixture(tmp_path)
    state = json.loads(fx["state"].read_text(encoding="utf-8"))
    state["is_playing"] = True
    fx["state"].write_text(json.dumps(state), encoding="utf-8")

    result = _run(fx, "backup")

    assert result.returncode == 2
    assert "UI-DOT4-UNITY-PLAYING" in result.stdout
    assert not fx["backup"].exists()


def test_tu_choi_source_sai_sha_truoc_khi_backup(tmp_path: Path) -> None:
    """Bắt đột biến dùng một PNG cùng tên/cỡ nhưng không phải bản canonical."""
    fx = _fixture(tmp_path)
    _png(fx["source"] / "ui_hud_khung_trai_x4.png", (12, 8), (99, 40, 20, 255))

    result = _run(fx, "backup")

    assert result.returncode == 2
    assert "UI-DOT4-SOURCE-SHA" in result.stdout
    assert not fx["backup"].exists()


def test_doc_duoc_state_tu_lop_boc_json_cua_mcp_moi(tmp_path: Path) -> None:
    """Khóa đúng output execute_code: state nằm trong data.result."""
    fx = _fixture(tmp_path)
    bridge = tmp_path / "fake_unity_bridge.py"
    state = json.loads(fx["state"].read_text(encoding="utf-8"))
    bridge.write_text(
        "import json\n"
        f"state = {state!r}\n"
        "print(json.dumps({'success': True, 'data': "
        "{'result': 'UI_DOT4_STATE=' + json.dumps(state), 'compiler': 'codedom'}}))\n",
        encoding="utf-8",
    )

    result = subprocess.run(
        [
            sys.executable,
            str(SCRIPT),
            "backup",
            "--source-root",
            str(fx["source"]),
            "--zoom-root",
            str(fx["zoom"]),
            "--unity-root",
            str(fx["unity"]),
            "--backup-root",
            str(fx["backup"]),
            "--staging-root",
            str(fx["staging"]),
            "--contract",
            str(fx["contract"]),
            "--unity-bridge",
            str(bridge),
            "--unity-state-code",
            str(tmp_path / "state.cs"),
        ],
        text=True,
        encoding="utf-8",
        capture_output=True,
        check=False,
    )

    assert result.returncode == 0, result.stdout + result.stderr
    assert "UI-DOT4-BACKUP-CREATED count=15" in result.stdout


def test_backup_stage_install_dung_15_target_va_khong_tao_x1(tmp_path: Path) -> None:
    """Khóa đúng ba asset logic thành 15 file runtime x2/x3/x4."""
    fx = _fixture(tmp_path)
    original = {
        p.relative_to(fx["unity"]): _sha(p)
        for p in fx["unity"].rglob("*.png")
    }

    staged = _run(fx, "stage")
    backed_up = _run(fx, "backup")
    installed = _run(fx, "install")

    assert staged.returncode == 0, staged.stdout + staged.stderr
    assert backed_up.returncode == 0, backed_up.stdout + backed_up.stderr
    assert installed.returncode == 0, installed.stdout + installed.stderr
    assert "UI-DOT4-INSTALLED count=15" in installed.stdout
    assert len(list(fx["staging"].rglob("*.png"))) == 15
    assert not (fx["unity"] / "Assets" / "Resources" / "res" / "x1" / "mainimage").exists()

    manifest = json.loads((fx["backup"] / "backup-manifest.json").read_text(encoding="utf-8"))
    assert len(manifest["files"]) == 15
    for item in manifest["files"]:
        relative = Path(item["relative_path"])
        assert item["sha256"] == original[relative]
        assert _sha(fx["backup"] / relative) == original[relative]

    for level in ("x2", "x3", "x4"):
        runtime = fx["unity"] / "Assets" / "Resources" / "res" / level / "mainimage"
        with Image.open(runtime / "myTexture2dimgPopup.png") as popup:
            assert popup.size == (round(4 * ZOOMS[level]), round(32 * ZOOMS[level]))
        for name in ("btn0left.png", "btn0mid.png", "btn0right.png"):
            with Image.open(runtime / name) as button_piece:
                assert button_piece.size == (round(4 * ZOOMS[level]), round(8 * ZOOMS[level]))
