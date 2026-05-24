import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ProgramCard } from '@/components/programs/program-card';
import type { LoanProgram } from '@pkt/shared';

vi.mock('next/link', () => ({
  default: ({
    href,
    children,
    className,
  }: {
    href: string;
    children: React.ReactNode;
    className?: string;
  }) => (
    <a href={href} className={className}>
      {children}
    </a>
  ),
}));

const baseProgram: LoanProgram = {
  id: 'prog-1',
  name: 'Агро 2025',
  rate: 7.5,
  minAmount: 100_000,
  maxAmount: 5_000_000,
  minTermMonths: 6,
  maxTermMonths: 60,
  activityTypes: [],
  isActive: true,
};

describe('ProgramCard', () => {
  it('renders program name', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.getByText('Агро 2025')).toBeInTheDocument();
  });

  it('renders rate', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.getByText('7.5%')).toBeInTheDocument();
  });

  it('renders term range', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.getByText('6–60')).toBeInTheDocument();
  });

  it('renders min and max amount labels', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.getByText('Мин. сумма')).toBeInTheDocument();
    expect(screen.getByText('Макс. сумма')).toBeInTheDocument();
  });

  it('links to program detail page', () => {
    render(<ProgramCard program={baseProgram} />);
    const link = screen.getByRole('link', { name: /подробнее/i });
    expect(link).toHaveAttribute('href', '/programs/prog-1');
  });

  it('renders "Подробнее" button', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.getByRole('button', { name: 'Подробнее' })).toBeInTheDocument();
  });

  it('does not render compare button when onCompare not provided', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.queryByTitle(/сравнени/i)).not.toBeInTheDocument();
  });

  it('renders compare button when onCompare is provided', () => {
    render(<ProgramCard program={baseProgram} onCompare={() => {}} />);
    expect(screen.getByTitle('Добавить к сравнению')).toBeInTheDocument();
  });

  it('compare button shows remove title when selected', () => {
    render(<ProgramCard program={baseProgram} selected onCompare={() => {}} />);
    expect(screen.getByTitle('Убрать из сравнения')).toBeInTheDocument();
  });

  it('calls onCompare with program id when compare button clicked', () => {
    const onCompare = vi.fn();
    render(<ProgramCard program={baseProgram} onCompare={onCompare} />);
    screen.getByTitle('Добавить к сравнению').click();
    expect(onCompare).toHaveBeenCalledWith('prog-1');
  });

  it('renders activity type badge when activityTypes present', () => {
    const p = { ...baseProgram, activityTypes: ['crop_farming' as const] };
    render(<ProgramCard program={p} />);
    expect(screen.getByText('Растениеводство')).toBeInTheDocument();
  });

  it('does not render activity type badge when activityTypes is empty', () => {
    render(<ProgramCard program={baseProgram} />);
    expect(screen.queryByText('Растениеводство')).not.toBeInTheDocument();
  });
});
