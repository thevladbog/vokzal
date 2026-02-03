import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
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
} from "@fluentui/react-components";
import { authService } from "@/services/auth";
import { useAuthStore } from "@/stores/authStore";

const useStyles = makeStyles({
  container: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    minHeight: "100vh",
    backgroundColor: tokens.colorNeutralBackground2,
  },
  card: {
    width: "400px",
    padding: "32px",
  },
  title: {
    marginBottom: "24px",
    textAlign: "center",
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: "16px",
  },
  error: {
    color: tokens.colorPaletteRedForeground1,
    fontSize: "12px",
  },
});

export const LoginPage: React.FC = () => {
  const styles = useStyles();
  const { t } = useTranslation();
  const navigate = useNavigate();
  const setUser = useAuthStore((state) => state.setUser);

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await authService.login(username, password);
      setUser(response.user);
      navigate("/");
    } catch (err: unknown) {
      const msg =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { error?: string } } }).response?.data
              ?.error
          : null;
      setError(msg || "Ошибка входа. Проверьте логин и пароль.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <FluentProvider theme={webLightTheme}>
      <div className={styles.container}>
        <Card className={styles.card}>
          <Title2 className={styles.title}>{t('login.title')}</Title2>
          <Text block style={{ textAlign: "center", marginBottom: "24px" }}>
            {t('login.subtitle')}
          </Text>

          <form onSubmit={handleSubmit} className={styles.form}>
            <Input
              type="text"
              placeholder={t('login.username')}
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <Input
              type="password"
              placeholder={t('login.password')}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />

            {error && <Text className={styles.error}>{error}</Text>}

            <Button appearance="primary" type="submit" disabled={loading}>
              {loading ? t('login.submitting') : t('login.submit')}
            </Button>
          </form>
        </Card>
      </div>
    </FluentProvider>
  );
};
