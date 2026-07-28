import argparse
import json
import sys
from pathlib import Path

from .manifest import ManifestError, validate_manifest, verify_legacy_inputs
from .build import BuildError, build


def main(argv=None):
    parser = argparse.ArgumentParser(prog="python -m assetpack")
    commands = parser.add_subparsers(dest="command", required=True)
    validate = commands.add_parser("validate")
    validate.add_argument("--manifest", type=Path, required=True)
    validate.add_argument("--legacy-ledger", type=Path)
    validate.add_argument("--version-ledger", type=Path)
    validate.add_argument("--art-root", type=Path)
    legacy = commands.add_parser("verify-legacy-inputs")
    legacy.add_argument("--ledger", type=Path, required=True)
    build_command = commands.add_parser("build")
    build_command.add_argument("--manifest", type=Path, required=True)
    build_command.add_argument("--slice", required=True)
    build_command.add_argument("--out", type=Path, required=True)
    build_command.add_argument("--review", type=Path, required=True)
    args = parser.parse_args(argv)
    if args.command == "validate":
        try:
            if args.legacy_ledger is None:
                candidate = args.manifest.parent / "legacy-input-ledger.json"
                args.legacy_ledger = candidate if candidate.is_file() else None
            if args.version_ledger is None:
                candidate = args.manifest.parent / "asset-version-ledger.json"
                args.version_ledger = candidate if candidate.is_file() else None
            occupied_ids = ()
            if args.legacy_ledger:
                legacy_data = json.loads(args.legacy_ledger.read_text(encoding="utf-8"))
                occupied_ids = legacy_data.get("expected", {}).get("occupied_ids", ())
            previous_runtime_version = None
            if args.version_ledger:
                version_data = json.loads(args.version_ledger.read_text(encoding="utf-8"))
                previous_runtime_version = version_data.get("runtime_slice_version")
            validate_manifest(
                json.loads(args.manifest.read_text(encoding="utf-8")),
                base_dir=args.art_root or args.manifest.parent,
                occupied_ids=occupied_ids,
                previous_runtime_version=previous_runtime_version,
            )
        except (ManifestError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    if args.command == "verify-legacy-inputs":
        try:
            verify_legacy_inputs(args.ledger, Path.cwd())
        except (ManifestError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    if args.command == "build":
        try:
            manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
            if manifest.get("phase", {}).get("name", manifest.get("slice")) != args.slice:
                raise BuildError("ASSET-BUILD-011: slice không khớp manifest")
            build(manifest, args.out, args.review, base_dir=args.manifest.parent)
        except (BuildError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
