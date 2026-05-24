/**
 * Public endpoints — no auth required.
 * Simulates anonymous visitors browsing the site.
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, WEB_URL, stages, thresholds } from '../config.js';

export const options = {
  stages,
  thresholds,
};

export default function () {
  // Health check
  {
    const res = http.get(`${BASE_URL}/api/v1/health`);
    check(res, {
      'health: status 200': (r) => r.status === 200,
    });
  }

  sleep(0.3);

  // Programs list
  {
    const res = http.get(`${BASE_URL}/api/v1/programs`);
    check(res, {
      'programs list: status 200': (r) => r.status === 200,
      'programs list: is array': (r) => {
        try {
          return Array.isArray(JSON.parse(r.body));
        } catch {
          return false;
        }
      },
    });

    // Fetch first program detail if list returned items
    try {
      const programs = JSON.parse(res.body);
      if (Array.isArray(programs) && programs.length > 0) {
        const id = programs[0].id;
        const detail = http.get(`${BASE_URL}/api/v1/programs/${id}`);
        check(detail, {
          'program detail: status 200': (r) => r.status === 200,
        });
      }
    } catch (_) {
      /* ignore parse errors under load */
    }
  }

  sleep(0.5);

  // Next.js public pages (web server)
  const pages = ['/', '/programs', '/about', '/contacts'];
  const page = pages[Math.floor(Math.random() * pages.length)];
  {
    const res = http.get(`${WEB_URL}${page}`);
    check(res, {
      [`page ${page}: status 200`]: (r) => r.status === 200,
    });
  }

  sleep(1);
}
