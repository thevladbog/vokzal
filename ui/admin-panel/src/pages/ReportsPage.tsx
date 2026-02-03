import React from 'react';
import { Title2, Text, Card, makeStyles, Link } from '@fluentui/react-components';

const useStyles = makeStyles({
  container: { padding: '24px' },
  card: { padding: '24px', marginBottom: '16px' },
});

export const ReportsPage: React.FC = () => {
  const styles = useStyles();
  const grafanaUrl = import.meta.env.VITE_GRAFANA_URL || 'http://localhost:3000';

  return (
    <div className={styles.container}>
      <Title2 style={{ marginBottom: '24px' }}>Отчёты</Title2>
      <Card className={styles.card}>
        <Text block style={{ marginBottom: '16px' }}>
          Сводные отчёты по продажам, возвратам и фискализации доступны в Grafana.
        </Text>
        <Link href={`${grafanaUrl}/d/business-metrics`} target="_blank" rel="noopener noreferrer">
          Открыть дашборд «Бизнес-метрики»
        </Link>
      </Card>
      <Card className={styles.card}>
        <Text block style={{ marginBottom: '8px' }}>
          <strong>Типы отчётов:</strong>
        </Text>
        <Text block>— Продажи за период</Text>
        <Text block>— Возвраты</Text>
        <Text block>— Z-отчёты (54-ФЗ)</Text>
      </Card>
    </div>
  );
};
