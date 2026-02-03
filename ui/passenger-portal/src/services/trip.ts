import api from './api';
import { Trip, TripSearchRequest } from '@/types';

export const tripService = {
  async search(params: TripSearchRequest): Promise<Trip[]> {
    const response = await api.post<{ data: Trip[] }>('/schedule/trips/search', params);
    return response.data.data;
  },

  async getById(id: string): Promise<Trip> {
    const response = await api.get<{ data: Trip }>(`/schedule/trips/${id}`);
    return response.data.data;
  },
};
