// Package khonoidung chi lam MOT viec: NHUNG cac file YAML cua kho noi dung vao
// file thuc thi bang //go:embed.
//
// VI SAO co goi rieng nay thay vi de YAML nam trong internal/noidung/: chi thi
// //go:embed KHONG voi ra ngoai thu muc cua goi duoc (khong co "../"). Ma thu
// muc YAML thi phai nam o `server/kho-noi-dung/` — ngay canh go.mod, cho chu du
// an mo ra thay ngay — chu khong chon sau trong internal/. Goi nay la cai cau
// mot dong giua hai rang buoc do.
//
// VI SAO NHUNG (embed) chu khong doc tu dia luc khoi dong:
//
//  1. May chu KHONG BAO GIO khoi dong duoc voi kho noi dung thieu hay lac cho.
//     Doc tu dia nghia la co mot kieu hong moi: chay dung binary nhung sai thu
//     muc lam viec -> khong tim thay file -> chet luc khoi dong, hoac te hon la
//     nap nham mot ban cu con sot lai tren VPS.
//  2. Ban dang chay LUON khop voi ma nguon da bien dich — tuc `git log` cua thu
//     muc nay la nhat ky day du cua "ai doi so gi, luc nao" (thiet ke §1.2).
//
// Ban NHUNG la ban KHOI DONG. Nap nong (task 2/10) se doc tu DIA de lay ban moi
// hon — do la ly do noidung.Nap nhan mot fs.FS bat ky chu khong bam cung vao
// embed.FS nay.
package khonoidung

import "embed"

// FS chua toan bo file .yaml cua kho noi dung, duong dan phang (vi du "chieu.yaml").
//
//go:embed *.yaml
var FS embed.FS
