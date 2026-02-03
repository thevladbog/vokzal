import { apiClient } from './api';
import type { Ticket } from '@/types';

export const ticketService = {
  getTickets: async (params?: {
    trip_id?: string;
    status?: string;
    from_date?: string;
    to_date?: string;
  }) => {
    const response = await apiClient.get<Ticket[]>('/tickets', { params });
    return response.data;
  },

  getTicket: async (id: string) => {
    const response = await apiClient.get<Ticket>(`/tickets/${id}`);
    return response.data;
  },

  sellTicket: async (data: {
    trip_id: string;
    passenger_fio?: string;
    passenger_phone?: string;
    seat_id?: string;
  }) => {
    const response = await apiClient.post<Ticket>('/tickets/sell', data);
    return response.data;
  },

  returnTicket: async (id: string) => {
    const response = await apiClient.post<Ticket>(`/tickets/${id}/return`);
    return response.data;
  },

  getSalesReport: async (fromDate: string, toDate: string) => {
    const response = await apiClient.get('/tickets/reports/sales', {
      params: { from_date: fromDate, to_date: toDate },
    });
    return response.data;
  },
};
