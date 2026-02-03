import api from './api';
import { ApiResponse, Ticket, BoardingRequest, BoardingResponse } from '@/types';

export const ticketService = {
  async getById(id: string): Promise<Ticket> {
    const response = await api.get<ApiResponse<Ticket>>(`/tickets/${id}`);
    return response.data.data;
  },

  async getByQR(qrCode: string): Promise<Ticket> {
    const response = await api.get<ApiResponse<Ticket>>('/tickets/by-qr', {
      params: { qr: qrCode },
    });
    return response.data.data;
  },

  async markBoarding(data: BoardingRequest): Promise<BoardingResponse> {
    const response = await api.post<ApiResponse<BoardingResponse>>(
      '/tickets/boarding',
      data
    );
    return response.data.data;
  },

  async getByTrip(tripId: string): Promise<Ticket[]> {
    const response = await api.get<ApiResponse<Ticket[]>>('/tickets', {
      params: { tripId },
    });
    return response.data.data;
  },
};
