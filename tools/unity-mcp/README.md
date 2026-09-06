# Điều khiển Unity Editor từ Claude Code — cầu HTTP tới MCP for Unity

Dùng khi muốn **tự lái client Unity thay chủ dự án**: bấm Play/Stop, đăng nhập, đọc
trạng thái nhân vật/quái/NPC, chụp màn hình, đọc Console — tất cả không cần ai ngồi
trước máy. Đã dùng thật 05/09/2026 để tìm ra 8 lỗi (xem `BAT-DAU-TU-DAY.md` hàng
"Chơi thử thật").

## Vì sao có file này

Unity (gói "MCP for Unity") mở một máy chủ MCP **HTTP tại `127.0.0.1:8080/mcp`** khi
project `D:\NRO_Project` đang mở. Claude Code có cấu hình MCP `UnityMCP` trỏ tới đó,
nhưng nếu Unity chưa mở lúc phiên Claude Code bắt đầu thì MCP đó **mất kết nối cả
phiên** và không nối lại được. `unity_mcp.py` gọi thẳng JSON-RPC (streamable HTTP) nên
dùng được bất cứ lúc nào Unity đang mở.

## Dùng

```powershell
# Unity phải đang mở project D:\NRO_Project (cổng 8080 đang nghe).
python tools\unity-mcp\unity_mcp.py --list                       # liệt kê tool
python tools\unity-mcp\unity_mcp.py --schema manage_camera        # xem tham số một tool
python tools\unity-mcp\unity_mcp.py --code tools\unity-mcp\cs\trangthai.cs   # chạy C# trong Editor
python tools\unity-mcp\unity_mcp.py --json tools\unity-mcp\ed-play.json      # gọi tool bất kỳ
```

`--code <file.cs>`: thân hàm C# chạy trong Editor (có `UnityEngine`, `UnityEditor`,
truy cập được `Game1.*` lúc Play), `return` chuỗi để đọc kết quả. Tránh quoting
PowerShell bằng cách để code trong file.

## Kịch bản C# sẵn (`cs/`)

| File | Làm gì |
|---|---|
| `trangthai.cs` | Đang Play? scene nào? cỡ màn hình, Editor có focus |
| `ketnoi.cs` | Nối máy chủ chưa, màn hình hiện tại, nhân vật của mình (hp, đầu/thân/chân) |
| `dangnhap.cs` | Ghi `acc`/`pass` vào RMS rồi `LoginScr.doLogin()` (tài khoản `naruto/1`) |
| `nhanvat.cs` | Nhân vật + part đang dùng (`dy` từng mảnh — biết client đã lấy `part.bin` mới chưa) |
| `vitri.cs` | Vị trí quái (`y`, `yFirst`), NPC, người khác trong map |
| `npc.cs` | `arrNpcTemplate`, `taskMaint`, `vNpc`, waypoint, bóng đổ — chẩn đoán "NPC không hiện" |
| `npc-mo.cs` / `npc-menu.cs` / `panel.cs` | Bấm NPC 0 → chọn "Cửa hàng" → đọc panel |
| `thu-npc.cs` | Tạo tay `Mob`/`Npc` để bắt ngoại lệ mà client nuốt |
| `chup.cs` | `ScreenCapture` Game view ra PNG (đọc file ở lượt sau, ảnh xuất hiện sau ~1 khung) |
| `iruka-part.cs` | Part 30/57/58/2189/2190/2191 client đang nạp (id, dx, dy 3 khung đầu) + NPC trong map (template head/body/leg, cf) — kiểm client đã nhận `part.bin` mới (06/09) |
| `nhay.cs` | Giữ phím nhảy (`keyHold[21|2]`) — kẹt vách phải nhảy rồi mới giữ trái/phải (chuỗi else-if) |
| `co-cho.cs` | Các cờ làm client vẽ "Xin chờ" (`isLoading`, `isLoadingMap`, `msgdlg`, `timeBreakLoading`, hp/statusMe) — chẩn đoán treo/rớt |
| `chay-phai.cs` / `chay-trai.cs` / `tha-phim.cs` | Giữ phím chạy phải/trái (`GameCanvas.keyHold[24|6]` / `[23|4]`, GameScr.cs:2737–2845) rồi thả — để chụp vòng chạy cf 2→6 và hướng nhìn trái (client lật gương) |

Chụp hoạt ảnh nhiều khung: `python tools\unity-mcp\chup-loat.py --out <thư mục> --ten chay --n 8 --cach 0.12 --truoc tools\unity-mcp\cs\chay-phai.cs --sau tools\unity-mcp\cs\tha-phim.cs` (giữ phím → chụp N ảnh cách 120 ms → thả phím).

Lệnh Editor (`--json`): `{"name":"manage_editor","arguments":{"action":"stop"}}` /
`"play"`; `{"name":"refresh_unity","arguments":{}}` sau khi chép asset vào
`Assets/Resources`; `{"name":"read_console","arguments":{"action":"get","count":"40","types":["error"]}}`.

## Bẫy đã vấp

- **Unity không có focus thì Game view không vẽ lại** → `ScreenCapture` trả khung cũ (thấy màn "Xin chờ" thay vì map). `chup.cs`/`chup-loat.py` gọi `InternalEditorUtility.RepaintAllViews()` trước khi chụp.
- **Bấm Play ngay sau `refresh_unity`** (đang nhập lại nhiều ảnh) → phiên chơi bị reset về màn chọn máy chủ. Chờ nhập xong.
- **Boss sinh 10 phút sau khi máy chủ khởi động** ngay điểm sinh của tài khoản thử → chết trong ~9 s, rồi client rớt. Chụp/đo trong 10 phút đầu sau khi khởi động lại máy chủ.
- `execute_code` chặn `System.IO.File.Delete` (safety) — `ScreenCapture` ghi đè được, không cần xoá.
- `types` của `read_console` phải là **mảng** (`["error"]`), không phải chuỗi.
- Client nuốt ngoại lệ trong `try/catch` không log gì → Console 0 lỗi mà vẫn hỏng. Khi đó
  **giải mã gói thật** (`tools/goi-tin/giai-ma-24.py` + log Debug hex của máy chủ) chứ đừng đoán.
- Đang Play mà chép ảnh vào `Assets/Resources` là Unity reimport giữa chừng: **Stop → chép →
  refresh → Play**.
- Tài khoản thử `naruto` đứng ở điểm sinh boss; đứng yên là bị boss đánh chết — không phải lỗi.
