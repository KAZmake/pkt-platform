/**
 * Test-only endpoint: mints a NextAuth session cookie for a given role.
 * Gated by E2E_TEST=true — returns 404 in production.
 */
import { NextResponse } from 'next/server';
import { encode } from '@auth/core/jwt';

export const dynamic = 'force-dynamic';

const COOKIE_NAME = 'authjs.session-token';

export async function GET(req: Request) {
  if (process.env.E2E_TEST !== 'true') {
    return NextResponse.json({ error: 'not found' }, { status: 404 });
  }

  const { searchParams } = new URL(req.url);
  const role = searchParams.get('role') ?? 'borrower';
  const secret = process.env.NEXTAUTH_SECRET ?? 'change-me-in-production';

  const token = await encode({
    token: {
      sub: 'e2e-test-user',
      name: 'E2E User',
      email: 'e2e@test.pkt',
      picture: null,
      role,
      accessToken: 'e2e-access-token',
      refreshToken: 'e2e-refresh-token',
      expiresAt: Math.floor(Date.now() / 1000) + 7200,
    },
    secret,
    maxAge: 7200,
    salt: COOKIE_NAME,
  });

  const res = NextResponse.json({ ok: true, role });
  res.cookies.set(COOKIE_NAME, token, {
    httpOnly: true,
    path: '/',
    sameSite: 'lax',
    maxAge: 7200,
  });
  return res;
}
