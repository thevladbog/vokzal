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
import type { Bus, Station } from '@/types';

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

const STATUS_OPTIONS: Bus['status'][] = ['active', 'maintenance', 'out_of_service'];

export const BusesPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const queryClient = useQueryClient();
  const [stationFilter, setStationFilter] = useState('');
  const [createOpen, setCreateOpen] = useState(false);
  const [editBus, setEditBus] = useState<Bus | null>(null);
  const [deleteBus, setDeleteBus] = useState<Bus | null>(null);
  const [createPlate, setCreatePlate] = useState('');
  const [createModel, setCreateModel] = useState('');
  const [createCapacity, setCreateCapacity] = useState(45);
  const [createStationId, setCreateStationId] = useState('');
  const [createStatus, setCreateStatus] = useState<Bus['status']>('active');
  const [editPlate, setEditPlate] = useState('');
  const [editModel, setEditModel] = useState('');
  const [editCapacity, setEditCapacity] = useState(0);
  const [editStatus, setEditStatus] = useState<Bus['status']>('active');

  const { data: stations = [] } = useQuery<Station[]>({
    queryKey: ['stations'],
    queryFn: () => scheduleService.getStations(),
  });

  const { data: buses = [], isLoading, error } = useQuery<Bus[]>({
    queryKey: ['buses', stationFilter],
    queryFn: () =>
      scheduleService.getBuses(stationFilter ? { station_id: stationFilter } : {}),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<Bus, 'id' | 'created_at' | 'updated_at'>) =>
      scheduleService.createBus(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buses'] });
      setCreateOpen(false);
      setCreatePlate('');
      setCreateModel('');
      setCreateCapacity(45);
      setCreateStationId('');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Bus> }) =>
      scheduleService.updateBus(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buses'] });
      setEditBus(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => scheduleService.deleteBus(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buses'] });
      setDeleteBus(null);
    },
  });

  const stationName = (id: string) => stations.find((s) => s.id === id)?.name ?? id.slice(0, 8);

  const openEdit = (bus: Bus) => {
    setEditBus(bus);
    setEditPlate(bus.plate_number);
    setEditModel(bus.model);
    setEditCapacity(bus.capacity);
    setEditStatus(bus.status);
  };

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label={t('buses.loading')} />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>{t('buses.loadError')}</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>{t('buses.title')}</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              {t('buses.addBus')}
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('buses.createBusTitle')}</DialogTitle>
              <DialogContent>
                <div className={styles.formRow}>
                  <Label>{t('buses.plateNumber')}</Label>
                  <Input
                    value={createPlate}
                    onChange={(_, v) => setCreatePlate(v.value)}
                    placeholder={t('buses.platePlaceholder')}
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('buses.model')}</Label>
                  <Input
                    value={createModel}
                    onChange={(_, v) => setCreateModel(v.value)}
                    placeholder={t('buses.modelPlaceholder')}
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('buses.capacity')}</Label>
                  <Input
                    type="number"
                    min={1}
                    value={String(createCapacity)}
                    onChange={(_, v) => setCreateCapacity(parseInt(v.value, 10) || 0)}
                  />
                </div>
                <div className={styles.formRow}>
                  <Label>{t('buses.station')}</Label>
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
                <div className={styles.formRow}>
                  <Label>{t('buses.status')}</Label>
                  <Select
                    value={createStatus}
                    onChange={(_, d) => setCreateStatus((d.value as Bus['status']) ?? 'active')}
                    style={{ width: '100%' }}
                  >
                    {STATUS_OPTIONS.map((s) => (
                      <Option key={s} value={s}>
                        {s === 'active' ? t('buses.statusActive') : s === 'maintenance' ? t('buses.statusMaintenance') : t('buses.statusOutOfService')}
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
                  disabled={!createPlate.trim() || !createModel.trim() || createCapacity < 1 || !createStationId}
                  onClick={() =>
                    createMutation.mutate({
                      plate_number: createPlate.trim(),
                      model: createModel.trim(),
                      capacity: createCapacity,
                      station_id: createStationId,
                      status: createStatus,
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
        <Label>{t('buses.station')}:</Label>
        <Select
          value={stationFilter}
          onChange={(_, d) => setStationFilter(d.value ?? '')}
          style={{ minWidth: '200px' }}
        >
          <Option value="">{t('buses.allStations')}</Option>
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
              <TableHeaderCell>{t('buses.plateNumber')}</TableHeaderCell>
              <TableHeaderCell>{t('buses.model')}</TableHeaderCell>
              <TableHeaderCell>{t('buses.capacity')}</TableHeaderCell>
              <TableHeaderCell>{t('buses.station')}</TableHeaderCell>
              <TableHeaderCell>{t('buses.status')}</TableHeaderCell>
              <TableHeaderCell></TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {buses.map((bus) => (
              <TableRow key={bus.id}>
                <TableCell>{bus.plate_number}</TableCell>
                <TableCell>{bus.model}</TableCell>
                <TableCell>{bus.capacity}</TableCell>
                <TableCell>{stationName(bus.station_id)}</TableCell>
                <TableCell>
                  {bus.status === 'active'
                    ? t('buses.statusActive')
                    : bus.status === 'maintenance'
                      ? t('buses.statusMaintenance')
                      : t('buses.statusOutOfService')}
                </TableCell>
                <TableCell>
                  <Button
                    appearance="subtle"
                    icon={<Edit24Regular />}
                    onClick={() => openEdit(bus)}
                    aria-label={t('common.edit')}
                  />
                  <Button
                    appearance="subtle"
                    icon={<Delete24Regular />}
                    onClick={() => setDeleteBus(bus)}
                    aria-label={t('common.delete')}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>

      <Dialog open={!!editBus} onOpenChange={(_, d) => (!d.open && setEditBus(null))}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>{t('buses.editBusTitle')}</DialogTitle>
            <DialogContent>
              {editBus && (
                <>
                  <div className={styles.formRow}>
                    <Label>{t('buses.plateNumber')}</Label>
                    <Input
                      value={editPlate}
                      onChange={(_, v) => setEditPlate(v.value)}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('buses.model')}</Label>
                    <Input value={editModel} onChange={(_, v) => setEditModel(v.value)} />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('buses.capacity')}</Label>
                    <Input
                      type="number"
                      min={1}
                      value={String(editCapacity)}
                      onChange={(_, v) => setEditCapacity(parseInt(v.value, 10) || 0)}
                    />
                  </div>
                  <div className={styles.formRow}>
                    <Label>{t('buses.status')}</Label>
                    <Select
                      value={editStatus}
                      onChange={(_, d) => setEditStatus((d.value as Bus['status']) ?? 'active')}
                      style={{ width: '100%' }}
                    >
                      {STATUS_OPTIONS.map((s) => (
                        <Option key={s} value={s}>
                          {s === 'active' ? t('buses.statusActive') : s === 'maintenance' ? t('buses.statusMaintenance') : t('buses.statusOutOfService')}
                        </Option>
                      ))}
                    </Select>
                  </div>
                </>
              )}
            </DialogContent>
            <DialogActions>
              <Button appearance="secondary" onClick={() => setEditBus(null)}>
                {t('common.cancel')}
              </Button>
              <Button
                appearance="primary"
                disabled={!editBus || !editPlate.trim() || !editModel.trim() || editCapacity < 1}
                onClick={() =>
                  editBus &&
                  updateMutation.mutate({
                    id: editBus.id,
                    data: {
                      plate_number: editPlate.trim(),
                      model: editModel.trim(),
                      capacity: editCapacity,
                      status: editStatus,
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

      <Dialog open={!!deleteBus} onOpenChange={(_, d) => (!d.open && setDeleteBus(null))}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>{t('buses.deleteBusTitle')}</DialogTitle>
            <DialogContent>
              {deleteBus && (
                <Text>
                  {t('buses.deleteConfirm', { plate: deleteBus.plate_number, model: deleteBus.model })}
                </Text>
              )}
            </DialogContent>
            <DialogActions>
              <Button appearance="secondary" onClick={() => setDeleteBus(null)}>
                {t('common.cancel')}
              </Button>
              <Button
                appearance="primary"
                onClick={() => deleteBus && deleteMutation.mutate(deleteBus.id)}
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
