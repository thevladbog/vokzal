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
  Select,
  Option,
} from '@fluentui/react-components';
import { Add24Regular, Delete24Regular, Edit24Regular } from '@fluentui/react-icons';
import { scheduleService } from '@/services/schedule';
import type { Driver, Station } from '@/types';

const useStyles = makeStyles({
  container: { padding: '24px' },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  filters: { display: 'flex', gap: '12px', alignItems: 'center', marginBottom: '16px' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
  formRow: { marginBottom: '16px' },
});

export const DriversPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
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
      <div className={styles.loading}>
        <Spinner label={t('drivers.loading')} />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>{t('drivers.loadError')}</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
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
                <div className={styles.formRow}>
                  <Label>{t('drivers.fullName')}</Label>
                  <Input
                    value={createFullName}
                    onChange={(_, v) => setCreateFullName(v.value)}
                    placeholder={t('drivers.fullNamePlaceholder')}
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('drivers.licenseNumber')}</Label>
                  <Input
                    value={createLicense}
                    onChange={(_, v) => setCreateLicense(v.value)}
                    placeholder={t('drivers.licensePlaceholder')}
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('drivers.experience')}</Label>
                  <Input
                    type="number"
                    min={0}
                    value={String(createExperience)}
                    onChange={(_, v) =>
                      setCreateExperience(v.value === '' ? '' : Math.max(0, parseInt(v.value, 10) || 0))
                    }
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('drivers.phone')}</Label>
                  <Input
                    value={createPhone}
                    onChange={(_, v) => setCreatePhone(v.value)}
                    placeholder={t('drivers.phonePlaceholder')}
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('drivers.station')}</Label>
                  <Select
                    value={createStationId}
                    onChange={(_, d) => setCreateStationId(d.value ?? '')}
                    style={{ width: '100%' }}
                  >
                    {stations.map((s) => (
                      <Option key={s.id} value={s.id}>
                        {s.name}
                      </Option>
                    ))}
                  </Select>
                </div>
              </DialogContent>
              <DialogActions>
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
              </DialogActions>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      </div>

      <div className={styles.filters}>
        <Label>{t('drivers.station')}:</Label>
        <Select
          value={stationFilter}
          onChange={(_, d) => setStationFilter(d.value ?? '')}
          style={{ minWidth: '200px' }}
        >
          <Option value="">{t('drivers.allStations')}</Option>
          {stations.map((s) => (
            <Option key={s.id} value={s.id}>
              {s.name}
            </Option>
          ))}
        </Select>
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
                <>
                  <div className={styles.formRow}>
                    <Label>{t('drivers.fullName')}</Label>
                    <Input value={editFullName} onChange={(_, v) => setEditFullName(v.value)} />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('drivers.licenseNumber')}</Label>
                    <Input value={editLicense} onChange={(_, v) => setEditLicense(v.value)} />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('drivers.experience')}</Label>
                    <Input
                      type="number"
                      min={0}
                      value={String(editExperience)}
                      onChange={(_, v) =>
                        setEditExperience(v.value === '' ? '' : Math.max(0, parseInt(v.value, 10) || 0))
                      }
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('drivers.phone')}</Label>
                    <Input value={editPhone} onChange={(_, v) => setEditPhone(v.value)} />
                  </div>
                </>
              )}
            </DialogContent>
            <DialogActions>
              <Button appearance="secondary" onClick={() => setEditDriver(null)}>
                {t('common.cancel')}
              </Button>
              <Button
                appearance="primary"
                disabled={!editDriver || !editFullName.trim() || !editLicense.trim()}
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
              <Button appearance="secondary" onClick={() => setDeleteDriver(null)}>
                {t('common.cancel')}
              </Button>
              <Button
                appearance="primary"
                onClick={() => deleteDriver && deleteMutation.mutate(deleteDriver.id)}
              >
                {t('common.delete')}
              </Button>
            </DialogActions>
          </DialogBody>
        </DialogSurface>
      </Dialog>
    </div>
  );
};
