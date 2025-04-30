import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    domains: ['i.ytimg.com'],
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://103.54.18.35/api/:path*' // Proxy to Backend
      }
    ]
  }
};

export default nextConfig;
