import hashlib
import json
import tempfile
import unittest
from pathlib import Path
from unittest import mock

from assetpack.deploy import DeployError, deploy


def _sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _pair_sha(root: Path) -> str:
    digest = hashlib.sha256()
    paths = sorted((root / "server").rglob("*")) + sorted((root / "client").rglob("*"))
    for path in paths:
        if path.is_file():
            digest.update(path.relative_to(root).as_posix().encode("utf-8"))
            digest.update(path.read_bytes())
    return digest.hexdigest()


def _write_json(path: Path, payload: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(
        json.dumps(payload, ensure_ascii=False, sort_keys=True),
        encoding="utf-8",
    )


class DeployArtifactTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.artifact = self.root / "artifact"
        self.release_root = self.root / "releases"
        self.snapshot = self.root / "unity-snapshot"
        self.snapshot.mkdir()
        _write_json(self.snapshot / ".nro-e2e-overlay.json", {"schema_version": 1})
        self.runtime_map = self.root / "runtime-key-map.json"
        _write_json(
            self.runtime_map,
            {
                "schema_version": 1,
                "entries": [
                    {
                        "source_asset": f"ui-{index}.png",
                        "mode": "direct",
                        "logical_keys": [f"/ui/{index}"],
                        "future_logical_key": None,
                        "physical_paths": {
                            f"x{zoom}": [f"Assets/Resources/res/x{zoom}/ui/{index}.png"]
                            for zoom in range(1, 5)
                        },
                    }
                    for index in range(8)
                ],
            },
        )
        self.ledger = self.root / "asset-version-ledger.json"
        _write_json(
            self.ledger,
            {
                "runtime_slice_version": 0,
                "data_version": 0,
                "map_version": 0,
                "resource_version": 0,
            },
        )
        self._make_artifact(self.artifact)

    def tearDown(self):
        self.temp.cleanup()

    def _make_artifact(self, root: Path, *, versions=None) -> dict:
        files = {
            "server/part.bin": b"part-v1",
            "server/skill.bin": b"skill-v1",
            "server/tile_set_info.bin": b"tiles-v1",
            "server/map-1-collision.bin": b"collision-v1",
            "server/catalog.json": b'{"maps":[]}\n',
            "client/mymap/1.bytes": b"\x01\x01\x00",
        }
        for zoom in range(1, 5):
            prefix = f"client/x{zoom}"
            files[f"{prefix}/smallimage/Small27041.png"] = f"small-x{zoom}".encode()
            files[f"{prefix}/bg/b00.png"] = f"bg-x{zoom}".encode()
            files[f"{prefix}/mainImage/naruto-hp-bg.png"] = f"hud-x{zoom}".encode()
            files[f"{prefix}/Mob/0/img.png"] = f"mob-x{zoom}".encode()
            for index in range(8):
                files[f"{prefix}/ui/{index}.png"] = f"ui-{index}-x{zoom}".encode()
        files["client/x1/t/1/t_01.png"] = b"tile-x1"
        for path, data in files.items():
            destination = root / path
            destination.parent.mkdir(parents=True, exist_ok=True)
            destination.write_bytes(data)
        outputs = [
            {"path": path, "sha256": _sha(data), "bytes": len(data)}
            for path, data in sorted(files.items())
        ]
        lock = {
            "schema_version": 1,
            "slice": "lang-la",
            "versions": versions
            or {"runtime_slice": 1, "data": 1, "map": 1, "resource": 1},
            "outputs": outputs,
            "provenance": {
                "runtime_key_map": {
                    "sha256": _sha(self.runtime_map.read_bytes()),
                }
            },
            "server_client_pair_sha256": _pair_sha(root),
        }
        _write_json(root / "asset-lock.json", lock)
        return lock

    def _preflight_ok(
        self,
        snapshot: Path,
        runtime_map: Path,
        contract: Path,
        log_path: Path,
    ) -> int:
        self.assertEqual(snapshot.resolve(), self.snapshot.resolve())
        self.assertEqual(runtime_map.resolve(), self.runtime_map.resolve())
        payload = json.loads(contract.read_text(encoding="utf-8"))
        self.assertEqual(payload["tile_count"], 100)
        self.assertEqual(
            {resource["kind"] for resource in payload["resources"]},
            {"small", "background", "mob", "hud", "skill"},
        )
        log_path.write_text("RESOURCE-PREFLIGHT-OK\n", encoding="utf-8")
        return 0

    def test_tu_choi_server_client_khac_lock(self):
        client_artifact = self.root / "client-artifact"
        self._make_artifact(
            client_artifact,
            versions={"runtime_slice": 2, "data": 2, "map": 2, "resource": 2},
        )

        with self.assertRaisesRegex(DeployError, "ASSET-LOCK-MISMATCH-017"):
            deploy(
                self.artifact,
                client_artifact,
                self.release_root,
                self.snapshot,
                self.ledger,
                self.runtime_map,
                preflight=self._preflight_ok,
            )

    def test_tu_choi_runtime_doi_byte_nhung_version_khong_tang(self):
        previous_lock = self.root / "previous-lock.json"
        previous_lock.write_bytes((self.artifact / "asset-lock.json").read_bytes())
        part = self.artifact / "server/part.bin"
        mutated = b"Part-v1"
        part.write_bytes(mutated)
        lock = json.loads((self.artifact / "asset-lock.json").read_text(encoding="utf-8"))
        for output in lock["outputs"]:
            if output["path"] == "server/part.bin":
                output["sha256"] = _sha(mutated)
                output["bytes"] = len(mutated)
        lock["server_client_pair_sha256"] = _pair_sha(self.artifact)
        _write_json(self.artifact / "asset-lock.json", lock)

        with self.assertRaisesRegex(DeployError, "ASSET-VERSION-015"):
            deploy(
                self.artifact,
                self.artifact,
                self.release_root,
                self.snapshot,
                self.ledger,
                self.runtime_map,
                previous_lock_path=previous_lock,
                preflight=self._preflight_ok,
            )

    def test_tu_choi_byte_lech_sha_trong_lock(self):
        part = self.artifact / "server/part.bin"
        part.write_bytes(b"Part-v1")

        with self.assertRaisesRegex(DeployError, "ASSET-RUNTIME-SHA-018"):
            deploy(
                self.artifact,
                self.artifact,
                self.release_root,
                self.snapshot,
                self.ledger,
                self.runtime_map,
                preflight=self._preflight_ok,
            )

    def test_tu_choi_runtime_key_map_doi_sau_build(self):
        self.runtime_map.write_bytes(self.runtime_map.read_bytes() + b" ")

        with self.assertRaisesRegex(DeployError, "ASSET-RUNTIME-SHA-018"):
            deploy(
                self.artifact,
                self.artifact,
                self.release_root,
                self.snapshot,
                self.ledger,
                self.runtime_map,
                preflight=self._preflight_ok,
            )

    def test_tu_choi_ghi_de_snapshot_dang_chay(self):
        _write_json(self.snapshot / ".nro-asset-lock.json", {"active": True})

        with self.assertRaisesRegex(DeployError, "ASSET-DEPLOY-ACTIVE-019"):
            deploy(
                self.artifact,
                self.artifact,
                self.release_root,
                self.snapshot,
                self.ledger,
                self.runtime_map,
                preflight=self._preflight_ok,
            )

    def test_tu_choi_unity_goc_va_snapshot_thieu_overlay(self):
        unity_source = self.root / "unity-source"
        unity_source.mkdir()
        with self.assertRaisesRegex(DeployError, "ASSET-DEPLOY-SOURCE-020"):
            deploy(
                self.artifact,
                self.artifact,
                self.release_root,
                unity_source,
                self.ledger,
                self.runtime_map,
                protected_unity_root=unity_source,
                preflight=self._preflight_ok,
            )

        bare = self.root / "bare-snapshot"
        bare.mkdir()
        with self.assertRaisesRegex(DeployError, "ASSET-DEPLOY-SNAPSHOT-021"):
            deploy(
                self.artifact,
                self.artifact,
                self.release_root,
                bare,
                self.ledger,
                self.runtime_map,
                preflight=self._preflight_ok,
            )

    def test_deploy_cung_lock_vao_snapshot_va_luu_preflight_log(self):
        release = deploy(
            self.artifact,
            self.artifact,
            self.release_root,
            self.snapshot,
            self.ledger,
            self.runtime_map,
            preflight=self._preflight_ok,
        )

        self.assertTrue((release / "server/part.bin").is_file())
        self.assertEqual(
            (self.snapshot / "Assets/Resources/res/mymap/1.bytes").read_bytes(),
            b"\x01\x01\x00",
        )
        self.assertEqual(
            (release / "preflight.log").read_text(encoding="utf-8"),
            "RESOURCE-PREFLIGHT-OK\n",
        )
        marker = json.loads(
            (self.snapshot / ".nro-asset-lock.json").read_text(encoding="utf-8")
        )
        self.assertEqual(marker["lock_sha256"], release.name)

    def test_publish_release_loi_phai_rollback_ca_marker_va_client(self):
        deployed = self.snapshot / "Assets/Resources/res/mymap/1.bytes"
        with mock.patch.object(Path, "replace", side_effect=OSError("publish failed")):
            with self.assertRaisesRegex(OSError, "publish failed"):
                deploy(
                    self.artifact,
                    self.artifact,
                    self.release_root,
                    self.snapshot,
                    self.ledger,
                    self.runtime_map,
                    preflight=self._preflight_ok,
                )

        self.assertFalse((self.snapshot / ".nro-asset-lock.json").exists())
        self.assertFalse(deployed.exists())


class ProductionVersionContractTests(unittest.TestCase):
    def test_manifest_ghim_moc_21_va_khong_trung_legacy(self):
        manifest = json.loads(
            Path("docs/my-thuat/mapping-id.json").read_text(encoding="utf-8")
        )
        versions = manifest["versions"]
        self.assertEqual(
            versions,
            {"runtime_slice": 21, "data": 21, "map": 21, "resource": 21},
        )
        self.assertNotEqual(versions["data"], 10)
        self.assertNotEqual(versions["map"], 2)
        self.assertNotEqual(versions["resource"], 1)
        self.assertNotEqual(versions["runtime_slice"], 4)


if __name__ == "__main__":
    unittest.main()
