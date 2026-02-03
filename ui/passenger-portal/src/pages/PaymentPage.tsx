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
import { paymentService } from '@/services/payment';
import { Payment } from '@/types';
import { QRCodeCanvas } from 'qrcode.react';
import { formatPrice } from '@/utils/format';

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
    maxWidth: '500px',
    width: '100%',
  },
  card: {
    padding: tokens.spacingVerticalXL,
    textAlign: 'center',
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalM,
  },
  qrCode: {
    margin: `${tokens.spacingVerticalXL} 0`,
  },
  info: {
    marginBottom: tokens.spacingVerticalM,
  },
  buttonGroup: {
    display: 'flex',
    gap: tokens.spacingHorizontalM,
    marginTop: tokens.spacingVerticalXL,
  },
});

export const PaymentPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const paymentId = searchParams.get('paymentId');

  const [payment, setPayment] = useState<Payment | null>(null);
  const [loading, setLoading] = useState(true);
  const [checking, setChecking] = useState(false);

  useEffect(() => {
    const loadPayment = async () => {
      if (!paymentId) {
        setLoading(false);
        return;
      }

      try {
        const data = await paymentService.getById(paymentId);
        setPayment(data);

        // If payment is completed, redirect to confirmation
        if (data.status === 'completed') {
          navigate(`/confirmation?ticketIds=${data.ticketId}`);
        }
      } catch (error) {
        console.error('Failed to load payment:', error);
      } finally {
        setLoading(false);
      }
    };

    loadPayment();
  }, [paymentId, navigate]);

  const handleCheckStatus = async () => {
    if (!paymentId) return;

    setChecking(true);
    try {
      const data = await paymentService.checkStatus(paymentId);
      setPayment(data);

      if (data.status === 'completed') {
        navigate(`/confirmation?ticketIds=${data.ticketId}`);
      }
    } catch (error) {
      console.error('Failed to check payment status:', error);
    } finally {
      setChecking(false);
    }
  };

  if (loading) {
    return (
      <div className={styles.container}>
        <Spinner label="Загрузка информации об оплате..." />
      </div>
    );
  }

  if (!payment) {
    return (
      <div className={styles.container}>
        <div className={styles.content}>
          <Text>Платёж не найден</Text>
          <Button onClick={() => navigate('/')} style={{ marginTop: tokens.spacingVerticalM }}>
            Вернуться на главную
          </Button>
        </div>
      </div>
    );
  }

  const renderPaymentContent = () => {
    if (payment.method === 'sbp' && payment.qrCodeUrl) {
      return (
        <>
          <Text className={styles.info}>Отсканируйте QR-код для оплаты через СБП</Text>
          <div className={styles.qrCode}>
            <QRCodeCanvas value={payment.qrCodeUrl} size={300} />
          </div>
          <Text className={styles.info}>Сумма к оплате: {formatPrice(payment.amount)}</Text>
        </>
      );
    }

    if (payment.method === 'card' && payment.paymentUrl) {
      return (
        <>
          <Text className={styles.info}>
            Вы будете перенаправлены на страницу оплаты банковской картой
          </Text>
          <Text className={styles.info}>Сумма к оплате: {formatPrice(payment.amount)}</Text>
          <Button
            appearance="primary"
            onClick={() => window.location.href = payment.paymentUrl!}
          >
            Перейти к оплате
          </Button>
        </>
      );
    }

    return <Text>Неизвестный способ оплаты</Text>;
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <Card className={styles.card}>
          <Text className={styles.title}>Оплата билета</Text>

          {payment.status === 'pending' && renderPaymentContent()}

          {payment.status === 'failed' && (
            <>
              <Text style={{ color: tokens.colorPaletteRedForeground1 }}>
                Ошибка при оплате. Попробуйте снова или выберите другой способ оплаты.
              </Text>
              <Button onClick={() => navigate('/')} style={{ marginTop: tokens.spacingVerticalM }}>
                Вернуться на главную
              </Button>
            </>
          )}

          {payment.status === 'pending' && (
            <div className={styles.buttonGroup}>
              <Button onClick={() => navigate('/')}>Отмена</Button>
              <Button onClick={handleCheckStatus} disabled={checking}>
                {checking ? <Spinner size="tiny" /> : 'Проверить статус'}
              </Button>
            </div>
          )}
        </Card>
      </div>
    </div>
  );
};
