# -*- coding: utf-8 -*-
"""stop -> refresh -> sua import -> play -> bam Choi TK -> cho vao map -> chup."""
from unity import U
import os, sys, time

ten = sys.argv[1] if len(sys.argv) > 1 else "shot"
u = U()
u.stop(); time.sleep(4)
u.refresh(); time.sleep(12)
fix = ('int n=0; foreach (string z in new string[]{"x1","x2","x3","x4"}) {'
       ' foreach (string sub in new string[]{"bg","t"}) {'
       '  string dir = "Assets/Resources/res/" + z + "/" + sub;'
       '  if (!System.IO.Directory.Exists(dir)) continue;'
       '  foreach (string f in System.IO.Directory.GetFiles(dir, "*.png")) {'
       '   string p = f.Replace("\\\\", "/");'
       '   var ti = UnityEditor.AssetImporter.GetAtPath(p) as UnityEditor.TextureImporter;'
       '   if (ti == null) continue;'
       '   bool can = ti.npotScale != UnityEditor.TextureImporterNPOTScale.None || ti.mipmapEnabled'
       '     || ti.textureCompression != UnityEditor.TextureImporterCompression.Uncompressed'
       '     || (sub == "bg" && !ti.isReadable);'
       '   if (!can) continue;'
       '   ti.npotScale = UnityEditor.TextureImporterNPOTScale.None;'
       '   ti.textureCompression = UnityEditor.TextureImporterCompression.Uncompressed;'
       '   ti.mipmapEnabled = false; if (sub == "bg") ti.isReadable = true;'
       '   ti.SaveAndReimport(); n++; } } }'
       'return "sua " + n + " anh";')
print(" import:", u.cs(fix, timeout=180))
time.sleep(6)
u.play(); time.sleep(8)
s = u.cs('var sc = Game1.GameCanvas.currentScreen; return sc == null ? "null" : sc.GetType().Name;')
if s == "ServerListScreen":
    u.cs('(Game1.GameCanvas.currentScreen as Game1.IActionListener).perform(3, null); return "bam";')
vao = False
for i in range(12):
    time.sleep(5)
    kq = u.cs('return (Game1.GameCanvas.currentScreen == null ? "null" : Game1.GameCanvas.currentScreen.GetType().Name) + "|x=" + Game1.Char.myCharz().cx;')
    if isinstance(kq, str) and "GameScr" in kq:
        print(" vao game:", kq); vao = True; break
if not vao:
    print(" KHONG VAO DUOC GAME"); sys.exit(1)
time.sleep(4)
for f in [ten + "A", ten + "B"]:
    p = f"D:/NRO_Project/{f}.png"
    try: os.remove(p)
    except: pass
    u.chup(p, 1); time.sleep(7)
    print(f" {f}:", os.path.getsize(p) if os.path.exists(p) else "KHONG")
