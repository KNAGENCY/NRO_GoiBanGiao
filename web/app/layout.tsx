import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "NRO — Naruto Online",
  description: "Web chinh thuc cua game Naruto Online: dang ky, dang nhap, nap the, bang xep hang.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="vi">
      <body className="min-h-screen antialiased">{children}</body>
    </html>
  );
}
