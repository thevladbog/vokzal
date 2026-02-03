import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Title2,
  Table,
  TableHeader,
  TableRow,
  TableHeaderCell,
  TableBody,
  TableCell,
  makeStyles,
  Spinner,
  Text,
  Input,
  Label,
} from '@fluentui/react-components';
import { scheduleService } from '@/services/schedule';
import type { Trip } from '@/types';
import { formatDate } from '@/utils/format';

const useStyles = makeStyles({
  container: { padding: '24px' },
  header: { marginBottom: '24px' },
  filters: { display: 'flex', gap: '12px', alignItems: 'center', marginBottom: '16px' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
});

export const TripsPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const [date, setDate] = useState(() => new Date().toISOString().slice(0, 10));

  const getStatusLabel = (status: Trip['status']): string => {
    const key = `trips.status_${status}` as const;
    return t(key, { defaultValue: status });
  };

  const { data: trips = [], isLoading, error } = useQuery<Trip[]>({
    queryKey: ['trips', date],
    queryFn: () => scheduleService.getTrips({ date }),
  });

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label={t('trips.loading')} />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>{t('trips.loadError')}</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>{t('trips.title')}</Title2>
      </div>
      <div className={styles.filters}>
        <Label htmlFor="trip-date">{t('trips.date')}</Label>
        <Input
          id="trip-date"
          type="date"
          value={date}
          onChange={(_, v) => setDate(v.value)}
          style={{ width: '160px' }}
        />
      </div>
      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('trips.date')}</TableHeaderCell>
              <TableHeaderCell>{t('trips.routeSchedule')}</TableHeaderCell>
              <TableHeaderCell>{t('trips.status')}</TableHeaderCell>
              <TableHeaderCell>{t('trips.delay')}</TableHeaderCell>
              <TableHeaderCell>{t('trips.platform')}</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {trips.map((trip) => (
              <TableRow key={trip.id}>
                <TableCell>
                  {trip.departure_datetime
                    ? formatDate(trip.departure_datetime)
                    : trip.date ?? trip.created_at}
                </TableCell>
                <TableCell>
                  {trip.schedule_id?.slice(0, 8) ?? trip.id.slice(0, 8)}…
                </TableCell>
                <TableCell>{getStatusLabel(trip.status)}</TableCell>
                <TableCell>
                  {trip.delay_minutes ? t('trips.delayMinutes', { count: trip.delay_minutes }) : '—'}
                </TableCell>
                <TableCell>{trip.platform ?? '—'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {trips.length === 0 && (
          <div style={{ padding: '24px', textAlign: 'center' }}>
            <Text>{t('trips.noTrips')}</Text>
          </div>
        )}
      </Card>
    </div>
  );
};
