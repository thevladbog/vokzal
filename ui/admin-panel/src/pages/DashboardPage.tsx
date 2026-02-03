import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  FluentProvider,
  webLightTheme,
  Text,
  Title2,
  Card,
  makeStyles,
  Button,
  Select,
  Option,
} from '@fluentui/react-components';
import { useAuthStore } from '@/stores/authStore';
import i18n from '@/i18n';

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
  nav: {
    display: 'flex',
    gap: '12px',
    marginBottom: '24px',
    flexWrap: 'wrap',
  },
  langSwitcher: {
    marginLeft: 'auto',
    minWidth: '100px',
  },
});

export const DashboardPage: React.FC = () => {
  const styles = useStyles();
  const { t } = useTranslation();
  const user = useAuthStore((state) => state.user);

  return (
    <FluentProvider theme={webLightTheme}>
      <div className={styles.container}>
        <Card className={styles.card}>
          <Title2>
            {t('dashboard.welcome')}, {user?.full_name || user?.fio || user?.username}!
          </Title2>
          <Text>{t('dashboard.role')}: {user?.role}</Text>
        </Card>

        <div className={styles.nav}>
          <Link to="/schedules">
            <Button appearance="secondary">{t('nav.schedules')}</Button>
          </Link>
          <Link to="/stations">
            <Button appearance="secondary">{t('nav.stations')}</Button>
          </Link>
          <Link to="/routes">
            <Button appearance="secondary">{t('nav.routes')}</Button>
          </Link>
          <Link to="/trips">
            <Button appearance="secondary">{t('nav.trips')}</Button>
          </Link>
          <Link to="/reports">
            <Button appearance="secondary">{t('nav.reports')}</Button>
          </Link>
          <Link to="/monitoring">
            <Button appearance="secondary">{t('nav.monitoring')}</Button>
          </Link>
          {user?.role === 'admin' && (
            <>
              <Link to="/users">
                <Button appearance="secondary">{t('nav.users')}</Button>
              </Link>
              <Link to="/audit">
                <Button appearance="secondary">{t('nav.audit')}</Button>
              </Link>
            </>
          )}
          <Select
            value={i18n.language}
            onChange={(_, v) => v.value && i18n.changeLanguage(v.value)}
            className={styles.langSwitcher}
            aria-label="Язык"
          >
            <Option value="ru" text="RU">RU</Option>
            <Option value="en" text="EN">EN</Option>
          </Select>
        </div>

        <Title2 style={{ marginBottom: '16px' }}>Статистика за сегодня</Title2>

        <div className={styles.stats}>
          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0</Text>
            <Text block>{t('dashboard.statsTrips')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0</Text>
            <Text block>{t('dashboard.statsTickets')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0 ₽</Text>
            <Text block>{t('dashboard.statsRevenue')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text size={600} weight="bold">0%</Text>
            <Text block>{t('dashboard.statsOccupancy')}</Text>
          </Card>
        </div>
      </div>
    </FluentProvider>
  );
};
