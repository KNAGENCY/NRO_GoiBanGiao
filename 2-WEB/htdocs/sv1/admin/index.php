<?php
// ---------------------------------------------------------------------------
//  /sv1/admin/  -> cua truoc cua khu quan tri
//
//  Truoc day file nay day BAT KY AI da dang nhap vao thang panel, khong hoi
//  co phai admin khong. Nguoi choi thuong go /sv1/admin/ la bi nem vao trang
//  quan tri roi moi bi tu choi - vua kho hieu voi ho, vua lo ra la co khu
//  quan tri o day.
//
//  Gio hoi thang CSDL: dung admin thi vao panel, khong phai thi ve trang chu.
// ---------------------------------------------------------------------------
session_start();

if (empty($_SESSION['account_id'])) {
    header('Location: /sv1/trangchu');
    exit;
}

require_once $_SERVER['DOCUMENT_ROOT'] . '/sv1/config.php';

// Doc quyen tu CSDL chu KHONG tin vao bat cu gi luu trong phien dang nhap.
// Ly do: neu sau nay ai do go quyen admin cua mot tai khoan, nguoi do phai
// mat quyen NGAY, khong doi den luc ho dang xuat roi dang nhap lai.
$stmt = $pdo->prepare('SELECT is_admin FROM account WHERE id = ? LIMIT 1');
$stmt->execute([(int)$_SESSION['account_id']]);
$laAdmin = ((int)$stmt->fetchColumn() === 1);

header('Location: ' . ($laAdmin ? '/sv1/admin/panel' : '/sv1/trangchu'));
exit;
