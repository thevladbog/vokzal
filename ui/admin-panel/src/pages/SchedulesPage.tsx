import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Title2,
  Button,
  Table,
  TableHeader,
  TableRow,
  TableHeaderCell,
  TableBody,
  TableCell,
  makeStyles,
  Spinner,
  Text,
} from '@fluentui/react-components';
import { Add24Regular } from '@fluentui/react-icons';
import { scheduleService } from '@/services/schedule';
import type { Schedule } from '@/types';

const useStyles = makeStyles({
  container: {
    padding: '24px',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  loading: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '400px',
  },
});

export const SchedulesPage: React.FC = () => {
  const styles = useStyles();

  const { data: schedules, isLoading, error } = useQuery<Schedule[]>({
    queryKey: ['schedules'],
    queryFn: () => scheduleService.getSchedules(),
  });

  const formatDaysOfWeek = (days: number[]) => {
    const dayNames = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];
    return days.map((d) => dayNames[d - 1]).join(', ');
  };

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label="Загрузка расписаний..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <Text>Ошибка загрузки расписаний</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>Расписания</Title2>
        <Button appearance="primary" icon={<Add24Regular />}>
          Создать расписание
        </Button>
      </div>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>Маршрут</TableHeaderCell>
              <TableHeaderCell>Время отправления</TableHeaderCell>
              <TableHeaderCell>Дни недели</TableHeaderCell>
              <TableHeaderCell>Цена</TableHeaderCell>
              <TableHeaderCell>Статус</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {schedules?.map((schedule) => (
              <TableRow key={schedule.id}>
                <TableCell>{schedule.route?.name ?? schedule.route_id}</TableCell>
                <TableCell>{schedule.departure_time}</TableCell>
                <TableCell>{formatDaysOfWeek(schedule.days_of_week ?? [])}</TableCell>
                <TableCell>{schedule.price != null ? `${schedule.price} ₽` : 'N/A'}</TableCell>
                <TableCell>{schedule.is_active ? 'Активно' : 'Неактивно'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>
    </div>
  );
};
