import { create } from 'zustand';
import type { User } from '@/types';

/**
 * Auth state is kept in memory only (no zustand persist) for security:
 * access tokens are not written to localStorage/sessionStorage.
 *
 * Session survives page refresh via a separate mechanism: the refresh token
 * is stored in sessionStorage (see constants/auth.ts), and authService.restoreSession()
 * is called on app init (App.tsx). That calls POST /auth/refresh, then setAuth()
 * to repopulate this store. So users stay logged in across refreshes within
 * the same tab.
 */

interface AuthStore {
  user: User | null;
  accessToken: string | null;
  setAuth: (user: User | null, accessToken: string | null) => void;
  setAccessToken: (token: string | null) => void;
  setUser: (user: User | null) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthStore>()((set) => ({
  user: null,
  accessToken: null,

  setAuth: (user, accessToken) => set({ user, accessToken }),

  setAccessToken: (token) => set({ accessToken: token }),

  setUser: (user) => set({ user }),

  clearAuth: () => set({ user: null, accessToken: null }),
}));
