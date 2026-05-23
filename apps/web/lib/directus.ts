/**
 * Directus REST client for public content: news, FAQ, projects, team.
 * Uses ISR — data is revalidated every 60 seconds.
 */

const DIRECTUS_URL = process.env.DIRECTUS_URL ?? 'http://localhost:8055';
const DIRECTUS_TOKEN = process.env.DIRECTUS_TOKEN ?? '';

// ─── Types ────────────────────────────────────────────────────────────────────

export interface NewsItem {
  id: string;
  title: string;
  slug: string;
  summary: string;
  content: string;
  cover_image?: string;
  published_at: string;
  category?: string;
}

export interface FaqItem {
  id: string;
  question: string;
  answer: string;
  category?: string;
  sort?: number;
}

export interface Project {
  id: string;
  title: string;
  description: string;
  images: string[];
  borrower_name?: string;
  activity_type?: string;
  year?: number;
}

export interface TeamMember {
  id: string;
  name: string;
  position: string;
  photo?: string;
}

// ─── Client ───────────────────────────────────────────────────────────────────

async function directusFetch<T>(path: string): Promise<T> {
  const headers: HeadersInit = { Accept: 'application/json' };
  if (DIRECTUS_TOKEN) headers['Authorization'] = `Bearer ${DIRECTUS_TOKEN}`;

  const res = await fetch(`${DIRECTUS_URL}${path}`, {
    headers,
    next: { revalidate: 60 },
  });

  if (!res.ok) {
    // Fail gracefully — return empty array / null so pages still render
    console.error(`Directus ${path} returned ${res.status}`);
    return (Array.isArray([] as unknown) ? [] : null) as T;
  }

  const json = await res.json();
  return json.data as T;
}

// ─── Queries ──────────────────────────────────────────────────────────────────

export function getNews(limit = 20) {
  return directusFetch<NewsItem[]>(
    `/items/news?filter[status][_eq]=published&sort=-published_at&limit=${limit}&fields=id,title,slug,summary,cover_image,published_at,category`,
  );
}

export function getNewsItem(slug: string) {
  return directusFetch<NewsItem[]>(
    `/items/news?filter[slug][_eq]=${slug}&filter[status][_eq]=published&fields=*&limit=1`,
  ).then((items) => items?.[0] ?? null);
}

export function getFaq() {
  return directusFetch<FaqItem[]>(
    `/items/faq?filter[status][_eq]=published&sort=sort,id&fields=id,question,answer,category`,
  );
}

export function getProjects(limit = 20) {
  return directusFetch<Project[]>(
    `/items/projects?filter[status][_eq]=published&limit=${limit}&fields=id,title,description,images,borrower_name,activity_type,year`,
  );
}

export function getTeam() {
  return directusFetch<TeamMember[]>(
    `/items/team?filter[status][_eq]=published&sort=sort&fields=id,name,position,photo`,
  );
}

/** Returns a full URL for a Directus asset. */
export function assetUrl(fileId: string) {
  return `${DIRECTUS_URL}/assets/${fileId}`;
}
