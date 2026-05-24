import { describe, it, expect } from 'vitest';
import { calcAnnuity, calcDifferentiated } from '@/components/calculator/loan-calculator';

describe('calcAnnuity', () => {
  it('returns correct number of rows', () => {
    const rows = calcAnnuity(1_000_000, 12, 12);
    expect(rows).toHaveLength(12);
  });

  it('rows have sequential period numbers', () => {
    const rows = calcAnnuity(1_000_000, 6, 12);
    rows.forEach((r, i) => expect(r.period).toBe(i + 1));
  });

  it('last row balance is ~0', () => {
    const rows = calcAnnuity(1_000_000, 12, 12);
    expect(rows[rows.length - 1].balance).toBeCloseTo(0, 0);
  });

  it('sum of principal payments equals original amount', () => {
    const amount = 2_000_000;
    const rows = calcAnnuity(amount, 24, 10);
    const totalPrincipal = rows.reduce((s, r) => s + r.principal, 0);
    expect(totalPrincipal).toBeCloseTo(amount, 0);
  });

  it('annuity payment is constant (all total values equal)', () => {
    const rows = calcAnnuity(500_000, 12, 12);
    const firstTotal = rows[0].total;
    rows.forEach((r) => expect(r.total).toBeCloseTo(firstTotal, 1));
  });

  it('handles zero rate — equal principal, zero interest', () => {
    const rows = calcAnnuity(120_000, 12, 0);
    rows.forEach((r) => {
      expect(r.interest).toBe(0);
      expect(r.principal).toBeCloseTo(10_000, 1);
    });
  });

  it('interest decreases over time with positive rate', () => {
    const rows = calcAnnuity(1_000_000, 12, 12);
    for (let i = 1; i < rows.length; i++) {
      expect(rows[i].interest).toBeLessThan(rows[i - 1].interest);
    }
  });

  it('principal increases over time with positive rate', () => {
    const rows = calcAnnuity(1_000_000, 12, 12);
    for (let i = 1; i < rows.length; i++) {
      expect(rows[i].principal).toBeGreaterThan(rows[i - 1].principal);
    }
  });
});

describe('calcDifferentiated', () => {
  it('returns correct number of rows', () => {
    const rows = calcDifferentiated(1_000_000, 12, 12);
    expect(rows).toHaveLength(12);
  });

  it('rows have sequential period numbers', () => {
    const rows = calcDifferentiated(1_000_000, 6, 12);
    rows.forEach((r, i) => expect(r.period).toBe(i + 1));
  });

  it('last row balance is ~0', () => {
    const rows = calcDifferentiated(1_000_000, 12, 12);
    expect(rows[rows.length - 1].balance).toBeCloseTo(0, 0);
  });

  it('principal payment is constant', () => {
    const amount = 1_200_000;
    const rows = calcDifferentiated(amount, 12, 12);
    const expectedPrincipal = amount / 12;
    rows.forEach((r) => expect(r.principal).toBeCloseTo(expectedPrincipal, 1));
  });

  it('interest decreases over time', () => {
    const rows = calcDifferentiated(1_000_000, 12, 12);
    for (let i = 1; i < rows.length; i++) {
      expect(rows[i].interest).toBeLessThan(rows[i - 1].interest);
    }
  });

  it('total payment decreases over time', () => {
    const rows = calcDifferentiated(1_000_000, 12, 12);
    for (let i = 1; i < rows.length; i++) {
      expect(rows[i].total).toBeLessThan(rows[i - 1].total);
    }
  });

  it('sum of principal payments equals original amount', () => {
    const amount = 3_000_000;
    const rows = calcDifferentiated(amount, 24, 15);
    const totalPrincipal = rows.reduce((s, r) => s + r.principal, 0);
    expect(totalPrincipal).toBeCloseTo(amount, 0);
  });

  it('total = principal + interest for each row', () => {
    const rows = calcDifferentiated(1_000_000, 12, 10);
    rows.forEach((r) => {
      expect(r.total).toBeCloseTo(r.principal + r.interest, 6);
    });
  });
});
