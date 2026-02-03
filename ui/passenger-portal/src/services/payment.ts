import api from './api';
import { Payment } from '@/types';

export const paymentService = {
  async getById(id: string): Promise<Payment> {
    const response = await api.get<{ data: Payment }>(`/payment/payments/${id}`);
    return response.data.data;
  },

  async checkStatus(id: string): Promise<Payment> {
    const response = await api.get<{ data: Payment }>(`/payment/payments/${id}/status`);
    return response.data.data;
  },
};
