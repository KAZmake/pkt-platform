import { describe, it, expect } from 'vitest';
import { cn, formatMoney, formatDate, getInitials } from '@/lib/utils';

describe('cn', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar');
  });

  it('resolves tailwind conflicts (last wins)', () => {
    const result = cn('px-2', 'px-4');
    expect(result).toBe('px-4');
  });

  it('ignores falsy values', () => {
    expect(cn('foo', undefined, false, null, 'bar')).toBe('foo bar');
  });

  it('handles empty input', () => {
    expect(cn()).toBe('');
  });
});

describe('formatMoney', () => {
  it('formats KZT with currency symbol', () => {
    const result = formatMoney(1_000_000);
    expect(result).toContain('1');
    expect(result).toContain('000');
  });

  it('rounds to whole number (no decimals)', () => {
    const result = formatMoney(500_500.99);
    expect(result).not.toContain('.');
    expect(result).not.toContain(',99');
  });

  it('formats zero', () => {
    const result = formatMoney(0);
    expect(result).toContain('0');
  });

  it('uses KZT by default', () => {
    const result = formatMoney(100);
    expect(result).toMatch(/₸|KZT|тг/);
  });
});

describe('formatDate', () => {
  it('formats ISO date as dd.MM.yyyy', () => {
    const result = formatDate('2024-03-15');
    expect(result).toBe('15.03.2024');
  });

  it('formats another date correctly', () => {
    const result = formatDate('2025-01-01');
    expect(result).toBe('01.01.2025');
  });

  it('formats end of year', () => {
    const result = formatDate('2023-12-31');
    expect(result).toBe('31.12.2023');
  });
});

describe('getInitials', () => {
  it('returns first chars of first and last name uppercased', () => {
    expect(getInitials('Иван', 'Петров')).toBe('ИП');
  });

  it('uppercases lowercase input', () => {
    expect(getInitials('алия', 'сейткали')).toBe('АС');
  });

  it('works with single-char names', () => {
    expect(getInitials('A', 'B')).toBe('AB');
  });

  it('handles latin chars', () => {
    expect(getInitials('John', 'Doe')).toBe('JD');
  });
});
