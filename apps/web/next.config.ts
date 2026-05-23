import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  // Allow images from MinIO and external sources
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '9000',
        pathname: '/**',
      },
    ],
  },
  // Expose env vars to the client
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
    NEXT_PUBLIC_SYNC_URL: process.env.NEXT_PUBLIC_SYNC_URL,
    NEXT_PUBLIC_ASSISTANT_URL: process.env.NEXT_PUBLIC_ASSISTANT_URL,
  },
};

export default nextConfig;
