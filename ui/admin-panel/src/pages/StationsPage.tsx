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
} from '@fluentui/react-components';
import { Add24Regular, Delete24Regular, Edit24Regular } from '@fluentui/react-icons';
import { scheduleService } from '@/services/schedule';
import type { Station } from '@/types';

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
  actions: { display: 'flex', gap: '8px' },
});

export const StationsPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const queryClient = useQueryClient();
  const [cityFilter, setCityFilter] = useState('');
  const [createOpen, setCreateOpen] = useState(false);
  const [editStation, setEditStation] = useState<Station | null>(null);
  const [deleteStation, setDeleteStation] = useState<Station | null>(null);

  const { data: stationsRaw, isLoading, error } = useQuery({
    queryKey: ['stations', cityFilter],
    queryFn: () => scheduleService.getStations(cityFilter ? { city: cityFilter } : undefined),
  });
  const stations = Array.isArray(stationsRaw) ? stationsRaw : [];

  const createMutation = useMutation({
    mutationFn: (data: Omit<Station, 'id' | 'created_at' | 'updated_at'>) =>
      scheduleService.createStation(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stations'] });
      setCreateOpen(false);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Station> }) =>
      scheduleService.updateStation(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stations'] });
      setEditStation(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => scheduleService.deleteStation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stations'] });
      setDeleteStation(null);
    },
  });

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label={t('stations.loading')} />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>{t('stations.loadError')}</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>{t('stations.title')}</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              {t('stations.addStation')}
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('stations.createStationTitle')}</DialogTitle>
              <DialogContent>
                <StationForm
                  onSubmit={(formData) => createMutation.mutate(formData)}
                  onCancel={() => setCreateOpen(false)}
                  isLoading={createMutation.isPending}
                />
              </DialogContent>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      </div>

      <div className={styles.filters}>
        <Label htmlFor="city-filter">{t('stations.cityFilterLabel')}</Label>
        <Input
          id="city-filter"
          value={cityFilter}
          onChange={(_, v) => setCityFilter(v.value)}
          placeholder={t('stations.filterPlaceholder')}
          style={{ minWidth: '200px' }}
        />
      </div>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('stations.code')}</TableHeaderCell>
              <TableHeaderCell>{t('stations.name')}</TableHeaderCell>
              <TableHeaderCell>{t('stations.address')}</TableHeaderCell>
              <TableHeaderCell>{t('stations.timezone')}</TableHeaderCell>
              <TableHeaderCell>{t('stations.actions')}</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {stations.map((s: Station) => (
              <TableRow key={s.id}>
                <TableCell>{s.code}</TableCell>
                <TableCell>{s.name}</TableCell>
                <TableCell>{s.address ?? '—'}</TableCell>
                <TableCell>{s.timezone ?? t('stations.timezonePlaceholder')}</TableCell>
                <TableCell>
                  <div className={styles.actions}>
                    <Button
                      appearance="subtle"
                      icon={<Edit24Regular />}
                      onClick={() => setEditStation(s)}
                      aria-label={t('common.edit')}
                    />
                    <Button
                      appearance="subtle"
                      icon={<Delete24Regular />}
                      onClick={() => setDeleteStation(s)}
                      aria-label={t('common.delete')}
                    />
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>

      {editStation && (
        <Dialog open={!!editStation} onOpenChange={(_, d) => !d.open && setEditStation(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('stations.editStationTitle')}</DialogTitle>
              <DialogContent>
                <StationForm
                  initial={editStation}
                  onSubmit={(formData) =>
                    updateMutation.mutate({ id: editStation.id, data: formData })
                  }
                  onCancel={() => setEditStation(null)}
                  isLoading={updateMutation.isPending}
                  isEdit
                />
              </DialogContent>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      )}

      {deleteStation && (
        <Dialog open={!!deleteStation} onOpenChange={(_, d) => !d.open && setDeleteStation(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('stations.deleteStationTitle')}</DialogTitle>
              <DialogContent>
                <Text>
                  {t('stations.deleteConfirm', {
                    name: deleteStation.name,
                    code: deleteStation.code,
                  })}
                </Text>
              </DialogContent>
              <DialogActions>
                <Button appearance="secondary" onClick={() => setDeleteStation(null)}>
                  {t('common.cancel')}
                </Button>
                <Button
                  appearance="primary"
                  onClick={() => deleteMutation.mutate(deleteStation.id)}
                  disabled={deleteMutation.isPending}
                >
                  {t('common.delete')}
                </Button>
              </DialogActions>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      )}
    </div>
  );
};

const StationForm: React.FC<{
  initial?: Station;
  onSubmit: (data: Partial<Station> & { name: string; code: string }) => void;
  onCancel: () => void;
  isLoading: boolean;
  isEdit?: boolean;
}> = ({ initial, onSubmit, onCancel, isLoading, isEdit }) => {
  const { t } = useTranslation();
  const styles = useStyles();
  const [name, setName] = useState(initial?.name ?? '');
  const [code, setCode] = useState(initial?.code ?? '');
  const [address, setAddress] = useState(initial?.address ?? '');
  const [timezone, setTimezone] = useState(initial?.timezone ?? 'Europe/Moscow');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !code.trim()) return;
    onSubmit({
      name: name.trim(),
      code: code.trim(),
      address: address.trim() || undefined,
      timezone: timezone.trim() || 'Europe/Moscow',
    });
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className={styles.formRow}>
        <Label htmlFor="station-name">{t('stations.nameRequired')}</Label>
        <Input
          id="station-name"
          value={name}
          onChange={(_, v) => setName(v.value)}
          required
          maxLength={100}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="station-code">{t('stations.codeRequired')}</Label>
        <Input
          id="station-code"
          value={code}
          onChange={(_, v) => setCode(v.value)}
          required
          maxLength={10}
          disabled={!!isEdit}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="station-address">{t('stations.address')}</Label>
        <Input
          id="station-address"
          value={address}
          onChange={(_, v) => setAddress(v.value)}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="station-timezone">{t('stations.timezone')}</Label>
        <Input
          id="station-timezone"
          value={timezone}
          onChange={(_, v) => setTimezone(v.value)}
          placeholder={t('stations.timezonePlaceholder')}
        />
      </div>
      <DialogActions>
        <Button type="button" appearance="secondary" onClick={onCancel}>
          {t('common.cancel')}
        </Button>
        <Button type="submit" appearance="primary" disabled={isLoading}>
          {isEdit ? t('common.save') : t('common.create')}
        </Button>
      </DialogActions>
    </form>
  );
};
