import { API_PREFIX } from '@pkt/shared';
import { config } from '../config';
import type { ApiResponse } from '@pkt/shared';

async function apiFetch<T>(
  url: string,
  options: RequestInit & { token?: string } = {},
): Promise<T> {
  const { token, ...rest } = options;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(rest.headers as Record<string, string> | undefined),
  };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(url, { ...rest, headers });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: `HTTP ${res.status}` }));
    throw new Error(String(err.message ?? `HTTP ${res.status}`));
  }
  return res.json() as Promise<T>;
}

export const apiClient = {
  get: <T>(path: string, token?: string) =>
    apiFetch<ApiResponse<T>>(`${config.apiUrl}${API_PREFIX}${path}`, { token }),

  post: <T>(path: string, body: unknown, token?: string) =>
    apiFetch<ApiResponse<T>>(`${config.apiUrl}${API_PREFIX}${path}`, {
      method: 'POST',
      body: JSON.stringify(body),
      token,
    }),

  patch: <T>(path: string, body: unknown, token?: string) =>
    apiFetch<ApiResponse<T>>(`${config.apiUrl}${API_PREFIX}${path}`, {
      method: 'PATCH',
      body: JSON.stringify(body),
      token,
    }),
};
