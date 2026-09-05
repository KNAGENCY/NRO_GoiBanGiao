import re, sys, subprocess
def chay(args):
    out=subprocess.run([sys.executable, r"D:\NRO_GoiBanGiao\tools\nhan-vat\ghep-khung.py"]+args, capture_output=True, text=True, encoding='utf-8', errors='replace').stdout
    rows=[]
    for ln in out.splitlines():
        m=re.match(r"\s*(\d+)\s*\|\s*(\d+)\s+(\d+)\s+(\d+)\s*\|\s*([+-]?\d+)\s*\|\s*([+-]?\d+)\s*\|\s*([+-]?\d+)", ln)
        if m: rows.append(tuple(int(x) for x in m.groups()))
    return rows
sp=sys.argv[1]
cases={
 "hien-trang x1": ["--out", sp+r"\wf-critic\o1.png"],
 "hien-trang x4": ["--zoom","4","--sprites", r"D:\NRO_Project\Assets\Resources\res\x4\smallimage","--out", sp+r"\wf-critic\o2.png"],
 "V1 x1": ["--part-bin", sp+r"\wf-thiet-ke\sprites\V1\part.bin","--sprites", sp+r"\wf-thiet-ke\sprites\V1","--out", sp+r"\wf-critic\o3.png"],
 "V1 x4": ["--zoom","4","--part-bin", sp+r"\wf-thiet-ke\sprites\V1\part.bin","--sprites", sp+r"\wf-thiet-ke\sprites\V1\x4","--out", sp+r"\wf-critic\o4.png"],
 "V1-siet(GK1) x4": ["--zoom","4","--part-bin", sp+r"\wf-giam-khao-1\part-V1-kl-2-d5.5.bin","--sprites", sp+r"\wf-thiet-ke\sprites\V1\x4","--out", sp+r"\wf-critic\o5.png"],
 "V1-siet(GK1) x1": ["--part-bin", sp+r"\wf-giam-khao-1\part-V1-kl-2-d5.5.bin","--sprites", sp+r"\wf-thiet-ke\sprites\V1","--out", sp+r"\wf-critic\o6.png"],
 "V2-B x1": ["--charinfo", sp+r"\wf-charinfo\charinfo-B.json","--part-bin", sp+r"\wf-charinfo\part-B.bin","--sprites", sp+r"\wf-charinfo\sprites-B","--out", sp+r"\wf-critic\o7.png"],
}
for ten,args in cases.items():
    rows=chay(args)
    if not rows: print(ten, "KHONG DOC DUOC BANG"); continue
    ht=[r[4] for r in rows]; tc=[r[5] for r in rows]; ch=[r[6] for r in rows]
    z=4 if "x4" in ten else 1
    print(f"{ten:18s} n={len(rows)} dau-than [{min(ht)},{max(ht)}] ho={sum(1 for v in ht if v>0)} | than-chan [{min(tc)},{max(tc)}] ho={sum(1 for v in tc if v>0)} cf_ho={[r[0] for r in rows if r[5]>0]} | chan {sorted(set(ch))}")
