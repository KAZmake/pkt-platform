import type { Metadata } from 'next';

export const metadata: Metadata = { title: 'Почта и диск' };

/**
 * Deeplinks to Microsoft 365 or Google Workspace.
 * Configure via env:
 *   NEXT_PUBLIC_MAIL_PROVIDER = "m365" | "gw"  (default: m365)
 *   NEXT_PUBLIC_M365_TENANT   = "pkt.kz"       (used for onedrive/sharepoint)
 *   NEXT_PUBLIC_SHAREPOINT_URL = "https://pkt.sharepoint.com/sites/intranet"
 */

const provider = (process.env.NEXT_PUBLIC_MAIL_PROVIDER as 'm365' | 'gw' | undefined) ?? 'm365';

const tenant = process.env.NEXT_PUBLIC_M365_TENANT ?? 'pkt.kz';
const sharepointUrl =
  process.env.NEXT_PUBLIC_SHAREPOINT_URL ?? `https://${tenant.split('.')[0]}.sharepoint.com`;

const M365_SERVICES = [
  {
    id: 'outlook',
    label: 'Почта (Outlook)',
    desc: 'Корпоративная email-почта',
    icon: '📧',
    url: 'https://outlook.office365.com',
    color: 'bg-blue-50 border-blue-200',
    badge: 'Microsoft 365',
  },
  {
    id: 'onedrive',
    label: 'OneDrive',
    desc: 'Личное хранилище файлов',
    icon: '☁️',
    url: `https://onedrive.live.com?authupn=admin@${tenant}`,
    color: 'bg-sky-50 border-sky-200',
    badge: 'Microsoft 365',
  },
  {
    id: 'sharepoint',
    label: 'SharePoint / Интранет',
    desc: 'Общие документы и новости компании',
    icon: '📁',
    url: sharepointUrl,
    color: 'bg-indigo-50 border-indigo-200',
    badge: 'Microsoft 365',
  },
  {
    id: 'teams',
    label: 'Microsoft Teams',
    desc: 'Чаты, совещания, звонки',
    icon: '💬',
    url: 'https://teams.microsoft.com',
    color: 'bg-purple-50 border-purple-200',
    badge: 'Microsoft 365',
  },
];

const GW_SERVICES = [
  {
    id: 'gmail',
    label: 'Почта (Gmail)',
    desc: 'Корпоративная email-почта',
    icon: '📧',
    url: 'https://mail.google.com',
    color: 'bg-red-50 border-red-200',
    badge: 'Google Workspace',
  },
  {
    id: 'drive',
    label: 'Google Drive',
    desc: 'Хранилище файлов и документов',
    icon: '☁️',
    url: 'https://drive.google.com',
    color: 'bg-yellow-50 border-yellow-200',
    badge: 'Google Workspace',
  },
  {
    id: 'docs',
    label: 'Google Docs / Sheets',
    desc: 'Документы, таблицы, презентации',
    icon: '📄',
    url: 'https://docs.google.com',
    color: 'bg-green-50 border-green-200',
    badge: 'Google Workspace',
  },
  {
    id: 'meet',
    label: 'Google Meet',
    desc: 'Видеозвонки и совещания',
    icon: '📹',
    url: 'https://meet.google.com',
    color: 'bg-teal-50 border-teal-200',
    badge: 'Google Workspace',
  },
];

const SERVICES = provider === 'gw' ? GW_SERVICES : M365_SERVICES;
const PROVIDER_LABEL = provider === 'gw' ? 'Google Workspace' : 'Microsoft 365';

export default function MailPage() {
  return (
    <div className="max-w-3xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Почта и диск</h1>
        <p className="text-sm text-gray-500 mt-1">
          Быстрый доступ к корпоративным сервисам · {PROVIDER_LABEL}
        </p>
      </div>

      <div className="grid sm:grid-cols-2 gap-4">
        {SERVICES.map((svc) => (
          <a
            key={svc.id}
            href={svc.url}
            target="_blank"
            rel="noopener noreferrer"
            className={`flex items-start gap-4 rounded-xl border p-4 hover:shadow-card transition-all group ${svc.color}`}
          >
            <span className="text-3xl shrink-0 mt-0.5">{svc.icon}</span>
            <div className="flex-1 min-w-0">
              <p className="font-medium text-gray-900 group-hover:text-brand-green">{svc.label}</p>
              <p className="text-xs text-gray-500 mt-0.5">{svc.desc}</p>
            </div>
            <span className="text-gray-300 group-hover:text-brand-green text-sm shrink-0 mt-0.5">
              ↗
            </span>
          </a>
        ))}
      </div>

      <div className="mt-8 bg-gray-50 rounded-xl border border-gray-200 p-4">
        <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">
          Единый вход (SSO)
        </p>
        <p className="text-sm text-gray-600">
          Аутентификация через Keycloak SAML 2.0. При первом переходе вы будете перенаправлены на
          страницу входа {PROVIDER_LABEL} — используйте корпоративный аккаунт{' '}
          <strong>@{tenant}</strong>.
        </p>
        <p className="text-xs text-gray-400 mt-2">
          Настройка провайдера: переменная{' '}
          <code className="bg-white px-1.5 py-0.5 rounded font-mono">
            NEXT_PUBLIC_MAIL_PROVIDER
          </code>{' '}
          = <strong>m365</strong> | <strong>gw</strong>
        </p>
      </div>
    </div>
  );
}
