import type { Metadata, Viewport } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { Providers } from "./providers"

const inter = Inter({ subsets: ["latin"], variable: "--font-sans" })

export const metadata: Metadata = {
  title: "todos",
  description: "シンプルなタスク管理アプリ",
  icons: { icon: "/favicon.svg" },
}

// iOS Safari の自動ズームを防ぐため initial-scale=1 を明示
// maximum-scale=1 は使わない（アクセシビリティ上ユーザーによるズームを妨げるため）
export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="ja" className={inter.variable}>
      <body className="min-h-screen bg-background antialiased">
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
