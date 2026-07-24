// Package game — phan dong bo du lieu tinh (Ke hoach 3).
//
// Khi vao the gioi, client so 4 phien ban (data/map/skill/item) may chu gui
// trong -28/4 voi cache cua no; lech thi tu xin cap nhat. Client MOI (RMS trong)
// luon xin lai ca 4 du phien ban co khop hay khong, nen may chu PHAI phuc vu du
// bon o lan dau. Day la cho ta thay 3 lop Dragon Ball cu bang 5 lang Naruto.
package game

import (
	_ "embed"
)

// duLieuPart: du lieu sprite nhan vat (dau/than/chan). Day la PLACEHOLDER cua
// ban Dragon Ball cu — nhan vat hien ra bang sprite cu; thay bang part Naruto
// khi duong ong my thuat giao. Client ve nhan vat bang parts[head/body/leg]
// KHONG kiem bien, nen part rong = sap khi ve (Ke hoach 4). File nhung phai phu
// du chi so head/body/leg cua moi nhan vat -> forward nguyen file cu la an toan.
//
//go:embed part.bin
var duLieuPart []byte

// Phien ban du lieu tinh cua ban Naruto. PHAI khac ban Dragon Ball cu
// (vsData=10, vsMap=2, vsSkill=1, vsItem=4) de client khong dung cache cu — nhat
// la vsSkill != 1 de khong dung cache 3 lop cu. 20 la gia tri tuy chon != cu va
// != 0. Bump khi doi noi dung tuong ung o cac ke hoach sau.
const (
	vsData  byte = 20
	vsMap   byte = 20
	vsSkill byte = 20
	vsItem  byte = 20
)

// lenhDuLieu: goi data. Client gui -87 than rong de xin, may chu tra -87 day du.
const lenhDuLieu int8 = -87

// Cac lenh con cua khung -28 lien quan dong bo tinh.
const (
	lenhConVersionGame  byte = 4  // may chu -> client
	lenhConCatalogMap   byte = 6  // client xin map / may chu tra
	lenhConCatalogSkill byte = 7  // client xin skill / may chu tra
	lenhConCatalogItem  byte = 8  // client xin item / may chu tra
	lenhConClientOk     byte = 13 // client bao san sang
)
