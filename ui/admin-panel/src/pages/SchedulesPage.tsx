import { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
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
  Dialog,
  DialogTrigger,
  DialogSurface,
  DialogTitle,
  DialogBody,
  DialogActions,
  DialogContent,
  Label,
  Input,
  Select,
  Checkbox,
} from '@fluentui/react-components';
import { Add24Regular, CalendarLtr24Regular, Edit24Regular } from '@fluentui/react-icons';
import { scheduleService } from '@/services/schedule';
import type { Schedule, Route } from '@/types';

const DAY_KEYS = ['schedules.dayMon', 'schedules.dayTue', 'schedules.dayWed', 'schedules.dayThu', 'schedules.dayFri', 'schedules.daySat', 'schedules.daySun'] as const;
const DAYS_OF_WEEK = [1, 2, 3, 4, 5, 6, 7] as const;

const useStyles = makeStyles({
  container: {
    padding: '24px',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
    flexWrap: 'wrap',
    gap: '12px',
  },
  loading: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '400px',
  },
  formRow: { marginBottom: '16px' },
  dayCheckbox: { marginRight: '12px' },
});

export const SchedulesPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const queryClient = useQueryClient();
  const [routeFilter, setRouteFilter] = useState<string>('');
  const hasInitializedRoute = useRef(false);
  const [createOpen, setCreateOpen] = useState(false);
  const [editSchedule, setEditSchedule] = useState<Schedule | null>(null);
  const [createRouteId, setCreateRouteId] = useState('');
  const [createDepartureTime, setCreateDepartureTime] = useState('08:00');
  const [createDaysOfWeek, setCreateDaysOfWeek] = useState<number[]>([1, 2, 3, 4, 5]);
  const [createPlatform, setCreatePlatform] = useState('');
  const [editDepartureTime, setEditDepartureTime] = useState('');
  const [editDaysOfWeek, setEditDaysOfWeek] = useState<number[]>([]);
  const [editPlatform, setEditPlatform] = useState('');
  const [editIsActive, setEditIsActive] = useState(true);
  const [generateOpen, setGenerateOpen] = useState(false);
  const [generateRouteId, setGenerateRouteId] = useState('');
  const [generateScheduleId, setGenerateScheduleId] = useState('');
  const [generateFromDate, setGenerateFromDate] = useState(() =>
    new Date().toISOString().slice(0, 10)
  );
  const [generateToDate, setGenerateToDate] = useState(() => {
    const d = new Date();
    d.setDate(d.getDate() + 7);
    return d.toISOString().slice(0, 10);
  });

  const { data: routes = [] } = useQuery<Route[]>({
    queryKey: ['routes'],
    queryFn: () => scheduleService.getRoutes(),
  });

  const { data: schedules = [], isLoading, error } = useQuery<Schedule[]>({
    queryKey: ['schedules', routeFilter],
    queryFn: () => scheduleService.getSchedules({ route_id: routeFilter }),
    enabled: !!routeFilter,
  });

  const { data: schedulesForGenerate = [] } = useQuery<Schedule[]>({
    queryKey: ['schedules', generateRouteId],
    queryFn: () => scheduleService.getSchedules({ route_id: generateRouteId }),
    enabled: !!generateRouteId && generateOpen,
  });

  const createScheduleMutation = useMutation({
    mutationFn: (data: {
      route_id: string;
      departure_time: string;
      days_of_week: number[];
      platform?: string;
    }) =>
      scheduleService.createSchedule({
        route_id: data.route_id,
        departure_time: data.departure_time,
        days_of_week: data.days_of_week,
        is_active: true,
        platform: data.platform || undefined,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schedules'] });
      setCreateOpen(false);
      setCreateRouteId('');
      setCreateDepartureTime('08:00');
      setCreateDaysOfWeek([1, 2, 3, 4, 5]);
      setCreatePlatform('');
    },
  });

  const updateScheduleMutation = useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: Partial<Schedule>;
    }) => scheduleService.updateSchedule(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schedules'] });
      setEditSchedule(null);
    },
  });

  const generateMutation = useMutation({
    mutationFn: ({
      scheduleId,
      fromDate,
      toDate,
    }: {
      scheduleId: string;
      fromDate: string;
      toDate: string;
    }) => scheduleService.generateTrips(scheduleId, fromDate, toDate),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trips'] });
      queryClient.invalidateQueries({ queryKey: ['schedules'] });
      setGenerateOpen(false);
      setGenerateScheduleId('');
      setGenerateRouteId('');
    },
  });

  const formatDaysOfWeek = (days: number[] | string) => {
    if (Array.isArray(days)) {
      return days.map((d: number) => t(DAY_KEYS[d - 1])).join(', ');
    }
    // PostgreSQL may return JSONB as string or already parsed
    if (typeof days === 'string') {
      try {
        const parsed = JSON.parse(days);
        return parsed.map((d: number) => t(DAY_KEYS[d - 1])).join(', ');
      } catch {
        // If parsing fails, try to use as-is or return empty
        return '';
      }
    }
    return '';
  };

  const toggleCreateDay = (day: number) => {
    setCreateDaysOfWeek((prev) =>
      prev.includes(day) ? prev.filter((d) => d !== day) : [...prev, day].sort((a, b) => a - b)
    );
  };

  const toggleEditDay = (day: number) => {
    setEditDaysOfWeek((prev) =>
      prev.includes(day) ? prev.filter((d) => d !== day) : [...prev, day].sort((a, b) => a - b)
    );
  };

  const openEdit = (schedule: Schedule) => {
    setEditSchedule(schedule);
    const timeStr = schedule.departure_time ?? '';
    setEditDepartureTime(timeStr.length >= 5 ? timeStr.slice(0, 5) : timeStr);
    setEditDaysOfWeek(schedule.days_of_week ?? []);
    setEditPlatform(schedule.platform ?? '');
    setEditIsActive(schedule.is_active ?? true);
  };

  const handleCreateSubmit = () => {
    if (!createRouteId || createDaysOfWeek.length === 0) return;
    createScheduleMutation.mutate({
      route_id: createRouteId,
      departure_time: createDepartureTime,
      days_of_week: createDaysOfWeek,
      platform: createPlatform || undefined,
    });
  };

  const handleEditSubmit = () => {
    if (!editSchedule) return;
    updateScheduleMutation.mutate({
      id: editSchedule.id,
      data: {
        departure_time: editDepartureTime,
        days_of_week: editDaysOfWeek,
        platform: editPlatform || undefined,
        is_active: editIsActive,
      },
    });
  };

  // Set first route only on initial load when routes become available; do not overwrite on refetch.
  useEffect(() => {
    if (hasInitializedRoute.current || routes.length === 0) return;
    setRouteFilter(routes[0].id);
    hasInitializedRoute.current = true;
  }, [routes]);

  const isGenerateRangeInvalid = generateFromDate > generateToDate;

  const handleGenerateSubmit = () => {
    if (!generateScheduleId || !generateFromDate || !generateToDate) return;
    if (isGenerateRangeInvalid) return;
    generateMutation.mutate({
      scheduleId: generateScheduleId,
      fromDate: generateFromDate,
      toDate: generateToDate,
    });
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>{t('schedules.title')}</Title2>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
          <Label htmlFor="schedule-route-filter">{t('schedules.route')}:</Label>
          <Select
            id="schedule-route-filter"
            value={routeFilter}
            onChange={(_, data) => setRouteFilter(data.value ?? '')}
            style={{ minWidth: '200px' }}
          >
            {routes.map((r) => (
              <option key={r.id} value={r.id}>
                {r.name}
              </option>
            ))}
          </Select>
          <Dialog open={generateOpen} onOpenChange={(_, v) => setGenerateOpen(v.open)}>
            <DialogTrigger disableButtonEnhancement>
              <Button
                appearance="secondary"
                icon={<CalendarLtr24Regular />}
                onClick={() => setGenerateOpen(true)}
              >
                {t('schedules.generateTrips')}
              </Button>
            </DialogTrigger>
            <DialogSurface>
              <DialogBody>
                <DialogTitle>{t('schedules.generateTripsTitle')}</DialogTitle>
                <DialogContent>
                  <div className={styles.formRow}>
                    <Label>{t('schedules.route')}</Label>
                    <Select
                      value={generateRouteId}
                      onChange={(_, data) => {
                        setGenerateRouteId(data.value ?? '');
                        setGenerateScheduleId('');
                      }}
                      style={{ width: '100%' }}
                    >
                      {routes.map((r) => (
                        <option key={r.id} value={r.id}>
                          {r.name}
                        </option>
                      ))}
                    </Select>
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('schedules.schedule')}</Label>
                    <Select
                      value={generateScheduleId}
                      onChange={(_, data) => setGenerateScheduleId(data.value ?? '')}
                      style={{ width: '100%' }}
                      disabled={!generateRouteId}
                    >
                      {schedulesForGenerate.map((s) => (
                        <option key={s.id} value={s.id}>
                          {s.departure_time} — {formatDaysOfWeek(s.days_of_week ?? [])}
                        </option>
                      ))}
                    </Select>
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="gen-from">{t('schedules.fromDate')}</Label>
                    <Input
                      id="gen-from"
                      type="date"
                      value={generateFromDate}
                      onChange={(_, v) => setGenerateFromDate(v.value)}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="gen-to">{t('schedules.toDate')}</Label>
                    <Input
                      id="gen-to"
                      type="date"
                      value={generateToDate}
                      onChange={(_, v) => setGenerateToDate(v.value)}
                    />
                  </div>
                </DialogContent>
                <DialogActions>
                  <DialogTrigger disableButtonEnhancement>
                    <Button appearance="secondary">{t('common.cancel')}</Button>
                  </DialogTrigger>
                  <Button
                    appearance="primary"
                    onClick={handleGenerateSubmit}
                    disabled={
                      !generateScheduleId || !generateFromDate || !generateToDate || isGenerateRangeInvalid || generateMutation.isPending
                    }
                  >
                    {generateMutation.isPending ? t('schedules.creating') : t('schedules.generateTrips')}
                  </Button>
                </DialogActions>
              </DialogBody>
            </DialogSurface>
          </Dialog>
          <Dialog open={createOpen} onOpenChange={(_, v) => setCreateOpen(v.open)}>
            <DialogTrigger disableButtonEnhancement>
              <Button appearance="primary" icon={<Add24Regular />} onClick={() => setCreateOpen(true)}>
                {t('schedules.createSchedule')}
              </Button>
            </DialogTrigger>
            <DialogSurface>
              <DialogBody>
                <DialogTitle>{t('schedules.createScheduleTitle')}</DialogTitle>
                <DialogContent>
                  <div className={styles.formRow}>
                    <Label>{t('schedules.route')}</Label>
                    <Select
                      value={createRouteId}
                      onChange={(_, data) => setCreateRouteId(data.value ?? '')}
                      style={{ width: '100%' }}
                    >
                      {routes.map((r) => (
                        <option key={r.id} value={r.id}>
                          {r.name}
                        </option>
                      ))}
                    </Select>
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="create-time">{t('schedules.departureTime')}</Label>
                    <Input
                      id="create-time"
                      type="time"
                      value={createDepartureTime}
                      onChange={(_, v) => setCreateDepartureTime(v.value)}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('schedules.daysOfWeek')}</Label>
                    <div style={{ display: 'flex', flexWrap: 'wrap' }}>
                      {DAYS_OF_WEEK.map((day) => (
                        <Checkbox
                          key={day}
                          className={styles.dayCheckbox}
                          label={t(DAY_KEYS[day - 1])}
                          checked={createDaysOfWeek.includes(day)}
                          onChange={() => toggleCreateDay(day)}
                        />
                      ))}
                    </div>
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="create-platform">{t('schedules.platformOptional')}</Label>
                    <Input
                      id="create-platform"
                      value={createPlatform}
                      onChange={(_, v) => setCreatePlatform(v.value)}
                      placeholder={t('schedules.platformPlaceholder')}
                    />
                  </div>
                </DialogContent>
                <DialogActions>
                  <DialogTrigger disableButtonEnhancement>
                    <Button appearance="secondary">{t('common.cancel')}</Button>
                  </DialogTrigger>
                  <Button
                    appearance="primary"
                    onClick={handleCreateSubmit}
                    disabled={
                      !createRouteId || createDaysOfWeek.length === 0 || createScheduleMutation.isPending
                    }
                  >
                    {createScheduleMutation.isPending ? t('schedules.creating') : t('common.create')}
                  </Button>
                </DialogActions>
              </DialogBody>
            </DialogSurface>
          </Dialog>
        </div>
      </div>

      <Dialog open={!!editSchedule} onOpenChange={(_, v) => (!v.open && setEditSchedule(null))}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>{t('schedules.editScheduleTitle')}</DialogTitle>
            <DialogContent>
              {editSchedule && (
                <>
                  <div className={styles.formRow}>
                    <Label htmlFor="edit-time">{t('schedules.departureTime')}</Label>
                    <Input
                      id="edit-time"
                      type="time"
                      value={editDepartureTime}
                      onChange={(_, v) => setEditDepartureTime(v.value)}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('schedules.daysOfWeek')}</Label>
                    <div style={{ display: 'flex', flexWrap: 'wrap' }}>
                      {DAYS_OF_WEEK.map((day) => (
                        <Checkbox
                          key={day}
                          className={styles.dayCheckbox}
                          label={t(DAY_KEYS[day - 1])}
                          checked={editDaysOfWeek.includes(day)}
                          onChange={() => toggleEditDay(day)}
                        />
                      ))}
                    </div>
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="edit-platform">{t('schedules.platformOptional')}</Label>
                    <Input
                      id="edit-platform"
                      value={editPlatform}
                      onChange={(_, v) => setEditPlatform(v.value)}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Checkbox
                      label={t('schedules.active')}
                      checked={editIsActive}
                      onChange={(_, v) => setEditIsActive(!!v.checked)}
                    />
                  </div>
                </>
              )}
            </DialogContent>
            <DialogActions>
              <Button appearance="secondary" onClick={() => setEditSchedule(null)}>
                {t('common.cancel')}
              </Button>
              <Button
                appearance="primary"
                onClick={handleEditSubmit}
                disabled={
                  !editSchedule ||
                  editDaysOfWeek.length === 0 ||
                  updateScheduleMutation.isPending
                }
              >
                {updateScheduleMutation.isPending ? t('schedules.saving') : t('common.save')}
              </Button>
            </DialogActions>
          </DialogBody>
        </DialogSurface>
      </Dialog>

      {!routeFilter && routes.length === 0 && !isLoading && (
        <Text>{t('schedules.noSchedulesHint')}</Text>
      )}

      {routeFilter && isLoading && (
        <div className={styles.loading}>
          <Spinner label={t('schedules.loading')} />
        </div>
      )}

      {routeFilter && error && (
        <Text>{t('schedules.loadError')}</Text>
      )}

      {routeFilter && !isLoading && !error && (
        <Card>
          <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('schedules.route')}</TableHeaderCell>
              <TableHeaderCell>{t('schedules.departureTime')}</TableHeaderCell>
              <TableHeaderCell>{t('schedules.daysOfWeek')}</TableHeaderCell>
              <TableHeaderCell>{t('schedules.price')}</TableHeaderCell>
              <TableHeaderCell>{t('schedules.status')}</TableHeaderCell>
              <TableHeaderCell></TableHeaderCell>
            </TableRow>
          </TableHeader>
            <TableBody>
              {schedules.map((schedule) => (
                <TableRow key={schedule.id}>
                  <TableCell>{schedule.route?.name ?? schedule.route_id}</TableCell>
                  <TableCell>{schedule.departure_time}</TableCell>
                  <TableCell>{formatDaysOfWeek(schedule.days_of_week ?? [])}</TableCell>
                <TableCell>{schedule.price != null ? `${schedule.price} ₽` : t('common.notAvailable')}</TableCell>
                <TableCell>{schedule.is_active ? t('schedules.active') : t('schedules.inactive')}</TableCell>
                  <TableCell>
                    <Button
                      appearance="subtle"
                      icon={<Edit24Regular />}
                      aria-label={t('common.edit')}
                      onClick={() => openEdit(schedule)}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}
    </div>
  );
};
