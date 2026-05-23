import { auth } from '@/auth';
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

// Routes that require authentication (any role)
const PROTECTED_PREFIXES = ['/cabinet', '/expertise', '/analytics'];

// Routes that require specific roles
const ROLE_ROUTES: { prefix: string; roles: string[] }[] = [
  { prefix: '/expertise', roles: ['employee', 'expert', 'admin'] },
  { prefix: '/analytics', roles: ['employee', 'expert', 'admin'] },
  { prefix: '/cabinet', roles: ['borrower', 'employee', 'expert', 'admin'] },
];

export default auth((req: NextRequest & { auth: { role?: string } | null }) => {
  const { pathname } = req.nextUrl;
  const session = req.auth;

  const isProtected = PROTECTED_PREFIXES.some((p) => pathname.startsWith(p));
  if (!isProtected) return NextResponse.next();

  // Not authenticated → redirect to login
  if (!session) {
    const loginUrl = new URL('/login', req.url);
    loginUrl.searchParams.set('callbackUrl', pathname);
    return NextResponse.redirect(loginUrl);
  }

  // Token refresh error → force re-login
  if ((session as Record<string, unknown>).error === 'RefreshAccessTokenError') {
    const loginUrl = new URL('/login', req.url);
    return NextResponse.redirect(loginUrl);
  }

  const role = session.role ?? 'public';
  const routeRule = ROLE_ROUTES.find((r) => pathname.startsWith(r.prefix));
  if (routeRule && !routeRule.roles.includes(role)) {
    return NextResponse.redirect(new URL('/403', req.url));
  }

  return NextResponse.next();
});

export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static, _next/image, favicon.ico
     * - public assets
     * - api/auth (NextAuth internal routes)
     */
    '/((?!_next/static|_next/image|favicon.ico|api/auth|.*\\.(?:svg|png|jpg|ico|webp)).*)',
  ],
};
