import Link from 'next/link';
import type { FarmProperties } from '@/lib/map/types';
import { ACTIVITY_COLORS, LOAN_STATUS_COLORS } from '@/lib/map/config';
import { formatMoney } from '@/lib/utils';

const ACTIVITY_LABELS: Record<string, string> = {
  crop_farming: 'Растениеводство',
  livestock: 'Животноводство',
  mixed: 'Смешанное',
};

const LOAN_STATUS_LABELS: Record<string, string> = {
  active: 'Активный',
  overdue: 'Просрочен',
  closed: 'Закрыт',
};

interface FarmPopupProps {
  farm: FarmProperties;
  isEmployee: boolean;
  onClose: () => void;
}

export function FarmPopup({ farm, isEmployee, onClose }: FarmPopupProps) {
  const activityColor = ACTIVITY_COLORS[farm.activity_type] ?? '#6b7280';
  const loanColor = farm.loan_status ? (LOAN_STATUS_COLORS[farm.loan_status] ?? '#6b7280') : null;

  return (
    <div className="bg-white rounded-xl shadow-xl border border-gray-200 w-64 overflow-hidden">
      {/* Header */}
      <div
        className="px-4 py-3 flex items-center justify-between"
        style={{ borderLeft: `4px solid ${activityColor}` }}
      >
        <div>
          <p className="font-semibold text-gray-900 text-sm leading-tight">{farm.name}</p>
          <p className="text-xs text-gray-400 mt-0.5">{farm.district}</p>
        </div>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-600 text-lg leading-none"
          aria-label="Закрыть"
        >
          ×
        </button>
      </div>

      {/* Details */}
      <div className="px-4 py-3 space-y-2 border-t border-gray-100">
        <Row
          label="Направление"
          value={ACTIVITY_LABELS[farm.activity_type] ?? farm.activity_type}
        />
        {farm.area_ha && <Row label="Площадь" value={`${farm.area_ha} га`} />}

        {/* Employee-only data */}
        {isEmployee && farm.loan_status && (
          <>
            <div className="border-t border-dashed border-gray-100 pt-2 mt-2">
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">
                Кредит
              </p>
              <div className="flex items-center gap-2 mb-1">
                <span
                  className="h-2 w-2 rounded-full shrink-0"
                  style={{ background: loanColor ?? '#6b7280' }}
                />
                <span className="text-sm text-gray-700">
                  {LOAN_STATUS_LABELS[farm.loan_status] ?? farm.loan_status}
                </span>
              </div>
              {farm.loan_amount && <Row label="Сумма" value={formatMoney(farm.loan_amount)} />}
            </div>
            {farm.borrower_id && (
              <Link
                href={`/expertise?borrower_id=${farm.borrower_id}`}
                className="block mt-3 text-center bg-brand-green text-white text-xs font-medium py-2 rounded-lg hover:bg-brand-green-dark transition-colors"
              >
                Открыть досье →
              </Link>
            )}
          </>
        )}
      </div>
    </div>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-2">
      <span className="text-xs text-gray-400">{label}</span>
      <span className="text-xs font-medium text-gray-700">{value}</span>
    </div>
  );
}
