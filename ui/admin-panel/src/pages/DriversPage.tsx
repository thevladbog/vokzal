import React, { useState } from 'react';
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
  Input,
  Label,
  Field,
  Dropdown,
  Option,
} from '@fluentui/react-components';
import { Add24Regular, Delete24Regular, Edit24Regular } from '@fluentui/react-icons';
import { AppLayout } from '@/components/layout/AppLayout';
import { useDialogFormStyles } from '@/styles/dialogFormStyles';
import { useDialogActionsStyles } from '@/styles/dialogActionsStyles';
import { scheduleService } from '@/services/schedule';
import type { Driver, Station } from '@/types';

const useStyles = makeStyles({
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  filters: { display: 'flex', gap: '12px', alignItems: 'center', marginBottom: '16px' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
  actions: { display: 'flex', gap: '8px' },
});

export const DriversPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const formStyles = useDialogFormStyles();
  const actionsStyles = useDialogActionsStyles();
  const queryClient = useQueryClient();
  const [stationFilter, setStationFilter] = useState('');
  const [createOpen, setCreateOpen] = useState(false);
  const [editDriver, setEditDriver] = useState<Driver | null>(null);
  const [deleteDriver, setDeleteDriver] = useState<Driver | null>(null);
  const [createFullName, setCreateFullName] = useState('');
  const [createLicense, setCreateLicense] = useState('');
  const [createExperience, setCreateExperience] = useState<number | ''>('');
  const [createPhone, setCreatePhone] = useState('');
  const [createStationId, setCreateStationId] = useState('');
  const [editFullName, setEditFullName] = useState('');
  const [editLicense, setEditLicense] = useState('');
  const [editExperience, setEditExperience] = useState<number | ''>('');
  const [editPhone, setEditPhone] = useState('');

  const { data: stations = [] } = useQuery<Station[]>({
    queryKey: ['stations'],
    queryFn: () => scheduleService.getStations(),
  });

  const { data: drivers = [], isLoading, error } = useQuery<Driver[]>({
    queryKey: ['drivers', stationFilter],
    queryFn: () =>
      scheduleService.getDrivers(stationFilter ? { station_id: stationFilter } : {}),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<Driver, 'id' | 'created_at' | 'updated_at'>) =>
      scheduleService.createDriver(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['drivers'] });
      setCreateOpen(false);
      setCreateFullName('');
      setCreateLicense('');
      setCreateExperience('');
      setCreatePhone('');
      setCreateStationId('');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Driver> }) =>
      scheduleService.updateDriver(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['drivers'] });
      setEditDriver(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => scheduleService.deleteDriver(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['drivers'] });
      setDeleteDriver(null);
    },
  });

  const stationName = (id: string) => stations.find((s) => s.id === id)?.name ?? id.slice(0, 8);

  const openEdit = (driver: Driver) => {
    setEditDriver(driver);
    setEditFullName(driver.full_name);
    setEditLicense(driver.license_number);
    setEditExperience(driver.experience_years ?? '');
    setEditPhone(driver.phone ?? '');
  };

  if (isLoading) {
    return (
      <AppLayout>
        <div className={styles.loading}>
          <Spinner label={t('drivers.loading')} />
        </div>
      </AppLayout>
    );
  }
  if (error) {
    return (
      <AppLayout>
        <Text>{t('drivers.loadError')}</Text>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className={styles.header}>
        <Title2>{t('drivers.title')}</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              {t('drivers.addDriver')}
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('drivers.createDriverTitle')}</DialogTitle>
              <DialogContent>
                <div className={formStyles.formContainer}>
                  <Field label={t('drivers.fullName')} required>
                    <Input
                      value={createFullName}
                      onChange={(_, v) => setCreateFullName(v.value)}
                      placeholder={t('drivers.fullNamePlaceholder')}
                    />
                  </Field>
                  <Field label={t('drivers.licenseNumber')} required>
                    <Input
                      value={createLicense}
                      onChange={(_, v) => setCreateLicense(v.value)}
                      placeholder={t('drivers.licensePlaceholder')}
                    />
                  </Field>
                  <Field label={t('drivers.experience')} required>
                    <Input
                      type="number"
                      min={0}
                      value={String(createExperience)}
                      onChange={(_, v) =>
                        setCreateExperience(v.value === '' ? '' : Math.max(0, parseInt(v.value, 10) || 0))
                      }
                    />
                  </Field>
                  <Field label={t('drivers.phone')}>
                    <Input
                      value={createPhone}
                      onChange={(_, v) => setCreatePhone(v.value)}
                      placeholder={t('drivers.phonePlaceholder')}
                    />
                  </Field>
                  <Field label={t('drivers.station')} required>
                    <Dropdown
                      placeholder={t('drivers.selectStation')}
                      value={stations.find(s => s.id === createStationId)?.name || ''}
                      selectedOptions={[createStationId]}
                      onOptionSelect={(_, d) => setCreateStationId(d.optionValue ?? '')}
                    >
                      {stations.map((s) => (
                        <Option key={s.id} value={s.id}>
                          {s.name}
                        </Option>
                      ))}
                    </Dropdown>
                  </Field>
                </div>
              </DialogContent>
              <DialogActions>
                <div className={actionsStyles.wrapper}>
                  <Button appearance="secondary" onClick={() => setCreateOpen(false)}>
                    {t('common.cancel')}
                  </Button>
                  <Button
                    appearance="primary"
                  disabled={!createFullName.trim() || !createLicense.trim() || !createStationId}
                  onClick={() =>
                    createMutation.mutate({
                      full_name: createFullName.trim(),
                      license_number: createLicense.trim(),
                      experience_years:
                        createExperience === '' ? undefined : (createExperience as number),
                      phone: createPhone.trim() || undefined,
                      station_id: createStationId,
                    })
                  }
                >
                  {t('common.create')}
                </Button>
                </div>
              </DialogActions>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      </div>

      <div className={styles.filters}>
        <Label>{t('drivers.station')}:</Label>
        <Dropdown
          placeholder={t('drivers.allStations')}
          value={stationFilter ? stations.find(s => s.id === stationFilter)?.name : t('drivers.allStations')}
          selectedOptions={[stationFilter]}
          onOptionSelect={(_, d) => setStationFilter(d.optionValue ?? '')}
          style={{ minWidth: '220px' }}
        >
          <Option value="">{t('drivers.allStations')}</Option>
          {stations.map((s) => (
            <Option key={s.id} value={s.id}>
              {s.name}
            </Option>
          ))}
        </Dropdown>
      </div>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('drivers.fullName')}</TableHeaderCell>
              <TableHeaderCell>{t('drivers.licenseNumber')}</TableHeaderCell>
              <TableHeaderCell>{t('drivers.experience')}</TableHeaderCell>
              <TableHeaderCell>{t('drivers.phone')}</TableHeaderCell>
              <TableHeaderCell>{t('drivers.station')}</TableHeaderCell>
              <TableHeaderCell></TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {drivers.map((d) => (
              <TableRow key={d.id}>
                <TableCell>{d.full_name}</TableCell>
                <TableCell>{d.license_number}</TableCell>
                <TableCell>{d.experience_years ?? '—'}</TableCell>
                <TableCell>{d.phone ?? '—'}</TableCell>
                <TableCell>{stationName(d.station_id)}</TableCell>
                <TableCell>
                  <Button
                    appearance="subtle"
                    icon={<Edit24Regular />}
                    onClick={() => openEdit(d)}
                    aria-label={t('common.edit')}
                  />
                  <Button
                    appearance="subtle"
                    icon={<Delete24Regular />}
                    onClick={() => setDeleteDriver(d)}
                    aria-label={t('common.delete')}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>

      <Dialog open={!!editDriver} onOpenChange={(_, d) => (!d.open && setEditDriver(null))}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>{t('drivers.editDriverTitle')}</DialogTitle>
            <DialogContent>
              {editDriver && (
                <div className={formStyles.formContainer}>
                  <Field label={t('drivers.fullName')} required>
                    <Input value={editFullName} onChange={(_, v) => setEditFullName(v.value)} />
                  </Field>
                  <Field label={t('drivers.licenseNumber')} required>
                    <Input value={editLicense} onChange={(_, v) => setEditLicense(v.value)} />
                  </Field>
                  <Field label={t('drivers.experience')} required>
                    <Input
                      type="number"
                      min={0}
                      value={String(editExperience)}
                      onChange={(_, v) =>
                        setEditExperience(v.value === '' ? '' : Math.max(0, parseInt(v.value, 10) || 0))
                      }
                    />
                  </Field>
                  <Field label={t('drivers.phone')}>
                    <Input value={editPhone} onChange={(_, v) => setEditPhone(v.value)} />
                  </Field>
                </div>
              )}
            </DialogContent>
            <DialogActions>
              <div className={actionsStyles.wrapper}>
                <Button appearance="secondary" onClick={() => setEditDriver(null)}>
                  {t('common.cancel')}
                </Button>
                <Button
                  appearance="primary"
                  disabled={!editDriver || !editFullName.trim() || !editLicense.trim() || updateMutation.isPending}
                onClick={() =>
                  editDriver &&
                  updateMutation.mutate({
                    id: editDriver.id,
                    data: {
                      full_name: editFullName.trim(),
                      license_number: editLicense.trim(),
                      experience_years:
                        editExperience === '' ? undefined : (editExperience as number),
                      phone: editPhone.trim() || undefined,
                    },
                  })
                }
              >
                {t('common.save')}
              </Button>
              </div>
            </DialogActions>
          </DialogBody>
        </DialogSurface>
      </Dialog>

      <Dialog open={!!deleteDriver} onOpenChange={(_, d) => (!d.open && setDeleteDriver(null))}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>{t('drivers.deleteDriverTitle')}</DialogTitle>
            <DialogContent>
              {deleteDriver && (
                <Text>
                  {t('drivers.deleteConfirm', { name: deleteDriver.full_name, license: deleteDriver.license_number })}
                </Text>
              )}
            </DialogContent>
            <DialogActions>
              <div className={actionsStyles.wrapper}>
                <Button appearance="secondary" onClick={() => setDeleteDriver(null)}>
                  {t('common.cancel')}
                </Button>
                <Button
                  appearance="primary"
                  disabled={deleteMutation.isPending}
                onClick={() => deleteDriver && deleteMutation.mutate(deleteDriver.id)}
              >
                {t('common.delete')}
              </Button>
              </div>
            </DialogActions>
          </DialogBody>
        </DialogSurface>
      </Dialog>
    </AppLayout>
  );
};
