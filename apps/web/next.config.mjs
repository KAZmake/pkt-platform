/** @type {import('next').NextConfig} */
const nextConfig = {
  output: process.env.NEXT_BUILD_STANDALONE === 'true' ? 'standalone' : undefined,
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
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
    NEXT_PUBLIC_SYNC_URL: process.env.NEXT_PUBLIC_SYNC_URL,
    NEXT_PUBLIC_ASSISTANT_URL: process.env.NEXT_PUBLIC_ASSISTANT_URL,
  },
};

export default nextConfig;
