import type { NextConfig } from "next";

const API_URL = `${process.env.API_URL ?? "http://localhost:8080"}`

const nextConfig: NextConfig = {
transpilePackages: ["@giftexchanger/types"],

async rewrites() {
  return [
    {
      source: "/api/:path*",
      destination: `${API_URL}/api/:path*`,
    },
  ];
}
};

export default nextConfig;
