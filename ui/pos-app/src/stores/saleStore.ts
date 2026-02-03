import { create } from 'zustand';
import type { Ticket } from '@/types';

interface SaleStore {
  currentTicket: Ticket | null;
  saleInProgress: boolean;
  setCurrentTicket: (ticket: Ticket | null) => void;
  setSaleInProgress: (inProgress: boolean) => void;
  clearSale: () => void;
}

export const useSaleStore = create<SaleStore>((set) => ({
  currentTicket: null,
  saleInProgress: false,
  setCurrentTicket: (ticket) => set({ currentTicket: ticket }),
  setSaleInProgress: (inProgress) => set({ saleInProgress: inProgress }),
  clearSale: () => set({ currentTicket: null, saleInProgress: false }),
}));
