import axios from 'axios';
import { useAuthStore } from '@/stores/authStore';

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

// Response interceptor: on 401, try refresh via HttpOnly cookie then retry
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Refresh token is sent via HttpOnly cookie (withCredentials: true)
        const response = await apiClient.post<{ access_token?: string }>(
          '/auth/refresh',
          {},
          { withCredentials: true }
        );

        const accessToken = response.data?.access_token;
        if (accessToken) {
          useAuthStore.getState().setAccessToken(accessToken);
          originalRequest.headers.Authorization = `Bearer ${accessToken}`;
          return apiClient(originalRequest);
        }
      } catch {
        // Refresh failed or no cookie: clear in-memory state and redirect to login
        useAuthStore.getState().clearAuth();
        window.location.href = '/login';
        return Promise.reject(error);
      }

      useAuthStore.getState().clearAuth();
      window.location.href = '/login';
    }

    return Promise.reject(error);
  }
);
