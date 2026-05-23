import type { LoanProgram } from '@pkt/shared';
import { apiClient } from './client';

export interface ProgramsResponse {
  programs: LoanProgram[];
}

/** Fetch all active loan programs (public). */
export function getPrograms(token?: string) {
  return apiClient.get<LoanProgram[]>('/programs', { token });
}

/** Fetch a single program by ID. */
export function getProgram(id: string, token?: string) {
  return apiClient.get<LoanProgram>(`/programs/${id}`, { token });
}
