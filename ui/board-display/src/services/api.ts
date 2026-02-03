import axios from 'axios';
import type { BoardTrip } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

export const boardService = {
  // Получить публичное табло
  getPublicBoard: async (): Promise<BoardTrip[]> => {
    const response = await apiClient.get<BoardTrip[]>('/board/public');
    return response.data;
  },

  // Получить табло для конкретного перрона
  getPlatformBoard: async (platformId: string): Promise<BoardTrip[]> => {
    const response = await apiClient.get<BoardTrip[]>(`/board/platform/${platformId}`);
    return response.data;
  },

  // Получить статистику
  getStats: async () => {
    const response = await apiClient.get('/board/stats');
    return response.data;
  },
};
