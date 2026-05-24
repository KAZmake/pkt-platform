export interface TokenSet {
  accessToken: string;
  refreshToken: string;
  idToken: string;
  expiresAt: number;
}

export interface AuthSession {
  tokens: TokenSet;
  userId: string;
  email: string;
  role: string;
  firstName: string;
  lastName: string;
}
