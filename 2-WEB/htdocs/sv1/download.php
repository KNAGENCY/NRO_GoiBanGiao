<?php
declare(strict_types=1);

$allowedFiles = [
    'NRO.apk' => 'application/vnd.android.package-archive',
    // Dung .zip chu khong .rar: Windows mo duoc zip san, khong phai cai them
    // WinRAR. Nguoi choi tai ve la bam chuot phai > Extract All duoc ngay.
    'NROKAME.zip' => 'application/zip',
    'IPA.ipa' => 'application/octet-stream',
];

$requested = basename((string)($_GET['file'] ?? ''));

if (!isset($allowedFiles[$requested])) {
    http_response_code(404);
    exit('File not found');
}

$filePath = __DIR__ . '/down/' . $requested;

if (!is_file($filePath) || !is_readable($filePath)) {
    http_response_code(404);
    exit('File not found');
}

$size = filesize($filePath);
$mtime = filemtime($filePath);
$ctime = filectime($filePath);
$etag = '"' . sha1($requested . '|' . $mtime . '|' . $ctime . '|' . $size) . '"';

while (ob_get_level() > 0) {
    ob_end_clean();
}

header('Content-Type: ' . $allowedFiles[$requested]);
header('Content-Length: ' . $size);
header('Content-Disposition: attachment; filename="' . addcslashes($requested, '"\\') . '"');
header('Content-Transfer-Encoding: binary');
header('Cache-Control: no-store, no-cache, must-revalidate, max-age=0');
header('Cache-Control: post-check=0, pre-check=0', false);
header('Pragma: no-cache');
header('Expires: 0');
header('Last-Modified: ' . gmdate('D, d M Y H:i:s', $mtime) . ' GMT');
header('ETag: ' . $etag);
header('X-Download-Version: ' . $mtime . '-' . $ctime . '-' . $size);

readfile($filePath);
exit;
