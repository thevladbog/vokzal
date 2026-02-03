import React, { useState } from 'react';
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
  useId,
  Toaster,
  useToastController,
  Toast,
  ToastTitle,
  ToastBody,
} from '@fluentui/react-components';
import { Add24Regular, Delete24Regular, Edit24Regular } from '@fluentui/react-icons';
import { usersService } from '@/services/users';
import type { UserAdmin, CreateUserRequest, UpdateUserRequest } from '@/types';

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

const ROLES = [
  { value: 'admin', label: 'Администратор' },
  { value: 'dispatcher', label: 'Диспетчер' },
  { value: 'cashier', label: 'Кассир' },
  { value: 'controller', label: 'Контролёр' },
] as const;

export const UsersPage: React.FC = () => {
  const styles = useStyles();
  const queryClient = useQueryClient();
  const { dispatchToast } = useToastController();
  const [page, setPage] = useState(1);
  const [roleFilter, setRoleFilter] = useState<string>('');
  const [createOpen, setCreateOpen] = useState(false);
  const [editUser, setEditUser] = useState<UserAdmin | null>(null);
  const [deleteUser, setDeleteUser] = useState<UserAdmin | null>(null);

  const listId = useId('list-users');
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
          <ToastTitle>Пользователь создан</ToastTitle>
          <ToastBody>Пользователь успешно добавлен.</ToastBody>
        </Toast>,
        { intent: 'success' }
      );
    },
    onError: (err: { response?: { data?: { error?: string } } }) => {
      dispatchToast(
        <Toast>
          <ToastTitle>Ошибка</ToastTitle>
          <ToastBody>{err.response?.data?.error ?? 'Не удалось создать пользователя'}</ToastBody>
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
          <ToastTitle>Пользователь обновлён</ToastTitle>
          <ToastBody>Изменения сохранены.</ToastBody>
        </Toast>,
        { intent: 'success' }
      );
    },
    onError: (err: { response?: { data?: { error?: string } } }) => {
      dispatchToast(
        <Toast>
          <ToastTitle>Ошибка</ToastTitle>
          <ToastBody>{err.response?.data?.error ?? 'Не удалось обновить пользователя'}</ToastBody>
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
          <ToastTitle>Пользователь удалён</ToastTitle>
          <ToastBody>Пользователь удалён из системы.</ToastBody>
        </Toast>,
        { intent: 'success' }
      );
    },
    onError: (err: { response?: { data?: { error?: string } } }) => {
      dispatchToast(
        <Toast>
          <ToastTitle>Ошибка</ToastTitle>
          <ToastBody>{err.response?.data?.error ?? 'Не удалось удалить пользователя'}</ToastBody>
        </Toast>,
        { intent: 'error' }
      );
    },
  });

  if (isLoading) {
    return (
      <div className={styles.loading}>
        <Spinner label="Загрузка пользователей..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <Text>Ошибка загрузки пользователей</Text>
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
        <Title2>Пользователи</Title2>
        <Dialog open={createOpen} onOpenChange={(_, d) => setCreateOpen(d.open)}>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="primary" icon={<Add24Regular />}>
              Добавить пользователя
            </Button>
          </DialogTrigger>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>Новый пользователь</DialogTitle>
              <DialogContent>
                <CreateUserForm
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
        <Label htmlFor="role-filter">Роль:</Label>
        <Select
          id="role-filter"
          value={roleFilter}
          onChange={(_, v) => setRoleFilter(v.value ?? '')}
          style={{ minWidth: '160px' }}
        >
          <Option value="">Все</Option>
          {ROLES.map((r) => (
            <Option key={r.value} value={r.value}>
              {r.label}
            </Option>
          ))}
        </Select>
      </div>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>Логин</TableHeaderCell>
              <TableHeaderCell>ФИО</TableHeaderCell>
              <TableHeaderCell>Роль</TableHeaderCell>
              <TableHeaderCell>Статус</TableHeaderCell>
              <TableHeaderCell>Действия</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.username}</TableCell>
                <TableCell>{user.full_name}</TableCell>
                <TableCell>{ROLES.find((r) => r.value === user.role)?.label ?? user.role}</TableCell>
                <TableCell>{user.is_active ? 'Активен' : 'Неактивен'}</TableCell>
                <TableCell>
                  <div className={styles.actions}>
                    <Button
                      appearance="subtle"
                      icon={<Edit24Regular />}
                      onClick={() => setEditUser(user)}
                      aria-label="Редактировать"
                    />
                    <Button
                      appearance="subtle"
                      icon={<Delete24Regular />}
                      onClick={() => setDeleteUser(user)}
                      aria-label="Удалить"
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
              Назад
            </Button>
            <Text>
              Страница {page} из {totalPages}
            </Text>
            <Button
              appearance="subtle"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              Вперёд
            </Button>
          </div>
        )}
      </Card>

      {editUser && (
        <Dialog open={!!editUser} onOpenChange={(_, d) => !d.open && setEditUser(null)}>
          <DialogSurface>
            <DialogBody>
              <DialogTitle>Редактировать пользователя</DialogTitle>
              <DialogContent>
                <EditUserForm
                  user={editUser}
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
              <DialogTitle>Удалить пользователя?</DialogTitle>
              <DialogContent>
                <Text>
                  Вы уверены, что хотите удалить пользователя {deleteUser.username} (
                  {deleteUser.full_name})? Это действие нельзя отменить.
                </Text>
              </DialogContent>
              <DialogActions>
                <DialogTrigger disableButtonEnhancement>
                  <Button appearance="secondary" onClick={() => setDeleteUser(null)}>
                    Отмена
                  </Button>
                </DialogTrigger>
                <Button
                  appearance="primary"
                  onClick={() => deleteMutation.mutate(deleteUser.id)}
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

// Create user form (inline in dialog)
const CreateUserForm: React.FC<{
  onSubmit: (data: CreateUserRequest) => void;
  onCancel: () => void;
  isLoading: boolean;
}> = ({ onSubmit, onCancel, isLoading }) => {
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
        <Label htmlFor="create-username">Логин *</Label>
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
        <Label htmlFor="create-password">Пароль *</Label>
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
        <Label htmlFor="create-fullname">ФИО *</Label>
        <Input
          id="create-fullname"
          value={fullName}
          onChange={(_, v) => setFullName(v.value)}
          required
          maxLength={100}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="create-role">Роль *</Label>
        <Select
          id="create-role"
          value={role}
          onChange={(_, v) => setRole((v.value ?? 'cashier') as CreateUserRequest['role'])}
        >
          {ROLES.map((r) => (
            <Option key={r.value} value={r.value}>
              {r.label}
            </Option>
          ))}
        </Select>
      </div>
      <DialogActions>
        <Button type="button" appearance="secondary" onClick={onCancel}>
          Отмена
        </Button>
        <Button type="submit" appearance="primary" disabled={isLoading}>
          Создать
        </Button>
      </DialogActions>
    </form>
  );
};

// Edit user form
const EditUserForm: React.FC<{
  user: UserAdmin;
  onSubmit: (data: UpdateUserRequest) => void;
  onCancel: () => void;
  isLoading: boolean;
}> = ({ user, onSubmit, onCancel, isLoading }) => {
  const styles = useStyles();
  const [fullName, setFullName] = useState(user.full_name);
  const [password, setPassword] = useState('');
  const [role, setRole] = useState<UpdateUserRequest['role']>(user.role);
  const [isActive, setIsActive] = useState(user.is_active);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const data: UpdateUserRequest = {
      full_name: fullName,
      role,
      is_active: isActive,
    };
    if (password.trim().length >= 8) data.password = password.trim();
    onSubmit(data);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className={styles.formRow}>
        <Label>Логин</Label>
        <Text block>{user.username}</Text>
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="edit-fullname">ФИО *</Label>
        <Input
          id="edit-fullname"
          value={fullName}
          onChange={(_, v) => setFullName(v.value)}
          required
          maxLength={100}
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="edit-password">Новый пароль (оставьте пустым, чтобы не менять)</Label>
        <Input
          id="edit-password"
          type="password"
          value={password}
          onChange={(_, v) => setPassword(v.value)}
          minLength={8}
          placeholder="Минимум 8 символов"
        />
      </div>
      <div className={styles.formRow}>
        <Label htmlFor="edit-role">Роль *</Label>
        <Select
          id="edit-role"
          value={role}
          onChange={(_, v) => setRole((v.value ?? user.role) as UpdateUserRequest['role'])}
        >
          {ROLES.map((r) => (
            <Option key={r.value} value={r.value}>
              {r.label}
            </Option>
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
          Активен
        </Label>
      </div>
      <DialogActions>
        <Button type="button" appearance="secondary" onClick={onCancel}>
          Отмена
        </Button>
        <Button type="submit" appearance="primary" disabled={isLoading}>
          Сохранить
        </Button>
      </DialogActions>
    </form>
  );
};
