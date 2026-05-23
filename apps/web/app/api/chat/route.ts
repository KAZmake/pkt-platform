/**
 * 3.6.2 — Proxy to assistant-svc (:8083).
 * Keeps ASSISTANT_URL as a server-side secret (no NEXT_PUBLIC_),
 * so the internal service address is never exposed to the browser.
 */
import { type NextRequest, NextResponse } from 'next/server';

const ASSISTANT_URL = process.env.ASSISTANT_URL ?? 'http://localhost:8083';

export async function POST(req: NextRequest) {
  let body: unknown;
  try {
    body = await req.json();
  } catch {
    return NextResponse.json({ error: 'Invalid JSON' }, { status: 400 });
  }

  try {
    const upstream = await fetch(`${ASSISTANT_URL}/api/v1/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });

    if (!upstream.ok) {
      const text = await upstream.text();
      return NextResponse.json(
        { error: 'Upstream error', detail: text },
        { status: upstream.status },
      );
    }

    const data = await upstream.json();
    return NextResponse.json(data);
  } catch {
    return NextResponse.json({ error: 'Assistant service unavailable' }, { status: 503 });
  }
}
