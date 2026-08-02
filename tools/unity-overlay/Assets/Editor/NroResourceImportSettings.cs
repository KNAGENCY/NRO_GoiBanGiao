#if UNITY_EDITOR
using UnityEditor;
using UnityEngine;

namespace NroE2E
{
    public sealed class NroResourceImportSettings : AssetPostprocessor
    {
        private void OnPreprocessTexture()
        {
            string normalized = assetPath.Replace('\\', '/');
            if (!normalized.StartsWith(
                    "Assets/Resources/res/x",
                    System.StringComparison.OrdinalIgnoreCase))
            {
                return;
            }

            TextureImporter importer = (TextureImporter)assetImporter;
            importer.npotScale = TextureImporterNPOTScale.None;
            importer.mipmapEnabled = false;
            importer.alphaIsTransparency = true;
            importer.filterMode = FilterMode.Point;
            importer.wrapMode = TextureWrapMode.Clamp;
        }
    }
}
#endif
