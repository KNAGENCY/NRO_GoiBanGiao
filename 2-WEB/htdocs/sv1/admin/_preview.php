<?php
// TRANG XEM TRUOC TAM - dung dung 3 file CSS that, khong can dang nhap.
// Chi de duyet giao dien. XOA sau khi duyet xong.
require_once __DIR__ . '/admin_layout.php';
admin_head('Xem trước giao diện');
admin_bar('Danh sách người chơi', 'users');
?>
<div class="kadmin-body">

  <div class="kadmin-stats">
    <div class="kadmin-stat"><p class="kadmin-stat__label">Đang online</p><div class="kadmin-stat__value">137</div></div>
    <div class="kadmin-stat"><p class="kadmin-stat__label">Tổng tài khoản</p><div class="kadmin-stat__value">2.480</div></div>
    <div class="kadmin-stat"><p class="kadmin-stat__label">Nạp hôm nay</p><div class="kadmin-stat__value">1,2 Tr</div></div>
    <div class="kadmin-stat"><p class="kadmin-stat__label">Nhân vật mới</p><div class="kadmin-stat__value">54</div></div>
  </div>

  <div class="kadmin-card">
    <div class="kadmin-card__head"><i class="ti ti-list" aria-hidden="true"></i> Người chơi</div>
    <div class="kadmin-toolbar">
      <input class="kame-input kadmin-toolbar__grow" placeholder="Tìm tên hoặc ID...">
      <button class="kame-button"><i class="ti ti-search" aria-hidden="true"></i> Tìm</button>
      <button class="kame-button kame-button--secondary"><i class="ti ti-gift" aria-hidden="true"></i> Tặng vật phẩm</button>
    </div>
    <div class="kadmin-table-wrap">
      <table class="kadmin-table">
        <thead><tr><th>ID</th><th>Tên</th><th>Sức mạnh</th><th>Vàng</th><th>Ngọc</th><th></th></tr></thead>
        <tbody>
          <tr><td class="num">1</td><td>cuhan123</td><td class="num">51.510</td><td class="num">9,5 Tỷ</td><td class="num">5 Tr</td><td class="act"><button class="kame-button kame-button--ghost"><i class="ti ti-edit" aria-hidden="true"></i></button></td></tr>
          <tr><td class="num">2</td><td>songoku99</td><td class="num">128.400</td><td class="num">2,1 Tỷ</td><td class="num">120 N</td><td class="act"><button class="kame-button kame-button--ghost"><i class="ti ti-edit" aria-hidden="true"></i></button></td></tr>
          <tr><td class="num">3</td><td>vegeta_sm</td><td class="num">96.220</td><td class="num">870 Tr</td><td class="num">44 N</td><td class="act"><button class="kame-button kame-button--ghost"><i class="ti ti-edit" aria-hidden="true"></i></button></td></tr>
          <tr><td class="num">4</td><td>picolo_x</td><td class="num">33.900</td><td class="num">210 Tr</td><td class="num">8 N</td><td class="act"><button class="kame-button kame-button--ghost"><i class="ti ti-edit" aria-hidden="true"></i></button></td></tr>
        </tbody>
      </table>
    </div>
  </div>

  <div class="kadmin-card" style="max-width:520px">
    <div class="kadmin-card__head"><i class="ti ti-alert-triangle" aria-hidden="true"></i> Bật chế độ bảo trì</div>
    <p class="kame-muted" style="margin:0 0 var(--kame-space-4)">Nhấn nút dưới để bắt đầu chế độ bảo trì. Tất cả người chơi sẽ bị kick.</p>
    <div class="kame-alert kame-alert--success" style="margin-bottom:var(--kame-space-4)">Lệnh bảo trì đã được gửi thành công.</div>
    <button class="kame-button kame-button--danger" style="width:100%"><i class="ti ti-player-power" aria-hidden="true"></i> Bắt Đầu Bảo Trì</button>
  </div>

</div>
<?php
admin_foot();
