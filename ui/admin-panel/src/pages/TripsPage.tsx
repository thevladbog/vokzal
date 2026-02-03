import React, { useState } from 'react';
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

const TRIP_STATUS_LABELS: Record<string, string> = {
  scheduled: 'Запланирован',
  delayed: 'Задержан',
  cancelled: 'Отменён',
  departed: 'Отправлен',
  arrived: 'Прибыл',
};

export const TripsPage: React.FC = () => {
  const styles = useStyles();
  const [date, setDate] = useState(() => new Date().toISOString().slice(0, 10));

  const { data: trips = [], isLoading, error } = useQuery<Trip[]>({
    queryKey: ['trips', date],
    queryFn: () => scheduleService.getTrips({ date }),
  });

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label="Загрузка рейсов..." />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>Ошибка загрузки рейсов</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>Рейсы</Title2>
      </div>
      <div className={styles.filters}>
        <Label htmlFor="trip-date">Дата</Label>
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
              <TableHeaderCell>Дата</TableHeaderCell>
              <TableHeaderCell>Маршрут / расписание</TableHeaderCell>
              <TableHeaderCell>Статус</TableHeaderCell>
              <TableHeaderCell>Задержка</TableHeaderCell>
              <TableHeaderCell>Перрон</TableHeaderCell>
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
                <TableCell>{TRIP_STATUS_LABELS[trip.status] ?? trip.status}</TableCell>
                <TableCell>{trip.delay_minutes ? `${trip.delay_minutes} мин` : '—'}</TableCell>
                <TableCell>{trip.platform ?? '—'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {trips.length === 0 && (
          <div style={{ padding: '24px', textAlign: 'center' }}>
            <Text>Рейсов на выбранную дату нет</Text>
          </div>
        )}
      </Card>
    </div>
  );
};
