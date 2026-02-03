// API Response types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
}

// Entity types
export interface User {
  id: string;
  username: string;
  email: string;
  fullName: string;
  role: 'cashier' | 'dispatcher' | 'controller' | 'accountant' | 'admin';
  phone?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Station {
  id: string;
  name: string;
  code: string;
  city: string;
  address: string;
  latitude: number;
  longitude: number;
  timezone: string;
  active: boolean;
}

export interface Route {
  id: string;
  fromStationId: string;
  toStationId: string;
  distance: number;
  duration: number;
  active: boolean;
  fromStation?: Station;
  toStation?: Station;
}

export interface Trip {
  id: string;
  routeId: string;
  scheduleId: string;
  departureTime: string;
  arrivalTime: string;
  busNumber: string;
  driverName: string;
  totalSeats: number;
  availableSeats: number;
  price: number;
  status: 'scheduled' | 'boarding' | 'departed' | 'arrived' | 'cancelled';
  platform?: string;
  gate?: string;
  delayMinutes?: number;
  route?: Route;
}

export interface Passenger {
  id?: string;
  fullName: string;
  documentType: 'passport' | 'birth_certificate' | 'foreign_passport';
  documentNumber: string;
  phone?: string;
  email?: string;
  dateOfBirth?: string;
  benefits?: 'none' | 'child' | 'student' | 'senior' | 'disabled';
}

export interface Ticket {
  id: string;
  tripId: string;
  passengerId: string;
  seatNumber: string;
  price: number;
  status: 'active' | 'boarded' | 'cancelled' | 'returned';
  qrCode: string;
  barcode: string;
  soldAt: string;
  soldBy: string;
  boardedAt?: string;
  boardedBy?: string;
  returnedAt?: string;
  passenger?: Passenger;
  trip?: Trip;
}

export interface BoardingRequest {
  ticketId: string;
  qrCode: string;
}

export interface BoardingResponse {
  success: boolean;
  ticket: Ticket;
  message: string;
}

export interface TripStats {
  tripId: string;
  totalSeats: number;
  soldTickets: number;
  boardedTickets: number;
  availableSeats: number;
  boardingProgress: number; // percentage
}

// Store types
export interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  setUser: (user: User | null) => void;
  setTokens: (accessToken: string, refreshToken: string) => void;
}

export interface ScanState {
  currentTrip: Trip | null;
  stats: TripStats | null;
  recentScans: Ticket[];
  isScanning: boolean;
  setCurrentTrip: (trip: Trip | null) => void;
  setStats: (stats: TripStats | null) => void;
  addRecentScan: (ticket: Ticket) => void;
  clearRecentScans: () => void;
  setIsScanning: (isScanning: boolean) => void;
}
