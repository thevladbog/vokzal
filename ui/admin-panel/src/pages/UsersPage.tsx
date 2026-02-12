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
  useId,
  Toaster,
  useToastController,
  Toast,
  ToastTitle,
  ToastBody,
} from '@fluentui/react-components';
import { Add24Regular, Delete24Regular, Edit24Regular } from '@fluentui/react-icons';
import { usersService } from '@/services/users';
import type { UserAdmin, CreateUserRequest, UpdateUserRequest, UserRole } from '@/types';

const useStyles = makeStyles({
  container: {
    padding: '24px',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  filters: {
    display: 'flex',
    gap: '12px',
    alignItems: 'center',
    marginBottom: '16px',
  },
  loading: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '400px',
  },
  formRow: {
    marginBottom: '16px',
  },
  actions: {
    display: 'flex',
    gap: '8px',
  },
});

const ROLE_VALUES: UserRole[] = ['admin', 'dispatcher', 'cashier', 'controller', 'accountant'];

export const UsersPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const queryClient = useQueryClient();
  const listId = useId('list-users');
  const { dispatchToast } = useToastController(listId);
  const getRoleLabel = (role: UserRole) => t(`users.role_${role}`, { defaultValue: role });
  const [page, setPage] = useState(1);
  const [roleFilter, setRoleFilter] = useState<string>('');
  const [createOpen, setCreateOpen] = useState(false);
  const [editUser, setEditUser] = useState<UserAdmin | null>(null);
  const [deleteUser, setDeleteUser] = useState<UserAdmin | null>(null);
  const { data, isLoading, error } = useQuery({
    queryKey: ['users', page, roleFilter],
    queryFn: () =>
      usersService.list({
        page,
        limit: 20,
        ...(roleFilter ? { role: roleFilter } : {}),
      }),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateUserRequest) => usersService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setCreateOpen(false);
      dispatchToast(
        <Toast>
          <ToastTitle>{t('users.createSuccess')}</ToastTitle>
          <ToastBody>{t('users.createSuccessBody')}</ToastBody>
        </Toast>,
        { intent: 'success' }
      );
    },
    onError: (err: { response?: { data?: { error?: string } } }) => {
      dispatchToast(
        <Toast>
          <ToastTitle>{t('common.error')}</ToastTitle>
          <ToastBody>{err.response?.data?.error ?? t('users.createError')}</ToastBody>
        </Toast>,
        { intent: 'error' }
      );
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserRequest }) =>
      usersService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setEditUser(null);
      dispatchToast(
        <Toast>
          <ToastTitle>{t('users.updateSuccess')}</ToastTitle>
          <ToastBody>{t('users.updateSuccessBody')}</ToastBody>
        </Toast>,
        { intent: 'success' }
      );
    },
    onError: (err: { response?: { data?: { error?: string } } }) => {
      dispatchToast(
        <Toast>
          <ToastTitle>{t('common.error')}</ToastTitle>
          <ToastBody>{err.response?.data?.error ?? t('users.updateError')}</ToastBody>
        </Toast>,
        { intent: 'error' }
      );
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => usersService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setDeleteUser(null);
      dispatchToast(
        <Toast>
          <ToastTitle>{t('users.deleteSuccess')}</ToastTitle>
          <ToastBody>{t('users.deleteSuccessBody')}</ToastBody>
        </Toast>,
        { intent: 'success' }
      );
    },
    onError: (err: { response?: { data?: { error?: string } } }) => {
      dispatchToast(
        <Toast>
          <ToastTitle>{t('common.error')}</ToastTitle>
          <ToastBody>{err.response?.data?.error ?? t('users.deleteError')}</ToastBody>
        </Toast>,
        { intent: 'error' }
      );
    },
  });

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label={t('users.loading')} />
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <Text>{t('users.loadError')}</Text>
      </div>
    );
  }

  const users = data?.users ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.ceil(total / (data?.limit ?? 20));

  return (
    <div className={styles.container}>
      <Toaster toasterId={listId} />
      <div className={styles.header}>
        <Title2>{t('users.title')}</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              {t('users.addUser')}
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('users.createUserTitle')}</DialogTitle>
              <DialogContent>
                <CreateUserForm
                  getRoleLabel={getRoleLabel}
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
        <Label htmlFor="role-filter">{t('users.roleFilter')}</Label>
        <Select
          id="role-filter"
          value={roleFilter}
          onChange={(_, v) => setRoleFilter(v.value ?? '')}
          style={{ minWidth: '160px' }}
        >
          <option value="">{t('users.allRoles')}</option>
          {ROLE_VALUES.map((r) => (
            <option key={r} value={r}>
              {getRoleLabel(r)}
            </option>
          ))}
        </Select>
      </div>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>{t('users.username')}</TableHeaderCell>
              <TableHeaderCell>{t('users.fullName')}</TableHeaderCell>
              <TableHeaderCell>{t('users.role')}</TableHeaderCell>
              <TableHeaderCell>{t('users.status')}</TableHeaderCell>
              <TableHeaderCell>{t('users.actions')}</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.username}</TableCell>
                <TableCell>{user.full_name}</TableCell>
                <TableCell>{getRoleLabel(user.role)}</TableCell>
                <TableCell>{user.is_active ? t('users.active') : t('users.inactive')}</TableCell>
                <TableCell>
                  <div className={styles.actions}>
                    <Button
                      appearance="subtle"
                      icon={<Edit24Regular />}
                      onClick={() => setEditUser(user)}
                      aria-label={t('users.edit')}
                    />
                    <Button
                      appearance="subtle"
                      icon={<Delete24Regular />}
                      onClick={() => setDeleteUser(user)}
                      aria-label={t('common.delete')}
                    />
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {totalPages > 1 && (
          <div className={styles.filters} style={{ marginTop: '16px' }}>
            <Button
              appearance="subtle"
              disabled={page <= 1}
              onClick={() => setPage((p) => p - 1)}
            >
              {t('users.back')}
            </Button>
            <Text>{t('users.pageOf', { page, total: totalPages })}</Text>
            <Button
              appearance="subtle"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              {t('users.next')}
            </Button>
          </div>
        )}
      </Card>

      {editUser && (
        <Dialog open={!!editUser} onOpenChange={(_, d) => !d.open && setEditUser(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('users.editUserTitle')}</DialogTitle>
              <DialogContent>
                <EditUserForm
                  user={editUser}
                  getRoleLabel={getRoleLabel}
                  onSubmit={(formData) =>
                    updateMutation.mutate({ id: editUser.id, data: formData })
                  }
                  onCancel={() => setEditUser(null)}
                  isLoading={updateMutation.isPending}
                />
              </DialogContent>
            </DialogBody>
          </DialogSurface>
        </Dialog>
      )}

      {deleteUser && (
        <Dialog open={!!deleteUser} onOpenChange={(_, d) => !d.open && setDeleteUser(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>{t('users.deleteUserTitle')}</DialogTitle>
              <DialogContent>
                <Text>
                  {t('users.deleteConfirm', {
                    username: deleteUser.username,
                    fullName: deleteUser.full_name,
                  })}
                </Text>
              </DialogContent>
              <DialogActions>
                <DialogTrigger disableButtonEnhancement>
                  <Button appearance="secondary" onClick={() => setDeleteUser(null)}>
                    {t('common.cancel')}
                  </Button>
                </DialogTrigger>
                <Button
                  appearance="primary"
                  onClick={() => deleteMutation.mutate(deleteUser.id)}
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

// Create user form (inline in dialog)
const CreateUserForm: React.FC<{
  getRoleLabel: (role: UserRole) => string;
  onSubmit: (data: CreateUserRequest) => void;
  onCancel: () => void;
  isLoading: boolean;
}> = ({ getRoleLabel, onSubmit, onCancel, isLoading }) => {
  const { t } = useTranslation();
  const styles = useStyles();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [fullName, setFullName] = useState('');
  const [role, setRole] = useState<CreateUserRequest['role']>('cashier');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !password.trim() || !fullName.trim()) return;
    onSubmit({ username: username.trim(), password, full_name: fullName.trim(), role });
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className={styles.formRow}>
        <Label htmlFor="create-username">{t('users.username')} *</Label>
        <Input
          id="create-username"
          value={username}
          onChange={(_, v) => setUsername(v.value)}
          required
          minLength={3}
          maxLength={50}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="create-password">{t('users.password')} *</Label>
        <Input
          id="create-password"
          type="password"
          value={password}
          onChange={(_, v) => setPassword(v.value)}
          required
          minLength={8}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="create-fullname">{t('users.fullName')} *</Label>
        <Input
          id="create-fullname"
          value={fullName}
          onChange={(_, v) => setFullName(v.value)}
          required
          maxLength={100}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="create-role">{t('users.role')} *</Label>
        <Select
          id="create-role"
          value={role}
          onChange={(_, v) => setRole((v.value ?? 'cashier') as CreateUserRequest['role'])}
        >
          {ROLE_VALUES.map((r) => (
            <option key={r} value={r}>
              {getRoleLabel(r)}
            </option>
          ))}
        </Select>
      </div>
      <DialogActions>
        <Button type="button" appearance="secondary" onClick={onCancel}>
          {t('common.cancel')}
        </Button>
        <Button type="submit" appearance="primary" disabled={isLoading}>
          {t('common.create')}
        </Button>
      </DialogActions>
    </form>
  );
};

// Edit user form
const EditUserForm: React.FC<{
  user: UserAdmin;
  getRoleLabel: (role: UserRole) => string;
  onSubmit: (data: UpdateUserRequest) => void;
  onCancel: () => void;
  isLoading: boolean;
}> = ({ user, getRoleLabel, onSubmit, onCancel, isLoading }) => {
  const { t } = useTranslation();
  const styles = useStyles();
  const [fullName, setFullName] = useState(user.full_name);
  const [password, setPassword] = useState('');
  const [role, setRole] = useState<UpdateUserRequest['role']>(user.role);
  const [isActive, setIsActive] = useState(user.is_active);
  const [passwordError, setPasswordError] = useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordError(null);
    const trimmed = password.trim();
    if (trimmed.length > 0 && trimmed.length < 8) {
      setPasswordError(t('users.passwordMinError'));
      return;
    }
    const data: UpdateUserRequest = {
      full_name: fullName,
      role,
      is_active: isActive,
    };
    if (trimmed.length >= 8) data.password = trimmed;
    onSubmit(data);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className={styles.formRow}>
        <Label>{t('users.username')}</Label>
        <Text block>{user.username}</Text>
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="edit-fullname">{t('users.fullName')} *</Label>
        <Input
          id="edit-fullname"
          value={fullName}
          onChange={(_, v) => setFullName(v.value)}
          required
          maxLength={100}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="edit-password">{t('users.newPasswordPlaceholder')}</Label>
        <Input
          id="edit-password"
          type="password"
          value={password}
          onChange={(_, v) => {
            setPassword(v.value);
            setPasswordError(null);
          }}
          minLength={8}
          placeholder={t('users.passwordMinHint')}
          aria-invalid={!!passwordError}
          aria-describedby={passwordError ? 'edit-password-error' : undefined}
        />
        {passwordError && (
          <Text id="edit-password-error" style={{ color: 'var(--colorPaletteRedForeground1)', marginTop: '4px' }} role="alert">
            {passwordError}
          </Text>
        )}
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="edit-role">{t('users.role')} *</Label>
        <Select
          id="edit-role"
          value={role}
          onChange={(_, v) => setRole((v.value ?? user.role) as UpdateUserRequest['role'])}
        >
          {ROLE_VALUES.map((r) => (
            <option key={r} value={r}>
              {getRoleLabel(r)}
            </option>
          ))}
        </Select>
      </div>
      <div className={styles.formRow}>
        <Label>
          <input
            type="checkbox"
            checked={isActive}
            onChange={(e) => setIsActive(e.target.checked)}
          />{' '}
          {t('users.activeCheckbox')}
        </Label>
      </div>
      <DialogActions>
        <Button type="button" appearance="secondary" onClick={onCancel}>
          {t('common.cancel')}
        </Button>
        <Button type="submit" appearance="primary" disabled={isLoading}>
          {t('common.save')}
        </Button>
      </DialogActions>
    </form>
  );
};
