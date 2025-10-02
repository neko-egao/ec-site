import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Header from "@/components/Header";
import { AuthProvider } from "@/contexts/AuthContext"; // ★ AuthProviderをインポート

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "ECサイト",
  description: "GoとNext.jsで作るECサイト",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className={`${inter.className} bg-gray-50`}>
        {/* ★★★ AuthProviderで全体を囲む ★★★ */}
        <AuthProvider>
          <Header />
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}