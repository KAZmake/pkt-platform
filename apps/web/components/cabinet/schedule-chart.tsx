'use client';

interface ChartRow {
  month: number;
  principal: number;
  interest: number;
  total: number;
  status: 'paid' | 'upcoming' | 'overdue';
}

export function ScheduleChart({ rows }: { rows: ChartRow[] }) {
  const maxTotal = Math.max(...rows.map((r) => r.total));

  return (
    <div className="flex items-end gap-px h-28 w-full overflow-hidden">
      {rows.map((row) => {
        const principalPct = (row.principal / maxTotal) * 100;
        const interestPct = (row.interest / maxTotal) * 100;
        const isPaid = row.status === 'paid';
        return (
          <div
            key={row.month}
            className="flex-1 flex flex-col justify-end cursor-default"
            title={`Мес. ${row.month}: ${row.total.toLocaleString('ru-KZ')} ₸`}
          >
            <div
              className="rounded-t-sm"
              style={{
                height: `${interestPct}%`,
                backgroundColor: isPaid ? '#c8921a' : '#f5e0b0',
              }}
            />
            <div
              style={{
                height: `${principalPct}%`,
                backgroundColor: isPaid ? '#1a5c36' : '#c8e9d4',
              }}
            />
          </div>
        );
      })}
    </div>
  );
}
