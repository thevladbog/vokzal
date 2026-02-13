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
  Field,
  Dropdown,
  Option,
} from '@fluentui/react-components';
import { Edit24Regular } from '@fluentui/react-icons';
import { AppLayout } from '@/components/layout/AppLayout';
import { VokzalDatePicker } from '@/components/common/VokzalDatePicker';
import { useDialogFormStyles } from '@/styles/dialogFormStyles';
import { useDialogActionsStyles } from '@/styles/dialogActionsStyles';
import { scheduleService } from '@/services/schedule';
import type { Trip, Bus, Driver } from '@/types';
import { formatDate } from '@/utils/format';

const TRIP_STATUSES = ['scheduled', 'boarding', 'delayed', 'cancelled', 'departed', 'arrived'] as const;

const useStyles = makeStyles({
  header: { marginBottom: '24px' },
  filters: { display: 'flex', gap: '12px', alignItems: 'center', marginBottom: '16px' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
});

export const TripsPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const formStyles = useDialogFormStyles();
  const actionsStyles = useDialogActionsStyles();
  const queryClient = useQueryClient();
  const [date, setDate] = useState<Date | null>(() => new Date());
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

  // Format date to YYYY-MM-DD using local timezone to avoid UTC offset issues
  const formatLocalDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const { data: trips = [], isLoading, error } = useQuery<Trip[]>({
    queryKey: ['trips', date],
    queryFn: () => {
      if (!date) return Promise.resolve([]);
      const dateStr = formatLocalDate(date);
      return scheduleService.getTrips({ date: dateStr });
    },
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
      <AppLayout>
        <div className={styles.loading}>
          <Spinner label={t('trips.loading')} />
        </div>
      </AppLayout>
    );
  }
  if (error) {
    return (
      <AppLayout>
        <Text>{t('trips.loadError')}</Text>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className={styles.header}>
        <Title2>{t('trips.title')}</Title2>
      </div>
      <div className={styles.filters}>
        <Label htmlFor="trip-date">{t('trips.date')}</Label>
        <VokzalDatePicker
          value={date}
          onSelectDate={(selectedDate) => setDate(selectedDate ?? null)}
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
                <div className={formStyles.formContainer}>
                  <Field label={t('trips.status')} required>
                    <Dropdown
                      value={getStatusLabel(editStatus)}
                      selectedOptions={[editStatus]}
                      onOptionSelect={(_, data) => setEditStatus(data.optionValue ?? '')}
                    >
                      {TRIP_STATUSES.map((s) => (
                        <Option key={s} value={s}>
                          {getStatusLabel(s)}
                        </Option>
                      ))}
                    </Dropdown>
                  </Field>
                  <Field 
                    label={`${t('trips.delay')} (${t('trips.minutes')})`}
                    validationMessage={editDelayError}
                    validationState={editDelayError ? 'error' : undefined}
                  >
                    <Input
                      id="edit-delay"
                      type="text"
                      inputMode="numeric"
                      value={editDelayInput}
                      onChange={(_, v) => handleDelayChange(v.value)}
                    />
                  </Field>
                  <Field label={t('trips.platform')}>
                    <Input
                      id="edit-platform"
                      value={editPlatform}
                      onChange={(_, v) => setEditPlatform(v.value)}
                      placeholder={t('trips.platformPlaceholder')}
                    />
                  </Field>
                  <Field label={t('trips.bus')}>
                    <Dropdown
                      placeholder="—"
                      value={buses.find(b => b.id === editBusId)?.plate_number || ''}
                      selectedOptions={[editBusId]}
                      onOptionSelect={(_, data) => setEditBusId(data.optionValue ?? '')}
                    >
                      <Option value="" text="—">—</Option>
                      {buses.map((b) => (
                        <Option key={b.id} value={b.id} text={`${b.plate_number} (${b.model})`}>
                          {b.plate_number} ({b.model})
                        </Option>
                      ))}
                    </Dropdown>
                  </Field>
                  <Field label={t('trips.driver')}>
                    <Dropdown
                      placeholder="—"
                      value={drivers.find(d => d.id === editDriverId)?.full_name || ''}
                      selectedOptions={[editDriverId]}
                      onOptionSelect={(_, data) => setEditDriverId(data.optionValue ?? '')}
                    >
                      <Option value="" text="—">—</Option>
                      {drivers.map((d) => (
                        <Option key={d.id} value={d.id} text={`${d.full_name} (${d.license_number})`}>
                          {d.full_name} ({d.license_number})
                        </Option>
                      ))}
                    </Dropdown>
                  </Field>
                </div>
              )}
            </DialogContent>
            <DialogActions>
              <div className={actionsStyles.wrapper}>
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
              </div>
            </DialogActions>
          </DialogBody>
        </DialogSurface>
      </Dialog>
    </AppLayout>
  );
};
