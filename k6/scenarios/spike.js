/**
 * Spike test — sudden burst to 200 VU, then drop.
 * Verifies the system recovers gracefully.
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, WEB_URL, thresholds } from '../config.js';

export const options = {
  stages: [
    { duration: '10s', target: 5 }, // baseline
    { duration: '10s', target: 200 }, // spike
    { duration: '1m', target: 200 }, // hold
    { duration: '10s', target: 5 }, // recovery
    { duration: '30s', target: 5 }, // verify recovery
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // looser during spike
    http_req_failed: ['rate<0.05'], // allow up to 5% errors during spike
  },
};

export default function () {
  const res = http.get(`${BASE_URL}/api/v1/programs`);
  check(res, { 'programs: not 500': (r) => r.status !== 500 });
  sleep(0.2);

  const web = http.get(`${WEB_URL}/`);
  check(web, { 'home page: not 500': (r) => r.status !== 500 });
  sleep(0.5);
}
