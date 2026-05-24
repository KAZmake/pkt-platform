import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Badge, ApplicationStatusBadge } from '@/components/ui/badge';
import type { ApplicationStatus } from '@pkt/shared';

describe('Badge', () => {
  it('renders children', () => {
    render(<Badge>Active</Badge>);
    expect(screen.getByText('Active')).toBeInTheDocument();
  });

  it('applies default variant', () => {
    const { container } = render(<Badge>Default</Badge>);
    expect(container.firstChild).toHaveClass('bg-gray-100');
  });

  it('applies green variant', () => {
    const { container } = render(<Badge variant="green">Green</Badge>);
    expect(container.firstChild).toHaveClass('bg-brand-green-100');
  });

  it('applies red variant', () => {
    const { container } = render(<Badge variant="red">Red</Badge>);
    expect(container.firstChild).toHaveClass('bg-red-100');
  });

  it('applies blue variant', () => {
    const { container } = render(<Badge variant="blue">Blue</Badge>);
    expect(container.firstChild).toHaveClass('bg-blue-100');
  });

  it('applies yellow variant', () => {
    const { container } = render(<Badge variant="yellow">Yellow</Badge>);
    expect(container.firstChild).toHaveClass('bg-yellow-100');
  });

  it('applies purple variant', () => {
    const { container } = render(<Badge variant="purple">Purple</Badge>);
    expect(container.firstChild).toHaveClass('bg-purple-100');
  });

  it('merges custom className', () => {
    const { container } = render(<Badge className="my-class">X</Badge>);
    expect(container.firstChild).toHaveClass('my-class');
  });
});

describe('ApplicationStatusBadge', () => {
  const cases: Array<{ status: ApplicationStatus; label: string }> = [
    { status: 'received', label: 'Получена' },
    { status: 'primary_scoring', label: 'Первичный скоринг' },
    { status: 'security_check', label: 'Проверка СБ' },
    { status: 'collateral_expertise', label: 'Экспертиза залога' },
    { status: 'legal_check', label: 'Юридическая проверка' },
    { status: 'credit_analysis', label: 'Кредитный анализ' },
    { status: 'credit_committee', label: 'Кредитный комитет' },
    { status: 'approved', label: 'Одобрена' },
    { status: 'rejected', label: 'Отказ' },
    { status: 'revision', label: 'На доработке' },
    { status: 'documentation', label: 'Оформление' },
    { status: 'issued', label: 'Выдан' },
  ];

  for (const { status, label } of cases) {
    it(`renders label for status "${status}"`, () => {
      render(<ApplicationStatusBadge status={status} />);
      expect(screen.getByText(label)).toBeInTheDocument();
    });
  }

  it('approved renders green variant', () => {
    const { container } = render(<ApplicationStatusBadge status="approved" />);
    expect(container.firstChild).toHaveClass('bg-brand-green-100');
  });

  it('rejected renders red variant', () => {
    const { container } = render(<ApplicationStatusBadge status="rejected" />);
    expect(container.firstChild).toHaveClass('bg-red-100');
  });

  it('revision renders yellow variant', () => {
    const { container } = render(<ApplicationStatusBadge status="revision" />);
    expect(container.firstChild).toHaveClass('bg-yellow-100');
  });
});
