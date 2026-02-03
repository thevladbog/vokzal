import { makeStyles, tokens } from '@fluentui/react-components';
import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';
import type { BoardTrip } from '@/types';

const useStyles = makeStyles({
  row: {
    display: 'grid',
    gridTemplateColumns: '1fr 2fr 2fr 1fr 1fr 1fr',
    padding: '20px 16px',
    borderBottom: `2px solid ${tokens.colorNeutralStroke1}`,
    fontSize: '24px',
    alignItems: 'center',
    transition: 'background-color 0.3s',
  },
  time: {
    fontWeight: 'bold',
    fontSize: '28px',
  },
  route: {
    fontWeight: '600',
  },
  platform: {
    fontSize: '32px',
    fontWeight: 'bold',
    color: tokens.colorBrandForeground1,
    textAlign: 'center',
  },
  statusScheduled: {
    color: tokens.colorNeutralForeground1,
  },
  statusBoarding: {
    color: tokens.colorPaletteGreenForeground1,
    fontWeight: 'bold',
    animation: 'pulse 2s infinite',
  },
  statusDeparted: {
    color: tokens.colorNeutralForeground3,
    textDecoration: 'line-through',
  },
  statusCancelled: {
    color: tokens.colorPaletteRedForeground1,
    fontWeight: 'bold',
  },
  statusDelayed: {
    color: tokens.colorPaletteYellowForeground1,
    fontWeight: 'bold',
  },
  delay: {
    color: tokens.colorPaletteRedForeground1,
    fontSize: '18px',
  },
});

interface TripRowProps {
  trip: BoardTrip;
}

export const TripRow: React.FC<TripRowProps> = ({ trip }) => {
  const styles = useStyles();

  const getStatusText = (status: BoardTrip['status']) => {
    switch (status) {
      case 'scheduled':
        return 'По расписанию';
      case 'boarding':
        return 'ПОСАДКА';
      case 'departed':
        return 'Отправился';
      case 'cancelled':
        return 'ОТМЕНЁН';
      case 'delayed':
        return 'Задержка';
      default:
        return status;
    }
  };

  const getStatusClass = (status: BoardTrip['status']) => {
    switch (status) {
      case 'boarding':
        return styles.statusBoarding;
      case 'departed':
        return styles.statusDeparted;
      case 'cancelled':
        return styles.statusCancelled;
      case 'delayed':
        return styles.statusDelayed;
      default:
        return styles.statusScheduled;
    }
  };

  const departureTime = format(parseISO(trip.departure_datetime), 'HH:mm', { locale: ru });
  const departureDate = format(parseISO(trip.departure_datetime), 'dd.MM', { locale: ru });

  return (
    <div className={styles.row}>
      <div className={styles.time}>
        {departureTime}
        <div style={{ fontSize: '16px', opacity: 0.7 }}>{departureDate}</div>
      </div>
      
      <div className={styles.route}>{trip.route_name}</div>
      
      <div>{trip.arrival_station}</div>
      
      <div className={styles.platform}>{trip.platform || '-'}</div>
      
      <div className={getStatusClass(trip.status)}>
        {getStatusText(trip.status)}
        {trip.delay_minutes && trip.delay_minutes > 0 && (
          <div className={styles.delay}>+{trip.delay_minutes} мин</div>
        )}
      </div>
      
      <div style={{ textAlign: 'center', opacity: 0.7 }}>
        {trip.available_seats !== undefined ? `${trip.available_seats} мест` : ''}
      </div>
    </div>
  );
};
