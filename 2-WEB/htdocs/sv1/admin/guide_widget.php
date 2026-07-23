<?php
$guidePage = basename($_SERVER['SCRIPT_NAME'] ?? '', '.php');
$guideData = [
    'account' => ['Account', '* Những Account là Admin (cột Admin = 1) thì sẽ có quyền chỉnh sửa Panel Web, hãy bảo mật những Account này.<br><br>* Cột Thành Viên = 1 tức là Account Đã Mở Thành Viên.<br><br>* Cột Khóa = 1 tức là Account Đã Bị Khóa.'],
    'buff_item' => ['Buff Item', '* Cộng thẳng cho nhân vật.<br><br>* Nhân vật cần Online.'],
    'player' => ['Nhân Vật', '* Khi Ấn Sửa = Nhân Vật trong game tự bị Kick, phải Ấn Lưu hoặc Hủy thì mới vào lại được game.<br><br>* Có thể Xóa trang bị, vật phẩm, sửa sức mạnh, ... tùy thích.'],
    'giftcode' => ['Giftcode', '* Ấn Tạo Code nhanh để tạo code cực nhanh và dễ.<br><br>* Code vừa tạo có thể sử dụng ngay.'],
    'phucloi' => ['Phúc Lợi', '* Khi bạn sửa quà ở đây thì sẽ cập nhật trực tiếp trong game mà không cần Bảo trì hoặc Chạy lại Server.<br><br>* Tuy nhiên, sau Bảo trì hoặc Chạy lại game thì quà sẽ bị Reset. Nếu bật Chế Độ Panel thì sẽ không bị Reset.'],
    'item' => ['Vật Phẩm', '* Tìm vật phẩm bằng Check Nhanh hoặc chọn Ô 1, 2, 3, 4.<br><br>* Up_To = 1 tức là Vật phẩm có thể cộng dồn x9999.<br><br>* Sửa Tên + Miêu Tả + ICON cần Bảo trì hoặc Chạy lại Server và Xóa dữ liệu đi Tải lại.'],
    'itemshop' => ['Cửa Hàng', '* Ấn Update Shop để cập nhật dữ liệu vừa chỉnh sửa mà không cần Bảo trì hoặc Chạy lại game.<br><br>* Bán bằng Hồng ngọc thì ghi 7743.'],
    'gioihan' => ['Chỉ Số', '* Cần Bảo trì và Chạy lại Server để Update.<br><br>* Không nên sửa Giới hạn Sức mạnh, chỉ sửa các chỉ số thôi.'],
    'map' => ['Bản Đồ', '* Cần Bảo trì và Chạy lại Server để Update.'],
    'giaodich' => ['Giao Dịch', '* Hãy Check tên nhân vật rồi Dò IP của nó trước để khóa cả IP thay vì chỉ khóa 1 Account.'],
    'noitai' => ['Nội Tại', '* Cần Bảo trì và Chạy lại Server để Update.'],
];
if (!isset($guideData[$guidePage])) return;
[$guideTitle, $guideText] = $guideData[$guidePage];
?>
<style>
.page-guide-btn{display:inline-flex!important;align-items:center!important;justify-content:center!important;gap:6px!important;padding:11px 20px!important;border:2px solid #fbbf24!important;border-radius:8px!important;background:linear-gradient(135deg,#8b5cf6,#6366f1)!important;color:#fff!important;box-shadow:0 3px 12px rgba(99,102,241,.35),inset 0 0 0 1px rgba(255,255,255,.2)!important;font:700 .9rem 'Plus Jakarta Sans',Arial,sans-serif!important;line-height:1!important;white-space:nowrap!important;cursor:pointer!important}
.page-guide-btn:hover{background:linear-gradient(135deg,#7c3aed,#4f46e5)!important;border-color:#fcd34d!important;transform:translateY(-2px)}
#pageGuideModal{display:none;position:fixed;z-index:200000;left:1vw;top:0;width:1500px;background:rgba(9,19,31,.65);align-items:center;justify-content:center;font-family:'Plus Jakarta Sans',Arial,sans-serif}
.page-guide-box{width:560px;max-width:90%;padding:32px;background:#fff;border:2px solid #fbbf24;border-radius:13px;box-shadow:0 12px 40px rgba(0,0,0,.35);color:#1a3a5c;text-align:left}
.page-guide-box h3{margin:0 0 16px;color:#6d3fd1;font-size:1.4rem;text-align:center}
.page-guide-content{font-size:1.05rem;line-height:1.75}
.page-guide-close{width:100%;margin-top:20px;padding:11px;border:2px solid #fbbf24;border-radius:8px;background:linear-gradient(135deg,#8b5cf6,#6366f1);color:#fff;font-weight:800;cursor:pointer}
</style>
<div id="pageGuideModal">
  <div class="page-guide-box" onclick="event.stopPropagation()">
    <h3>Hướng Dẫn <?= htmlspecialchars($guideTitle, ENT_QUOTES, 'UTF-8') ?></h3>
    <div class="page-guide-content"><?= $guideText ?></div>
    <button type="button" class="page-guide-close" onclick="document.getElementById('pageGuideModal').style.display='none'">Đóng</button>
  </div>
</div>
<script>
(function(){
  var modal=document.getElementById('pageGuideModal');
  var refresh=[].slice.call(document.querySelectorAll('button,a')).find(function(el){return /(làm\s*mới|f5)/i.test(el.textContent||'')});
  if(refresh){
    var btn=document.createElement('button'); btn.type='button'; btn.className='page-guide-btn'; btn.innerHTML='❓ <span>Hướng Dẫn</span>';
    btn.onclick=function(){modal.style.display='flex'};
    refresh.insertAdjacentElement('afterend',btn);
  }
  function scaleGuide(){
    var ratio=(window.innerWidth*.98)/1500;
    modal.style.width='1500px'; modal.style.height=(window.innerHeight/ratio)+'px';
    modal.style.transform='scale('+ratio+')'; modal.style.transformOrigin='top left'; modal.style.left='1vw'; modal.style.top='0';
  }
  scaleGuide(); window.addEventListener('resize',scaleGuide);
})();
</script>
