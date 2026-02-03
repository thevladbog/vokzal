import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { makeStyles, tokens, Spinner } from '@fluentui/react-components';
import { useBoardStore } from '@/stores/boardStore';
import { boardService } from '@/services/api';
import { useWebSocket } from '@/hooks/useWebSocket';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

const useStyles = makeStyles({
  container: {
    minHeight: '100vh',
    backgroundColor: tokens.colorNeutralBackground1,
    padding: '24px',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '32px',
    paddingBottom: '24px',
    borderBottom: `3px solid ${tokens.colorBrandStroke1}`,
  },
  platformTitle: {
    fontSize: '64px',
    fontWeight: 'bold',
    color: tokens.colorBrandForeground1,
  },
  clock: {
    fontSize: '36px',
    fontWeight: 'bold',
  },
  tableHeader: {
    display: 'grid',
    gridTemplateColumns: '1fr 3fr 3fr 1fr',
    padding: '16px',
    backgroundColor: tokens.colorBrandBackground,
    color: 'white',
    fontSize: '20px',
    fontWeight: 'bold',
    borderRadius: '8px 8px 0 0',
  },
  table: {
    backgroundColor: 'white',
    borderRadius: '0 0 8px 8px',
    boxShadow: tokens.shadow16,
  },
  loading: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '400px',
  },
  row: {
    display: 'grid',
    gridTemplateColumns: '1fr 3fr 3fr 1fr',
    padding: '24px 16px',
    borderBottom: `2px solid ${tokens.colorNeutralStroke1}`,
    fontSize: '28px',
    alignItems: 'center',
  },
  time: {
    fontWeight: 'bold',
    fontSize: '36px',
  },
  route: {
    fontWeight: '600',
    fontSize: '32px',
  },
  statusBoarding: {
    color: tokens.colorPaletteGreenForeground1,
    fontWeight: 'bold',
    fontSize: '32px',
    animation: 'pulse 2s infinite',
  },
  connectionStatus: {
    position: 'fixed',
    top: '16px',
    right: '16px',
    padding: '8px 16px',
    borderRadius: '8px',
    fontSize: '14px',
    fontWeight: 'bold',
  },
  connected: {
    backgroundColor: tokens.colorPaletteGreenBackground2,
    color: tokens.colorPaletteGreenForeground1,
  },
  disconnected: {
    backgroundColor: tokens.colorPaletteRedBackground2,
    color: tokens.colorPaletteRedForeground1,
  },
});

export const PlatformBoardPage = () => {
  const styles = useStyles();
  const [searchParams] = useSearchParams();
  const platformId = searchParams.get('platform') || '1';
  const platformName = searchParams.get('name') || `Перрон ${platformId}`;

  const { trips, setTrips, updateTrip, addTrip } = useBoardStore();
  const [loading, setLoading] = useState(true);
  const [currentTime, setCurrentTime] = useState(new Date());

  // WebSocket connection
  const { connected } = useWebSocket((newTrips) => {
    newTrips.forEach((trip) => {
      // Filter only trips for this platform
      if (trip.platform === platformId) {
        const existingTrip = trips.find((t) => t.id === trip.id);
        if (existingTrip) {
          updateTrip(trip);
        } else {
          addTrip(trip);
        }
      }
    });
  });

  // Load initial data
  useEffect(() => {
    const loadTrips = async () => {
      try {
        const data = await boardService.getPlatformBoard(platformId);
        setTrips(data);
      } catch (err) {
        console.error('Failed to load platform trips:', err);
      } finally {
        setLoading(false);
      }
    };

    loadTrips();

    // Refresh every minute
    const interval = setInterval(loadTrips, 60000);
    return () => clearInterval(interval);
  }, [platformId, setTrips]);

  // Update clock
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className={styles.loading}>
        <Spinner label="Загрузка табло..." size="huge" />
      </div>
    );
  }

  // Sort trips by departure time
  const sortedTrips = [...trips]
    .filter((trip) => trip.platform === platformId)
    .sort((a, b) => a.departure_datetime.localeCompare(b.departure_datetime));

  return (
    <div className={styles.container}>
      {/* Connection status */}
      <div className={`${styles.connectionStatus} ${connected ? styles.connected : styles.disconnected}`}>
        {connected ? '● Онлайн' : '○ Отключено'}
      </div>

      {/* Header */}
      <div className={styles.header}>
        <div className={styles.platformTitle}>📍 {platformName}</div>
        <div className={styles.clock}>
          {format(currentTime, 'HH:mm:ss', { locale: ru })}
          <div style={{ fontSize: '24px', opacity: 0.8 }}>
            {format(currentTime, 'dd MMMM yyyy', { locale: ru })}
          </div>
        </div>
      </div>

      {/* Table */}
      <div className={styles.tableHeader}>
        <div>Время</div>
        <div>Маршрут</div>
        <div>Направление</div>
        <div style={{ textAlign: 'center' }}>Статус</div>
      </div>

      <div className={styles.table}>
        {sortedTrips.length > 0 ? (
          sortedTrips.map((trip) => (
            <div key={trip.id} className={styles.row}>
              <div className={styles.time}>
                {format(new Date(trip.departure_datetime), 'HH:mm', { locale: ru })}
              </div>
              <div className={styles.route}>{trip.route_name}</div>
              <div style={{ fontSize: '24px' }}>{trip.arrival_station}</div>
              <div className={trip.status === 'boarding' ? styles.statusBoarding : ''} style={{ textAlign: 'center' }}>
                {trip.status === 'boarding' ? 'ПОСАДКА' : 'По расписанию'}
              </div>
            </div>
          ))
        ) : (
          <div style={{ padding: '48px', textAlign: 'center', fontSize: '28px', opacity: 0.5 }}>
            Нет рейсов с этого перрона
          </div>
        )}
      </div>
    </div>
  );
};
