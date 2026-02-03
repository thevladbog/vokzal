import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
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
} from '@fluentui/react-components';
import { auditService } from '@/services/audit';
import type { AuditLog } from '@/types';
import { formatDateTime } from '@/utils/format';

const useStyles = makeStyles({
  container: { padding: '24px' },
  header: { marginBottom: '24px' },
  filters: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '12px',
    alignItems: 'flex-end',
    marginBottom: '16px',
  },
  filterRow: { display: 'flex', alignItems: 'center', gap: '8px' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
});

export const AuditPage: React.FC = () => {
  const styles = useStyles();
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [entityType, setEntityType] = useState('');
  const [entityId, setEntityId] = useState('');
  const [userId, setUserId] = useState('');
  const [applyFilters, setApplyFilters] = useState(false);

  const { data: logs = [], isLoading, error } = useQuery<AuditLog[]>({
    queryKey: ['audit', applyFilters, from, to, entityType, entityId, userId],
    queryFn: async () => {
      if (applyFilters && from && to) {
        return auditService.getLogsByDateRange(from, to);
      }
      if (applyFilters && entityType && entityId) {
        return auditService.getLogsByEntity(entityType, entityId);
      }
      if (applyFilters && userId) {
        return auditService.getLogsByUser(userId, 100);
      }
      return auditService.listLogs(100);
    },
  });

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label="Загрузка журнала аудита..." />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>Ошибка загрузки журнала аудита</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>Журнал аудита (152-ФЗ)</Title2>
      </div>
      <div className={styles.filters}>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-from">С</Label>
          <Input
            id="audit-from"
            type="date"
            value={from}
            onChange={(_, v) => setFrom(v.value)}
            style={{ width: '140px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-to">По</Label>
          <Input
            id="audit-to"
            type="date"
            value={to}
            onChange={(_, v) => setTo(v.value)}
            style={{ width: '140px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-entity-type">Тип сущности</Label>
          <Input
            id="audit-entity-type"
            value={entityType}
            onChange={(_, v) => setEntityType(v.value)}
            placeholder="ticket"
            style={{ width: '120px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-entity-id">ID сущности</Label>
          <Input
            id="audit-entity-id"
            value={entityId}
            onChange={(_, v) => setEntityId(v.value)}
            placeholder="UUID"
            style={{ width: '200px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-user">ID пользователя</Label>
          <Input
            id="audit-user"
            value={userId}
            onChange={(_, v) => setUserId(v.value)}
            placeholder="UUID"
            style={{ width: '200px' }}
          />
        </div>
        <Button
          appearance={applyFilters ? 'primary' : 'secondary'}
          onClick={() => setApplyFilters(!applyFilters)}
        >
          {applyFilters ? 'Применить фильтры' : 'Показать все'}
        </Button>
      </div>
      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>Время</TableHeaderCell>
              <TableHeaderCell>Действие</TableHeaderCell>
              <TableHeaderCell>Сущность</TableHeaderCell>
              <TableHeaderCell>ID сущности</TableHeaderCell>
              <TableHeaderCell>Пользователь</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {logs.map((log) => (
              <TableRow key={log.id}>
                <TableCell>{formatDateTime(log.created_at)}</TableCell>
                <TableCell>{log.action}</TableCell>
                <TableCell>{log.entity_type}</TableCell>
                <TableCell>{log.entity_id.slice(0, 8)}…</TableCell>
                <TableCell>{log.user_id ? log.user_id.slice(0, 8) + '…' : '—'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {logs.length === 0 && (
          <div style={{ padding: '24px', textAlign: 'center' }}>
            <Text>Записей не найдено</Text>
          </div>
        )}
      </Card>
    </div>
  );
};
