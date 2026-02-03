import React from 'react';
import {
  FluentProvider,
  webLightTheme,
  Text,
  Title2,
  Card,
  makeStyles,
} from '@fluentui/react-components';
import { useAuthStore } from '@/stores/authStore';

const useStyles = makeStyles({
  container: {
    padding: '24px',
  },
  card: {
    padding: '24px',
    marginBottom: '16px',
  },
  stats: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
    gap: '16px',
    marginTop: '24px',
  },
  statCard: {
    padding: '20px',
    textAlign: 'center',
  },
});

export const DashboardPage: React.FC = () => {
  const styles = useStyles();
  const user = useAuthStore((state) => state.user);

  return (
    <FluentProvider theme={webLightTheme}>
      <div className={styles.container}>
        <Card className={styles.card}>
          <Title2>Добро пожаловать, {user?.fio || user?.username}!</Title2>
          <Text>Роль: {user?.role}</Text>
        </Card>

        <Title2 style={{ marginBottom: '16px' }}>Статистика за сегодня</Title2>

        <div className={styles.stats}>
          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0</Text>
            <Text block>Рейсов запланировано</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0</Text>
            <Text block>Билетов продано</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0 ₽</Text>
            <Text block>Выручка</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0%</Text>
            <Text block>Заполняемость</Text>
          </Card>
        </div>
      </div>
    </FluentProvider>
  );
};
