import React, { useEffect } from 'react';
import { makeStyles, tokens, Spinner } from '@fluentui/react-components';
import { TripRow } from '@/components/TripRow';
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
  logo: {
    height: '56px',
    width: 'auto',
    display: 'block',
  },
  clock: {
    fontSize: '36px',
    fontWeight: 'bold',
  },
  tableHeader: {
    display: 'grid',
    gridTemplateColumns: '1fr 2fr 2fr 1fr 1fr 1fr',
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
  footer: {
    marginTop: '32px',
    textAlign: 'center',
    fontSize: '18px',
    opacity: 0.7,
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

export const PublicBoardPage: React.FC = () => {
  const styles = useStyles();
  const { trips, setTrips, updateTrip, addTrip } = useBoardStore();
  const [loading, setLoading] = React.useState(true);
  const [currentTime, setCurrentTime] = React.useState(new Date());

  // WebSocket connection
  const { connected } = useWebSocket((newTrips) => {
    newTrips.forEach((trip) => {
      const existingTrip = trips.find((t) => t.id === trip.id);
      if (existingTrip) {
        updateTrip(trip);
      } else {
        addTrip(trip);
      }
    });
  });

  // Load initial data
  useEffect(() => {
    const loadTrips = async () => {
      try {
        const data = await boardService.getPublicBoard();
        setTrips(data);
      } catch (err) {
        console.error('Failed to load trips:', err);
      } finally {
        setLoading(false);
      }
    };

    loadTrips();

    // Refresh every minute
    const interval = setInterval(loadTrips, 60000);
    return () => clearInterval(interval);
  }, [setTrips]);

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
  const sortedTrips = [...trips].sort((a, b) =>
    a.departure_datetime.localeCompare(b.departure_datetime)
  );

  return (
    <div className={styles.container}>
      {/* Connection status */}
      <div className={`${styles.connectionStatus} ${connected ? styles.connected : styles.disconnected}`}>
        {connected ? '● Онлайн' : '○ Отключено'}
      </div>

      {/* Header */}
      <div className={styles.header}>
        <img src="/vokzal-logo.svg" alt="Вокзал.ТЕХ" className={styles.logo} />
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
        <div style={{ textAlign: 'center' }}>Перрон</div>
        <div style={{ textAlign: 'center' }}>Статус</div>
        <div style={{ textAlign: 'center' }}>Места</div>
      </div>

      <div className={styles.table}>
        {sortedTrips.length > 0 ? (
          sortedTrips.map((trip) => <TripRow key={trip.id} trip={trip} />)
        ) : (
          <div style={{ padding: '48px', textAlign: 'center', fontSize: '24px', opacity: 0.5 }}>
            Нет рейсов для отображения
          </div>
        )}
      </div>

      {/* Footer */}
      <div className={styles.footer}>
        © 2025 Вокзал.ТЕХ — Приятного путешествия!
      </div>
    </div>
  );
};
