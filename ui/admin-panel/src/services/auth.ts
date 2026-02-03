import { z } from 'zod';
import { apiClient } from './api';
import { useAuthStore } from '@/stores/authStore';
import type { AuthResponse, User } from '@/types';

const UserRoleSchema = z.enum([
  'admin',
  'dispatcher',
  'cashier',
  'controller',
  'accountant',
]);

const UserSchema = z.object({
  id: z.string(),
  username: z.string(),
  fio: z.string().optional(),
  full_name: z.string().optional(),
  role: UserRoleSchema,
  station_id: z.string().optional(),
});

export const AuthResponseSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
  expires_in: z.number(),
  user: UserSchema,
});

export type AuthResponseParsed = z.infer<typeof AuthResponseSchema>;

export const authService = {
  login: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await apiClient.post<{ success: boolean; data?: AuthResponse }>(
      '/auth/login',
      { username, password }
    );

    if (response.data?.success !== true) {
      const message =
        (response.data as { error?: string })?.error ?? 'Login failed: invalid response';
      throw new Error(message);
    }

    const payload = response.data.data ?? response.data;
    const result = AuthResponseSchema.safeParse(payload);

    if (!result.success) {
      const details = result.error.issues
        .map((e) => `${e.path.join('.')}: ${e.message}`)
        .join('; ');
      throw new Error(`Invalid auth response: ${details}`);
    }

    const data = result.data;
    useAuthStore.getState().setAuth(data.user, data.access_token);
    return data;
  },

  logout: async (): Promise<void> => {
    try {
      await apiClient.post('/auth/logout', {}, { withCredentials: true });
    } finally {
      useAuthStore.getState().clearAuth();
    }
  },

  getCurrentUser: (): User | null => {
    return useAuthStore.getState().user;
  },

  isAuthenticated: (): boolean => {
    return !!useAuthStore.getState().accessToken;
  },
};
