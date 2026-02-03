import { create } from 'zustand';
import { ScanState, Trip, TripStats, Ticket } from '@/types';

export const useScanStore = create<ScanState>((set) => ({
  currentTrip: null,
  stats: null,
  recentScans: [],
  isScanning: false,

  setCurrentTrip: (trip: Trip | null) => {
    set({ currentTrip: trip });
  },

  setStats: (stats: TripStats | null) => {
    set({ stats });
  },

  addRecentScan: (ticket: Ticket) => {
    set((state) => ({
      recentScans: [ticket, ...state.recentScans].slice(0, 20), // Keep last 20 scans
    }));
  },

  clearRecentScans: () => {
    set({ recentScans: [] });
  },

  setIsScanning: (isScanning: boolean) => {
    set({ isScanning });
  },
}));
