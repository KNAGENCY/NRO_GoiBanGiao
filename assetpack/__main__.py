import argparse
import json
import sys
from pathlib import Path

from .manifest import ManifestError, validate_manifest, verify_legacy_inputs


def main(argv=None):
    parser = argparse.ArgumentParser(prog="python -m assetpack")
    commands = parser.add_subparsers(dest="command", required=True)
    validate = commands.add_parser("validate")
    validate.add_argument("--manifest", type=Path, required=True)
    legacy = commands.add_parser("verify-legacy-inputs")
    legacy.add_argument("--ledger", type=Path, required=True)
    args = parser.parse_args(argv)
    if args.command == "validate":
        try:
            validate_manifest(json.loads(args.manifest.read_text(encoding="utf-8")))
        except (ManifestError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    if args.command == "verify-legacy-inputs":
        try:
            verify_legacy_inputs(args.ledger, Path.cwd())
        except (ManifestError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
