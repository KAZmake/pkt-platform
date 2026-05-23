import { GrafanaFrame } from '@/components/analytics/grafana-frame';
import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Аналитика' };

/**
 * GRAFANA_URL — внутренняя переменная (не NEXT_PUBLIC_), т.к. Grafana
 * доступна только внутри сети. Аутентификация через Keycloak OIDC SSO:
 * пользователь уже залогинен через Keycloak, Grafana подхватывает сессию.
 *
 * Настройка в Grafana:
 *   [auth.generic_oauth] enabled = true
 *   auth_url = https://keycloak.pkt.kz/realms/pkt/protocol/openid-connect/auth
 *   token_url = …/token   client_id = grafana   role_attribute_path = …
 */
const GRAFANA_URL = process.env.GRAFANA_URL ?? '';

const DASHBOARDS = [
  {
    id: 'portfolio',
    label: 'Портфель займов',
    path: '/d/portfolio/main?orgId=1&kiosk=tv',
  },
  {
    id: 'borrowers',
    label: 'Заёмщики',
    path: '/d/borrowers/overview?orgId=1&kiosk=tv',
  },
  {
    id: 'applications',
    label: 'Заявки и воронка',
    path: '/d/applications/funnel?orgId=1&kiosk=tv',
  },
  {
    id: 'overdue',
    label: 'Просроченные займы',
    path: '/d/overdue/main?orgId=1&kiosk=tv',
  },
];

export default function AnalyticsPage({ searchParams }: { searchParams: { dash?: string } }) {
  const activeDash = searchParams.dash ?? DASHBOARDS[0].id;
  const dashboard = DASHBOARDS.find((d) => d.id === activeDash) ?? DASHBOARDS[0];
  const grafanaConfigured = Boolean(GRAFANA_URL);

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)]">
      <div className="mb-4 flex items-center justify-between gap-4 shrink-0">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Аналитика</h1>
          <p className="text-sm text-gray-500 mt-0.5">Grafana + Prometheus · SSO через Keycloak</p>
        </div>

        {/* Dashboard tabs */}
        <div className="flex gap-1 bg-gray-100 rounded-lg p-1">
          {DASHBOARDS.map((d) => (
            <a
              key={d.id}
              href={`/analytics?dash=${d.id}`}
              className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors whitespace-nowrap ${
                d.id === activeDash
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              {d.label}
            </a>
          ))}
        </div>
      </div>

      {grafanaConfigured ? (
        <div className="flex-1 min-h-0">
          <GrafanaFrame src={`${GRAFANA_URL}${dashboard.path}`} />
        </div>
      ) : (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center max-w-md">
            <div className="w-16 h-16 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <span className="text-3xl">📊</span>
            </div>
            <h2 className="text-lg font-semibold text-gray-900 mb-2">Grafana не настроена</h2>
            <p className="text-sm text-gray-500 mb-4">
              Укажите переменную окружения{' '}
              <code className="bg-gray-100 px-1.5 py-0.5 rounded text-xs font-mono">
                GRAFANA_URL
              </code>{' '}
              с базовым URL вашего Grafana-инстанса.
            </p>
            <div className="bg-gray-900 text-green-400 rounded-lg px-4 py-3 text-xs font-mono text-left">
              GRAFANA_URL=https://grafana.pkt.kz
            </div>
            <p className="text-xs text-gray-400 mt-3">
              Grafana должна быть настроена с Keycloak OIDC SSO (generic_oauth).
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
