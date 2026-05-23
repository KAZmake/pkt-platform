'use client';

import { useSession } from 'next-auth/react';

/**
 * Returns the access token from the current session.
 * Use this in Client Components when calling the API.
 */
export function useSessionToken() {
  const { data: session, status } = useSession();
  return {
    token: (session as { accessToken?: string })?.accessToken,
    role: (session as { role?: string })?.role ?? 'public',
    status,
    isLoading: status === 'loading',
    isAuthenticated: status === 'authenticated',
  };
}
