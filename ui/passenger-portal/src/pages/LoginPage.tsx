import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import {
  Card,
  CardHeader,
  Input,
  Button,
  Text,
  makeStyles,
  tokens,
} from "@fluentui/react-components";
import { useAuthStore } from "@/stores/authStore";

const useStyles = makeStyles({
  container: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    minHeight: "100vh",
    backgroundColor: tokens.colorNeutralBackground2,
    padding: tokens.spacingVerticalXXL,
  },
  card: {
    width: "100%",
    maxWidth: "400px",
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalM,
    padding: tokens.spacingVerticalL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    textAlign: "center",
    marginBottom: tokens.spacingVerticalM,
  },
  error: {
    color: tokens.colorPaletteRedForeground1,
    fontSize: tokens.fontSizeBase300,
  },
  link: {
    textAlign: "center",
    marginTop: tokens.spacingVerticalS,
  },
});

export const LoginPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await login(email, password);
      navigate("/my-tickets");
    } catch {
      setError("Неверный email или пароль");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card}>
        <CardHeader
          header={<Text className={styles.title}>Вход в личный кабинет</Text>}
        />
        <form className={styles.form} onSubmit={handleSubmit}>
          <Input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
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
          <Button type="submit" appearance="primary" disabled={loading}>
            {loading ? "Вход..." : "Войти"}
          </Button>
          <div className={styles.link}>
            <Text>
              Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
            </Text>
          </div>
          <div className={styles.link}>
            <Link to="/">Вернуться на главную</Link>
          </div>
        </form>
      </Card>
    </div>
  );
};
