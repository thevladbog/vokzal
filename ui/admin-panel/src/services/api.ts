import axios from 'axios';
import { REFRESH_TOKEN_STORAGE_KEY } from '@/constants/auth';
import { useAuthStore } from '@/stores/authStore';
import type { User } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor: add in-memory access token to Authorization header
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().accessToken;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor: on 401, try refresh using token from sessionStorage then retry
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Don't retry if already retried, or if this IS the refresh endpoint itself
    if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/refresh')) {
      originalRequest._retry = true;

      let refreshToken: string | null = null;
      try {
        refreshToken = sessionStorage.getItem(REFRESH_TOKEN_STORAGE_KEY);
      } catch {
        // sessionStorage unavailable
      }

      let refreshed = false;
      if (refreshToken) {
        try {
          const response = await apiClient.post<{
            success?: boolean;
            data?: { access_token: string; refresh_token: string; user: unknown };
          }>('/auth/refresh', { refresh_token: refreshToken });

          const data = response.data?.data;
          if (response.data?.success === true && data?.access_token && data?.user) {
            useAuthStore.getState().setAuth(data.user as User, data.access_token);
            try {
              sessionStorage.setItem(REFRESH_TOKEN_STORAGE_KEY, data.refresh_token);
            } catch {
              // ignore
            }
            originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
            refreshed = true;
            return apiClient(originalRequest);
          }
        } catch {
          // Refresh failed (network or 4xx): fall through to clear and redirect
        }
      }

      // Only clear and redirect when we did not successfully refresh (no token, failed, or bad response)
      if (!refreshed) {
        try {
          sessionStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
        } catch {
          // ignore
        }
        useAuthStore.getState().clearAuth();
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);
