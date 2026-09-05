# -*- coding: utf-8 -*-
"""Cầu nối thẳng tới Unity MCP (streamable HTTP, 127.0.0.1:8080/mcp) bằng JSON-RPC.
Dùng: python unity_mcp.py <ten_tool> '<json_args>'   |  python unity_mcp.py --list
Vì phiên Claude Code này mất kết nối MCP 'UnityMCP' lúc mở (Unity chưa bật), gọi tay qua HTTP.
"""
import json, sys, urllib.request

URL = "http://127.0.0.1:8080/mcp"

def goi(method, params=None, sid=None, rid=1):
    body = {"jsonrpc": "2.0", "id": rid, "method": method}
    if params is not None:
        body["params"] = params
    req = urllib.request.Request(URL, data=json.dumps(body).encode("utf-8"), method="POST")
    req.add_header("Content-Type", "application/json")
    req.add_header("Accept", "application/json, text/event-stream")
    if sid:
        req.add_header("Mcp-Session-Id", sid)
    with urllib.request.urlopen(req, timeout=120) as r:
        sid2 = r.headers.get("Mcp-Session-Id") or sid
        raw = r.read().decode("utf-8", "replace")
        ct = r.headers.get("Content-Type", "")
    if "text/event-stream" in ct:
        # lấy khối data: cuối cùng
        du = [l[5:].strip() for l in raw.splitlines() if l.startswith("data:")]
        raw = du[-1] if du else "{}"
    return json.loads(raw) if raw.strip() else {}, sid2

def thong_bao(method, params, sid):
    body = {"jsonrpc": "2.0", "method": method, "params": params}
    req = urllib.request.Request(URL, data=json.dumps(body).encode("utf-8"), method="POST")
    req.add_header("Content-Type", "application/json")
    req.add_header("Accept", "application/json, text/event-stream")
    req.add_header("Mcp-Session-Id", sid)
    try:
        urllib.request.urlopen(req, timeout=30).read()
    except Exception:
        pass

def phien():
    r, sid = goi("initialize", {"protocolVersion": "2025-03-26", "capabilities": {},
                                "clientInfo": {"name": "claude-code-cau", "version": "1"}})
    if "error" in r:
        raise SystemExit("initialize loi: " + json.dumps(r["error"]))
    thong_bao("notifications/initialized", {}, sid)
    return sid

def main():
    sid = phien()
    if len(sys.argv) >= 2 and sys.argv[1] == "--list":
        r, _ = goi("tools/list", {}, sid, 2)
        for t in r.get("result", {}).get("tools", []):
            print(t["name"])
        return
    if len(sys.argv) >= 3 and sys.argv[1] == "--schema":
        r, _ = goi("tools/list", {}, sid, 2)
        for t in r.get("result", {}).get("tools", []):
            if t["name"] == sys.argv[2]:
                print(t.get("description", "")[:1500])
                print(json.dumps(t.get("inputSchema", {}), ensure_ascii=False, indent=1)[:4000])
        return
    if sys.argv[1] == "--code":
        # python unity_mcp.py --code file.cs  -> execute_code voi than ham C# doc tu file (tranh quoting PowerShell)
        ten = "execute_code"
        args = {"action": "execute", "code": open(sys.argv[2], encoding="utf-8").read()}
    elif sys.argv[1] == "--json":
        # python unity_mcp.py --json file.json -> {"name":..., "arguments":{...}}
        d = json.load(open(sys.argv[2], encoding="utf-8"))
        ten, args = d["name"], d.get("arguments", {})
    else:
        ten = sys.argv[1]
        args = json.loads(sys.argv[2]) if len(sys.argv) >= 3 else {}
    r, _ = goi("tools/call", {"name": ten, "arguments": args}, sid, 3)
    if "error" in r:
        print("LOI:", json.dumps(r["error"], ensure_ascii=False)); return
    res = r.get("result", {})
    for c in res.get("content", []):
        if c.get("type") == "text":
            print(c["text"])
        else:
            print(json.dumps({k: (v[:80] + "..." if isinstance(v, str) and len(v) > 80 else v) for k, v in c.items()}, ensure_ascii=False))
    if res.get("isError"):
        print("(isError)")

if __name__ == "__main__":
    main()
