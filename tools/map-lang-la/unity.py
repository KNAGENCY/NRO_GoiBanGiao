# -*- coding: utf-8 -*-
"""Goi THANG HTTP toi MCP Unity — bo tang cau noi mcp-remote (hay chet, cham).

Dung:  from unity import U;  u = U();  print(u.cs('return "x";'))
"""
import json, re, time, urllib.request, uuid

URL = "http://127.0.0.1:8080/mcp"


class U:
    def __init__(self, url=URL):
        self.url = url
        self.sid = None
        st, raw, sid = self._goi("initialize", {
            "protocolVersion": "2024-11-05", "capabilities": {},
            "clientInfo": {"name": "truc-tiep", "version": "1"}})
        if not sid:
            raise RuntimeError(f"khong lay duoc session: {raw[:200]}")
        self.sid = sid
        self._goi("notifications/initialized", {})

    def _goi(self, method, params=None, timeout=60):
        body = {"jsonrpc": "2.0", "id": str(uuid.uuid4())[:8], "method": method}
        if params is not None:
            body["params"] = params
        req = urllib.request.Request(self.url, data=json.dumps(body).encode(), method="POST")
        req.add_header("Content-Type", "application/json")
        req.add_header("Accept", "application/json, text/event-stream")
        if self.sid:
            req.add_header("Mcp-Session-Id", self.sid)
        try:
            with urllib.request.urlopen(req, timeout=timeout) as r:
                return r.status, r.read().decode("utf-8", "replace"), r.headers.get("Mcp-Session-Id")
        except Exception as e:
            return None, f"LOI: {e}", None

    def tool(self, ten, args, timeout=60):
        st, raw, _ = self._goi("tools/call", {"name": ten, "arguments": args}, timeout)
        m = re.search(r'^data: (.+)$', raw, re.M)
        if not m:
            return {"loi": raw[:300]}
        try:
            o = json.loads(m.group(1))
        except Exception:
            return {"loi": m.group(1)[:300]}
        res = o.get("result", o)
        sc = res.get("structuredContent")
        if isinstance(sc, dict):
            return sc
        return res

    def cs(self, code, timeout=60):
        """Chay mot doan C# trong Unity, tra ve chuoi ket qua (hoac dict loi)."""
        r = self.tool("execute_code", {"action": "execute", "code": code, "compiler": "codedom"}, timeout)
        if isinstance(r, dict) and r.get("success") and isinstance(r.get("data"), dict):
            return r["data"].get("result")
        return r

    def play(self):
        return self.tool("manage_editor", {"action": "play"}, timeout=90)

    def stop(self):
        return self.tool("manage_editor", {"action": "stop"}, timeout=90)

    def refresh(self):
        return self.tool("refresh_unity", {"mode": "force", "scope": "assets"}, timeout=180)

    def chup(self, duong_dan, super_size=2):
        return self.cs(f'UnityEngine.ScreenCapture.CaptureScreenshot("{duong_dan}", {super_size}); return "ok";')

    def cho_vao_map(self, so_lan=24, giay=5):
        """Cho toi khi client vao map that (mapID != 0). Tra ve chuoi trang thai."""
        for _ in range(so_lan):
            kq = self.cs('return "map=" + Game1.TileMap.mapID + "|x=" + Game1.Char.myCharz().cx '
                         '+ "|y=" + Game1.Char.myCharz().cy + "|zoom=" + Game1.mGraphics.zoomLevel;')
            if isinstance(kq, str) and "map=0|" not in kq and kq.startswith("map="):
                return kq
            time.sleep(giay)
        return None
