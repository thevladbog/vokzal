import { apiClient } from './api';
import type { Schedule, Trip, Route, Station } from '@/types';

export const scheduleService = {
  // Stations
  getStations: async (params?: { city?: string }): Promise<Station[]> => {
    const response = await apiClient.get<{ data: Station[] }>('/schedule/stations', { params });
    const body = response.data as { data?: Station[] };
    return body.data ?? (response.data as unknown as Station[]);
  },
  getStation: async (id: string) => {
    const response = await apiClient.get<{ data: Station }>(`/schedule/stations/${id}`);
    const body = response.data as { data?: Station };
    return body.data ?? (response.data as unknown as Station);
  },
  createStation: async (data: Omit<Station, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Station }>('/schedule/stations', data);
    const body = response.data as { data?: Station };
    return body.data ?? (response.data as unknown as Station);
  },
  updateStation: async (id: string, data: Partial<Station>) => {
    const response = await apiClient.patch<{ data: Station }>(`/schedule/stations/${id}`, data);
    const body = response.data as { data?: Station };
    return body.data ?? (response.data as unknown as Station);
  },
  deleteStation: async (id: string) => {
    await apiClient.delete(`/schedule/stations/${id}`);
  },

  // Schedules
  getSchedules: async (params?: { route_id?: string; is_active?: boolean }): Promise<Schedule[]> => {
    const response = await apiClient.get<Schedule[] | { data: Schedule[] }>('/schedule/schedules', { params });
    const body = response.data;
    return Array.isArray(body) ? body : (body as { data: Schedule[] }).data ?? [];
  },

  getSchedule: async (id: string) => {
    const response = await apiClient.get<Schedule>(`/schedule/schedules/${id}`);
    return response.data;
  },

  createSchedule: async (data: Omit<Schedule, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<Schedule>('/schedule/schedules', data);
    return response.data;
  },

  updateSchedule: async (id: string, data: Partial<Schedule>) => {
    const response = await apiClient.put<Schedule>(`/schedule/schedules/${id}`, data);
    return response.data;
  },

  deleteSchedule: async (id: string) => {
    await apiClient.delete(`/schedule/schedules/${id}`);
  },

  // Routes
  getRoutes: async (params?: { is_active?: boolean }): Promise<Route[]> => {
    const response = await apiClient.get<Route[] | { data: Route[] }>('/schedule/routes', { params });
    const body = response.data;
    return Array.isArray(body) ? body : (body as { data: Route[] }).data ?? [];
  },

  getRoute: async (id: string) => {
    const response = await apiClient.get<{ data: Route }>(`/schedule/routes/${id}`);
    const body = response.data as { data?: Route };
    return body.data ?? (response.data as unknown as Route);
  },

  createRoute: async (data: Omit<Route, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Route }>('/schedule/routes', data);
    const body = response.data as { data?: Route };
    return body.data ?? (response.data as unknown as Route);
  },

  updateRoute: async (id: string, data: Partial<Route>) => {
    const response = await apiClient.patch<{ data: Route }>(`/schedule/routes/${id}`, data);
    const body = response.data as { data?: Route };
    return body.data ?? (response.data as unknown as Route);
  },
  deleteRoute: async (id: string) => {
    await apiClient.delete(`/schedule/routes/${id}`);
  },

  // Trips (backend accepts date=YYYY-MM-DD for single day, or from_date/to_date if supported)
  getTrips: async (params?: {
    schedule_id?: string;
    status?: string;
    date?: string;
    from_date?: string;
    to_date?: string;
  }): Promise<Trip[]> => {
    const response = await apiClient.get<Trip[] | { data: Trip[] }>('/schedule/trips', {
      params: params?.date ? { date: params.date } : params,
    });
    const body = response.data;
    return Array.isArray(body) ? body : (body as { data: Trip[] }).data ?? [];
  },

  getTrip: async (id: string) => {
    const response = await apiClient.get<Trip>(`/schedule/trips/${id}`);
    return response.data;
  },

  updateTrip: async (id: string, data: Partial<Trip>) => {
    const response = await apiClient.put<Trip>(`/schedule/trips/${id}`, data);
    return response.data;
  },

  generateTrips: async (scheduleId: string, fromDate: string, toDate: string) => {
    const response = await apiClient.post(`/schedule/schedules/${scheduleId}/generate-trips`, {
      from_date: fromDate,
      to_date: toDate,
    });
    return response.data;
  },
};
