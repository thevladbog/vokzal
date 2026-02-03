import { apiClient } from './api';
import type { AuditLog } from '@/types';

const baseUrl = '/audit';

const unwrap = <T>(res: { data?: T }): T =>
  (res as { data: T }).data ?? (res as unknown as T);

export const auditService = {
  getLog: async (id: string): Promise<AuditLog> => {
    const response = await apiClient.get<{ data: AuditLog }>(`${baseUrl}/${id}`);
    return unwrap(response.data as { data: AuditLog });
  },
  getLogsByEntity: async (entityType: string, entityId: string): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/entity`, {
      params: { entity_type: entityType, entity_id: entityId },
    });
    return unwrap(response.data as { data?: AuditLog[] }) ?? [];
  },
  getLogsByUser: async (userId: string, limit = 100): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/user`, {
      params: { user_id: userId, limit },
    });
    return unwrap(response.data as { data?: AuditLog[] }) ?? [];
  },
  getLogsByDateRange: async (from: string, to: string): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/date-range`, {
      params: { from, to },
    });
    return unwrap(response.data as { data?: AuditLog[] }) ?? [];
  },
  listLogs: async (limit = 100): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/list`, {
      params: { limit },
    });
    return unwrap(response.data as { data?: AuditLog[] }) ?? [];
  },
};
