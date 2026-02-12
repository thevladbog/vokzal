import { z } from 'zod';
import { REFRESH_TOKEN_STORAGE_KEY } from '@/constants/auth';
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
    try {
      sessionStorage.setItem(REFRESH_TOKEN_STORAGE_KEY, data.refresh_token);
    } catch {
      // sessionStorage may be unavailable (private mode, etc.)
    }
    return data;
  },

  logout: async (): Promise<void> => {
    try {
      const refreshToken = sessionStorage.getItem(REFRESH_TOKEN_STORAGE_KEY);
      if (refreshToken) {
        await apiClient.post(
          '/auth/logout',
          {},
          { withCredentials: true, headers: { 'X-Refresh-Token': refreshToken } }
        );
      }
    } finally {
      try {
        sessionStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
      } catch {
        // ignore
      }
      useAuthStore.getState().clearAuth();
    }
  },

  /**
   * Restore session on app load using refresh token from sessionStorage.
   * Call once before rendering protected routes. Returns true if session was restored.
   */
  restoreSession: async (): Promise<boolean> => {
    let refreshToken: string | null = null;
    try {
      refreshToken = sessionStorage.getItem(REFRESH_TOKEN_STORAGE_KEY);
    } catch {
      return false;
    }
    if (!refreshToken) return false;

    try {
      const response = await apiClient.post<{ success: boolean; data?: AuthResponse }>(
        '/auth/refresh',
        { refresh_token: refreshToken },
        { validateStatus: (status) => status < 500 } // Don't throw on 4xx
      );
      if (response.status !== 200 || response.data?.success !== true || !response.data?.data) {
        // Token expired or invalid, clean up
        try {
          sessionStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
        } catch {
          // ignore
        }
        return false;
      }

      const parsed = AuthResponseSchema.safeParse(response.data.data);
      if (!parsed.success) return false;

      const data = parsed.data;
      useAuthStore.getState().setAuth(data.user, data.access_token);
      try {
        sessionStorage.setItem(REFRESH_TOKEN_STORAGE_KEY, data.refresh_token);
      } catch {
        // ignore
      }
      return true;
    } catch {
      try {
        sessionStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
      } catch {
        // ignore
      }
      return false;
    }
  },

  getCurrentUser: (): User | null => {
    return useAuthStore.getState().user;
  },

  isAuthenticated: (): boolean => {
    return !!useAuthStore.getState().accessToken;
  },
};
