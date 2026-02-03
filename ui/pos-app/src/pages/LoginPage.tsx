import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Card,
  Input,
  Button,
  Text,
  Title2,
  makeStyles,
  tokens,
} from "@fluentui/react-components";
import { authService } from "@/services/api";
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

      if (response.user.role !== "cashier" && response.user.role !== "admin") {
        setError('Доступ запрещён. Требуется роль "Кассир".');
        await authService.logout();
        return;
      }

      setUser(response.user);
      navigate("/");
    } catch (err: unknown) {
      const message =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { error?: string } } }).response?.data
              ?.error
          : undefined;
      setError(message || "Ошибка входа. Проверьте логин и пароль.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card}>
        <Title2 className={styles.title}>Вокзал.ТЕХ POS</Title2>
        <Text block style={{ textAlign: "center", marginBottom: "24px" }}>
          Вход для кассира
        </Text>

        <form onSubmit={handleSubmit} className={styles.form}>
          <Input
            type="text"
            placeholder="Логин"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            autoFocus
          />
          <Input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />

          {error && <Text className={styles.error}>{error}</Text>}

          <Button appearance="primary" type="submit" disabled={loading}>
            {loading ? "Вход..." : "Войти"}
          </Button>
        </form>
      </Card>
    </div>
  );
};
