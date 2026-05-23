/**
 * Base API client with JWT auto-inject and uniform error handling.
 * Works both in Server Components (reads session server-side) and
 * Client Components (token passed explicitly or via cookie).
 */

export class ApiError extends Error {
  constructor(
    public status: number,
    public body: unknown,
  ) {
    super(`API error ${status}`);
    this.name = 'ApiError';
  }
}

interface FetchOptions extends RequestInit {
  token?: string;
  /** Base URL override (default: NEXT_PUBLIC_API_URL) */
  baseUrl?: string;
}

async function apiFetch<T>(path: string, options: FetchOptions = {}): Promise<T> {
  const { token, baseUrl, ...init } = options;
  const base = baseUrl ?? process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';
  const url = `${base}/api/v1${path}`;

  const headers = new Headers(init.headers);
  headers.set('Content-Type', 'application/json');
  if (token) headers.set('Authorization', `Bearer ${token}`);

  const res = await fetch(url, { ...init, headers });

  if (!res.ok) {
    let body: unknown;
    try {
      body = await res.json();
    } catch {
      body = await res.text();
    }
    throw new ApiError(res.status, body);
  }

  // 204 No Content
  if (res.status === 204) return undefined as T;

  return res.json() as Promise<T>;
}

// ─── Convenience methods ──────────────────────────────────────────────────────

export const apiClient = {
  get: <T>(path: string, options?: FetchOptions) =>
    apiFetch<T>(path, { method: 'GET', ...options }),

  post: <T>(path: string, body: unknown, options?: FetchOptions) =>
    apiFetch<T>(path, {
      method: 'POST',
      body: JSON.stringify(body),
      ...options,
    }),

  patch: <T>(path: string, body: unknown, options?: FetchOptions) =>
    apiFetch<T>(path, {
      method: 'PATCH',
      body: JSON.stringify(body),
      ...options,
    }),

  delete: <T>(path: string, options?: FetchOptions) =>
    apiFetch<T>(path, { method: 'DELETE', ...options }),
};

/** API client pointing at sync-svc (:8082). */
export const syncClient = {
  get: <T>(path: string, options?: FetchOptions) =>
    apiFetch<T>(path, {
      method: 'GET',
      baseUrl: process.env.NEXT_PUBLIC_SYNC_URL ?? 'http://localhost:8082',
      ...options,
    }),
};

/** API client pointing at assistant-svc (:8083). */
export const assistantClient = {
  post: <T>(path: string, body: unknown, options?: FetchOptions) =>
    apiFetch<T>(path, {
      method: 'POST',
      body: JSON.stringify(body),
      baseUrl: process.env.NEXT_PUBLIC_ASSISTANT_URL ?? 'http://localhost:8083',
      ...options,
    }),
};
