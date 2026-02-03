import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  Button,
  Text,
  Spinner,
  Badge,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import { SignOut24Regular } from '@fluentui/react-icons';
import { useAuthStore } from '@/stores/authStore';
import { ticketService } from '@/services/ticket';
import { Ticket } from '@/types';
import { QRCodeCanvas } from 'qrcode.react';
import { formatDateTime, formatPrice } from '@/utils/format';

const useStyles = makeStyles({
  container: {
    minHeight: '100vh',
    backgroundColor: tokens.colorNeutralBackground2,
  },
  header: {
    backgroundColor: tokens.colorBrandBackground,
    color: tokens.colorNeutralForegroundInverted,
    padding: tokens.spacingVerticalXL,
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  content: {
    maxWidth: '900px',
    margin: '0 auto',
    padding: tokens.spacingVerticalXXL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalXL,
  },
  ticketCard: {
    marginBottom: tokens.spacingVerticalM,
    padding: tokens.spacingVerticalL,
  },
  ticketHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: tokens.spacingVerticalM,
  },
  ticketInfo: {
    display: 'grid',
    gridTemplateColumns: '1fr auto',
    gap: tokens.spacingVerticalS,
  },
  qrCode: {
    textAlign: 'center',
    marginTop: tokens.spacingVerticalM,
    paddingTop: tokens.spacingVerticalM,
    borderTop: `1px solid ${tokens.colorNeutralStroke1}`,
  },
  emptyState: {
    textAlign: 'center',
    padding: tokens.spacingVerticalXXL,
    color: tokens.colorNeutralForeground2,
  },
});

export const MyTicketsPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();

  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadTickets = async () => {
      try {
        const data = await ticketService.getUserTickets();
        setTickets(data);
      } catch (error) {
        console.error('Failed to load tickets:', error);
      } finally {
        setLoading(false);
      }
    };

    loadTickets();
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  const getStatusBadge = (status: Ticket['status']) => {
    switch (status) {
      case 'sold':
        return <Badge appearance="filled" color="success">Куплен</Badge>;
      case 'boarded':
        return <Badge appearance="filled" color="informative">Посадка выполнена</Badge>;
      case 'returned':
        return <Badge appearance="filled" color="warning">Возвращён</Badge>;
      case 'expired':
        return <Badge appearance="filled" color="danger">Истёк</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Text weight="semibold">
          {user?.lastName} {user?.firstName}
        </Text>
        <Button
          appearance="subtle"
          icon={<SignOut24Regular />}
          onClick={handleLogout}
          style={{ color: tokens.colorNeutralForegroundInverted }}
        >
          Выйти
        </Button>
      </div>

      <div className={styles.content}>
        <Text className={styles.title}>Мои билеты</Text>

        {loading && (
          <div style={{ textAlign: 'center', padding: tokens.spacingVerticalXXL }}>
            <Spinner label="Загрузка билетов..." />
          </div>
        )}

        {!loading && tickets.length === 0 && (
          <div className={styles.emptyState}>
            <Text>У вас пока нет купленных билетов</Text>
            <Button
              appearance="primary"
              onClick={() => navigate('/')}
              style={{ marginTop: tokens.spacingVerticalM }}
            >
              Найти билеты
            </Button>
          </div>
        )}

        {!loading &&
          tickets.map((ticket) => (
            <Card key={ticket.id} className={styles.ticketCard}>
              <div className={styles.ticketHeader}>
                <Text weight="semibold">Билет №{ticket.number}</Text>
                {getStatusBadge(ticket.status)}
              </div>

              <div className={styles.ticketInfo}>
                <Text>Пассажир:</Text>
                <Text>
                  {ticket.passengerLastName} {ticket.passengerFirstName}{' '}
                  {ticket.passengerMiddleName}
                </Text>

                <Text>Маршрут:</Text>
                <Text>
                  {ticket.trip?.route?.fromStation?.name} →{' '}
                  {ticket.trip?.route?.toStation?.name}
                </Text>

                <Text>Отправление:</Text>
                <Text>{formatDateTime(ticket.trip?.departureTime || '')}</Text>

                {ticket.seatNumber && (
                  <>
                    <Text>Место:</Text>
                    <Text>{ticket.seatNumber}</Text>
                  </>
                )}

                <Text>Цена:</Text>
                <Text>{formatPrice(ticket.finalPrice)}</Text>

                <Text>Куплен:</Text>
                <Text>{formatDateTime(ticket.soldAt)}</Text>
              </div>

              {ticket.status === 'sold' && ticket.qrCode && (
                <div className={styles.qrCode}>
                  <QRCodeCanvas value={ticket.qrCode} size={150} />
                  <Text style={{ display: 'block', marginTop: tokens.spacingVerticalS }}>
                    QR-код для посадки
                  </Text>
                </div>
              )}

              {ticket.status === 'sold' && (
                <Button
                  appearance="subtle"
                  style={{ marginTop: tokens.spacingVerticalM }}
                  onClick={async () => {
                    if (
                      window.confirm(
                        'Вы уверены, что хотите вернуть этот билет? Может быть удержана комиссия.'
                      )
                    ) {
                      try {
                        await ticketService.requestReturn(ticket.id);
                        // Reload tickets
                        const data = await ticketService.getUserTickets();
                        setTickets(data);
                      } catch (error) {
                        alert('Ошибка при возврате билета');
                      }
                    }
                  }}
                >
                  Вернуть билет
                </Button>
              )}
            </Card>
          ))}
      </div>
    </div>
  );
};
