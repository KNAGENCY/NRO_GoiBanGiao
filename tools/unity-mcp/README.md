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

Lệnh Editor (`--json`): `{"name":"manage_editor","arguments":{"action":"stop"}}` /
`"play"`; `{"name":"refresh_unity","arguments":{}}` sau khi chép asset vào
`Assets/Resources`; `{"name":"read_console","arguments":{"action":"get","count":"40","types":["error"]}}`.

## Bẫy đã vấp

- `execute_code` chặn `System.IO.File.Delete` (safety) — `ScreenCapture` ghi đè được, không cần xoá.
- `types` của `read_console` phải là **mảng** (`["error"]`), không phải chuỗi.
- Client nuốt ngoại lệ trong `try/catch` không log gì → Console 0 lỗi mà vẫn hỏng. Khi đó
  **giải mã gói thật** (`tools/goi-tin/giai-ma-24.py` + log Debug hex của máy chủ) chứ đừng đoán.
- Đang Play mà chép ảnh vào `Assets/Resources` là Unity reimport giữa chừng: **Stop → chép →
  refresh → Play**.
- Tài khoản thử `naruto` đứng ở điểm sinh boss; đứng yên là bị boss đánh chết — không phải lỗi.
