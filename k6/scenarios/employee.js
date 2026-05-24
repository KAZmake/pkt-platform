/**
 * Employee / expert flow.
 * Hits expertise-svc, sync-svc, and assistant-svc endpoints.
 *
 * Usage: EMPLOYEE_TOKEN=<jwt> k6 run scenarios/employee.js
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, stages, thresholds, defaultHeaders } from '../config.js';

// Expertise service runs on a separate port (same host in staging)
const EXPERTISE_URL = __ENV.EXPERTISE_URL || BASE_URL.replace(':8080', ':8081');
const SYNC_URL = __ENV.SYNC_URL || BASE_URL.replace(':8080', ':8082');
const ASSISTANT_URL = __ENV.ASSISTANT_URL || BASE_URL.replace(':8080', ':8083');

export const options = {
  stages,
  thresholds,
};

const TOKEN = __ENV.EMPLOYEE_TOKEN || '';

export default function () {
  const headers = defaultHeaders(TOKEN);

  // core-api: admin/employee endpoints
  {
    const res = http.get(`${BASE_URL}/api/v1/applications`, { headers });
    check(res, { 'applications (employee): not 500': (r) => r.status !== 500 });
  }

  sleep(0.3);

  // expertise-svc: queue
  {
    const res = http.get(`${EXPERTISE_URL}/expertise/queue`, { headers });
    check(res, {
      'expertise queue: not 500': (r) => r.status !== 500,
    });

    try {
      const items = JSON.parse(res.body);
      if (Array.isArray(items) && items.length > 0) {
        const id = items[0].id;
        const detail = http.get(`${EXPERTISE_URL}/expertise/applications/${id}`, { headers });
        check(detail, { 'expertise app detail: not 500': (r) => r.status !== 500 });
      }
    } catch (_) {}
  }

  sleep(0.5);

  // sync-svc: health (sync runs on cron, just verify it's alive)
  {
    const res = http.get(`${SYNC_URL}/health`);
    check(res, { 'sync-svc health: status 200': (r) => r.status === 200 });
  }

  sleep(0.3);

  // assistant-svc: send a short chat message
  {
    const payload = JSON.stringify({ message: 'Какова процентная ставка по агропрограмме?' });
    const res = http.post(`${ASSISTANT_URL}/assistant/chat`, payload, { headers });
    check(res, {
      'assistant chat: not 500': (r) => r.status !== 500,
      'assistant chat: has response': (r) => {
        try {
          const body = JSON.parse(r.body);
          return typeof body.reply === 'string' || typeof body.content === 'string';
        } catch {
          return false;
        }
      },
    });
  }

  sleep(2); // AI response is slower — give extra think time
}
