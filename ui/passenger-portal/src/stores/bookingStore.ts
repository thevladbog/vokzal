import { create } from 'zustand';
import { BookingState } from '@/types';

export const useBookingStore = create<BookingState>((set, get) => ({
  selectedTrip: null,
  passengers: [],
  paymentMethod: 'card',
  contactPhone: '',
  contactEmail: '',

  selectTrip: (trip) => set({ selectedTrip: trip }),

  addPassenger: (passenger) => {
    const { passengers } = get();
    set({ passengers: [...passengers, passenger] });
  },

  removePassenger: (index) => {
    const { passengers } = get();
    set({ passengers: passengers.filter((_, i) => i !== index) });
  },

  updatePassenger: (index, passenger) => {
    const { passengers } = get();
    const updated = [...passengers];
    updated[index] = passenger;
    set({ passengers: updated });
  },

  setPaymentMethod: (method) => set({ paymentMethod: method }),

  setContactPhone: (phone) => set({ contactPhone: phone }),

  setContactEmail: (email) => set({ contactEmail: email }),

  reset: () =>
    set({
      selectedTrip: null,
      passengers: [],
      paymentMethod: 'card',
      contactPhone: '',
      contactEmail: '',
    }),
}));
