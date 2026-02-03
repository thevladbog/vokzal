import axios from 'axios';
import type { Trip, Payment, AuthResponse, Ticket } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Auth
export const authService = {
  login: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>('/auth/login', {
      username,
      password,
    });

    localStorage.setItem('access_token', response.data.access_token);
    localStorage.setItem('refresh_token', response.data.refresh_token);
    localStorage.setItem('user', JSON.stringify(response.data.user));

    return response.data;
  },

  logout: async (): Promise<void> => {
    const refreshToken = localStorage.getItem('refresh_token');
    if (refreshToken) {
      await apiClient.post('/auth/logout', { refresh_token: refreshToken });
    }
    localStorage.clear();
  },
};

// Trips
export const tripService = {
  getTrips: async (params?: { from_date?: string; to_date?: string; status?: string }) => {
    const response = await apiClient.get<Trip[]>('/schedule/trips', { params });
    return response.data;
  },

  getTrip: async (id: string) => {
    const response = await apiClient.get<Trip>(`/schedule/trips/${id}`);
    return response.data;
  },
};

// Tickets
export const ticketService = {
  getTicket: async (id: string): Promise<Ticket> => {
    const response = await apiClient.get<{ data: Ticket }>(`/tickets/${id}`);
    return response.data.data;
  },
};

// Payment
export const paymentService = {
  initPayment: async (ticketId: string, method: 'card' | 'sbp' | 'cash', amount: number) => {
    const response = await apiClient.post<Payment>('/payment/init', {
      ticket_id: ticketId,
      method,
      amount,
    });
    return response.data;
  },

  checkPaymentStatus: async (paymentId: string) => {
    const response = await apiClient.get<Payment>(`/payment/${paymentId}`);
    return response.data;
  },
};
