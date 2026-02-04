// User roles
export type UserRole = 'admin' | 'dispatcher' | 'cashier' | 'controller' | 'accountant';

// Auth (login response user)
export interface User {
  id: string;
  username: string;
  fio?: string;
  full_name?: string;
  role: UserRole;
  station_id?: string;
}

// User from Users API (admin CRUD)
export interface UserAdmin {
  id: string;
  username: string;
  full_name: string;
  role: UserRole;
  station_id: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ListUsersResponse {
  users: UserAdmin[];
  total: number;
  page: number;
  limit: number;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  full_name: string;
  role: UserRole;
  station_id?: string | null;
}

export interface UpdateUserRequest {
  full_name?: string;
  password?: string;
  role?: UserRole;
  station_id?: string | null;
  is_active?: boolean;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

// Station (автовокзал)
export interface Station {
  id: string;
  name: string;
  code: string;
  address?: string;
  timezone?: string;
  latitude?: number;
  longitude?: number;
  platforms_count?: number;
  is_active?: boolean;
  created_at: string;
  updated_at: string;
}

// Bus (backend: plate_number, capacity)
export interface Bus {
  id: string;
  plate_number: string;
  model: string;
  capacity: number;
  station_id: string;
  status: 'active' | 'maintenance' | 'out_of_service';
  created_at: string;
  updated_at: string;
}

// Driver
export interface Driver {
  id: string;
  full_name: string;
  license_number: string;
  experience_years?: number | null;
  phone?: string | null;
  station_id: string;
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

// Schedule (matches backend: id, route_id, departure_time, days_of_week, is_active, created_at, updated_at, platform, nested route)
export interface Schedule {
  id: string;
  route_id: string;
  departure_time: string;
  days_of_week: number[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
  platform?: string;
  route?: Route;
  // Optional/legacy fields not returned by API:
  price?: number;
  departure_station_id?: string;
  arrival_station_id?: string;
  bus_id?: string;
  driver_id?: string;
}

// Trip
export interface Trip {
  id: string;
  schedule_id: string;
  bus_id?: string;
  driver_id?: string;
  date?: string;
  departure_datetime?: string;
  arrival_datetime?: string;
  status: 'scheduled' | 'boarding' | 'departed' | 'arrived' | 'cancelled' | 'delayed';
  delay_minutes?: number;
  platform?: string;
  available_seats?: number;
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

// Audit log (152-ФЗ)
export interface AuditLog {
  id: string;
  entity_type: string;
  entity_id: string;
  action: string;
  created_at: string;
  user_id?: string | null;
  ip_address?: string | null;
  user_agent?: string | null;
  old_value?: unknown;
  new_value?: unknown;
}
