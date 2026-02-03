import { create } from 'zustand';
import { SearchState } from '@/types';
import { tripService } from '@/services/trip';
import { format } from 'date-fns';

export const useSearchStore = create<SearchState>((set, get) => ({
  fromStation: null,
  toStation: null,
  date: new Date(),
  trips: [],
  isSearching: false,

  setFromStation: (station) => set({ fromStation: station }),

  setToStation: (station) => set({ toStation: station }),

  setDate: (date) => set({ date }),

  searchTrips: async () => {
    const { fromStation, toStation, date } = get();
    if (!fromStation || !toStation) {
      throw new Error('Необходимо выбрать станции отправления и прибытия');
    }

    set({ isSearching: true });
    try {
      const trips = await tripService.search({
        fromStationId: fromStation.id,
        toStationId: toStation.id,
        date: format(date, 'yyyy-MM-dd'),
      });
      set({ trips, isSearching: false });
    } catch (error) {
      set({ isSearching: false });
      throw error;
    }
  },

  swapStations: () => {
    const { fromStation, toStation } = get();
    set({ fromStation: toStation, toStation: fromStation });
  },
}));
