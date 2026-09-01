# KH6b Task 1 (GATE task 11) - Client co chap nhan id bi -6 roi -5 quay lai khong?

TRA LOI: CO - client chap nhan mot id da bi goi -6 xoa xuat hien tro lai qua goi -5 voi cung id, khong co cache hay bang tra nao giu vet id cu can tro viec nay.

## 1. Xu ly goi -6 (roi khu vuc) - Controller.cs:4252-4272

Doan xu ly (Controller.cs, case -6, bat dau dong 4252):
- Doc charID tu message (dong 4255).
- Duyet vCharInMap tim char co charID trung (dong 4256-4259).
- Neu tim thay: them hieu ung bien mat (dong 4263), va neu KHONG co isUsePlane thi XOA phan tu khoi vCharInMap bang removeElementAt (dong 4265-4268).

Diem can luu y: dieu kien isUsePlane (dong 4265) - neu boss dang bat co isUsePlane = true thi goi -6 se KHONG xoa no khoi vCharInMap. Grep toan bo Char.cs cho thay isUsePlane chi duoc khai bao o Char.cs dong 25 va doc trong dieu kien isTeleport hoac isUsePlane (Char.cs dong 1678) va freeze check (Char.cs dong 4913) - khong thay logic nao gan rieng co nay cho boss, nen voi boss thong thuong phan tu se bi xoa that.

## 2. Xu ly goi -5 (xuat hien) - Controller.cs:4110 tro di

Doan xu ly (case -5, dong 4110-4157):
- Doc charID2 (dong 4112).
- Tao Char hoac Mabu MOI hoan toan, khong tai su dung object cu theo id (dong 4114-4122).
- Goi readCharInfo(char15, msg) tai dong 4127 (dinh nghia o Controller.cs dong 6368) - ghi toan bo thuoc tinh ten, cap, ngoai hinh, dau toc truc tiep vao object moi, nen bien hinh doi ten/cap/ngoai hinh la hoan toan hop le ma khong dinh du lieu cu.
- Chi them vao vCharInMap NEU chua ton tai id do: if (GameScr.findCharInMap(char15.charID) == null) GameScr.vCharInMap.addElement(char15) (dong 4154-4156).
- Vi -6 da xoa phan tu cu khoi vCharInMap truoc do, findCharInMap se tra null va object moi duoc them binh thuong.

## 3. findCharInMap - khong co cache rieng

GameScr.cs dong 5285-5294, findCharInMap chi la vong lap tuyen tinh tren vCharInMap hien tai, khong dung Hashtable hay Dictionary nao khac. Da grep toan bo GameScr.cs cho tu khoa Hashtable va Dictionary - khong co ket qua nao. Khong tim thay cache anh/ten/hieu ung nao giu theo charID ma khong bi don khi phan tu bi xoa khoi vCharInMap.

## 4. Bay can luu y

- Dieu kien isUsePlane (Controller.cs dong 4265): neu vi ly do nao do boss duoc danh dau isUsePlane = true, goi -6 se KHONG xoa no khoi vCharInMap. Khi do goi -5 tiep theo voi cung id se thay findCharInMap khac null (dong 4154) va KHONG add lai doi tuong moi - object cu (pha 1) van nam trong vCharInMap, con object moi (pha 2, da doc xong thong tin) bi rot, khong hien thi. Day la kich ban loi duy nhat tim thay trong code, nhung chi xay ra khi co isUsePlane bat - chua xac minh boss binh thuong co bao gio bat co nay hay khong.
- Khong thay yeu cau cho N ms giua -6 va -5 trong code client - case -6 xu ly dong bo ngay khi doc goi, goi -5 tiep theo duoc xu ly o luot doc message ke tiep nhu moi goi khac. Khong co debounce/delay rieng cho cap -6/-5.
- Khong can gui them goi phu nao de don client - viec xoa dien ra ngay trong xu ly -6.

## 5. Doi chieu ban Java goc

- Boss.java dong 664-676 (leaveMap()): khi boss con pha ke tiep (currentLevel < data.length - 1), server KHONG goi ChangeMapService.exitMap (khong roi han zone) ma nhay thang sang trang thai RESPAWN (dong 668) roi respawn() (dong 366-370: currentLevel++, initBase()) roi JOIN_MAP -> joinMap() (dong 380).
- joinMap() (Boss.java dong 380-397) goi joinMapByZone(zoneFinal) -> ChangeMapService.gI().changeMap(this, this.zone, x, y) (Boss.java dong 473) - day la CUNG MOT object Boss, cung id, chi doi zone/toa do va du lieu (data[nextLevel] qua initBase()), roi phat lai thong tin xuat hien.
- Noi cach khac: ban goc Java KHONG dung cach gui cap -6 roi -5 cho cung id de doi ngoai hinh ma khong doi zone - no dung co che changeMap (doi zone) dung chung co che broadcast roi/vao zone. Cach tiep can gui -6 xong gui lai -5 voi cung id de doi ngoai hinh la de xuat MOI cua server Go, chua co tien le giong het trong code Java - nhung phan tich client o tren cho thay client DU KHA NANG xu ly dung cach nay (mien la boss khong bi bat co isUsePlane).
- Ket luan: cach tiep can -6 roi -5 cung id dung duoc cho server Go, vi phia client khong cache gi ngan can, du no khong phai la cach ban Java goc tung lam.

## 6. Chua xac minh duoc (chi suy luan tu code, chua chay thu)

- Chua chay thuc te game de xac nhan hanh vi hien thi (khong co moi truong build/run client Unity trong phien lam viec nay) - toan bo ket luan dua tren doc tinh code C# va Java.
- Chua xac minh boss trong thu muc boss_manifest co bao gio dat isUsePlane = true hay khong (khong grep het toan bo gan 90 file boss manifest); neu co boss dac biet bat co nay thi kich ban loi o muc 4 co the xay ra.
- Chua kiem tra cac he thong phu khac co the tham chieu charID cu sau khi bien mat, vi du: thanh mau boss tren UI, minimap, hoac hieu ung dang bay toi boss cu (ServerEffect, Teleport.addTeleport dong 4137) - cac he thong nay khong duoc ra soat chi tiet trong lan doc nay.
- Chua xac minh thu tu doc goi phia client khi hai goi -6 va -5 den trong cung mot buffer/frame mang co gay race hay khong - code xu ly tuan tu theo message queue nhung chua kiem chung bang test thuc te.

## Kết luận sau khi đối chiếu phía máy chủ Go (02/09/2026)

Bẫy duy nhất ở trên — `isUsePlane == true` làm `-6` không xoá — **không chạm tới boss**:

- Client chỉ bật `isUsePlane` ở **một** chỗ: `Controller.cs:1339`, nằm trong `case -65` (gói dịch chuyển/teleport),
  **không phải** trong payload của gói `-5`. Nghĩa là gói `-5` xuất hiện boss **không thể** bật cờ này.
- Máy chủ Go **không bao giờ gửi gói `-65`**: grep toàn bộ `server/internal` (trừ test) thấy `-65` xuất hiện **0 lần**.

→ Với boss do server Go sinh, cờ này luôn `false`, nên `-6` rồi `-5` cùng id **chạy đúng**. **GATE mở — task 11 được làm.**

Còn phải giữ trong đầu khi thi công task 11: gửi `-6` **trước** `-5` trong cùng một lượt actor (cùng thứ tự FIFO
qua `XepGui`), và **giữ nguyên id**. Chưa kiểm chạy thật trong client — đó là việc nghiệm thu tay của chủ dự án.
