# -*- coding: utf-8 -*-
"""Giải mã thân gói -24 (hex từ log máy chủ) ĐÚNG như client Controller.cs case -24 + loadInfoMap đọc.
Dùng: python giai-ma-24.py <hex>   hoặc   python giai-ma-24.py --log <file log> (lấy dòng 'MAP INFO -24 hex' cuối)."""
import sys, struct, re, io

def main():
    if sys.argv[1] == "--log":
        s = io.open(sys.argv[2], encoding="utf-8", errors="replace").read()
        hx = re.findall(r'MAP INFO -24 hex.*?hex=([0-9a-f]+)', s)
        if not hx:
            raise SystemExit("khong thay dong hex trong log")
        h = hx[-1]
    else:
        h = sys.argv[1]
    d = bytes.fromhex(h)
    o = [0]
    def b():  v = d[o[0]]; o[0] += 1; return v
    def sb(): v = struct.unpack_from(">b", d, o[0])[0]; o[0] += 1; return v
    def sh(): v = struct.unpack_from(">h", d, o[0])[0]; o[0] += 2; return v
    def i4(): v = struct.unpack_from(">i", d, o[0])[0]; o[0] += 4; return v
    def utf():
        n = struct.unpack_from(">H", d, o[0])[0]; o[0] += 2
        s = d[o[0]:o[0]+n].decode("utf-8", "replace"); o[0] += n; return s
    print("tong", len(d), "byte")
    print("mapID", b(), "planet", sb(), "tile", sb(), "bg", sb(), "type", sb(), "name", utf(), "zone", sb())
    # loadInfoMap
    print("cx", sh(), "cy", sh())
    n = sb(); print("waypoints", n)
    for _ in range(n): print("  wp", sh(), sh(), sh(), sh(), b(), b(), utf())
    n = sb(); print("mobs", n, "@", o[0])
    for k in range(n):
        st = o[0]
        flags = [b() for _ in range(5)]
        tpl, sys_ = b(), b()
        hp = i4(); lv = sb(); mx = i4(); x = sh(); y = sh(); status = sb(); lvb = sb(); isBoss = b()
        print(f"  mob {k}: tpl={tpl} sys={sys_} hp={hp} lv={lv} max={mx} x={x} y={y} st={status} lvb={lvb} isBoss={isBoss} ({o[0]-st} byte)")
    n = sb(); print("du phong (doc nhu count) =", n)
    n = sb(); print("npcs", n, "@", o[0])
    for k in range(n):
        print(f"  npc {k}: status={sb()} cx={sh()} cy={sh()} tpl={sb()} avatar={sh()}")
    n = sb(); print("items", n)
    for k in range(n):
        print("  item", sh(), sh(), sh(), sh(), i4())
    n = sh(); print("bg items", n)
    for k in range(n): sh(); sh(); sh()
    n = sh(); print("effMap kv", n)
    for k in range(n): utf(); utf()
    print("bgType", sb(), "teleport", sb())
    print("con lai", len(d) - o[0], "byte:", d[o[0]:].hex())

main()
