// User roles
export type UserRole = 'admin' | 'dispatcher' | 'cashier' | 'controller' | 'accountant';

// Auth
export interface User {
  id: string;
  username: string;
  fio: string;
  role: UserRole;
  station_id?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

// Station
export interface Station {
  id: string;
  name: string;
  code: string;
  address: string;
  latitude: number;
  longitude: number;
  platforms_count: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Bus
export interface Bus {
  id: string;
  registration_number: string;
  model: string;
  seats_count: number;
  station_id: string;
  status: 'active' | 'maintenance' | 'inactive';
  created_at: string;
  updated_at: string;
}

// Route
export interface Route {
  id: string;
  name: string;
  stops: Array<{
    station_id: string;
    station_name?: string;
    order: number;
    distance_km?: number;
    duration_min?: number;
  }>;
  distance_km: number;
  duration_min: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Schedule
export interface Schedule {
  id: string;
  route_id: string;
  route_name?: string;
  departure_station_id: string;
  arrival_station_id: string;
  departure_time: string;
  days_of_week: number[];
  price: number;
  bus_id?: string;
  driver_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Trip
export interface Trip {
  id: string;
  schedule_id: string;
  bus_id?: string;
  driver_id?: string;
  departure_datetime: string;
  arrival_datetime?: string;
  status: 'scheduled' | 'boarding' | 'departed' | 'arrived' | 'cancelled';
  platform?: string;
  available_seats: number;
  created_at: string;
  updated_at: string;
}

// Ticket
export interface Ticket {
  id: string;
  trip_id: string;
  passenger_fio?: string;
  passenger_phone?: string;
  seat_id?: string;
  price: number;
  status: 'active' | 'returned' | 'used' | 'expired';
  qr_code: string;
  bar_code: string;
  issued_by_user_id?: string;
  created_at: string;
  updated_at: string;
}

// Payment
export interface Payment {
  id: string;
  ticket_id?: string;
  amount: number;
  method: 'card' | 'sbp' | 'cash';
  provider: 'tinkoff' | 'sbp' | 'manual';
  status: 'pending' | 'confirmed' | 'failed' | 'refunded';
  created_at: string;
  confirmed_at?: string;
}

// Report types
export interface SalesReport {
  date: string;
  tickets_sold: number;
  total_revenue: number;
  cash_revenue: number;
  card_revenue: number;
  sbp_revenue: number;
}

export interface OccupancyReport {
  trip_id: string;
  route_name: string;
  departure_datetime: string;
  total_seats: number;
  sold_seats: number;
  occupancy_percent: number;
}
