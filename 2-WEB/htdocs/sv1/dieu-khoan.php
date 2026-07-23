<?php
/* ==========================================================================
   NGỌC RỒNG KAME — ĐIỀU KHOẢN SỬ DỤNG

   Truy cập qua /sv1/dieu-khoan (mod_rewrite trong .htaccess tự thêm .php).
   Trang song sinh là /privacy-policy (chính sách bảo mật) ở gốc web;
   hai trang trỏ qua lại lẫn nhau ở cuối.

   CSS dùng chung tiền tố kpl- giống trang chính sách bảo mật, khai báo ngay
   trong trang này để hai file độc lập với nhau — sửa một trang không làm
   trang kia lệch giao diện.
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';

/* Ngày cập nhật hiển thị ở đầu trang. Sửa ở đây là đổi cả trang. */
const DK_NGAY_CAP_NHAT = '21/07/2026';

/* Địa chỉ liên hệ, giống trang chính sách bảo mật. */
const DK_EMAIL = 'support@nrokame.com';

/* --- Danh sách mục ------------------------------------------------------
   Khoá mảng là id của thẻ <section>, dùng cho cả mục lục lẫn tiêu đề trong
   thân trang, nên mục lục không bao giờ lệch với nội dung.
   ------------------------------------------------------------------------ */
$dkMuc = [
    'gioi-thieu'  => 'Trang này nói về chuyện gì',
    'dieu-kien'   => 'Ai được chơi',
    'tai-khoan'   => 'Tài khoản là của bạn, và bạn chịu trách nhiệm về nó',
    'ung-xu'      => 'Cách cư xử trong game',
    'phan-mem'    => 'Cấm phần mềm ngoài và khai thác lỗi',
    'vat-pham'    => 'Vật phẩm và tài sản trong game',
    'vi-pham'     => 'Vi phạm thì bị xử lý thế nào',
    'bao-tri'     => 'Bảo trì và tạm ngừng dịch vụ',
    'trach-nhiem' => 'Chúng tôi chịu trách nhiệm tới đâu',
    'thay-doi'    => 'Khi điều khoản này thay đổi',
    'lien-he'     => 'Liên hệ với chúng tôi',
];

/**
 * Mở một mục: in thẻ <section> kèm tiêu đề lấy từ $dkMuc.
 * Số thứ tự do CSS đếm, không gõ tay để khỏi sai khi thêm bớt mục.
 */
function dk_mo_muc(string $id): void
{
    global $dkMuc;
    $tieuDe = $dkMuc[$id] ?? $id;
    echo '<section class="kpl-muc" id="' . kame_e($id) . '">';
    echo '<h2 class="kame-heading kpl-muc__tieude">' . kame_e($tieuDe) . '</h2>';
}

/** Đóng mục vừa mở. */
function dk_dong_muc(): void
{
    echo '</section>';
}

kame_head(
    'Điều khoản sử dụng',
    'Quy định khi chơi Ngọc Rồng Kame: cách cư xử, cấm phần mềm ngoài, vật phẩm trong game và cách xử lý vi phạm.'
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
      <h1 class="kame-heading kame-heading--lg kpl-band__title">Điều khoản sử dụng</h1>
      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <span>Điều khoản sử dụng</span>
      </nav>
      <div class="kpl-band__ngay">Cập nhật ngày <?= kame_e(DK_NGAY_CAP_NHAT) ?></div>
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
                <?php foreach ($dkMuc as $id => $tieuDe): ?>
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

            <?php dk_mo_muc('gioi-thieu'); ?>
              <p>
                Ngọc Rồng Kame là một máy chủ trò chơi do cá nhân vận hành, mở cửa miễn phí cho mọi người
                cùng chơi. Để chỗ này còn vui thì cần vài quy định chung. Trang này liệt kê những quy định đó:
                <strong>bạn được làm gì, không được làm gì, và chuyện gì xảy ra nếu vi phạm</strong>.
              </p>
              <p>
                Khi bạn tạo tài khoản hoặc đăng nhập vào game, xem như bạn đã đọc và đồng ý với những điều
                ghi ở đây. Nếu có chỗ nào bạn không đồng ý, cách xử lý đơn giản là đừng dùng dịch vụ.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('dieu-kien'); ?>
              <p>
                Trò chơi này <strong>dành cho người từ 12 tuổi trở lên</strong>. Nếu bạn dưới 18 tuổi,
                hãy hỏi ý cha mẹ hoặc người giám hộ trước khi tạo tài khoản.
              </p>
              <p>
                Dịch vụ được cung cấp <strong>miễn phí và không kèm bảo đảm nào</strong>. Chúng tôi không
                bán tài khoản, không thu phí đăng ký, và có quyền từ chối phục vụ một người cụ thể nếu người
                đó liên tục phá hỏng trải nghiệm của người khác.
              </p>
              <p>
                Cũng nói thẳng để không ai hiểu nhầm: đây là <strong>máy chủ tư nhân</strong>, không liên kết
                với chủ sở hữu bản quyền Dragon Ball, không phải sản phẩm chính thức của bất kỳ nhà phát hành nào,
                và không đại diện cho ai ngoài chính ban quản trị máy chủ.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('tai-khoan'); ?>
              <ul>
                <li>Mỗi người tự chịu trách nhiệm về <strong>mọi hoạt động diễn ra trên tài khoản của mình</strong>,
                    kể cả khi do người khác làm.</li>
                <li><strong>Đừng cho mượn tài khoản.</strong> Bạn cho bạn bè mượn rồi người đó gian lận,
                    tài khoản của bạn vẫn bị khoá — chúng tôi không có cách nào phân biệt được.</li>
                <li>Ban quản trị <strong>không bao giờ hỏi mật khẩu của bạn</strong>. Ai nhắn tin hỏi mật khẩu
                    thì chắc chắn là kẻ lừa đảo, dù họ tự xưng là gì.</li>
                <li><strong>Mua bán, trao đổi tài khoản</strong> giữa người chơi với nhau là việc chúng tôi
                    không công nhận và không hỗ trợ. Mất tài khoản vì mua bán thì chúng tôi không lấy lại giúp được.</li>
                <li>Đặt tên nhân vật và tên tài khoản <strong>lịch sự</strong>. Tên tục tĩu, xúc phạm, hoặc cố ý
                    giả danh quản trị viên sẽ bị đổi hoặc khoá.</li>
              </ul>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('ung-xu'); ?>
              <p>Nguyên tắc chung rất đơn giản: <strong>đừng làm hỏng cuộc chơi của người khác.</strong></p>

              <h3>Những việc bị cấm</h3>
              <ul>
                <li><strong>Chửi bới, xúc phạm, phân biệt đối xử</strong> — nhắm vào cá nhân, gia đình, giới tính,
                    vùng miền, tôn giáo hay dân tộc.</li>
                <li><strong>Ngôn từ tục tĩu và nội dung khiêu dâm</strong> trong tên nhân vật, chat công khai
                    hoặc tin nhắn riêng.</li>
                <li><strong>Spam</strong> — lặp lại một câu liên tục, gửi hàng loạt tin nhắn riêng, hoặc phát tán
                    quảng cáo cho máy chủ khác trong khung chat.</li>
                <li><strong>Lừa đảo</strong> — hứa hẹn đổi vật phẩm rồi ăn quỵt, dụ người khác nhập tài khoản vào
                    trang giả mạo, giả danh quản trị viên để xin đồ hoặc xin mật khẩu.</li>
                <li><strong>Đe doạ và quấy rối</strong> — bám theo một người chơi để phá, doạ tấn công ngoài đời,
                    hoặc bêu riếu thông tin cá nhân của người khác.</li>
                <li><strong>Rao bán vật phẩm hoặc tài khoản lấy tiền thật</strong>, dưới bất kỳ hình thức nào.</li>
              </ul>

              <h3>Những việc nên làm</h3>
              <ul>
                <li>Báo cho ban quản trị khi thấy ai đó gian lận hoặc lừa đảo, kèm ảnh chụp màn hình nếu có.</li>
                <li>Báo lỗi game thay vì âm thầm khai thác nó.</li>
                <li>Nhường nhịn người mới. Ai cũng từng là người mới.</li>
              </ul>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('phan-mem'); ?>
              <p>
                Chỉ được chơi bằng <strong>bản game tải từ trang này</strong>. Cụ thể, những việc sau bị cấm:
              </p>
              <ul>
                <li>Dùng <strong>phần mềm hack, auto, bot, tool</strong> hoặc bất kỳ chương trình nào chơi thay bạn
                    hoặc cho bạn lợi thế mà người chơi bình thường không có.</li>
                <li><strong>Sửa đổi bản game</strong>, chỉnh bộ nhớ, giả lập nhiều tài khoản chạy tự động, hoặc
                    can thiệp vào dữ liệu gửi giữa máy bạn và máy chủ.</li>
                <li><strong>Khai thác lỗi</strong> — khi phát hiện một lỗi cho phép nhân đôi vật phẩm, lên cấp bất
                    thường hay chui qua địa hình, việc đúng đắn là <strong>báo cho ban quản trị</strong>.
                    Cố tình lợi dụng lỗi để trục lợi bị xử lý như gian lận.</li>
                <li><strong>Tấn công hệ thống</strong> — dò mật khẩu, làm nghẽn máy chủ, thử xâm nhập trang web
                    hoặc cơ sở dữ liệu.</li>
                <li>Dùng lại mã nguồn, hình ảnh, âm thanh của máy chủ này để dựng máy chủ khác.</li>
              </ul>
              <p>
                Với các trường hợp gian lận, chúng tôi có thể <strong>thu hồi toàn bộ vật phẩm và chỉ số</strong>
                có được nhờ gian lận, kể cả khi chúng đã được chuyển sang tài khoản khác.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('vat-pham'); ?>
              <p>
                Đây là mục hay gây hiểu nhầm nhất, nên xin nói thật rõ:
              </p>
              <ul>
                <li>Mọi thứ trong game — nhân vật, vật phẩm, ngọc, vàng, cấp độ, sức mạnh — đều là
                    <strong>dữ liệu nằm trên máy chủ của chúng tôi</strong>, không phải tài sản riêng của bạn
                    theo nghĩa pháp lý.</li>
                <li>Bạn được <strong>quyền sử dụng</strong> những thứ đó trong phạm vi trò chơi, chừng nào
                    tài khoản còn hoạt động và máy chủ còn chạy.</li>
                <li><strong>Vật phẩm ảo không quy đổi ngược ra tiền mặt.</strong> Không có chuyện hoàn tiền,
                    đổi ra tiền, hay bồi thường bằng tiền cho vật phẩm bị mất.</li>
                <li>Ban quản trị có quyền <strong>chỉnh sửa, thu hồi hoặc xoá</strong> vật phẩm và chỉ số khi cần
                    — ví dụ để cân bằng game, sửa lỗi, hoặc xử lý gian lận.</li>
                <li>Nếu bạn có nạp tiền ủng hộ máy chủ, khoản đó được xem là <strong>ủng hộ tự nguyện</strong>
                    để duy trì máy chủ, không phải giao dịch mua bán hàng hoá và không được hoàn lại.</li>
              </ul>
              <div class="kame-alert kame-alert--info" style="margin-top:var(--kame-space-4)">
                <span aria-hidden="true">!</span>
                <span>
                  Nói ngắn gọn: hãy coi thời gian bạn bỏ vào game là <strong>để giải trí</strong>,
                  đừng coi nó là một khoản đầu tư.
                </span>
                <span></span>
              </div>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('vi-pham'); ?>
              <p>
                Tuỳ mức độ và tuỳ việc tái phạm, chúng tôi có thể áp dụng một hoặc nhiều biện pháp sau:
              </p>
              <ul>
                <li><strong>Nhắc nhở</strong> — với lỗi nhẹ, lần đầu.</li>
                <li><strong>Cấm chat</strong> có thời hạn — với spam và ngôn từ không phù hợp.</li>
                <li><strong>Thu hồi vật phẩm hoặc chỉ số</strong> có được không đúng cách.</li>
                <li><strong>Khoá tài khoản có thời hạn</strong> — vài ngày đến vài tuần.</li>
                <li><strong>Khoá tài khoản vĩnh viễn</strong> — với gian lận nặng, lừa đảo, hoặc tấn công hệ thống.</li>
              </ul>
              <p>
                Chúng tôi cố gắng xử lý công bằng, nhưng đây là <strong>quyết định của ban quản trị</strong>
                và trong nhiều trường hợp không thể công bố chi tiết bằng chứng (vì làm vậy sẽ chỉ giúp
                người gian lận biết cách né lần sau).
              </p>
              <p>
                Nếu bạn tin là mình bị xử lý oan, hãy khiếu nại qua trang
                <a class="kame-link" href="/sv1/ho-tro">Hỗ trợ</a> hoặc gửi thư tới
                <a class="kame-link" href="mailto:<?= kame_e(DK_EMAIL) ?>"><?= kame_e(DK_EMAIL) ?></a>.
                Chúng tôi sẽ xem lại. Xin đừng lập tài khoản mới để tiếp tục — làm vậy chỉ khiến
                trường hợp của bạn nặng thêm.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('bao-tri'); ?>
              <p>
                Máy chủ cần được bảo trì để cập nhật nội dung và sửa lỗi. Trong lúc bảo trì bạn sẽ
                <strong>không đăng nhập được</strong>. Chúng tôi cố gắng báo trước ở trang
                <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a>, nhưng khi có sự cố khẩn cấp thì
                có thể phải tắt máy chủ ngay mà không kịp báo.
              </p>
              <p>
                Chúng tôi cũng có thể <strong>thay đổi, tạm ngừng hoặc đóng cửa dịch vụ</strong> bất cứ lúc nào.
                Đây là máy chủ do cá nhân duy trì bằng chi phí tự bỏ ra, nên không thể hứa nó tồn tại mãi mãi.
                Nếu quyết định đóng cửa, chúng tôi sẽ thông báo trước trong khả năng có thể.
              </p>
              <p>
                Việc mất dữ liệu nhân vật do sự cố kỹ thuật, hỏng ổ cứng hoặc lỗi cập nhật là điều
                <strong>có thể xảy ra</strong>. Chúng tôi cố gắng sao lưu thường xuyên, nhưng không cam kết
                khôi phục được mọi thứ.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('trach-nhiem'); ?>
              <p>
                Dịch vụ được cung cấp <strong>"nguyên trạng"</strong>, miễn phí, không kèm bảo đảm.
                Trong phạm vi pháp luật cho phép, chúng tôi <strong>không chịu trách nhiệm bồi thường</strong>
                cho các thiệt hại phát sinh từ việc bạn dùng dịch vụ, bao gồm nhưng không giới hạn ở:
              </p>
              <ul>
                <li>Mất nhân vật, vật phẩm hoặc tiến độ chơi vì bất kỳ lý do gì.</li>
                <li>Máy chủ ngừng hoạt động, giật lag, hoặc đóng cửa.</li>
                <li>Thiệt hại do bạn bị người chơi khác lừa đảo hoặc quấy rối.</li>
                <li>Thiệt hại do bạn tự để lộ mật khẩu hoặc cho người khác mượn tài khoản.</li>
                <li>Thời gian, công sức bạn đã bỏ vào trò chơi.</li>
              </ul>
              <p>
                Điều này không loại trừ trách nhiệm mà pháp luật bắt buộc phải có. Nó chỉ nói rõ một chuyện
                thực tế: đây là một máy chủ nhỏ do cá nhân vận hành, không phải một công ty có bộ phận
                bồi thường khách hàng.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('thay-doi'); ?>
              <p>
                Điều khoản này có thể được sửa khi game thay đổi hoặc khi xuất hiện tình huống chưa được
                quy định. Mỗi lần sửa, chúng tôi đổi <strong>ngày cập nhật</strong> ghi ở đầu trang.
              </p>
              <p>
                Với những thay đổi đáng kể, chúng tôi sẽ đăng thông báo ở trang
                <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a>. Bạn tiếp tục chơi sau khi thay đổi
                được đăng thì xem như bạn đồng ý với bản mới.
              </p>
            <?php dk_dong_muc(); ?>

            <?php dk_mo_muc('lien-he'); ?>
              <p>
                Có gì chưa rõ trong điều khoản này, muốn khiếu nại một quyết định xử lý, hoặc muốn báo cáo
                một người chơi vi phạm, hãy liên hệ:
              </p>
              <ul>
                <li><strong>Thư điện tử:</strong>
                    <a class="kame-link" href="mailto:<?= kame_e(DK_EMAIL) ?>"><?= kame_e(DK_EMAIL) ?></a></li>
                <li><strong>Trang hỗ trợ:</strong>
                    <a class="kame-link" href="/sv1/ho-tro">Trung tâm hỗ trợ Ngọc Rồng Kame</a></li>
                <li><strong>Kênh cộng đồng:</strong> Zalo, Facebook, YouTube và TikTok — địa chỉ các kênh
                    được liệt kê ở trang <a class="kame-link" href="/sv1/ho-tro">Hỗ trợ</a>.</li>
              </ul>
              <p>
                Chúng tôi trả lời trong khả năng của một đội ngũ nhỏ, nên có thể mất vài ngày.
                Mong bạn thông cảm.
              </p>
            <?php dk_dong_muc(); ?>

          </div>
        </div>

      </div>
    </div>
  </section>

  <!-- ===== Liên kết sang trang chính sách bảo mật ======================= -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container">
      <div class="kpl-cheo">
        <div>
          <h2 class="kame-heading kame-heading--md">Còn một văn bản nữa</h2>
          <p>
            Trang này chỉ nói về quy định khi chơi. Chuyện chúng tôi giữ những thông tin nào của bạn,
            dùng để làm gì và bạn có quyền gì với chúng nằm ở trang Chính sách bảo mật.
          </p>
        </div>
        <div class="kame-cluster">
          <a class="kame-button" href="/privacy-policy">Xem chính sách bảo mật</a>
          <a class="kame-button kame-button--secondary" href="/sv1/ho-tro">Cần hỗ trợ</a>
        </div>
      </div>
    </div>
  </section>

</main>

<?php
kame_auth_modals();
kame_footer();
