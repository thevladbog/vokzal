import React from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import {
  Text,
  Title2,
  Card,
  makeStyles,
  Spinner,
  tokens,
  shorthands,
} from '@fluentui/react-components';
import { AppLayout } from '@/components/layout/AppLayout';
import { useAuthStore } from '@/stores/authStore';
import { scheduleService } from '@/services/schedule';
import { ticketService } from '@/services/ticket';
import i18n from '@/i18n';

const useStyles = makeStyles({
  welcomeCard: {
    ...shorthands.padding('24px'),
    marginBottom: '24px',
  },
  stats: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
    ...shorthands.gap('16px'),
  },
  statCard: {
    ...shorthands.padding('20px'),
    textAlign: 'center',
    ...shorthands.transition('transform', '150ms', 'ease-in-out'),
    ':hover': {
      transform: 'translateY(-2px)',
      boxShadow: tokens.shadow8,
    },
  },
  statValue: {
    display: 'block',
    fontSize: tokens.fontSizeHero900,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorBrandForeground1,
    marginBottom: '8px',
  },
  statLabel: {
    color: tokens.colorNeutralForeground3,
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
    <AppLayout>
      <Card className={styles.welcomeCard}>
        <Title2>
          {t('dashboard.welcome')}, {user?.full_name || user?.fio || user?.username}!
        </Title2>
        <Text>{t('dashboard.role')}: {user?.role}</Text>
      </Card>

      <Title2 style={{ marginBottom: '16px' }}>{t('dashboard.statsTitle')}</Title2>

      {isLoading ? (
        <Spinner label={t('dashboard.loading')} />
      ) : (
        <div className={styles.stats}>
          <Card className={styles.statCard}>
            <Text className={styles.statValue}>{tripsTotal}</Text>
            <Text className={styles.statLabel}>{t('dashboard.statsTrips')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text className={styles.statValue}>{ticketsSold}</Text>
            <Text className={styles.statLabel}>{t('dashboard.statsTickets')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text className={styles.statValue}>{revenue.toLocaleString(i18n.language === 'en' ? 'en-US' : 'ru-RU')} ₽</Text>
            <Text className={styles.statLabel}>{t('dashboard.statsRevenue')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text className={styles.statValue}>{ticketsReturned}</Text>
            <Text className={styles.statLabel}>{t('dashboard.statsReturns')}</Text>
          </Card>

          <Card className={styles.statCard}>
            <Text className={styles.statValue}>{occupancyPercent}%</Text>
            <Text className={styles.statLabel}>{t('dashboard.statsOccupancy')}</Text>
          </Card>
        </div>
      )}
    </AppLayout>
  );
};
