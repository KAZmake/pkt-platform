// Shared configuration for all k6 scenarios

export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
export const WEB_URL = __ENV.WEB_URL || 'http://localhost:3100';

// Ramp: 0 → 50 → 200 → 0 over ~5 min
export const stages = [
  { duration: '30s', target: 10 }, // warm-up
  { duration: '1m', target: 50 }, // low load
  { duration: '2m', target: 200 }, // peak load
  { duration: '30s', target: 0 }, // ramp-down
];

export const thresholds = {
  http_req_duration: ['p(95)<500', 'p(99)<1000'],
  http_req_failed: ['rate<0.01'],
};

export function defaultHeaders(token) {
  const h = { 'Content-Type': 'application/json' };
  if (token) h['Authorization'] = `Bearer ${token}`;
  return h;
}
