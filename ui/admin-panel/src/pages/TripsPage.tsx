import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
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
  Button,
  Dialog,
  DialogSurface,
  DialogTitle,
  DialogBody,
  DialogActions,
  DialogContent,
  Select,
  Option,
} from '@fluentui/react-components';
import { Edit24Regular } from '@fluentui/react-icons';
import { scheduleService } from '@/services/schedule';
import type { Trip, Bus, Driver } from '@/types';
import { formatDate } from '@/utils/format';

const TRIP_STATUSES = ['scheduled', 'boarding', 'delayed', 'cancelled', 'departed', 'arrived'] as const;

const useStyles = makeStyles({
  container: { padding: '24px' },
  header: { marginBottom: '24px' },
  filters: { display: 'flex', gap: '12px', alignItems: 'center', marginBottom: '16px' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
  formRow: { marginBottom: '16px' },
});

export const TripsPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const queryClient = useQueryClient();
  const [date, setDate] = useState(() => new Date().toISOString().slice(0, 10));
  const [editTrip, setEditTrip] = useState<Trip | null>(null);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editStatus, setEditStatus] = useState<string>('');
  const [editDelayInput, setEditDelayInput] = useState<string>('0');
  const [editDelayError, setEditDelayError] = useState<string>('');
  const [editPlatform, setEditPlatform] = useState('');
  const [editBusId, setEditBusId] = useState('');
  const [editDriverId, setEditDriverId] = useState('');

  // Fetched on page load so edit dialog Selects have data ready when opened.
  const { data: buses = [] } = useQuery<Bus[]>({
    queryKey: ['buses'],
    queryFn: () => scheduleService.getBuses({}),
  });

  const { data: drivers = [] } = useQuery<Driver[]>({
    queryKey: ['drivers'],
    queryFn: () => scheduleService.getDrivers({}),
  });

  const getStatusLabel = (status: string): string => {
    const key = `trips.status_${status}` as const;
    return t(key, { defaultValue: status });
  };

  const { data: trips = [], isLoading, error } = useQuery<Trip[]>({
    queryKey: ['trips', date],
    queryFn: () => scheduleService.getTrips({ date }),
  });

  // Two API calls (status/delay then platform/bus/driver); on second failure we rollback the first.
  const updateTripMutation = useMutation({
    mutationFn: async (payload: {
      id: string;
      status: string;
      delay_minutes: number;
      platform?: string;
      bus_id?: string | null;
      driver_id?: string | null;
      previous_status: string;
      previous_delay_minutes: number;
    }) => {
      await scheduleService.updateTripStatus(payload.id, {
        status: payload.status,
        delay_minutes: payload.delay_minutes,
      });
      try {
        await scheduleService.updateTrip(payload.id, {
          platform: payload.platform || undefined,
          bus_id: payload.bus_id ?? undefined,
          driver_id: payload.driver_id ?? undefined,
        });
      } catch (err) {
        // Compensating rollback: restore original status/delay so trip state stays consistent.
        await scheduleService.updateTripStatus(payload.id, {
          status: payload.previous_status,
          delay_minutes: payload.previous_delay_minutes,
        }).catch(() => {});
        throw err;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trips'] });
      setEditDialogOpen(false);
      setEditTrip(null);
    },
  });

  const parseDelay = (raw: string): { valid: boolean; value: number } => {
    const trimmed = raw.trim();
    if (trimmed === '') return { valid: true, value: 0 };
    const n = parseInt(trimmed, 10);
    if (Number.isNaN(n) || n < 0 || !/^\d+$/.test(trimmed)) {
      return { valid: false, value: 0 };
    }
    return { valid: true, value: n };
  };

  const openEdit = (trip: Trip) => {
    const delay = trip.delay_minutes ?? 0;
    setEditTrip(trip);
    setEditStatus(trip.status);
    setEditDelayInput(String(delay));
    setEditDelayError('');
    setEditPlatform(trip.platform ?? '');
    setEditBusId(trip.bus_id ?? '');
    setEditDriverId(trip.driver_id ?? '');
    setEditDialogOpen(true);
  };

  const closeEditDialog = () => {
    setEditDialogOpen(false);
    setEditTrip(null);
  };

  const handleDelayChange = (value: string) => {
    setEditDelayInput(value);
    const { valid } = parseDelay(value);
    if (valid) {
      setEditDelayError('');
    } else {
      setEditDelayError(t('trips.delayInvalid'));
    }
  };

  const handleEditSubmit = () => {
    if (!editTrip) return;
    const { valid, value: delayMinutes } = parseDelay(editDelayInput);
    if (!valid) {
      setEditDelayError(t('trips.delayInvalid'));
      return;
    }
    updateTripMutation.mutate({
      id: editTrip.id,
      status: editStatus,
      delay_minutes: delayMinutes,
      platform: editPlatform.trim() || undefined,
      bus_id: editBusId === '' ? null : editBusId,
      driver_id: editDriverId === '' ? null : editDriverId,
      previous_status: editTrip.status,
      previous_delay_minutes: editTrip.delay_minutes ?? 0,
    });
  };

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
              <TableHeaderCell></TableHeaderCell>
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
                <TableCell>
                  <Button
                    appearance="subtle"
                    icon={<Edit24Regular />}
                    aria-label={t('common.edit')}
                    onClick={() => openEdit(trip)}
                  />
                </TableCell>
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

      <Dialog open={editDialogOpen} onOpenChange={(_, v) => (!v.open && closeEditDialog())}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>{t('trips.editTrip')}</DialogTitle>
            <DialogContent>
              {editTrip && (
                <>
                  <div className={styles.formRow}>
                    <Label>{t('trips.status')}</Label>
                    <Select
                      value={editStatus}
                      onChange={(_, data) => setEditStatus(data.value ?? '')}
                      style={{ width: '100%' }}
                    >
                      {TRIP_STATUSES.map((s) => (
                        <Option key={s} value={s} text={getStatusLabel(s)}>
                          {getStatusLabel(s)}
                        </Option>
                      ))}
                    </Select>
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="edit-delay">
                      {t('trips.delay')} ({t('trips.minutes')})
                    </Label>
                    <Input
                      id="edit-delay"
                      type="text"
                      inputMode="numeric"
                      value={editDelayInput}
                      onChange={(_, v) => handleDelayChange(v.value)}
                      aria-invalid={!!editDelayError}
                      aria-describedby={editDelayError ? 'edit-delay-error' : undefined}
                    />
                    {editDelayError && (
                      <Text id="edit-delay-error" style={{ color: 'var(--colorPaletteRedForeground1)', fontSize: '12px', marginTop: '4px' }}>
                        {editDelayError}
                      </Text>
                    )}
                  </div>
                  <div className={styles.formRow}>
                    <Label htmlFor="edit-platform">{t('trips.platform')}</Label>
                    <Input
                      id="edit-platform"
                      value={editPlatform}
                      onChange={(_, v) => setEditPlatform(v.value)}
                      placeholder={t('trips.platformPlaceholder')}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('trips.bus')}</Label>
                    <Select
                      value={editBusId}
                      onChange={(_, d) => setEditBusId(d.value ?? '')}
                      style={{ width: '100%' }}
                    >
                      <Option value="" text="—">—</Option>
                      {buses.map((b) => (
                        <Option key={b.id} value={b.id} text={`${b.plate_number} (${b.model})`}>
                          {b.plate_number} ({b.model})
                        </Option>
                      ))}
                    </Select>
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('trips.driver')}</Label>
                    <Select
                      value={editDriverId}
                      onChange={(_, d) => setEditDriverId(d.value ?? '')}
                      style={{ width: '100%' }}
                    >
                      <Option value="" text="—">—</Option>
                      {drivers.map((d) => (
                        <Option key={d.id} value={d.id} text={`${d.full_name} (${d.license_number})`}>
                          {d.full_name} ({d.license_number})
                        </Option>
                      ))}
                    </Select>
                  </div>
                </>
              )}
            </DialogContent>
            <DialogActions>
              <Button appearance="secondary" onClick={closeEditDialog}>
                {t('common.cancel')}
              </Button>
              <Button
                appearance="primary"
                onClick={handleEditSubmit}
                disabled={!editTrip || updateTripMutation.isPending || !!editDelayError}
              >
                {updateTripMutation.isPending ? t('common.saving') : t('common.save')}
              </Button>
            </DialogActions>
          </DialogBody>
        </DialogSurface>
      </Dialog>
    </div>
  );
};
