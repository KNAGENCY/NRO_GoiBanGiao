# Công cụ dựng map Làng Lá

Chạy bằng `python` (Windows). Đích tác động: `D:\NRO_Project\Assets\Resources\res\`.

| File | Làm gì | Khi nào chạy |
|---|---|---|
| `unity.py` | Nối thẳng HTTP tới MCP Unity (bỏ cầu `npx mcp-remote` hay chết) | thư viện, các file khác `import` |
| `vao-game-chup.py <tên>` | Dừng → nạp lại asset → sửa cài đặt import → chạy game → vào map → chụp 2 ảnh `<tên>A/B.png` | mỗi lần muốn xem kết quả thật |
| `ve-bo-cuc-map.py` | Dọn tile cũ + vẽ lại toàn bộ bố cục map 1 | khi đổi bố cục |
| `do-vien-tile.py` | Dò viền trắng/đen nằm trong ảnh tile | khi thấy ô có khung lạ |
| `xoa-vien-tile.py` | Xoá viền đó (kéo pixel sạch đè ra mép), có sao lưu `.bak-vien` | sau khi `do-vien-tile.py` báo bẩn |

Lỗi đã vấp và cách sửa: [`../../docs/my-thuat/map-1-lang-la/LOI-DA-GAP.md`](../../docs/my-thuat/map-1-lang-la/LOI-DA-GAP.md)
