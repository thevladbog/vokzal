import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
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
    marginBottom: '8px',
  },
  filterRow: { display: 'flex', alignItems: 'center', gap: '8px' },
  filterHint: { marginBottom: '8px', color: 'var(--colorNeutralForeground3)' },
  filterError: { marginBottom: '16px', color: 'var(--colorPaletteRedForeground1)' },
  loading: { display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' },
});

export const AuditPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');
  const [entityType, setEntityType] = useState('');
  const [entityId, setEntityId] = useState('');
  const [userId, setUserId] = useState('');
  const [applyFilters, setApplyFilters] = useState(false);
  const [filterError, setFilterError] = useState('');

  const handleApplyToggle = () => {
    if (!applyFilters) {
      const dateRangeFilled = Boolean(from.trim() && to.trim());
      const entityFilled = Boolean(entityType.trim() && entityId.trim());
      const userFilled = Boolean(userId.trim());
      const filledCount = [dateRangeFilled, entityFilled, userFilled].filter(Boolean).length;
      if (filledCount > 1) {
        setFilterError(t('auditPage.onlyOneFilter'));
        return;
      }
      setFilterError('');
    }
    setApplyFilters(!applyFilters);
  };

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
        <Spinner label={t('auditPage.loading')} />
      </div>
    );
  }
  if (error) {
    return (
      <div className={styles.container}>
        <Text>{t('auditPage.loadError')}</Text>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title2>{t('auditPage.title')}</Title2>
      </div>
      <Text className={styles.filterHint} as="p">
        {t('auditPage.filterHint')}
      </Text>
      <div className={styles.filters}>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-from">{t('auditPage.dateFrom')}</Label>
          <Input
            id="audit-from"
            type="date"
            value={from}
            onChange={(_, v) => {
              setFrom(v.value);
              setFilterError('');
            }}
            style={{ width: '140px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-to">{t('auditPage.dateTo')}</Label>
          <Input
            id="audit-to"
            type="date"
            value={to}
            onChange={(_, v) => {
              setTo(v.value);
              setFilterError('');
            }}
            style={{ width: '140px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-entity-type">{t('auditPage.entityType')}</Label>
          <Input
            id="audit-entity-type"
            value={entityType}
            onChange={(_, v) => {
              setEntityType(v.value);
              setFilterError('');
            }}
            placeholder={t('auditPage.entityTypePlaceholder')}
            style={{ width: '120px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-entity-id">{t('auditPage.entityId')}</Label>
          <Input
            id="audit-entity-id"
            value={entityId}
            onChange={(_, v) => {
              setEntityId(v.value);
              setFilterError('');
            }}
            placeholder={t('auditPage.entityIdPlaceholder')}
            style={{ width: '200px' }}
          />
        </div>
        <div className={styles.filterRow}>
          <Label htmlFor="audit-user">{t('auditPage.userId')}</Label>
          <Input
            id="audit-user"
            value={userId}
            onChange={(_, v) => {
              setUserId(v.value);
              setFilterError('');
            }}
            placeholder={t('auditPage.userIdPlaceholder')}
            style={{ width: '200px' }}
          />
        </div>
        <Button
          appearance={applyFilters ? 'primary' : 'secondary'}
          onClick={handleApplyToggle}
        >
          {applyFilters ? t('auditPage.applyFilters') : t('auditPage.showAll')}
        </Button>
      </div>
      {filterError && (
        <Text className={styles.filterError} as="p">
          {filterError}
        </Text>
      )}
      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('auditPage.time')}</TableHeaderCell>
              <TableHeaderCell>{t('auditPage.action')}</TableHeaderCell>
              <TableHeaderCell>{t('auditPage.entity')}</TableHeaderCell>
              <TableHeaderCell>{t('auditPage.entityIdCol')}</TableHeaderCell>
              <TableHeaderCell>{t('auditPage.user')}</TableHeaderCell>
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
            <Text>{t('auditPage.noRecords')}</Text>
          </div>
        )}
      </Card>
    </div>
  );
};
