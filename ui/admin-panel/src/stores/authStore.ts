import { create } from 'zustand';
import type { User } from '@/types';

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
