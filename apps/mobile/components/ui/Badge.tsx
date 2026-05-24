import { Text, View } from 'react-native';
import type { ApplicationStatus } from '@pkt/shared';

const STATUS_CONFIG: Record<ApplicationStatus, { label: string; bg: string; text: string }> = {
  received: { label: 'Получена', bg: 'bg-gray-100', text: 'text-gray-600' },
  primary_scoring: { label: 'Скоринг', bg: 'bg-blue-100', text: 'text-blue-700' },
  security_check: { label: 'Безопасность', bg: 'bg-purple-100', text: 'text-purple-700' },
  collateral_expertise: { label: 'Экспертиза', bg: 'bg-amber-100', text: 'text-amber-700' },
  legal_check: { label: 'Юр. проверка', bg: 'bg-pink-100', text: 'text-pink-700' },
  credit_analysis: { label: 'Анализ', bg: 'bg-cyan-100', text: 'text-cyan-700' },
  credit_committee: { label: 'Комитет', bg: 'bg-orange-100', text: 'text-orange-700' },
  approved: { label: 'Одобрено', bg: 'bg-green-100', text: 'text-green-700' },
  rejected: { label: 'Отклонено', bg: 'bg-red-100', text: 'text-red-700' },
  revision: { label: 'Доработка', bg: 'bg-yellow-100', text: 'text-yellow-700' },
  documentation: { label: 'Документы', bg: 'bg-indigo-100', text: 'text-indigo-700' },
  issued: { label: 'Выдан', bg: 'bg-emerald-100', text: 'text-emerald-700' },
};

export function StatusBadge({ status }: { status: ApplicationStatus }) {
  const cfg = STATUS_CONFIG[status];
  return (
    <View className={`px-2.5 py-1 rounded-full self-start ${cfg.bg}`}>
      <Text className={`text-xs font-medium ${cfg.text}`}>{cfg.label}</Text>
    </View>
  );
}

export function TagBadge({ label, className }: { label: string; className?: string }) {
  return (
    <View className={`px-2.5 py-1 rounded-full self-start bg-brand-green-50 ${className ?? ''}`}>
      <Text className="text-xs font-medium text-brand-green-600">{label}</Text>
    </View>
  );
}
