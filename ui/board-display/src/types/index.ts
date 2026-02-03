// Trip for board display
export interface BoardTrip {
  id: string;
  route_name: string;
  departure_station: string;
  arrival_station: string;
  departure_datetime: string;
  arrival_datetime?: string;
  status: 'scheduled' | 'boarding' | 'departed' | 'arrived' | 'cancelled' | 'delayed';
  platform?: string;
  delay_minutes?: number;
  available_seats?: number;
}

// WebSocket message types
export interface WSMessage {
  type: 'trip_update' | 'trip_created' | 'trip_deleted' | 'status_changed';
  data: BoardTrip | BoardTrip[];
}

// Board types
export type BoardType = 'public' | 'platform';

// Platform board config
export interface PlatformConfig {
  platformId: string;
  platformName: string;
}
