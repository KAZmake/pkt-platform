/**
 * Borrower flow — authenticated user browsing their cabinet.
 *
 * Usage: pass a valid Keycloak JWT via env var:
 *   BORROWER_TOKEN=<jwt> k6 run scenarios/borrower.js
 *
 * For staging: obtain token via Keycloak Resource Owner Password grant
 * (test accounts only — never use production credentials in scripts).
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, WEB_URL, stages, thresholds, defaultHeaders } from '../config.js';

export const options = {
  stages,
  thresholds: {
    ...thresholds,
    // Cabinet endpoints can be slightly slower (DB joins)
    'http_req_duration{scenario:borrower}': ['p(95)<800'],
  },
};

const TOKEN = __ENV.BORROWER_TOKEN || '';

export default function () {
  const headers = defaultHeaders(TOKEN);

  // GET /me — profile
  {
    const res = http.get(`${BASE_URL}/api/v1/me`, { headers });
    check(res, {
      'me: not 500': (r) => r.status !== 500,
      'me: auth required without token': (r) => (TOKEN ? r.status === 200 : r.status === 401),
    });
  }

  sleep(0.3);

  // GET /applications — loan list
  {
    const res = http.get(`${BASE_URL}/api/v1/applications`, { headers });
    check(res, {
      'applications: not 500': (r) => r.status !== 500,
    });

    // Drill into first application
    try {
      const apps = JSON.parse(res.body);
      if (Array.isArray(apps) && apps.length > 0) {
        const id = apps[0].id;

        const schedule = http.get(`${BASE_URL}/api/v1/applications/${id}/schedule`, { headers });
        check(schedule, { 'schedule: not 500': (r) => r.status !== 500 });

        const docs = http.get(`${BASE_URL}/api/v1/applications/${id}/documents`, { headers });
        check(docs, { 'documents: not 500': (r) => r.status !== 500 });
      }
    } catch (_) {}
  }

  sleep(0.5);

  // GET /notifications
  {
    const res = http.get(`${BASE_URL}/api/v1/notifications`, { headers });
    check(res, { 'notifications: not 500': (r) => r.status !== 500 });

    const unread = http.get(`${BASE_URL}/api/v1/notifications/unread-count`, { headers });
    check(unread, { 'unread-count: not 500': (r) => r.status !== 500 });
  }

  sleep(0.3);

  // GET /tickets — support tickets
  {
    const res = http.get(`${BASE_URL}/api/v1/tickets`, { headers });
    check(res, { 'tickets: not 500': (r) => r.status !== 500 });
  }

  sleep(1);

  // Cabinet page via Next.js (HTML)
  {
    const res = http.get(`${WEB_URL}/cabinet`);
    check(res, { 'cabinet page: status not 500': (r) => r.status !== 500 });
  }

  sleep(1);
}
