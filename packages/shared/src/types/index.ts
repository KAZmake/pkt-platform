// ─── User Roles ──────────────────────────────────────────────────────────────

export type UserRole = 'public' | 'borrower' | 'employee' | 'expert' | 'admin';

export interface User {
  id: string;
  email: string;
  role: UserRole;
  firstName: string;
  lastName: string;
  createdAt: string;
}

// ─── API Response ────────────────────────────────────────────────────────────

export interface ApiResponse<T> {
  data: T;
  meta?: {
    page?: number;
    perPage?: number;
    total?: number;
  };
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

// ─── Application (Заявка) ────────────────────────────────────────────────────

export type ApplicationStatus =
  | 'received'
  | 'primary_scoring'
  | 'security_check'
  | 'collateral_expertise'
  | 'legal_check'
  | 'credit_analysis'
  | 'credit_committee'
  | 'approved'
  | 'rejected'
  | 'revision'
  | 'documentation'
  | 'issued';

export interface Application {
  id: string;
  borrowerId: string;
  status: ApplicationStatus;
  programId: string;
  amount: number;
  termMonths: number;
  createdAt: string;
  updatedAt: string;
}

// ─── Loan Program (Программа кредитования) ───────────────────────────────────

export interface LoanProgram {
  id: string;
  name: string;
  nameKz?: string;
  nameEn?: string;
  rate: number;
  minAmount: number;
  maxAmount: number;
  minTermMonths: number;
  maxTermMonths: number;
  activityTypes: ActivityType[];
  isActive: boolean;
}

export type ActivityType = 'crop_farming' | 'livestock' | 'mixed';

// ─── Collateral (Залог) ──────────────────────────────────────────────────────

export interface Collateral {
  id: string;
  type: CollateralType;
  description: string;
  estimatedValue: number;
  insuranceExpiry?: string;
  lastInventoryDate?: string;
  applicationIds: string[];
}

export type CollateralType = 'land' | 'equipment' | 'livestock' | 'real_estate' | 'other';

// ─── Locale ──────────────────────────────────────────────────────────────────

export type Locale = 'ru' | 'kz' | 'en';
