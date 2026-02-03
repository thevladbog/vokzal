import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import {
  Card,
  Button,
  Text,
  Spinner,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import { Checkmark24Filled } from '@fluentui/react-icons';
import { ticketService } from '@/services/ticket';
import { Ticket } from '@/types';
import { QRCodeCanvas } from 'qrcode.react';
import { formatDateTime, formatPrice } from '@/utils/format';

const useStyles = makeStyles({
  container: {
    minHeight: '100vh',
    backgroundColor: tokens.colorNeutralBackground2,
    padding: tokens.spacingVerticalXXL,
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    maxWidth: '600px',
    width: '100%',
  },
  successIcon: {
    color: tokens.colorPaletteGreenForeground1,
    fontSize: '64px',
    textAlign: 'center',
    marginBottom: tokens.spacingVerticalL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    textAlign: 'center',
    marginBottom: tokens.spacingVerticalM,
  },
  ticketCard: {
    marginBottom: tokens.spacingVerticalM,
    padding: tokens.spacingVerticalL,
  },
  ticketInfo: {
    display: 'grid',
    gap: tokens.spacingVerticalS,
    marginBottom: tokens.spacingVerticalM,
  },
  qrCode: {
    textAlign: 'center',
    marginTop: tokens.spacingVerticalM,
  },
  buttonGroup: {
    display: 'flex',
    gap: tokens.spacingHorizontalM,
    marginTop: tokens.spacingVerticalXL,
  },
});

export const ConfirmationPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const ticketIds = searchParams.get('ticketIds')?.split(',') || [];

  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadTickets = async () => {
      try {
        const ticketPromises = ticketIds.map((id) => ticketService.getById(id));
        const loadedTickets = await Promise.all(ticketPromises);
        setTickets(loadedTickets);
      } catch (error) {
        console.error('Failed to load tickets:', error);
      } finally {
        setLoading(false);
      }
    };

    if (ticketIds.length > 0) {
      loadTickets();
    } else {
      setLoading(false);
    }
  }, []);

  if (loading) {
    return (
      <div className={styles.container}>
        <Spinner label="Загрузка билетов..." />
      </div>
    );
  }

  if (tickets.length === 0) {
    return (
      <div className={styles.container}>
        <div className={styles.content}>
          <Text>Билеты не найдены</Text>
          <Button onClick={() => navigate('/')} style={{ marginTop: tokens.spacingVerticalM }}>
            Вернуться на главную
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <div className={styles.successIcon}>
          <Checkmark24Filled />
        </div>
        <Text className={styles.title}>Билеты успешно оформлены!</Text>

        {tickets.map((ticket) => (
          <Card key={ticket.id} className={styles.ticketCard}>
            <div className={styles.ticketInfo}>
              <Text weight="semibold">Билет №{ticket.number}</Text>
              <Text>
                Пассажир: {ticket.passengerLastName} {ticket.passengerFirstName}{' '}
                {ticket.passengerMiddleName}
              </Text>
              <Text>
                Маршрут: {ticket.trip?.route?.fromStation?.name} →{' '}
                {ticket.trip?.route?.toStation?.name}
              </Text>
              <Text>Отправление: {formatDateTime(ticket.trip?.departureTime || '')}</Text>
              <Text>Цена: {formatPrice(ticket.finalPrice)}</Text>
              {ticket.seatNumber && <Text>Место: {ticket.seatNumber}</Text>}
            </div>

            {ticket.qrCode && (
              <div className={styles.qrCode}>
                <QRCodeCanvas value={ticket.qrCode} size={200} />
                <Text style={{ display: 'block', marginTop: tokens.spacingVerticalS }}>
                  Предъявите этот QR-код при посадке
                </Text>
              </div>
            )}
          </Card>
        ))}

        <div className={styles.buttonGroup}>
          <Button onClick={() => navigate('/')} style={{ flex: 1 }}>
            Вернуться на главную
          </Button>
          <Button
            appearance="primary"
            onClick={() => window.print()}
            style={{ flex: 1 }}
          >
            Печать билетов
          </Button>
        </div>
      </div>
    </div>
  );
};
