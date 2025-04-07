import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    domains: ['i.ytimg.com'],
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://back:9999/api/:path*' // Proxy to Backend
      }
    ]
  }
};

export default nextConfig;
