"""Generate provenance-bound tile crops from the reviewed Naruto v2 sheet."""

from __future__ import annotations

import argparse
import hashlib
import json
import shutil
import tempfile
from pathlib import Path

from PIL import Image


class TileSheetError(ValueError):
    """A reviewed tile-sheet input or generated output is invalid."""


ZOOMS = {"x1": 24, "x2": 48, "x3": 72, "x4": 96}
GRID_COLUMNS = 10
GRID_ROWS = 10
EXPECTED_SIZE = (1254, 1254)


def _sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for block in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(block)
    return digest.hexdigest()


def _error(code: str, detail: str) -> TileSheetError:
    return TileSheetError(f"{code}: {detail}")


def _lattice_boxes(width: int, height: int) -> list[tuple[int, int, int, int]]:
    boxes = []
    for row in range(GRID_ROWS):
        y0 = round(row * height / GRID_ROWS)
        y1 = round((row + 1) * height / GRID_ROWS)
        for column in range(GRID_COLUMNS):
            x0 = round(column * width / GRID_COLUMNS)
            x1 = round((column + 1) * width / GRID_COLUMNS)
            boxes.append((x0, y0, x1, y1))
    return boxes


def _nonblack_ratio(tile: Image.Image) -> float:
    rgb = tile.convert("RGB")
    pixels = list(rgb.get_flattened_data())
    if not pixels:
        return 0.0
    return sum(1 for pixel in pixels if max(pixel) > 8) / len(pixels)


def validate_source(source: Path, expected_sha256: str) -> dict:
    if not source.is_file():
        raise _error("ASSET-TILE-001", f"thiếu tile sheet đã review: {source}")
    actual_sha256 = _sha256(source)
    if actual_sha256 != expected_sha256:
        raise _error(
            "ASSET-TILE-001",
            f"SHA tile sheet lệch: expected {expected_sha256}, actual {actual_sha256}",
        )

    try:
        image = Image.open(source)
        image.load()
    except Exception as exc:  # pragma: no cover - Pillow supplies the detail
        raise _error("ASSET-TILE-002", f"không đọc được tile sheet: {source}") from exc
    if image.size != EXPECTED_SIZE:
        raise _error(
            "ASSET-TILE-002",
            f"kích thước tile sheet phải là {EXPECTED_SIZE}, nhận {image.size}",
        )
    if image.mode not in {"RGB", "RGBA"}:
        raise _error("ASSET-TILE-002", f"mode không được hỗ trợ: {image.mode}")

    ratios = []
    boxes = _lattice_boxes(*image.size)
    for index, box in enumerate(boxes):
        ratio = _nonblack_ratio(image.crop(box))
        ratios.append(ratio)
        if ratio < 0.10:
            raise _error(
                "ASSET-TILE-002",
                f"ô lattice {index} gần rỗng: nonblack_ratio={ratio:.6f}",
            )
    return {
        "source_sha256": actual_sha256,
        "source_dimensions": list(image.size),
        "lattice": {
            "columns": GRID_COLUMNS,
            "rows": GRID_ROWS,
            "all_cells_nonempty": True,
            "nonblack_ratio_min": round(min(ratios), 6),
            "nonblack_ratio_max": round(max(ratios), 6),
        },
        "boxes": boxes,
    }


def generate_tile_sheet(
    source: Path | str, output: Path | str, *, expected_sha256: str
) -> dict:
    """Validate the source, then atomically generate 100 tiles at four zooms."""

    source = Path(source)
    output = Path(output)
    if output.exists():
        raise _error("ASSET-TILE-003", f"output đã tồn tại, không ghi đè: {output}")

    provenance = validate_source(source, expected_sha256)
    output.parent.mkdir(parents=True, exist_ok=True)
    temporary = Path(tempfile.mkdtemp(prefix=f".{output.name}-", dir=output.parent))
    try:
        with Image.open(source) as image:
            image.load()
            tiles = []
            for index, box in enumerate(provenance.pop("boxes")):
                row, column = divmod(index, GRID_COLUMNS)
                tile = image.crop(box).convert("RGBA")
                zooms = {}
                for zoom, size in ZOOMS.items():
                    relative_path = Path("tiles") / f"{index:03d}" / zoom / "tile.png"
                    destination = temporary / relative_path
                    destination.parent.mkdir(parents=True, exist_ok=True)
                    tile.resize((size, size), Image.Resampling.NEAREST).save(
                        destination, format="PNG"
                    )
                    zooms[zoom] = {
                        "path": relative_path.as_posix(),
                        "sha256": _sha256(destination),
                        "dimensions": [size, size],
                    }
                x0, y0, x1, y1 = box
                tiles.append(
                    {
                        "index": index,
                        "row": row,
                        "column": column,
                        "source_box": [x0, y0, x1 - x0, y1 - y0],
                        "zooms": zooms,
                    }
                )

        generated_index = {
            "schema_version": 1,
            "generator": "tools.artqa.generate_tile_sheet",
            **provenance,
            "tiles": tiles,
        }
        index_path = temporary / "generated-index.json"
        index_path.write_text(
            json.dumps(generated_index, ensure_ascii=False, indent=2) + "\n",
            encoding="utf-8",
        )
        temporary.replace(output)
        return generated_index
    except Exception:
        shutil.rmtree(temporary, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--sha256", required=True)
    args = parser.parse_args()
    try:
        generate_tile_sheet(args.source, args.output, expected_sha256=args.sha256)
    except TileSheetError as exc:
        parser.error(str(exc))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
