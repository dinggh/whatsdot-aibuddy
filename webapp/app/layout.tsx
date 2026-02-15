import type { Metadata } from "next";
import { Outfit } from "next/font/google";
import type { ReactNode } from "react";

import "./globals.css";

const outfit = Outfit({
  subsets: ["latin"],
  display: "swap"
});

export const metadata: Metadata = {
  title: "微点辅导助手",
  description: "Pencil 设计还原 - shadcn + Tailwind"
};

export default function RootLayout({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="zh-CN">
      <body className={outfit.className}>{children}</body>
    </html>
  );
}
