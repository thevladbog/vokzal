import api from './api';
import { Ticket, TicketSaleRequest, Payment } from '@/types';

export const ticketService = {
  async sell(request: TicketSaleRequest): Promise<{ tickets: Ticket[]; payment: Payment }> {
    const response = await api.post<{ data: { tickets: Ticket[]; payment: Payment } }>(
      '/ticket/sell',
      request
    );
    return response.data.data;
  },

  async getById(id: string): Promise<Ticket> {
    const response = await api.get<{ data: Ticket }>(`/ticket/tickets/${id}`);
    return response.data.data;
  },

  async getByNumber(number: string): Promise<Ticket> {
    const response = await api.get<{ data: Ticket }>(`/ticket/tickets/number/${number}`);
    return response.data.data;
  },

  async getUserTickets(): Promise<Ticket[]> {
    const response = await api.get<{ data: Ticket[] }>('/ticket/my-tickets');
    return response.data.data;
  },

  async requestReturn(ticketId: string): Promise<Ticket> {
    const response = await api.post<{ data: Ticket }>(`/ticket/return/${ticketId}`);
    return response.data.data;
  },
};
