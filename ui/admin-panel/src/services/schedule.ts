import { apiClient } from './api';
import type { Schedule, Trip, Route } from '@/types';

export const scheduleService = {
  // Schedules
  getSchedules: async (params?: { route_id?: string; is_active?: boolean }) => {
    const response = await apiClient.get<Schedule[]>('/schedule/schedules', { params });
    return response.data;
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
  getRoutes: async (params?: { is_active?: boolean }) => {
    const response = await apiClient.get<Route[]>('/schedule/routes', { params });
    return response.data;
  },

  getRoute: async (id: string) => {
    const response = await apiClient.get<Route>(`/schedule/routes/${id}`);
    return response.data;
  },

  createRoute: async (data: Omit<Route, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await apiClient.post<Route>('/schedule/routes', data);
    return response.data;
  },

  updateRoute: async (id: string, data: Partial<Route>) => {
    const response = await apiClient.put<Route>(`/schedule/routes/${id}`, data);
    return response.data;
  },

  // Trips
  getTrips: async (params?: { 
    schedule_id?: string; 
    status?: string; 
    from_date?: string; 
    to_date?: string;
  }) => {
    const response = await apiClient.get<Trip[]>('/schedule/trips', { params });
    return response.data;
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
