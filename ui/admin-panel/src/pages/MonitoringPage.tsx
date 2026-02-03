import React from 'react';
import { Title2, Text, Card, makeStyles, Link } from '@fluentui/react-components';

const useStyles = makeStyles({
  container: { padding: '24px' },
  card: { padding: '24px', marginBottom: '16px' },
});

export const MonitoringPage: React.FC = () => {
  const styles = useStyles();
  const grafanaUrl = import.meta.env.VITE_GRAFANA_URL || 'http://localhost:3000';

  return (
    <div className={styles.container}>
      <Title2 style={{ marginBottom: '24px' }}>Мониторинг</Title2>
      <Card className={styles.card}>
        <Text block style={{ marginBottom: '16px' }}>
          Обзор состояния сервисов и инфраструктуры — в Grafana.
        </Text>
        <Link href={`${grafanaUrl}/d/services-overview`} target="_blank" rel="noopener noreferrer">
          Открыть дашборд «Обзор сервисов»
        </Link>
      </Card>
      <Card className={styles.card}>
        <Link href={`${grafanaUrl}/d/database-monitoring`} target="_blank" rel="noopener noreferrer">
          Мониторинг БД
        </Link>
      </Card>
    </div>
  );
};
