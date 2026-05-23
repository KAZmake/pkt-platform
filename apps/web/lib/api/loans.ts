import { syncClient } from './client';

export interface Loan {
  id: string;
  one_c_id: string;
  borrower_id: string;
  amount: number;
  issued_at: string;
  status: 'active' | 'overdue' | 'closed';
  program_name: string;
  term_months: number;
  rate: number;
}

export interface ScheduleItem {
  id: string;
  loan_id: string;
  period: number;
  due_date: string;
  principal: number;
  interest: number;
  total: number;
  is_paid: boolean;
}

export interface LoanDebt {
  id: string;
  loan_id: string;
  debt_type: string;
  amount: number;
  due_date: string;
}

/** Fetch loans for a borrower from sync-svc. */
export function getBorrowerLoans(borrowerId: string, token?: string) {
  return syncClient.get<Loan[]>(`/borrowers/${borrowerId}/loans`, { token });
}

/** Fetch a single loan. */
export function getLoan(id: string, token?: string) {
  return syncClient.get<Loan>(`/loans/${id}`, { token });
}

/** Fetch payment schedule for a loan. */
export function getLoanSchedule(id: string, token?: string) {
  return syncClient.get<ScheduleItem[]>(`/loans/${id}/schedule`, { token });
}

/** Fetch debts for a loan. */
export function getLoanDebts(id: string, token?: string) {
  return syncClient.get<LoanDebt[]>(`/loans/${id}/debts`, { token });
}
