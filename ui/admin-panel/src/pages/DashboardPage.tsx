import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
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
  Spinner,
} from '@fluentui/react-components';
import { useAuthStore } from '@/stores/authStore';
import { scheduleService } from '@/services/schedule';
import { ticketService } from '@/services/ticket';
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

const today = () => new Date().toISOString().slice(0, 10);

export const DashboardPage: React.FC = () => {
  const styles = useStyles();
  const { t } = useTranslation();
  const user = useAuthStore((state) => state.user);
  const todayStr = today();

  const { data: scheduleStats, isLoading: scheduleLoading } = useQuery({
    queryKey: ['dashboard', 'schedule', todayStr],
    queryFn: () => scheduleService.getDashboardStats(todayStr),
  });

  const { data: ticketStats, isLoading: ticketLoading } = useQuery({
    queryKey: ['dashboard', 'ticket', todayStr],
    queryFn: () => ticketService.getDashboardStats(todayStr),
  });

  const isLoading = scheduleLoading || ticketLoading;
  const tripsTotal = scheduleStats?.trips_total ?? 0;
  const totalCapacity = scheduleStats?.total_capacity ?? 0;
  const ticketsSold = ticketStats?.tickets_sold ?? 0;
  const ticketsReturned = ticketStats?.tickets_returned ?? 0;
  const revenue = ticketStats?.revenue ?? 0;
  const totalSeats = totalCapacity > 0 ? totalCapacity : (tripsTotal > 0 ? tripsTotal * 40 : 40);
  const occupancyPercent =
    (ticketStats as { occupancy?: number } | undefined)?.occupancy ??
    (scheduleStats as { occupancy?: number } | undefined)?.occupancy ??
    (totalSeats > 0 ? Math.round((ticketsSold / totalSeats) * 100) : 0);

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
          <Link to="/buses">
            <Button appearance="secondary">{t('nav.buses')}</Button>
          </Link>
          <Link to="/drivers">
            <Button appearance="secondary">{t('nav.drivers')}</Button>
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
            aria-label={t('common.language')}
          >
            <Option value="ru" text="RU">RU</Option>
            <Option value="en" text="EN">EN</Option>
          </Select>
        </div>

        <Title2 style={{ marginBottom: '16px' }}>{t('dashboard.statsTitle')}</Title2>

        {isLoading ? (
          <Spinner label={t('dashboard.loading')} />
        ) : (
          <div className={styles.stats}>
            <Card className={styles.statCard}>
              <Text size={600} weight="bold">{tripsTotal}</Text>
              <Text block>{t('dashboard.statsTrips')}</Text>
            </Card>

            <Card className={styles.statCard}>
              <Text size={600} weight="bold">{ticketsSold}</Text>
              <Text block>{t('dashboard.statsTickets')}</Text>
            </Card>

            <Card className={styles.statCard}>
              <Text size={600} weight="bold">{revenue.toLocaleString(i18n.language === 'en' ? 'en-US' : 'ru-RU')} ₽</Text>
              <Text block>{t('dashboard.statsRevenue')}</Text>
            </Card>

            <Card className={styles.statCard}>
              <Text size={600} weight="bold">{ticketsReturned}</Text>
              <Text block>{t('dashboard.statsReturns')}</Text>
            </Card>

            <Card className={styles.statCard}>
              <Text size={600} weight="bold">{occupancyPercent}%</Text>
              <Text block>{t('dashboard.statsOccupancy')}</Text>
            </Card>
          </div>
        )}
      </div>
    </FluentProvider>
  );
};
