import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  FluentProvider,
  webLightTheme,
  Input,
  Button,
  Text,
  Title2,
  Card,
  makeStyles,
  tokens,
} from '@fluentui/react-components';
import { authService } from '@/services/auth';
import { useAuthStore } from '@/stores/authStore';

const useStyles = makeStyles({
  container: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    backgroundColor: tokens.colorNeutralBackground2,
  },
  card: {
    width: '400px',
    padding: '32px',
  },
  title: {
    marginBottom: '24px',
    textAlign: 'center',
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
  error: {
    color: tokens.colorPaletteRedForeground1,
    fontSize: '12px',
  },
});

export const LoginPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const setUser = useAuthStore((state) => state.setUser);

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await authService.login(username, password);
      setUser(response.user);
      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка входа. Проверьте логин и пароль.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <FluentProvider theme={webLightTheme}>
      <div className={styles.container}>
        <Card className={styles.card}>
          <Title2 className={styles.title}>Вокзал.ТЕХ</Title2>
          <Text block style={{ textAlign: 'center', marginBottom: '24px' }}>
            Вход в админ-панель
          </Text>

          <form onSubmit={handleSubmit} className={styles.form}>
            <Input
              type="text"
              placeholder="Логин"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <Input
              type="password"
              placeholder="Пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />

            {error && <Text className={styles.error}>{error}</Text>}

            <Button 
              appearance="primary" 
              type="submit" 
              disabled={loading}
            >
              {loading ? 'Вход...' : 'Войти'}
            </Button>
          </form>
        </Card>
      </div>
    </FluentProvider>
  );
};
