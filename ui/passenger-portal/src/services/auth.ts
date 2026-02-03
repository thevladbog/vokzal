import api from './api';
import { AuthResponse, LoginRequest, RegisterRequest, User } from '@/types';

export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await api.post<{ data: AuthResponse }>('/auth/login', credentials);
    return response.data.data;
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await api.post<{ data: AuthResponse }>('/auth/register', data);
    return response.data.data;
  },

  async refresh(refreshToken: string): Promise<AuthResponse> {
    const response = await api.post<{ data: AuthResponse }>('/auth/refresh', {
      refreshToken,
    });
    return response.data.data;
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout');
  },

  async getCurrentUser(): Promise<User> {
    const response = await api.get<{ data: User }>('/auth/me');
    return response.data.data;
  },
};
