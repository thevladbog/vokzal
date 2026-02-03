// import { invoke } from '@tauri-apps/api/tauri';
import type { Ticket, SaleRequest, Receipt } from '@/types';

// Temporary mock for Tauri invoke
const invoke = async <T>(command: string, args?: any): Promise<T> => {
  console.warn(`Tauri invoke called: ${command}`, args);
  throw new Error('Tauri not available - this is a web build');
};

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1';
const AGENT_URL = import.meta.env.VITE_AGENT_URL || 'http://localhost:8081';

export const posService = {
  // Продажа билета
  sellTicket: async (request: SaleRequest): Promise<Ticket> => {
    const token = localStorage.getItem('access_token');
    if (!token) throw new Error('Не авторизован');

    const ticket = await invoke<Ticket>('sell_ticket', {
      apiUrl: API_URL,
      token,
      request,
    });

    return ticket;
  },

  // Возврат билета
  returnTicket: async (ticketId: string): Promise<Ticket> => {
    const token = localStorage.getItem('access_token');
    if (!token) throw new Error('Не авторизован');

    const ticket = await invoke<Ticket>('return_ticket', {
      apiUrl: API_URL,
      token,
      ticketId,
    });

    return ticket;
  },

  // Печать билета
  printTicket: async (ticketData: any): Promise<boolean> => {
    const success = await invoke<boolean>('print_ticket', {
      agentUrl: AGENT_URL,
      ticketData,
    });

    return success;
  },

  // Печать чека
  printReceipt: async (receiptData: any): Promise<Receipt> => {
    const receipt = await invoke<Receipt>('print_receipt', {
      agentUrl: AGENT_URL,
      receiptData,
    });

    return receipt;
  },

  // Открыть экран покупателя
  openCustomerDisplay: async (): Promise<void> => {
    await invoke('open_customer_display');
  },
};
