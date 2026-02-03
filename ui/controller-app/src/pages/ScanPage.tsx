import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Button,
  Card,
  Text,
  Badge,
  ProgressBar,
  Divider,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import {
  ArrowLeft24Regular,
  Checkmark24Filled,
  Dismiss24Filled,
} from '@fluentui/react-icons';
import { useMutation } from '@tanstack/react-query';
import { useScanStore } from '@/stores/scanStore';
import { ticketService } from '@/services/ticket';
import { tripService } from '@/services/trip';
import { QRScanner } from '@/components/QRScanner';
import { formatTime, formatDateTime, playSuccessSound, playErrorSound, vibratePhone } from '@/utils/format';
import type { Ticket } from '@/types';

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
    marginBottom: tokens.spacingVerticalL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorBrandForeground1,
  },
  tripInfo: {
    padding: tokens.spacingVerticalL,
    marginBottom: tokens.spacingVerticalL,
  },
  routeInfo: {
    display: 'flex',
    alignItems: 'center',
    gap: tokens.spacingHorizontalM,
    marginBottom: tokens.spacingVerticalM,
  },
  time: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
  },
  statsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(2, 1fr)',
    gap: tokens.spacingVerticalM,
    marginTop: tokens.spacingVerticalM,
  },
  statItem: {
    display: 'flex',
    flexDirection: 'column',
    gap: tokens.spacingVerticalXS,
  },
  statLabel: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground2,
  },
  statValue: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
  },
  scannerCard: {
    padding: tokens.spacingVerticalL,
    marginBottom: tokens.spacingVerticalL,
  },
  resultCard: {
    padding: tokens.spacingVerticalL,
    marginBottom: tokens.spacingVerticalL,
  },
  successCard: {
    backgroundColor: tokens.colorPaletteGreenBackground2,
  },
  errorCard: {
    backgroundColor: tokens.colorPaletteRedBackground2,
  },
  resultHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: tokens.spacingHorizontalS,
    marginBottom: tokens.spacingVerticalM,
  },
  resultTitle: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
  },
  passengerInfo: {
    display: 'flex',
    flexDirection: 'column',
    gap: tokens.spacingVerticalXS,
  },
  recentScans: {
    padding: tokens.spacingVerticalL,
  },
  scansList: {
    display: 'flex',
    flexDirection: 'column',
    gap: tokens.spacingVerticalS,
    marginTop: tokens.spacingVerticalM,
  },
  scanItem: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: tokens.spacingVerticalS,
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
  },
});

export const ScanPage = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const { currentTrip, stats, recentScans, addRecentScan, setStats } = useScanStore();
  
  const [scanResult, setScanResult] = useState<{
    success: boolean;
    ticket?: Ticket;
    message: string;
  } | null>(null);
  const [isScanning, setIsScanning] = useState(true);

  useEffect(() => {
    if (!currentTrip) {
      navigate('/trips');
    }
  }, [currentTrip, navigate]);

  useEffect(() => {
    // Refresh stats every 10 seconds
    const interval = setInterval(async () => {
      if (currentTrip) {
        try {
          const updatedStats = await tripService.getStats(currentTrip.id);
          setStats(updatedStats);
        } catch (err) {
          console.error('Failed to refresh stats:', err);
        }
      }
    }, 10000);

    return () => clearInterval(interval);
  }, [currentTrip, setStats]);

  const boardingMutation = useMutation({
    mutationFn: (qrCode: string) => ticketService.markBoarding({
      ticketId: '',
      qrCode,
    }),
    onSuccess: (data) => {
      setScanResult({
        success: data.success,
        ticket: data.ticket,
        message: data.message,
      });
      
      if (data.success) {
        playSuccessSound();
        vibratePhone([100, 50, 100]);
        addRecentScan(data.ticket);
        
        // Refresh stats
        if (currentTrip) {
          tripService.getStats(currentTrip.id).then(setStats);
        }
      } else {
        playErrorSound();
        vibratePhone(500);
      }

      // Auto-clear result and resume scanning after 3 seconds
      setTimeout(() => {
        setScanResult(null);
        setIsScanning(true);
      }, 3000);
    },
    onError: (error) => {
      setScanResult({
        success: false,
        message: error instanceof Error ? error.message : 'Ошибка проверки билета',
      });
      playErrorSound();
      vibratePhone(500);

      // Auto-clear error and resume scanning after 3 seconds
      setTimeout(() => {
        setScanResult(null);
        setIsScanning(true);
      }, 3000);
    },
  });

  const handleQRScan = (qrCode: string) => {
    setIsScanning(false);
    setScanResult(null);
    boardingMutation.mutate(qrCode);
  };

  const handleBack = () => {
    navigate('/trips');
  };

  if (!currentTrip) {
    return null;
  }

  const progress = stats ? (stats.boardedTickets / stats.soldTickets) * 100 : 0;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Button
          appearance="subtle"
          icon={<ArrowLeft24Regular />}
          onClick={handleBack}
        >
          Назад
        </Button>
        <Text className={styles.title}>Контроль посадки</Text>
      </div>

      {/* Trip Info */}
      <Card className={styles.tripInfo}>
        <div className={styles.routeInfo}>
          <div>
            <Text className={styles.time}>{formatTime(currentTrip.departureTime)}</Text>
            <Text>{currentTrip.route?.fromStation?.name}</Text>
          </div>
          <Text>→</Text>
          <div>
            <Text className={styles.time}>{formatTime(currentTrip.arrivalTime)}</Text>
            <Text>{currentTrip.route?.toStation?.name}</Text>
          </div>
        </div>

        <Divider />

        <div className={styles.statsGrid}>
          <div className={styles.statItem}>
            <Text className={styles.statLabel}>Продано билетов</Text>
            <Text className={styles.statValue}>
              {stats?.soldTickets || 0}
            </Text>
          </div>
          <div className={styles.statItem}>
            <Text className={styles.statLabel}>Прошли посадку</Text>
            <Text className={styles.statValue} style={{ color: tokens.colorPaletteGreenForeground1 }}>
              {stats?.boardedTickets || 0}
            </Text>
          </div>
          <div className={styles.statItem}>
            <Text className={styles.statLabel}>Прогресс посадки</Text>
            <ProgressBar value={progress / 100} />
            <Text style={{ fontSize: tokens.fontSizeBase200 }}>
              {Math.round(progress)}%
            </Text>
          </div>
          <div className={styles.statItem}>
            <Text className={styles.statLabel}>Автобус / Платформа</Text>
            <Text className={styles.statValue}>
              {currentTrip.busNumber} / {currentTrip.platform || '—'}
            </Text>
          </div>
        </div>
      </Card>

      {/* QR Scanner */}
      {!scanResult && (
        <Card className={styles.scannerCard}>
          <QRScanner
            onScan={handleQRScan}
            isActive={isScanning}
          />
        </Card>
      )}

      {/* Scan Result */}
      {scanResult && (
        <Card
          className={`${styles.resultCard} ${
            scanResult.success ? styles.successCard : styles.errorCard
          }`}
        >
          <div className={styles.resultHeader}>
            {scanResult.success ? (
              <Checkmark24Filled style={{ color: tokens.colorPaletteGreenForeground1 }} />
            ) : (
              <Dismiss24Filled style={{ color: tokens.colorPaletteRedForeground1 }} />
            )}
            <Text className={styles.resultTitle}>{scanResult.message}</Text>
          </div>

          {scanResult.ticket && (
            <div className={styles.passengerInfo}>
              <Text>
                <strong>Пассажир:</strong> {scanResult.ticket.passenger?.fullName}
              </Text>
              <Text>
                <strong>Место:</strong> {scanResult.ticket.seatNumber}
              </Text>
              <Text>
                <strong>Документ:</strong> {scanResult.ticket.passenger?.documentType}{' '}
                {scanResult.ticket.passenger?.documentNumber}
              </Text>
              {scanResult.ticket.boardedAt && (
                <Text>
                  <strong>Время посадки:</strong> {formatDateTime(scanResult.ticket.boardedAt)}
                </Text>
              )}
            </div>
          )}
        </Card>
      )}

      {/* Recent Scans */}
      {recentScans.length > 0 && (
        <Card className={styles.recentScans}>
          <Text style={{ fontWeight: tokens.fontWeightSemibold, marginBottom: tokens.spacingVerticalM }}>
            Недавние проверки ({recentScans.length})
          </Text>
          <div className={styles.scansList}>
            {recentScans.slice(0, 5).map((ticket) => (
              <div key={ticket.id} className={styles.scanItem}>
                <div>
                  <Text style={{ fontWeight: tokens.fontWeightSemibold }}>
                    {ticket.passenger?.fullName}
                  </Text>
                  <Text style={{ fontSize: tokens.fontSizeBase200, color: tokens.colorNeutralForeground2 }}>
                    Место {ticket.seatNumber}
                  </Text>
                </div>
                <Badge appearance="filled" color="success" icon={<Checkmark24Filled />}>
                  {formatTime(ticket.boardedAt || '')}
                </Badge>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
};
