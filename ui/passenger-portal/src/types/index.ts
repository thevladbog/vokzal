// API Response Types
export interface ApiResponse<T> {
  data: T;
  message?: string;
  error?: string;
}

// Station
export interface Station {
  id: string;
  name: string;
  code: string;
  address: string;
  city: string;
  region: string;
  latitude?: number;
  longitude?: number;
  active: boolean;
}

// Route
export interface Route {
  id: string;
  code: string;
  name: string;
  fromStationId: string;
  toStationId: string;
  fromStation?: Station;
  toStation?: Station;
  distance: number;
  duration: number;
  active: boolean;
}

// Trip
export interface Trip {
  id: string;
  routeId: string;
  route?: Route;
  scheduleId: string;
  departureTime: string;
  arrivalTime: string;
  busId: string;
  busNumber?: string;
  driverId: string;
  driverName?: string;
  price: number;
  totalSeats: number;
  availableSeats: number;
  status: 'scheduled' | 'boarding' | 'departed' | 'arrived' | 'cancelled';
  platform?: string;
  gate?: string;
}

// Search Request
export interface TripSearchRequest {
  fromStationId: string;
  toStationId: string;
  date: string; // YYYY-MM-DD
}

// Passenger
export interface Passenger {
  lastName: string;
  firstName: string;
  middleName?: string;
  birthDate: string; // YYYY-MM-DD
  documentType: 'passport' | 'birth_certificate' | 'foreign_passport';
  documentSeries?: string;
  documentNumber: string;
  phone?: string;
  email?: string;
  benefitType?: 'none' | 'child' | 'student' | 'pensioner' | 'disabled';
  discount?: number;
}

// Ticket Sale Request
export interface TicketSaleRequest {
  tripId: string;
  passengers: Passenger[];
  paymentMethod: 'card' | 'sbp' | 'cash';
  contactPhone: string;
  contactEmail?: string;
}

// Ticket
export interface Ticket {
  id: string;
  number: string;
  tripId: string;
  trip?: Trip;
  passengerLastName: string;
  passengerFirstName: string;
  passengerMiddleName?: string;
  passengerBirthDate: string;
  passengerDocumentType: string;
  passengerDocumentSeries?: string;
  passengerDocumentNumber: string;
  seatNumber?: string;
  price: number;
  discount: number;
  finalPrice: number;
  status: 'sold' | 'returned' | 'boarded' | 'expired';
  paymentMethod: string;
  soldAt: string;
  soldBy?: string;
  returnedAt?: string;
  returnPenalty?: number;
  boardedAt?: string;
  qrCode?: string;
}

// Payment
export interface Payment {
  id: string;
  ticketId?: string;
  amount: number;
  method: 'card' | 'sbp' | 'cash';
  status: 'pending' | 'completed' | 'failed' | 'refunded';
  providerTransactionId?: string;
  qrCodeUrl?: string;
  paymentUrl?: string;
  createdAt: string;
  completedAt?: string;
}

// User (for registered passengers)
export interface User {
  id: string;
  email: string;
  phone?: string;
  lastName: string;
  firstName: string;
  middleName?: string;
  birthDate?: string;
  role: 'passenger';
  createdAt: string;
}

// Auth
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  phone?: string;
  password: string;
  lastName: string;
  firstName: string;
  middleName?: string;
  birthDate?: string;
}

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: User;
}

// Store States
export interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshAuth: () => Promise<void>;
}

export interface SearchState {
  fromStation: Station | null;
  toStation: Station | null;
  date: Date;
  trips: Trip[];
  isSearching: boolean;
  setFromStation: (station: Station | null) => void;
  setToStation: (station: Station | null) => void;
  setDate: (date: Date) => void;
  searchTrips: () => Promise<void>;
  swapStations: () => void;
}

export interface BookingState {
  selectedTrip: Trip | null;
  passengers: Passenger[];
  paymentMethod: 'card' | 'sbp' | 'cash';
  contactPhone: string;
  contactEmail: string;
  selectTrip: (trip: Trip) => void;
  addPassenger: (passenger: Passenger) => void;
  removePassenger: (index: number) => void;
  updatePassenger: (index: number, passenger: Passenger) => void;
  setPaymentMethod: (method: 'card' | 'sbp' | 'cash') => void;
  setContactPhone: (phone: string) => void;
  setContactEmail: (email: string) => void;
  reset: () => void;
}
