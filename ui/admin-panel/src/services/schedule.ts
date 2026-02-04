import { apiClient } from './api';
import type { Schedule, Trip, Route, Station, Bus, Driver } from '@/types';

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
    const response = await apiClient.patch<{ data: Schedule }>(`/schedule/schedules/${id}`, data);
    const body = response.data as { data?: Schedule };
    return body.data ?? (response.data as unknown as Schedule);
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

  updateTripStatus: async (
    id: string,
    data: { status: string; delay_minutes?: number }
  ): Promise<Trip> => {
    const response = await apiClient.patch<{ data: Trip }>(`/schedule/trips/${id}/status`, data);
    const body = response.data as { data?: Trip };
    return body.data ?? (response.data as unknown as Trip);
  },

  generateTrips: async (scheduleId: string, fromDate: string, toDate: string): Promise<void> => {
    await apiClient.post('/schedule/trips/generate', {
      schedule_id: scheduleId,
      from_date: fromDate,
      to_date: toDate,
    });
  },

  updateTrip: async (
    id: string,
    data: { platform?: string; bus_id?: string; driver_id?: string }
  ): Promise<Trip> => {
    const response = await apiClient.patch<{ data: Trip }>(`/schedule/trips/${id}`, data);
    const body = response.data as { data?: Trip };
    return body.data ?? (response.data as unknown as Trip);
  },

  // Buses
  getBuses: async (params?: { station_id?: string; status?: string }): Promise<Bus[]> => {
    const response = await apiClient.get<{ data: Bus[] }>('/schedule/buses', { params });
    const body = response.data as { data?: Bus[] };
    return body.data ?? (response.data as unknown as Bus[]) ?? [];
  },
  createBus: async (data: Omit<Bus, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Bus }>('/schedule/buses', data);
    const body = response.data as { data?: Bus };
    return body.data ?? (response.data as unknown as Bus);
  },
  updateBus: async (id: string, data: Partial<Bus>) => {
    const response = await apiClient.patch<{ data: Bus }>(`/schedule/buses/${id}`, data);
    const body = response.data as { data?: Bus };
    return body.data ?? (response.data as unknown as Bus);
  },
  deleteBus: async (id: string) => {
    await apiClient.delete(`/schedule/buses/${id}`);
  },

  // Drivers
  getDrivers: async (params?: { station_id?: string }): Promise<Driver[]> => {
    const response = await apiClient.get<{ data: Driver[] }>('/schedule/drivers', { params });
    const body = response.data as { data?: Driver[] };
    return body.data ?? (response.data as unknown as Driver[]) ?? [];
  },
  createDriver: async (data: Omit<Driver, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Driver }>('/schedule/drivers', data);
    const body = response.data as { data?: Driver };
    return body.data ?? (response.data as unknown as Driver);
  },
  updateDriver: async (id: string, data: Partial<Driver>) => {
    const response = await apiClient.patch<{ data: Driver }>(`/schedule/drivers/${id}`, data);
    const body = response.data as { data?: Driver };
    return body.data ?? (response.data as unknown as Driver);
  },
  deleteDriver: async (id: string) => {
    await apiClient.delete(`/schedule/drivers/${id}`);
  },

  getDashboardStats: async (date?: string) => {
    const response = await apiClient.get<{ data: { trips_total: number; trips_scheduled: number; trips_departed: number; trips_cancelled: number; trips_delayed: number; trips_arrived: number } }>(
      '/schedule/stats/dashboard',
      { params: date ? { date } : {} }
    );
    const body = response.data as { data?: { trips_total: number; trips_scheduled: number; trips_departed: number; trips_cancelled: number; trips_delayed: number; trips_arrived: number } };
    return body.data ?? { trips_total: 0, trips_scheduled: 0, trips_departed: 0, trips_cancelled: 0, trips_delayed: 0, trips_arrived: 0 };
  },
};
