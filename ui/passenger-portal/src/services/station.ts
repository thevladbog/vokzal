import api from './api';
import { Station } from '@/types';

export const stationService = {
  async getAll(): Promise<Station[]> {
    const response = await api.get<{ data: Station[] }>('/schedule/stations');
    return response.data.data;
  },

  async getById(id: string): Promise<Station> {
    const response = await api.get<{ data: Station }>(`/schedule/stations/${id}`);
    return response.data.data;
  },

  async search(query: string): Promise<Station[]> {
    const response = await api.get<{ data: Station[] }>('/schedule/stations/search', {
      params: { q: query },
    });
    return response.data.data;
  },
};
