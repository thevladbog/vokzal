import { apiClient } from './api';
import type {
  UserAdmin,
  ListUsersResponse,
  CreateUserRequest,
  UpdateUserRequest,
} from '@/types';

const baseUrl = '/users';

export const usersService = {
  list: async (params?: {
    role?: string;
    station_id?: string;
    page?: number;
    limit?: number;
  }): Promise<ListUsersResponse> => {
    const response = await apiClient.get<{ success: boolean; data: ListUsersResponse }>(baseUrl, {
      params,
    });
    return response.data.data;
  },

  get: async (id: string): Promise<UserAdmin> => {
    const response = await apiClient.get<{ success: boolean; data: UserAdmin }>(`${baseUrl}/${id}`);
    return response.data.data;
  },

  create: async (data: CreateUserRequest): Promise<UserAdmin> => {
    const response = await apiClient.post<{ success: boolean; data: UserAdmin }>(baseUrl, data);
    return response.data.data;
  },

  update: async (id: string, data: UpdateUserRequest): Promise<UserAdmin> => {
    const response = await apiClient.put<{ success: boolean; data: UserAdmin }>(`${baseUrl}/${id}`, data);
    return response.data.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`${baseUrl}/${id}`);
  },
};
