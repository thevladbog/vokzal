import React from 'react';
import {
  Card,
  CardHeader,
  Text,
  Badge,
  Button,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import { ArrowRight24Regular, People24Regular, Seat24Regular } from '@fluentui/react-icons';
import { Trip } from '@/types';
import { formatTime, formatDuration, formatPrice } from '@/utils/format';

const useStyles = makeStyles({
  card: {
    marginBottom: tokens.spacingVerticalM,
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    ':hover': {
      boxShadow: tokens.shadow8,
      transform: 'translateY(-2px)',
    },
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  routeInfo: {
    display: 'flex',
    alignItems: 'center',
    gap: tokens.spacingHorizontalM,
    marginBottom: tokens.spacingVerticalS,
  },
  time: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
  },
  station: {
    fontSize: tokens.fontSizeBase400,
    color: tokens.colorNeutralForeground2,
  },
  details: {
    display: 'flex',
    gap: tokens.spacingHorizontalL,
    marginTop: tokens.spacingVerticalS,
    flexWrap: 'wrap',
  },
  detail: {
    display: 'flex',
    alignItems: 'center',
    gap: tokens.spacingHorizontalXS,
  },
  price: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorBrandForeground1,
  },
  button: {
    marginTop: tokens.spacingVerticalM,
  },
});

interface TripCardProps {
  trip: Trip;
  onSelect: (trip: Trip) => void;
}

export const TripCard: React.FC<TripCardProps> = ({ trip, onSelect }) => {
  const styles = useStyles();

  const getStatusBadge = (status: Trip['status']) => {
    switch (status) {
      case 'scheduled':
        return <Badge appearance="filled" color="success">По расписанию</Badge>;
      case 'boarding':
        return <Badge appearance="filled" color="warning">Посадка</Badge>;
      case 'departed':
        return <Badge appearance="filled" color="informative">Отправлен</Badge>;
      case 'cancelled':
        return <Badge appearance="filled" color="danger">Отменён</Badge>;
      default:
        return null;
    }
  };

  const isBookingAvailable = trip.status === 'scheduled' && trip.availableSeats > 0;

  return (
    <Card className={styles.card} onClick={() => isBookingAvailable && onSelect(trip)}>
      <CardHeader
        header={
          <div>
            <div className={styles.routeInfo}>
              <div>
                <Text className={styles.time}>{formatTime(trip.departureTime)}</Text>
                <Text className={styles.station}>{trip.route?.fromStation?.name}</Text>
              </div>
              <ArrowRight24Regular />
              <div>
                <Text className={styles.time}>{formatTime(trip.arrivalTime)}</Text>
                <Text className={styles.station}>{trip.route?.toStation?.name}</Text>
              </div>
            </div>

            <div className={styles.details}>
              <div className={styles.detail}>
                <Text>В пути: {formatDuration(trip.route?.duration || 0)}</Text>
              </div>
              <div className={styles.detail}>
                <People24Regular />
                <Text>Автобус: {trip.busNumber}</Text>
              </div>
              <div className={styles.detail}>
                <Seat24Regular />
                <Text>
                  Мест: {trip.availableSeats} из {trip.totalSeats}
                </Text>
              </div>
              {trip.platform && (
                <div className={styles.detail}>
                  <Text>Платформа: {trip.platform}</Text>
                </div>
              )}
            </div>

            <div className={styles.details}>
              <Text className={styles.price}>{formatPrice(trip.price)}</Text>
              {getStatusBadge(trip.status)}
            </div>

            {isBookingAvailable && (
              <Button
                className={styles.button}
                appearance="primary"
                onClick={() => onSelect(trip)}
              >
                Выбрать рейс
              </Button>
            )}
            {!isBookingAvailable && trip.status === 'scheduled' && (
              <Text>Мест нет в наличии</Text>
            )}
          </div>
        }
      />
    </Card>
  );
};
