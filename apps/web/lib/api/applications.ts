import type { Application } from '@pkt/shared';
import { apiClient } from './client';

export interface CreateApplicationPayload {
  programId: string;
  amount: number;
  termMonths: number;
  purpose?: string;
}

/** List applications for the authenticated user. */
export function getMyApplications(token: string) {
  return apiClient.get<Application[]>('/applications/my', { token });
}

/** Get a single application. */
export function getApplication(id: string, token: string) {
  return apiClient.get<Application>(`/applications/${id}`, { token });
}

/** Create a new application. */
export function createApplication(payload: CreateApplicationPayload, token: string) {
  return apiClient.post<Application>('/applications', payload, { token });
}
