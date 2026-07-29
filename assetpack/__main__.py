import argparse
import json
import sys
from pathlib import Path

from .manifest import ManifestError, validate_manifest, verify_legacy_inputs
from .build import BuildError, build
from .deploy import DeployError, deploy, unity_preflight


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
    build_command.add_argument("--base-dir", type=Path)
    build_command.add_argument("--out", type=Path, required=True)
    build_command.add_argument("--review", type=Path, required=True)
    deploy_command = commands.add_parser("deploy")
    deploy_command.add_argument("--server-artifact", type=Path, required=True)
    deploy_command.add_argument("--client-artifact", type=Path)
    deploy_command.add_argument("--release-root", type=Path, required=True)
    deploy_command.add_argument("--snapshot", type=Path, required=True)
    deploy_command.add_argument("--version-ledger", type=Path, required=True)
    deploy_command.add_argument("--previous-lock", type=Path)
    deploy_command.add_argument("--runtime-key-map", type=Path, required=True)
    deploy_command.add_argument("--unity-exe", type=Path, required=True)
    deploy_command.add_argument("--protected-unity-root", type=Path)
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
            build(manifest, args.out, args.review, base_dir=args.base_dir or args.manifest.parent)
        except (BuildError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    if args.command == "deploy":
        try:
            release = deploy(
                args.server_artifact,
                args.client_artifact or args.server_artifact,
                args.release_root,
                args.snapshot,
                args.version_ledger,
                args.runtime_key_map,
                previous_lock_path=args.previous_lock,
                protected_unity_root=args.protected_unity_root,
                preflight=lambda snapshot, runtime_map, contract, log: unity_preflight(
                    args.unity_exe,
                    snapshot,
                    runtime_map,
                    contract,
                    log,
                ),
            )
            print(release)
        except (DeployError, OSError, json.JSONDecodeError) as error:
            print(error, file=sys.stderr)
            return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
