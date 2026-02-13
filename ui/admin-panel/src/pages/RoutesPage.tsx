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
  Field,
  Checkbox,
} from '@fluentui/react-components';
import { Dropdown, Option } from '@fluentui/react-combobox';
import { Add24Regular, Delete24Regular, Edit24Regular } from '@fluentui/react-icons';
import { AppLayout } from '@/components/layout/AppLayout';
import { useDialogFormStyles } from '@/styles/dialogFormStyles';
import { useDialogActionsStyles } from '@/styles/dialogActionsStyles';
import { scheduleService } from '@/services/schedule';
import type { Route, Station } from '@/types';

const useStyles = makeStyles({
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
  actions: { display: 'flex', gap: '8px' },
});

export const RoutesPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const actionsStyles = useDialogActionsStyles();
  const queryClient = useQueryClient();
  const [createOpen, setCreateOpen] = useState(false);
  const [editRoute, setEditRoute] = useState<Route | null>(null);
  const [deleteRoute, setDeleteRoute] = useState<Route | null>(null);

  const { data: routes = [], isLoading, error } = useQuery({
    queryKey: ['routes'],
    queryFn: () => scheduleService.getRoutes(),
  });

  const createMutation = useMutation({
    mutationFn: (data: { name: string; stops: Array<{ station_id: string; order: number; arrival_offset_min?: number }>; distance_km?: number; duration_min?: number }) =>
      scheduleService.createRoute({
        name: data.name,
        stops: data.stops,
        distance_km: data.distance_km ?? 0,
        duration_min: data.duration_min ?? 0,
        is_active: true,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routes'] });
      setCreateOpen(false);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Route> }) =>
      scheduleService.updateRoute(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routes'] });
      setEditRoute(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => scheduleService.deleteRoute(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['routes'] });
      setDeleteRoute(null);
    },
  });

  if (isLoading) {
    return (
      <AppLayout>
        <div className={styles.loading}>
          <Spinner label={t('routes.loading')} />
        </div>
      </AppLayout>
    );
  }
  if (error) {
    return (
      <AppLayout>
        <Text>{t('routes.loadError')}</Text>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className={styles.header}>
        <Title2>{t('routes.title')}</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              {t('routes.addRoute')}
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('routes.createRouteTitle')}</DialogTitle>
              <DialogContent>
                <RouteForm
                  onSubmit={(formData) => {
                    if (formData.stops && formData.stops.length > 0) {
                      createMutation.mutate({
                        name: formData.name,
                        stops: formData.stops,
                        distance_km: formData.distance_km,
                        duration_min: formData.duration_min,
                      });
                    }
                  }}
                  onCancel={() => setCreateOpen(false)}
                  isLoading={createMutation.isPending}
                />
              </DialogContent>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      </div>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('routes.name')}</TableHeaderCell>
              <TableHeaderCell>{t('routes.distanceKm')}</TableHeaderCell>
              <TableHeaderCell>{t('routes.durationMin')}</TableHeaderCell>
              <TableHeaderCell>{t('routes.stopsCount')}</TableHeaderCell>
              <TableHeaderCell>{t('routes.status')}</TableHeaderCell>
              <TableHeaderCell>{t('routes.actions')}</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {routes.map((r) => (
              <TableRow key={r.id}>
                <TableCell>{r.name}</TableCell>
                <TableCell>{r.distance_km ?? '—'}</TableCell>
                <TableCell>{r.duration_min ?? '—'}</TableCell>
                <TableCell>{Array.isArray(r.stops) ? r.stops.length : 0}</TableCell>
                <TableCell>{r.is_active ? t('routes.active') : t('routes.inactive')}</TableCell>
                <TableCell>
                  <div className={styles.actions}>
                    <Button
                      appearance="subtle"
                      icon={<Edit24Regular />}
                      onClick={() => setEditRoute(r)}
                      aria-label={t('common.edit')}
                    />
                    <Button
                      appearance="subtle"
                      icon={<Delete24Regular />}
                      onClick={() => setDeleteRoute(r)}
                      aria-label={t('common.delete')}
                    />
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Card>

      {editRoute && (
        <Dialog open={!!editRoute} onOpenChange={(_, d) => !d.open && setEditRoute(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('routes.editRouteTitle')}</DialogTitle>
              <DialogContent>
                <RouteForm
                  initial={editRoute}
                  onSubmit={(formData) =>
                    updateMutation.mutate({ id: editRoute.id, data: formData })
                  }
                  onCancel={() => setEditRoute(null)}
                  isLoading={updateMutation.isPending}
                  isEdit
                />
              </DialogContent>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      )}

      {deleteRoute && (
        <Dialog open={!!deleteRoute} onOpenChange={(_, d) => !d.open && setDeleteRoute(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('routes.deleteRouteTitle')}</DialogTitle>
              <DialogContent>
                <Text>
                  {t('routes.deleteConfirm', { name: deleteRoute.name })}
                </Text>
              </DialogContent>
              <DialogActions>
                <div className={actionsStyles.wrapper}>
                  <Button appearance="secondary" onClick={() => setDeleteRoute(null)}>
                    {t('common.cancel')}
                </Button>
                <Button
                  appearance="primary"
                  onClick={() => deleteMutation.mutate(deleteRoute.id)}
                  disabled={deleteMutation.isPending}
                >
                  {t('common.delete')}
                </Button>
                </div>
              </DialogActions>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      )}
    </AppLayout>
  );
};

type RouteFormData = {
  name: string;
  distance_km?: number;
  duration_min?: number;
  stops?: Array<{ station_id: string; order: number; arrival_offset_min?: number }>;
  is_active?: boolean;
};

const RouteForm: React.FC<{
  initial?: Route;
  onSubmit: (data: RouteFormData) => void;
  onCancel: () => void;
  isLoading: boolean;
  isEdit?: boolean;
}> = ({ initial, onSubmit, onCancel, isLoading, isEdit }) => {
  const { t } = useTranslation();
  const formStyles = useDialogFormStyles();
  const actionsStyles = useDialogActionsStyles();
  const { data: stationsRaw } = useQuery({
    queryKey: ['stations'],
    queryFn: () => scheduleService.getStations(),
  });
  const stations: Station[] = Array.isArray(stationsRaw) ? stationsRaw : [];
  const [name, setName] = useState(initial?.name ?? '');
  const [distanceKm, setDistanceKm] = useState(
    initial?.distance_km != null ? String(initial.distance_km) : ''
  );
  const [durationMin, setDurationMin] = useState(
    initial?.duration_min != null ? String(initial.duration_min) : ''
  );
  const [isActive, setIsActive] = useState(initial?.is_active ?? true);
  const [firstStationId, setFirstStationId] = useState(
    (Array.isArray(initial?.stops) && initial.stops[0]?.station_id) || ''
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    if (isEdit) {
      onSubmit({
        name: name.trim(),
        distance_km: distanceKm ? Number(distanceKm) : undefined,
        duration_min: durationMin ? Number(durationMin) : undefined,
        is_active: isActive,
      });
    } else {
      const sid = firstStationId || (stations[0]?.id ?? '');
      if (!sid) {
        return;
      }
      onSubmit({
        name: name.trim(),
        stops: [{ station_id: sid, order: 1, arrival_offset_min: 0 }],
        distance_km: distanceKm ? Number(distanceKm) : 0,
        duration_min: durationMin ? Number(durationMin) : 0,
      });
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className={formStyles.formContainer}>
        <Field label={t('routes.nameRequired')} required>
          <Input
            id="route-name"
            value={name}
            onChange={(_, v) => setName(v.value)}
            required
            maxLength={100}
          />
        </Field>
        <Field label={t('routes.distanceKm')}>
          <Input
            id="route-distance"
            type="number"
            min={0}
            step={0.1}
            value={distanceKm}
            onChange={(_, v) => setDistanceKm(v.value)}
          />
        </Field>
        <Field label={t('routes.durationLabel')}>
          <Input
            id="route-duration"
            type="number"
            min={0}
            value={durationMin}
            onChange={(_, v) => setDurationMin(v.value)}
          />
        </Field>
        {!isEdit && stations.length > 0 && (
          <Field label={t('routes.firstStationRequired')} required>
            <Dropdown
              placeholder={t('routes.selectStation')}
              value={stations.find(s => s.id === firstStationId)?.name || ''}
              selectedOptions={[firstStationId]}
              onOptionSelect={(_, data) => setFirstStationId(data.optionValue ?? '')}
            >
              {stations.map((s: Station) => (
                <Option key={s.id} value={s.id} text={s.name}>
                  {s.name} ({s.code})
                </Option>
              ))}
            </Dropdown>
          </Field>
        )}
        {isEdit && (
          <Field>
            <Checkbox
              label={t('routes.activeCheckbox')}
              checked={isActive}
              onChange={(_, data) => setIsActive(!!data.checked)}
            />
          </Field>
        )}
      </div>
      <DialogActions>
        <div className={actionsStyles.wrapper}>
          <Button type="button" appearance="secondary" onClick={onCancel}>
            {t('common.cancel')}
          </Button>
          <Button type="submit" appearance="primary" disabled={isLoading}>
            {isEdit ? t('common.save') : t('common.create')}
          </Button>
        </div>
      </DialogActions>
    </form>
  );
};
