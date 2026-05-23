import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import type { ApplicationStatus } from '@pkt/shared';

const badgeVariants = cva(
  'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
  {
    variants: {
      variant: {
        default: 'bg-gray-100 text-gray-700',
        green: 'bg-brand-green-100 text-brand-green-700',
        gold: 'bg-brand-gold-100 text-brand-gold-dark',
        red: 'bg-red-100 text-red-700',
        blue: 'bg-blue-100 text-blue-700',
        purple: 'bg-purple-100 text-purple-700',
        yellow: 'bg-yellow-100 text-yellow-700',
      },
    },
    defaultVariants: { variant: 'default' },
  },
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>, VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />;
}

// ─── Application status badge ─────────────────────────────────────────────────

const STATUS_LABELS: Record<ApplicationStatus, string> = {
  received: 'Получена',
  primary_scoring: 'Первичный скоринг',
  security_check: 'Проверка СБ',
  collateral_expertise: 'Экспертиза залога',
  legal_check: 'Юридическая проверка',
  credit_analysis: 'Кредитный анализ',
  credit_committee: 'Кредитный комитет',
  approved: 'Одобрена',
  rejected: 'Отказ',
  revision: 'На доработке',
  documentation: 'Оформление',
  issued: 'Выдан',
};

const STATUS_VARIANT: Record<ApplicationStatus, BadgeProps['variant']> = {
  received: 'default',
  primary_scoring: 'blue',
  security_check: 'purple',
  collateral_expertise: 'gold',
  legal_check: 'purple',
  credit_analysis: 'blue',
  credit_committee: 'gold',
  approved: 'green',
  rejected: 'red',
  revision: 'yellow',
  documentation: 'blue',
  issued: 'green',
};

export function ApplicationStatusBadge({ status }: { status: ApplicationStatus }) {
  return <Badge variant={STATUS_VARIANT[status]}>{STATUS_LABELS[status] ?? status}</Badge>;
}
