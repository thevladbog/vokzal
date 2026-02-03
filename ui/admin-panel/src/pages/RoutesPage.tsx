import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Select, Option } from '@fluentui/react-components';
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
import type { Route, Station } from '@/types';

const useStyles = makeStyles({
  container: { padding: '24px' },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
  formRow: { marginBottom: '16px' },
  actions: { display: 'flex', gap: '8px' },
});

export const RoutesPage: React.FC = () => {
  const styles = useStyles();
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
      <div className={styles.loading}>
        <Spinner label="Загрузка маршрутов..." />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>Ошибка загрузки маршрутов</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>Маршруты</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              Добавить маршрут
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>Новый маршрут</DialogTitle>
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
              <TableHeaderCell>Название</TableHeaderCell>
              <TableHeaderCell>Расстояние (км)</TableHeaderCell>
              <TableHeaderCell>Время (мин)</TableHeaderCell>
              <TableHeaderCell>Остановок</TableHeaderCell>
              <TableHeaderCell>Статус</TableHeaderCell>
              <TableHeaderCell>Действия</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {routes.map((r) => (
              <TableRow key={r.id}>
                <TableCell>{r.name}</TableCell>
                <TableCell>{r.distance_km ?? '—'}</TableCell>
                <TableCell>{r.duration_min ?? '—'}</TableCell>
                <TableCell>{Array.isArray(r.stops) ? r.stops.length : 0}</TableCell>
                <TableCell>{r.is_active ? 'Активен' : 'Неактивен'}</TableCell>
                <TableCell>
                  <div className={styles.actions}>
                    <Button
                      appearance="subtle"
                      icon={<Edit24Regular />}
                      onClick={() => setEditRoute(r)}
                      aria-label="Редактировать"
                    />
                    <Button
                      appearance="subtle"
                      icon={<Delete24Regular />}
                      onClick={() => setDeleteRoute(r)}
                      aria-label="Удалить"
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
              <DialogTitle>Редактировать маршрут</DialogTitle>
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
              <DialogTitle>Удалить маршрут?</DialogTitle>
              <DialogContent>
                <Text>
                  Удалить маршрут «{deleteRoute.name}»? Это действие нельзя отменить.
                </Text>
              </DialogContent>
              <DialogActions>
                <Button appearance="secondary" onClick={() => setDeleteRoute(null)}>
                  Отмена
                </Button>
                <Button
                  appearance="primary"
                  onClick={() => deleteMutation.mutate(deleteRoute.id)}
                  disabled={deleteMutation.isPending}
                >
                  Удалить
                </Button>
              </DialogActions>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      )}
    </div>
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
  const styles = useStyles();
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
      <div className={styles.formRow}>
        <Label htmlFor="route-name">Название *</Label>
        <Input
          id="route-name"
          value={name}
          onChange={(_, v) => setName(v.value)}
          required
          maxLength={100}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="route-distance">Расстояние (км)</Label>
        <Input
          id="route-distance"
          type="number"
          min={0}
          step={0.1}
          value={distanceKm}
          onChange={(_, v) => setDistanceKm(v.value)}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="route-duration">Время в пути (мин)</Label>
        <Input
          id="route-duration"
          type="number"
          min={0}
          value={durationMin}
          onChange={(_, v) => setDurationMin(v.value)}
        />
      </div>
      {!isEdit && stations.length > 0 && (
        <div className={styles.formRow}>
          <Label htmlFor="route-first-station">Первая остановка (обязательно)</Label>
          <Select
            id="route-first-station"
            value={firstStationId}
            onChange={(_, v) => setFirstStationId(v.value ?? '')}
            style={{ minWidth: '100%' }}
          >
            {stations.map((s: Station) => (
              <Option key={s.id} value={s.id} text={`${s.name} (${s.code})`}>
                {s.name} ({s.code})
              </Option>
            ))}
          </Select>
        </div>
      )}
      {isEdit && (
        <div className={styles.formRow}>
          <Label>
            <input
              type="checkbox"
              checked={isActive}
              onChange={(e) => setIsActive(e.target.checked)}
            />{' '}
            Активен
          </Label>
        </div>
      )}
      <DialogActions>
        <Button type="button" appearance="secondary" onClick={onCancel}>
          Отмена
        </Button>
        <Button type="submit" appearance="primary" disabled={isLoading}>
          {isEdit ? 'Сохранить' : 'Создать'}
        </Button>
      </DialogActions>
    </form>
  );
};
