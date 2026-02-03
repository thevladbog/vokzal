// Trip
export interface Trip {
  id: string;
  route_name: string;
  departure_station: string;
  arrival_station: string;
  departure_datetime: string;
  arrival_datetime?: string;
  price: number;
  available_seats: number;
  status: 'scheduled' | 'boarding' | 'departed' | 'arrived' | 'cancelled';
  platform?: string;
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
  refund_penalty?: number;
  created_at: string;
}

// Sale request
export interface SaleRequest {
  trip_id: string;
  passenger_fio?: string;
  passenger_phone?: string;
  seat_id?: string;
}

// Receipt
export interface Receipt {
  success: boolean;
  fiscal_sign: string;
  ofd_url: string;
  receipt_num?: number;
}

// Payment
export interface Payment {
  id: string;
  ticket_id?: string;
  amount: number;
  method: 'card' | 'sbp' | 'cash';
  provider: 'tinkoff' | 'sbp' | 'manual';
  status: 'pending' | 'confirmed' | 'failed' | 'refunded';
  payment_url?: string;
  qr_code?: string;
}

// Auth
export interface User {
  id: string;
  username: string;
  fio: string;
  role: 'cashier' | 'dispatcher' | 'admin';
  station_id?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}
