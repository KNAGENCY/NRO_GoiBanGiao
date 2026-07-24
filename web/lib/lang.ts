// Dich id lang/chung toc (cot nhan_vat.lang, 0-4) sang ten hien thi.
// Khop thu tu lang trong game goc — chi sua o day neu game doi thu tu id.
const TEN_LANG = ["Làng Lá", "Làng Cát", "Làng Sương Mù", "Làng Đá", "Làng Mây"];

export function tenLang(lang: number): string {
  return TEN_LANG[lang] ?? `Làng #${lang}`;
}
