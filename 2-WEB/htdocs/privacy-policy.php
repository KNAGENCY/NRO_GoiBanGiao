<?php
/* ==========================================================================
   NGỌC RỒNG KAME — CHÍNH SÁCH BẢO MẬT

   Đường dẫn /privacy-policy PHẢI giữ nguyên: kho ứng dụng Google Play và
   TestFlight trỏ thẳng tới URL này, đổi đường dẫn là hồ sơ ứng dụng gãy link.

   File nằm ở gốc web nên đường dẫn tới bộ khung là 'sv1/kame/layout.php'.
   Trang song sinh với trang này là /sv1/dieu-khoan (điều khoản sử dụng);
   hai trang trỏ qua lại lẫn nhau ở cuối.
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/sv1/kame/layout.php';

/* Ngày cập nhật hiển thị ở đầu trang. Sửa ở đây là đổi cả trang. */
const PP_NGAY_CAP_NHAT = '21/07/2026';

/* Địa chỉ liên hệ. Hòm thư này do ban quản trị tạo, đổi ở một chỗ duy nhất. */
const PP_EMAIL = 'support@nrokame.com';

/* --- Danh sách mục ------------------------------------------------------
   Khoá mảng chính là id của thẻ <section>, dùng cho cả mục lục lẫn tiêu đề
   trong thân trang. Nhờ vậy mục lục không bao giờ lệch với nội dung.
   ------------------------------------------------------------------------ */
$ppMuc = [
    'gioi-thieu'  => 'Trang này nói về chuyện gì',
    'thu-thap'    => 'Chúng tôi giữ những thông tin nào của bạn',
    'muc-dich'    => 'Chúng tôi dùng thông tin đó để làm gì',
    'luu-tru'     => 'Thông tin được giữ trong bao lâu',
    'chia-se'     => 'Ai được xem thông tin của bạn',
    'quyen'       => 'Bạn có quyền gì với thông tin của mình',
    'mat-khau'    => 'Mật khẩu của bạn được giữ thế nào',
    'cookie'      => 'Cookie trên trang web',
    'tre-em'      => 'Người chơi nhỏ tuổi',
    'thay-doi'    => 'Khi chính sách này thay đổi',
    'lien-he'     => 'Liên hệ với chúng tôi',
];

/**
 * Mở một mục: in thẻ <section> kèm tiêu đề lấy từ $ppMuc.
 * Số thứ tự do CSS đếm, không gõ tay để khỏi sai khi thêm bớt mục.
 */
function pp_mo_muc(string $id): void
{
    global $ppMuc;
    $tieuDe = $ppMuc[$id] ?? $id;
    echo '<section class="kpl-muc" id="' . kame_e($id) . '">';
    echo '<h2 class="kame-heading kpl-muc__tieude">' . kame_e($tieuDe) . '</h2>';
}

/** Đóng mục vừa mở. */
function pp_dong_muc(): void
{
    echo '</section>';
}

kame_head(
    'Chính sách bảo mật',
    'Ngọc Rồng Kame giữ những thông tin nào của người chơi, dùng để làm gì, giữ bao lâu và bạn có quyền gì với thông tin của mình.'
);
kame_header('');
?>

<style>
/* --- CSS riêng của hai trang pháp lý. Tiền tố kpl- (Kame Pháp Lý). ------
   Không sửa ui-kit.css; chỉ dùng lại biến màu và khoảng cách có sẵn.      */

/* Dải đầu trang: nền đỏ mận trơn, không cần ảnh nên không sợ thiếu asset */
.kpl-band {
  border-bottom: var(--kame-border);
  background: var(--kame-maroon-900);
}
.kpl-band__inner {
  width: min(var(--kame-container), calc(100% - 32px));
  margin-inline: auto;
  padding-block: var(--kame-space-7) var(--kame-space-6);
}
.kpl-band__title { color: var(--kame-white); text-shadow: 0 4px 0 rgb(0 0 0 / 35%); }
.kpl-band .kame-breadcrumb { margin-top: var(--kame-space-3); color: rgb(255 253 245 / 82%); }
.kpl-band .kame-breadcrumb a { color: var(--kame-gold); }
.kpl-band__ngay {
  display: inline-block;
  margin-top: var(--kame-space-4);
  padding: 6px 14px;
  border: 2px solid rgb(255 200 61 / 55%);
  border-radius: 999px;
  color: var(--kame-gold);
  font-size: 14px;
  font-weight: 800;
}

/* Mục lục */
.kpl-mucluc ol {
  display: grid;
  gap: 6px;
  margin: var(--kame-space-4) 0 0;
  padding-left: 22px;
  font-size: 15px;
  line-height: 1.7;
}
.kpl-mucluc li::marker { color: var(--kame-orange-600); font-weight: 800; }

/* Thân văn bản: bộ đếm cho số thứ tự các mục */
.kpl-than { counter-reset: kplmuc; }
.kpl-muc {
  margin-top: var(--kame-space-7);
  padding-top: var(--kame-space-5);
  border-top: 2px dashed rgb(96 0 24 / 24%);
  scroll-margin-top: 90px;              /* để tiêu đề không bị header che */
}
.kpl-muc:first-child { margin-top: 0; padding-top: 0; border-top: 0; }
.kpl-muc__tieude {
  display: flex;
  align-items: baseline;
  gap: var(--kame-space-3);
  font-size: clamp(20px, 3vw, 26px);
}
.kpl-muc__tieude::before {
  counter-increment: kplmuc;
  content: counter(kplmuc);
  color: var(--kame-orange-600);
  font-size: .82em;
  opacity: .75;
}
.kpl-muc p { margin: var(--kame-space-3) 0 0; font-size: 16px; line-height: 1.8; }
.kpl-muc ul {
  margin: var(--kame-space-3) 0 0;
  padding-left: 22px;
  font-size: 16px;
  line-height: 1.8;
}
.kpl-muc li { margin-bottom: 4px; }
.kpl-muc li::marker { color: var(--kame-orange-600); }
.kpl-muc h3 {
  margin: var(--kame-space-5) 0 0;
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 18px;
  font-weight: 800;
  color: var(--kame-maroon-700);
}

/* Khối liên kết sang trang pháp lý còn lại */
.kpl-cheo {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: var(--kame-space-5);
  padding: var(--kame-space-6);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 94%);
  box-shadow: var(--kame-shadow-sm);
}
.kpl-cheo p { margin: var(--kame-space-2) 0 0; max-width: 62ch; font-size: 15px; }

@media (max-width: 600px) {
  .kpl-band__inner { padding-block: var(--kame-space-6) var(--kame-space-5); }
  .kpl-cheo { padding: var(--kame-space-5); }
  .kpl-cheo .kame-button { width: 100%; }
}
</style>

<main class="kame-page-wrap">

  <!-- ===== Dải đầu trang ================================================ -->
  <section class="kpl-band">
    <div class="kpl-band__inner">
      <h1 class="kame-heading kame-heading--lg kpl-band__title">Chính sách bảo mật</h1>
      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <span>Chính sách bảo mật</span>
      </nav>
      <div class="kpl-band__ngay">Cập nhật ngày <?= kame_e(PP_NGAY_CAP_NHAT) ?></div>
    </div>
  </section>

  <section class="kame-section">
    <div class="kame-container">
      <div class="kame-grid">

        <!-- ===== Mục lục =================================================== -->
        <div class="kame-col-4">
          <article class="kame-card kpl-mucluc">
            <div class="kame-card__body">
              <h2 class="kame-heading kame-heading--md" style="font-size:22px">Trong trang này có gì</h2>
              <ol>
                <?php foreach ($ppMuc as $id => $tieuDe): ?>
                  <li><a class="kame-link" href="#<?= kame_e($id) ?>"><?= kame_e($tieuDe) ?></a></li>
                <?php endforeach; ?>
              </ol>
            </div>
          </article>

          <div class="kame-alert kame-alert--info" style="margin-top:var(--kame-space-4)">
            <span aria-hidden="true">!</span>
            <span>
              Đây là <strong>máy chủ do cá nhân vận hành</strong>, không thuộc và không liên kết
              với chủ sở hữu bản quyền Dragon Ball hay bất kỳ nhà phát hành nào.
            </span>
            <span></span>
          </div>
        </div>

        <!-- ===== Nội dung ================================================= -->
        <div class="kame-col-8">
          <div class="kpl-than">

            <?php pp_mo_muc('gioi-thieu'); ?>
              <p>
                Ngọc Rồng Kame là một máy chủ trò chơi do cá nhân vận hành. Để bạn đăng nhập được,
                giữ được nhân vật và chơi cùng người khác, chúng tôi buộc phải lưu lại một ít thông tin
                về bạn. Trang này viết rõ ra: <strong>giữ những gì, để làm gì, giữ bao lâu, ai xem được</strong>
                và bạn có thể yêu cầu chúng tôi làm gì với những thông tin đó.
              </p>
              <p>
                Chúng tôi cố gắng viết ngắn và dễ hiểu thay vì dùng câu chữ luật lệ. Khi bạn tạo tài khoản
                và chơi, xem như bạn đã đọc và đồng ý với những gì ghi ở đây.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('thu-thap'); ?>
              <p>Danh sách dưới đây là đầy đủ. Ngoài những thứ này chúng tôi không thu thập gì thêm.</p>

              <h3>Khi bạn đăng ký tài khoản</h3>
              <ul>
                <li><strong>Tên đăng nhập</strong> — bắt buộc, vì đó là cách hệ thống nhận ra bạn.</li>
                <li><strong>Mật khẩu</strong> — bắt buộc, nhưng được lưu ở dạng đã mã hoá một chiều
                    (xem mục "Mật khẩu của bạn được giữ thế nào").</li>
                <li><strong>Số điện thoại</strong> — không bắt buộc, dùng để lấy lại mật khẩu khi bạn quên.</li>
                <li><strong>Email</strong> — không bắt buộc, cũng chỉ dùng để hỗ trợ và lấy lại tài khoản.</li>
              </ul>

              <h3>Khi bạn đăng nhập và chơi</h3>
              <ul>
                <li><strong>Địa chỉ IP</strong> của máy bạn dùng để kết nối.</li>
                <li><strong>Thời điểm đăng nhập và đăng xuất</strong>.</li>
                <li><strong>Dữ liệu nhân vật trong game</strong>: tên nhân vật, hành tinh, cấp độ, sức mạnh,
                    vật phẩm, tiền trong game, tiến độ nhiệm vụ và các chỉ số tương tự.</li>
              </ul>

              <h3>Khi bạn nhắn cho ban quản trị</h3>
              <ul>
                <li>Nội dung tin nhắn và những gì bạn tự gửi kèm (ảnh chụp màn hình, mô tả lỗi...).</li>
              </ul>

              <p>
                Chúng tôi <strong>không</strong> lấy danh bạ, ảnh, vị trí GPS, micrô hay camera của bạn.
                Chúng tôi cũng không yêu cầu ảnh giấy tờ tuỳ thân.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('muc-dich'); ?>
              <p>Mỗi thứ thu thập đều gắn với một việc cụ thể:</p>
              <ul>
                <li><strong>Cho bạn đăng nhập và giữ nhân vật</strong> — nếu không lưu thì mỗi lần thoát game
                    là mất hết.</li>
                <li><strong>Lấy lại tài khoản khi bạn quên mật khẩu</strong> — đây là lý do duy nhất chúng tôi
                    hỏi số điện thoại và email.</li>
                <li><strong>Giữ trật tự trong game</strong> — địa chỉ IP và lịch sử đăng nhập giúp chúng tôi
                    xử lý các trường hợp gian lận, phá rối hoặc cố đăng nhập vào tài khoản người khác.</li>
                <li><strong>Sửa lỗi và cân bằng game</strong> — xem chỗ nào người chơi hay kẹt, chỗ nào quá dễ
                    hoặc quá khó.</li>
                <li><strong>Trả lời khi bạn cần hỗ trợ</strong>.</li>
              </ul>
              <p>
                Chúng tôi <strong>không dùng thông tin của bạn để quảng cáo</strong>, không gắn mã theo dõi của
                bên thứ ba, và <strong>không bán dữ liệu cho bất kỳ ai</strong>.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('luu-tru'); ?>
              <ul>
                <li><strong>Thông tin tài khoản và nhân vật</strong> được giữ chừng nào tài khoản của bạn
                    còn tồn tại.</li>
                <li><strong>Địa chỉ IP và lịch sử đăng nhập</strong> chỉ giữ trong thời gian cần thiết để
                    xử lý gian lận và sự cố, sau đó được xoá dần.</li>
                <li><strong>Khi bạn yêu cầu xoá tài khoản</strong>, chúng tôi xoá toàn bộ dữ liệu nhân vật và
                    thông tin cá nhân gắn với tài khoản đó.</li>
              </ul>
              <p>
                Có một ngoại lệ: nếu một tài khoản bị khoá vì gian lận hoặc phá hoại, chúng tôi có thể giữ lại
                dấu vết tối thiểu (ví dụ tên tài khoản đã bị khoá) để người đó không lập lại tài khoản mới rồi
                tiếp tục làm điều tương tự.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('chia-se'); ?>
              <p>
                Nói ngắn gọn: <strong>không chia sẻ với ai cả</strong>. Cụ thể hơn:
              </p>
              <ul>
                <li>Chúng tôi <strong>không bán, không trao đổi, không cho thuê</strong> dữ liệu người chơi.</li>
                <li>Không có công ty quảng cáo hay đơn vị phân tích bên thứ ba nào được nối vào hệ thống này.</li>
                <li>Chỉ những người trong ban quản trị máy chủ mới xem được dữ liệu, và chỉ khi cần để vận hành
                    hoặc xử lý sự cố.</li>
                <li>Nhà cung cấp máy chủ (nơi đặt hệ thống) về mặt kỹ thuật có quyền truy cập vào máy chủ,
                    giống như mọi dịch vụ lưu trữ khác.</li>
                <li>Nếu cơ quan nhà nước có thẩm quyền yêu cầu theo đúng quy định pháp luật, chúng tôi buộc
                    phải cung cấp thông tin trong phạm vi được yêu cầu.</li>
              </ul>
              <p>
                Một điều bạn nên biết: <strong>tên nhân vật và các chỉ số công khai trong game</strong>
                (cấp độ, sức mạnh, thứ hạng) thì người chơi khác nhìn thấy được — đó là chuyện bình thường
                của một trò chơi nhiều người. Tên đăng nhập, mật khẩu, số điện thoại và email của bạn thì không.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('quyen'); ?>
              <p>Với thông tin của chính mình, bạn có thể yêu cầu:</p>
              <ul>
                <li><strong>Xem lại</strong> — hỏi chúng tôi đang giữ những gì về tài khoản của bạn.</li>
                <li><strong>Sửa lại</strong> — đổi mật khẩu và số điện thoại bạn tự làm được ngay trên trang web.
                    Những thứ khác thì nhắn cho ban quản trị.</li>
                <li><strong>Xoá tài khoản</strong> — chúng tôi sẽ xoá tài khoản cùng dữ liệu nhân vật của bạn.</li>
                <li><strong>Rút lại sự đồng ý</strong> — bạn có thể ngừng chơi và yêu cầu xoá bất cứ lúc nào,
                    không cần giải thích lý do.</li>
              </ul>
              <p>
                Để yêu cầu, hãy gửi thư tới <a class="kame-link" href="mailto:<?= kame_e(PP_EMAIL) ?>"><?= kame_e(PP_EMAIL) ?></a>
                hoặc nhắn qua trang <a class="kame-link" href="/sv1/ho-tro">Hỗ trợ</a>. Chúng tôi sẽ hỏi vài câu
                để chắc chắn bạn đúng là chủ tài khoản — đây là để bảo vệ bạn, tránh việc người khác xoá tài khoản
                của bạn.
              </p>
              <div class="kame-alert kame-alert--danger" style="margin-top:var(--kame-space-4)">
                <span aria-hidden="true">!</span>
                <span>
                  <strong>Xoá tài khoản là không lấy lại được.</strong>
                  Nhân vật, vật phẩm, cấp độ và mọi thứ bạn đã làm được trong game sẽ mất vĩnh viễn.
                  Hãy cân nhắc kỹ trước khi gửi yêu cầu.
                </span>
                <span></span>
              </div>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('mat-khau'); ?>
              <p>
                Mật khẩu của bạn <strong>không được lưu ở dạng chữ đọc được</strong>. Chúng tôi lưu nó dưới dạng
                một chuỗi đã mã hoá một chiều bằng phương pháp PBKDF2 — nghĩa là từ chuỗi đó
                <strong>không thể suy ngược ra mật khẩu gốc</strong>. Khi bạn đăng nhập, hệ thống mã hoá lại
                mật khẩu bạn vừa gõ rồi so hai chuỗi với nhau.
              </p>
              <p>
                Hệ quả thực tế: <strong>ngay cả ban quản trị cũng không đọc được mật khẩu của bạn.</strong>
                Vì vậy nếu có ai đó nhắn tin tự xưng là quản trị viên và hỏi mật khẩu của bạn, đó chắc chắn
                là kẻ lừa đảo — chúng tôi không bao giờ hỏi mật khẩu.
              </p>
              <p>Vài điều bạn nên tự làm để giữ tài khoản an toàn:</p>
              <ul>
                <li>Đặt mật khẩu <strong>khác</strong> với mật khẩu bạn đang dùng ở Facebook, email hay ngân hàng.</li>
                <li>Không cho mượn tài khoản, kể cả bạn thân.</li>
                <li>Không nhập tài khoản vào các trang "tăng sức mạnh", "nhận quà miễn phí" bên ngoài.</li>
              </ul>
              <p>
                Chúng tôi cố gắng hết sức để bảo vệ hệ thống, nhưng không dịch vụ trực tuyến nào dám hứa
                an toàn tuyệt đối. Nếu xảy ra sự cố ảnh hưởng tới dữ liệu người chơi, chúng tôi sẽ thông báo
                công khai trên trang <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a>.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('cookie'); ?>
              <p>Trang web này dùng đúng hai loại cookie, cả hai đều cần thiết để trang chạy được:</p>
              <ul>
                <li><strong>Cookie phiên đăng nhập</strong> — để bạn không phải đăng nhập lại sau mỗi lần
                    chuyển trang.</li>
                <li><strong>Cookie chống spam</strong> — giúp máy chủ phân biệt người thật với chương trình
                    tự động gửi hàng nghìn yêu cầu.</li>
              </ul>
              <p>
                <strong>Không có cookie quảng cáo, không có mã theo dõi của bên thứ ba.</strong>
                Chúng tôi không theo dõi bạn sang các trang web khác.
              </p>
              <p>
                Bạn có thể xoá cookie trong phần cài đặt của trình duyệt bất cứ lúc nào. Hệ quả là bạn sẽ
                bị đăng xuất và phải đăng nhập lại.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('tre-em'); ?>
              <p>
                Trò chơi này <strong>dành cho người chơi từ 12 tuổi trở lên</strong>. Chúng tôi không chủ ý
                thu thập thông tin của trẻ nhỏ hơn độ tuổi đó.
              </p>
              <p>
                Nếu bạn là phụ huynh và phát hiện con mình dưới 12 tuổi đã tạo tài khoản, hãy nhắn cho chúng tôi
                qua <a class="kame-link" href="mailto:<?= kame_e(PP_EMAIL) ?>"><?= kame_e(PP_EMAIL) ?></a>.
                Chúng tôi sẽ xoá tài khoản đó cùng toàn bộ thông tin liên quan.
              </p>
              <p>
                Cũng xin nhắc: chơi game quá 180 phút mỗi ngày là có hại cho sức khoẻ, nhất là với người
                đang tuổi đi học.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('thay-doi'); ?>
              <p>
                Chính sách này có thể được sửa khi hệ thống thay đổi. Mỗi lần sửa, chúng tôi đổi
                <strong>ngày cập nhật</strong> ghi ở đầu trang.
              </p>
              <p>
                Nếu có thay đổi lớn — ví dụ chúng tôi bắt đầu thu thập thêm một loại thông tin mới — chúng tôi
                sẽ đăng thông báo ở trang <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a> chứ không lặng lẽ
                sửa. Bạn tiếp tục chơi sau khi thay đổi được đăng thì xem như bạn đồng ý với bản mới.
              </p>
            <?php pp_dong_muc(); ?>

            <?php pp_mo_muc('lien-he'); ?>
              <p>
                Có gì thắc mắc về chính sách này, hoặc muốn yêu cầu xem / sửa / xoá thông tin của mình,
                hãy liên hệ:
              </p>
              <ul>
                <li><strong>Thư điện tử:</strong>
                    <a class="kame-link" href="mailto:<?= kame_e(PP_EMAIL) ?>"><?= kame_e(PP_EMAIL) ?></a></li>
                <li><strong>Trang hỗ trợ:</strong>
                    <a class="kame-link" href="/sv1/ho-tro">Trung tâm hỗ trợ Ngọc Rồng Kame</a></li>
                <li><strong>Kênh cộng đồng:</strong> Zalo, Facebook, YouTube và TikTok — địa chỉ các kênh
                    được liệt kê ở trang <a class="kame-link" href="/sv1/ho-tro">Hỗ trợ</a>.</li>
              </ul>
              <p>
                Chúng tôi trả lời trong khả năng của một đội ngũ nhỏ, nên có thể mất vài ngày.
                Mong bạn thông cảm.
              </p>
            <?php pp_dong_muc(); ?>

          </div>
        </div>

      </div>
    </div>
  </section>

  <!-- ===== Liên kết sang trang điều khoản =============================== -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container">
      <div class="kpl-cheo">
        <div>
          <h2 class="kame-heading kame-heading--md">Còn một văn bản nữa</h2>
          <p>
            Trang này chỉ nói về thông tin cá nhân. Các quy định về cách cư xử trong game, vật phẩm ảo
            và những việc có thể khiến tài khoản bị khoá nằm ở trang Điều khoản sử dụng.
          </p>
        </div>
        <div class="kame-cluster">
          <a class="kame-button" href="/sv1/dieu-khoan">Xem điều khoản sử dụng</a>
          <a class="kame-button kame-button--secondary" href="/sv1/ho-tro">Cần hỗ trợ</a>
        </div>
      </div>
    </div>
  </section>

</main>

<?php
kame_auth_modals();
kame_footer();
