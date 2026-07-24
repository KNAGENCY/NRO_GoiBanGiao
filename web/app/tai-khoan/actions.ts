"use server";

import { redirect } from "next/navigation";
import { xoaPhien } from "@/lib/phien";

/** Dang xuat: xoa cookie phien roi ve trang chu. */
export async function dangXuat(): Promise<void> {
  await xoaPhien();
  redirect("/");
}
