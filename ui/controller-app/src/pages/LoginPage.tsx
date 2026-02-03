import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Button,
  Input,
  Text,
  Card,
  makeStyles,
  tokens,
  Spinner,
} from '@fluentui/react-components';
import { ArrowRight24Regular } from '@fluentui/react-icons';
import { useAuthStore } from '@/stores/authStore';

const useStyles = makeStyles({
  container: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    padding: tokens.spacingVerticalXXL,
    backgroundColor: tokens.colorNeutralBackground3,
  },
  card: {
    width: '100%',
    maxWidth: '400px',
    padding: tokens.spacingVerticalXXL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalL,
    textAlign: 'center',
    color: tokens.colorBrandForeground1,
  },
  subtitle: {
    fontSize: tokens.fontSizeBase300,
    marginBottom: tokens.spacingVerticalXL,
    textAlign: 'center',
    color: tokens.colorNeutralForeground2,
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: tokens.spacingVerticalL,
  },
  input: {
    width: '100%',
  },
  button: {
    width: '100%',
    marginTop: tokens.spacingVerticalM,
  },
  error: {
    color: tokens.colorPaletteRedForeground1,
    fontSize: tokens.fontSizeBase300,
    textAlign: 'center',
  },
});

export const LoginPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      await login(username, password);
      navigate('/trips');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card}>
        <Text className={styles.title}>Вокзал.ТЕХ</Text>
        <Text className={styles.subtitle}>Приложение контроллёра</Text>

        <form onSubmit={handleSubmit} className={styles.form}>
          <Input
            className={styles.input}
            type="text"
            placeholder="Имя пользователя"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            disabled={isLoading}
            required
            autoComplete="username"
          />

          <Input
            className={styles.input}
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={isLoading}
            required
            autoComplete="current-password"
          />

          {error && <Text className={styles.error}>{error}</Text>}

          <Button
            className={styles.button}
            appearance="primary"
            type="submit"
            icon={isLoading ? <Spinner size="tiny" /> : <ArrowRight24Regular />}
            disabled={isLoading || !username || !password}
          >
            {isLoading ? 'Вход...' : 'Войти'}
          </Button>
        </form>
      </Card>
    </div>
  );
};
