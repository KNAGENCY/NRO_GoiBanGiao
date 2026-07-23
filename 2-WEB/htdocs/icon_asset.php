<?php
declare(strict_types=1);

// Thu muc icon tren may ca nhan va VPS. Thu muc ton tai dau tien se duoc dung.
// Tren VPS game root la C:\zz1 (khong phai C:\Users\Administrator\Music\zz1) -
// day la ly do truoc kia icon 503 tren VPS. Liet ke x4 truoc, lui ve x3/x2.
$iconDirectories = [
    'C:/zz1/data/icon/x4',
    'C:/zz1/data/icon/x3',
    'C:/zz1/data/icon/x2',
    'C:/Users/Admin/Music/zz1/data/icon/x4',
    'C:/Users/Admin/Music/zz1/data/icon/x3',
    'C:/Users/Administrator/Music/zz1/data/icon/x4',
];

$iconRoot = null;
foreach ($iconDirectories as $directory) {
    if (is_dir($directory)) {
        $iconRoot = realpath($directory);
        break;
    }
}

if ($iconRoot === null) {
    http_response_code(503);
    header('Content-Type: text/plain; charset=UTF-8');
    echo 'Thu muc icon chua san sang.';
    exit;
}

$requestedFile = (string)($_GET['file'] ?? '');

// Icon game chi cho phep ten file don, ngan truy cap ra ngoai thu muc du lieu.
if ($requestedFile === '' || basename($requestedFile) !== $requestedFile || !preg_match('/^[A-Za-z0-9_-]+\.(?:png|jpe?g|gif|webp)$/i', $requestedFile)) {
    http_response_code(404);
    exit;
}

$iconFile = $iconRoot . DIRECTORY_SEPARATOR . $requestedFile;
if (!is_file($iconFile) || !is_readable($iconFile)) {
    http_response_code(404);
    exit;
}

$extension = strtolower((string)pathinfo($iconFile, PATHINFO_EXTENSION));
$contentTypes = [
    'png' => 'image/png',
    'jpg' => 'image/jpeg',
    'jpeg' => 'image/jpeg',
    'gif' => 'image/gif',
    'webp' => 'image/webp',
];

$size = filesize($iconFile);
$modified = filemtime($iconFile);
$etag = '"' . sha1($requestedFile . '|' . $size . '|' . $modified) . '"';

header('Content-Type: ' . ($contentTypes[$extension] ?? 'application/octet-stream'));
header('Content-Length: ' . $size);
header('Cache-Control: public, max-age=86400, immutable');
header('ETag: ' . $etag);
header('Last-Modified: ' . gmdate('D, d M Y H:i:s', $modified) . ' GMT');

$ifNoneMatch = trim((string)($_SERVER['HTTP_IF_NONE_MATCH'] ?? ''));
$ifModifiedSince = strtotime((string)($_SERVER['HTTP_IF_MODIFIED_SINCE'] ?? ''));
if ($ifNoneMatch === $etag || ($ifModifiedSince !== false && $ifModifiedSince >= $modified)) {
    http_response_code(304);
    exit;
}

readfile($iconFile);
