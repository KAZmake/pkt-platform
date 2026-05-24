import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import * as AuthSession from 'expo-auth-session';
import * as SecureStore from 'expo-secure-store';
import { config } from '../config';
import type { AuthSession as AuthSessionType, TokenSet } from './types';

interface AuthContextValue {
  session: AuthSessionType | null;
  isLoading: boolean;
  login: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue>({
  session: null,
  isLoading: true,
  login: async () => {},
  logout: async () => {},
});

const TOKEN_KEY = 'pkt_auth_tokens';

const discovery: AuthSession.DiscoveryDocument = {
  authorizationEndpoint: `${config.keycloakUrl}/realms/${config.keycloakRealm}/protocol/openid-connect/auth`,
  tokenEndpoint: `${config.keycloakUrl}/realms/${config.keycloakRealm}/protocol/openid-connect/token`,
  revocationEndpoint: `${config.keycloakUrl}/realms/${config.keycloakRealm}/protocol/openid-connect/logout`,
  userInfoEndpoint: `${config.keycloakUrl}/realms/${config.keycloakRealm}/protocol/openid-connect/userinfo`,
};

function parseJwtPayload(token: string): Record<string, unknown> {
  try {
    const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(Buffer.from(base64, 'base64').toString('utf8'));
  } catch {
    return {};
  }
}

function extractRole(payload: Record<string, unknown>): string {
  const realmAccess = payload.realm_access as { roles?: string[] } | undefined;
  const roles = realmAccess?.roles ?? [];
  if (roles.includes('admin')) return 'admin';
  if (roles.includes('expert')) return 'expert';
  if (roles.includes('employee')) return 'employee';
  if (roles.includes('borrower')) return 'borrower';
  return 'public';
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [session, setSession] = useState<AuthSessionType | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const redirectUri = AuthSession.makeRedirectUri({ scheme: 'pkt', path: 'auth' });

  const [request, response, promptAsync] = AuthSession.useAuthRequest(
    {
      clientId: config.keycloakClientId,
      scopes: ['openid', 'profile', 'email'],
      redirectUri,
      usePKCE: true,
    },
    discovery,
  );

  // Restore session from secure store on mount
  useEffect(() => {
    (async () => {
      try {
        const stored = await SecureStore.getItemAsync(TOKEN_KEY);
        if (stored) {
          const tokens: TokenSet = JSON.parse(stored);
          if (tokens.expiresAt > Date.now()) {
            const payload = parseJwtPayload(tokens.accessToken);
            setSession({
              tokens,
              userId: String(payload.sub ?? ''),
              email: String(payload.email ?? ''),
              role: extractRole(payload),
              firstName: String(payload.given_name ?? ''),
              lastName: String(payload.family_name ?? ''),
            });
          } else {
            await SecureStore.deleteItemAsync(TOKEN_KEY);
          }
        }
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  // Handle auth response from browser
  useEffect(() => {
    if (response?.type !== 'success') return;
    (async () => {
      const { code } = response.params;
      const tokenRes = await AuthSession.exchangeCodeAsync(
        {
          clientId: config.keycloakClientId,
          redirectUri,
          code,
          extraParams: { code_verifier: request?.codeVerifier ?? '' },
        },
        discovery,
      );
      const tokens: TokenSet = {
        accessToken: tokenRes.accessToken ?? '',
        refreshToken: tokenRes.refreshToken ?? '',
        idToken: tokenRes.idToken ?? '',
        expiresAt: Date.now() + (tokenRes.expiresIn ?? 3600) * 1000,
      };
      await SecureStore.setItemAsync(TOKEN_KEY, JSON.stringify(tokens));
      const payload = parseJwtPayload(tokens.accessToken);
      setSession({
        tokens,
        userId: String(payload.sub ?? ''),
        email: String(payload.email ?? ''),
        role: extractRole(payload),
        firstName: String(payload.given_name ?? ''),
        lastName: String(payload.family_name ?? ''),
      });
    })();
  }, [response]);

  const login = useCallback(async () => {
    await promptAsync();
  }, [promptAsync]);

  const logout = useCallback(async () => {
    await SecureStore.deleteItemAsync(TOKEN_KEY);
    setSession(null);
  }, []);

  return (
    <AuthContext.Provider value={{ session, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
