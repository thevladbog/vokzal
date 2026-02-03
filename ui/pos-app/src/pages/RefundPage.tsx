import React, { useState } from 'react';
import {
  Card,
  Title2,
  Button,
  Input,
  Text,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import { SearchRegular, ArrowUndoRegular } from '@fluentui/react-icons';
import { posService } from '@/services/pos';

const useStyles = makeStyles({
  container: {
    padding: '24px',
    maxWidth: '600px',
    margin: '0 auto',
  },
  card: {
    padding: '24px',
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
    marginTop: '24px',
  },
  searchBar: {
    display: 'flex',
    gap: '12px',
    alignItems: 'center',
  },
  ticketInfo: {
    padding: '16px',
    backgroundColor: tokens.colorNeutralBackground3,
    borderRadius: '8px',
    marginTop: '16px',
  },
  warning: {
    color: tokens.colorPaletteRedForeground1,
    fontWeight: 'bold',
  },
});

export const RefundPage: React.FC = () => {
  const styles = useStyles();

  const [ticketId, setTicketId] = useState('');
  const [ticket, setTicket] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [processing, setProcessing] = useState(false);

  const handleSearch = async () => {
    if (!ticketId.trim()) return;

    setLoading(true);
    try {
      // TODO: Implement search ticket by ID
      alert('Функция поиска билета не реализована');
    } catch (err: any) {
      alert(`Ошибка поиска: ${err}`);
    } finally {
      setLoading(false);
    }
  };

  const handleRefund = async () => {
    if (!ticketId) return;

    if (!confirm('Вы уверены, что хотите вернуть этот билет?')) {
      return;
    }

    setProcessing(true);
    try {
      const returnedTicket = await posService.returnTicket(ticketId);

      // Печать чека возврата
      await posService.printReceipt({
        operation: 'refund',
        items: [
          {
            name: 'Возврат билета',
            quantity: 1,
            price: returnedTicket.price - (returnedTicket.refund_penalty || 0),
            vat: 'vat20',
          },
        ],
        payment: {
          type: 'cash',
          amount: returnedTicket.price - (returnedTicket.refund_penalty || 0),
        },
      });

      alert(
        `Билет возвращён!\n` +
        `Сумма возврата: ${(returnedTicket.price - (returnedTicket.refund_penalty || 0)).toFixed(2)} ₽\n` +
        `${returnedTicket.refund_penalty ? `Штраф: ${returnedTicket.refund_penalty.toFixed(2)} ₽` : ''}`
      );

      setTicketId('');
      setTicket(null);
    } catch (err: any) {
      alert(`Ошибка возврата: ${err}`);
    } finally {
      setProcessing(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card}>
        <Title2>Возврат билета</Title2>

        <div className={styles.form}>
          <div className={styles.searchBar}>
            <Input
              placeholder="ID билета или штрихкод"
              value={ticketId}
              onChange={(e) => setTicketId(e.target.value)}
              style={{ flex: 1 }}
              onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
            />
            <Button
              appearance="primary"
              icon={<SearchRegular />}
              onClick={handleSearch}
              disabled={loading}
            >
              Найти
            </Button>
          </div>

          {ticket && (
            <div className={styles.ticketInfo}>
              <Text weight="bold" block>Информация о билете</Text>
              <Text block>ID: {ticket.id}</Text>
              <Text block>Маршрут: {ticket.route_name}</Text>
              <Text block>Цена: {ticket.price} ₽</Text>
              {ticket.refund_penalty && (
                <Text block className={styles.warning}>
                  Штраф за возврат: {ticket.refund_penalty} ₽
                </Text>
              )}
              <Text block weight="bold" style={{ marginTop: '12px' }}>
                Сумма к возврату: {(ticket.price - (ticket.refund_penalty || 0)).toFixed(2)} ₽
              </Text>
            </div>
          )}

          <Button
            appearance="primary"
            size="large"
            icon={<ArrowUndoRegular />}
            onClick={handleRefund}
            disabled={!ticketId || processing}
          >
            {processing ? 'Возврат...' : 'Вернуть билет'}
          </Button>
        </div>
      </Card>
    </div>
  );
};
