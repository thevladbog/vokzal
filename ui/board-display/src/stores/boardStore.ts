import { create } from 'zustand';
import type { BoardTrip } from '@/types';

interface BoardStore {
  trips: BoardTrip[];
  setTrips: (trips: BoardTrip[]) => void;
  updateTrip: (trip: BoardTrip) => void;
  addTrip: (trip: BoardTrip) => void;
  removeTrip: (tripId: string) => void;
}

export const useBoardStore = create<BoardStore>((set) => ({
  trips: [],
  
  setTrips: (trips) => set({ trips }),
  
  updateTrip: (updatedTrip) =>
    set((state) => ({
      trips: state.trips.map((trip) =>
        trip.id === updatedTrip.id ? updatedTrip : trip
      ),
    })),
  
  addTrip: (newTrip) =>
    set((state) => ({
      trips: [...state.trips, newTrip],
    })),
  
  removeTrip: (tripId) =>
    set((state) => ({
      trips: state.trips.filter((trip) => trip.id !== tripId),
    })),
}));
