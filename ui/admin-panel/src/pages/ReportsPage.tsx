import React from 'react';
import { useTranslation } from 'react-i18next';
import { Title2, Text, Card, makeStyles, Link } from '@fluentui/react-components';

const useStyles = makeStyles({
  container: { padding: '24px' },
  card: { padding: '24px', marginBottom: '16px' },
});

export const ReportsPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const grafanaUrl = import.meta.env.VITE_GRAFANA_URL || 'http://localhost:3000';

  return (
    <div className={styles.container}>
      <Title2 style={{ marginBottom: '24px' }}>{t('reports.title')}</Title2>
      <Card className={styles.card}>
        <Text block style={{ marginBottom: '16px' }}>
          {t('reports.description')}
        </Text>
        <Link href={`${grafanaUrl}/d/business-metrics`} target="_blank" rel="noopener noreferrer">
          {t('reports.openBusinessMetrics')}
        </Link>
      </Card>
      <Card className={styles.card}>
        <Text block style={{ marginBottom: '8px' }}>
          <strong>{t('reports.reportTypes')}</strong>
        </Text>
        <Text block>— {t('reports.typeSales')}</Text>
        <Text block>— {t('reports.typeReturns')}</Text>
        <Text block>— {t('reports.typeZReports')}</Text>
      </Card>
    </div>
  );
};
