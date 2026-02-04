import { apiClient } from './api';
import type { AuditLog } from '@/types';

const baseUrl = '/audit';

/** Extracts .data from a wrapped API response; uses runtime check so unexpected shapes return undefined. */
function unwrap<T>(res: unknown): T | undefined {
  if (res != null && typeof res === 'object' && 'data' in res) {
    const data = (res as { data: T }).data;
    if (data !== undefined) return data;
  }
  return undefined;
}

export const auditService = {
  getLog: async (id: string): Promise<AuditLog> => {
    const response = await apiClient.get<{ data: AuditLog }>(`${baseUrl}/${id}`);
    const data = unwrap<AuditLog>(response.data);
    if (data === undefined) {
      throw new Error('Unexpected audit response: missing data');
    }
    return data;
  },
  getLogsByEntity: async (entityType: string, entityId: string): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/entity`, {
      params: { entity_type: entityType, entity_id: entityId },
    });
    return unwrap<AuditLog[]>(response.data) ?? [];
  },
  getLogsByUser: async (userId: string, limit = 100): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/user`, {
      params: { user_id: userId, limit },
    });
    return unwrap<AuditLog[]>(response.data) ?? [];
  },
  getLogsByDateRange: async (from: string, to: string): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/date-range`, {
      params: { from, to },
    });
    return unwrap<AuditLog[]>(response.data) ?? [];
  },
  listLogs: async (limit = 100): Promise<AuditLog[]> => {
    const response = await apiClient.get<{ data: AuditLog[] }>(`${baseUrl}/list`, {
      params: { limit },
    });
    return unwrap<AuditLog[]>(response.data) ?? [];
  },
};
