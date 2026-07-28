#if UNITY_EDITOR
using System;
using System.IO;
using UnityEditor;
using UnityEngine;

namespace NroE2E
{
    public static class NroResourcePreflightTests
    {
        [Serializable]
        private sealed class RuntimeKeyMap
        {
            public int schema_version;
            public string unity_source_scope;
            public RuntimeBinding[] entries;
        }

        [Serializable]
        private sealed class RuntimeBinding
        {
            public string source_asset;
            public string mode;
            public string[] logical_keys;
            public string future_logical_key;
            public PhysicalPaths physical_paths;
        }

        [Serializable]
        private sealed class PhysicalPaths
        {
            public string[] x1;
            public string[] x2;
            public string[] x3;
            public string[] x4;

            public string[] ForZoom(int zoom)
            {
                switch (zoom)
                {
                    case 1: return x1;
                    case 2: return x2;
                    case 3: return x3;
                    case 4: return x4;
                    default: throw new ArgumentOutOfRangeException(nameof(zoom));
                }
            }
        }

        [Serializable]
        private sealed class ResourceContract
        {
            public int schema_version;
            public int tile_id;
            public int tile_count;
            public RequiredResource[] resources;
        }

        [Serializable]
        private sealed class RequiredResource
        {
            public string kind;
            public string name;
            public ZoomResource[] zooms;
        }

        [Serializable]
        private sealed class ZoomResource
        {
            public int zoom;
            public string physical_path;
            public int expected_width;
            public int expected_height;
        }

        private sealed class PreflightException : Exception
        {
            public readonly string Code;

            public PreflightException(string code, string message)
                : base(message)
            {
                Code = code;
            }
        }

        public static void RunBatch()
        {
            try
            {
                string runtimeMapPath = RequiredArgument("-runtime-key-map");
                string contractPath = RequiredArgument("-resource-contract");
                AssetDatabase.Refresh(ImportAssetOptions.ForceSynchronousImport);
                Game1.Main.typeClient = 4;

                RuntimeKeyMap runtimeMap = ReadJson<RuntimeKeyMap>(runtimeMapPath);
                ResourceContract contract = ReadJson<ResourceContract>(contractPath);
                ValidateRuntimeBindings(runtimeMap);
                ValidateTiles(contract);
                ValidateRequiredResources(contract);

                Debug.Log("RESOURCE-PREFLIGHT-OK");
                EditorApplication.Exit(0);
            }
            catch (PreflightException error)
            {
                Debug.LogError(error.Code + ": " + error.Message);
                EditorApplication.Exit(1);
            }
            catch (Exception error)
            {
                Debug.LogError("ASSET-PREFLIGHT-999: " + error);
                EditorApplication.Exit(1);
            }
        }

        private static void ValidateRuntimeBindings(RuntimeKeyMap runtimeMap)
        {
            if (runtimeMap == null || runtimeMap.schema_version != 1 ||
                runtimeMap.entries == null || runtimeMap.entries.Length != 8)
            {
                throw new PreflightException(
                    "ASSET-UI-007",
                    "Runtime key map must contain exactly eight schema-v1 bindings.");
            }

            for (int index = 0; index < runtimeMap.entries.Length; index++)
            {
                RuntimeBinding binding = runtimeMap.entries[index];
                bool overlayRequired = binding.mode == "overlay_required";
                bool directOrGroup = binding.mode == "direct" ||
                    binding.mode == "slice_group";
                if (!overlayRequired && !directOrGroup)
                {
                    throw new PreflightException(
                        "ASSET-UI-007",
                        "Unsupported runtime binding mode at entry " + index + ".");
                }
                if (overlayRequired)
                {
                    if (string.IsNullOrWhiteSpace(binding.future_logical_key))
                    {
                        throw new PreflightException(
                            "ASSET-UI-007",
                            "overlay_required entry has no future_logical_key.");
                    }
                }
                else if (binding.logical_keys == null || binding.logical_keys.Length == 0)
                {
                    throw new PreflightException(
                        "ASSET-UI-007",
                        "Direct/slice binding has no logical key.");
                }

                if (binding.physical_paths == null)
                {
                    throw new PreflightException(
                        "ASSET-UI-007",
                        "Binding has no physical paths: " + binding.source_asset);
                }
                for (int zoom = 1; zoom <= 4; zoom++)
                {
                    string[] paths = binding.physical_paths.ForZoom(zoom);
                    if (paths == null || paths.Length == 0)
                    {
                        throw new PreflightException(
                            "ASSET-UI-007",
                            "Binding is missing zoom x" + zoom + ": " + binding.source_asset);
                    }
                    foreach (string physicalPath in paths)
                    {
                        LoadTexture(
                            physicalPath,
                            ZoomErrorCode(zoom),
                            null,
                            null);
                    }
                }
            }
        }

        private static void ValidateTiles(ResourceContract contract)
        {
            if (contract == null || contract.schema_version != 1 ||
                contract.tile_id < 0 || contract.tile_count < 1)
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    "Invalid resource preflight contract.");
            }

            for (int zoom = 1; zoom <= 4; zoom++)
            {
                for (int tile = 1; tile <= contract.tile_count; tile++)
                {
                    string path;
                    if (zoom == 1)
                    {
                        path = string.Format(
                            "Assets/Resources/res/x1/t/{0}/t_{1:00}.png",
                            contract.tile_id,
                            tile);
                    }
                    else
                    {
                        path = string.Format(
                            "Assets/Resources/res/x{0}/t/{1}${2}.png",
                            zoom,
                            contract.tile_id,
                            tile);
                    }
                    LoadTexture(
                        path,
                        ZoomErrorCode(zoom),
                        24 * zoom,
                        24 * zoom);
                }
            }
        }

        private static void ValidateRequiredResources(ResourceContract contract)
        {
            if (contract.resources == null || contract.resources.Length == 0)
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    "Contract must include Small, background and HUD resources.");
            }

            bool hasSmall = false;
            bool hasBackground = false;
            bool hasHud = false;
            foreach (RequiredResource resource in contract.resources)
            {
                hasSmall |= resource.kind == "small";
                hasBackground |= resource.kind == "background";
                hasHud |= resource.kind == "hud";
                if (resource.zooms == null || resource.zooms.Length != 4)
                {
                    throw new PreflightException(
                        "ASSET-PREFLIGHT-001",
                        resource.name + " must declare four zooms.");
                }
                foreach (ZoomResource zoom in resource.zooms)
                {
                    LoadTexture(
                        zoom.physical_path,
                        ZoomErrorCode(zoom.zoom),
                        zoom.expected_width,
                        zoom.expected_height);
                }
            }
            if (!hasSmall || !hasBackground || !hasHud)
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    "Contract must cover small, background and hud kinds.");
            }
        }

        private static Texture2D LoadTexture(
            string physicalPath,
            string errorCode,
            int? expectedWidth,
            int? expectedHeight)
        {
            string key = ToResourceKey(physicalPath);
            Texture2D texture = Resources.Load<Texture2D>(key);
            if (texture == null)
            {
                throw new PreflightException(
                    errorCode,
                    "Resources.Load returned null for " + key + ".");
            }
            if ((expectedWidth.HasValue && texture.width != expectedWidth.Value) ||
                (expectedHeight.HasValue && texture.height != expectedHeight.Value))
            {
                throw new PreflightException(
                    errorCode,
                    string.Format(
                        "{0} expected {1}x{2}, actual {3}x{4}.",
                        key,
                        expectedWidth,
                        expectedHeight,
                        texture.width,
                        texture.height));
            }
            return texture;
        }

        private static string ToResourceKey(string physicalPath)
        {
            const string Prefix = "Assets/Resources/";
            string normalized = (physicalPath ?? string.Empty).Replace('\\', '/');
            if (!normalized.StartsWith(Prefix, StringComparison.Ordinal) ||
                !normalized.EndsWith(".png", StringComparison.OrdinalIgnoreCase))
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    "Physical path must be a PNG below Assets/Resources: " + physicalPath);
            }
            return normalized.Substring(
                Prefix.Length,
                normalized.Length - Prefix.Length - ".png".Length);
        }

        private static string RequiredArgument(string name)
        {
            string value = NroE2EConfig.GetArgumentValue(name);
            if (string.IsNullOrWhiteSpace(value) || !Path.IsPathRooted(value))
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    name + " must be an absolute path.");
            }
            return Path.GetFullPath(value);
        }

        private static T ReadJson<T>(string path)
        {
            if (!File.Exists(path))
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    "JSON input does not exist: " + path);
            }
            T value = JsonUtility.FromJson<T>(File.ReadAllText(path));
            if (value == null)
            {
                throw new PreflightException(
                    "ASSET-PREFLIGHT-001",
                    "Cannot parse JSON input: " + path);
            }
            return value;
        }

        private static string ZoomErrorCode(int zoom)
        {
            switch (zoom)
            {
                case 1: return "ASSET-ZOOM-001";
                case 2: return "ASSET-ZOOM-002";
                case 3: return "ASSET-ZOOM-003";
                case 4: return "ASSET-ZOOM-004";
                default:
                    throw new PreflightException(
                        "ASSET-PREFLIGHT-001",
                        "Zoom must be in 1..4.");
            }
        }
    }
}
#endif
