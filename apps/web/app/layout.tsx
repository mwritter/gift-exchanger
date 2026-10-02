import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";

import { Providers } from "./providers";
import "./globals.css";
import { CurrentUserProvider } from "@/components/CurrentUser/CurrentUserProvider";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "GiftExchanger",
  description: "A way to give gifts for friends, families, and small groups.",
};

export default function RootLayout({ children, login }: LayoutProps<"/">) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className={`${geistSans.className} min-h-full flex flex-col`}>
        <Providers>
          <CurrentUserProvider>
            {children}
            {login}
          </CurrentUserProvider>
        </Providers>
      </body>
    </html>
  );
}
