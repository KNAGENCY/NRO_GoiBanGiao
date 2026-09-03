//go:build race

package game

// raceEnabled bao cho test biet dang chay duoi `go test -race`. Mot so test DO
// THOI GIAN (rothoigian_test.go) vo nghia duoi race: bo do race lam PBKDF2 cham
// ~10 lan va nhieu tang len hang mili-giay, vuot nguong 1ms cua phep do — do
// KHONG PHAI lo hong bao mat, chi la dung cu do bi nhieu. Cac test do t.Skip khi
// raceEnabled. Cap file race_bat.go / race_tat.go la cach chuan cua Go de doc co
// -race (khong co API runtime nao cho biet).
const raceEnabled = true
