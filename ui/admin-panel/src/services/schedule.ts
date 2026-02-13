import { apiClient } from './api';
import type { Schedule, Trip, Route, Station, Bus, Driver } from '@/types';

/**
 * Schedule API returns all successful responses in a wrapper: `{ data: T }`.
 * We unwrap to T so callers receive the entity/list directly.
 */
function unwrap<T>(response: { data: { data?: T } }): T {
  const body = response.data;
  if (body && typeof body === 'object' && 'data' in body) {
    const data = body.data;
    if (data !== null && data !== undefined) {
      return data;
    }
    throw new Error('Unexpected API response: "data" wrapper contains null or undefined value');
  }
  throw new Error('Unexpected API response shape: missing "data" wrapper');
}

export const scheduleService = {
  // Stations
  getStations: async (params?: { city?: string }): Promise<Station[]> => {
    const response = await apiClient.get<{ data: Station[] }>('/stations', { params });
    return unwrap<Station[]>(response) ?? [];
  },
  getStation: async (id: string) => {
    const response = await apiClient.get<{ data: Station }>(`/stations/${id}`);
    return unwrap<Station>(response);
  },
  createStation: async (data: Omit<Station, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Station }>('/stations', data);
    return unwrap<Station>(response);
  },
  updateStation: async (id: string, data: Partial<Station>) => {
    const response = await apiClient.patch<{ data: Station }>(`/stations/${id}`, data);
    return unwrap<Station>(response);
  },
  deleteStation: async (id: string) => {
    await apiClient.delete(`/stations/${id}`);
  },

  // Schedules
  getSchedules: async (params?: { route_id?: string; is_active?: boolean }): Promise<Schedule[]> => {
    const response = await apiClient.get<{ data: Schedule[] }>('/schedules', { params });
    return unwrap<Schedule[]>(response) ?? [];
  },

  getSchedule: async (id: string) => {
    const response = await apiClient.get<{ data: Schedule }>(`/schedules/${id}`);
    return unwrap<Schedule>(response);
  },

  createSchedule: async (data: Omit<Schedule, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Schedule }>('/schedules', data);
    return unwrap<Schedule>(response);
  },

  updateSchedule: async (id: string, data: Partial<Schedule>) => {
    const response = await apiClient.patch<{ data: Schedule }>(`/schedules/${id}`, data);
    return unwrap<Schedule>(response);
  },

  deleteSchedule: async (id: string) => {
    await apiClient.delete(`/schedules/${id}`);
  },

  // Routes
  getRoutes: async (params?: { is_active?: boolean }): Promise<Route[]> => {
    const response = await apiClient.get<{ data: Route[] }>('/routes', { params });
    return unwrap<Route[]>(response) ?? [];
  },

  getRoute: async (id: string) => {
    const response = await apiClient.get<{ data: Route }>(`/routes/${id}`);
    return unwrap<Route>(response);
  },

  createRoute: async (data: Omit<Route, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Route }>('/routes', data);
    return unwrap<Route>(response);
  },

  updateRoute: async (id: string, data: Partial<Route>) => {
    const response = await apiClient.patch<{ data: Route }>(`/routes/${id}`, data);
    return unwrap<Route>(response);
  },
  deleteRoute: async (id: string) => {
    await apiClient.delete(`/routes/${id}`);
  },

  // Trips (backend accepts date=YYYY-MM-DD for single day, or from_date/to_date if supported)
  getTrips: async (params?: {
    schedule_id?: string;
    status?: string;
    date?: string;
    from_date?: string;
    to_date?: string;
  }): Promise<Trip[]> => {
    const response = await apiClient.get<{ data: Trip[] }>('/trips', {
      params: params,
    });
    return unwrap<Trip[]>(response) ?? [];
  },

  getTrip: async (id: string) => {
    const response = await apiClient.get<{ data: Trip }>(`/trips/${id}`);
    return unwrap<Trip>(response);
  },

  updateTripStatus: async (
    id: string,
    data: { status: string; delay_minutes?: number }
  ): Promise<Trip> => {
    const response = await apiClient.patch<{ data: Trip }>(`/trips/${id}/status`, data);
    return unwrap<Trip>(response);
  },

  generateTrips: async (scheduleId: string, fromDate: string, toDate: string): Promise<void> => {
    await apiClient.post('/trips/generate', {
      schedule_id: scheduleId,
      from_date: fromDate,
      to_date: toDate,
    });
  },

  updateTrip: async (
    id: string,
    data: { platform?: string; bus_id?: string | null; driver_id?: string | null }
  ): Promise<Trip> => {
    const response = await apiClient.patch<{ data: Trip }>(`/trips/${id}`, data);
    return unwrap<Trip>(response);
  },

  // Buses
  getBuses: async (params?: { station_id?: string; status?: string }): Promise<Bus[]> => {
    const response = await apiClient.get<{ data: Bus[] }>('/buses', { params });
    return unwrap<Bus[]>(response) ?? [];
  },
  createBus: async (data: Omit<Bus, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Bus }>('/buses', data);
    return unwrap<Bus>(response);
  },
  updateBus: async (id: string, data: Partial<Bus>) => {
    const response = await apiClient.patch<{ data: Bus }>(`/buses/${id}`, data);
    return unwrap<Bus>(response);
  },
  deleteBus: async (id: string) => {
    await apiClient.delete(`/buses/${id}`);
  },

  // Drivers
  getDrivers: async (params?: { station_id?: string }): Promise<Driver[]> => {
    const response = await apiClient.get<{ data: Driver[] }>('/drivers', { params });
    return unwrap<Driver[]>(response) ?? [];
  },
  createDriver: async (data: Omit<Driver, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<{ data: Driver }>('/drivers', data);
    return unwrap<Driver>(response);
  },
  updateDriver: async (id: string, data: Partial<Driver>) => {
    const response = await apiClient.patch<{ data: Driver }>(`/drivers/${id}`, data);
    return unwrap<Driver>(response);
  },
  deleteDriver: async (id: string) => {
    await apiClient.delete(`/drivers/${id}`);
  },

  getDashboardStats: async (date?: string) => {
    type Stats = {
      trips_total: number;
      trips_scheduled: number;
      trips_boarding: number;
      trips_departed: number;
      trips_cancelled: number;
      trips_delayed: number;
      trips_arrived: number;
      total_capacity: number;
    };
    const response = await apiClient.get<{ data: Stats }>('/schedule/stats/dashboard', {
      params: date ? { date } : {},
    });
    return unwrap<Stats>(response) ?? {
      trips_total: 0,
      trips_scheduled: 0,
      trips_boarding: 0,
      trips_departed: 0,
      trips_cancelled: 0,
      trips_delayed: 0,
      trips_arrived: 0,
      total_capacity: 0,
    };
  },
};
