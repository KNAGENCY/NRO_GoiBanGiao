<?php
/* ==========================================================================
   NGỌC RỒNG KAME — TIN TỨC
   Một file lo hai việc:
     · Không có ?bai=      → danh sách (nổi bật + lọc tag + tìm kiếm + phân trang)
     · Có  ?bai=<slug>     → chi tiết bài viết + 3 bài liên quan cùng tag
   Trang này KHÔNG viết SQL, chỉ gọi các hàm của kame/data.php.
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';

/* --- Hằng cấu hình ------------------------------------------------------- */

/** Số bài trong lưới mỗi trang. */
const TT_MOI_TRANG = 9;

/** Nhãn tiếng Việt cho mọi tag đang có trong bảng news. */
const TT_NHAN_TAG = [
    'su-kien'  => 'Sự kiện',
    'cap-nhat' => 'Cập nhật',
    'huong-dan' => 'Hướng dẫn',
    'giftcode' => 'Giftcode',
    'bao-tri'  => 'Bảo trì',
    'tin-tuc'  => 'Tin tức',
];

/** Các tab hiển thị trên thanh lọc, theo đúng thứ tự yêu cầu. */
const TT_TAB = ['su-kien', 'cap-nhat', 'huong-dan', 'giftcode', 'bao-tri'];

/** Màu badge cho từng tag, dùng lại biến thể có sẵn của ui-kit. */
const TT_MAU_TAG = [
    'su-kien'  => 'kame-badge--hot',
    'cap-nhat' => 'kame-badge--info',
    'huong-dan' => 'kame-badge--success',
    'giftcode' => 'tt-badge--gold',
    'bao-tri'  => '',
    'tin-tuc'  => '',
];

/* --- Tiện ích ------------------------------------------------------------ */

/** Nhãn tiếng Việt của một tag; tag lạ thì trả về chính nó. */
function tt_nhan_tag(?string $tag): string
{
    $tag = (string)$tag;
    return TT_NHAN_TAG[$tag] ?? ($tag === '' ? 'Tin tức' : $tag);
}

/** Lớp màu badge của một tag. */
function tt_mau_tag(?string $tag): string
{
    return TT_MAU_TAG[(string)$tag] ?? '';
}

/** Ngày kiểu 20.07.2026. Chuỗi hỏng thì trả rỗng, không in "01.01.1970". */
function tt_ngay(?string $sql): string
{
    $ts = strtotime((string)$sql);
    return $ts ? date('d.m.Y', $ts) : '';
}

/** Ngày dạng máy đọc cho thẻ <time datetime>. */
function tt_ngay_iso(?string $sql): string
{
    $ts = strtotime((string)$sql);
    return $ts ? date('Y-m-d', $ts) : '';
}

/** Tìm chuỗi con không phân biệt hoa thường, có hỗ trợ tiếng Việt. */
function tt_chua(string $chuoi, string $tim): bool
{
    if ($tim === '') {
        return true;
    }
    if (function_exists('mb_stripos')) {
        return mb_stripos($chuoi, $tim, 0, 'UTF-8') !== false;
    }
    return stripos($chuoi, $tim) !== false;
}

/**
 * Toàn bộ bài viết, lấy theo từng lô 100 bài (giới hạn của kame_news).
 * Lọc tag và tìm kiếm làm bằng PHP vì tầng dữ liệu không nhận điều kiện lọc.
 *
 * @return array<int,array>
 */
function tt_tat_ca_bai(): array
{
    static $cache = null;
    if ($cache !== null) {
        return $cache;
    }

    $cache = [];
    $tong  = kame_news_count();
    for ($offset = 0; $offset < $tong; $offset += 100) {
        $lo = kame_news(100, $offset);
        if (!$lo) {
            break;
        }
        foreach ($lo as $bai) {
            $cache[] = $bai;
        }
        if (count($lo) < 100) {
            break;
        }
    }
    return $cache;
}

/** Ảnh đại diện bài viết: ảnh thật nếu có, không thì ô giữ chỗ. */
function tt_anh_bai(array $bai, string $ty_le = '16/9'): string
{
    $anh = trim((string)($bai['image'] ?? ''));
    if ($anh !== '') {
        return '<img class="tt-anh" style="aspect-ratio:' . ($ty_le === '16/9' ? '16/9' : '800/450') . '"'
             . ' src="' . kame_e($anh) . '" alt="' . kame_e((string)($bai['title'] ?? '')) . '" loading="lazy">';
    }
    return kame_placeholder('tin-mac-dinh.jpg', $ty_le);
}

/* --- Đọc tham số URL ----------------------------------------------------- */

$slug     = isset($_GET['bai']) ? trim((string)$_GET['bai']) : '';
$tagLoc   = isset($_GET['tag']) ? trim((string)$_GET['tag']) : '';
$tuKhoa   = isset($_GET['q'])   ? trim((string)$_GET['q'])   : '';
$trang    = isset($_GET['trang']) ? max(1, (int)$_GET['trang']) : 1;

// Tag lạ coi như không lọc, tránh hiện tab rỗng vô nghĩa.
if ($tagLoc !== '' && !in_array($tagLoc, TT_TAB, true)) {
    $tagLoc = '';
}
// Cắt từ khoá quá dài cho gọn URL và ô nhập.
if (function_exists('mb_substr')) {
    $tuKhoa = mb_substr($tuKhoa, 0, 80, 'UTF-8');
} else {
    $tuKhoa = substr($tuKhoa, 0, 80);
}

/* ==========================================================================
   NHÁNH 1 — CHI TIẾT BÀI VIẾT
   ========================================================================== */
if ($slug !== '') {

    $bai = kame_news_one($slug);

    if ($bai === null) {
        // Phải trả đúng mã 404. Nếu trả 200, Google sẽ lập chỉ mục mọi slug rác
        // và trình duyệt tưởng bài viết có thật.
        http_response_code(404);
        kame_head('Không tìm thấy bài viết', 'Bài viết bạn tìm không tồn tại hoặc đã bị gỡ.');
        kame_header('tintuc');
        ?>
        <main class="kame-container kame-section kame-page-wrap">
          <nav class="kame-breadcrumb" aria-label="Đường dẫn">
            <a href="/sv1/trangchu">Trang chủ</a>
            <a href="/sv1/tin-tuc">Tin tức</a>
            <span class="kame-muted">Không tìm thấy</span>
          </nav>

          <div class="kame-alert kame-alert--danger tt-alert" style="margin-top:24px">
            <span aria-hidden="true">!</span>
            <span>
              <strong>Không tìm thấy bài viết</strong><br>
              Bài viết với đường dẫn <em><?= kame_e($slug) ?></em> không tồn tại hoặc đã được gỡ xuống.
            </span>
            <span></span>
          </div>

          <p style="margin-top:24px">
            <a class="kame-button" href="/sv1/tin-tuc">Quay lại tin tức</a>
          </p>
        </main>
        <?php
        kame_auth_modals();
        kame_footer();
        exit;
    }

    // 3 bài liên quan: cùng tag, bỏ chính nó ra.
    $lienQuan = [];
    foreach (tt_tat_ca_bai() as $khac) {
        if ((string)$khac['slug'] === (string)$bai['slug']) {
            continue;
        }
        if ((string)$khac['tag'] !== (string)$bai['tag']) {
            continue;
        }
        $lienQuan[] = $khac;
        if (count($lienQuan) >= 3) {
            break;
        }
    }

    kame_head((string)$bai['title'], (string)($bai['excerpt'] ?? ''));
    kame_header('tintuc');
    ?>
    <main class="kame-container kame-section kame-page-wrap">

      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <a href="/sv1/tin-tuc">Tin tức</a>
        <a href="/sv1/tin-tuc?tag=<?= kame_e(rawurlencode((string)$bai['tag'])) ?>"><?= kame_e(tt_nhan_tag((string)$bai['tag'])) ?></a>
        <span class="kame-muted"><?= kame_e((string)$bai['title']) ?></span>
      </nav>

      <article class="tt-bai">

        <header class="tt-bai__head">
          <span class="kame-badge <?= kame_e(tt_mau_tag((string)$bai['tag'])) ?>"><?= kame_e(tt_nhan_tag((string)$bai['tag'])) ?></span>
          <h1 class="kame-heading kame-heading--lg tt-bai__title"><?= kame_e((string)$bai['title']) ?></h1>
          <p class="kame-muted tt-bai__meta">
            Đăng ngày
            <time datetime="<?= kame_e(tt_ngay_iso((string)$bai['created_at'])) ?>"><?= kame_e(tt_ngay((string)$bai['created_at'])) ?></time>
          </p>
        </header>

        <div class="tt-bai__anh">
          <?= tt_anh_bai($bai, '800/450') ?>
        </div>

        <?php if (trim((string)($bai['excerpt'] ?? '')) !== ''): ?>
          <p class="tt-bai__sapo"><?= kame_e((string)$bai['excerpt']) ?></p>
        <?php endif; ?>

        <?php /* body đã là HTML sẵn trong CSDL — in nguyên, không escape */ ?>
        <div class="tt-body">
          <?= (string)($bai['body'] ?? '') ?>
        </div>

      </article>

      <p class="tt-quaylai">
        <a class="kame-button kame-button--secondary" href="/sv1/tin-tuc">« Về trang tin tức</a>
      </p>

      <section class="kame-section--tight" aria-labelledby="tt-lienquan">
        <div class="kame-section-title">
          <h2 class="kame-heading kame-heading--md" id="tt-lienquan">Tin liên quan</h2>
          <a class="kame-link" href="/sv1/tin-tuc?tag=<?= kame_e(rawurlencode((string)$bai['tag'])) ?>">Xem tất cả <?= kame_e(tt_nhan_tag((string)$bai['tag'])) ?> »</a>
        </div>

        <?php if (!$lienQuan): ?>
          <p class="kame-muted">Chưa có dữ liệu.</p>
        <?php else: ?>
          <div class="kame-grid">
            <?php foreach ($lienQuan as $lq): ?>
              <article class="kame-card tt-card tt-card--doc kame-col-4">
                <a class="tt-card__media" href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$lq['slug'])) ?>" tabindex="-1" aria-hidden="true">
                  <?= tt_anh_bai($lq) ?>
                </a>
                <div class="tt-card__body">
                  <div class="tt-card__top">
                    <span class="kame-badge <?= kame_e(tt_mau_tag((string)$lq['tag'])) ?>"><?= kame_e(tt_nhan_tag((string)$lq['tag'])) ?></span>
                    <time class="tt-card__date" datetime="<?= kame_e(tt_ngay_iso((string)$lq['created_at'])) ?>"><?= kame_e(tt_ngay((string)$lq['created_at'])) ?></time>
                  </div>
                  <h3 class="tt-card__title">
                    <a href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$lq['slug'])) ?>"><?= kame_e((string)$lq['title']) ?></a>
                  </h3>
                </div>
              </article>
            <?php endforeach; ?>
          </div>
        <?php endif; ?>
      </section>

    </main>
    <?php
    kame_auth_modals();
    kame_footer();
    exit;
}

/* ==========================================================================
   NHÁNH 2 — DANH SÁCH BÀI VIẾT
   ========================================================================== */

$tatCa = tt_tat_ca_bai();

// Lọc theo tag rồi lọc theo từ khoá (tìm trong tiêu đề và mô tả ngắn).
$locDuoc = [];
foreach ($tatCa as $bai) {
    if ($tagLoc !== '' && (string)$bai['tag'] !== $tagLoc) {
        continue;
    }
    if ($tuKhoa !== ''
        && !tt_chua((string)($bai['title'] ?? ''), $tuKhoa)
        && !tt_chua((string)($bai['excerpt'] ?? ''), $tuKhoa)) {
        continue;
    }
    $locDuoc[] = $bai;
}

$dangLoc = ($tagLoc !== '' || $tuKhoa !== '');

// Không lọc gì thì bài mới nhất tách ra làm thẻ nổi bật, không lặp lại trong lưới.
$noiBat = null;
$hangDoi = $locDuoc;
if (!$dangLoc && $hangDoi) {
    $noiBat  = array_shift($hangDoi);
}

$tongLuoi = count($hangDoi);
$tongTrang = max(1, (int)ceil($tongLuoi / TT_MOI_TRANG));
if ($trang > $tongTrang) {
    $trang = $tongTrang;
}
$danhSach = array_slice($hangDoi, ($trang - 1) * TT_MOI_TRANG, TT_MOI_TRANG);

/** Dựng URL giữ nguyên bộ lọc hiện tại, chỉ thay vài tham số. */
function tt_url(array $doi = []): string
{
    global $tagLoc, $tuKhoa, $trang;

    $tham = ['tag' => $tagLoc, 'q' => $tuKhoa, 'trang' => $trang];
    foreach ($doi as $k => $v) {
        $tham[$k] = $v;
    }
    if ((int)$tham['trang'] <= 1) {
        unset($tham['trang']);
    }
    $tham = array_filter($tham, static fn($v) => $v !== '' && $v !== null);

    return '/sv1/tin-tuc' . ($tham ? '?' . http_build_query($tham) : '');
}

$moTa = 'Tin tức, sự kiện, cập nhật và giftcode mới nhất của máy chủ Ngọc Rồng Kame.';
kame_head('Tin tức', $moTa);
kame_header('tintuc');
?>
<main class="kame-page-wrap">

  <!-- Đầu trang -->
  <section class="tt-dautrang">
    <div class="kame-container">
      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <span class="kame-muted">Tin tức</span>
      </nav>
      <h1 class="kame-heading kame-heading--lg tt-dautrang__title">Tin tức &amp; sự kiện</h1>
      <p class="kame-muted tt-dautrang__desc">
        Cập nhật phiên bản, lịch sự kiện, hướng dẫn chơi và giftcode của Ngọc Rồng Kame.
      </p>
    </div>
  </section>

  <div class="kame-container kame-section">

    <?php if (!$tatCa): ?>

      <div class="kame-alert tt-alert">
        <span aria-hidden="true">i</span>
        <span>Chưa có bài viết.</span>
        <span></span>
      </div>

    <?php else: ?>

      <!-- Bài nổi bật (chỉ ở trang 1 khi không lọc) -->
      <?php if ($noiBat !== null && $trang === 1): ?>
        <article class="kame-card tt-noibat">
          <a class="tt-noibat__media" href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$noiBat['slug'])) ?>" tabindex="-1" aria-hidden="true">
            <?= tt_anh_bai($noiBat, '800/450') ?>
          </a>
          <div class="tt-noibat__body">
            <div class="tt-card__top">
              <span class="kame-badge kame-badge--hot">Nổi bật</span>
              <span class="kame-badge <?= kame_e(tt_mau_tag((string)$noiBat['tag'])) ?>"><?= kame_e(tt_nhan_tag((string)$noiBat['tag'])) ?></span>
              <time class="tt-card__date" datetime="<?= kame_e(tt_ngay_iso((string)$noiBat['created_at'])) ?>"><?= kame_e(tt_ngay((string)$noiBat['created_at'])) ?></time>
            </div>
            <h2 class="kame-heading kame-heading--md tt-noibat__title">
              <a href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$noiBat['slug'])) ?>"><?= kame_e((string)$noiBat['title']) ?></a>
            </h2>
            <p class="tt-noibat__excerpt"><?= kame_e((string)($noiBat['excerpt'] ?? '')) ?></p>
            <p>
              <a class="kame-button" href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$noiBat['slug'])) ?>">Xem chi tiết ›</a>
            </p>
          </div>
        </article>
      <?php endif; ?>

      <!-- Thanh lọc: tab theo tag + ô tìm kiếm -->
      <div class="tt-loc">
        <div class="kame-tabs" role="tablist" aria-label="Lọc theo chuyên mục">
          <a class="kame-tab" role="tab" href="<?= kame_e(tt_url(['tag' => '', 'trang' => 1])) ?>"
             aria-selected="<?= $tagLoc === '' ? 'true' : 'false' ?>">Tất cả</a>
          <?php foreach (TT_TAB as $tag): ?>
            <a class="kame-tab" role="tab" href="<?= kame_e(tt_url(['tag' => $tag, 'trang' => 1])) ?>"
               aria-selected="<?= $tagLoc === $tag ? 'true' : 'false' ?>"><?= kame_e(tt_nhan_tag($tag)) ?></a>
          <?php endforeach; ?>
        </div>

        <form class="tt-timkiem" method="get" action="/sv1/tin-tuc" role="search">
          <?php if ($tagLoc !== ''): ?>
            <input type="hidden" name="tag" value="<?= kame_e($tagLoc) ?>">
          <?php endif; ?>
          <label class="tt-timkiem__label" for="tt-q">Tìm bài viết</label>
          <input class="kame-input" id="tt-q" name="q" type="search" value="<?= kame_e($tuKhoa) ?>"
                 placeholder="Tìm tin tức, sự kiện, giftcode...">
          <button class="kame-button kame-button--sm" type="submit">Tìm</button>
        </form>
      </div>

      <!-- Dòng tóm tắt kết quả -->
      <p class="kame-muted tt-tomtat">
        <?php if ($dangLoc): ?>
          <?php
          // Ghép sẵn thành một câu để dấu chấm không bị tách khỏi chữ.
          $cau = 'Tìm thấy <strong>' . (int)count($locDuoc) . '</strong> bài viết';
          if ($tagLoc !== '') {
              $cau .= ' trong mục <strong>' . kame_e(tt_nhan_tag($tagLoc)) . '</strong>';
          }
          if ($tuKhoa !== '') {
              $cau .= ' khớp từ khoá “<strong>' . kame_e($tuKhoa) . '</strong>”';
          }
          echo $cau . '.';
          ?>
          <a class="kame-link" href="/sv1/tin-tuc">Xoá bộ lọc</a>
        <?php else: ?>
          Tổng cộng <strong><?= (int)count($locDuoc) ?></strong> bài viết. Trang <?= (int)$trang ?>/<?= (int)$tongTrang ?>.
        <?php endif; ?>
      </p>

      <!-- Lưới bài viết -->
      <?php if (!$danhSach): ?>

        <div class="kame-alert tt-alert">
          <span aria-hidden="true">i</span>
          <span>
            <strong>Chưa có bài viết</strong><br>
            Không có bài nào khớp với lựa chọn hiện tại. Hãy thử chuyên mục khác hoặc
            <a class="kame-link" href="/sv1/tin-tuc">xem toàn bộ tin tức</a>.
          </span>
          <span></span>
        </div>

      <?php else: ?>

        <div class="kame-grid tt-luoi">
          <?php foreach ($danhSach as $bai): ?>
            <article class="kame-card tt-card kame-col-6">
              <a class="tt-card__media" href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$bai['slug'])) ?>" tabindex="-1" aria-hidden="true">
                <?= tt_anh_bai($bai) ?>
              </a>
              <div class="tt-card__body">
                <div class="tt-card__top">
                  <span class="kame-badge <?= kame_e(tt_mau_tag((string)$bai['tag'])) ?>"><?= kame_e(tt_nhan_tag((string)$bai['tag'])) ?></span>
                  <time class="tt-card__date" datetime="<?= kame_e(tt_ngay_iso((string)$bai['created_at'])) ?>"><?= kame_e(tt_ngay((string)$bai['created_at'])) ?></time>
                </div>
                <h2 class="tt-card__title">
                  <a href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$bai['slug'])) ?>"><?= kame_e((string)$bai['title']) ?></a>
                </h2>
                <p class="tt-card__excerpt"><?= kame_e((string)($bai['excerpt'] ?? '')) ?></p>
                <a class="tt-card__more" href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$bai['slug'])) ?>">Đọc thêm ›</a>
              </div>
            </article>
          <?php endforeach; ?>
        </div>

        <!-- Phân trang -->
        <?php if ($tongTrang > 1): ?>
          <nav class="kame-pagination tt-phantrang" aria-label="Phân trang tin tức">
            <?php if ($trang > 1): ?>
              <a class="kame-page" href="<?= kame_e(tt_url(['trang' => $trang - 1])) ?>" aria-label="Trang trước">«</a>
            <?php else: ?>
              <span class="kame-page tt-page--tat" aria-hidden="true">«</span>
            <?php endif; ?>

            <?php
            // Cửa sổ tối đa 7 số quanh trang hiện tại cho khỏi tràn hàng.
            $dau = max(1, $trang - 3);
            $cuoi = min($tongTrang, $dau + 6);
            $dau = max(1, $cuoi - 6);
            for ($i = $dau; $i <= $cuoi; $i++):
                if ($i === $trang): ?>
                  <span class="kame-page" aria-current="page"><?= (int)$i ?></span>
                <?php else: ?>
                  <a class="kame-page" href="<?= kame_e(tt_url(['trang' => $i])) ?>"><?= (int)$i ?></a>
                <?php endif;
            endfor;
            ?>

            <?php if ($trang < $tongTrang): ?>
              <a class="kame-page" href="<?= kame_e(tt_url(['trang' => $trang + 1])) ?>" aria-label="Trang sau">»</a>
            <?php else: ?>
              <span class="kame-page tt-page--tat" aria-hidden="true">»</span>
            <?php endif; ?>
          </nav>
        <?php endif; ?>

      <?php endif; ?>

    <?php endif; ?>

  </div>
</main>

<style>
/* ==========================================================================
   CSS riêng của trang tin tức. Đặt tại đây để không đụng theme.css dùng chung.
   Mọi màu đều lấy từ biến của ui-kit.
   ========================================================================== */

/* --- Đầu trang ----------------------------------------------------------- */
.tt-dautrang {
  padding-block: var(--kame-space-6);
  border-bottom: var(--kame-border);
  background: linear-gradient(180deg, rgb(255 200 61 / 22%), transparent);
}
.tt-dautrang__title { margin-top: var(--kame-space-3); }
.tt-dautrang__desc { max-width: 640px; margin: var(--kame-space-3) 0 0; }

/* --- Ảnh bài ------------------------------------------------------------- */
.tt-anh {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* --- Bài nổi bật --------------------------------------------------------- */
.tt-noibat {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 1fr);
  gap: 0;
  margin-bottom: var(--kame-space-6);
}
.tt-noibat__media { display: block; padding: var(--kame-space-4); }
.tt-noibat__body {
  display: grid;
  align-content: center;
  gap: var(--kame-space-3);
  padding: var(--kame-space-5);
  border-left: 2px dashed rgb(96 0 24 / 30%);
}
.tt-noibat__body > p { margin: 0; }
.tt-noibat__title a { color: inherit; text-decoration: none; }
.tt-noibat__title a:hover { color: var(--kame-orange-700); }
.tt-noibat__excerpt { color: rgb(40 21 15 / 78%); }

/* --- Thanh lọc ----------------------------------------------------------- */
.tt-loc {
  display: flex;
  flex-wrap: wrap;
  gap: var(--kame-space-4);
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--kame-space-5);
}
/* min-width:0 là bắt buộc: mặc định flex item không co nhỏ hơn nội dung,
   khiến hàng tab dài hơn màn hình và làm cả trang cuộn ngang ở 375px. */
.tt-loc > * { min-width: 0; }
.tt-loc .kame-tabs { flex: 1 1 420px; border-bottom-color: rgb(96 0 24 / 28%); }

/* kame-tab vốn dành cho <button>; chỉnh lại để thẻ <a> hiển thị y hệt. */
a.kame-tab {
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
  border-radius: var(--kame-radius-sm) var(--kame-radius-sm) 0 0;
  font-size: 14px;
  text-decoration: none;
}
a.kame-tab:hover { color: var(--kame-orange-700); }
a.kame-tab[aria-selected="true"]:hover { color: var(--kame-white); }

.tt-timkiem { display: flex; flex: 0 1 340px; gap: var(--kame-space-2); align-items: center; }
.tt-timkiem__label {
  position: absolute;
  width: 1px; height: 1px;
  overflow: hidden;
  clip-path: inset(50%);
  white-space: nowrap;
}
.tt-timkiem .kame-input { min-height: 44px; }

.tt-tomtat { margin: 0 0 var(--kame-space-5); font-size: 14px; }

/* --- Thẻ bài viết -------------------------------------------------------- */
.tt-card {
  display: grid;
  grid-template-columns: minmax(0, 200px) minmax(0, 1fr);
  gap: 0;
  align-content: start;
}
.tt-card__media { display: block; padding: var(--kame-space-3); }
.tt-card__body {
  display: grid;
  align-content: start;
  gap: var(--kame-space-2);
  padding: var(--kame-space-4) var(--kame-space-4) var(--kame-space-4) 0;
}
.tt-card__top { display: flex; flex-wrap: wrap; gap: var(--kame-space-2); align-items: center; }
.tt-card__date { font-size: 13px; color: rgb(40 21 15 / 60%); font-weight: 600; }

.tt-card__title { margin: 0; font-size: 18px; font-weight: 800; line-height: 1.3; }
.tt-card__title a { color: var(--kame-maroon-700); text-decoration: none; }
.tt-card__title a:hover { color: var(--kame-orange-700); text-decoration: underline; }

.tt-card__excerpt { margin: 0; font-size: 14px; color: rgb(40 21 15 / 74%); }
.tt-card__more { color: var(--kame-orange-700); font-size: 14px; font-weight: 800; text-decoration: none; }
.tt-card__more:hover { text-decoration: underline; }

/* Thẻ dọc: ảnh nằm trên, dùng cho khối "Tin liên quan" */
.tt-card--doc { grid-template-columns: 1fr; }
.tt-card--doc .tt-card__body { padding: 0 var(--kame-space-4) var(--kame-space-4); }

/* --- Badge giftcode ------------------------------------------------------ */
.tt-badge--gold { color: var(--kame-maroon-900); background: var(--kame-gold); }

/* --- Phân trang ---------------------------------------------------------- */
.tt-phantrang { margin-top: var(--kame-space-6); }
.tt-page--tat { opacity: .35; }

/* --- Trang chi tiết ------------------------------------------------------ */
.tt-bai { max-width: 860px; margin-top: var(--kame-space-5); }
.tt-bai__head { display: grid; gap: var(--kame-space-3); justify-items: start; }
.tt-bai__title { margin: 0; }
.tt-bai__meta { margin: 0; font-size: 14px; }
.tt-bai__anh {
  margin-block: var(--kame-space-5);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  overflow: hidden;
}
.tt-bai__anh .kame-ph { border: 0; border-radius: 0; }
.tt-bai__sapo {
  margin: 0 0 var(--kame-space-5);
  padding-left: var(--kame-space-4);
  border-left: 4px solid var(--kame-orange-600);
  font-size: 17px;
  font-weight: 600;
  color: rgb(40 21 15 / 82%);
}

.tt-body { font-size: 16px; }
.tt-body > :first-child { margin-top: 0; }
.tt-body p { margin: 0 0 var(--kame-space-4); }
.tt-body h2, .tt-body h3, .tt-body h4 {
  margin: var(--kame-space-6) 0 var(--kame-space-3);
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-style: italic;
  font-weight: 900;
  text-transform: uppercase;
}
.tt-body h2 { font-size: 26px; }
.tt-body h3 { font-size: 22px; }
.tt-body h4 { font-size: 18px; }
.tt-body ul, .tt-body ol { margin: 0 0 var(--kame-space-4); padding-left: var(--kame-space-5); }
.tt-body li { margin-bottom: 6px; }
.tt-body a { color: var(--kame-orange-700); font-weight: 700; }
.tt-body strong { color: var(--kame-maroon-700); }
.tt-body blockquote {
  margin: 0 0 var(--kame-space-4);
  padding: var(--kame-space-4);
  border-left: 4px solid var(--kame-gold);
  background: rgb(255 200 61 / 16%);
}
.tt-body img { max-width: 100%; height: auto; border-radius: var(--kame-radius-sm); }
.tt-body table { width: 100%; border-collapse: collapse; margin-bottom: var(--kame-space-4); }
.tt-body th, .tt-body td { padding: 10px 12px; border: 1px solid rgb(96 0 24 / 26%); text-align: left; }
.tt-body th { color: var(--kame-white); background: var(--kame-maroon-700); }

.tt-quaylai { margin: var(--kame-space-6) 0; }
.tt-alert a { color: inherit; }

/* --- Ba mức responsive: 1180 / 900 / 600 -------------------------------- */
@media (max-width: 1180px) {
  .tt-card { grid-template-columns: minmax(0, 170px) minmax(0, 1fr); }
  .tt-card__title { font-size: 17px; }
}

@media (max-width: 900px) {
  .tt-noibat { grid-template-columns: 1fr; }
  .tt-noibat__body { border-left: 0; border-top: 2px dashed rgb(96 0 24 / 30%); }
  .tt-loc { flex-direction: column; align-items: stretch; }
  /* Xếp dọc rồi thì phải ép width 100%: riêng align-items:stretch không kéo
     hàng tab (vốn là hộp cuộn ngang) co lại, làm cả trang bị cuộn ngang. */
  .tt-loc > * { width: 100%; }
  .tt-timkiem { flex: 1 1 auto; }
  .tt-bai { max-width: none; }
}

@media (max-width: 600px) {
  .tt-card { grid-template-columns: 1fr; }
  .tt-card__body { padding: 0 var(--kame-space-4) var(--kame-space-4); }
  .tt-card__media { padding: var(--kame-space-3) var(--kame-space-3) 0; }
  .tt-dautrang { padding-block: var(--kame-space-5); }
  .tt-body h2 { font-size: 22px; }
  .tt-body h3 { font-size: 19px; }
}
</style>
<?php
kame_auth_modals();
kame_footer();
