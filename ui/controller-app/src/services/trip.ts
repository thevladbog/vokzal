import api from './api';
import { ApiResponse, Trip, TripStats } from '@/types';

export const tripService = {
  async getById(id: string): Promise<Trip> {
    const response = await api.get<ApiResponse<Trip>>(`/trips/${id}`);
    return response.data.data;
  },

  async getActive(): Promise<Trip[]> {
    const response = await api.get<ApiResponse<Trip[]>>('/trips', {
      params: {
        status: 'boarding',
        limit: 50,
      },
    });
    return response.data.data;
  },

  async getStats(tripId: string): Promise<TripStats> {
    const response = await api.get<ApiResponse<TripStats>>(`/trips/${tripId}/stats`);
    return response.data.data;
  },
};
