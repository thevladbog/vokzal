import { useNavigate } from 'react-router-dom';
import {
  Button,
  Card,
  Text,
  Badge,
  Spinner,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import {
  ArrowRight24Regular,
  SignOut24Regular,
  Person24Regular,
} from '@fluentui/react-icons';
import { useQuery } from '@tanstack/react-query';
import { useAuthStore } from '@/stores/authStore';
import { useScanStore } from '@/stores/scanStore';
import { tripService } from '@/services/trip';
import { formatTime } from '@/utils/format';

const useStyles = makeStyles({
  container: {
    minHeight: '100vh',
    backgroundColor: tokens.colorNeutralBackground3,
    padding: tokens.spacingVerticalL,
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: tokens.spacingVerticalXL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorBrandForeground1,
  },
  userInfo: {
    display: 'flex',
    alignItems: 'center',
    gap: tokens.spacingHorizontalS,
    fontSize: tokens.fontSizeBase300,
    color: tokens.colorNeutralForeground2,
  },
  tripsList: {
    display: 'flex',
    flexDirection: 'column',
    gap: tokens.spacingVerticalM,
  },
  tripCard: {
    padding: tokens.spacingVerticalL,
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    ':hover': {
      boxShadow: tokens.shadow8,
      transform: 'translateY(-2px)',
    },
  },
  tripHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: tokens.spacingVerticalS,
  },
  routeInfo: {
    display: 'flex',
    alignItems: 'center',
    gap: tokens.spacingHorizontalM,
  },
  time: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
  },
  station: {
    fontSize: tokens.fontSizeBase300,
    color: tokens.colorNeutralForeground2,
  },
  tripDetails: {
    display: 'flex',
    gap: tokens.spacingHorizontalL,
    marginTop: tokens.spacingVerticalS,
    flexWrap: 'wrap',
    fontSize: tokens.fontSizeBase300,
  },
  empty: {
    textAlign: 'center',
    padding: tokens.spacingVerticalXXL,
    color: tokens.colorNeutralForeground2,
  },
  loading: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '200px',
  },
});

export const TripSelectionPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const { setCurrentTrip, setStats } = useScanStore();

  const { data: trips, isLoading, error } = useQuery({
    queryKey: ['active-trips'],
    queryFn: () => tripService.getActive(),
    refetchInterval: 30000, // Refresh every 30 seconds
  });

  const handleTripSelect = async (tripId: string) => {
    const trip = trips?.find((t) => t.id === tripId);
    if (trip) {
      setCurrentTrip(trip);
      
      try {
        const stats = await tripService.getStats(tripId);
        setStats(stats);
      } catch (err) {
        console.error('Failed to fetch trip stats:', err);
      }

      navigate('/scan');
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'boarding':
        return <Badge appearance="filled" color="success">Посадка</Badge>;
      case 'scheduled':
        return <Badge appearance="filled" color="informative">По расписанию</Badge>;
      case 'departed':
        return <Badge appearance="filled" color="subtle">Отправлен</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Text className={styles.title}>Выберите рейс</Text>
        <div style={{ display: 'flex', gap: tokens.spacingHorizontalM, alignItems: 'center' }}>
          <div className={styles.userInfo}>
            <Person24Regular />
            <span>{user?.fullName}</span>
          </div>
          <Button
            appearance="subtle"
            icon={<SignOut24Regular />}
            onClick={handleLogout}
          >
            Выход
          </Button>
        </div>
      </div>

      {isLoading && (
        <div className={styles.loading}>
          <Spinner label="Загрузка рейсов..." />
        </div>
      )}

      {error && (
        <div className={styles.empty}>
          <Text>Ошибка загрузки рейсов</Text>
        </div>
      )}

      {trips && trips.length === 0 && (
        <div className={styles.empty}>
          <Text>Нет активных рейсов для посадки</Text>
        </div>
      )}

      {trips && trips.length > 0 && (
        <div className={styles.tripsList}>
          {trips.map((trip) => (
            <Card
              key={trip.id}
              className={styles.tripCard}
              onClick={() => handleTripSelect(trip.id)}
            >
              <div className={styles.tripHeader}>
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
                {getStatusBadge(trip.status)}
              </div>

              <div className={styles.tripDetails}>
                <Text>Автобус: {trip.busNumber}</Text>
                <Text>Водитель: {trip.driverName}</Text>
                {trip.platform && <Text>Платформа: {trip.platform}</Text>}
                <Text>
                  Мест: {trip.availableSeats}/{trip.totalSeats}
                </Text>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};
