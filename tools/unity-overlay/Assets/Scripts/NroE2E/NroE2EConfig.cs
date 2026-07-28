using System;
using System.IO;

namespace NroE2E
{
    public static class NroE2EConfig
    {
        public static string GetArgumentValue(string name)
        {
            string[] args = Environment.GetCommandLineArgs();
            for (int i = 0; i < args.Length - 1; i++)
            {
                if (string.Equals(args[i], name, StringComparison.Ordinal))
                {
                    return args[i + 1];
                }
            }
            return null;
        }

        public static string ResolveProfileRoot(
            string e2eFlag,
            string requestedRoot,
            string productionRoot)
        {
            if (!string.Equals(e2eFlag, "1", StringComparison.Ordinal))
            {
                return productionRoot;
            }
            if (string.IsNullOrWhiteSpace(requestedRoot) ||
                !Path.IsPathRooted(requestedRoot))
            {
                throw new InvalidOperationException(
                    "NRO_E2E=1 requires an absolute -profile-root.");
            }

            string resolved = Path.GetFullPath(requestedRoot);
            Directory.CreateDirectory(resolved);
            return resolved;
        }

        public static bool ShouldUseLegacyEffect(int effectId)
        {
            return effectId == 0;
        }

        public static bool HasSkillEffect(int effectId)
        {
            return effectId > 0;
        }
    }
}
