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
import { RegisterRequest } from "@/types";

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
    maxWidth: "500px",
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
  row: {
    display: "grid",
    gridTemplateColumns: "1fr 1fr",
    gap: tokens.spacingHorizontalM,
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

export const RegisterPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const register = useAuthStore((state) => state.register);

  const [formData, setFormData] = useState<RegisterRequest>({
    email: "",
    phone: "",
    password: "",
    lastName: "",
    firstName: "",
    middleName: "",
    birthDate: "",
  });
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (field: keyof RegisterRequest, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.password !== confirmPassword) {
      setError("Пароли не совпадают");
      return;
    }

    if (formData.password.length < 8) {
      setError("Пароль должен содержать минимум 8 символов");
      return;
    }

    setLoading(true);
    try {
      await register(formData);
      navigate("/my-tickets");
    } catch {
      setError(
        "Ошибка регистрации. Возможно, пользователь с таким email уже существует."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.card}>
        <CardHeader
          header={<Text className={styles.title}>Регистрация</Text>}
        />
        <form className={styles.form} onSubmit={handleSubmit}>
          <Input
            type="email"
            placeholder="Email *"
            value={formData.email}
            onChange={(e) => handleChange("email", e.target.value)}
            required
          />
          <Input
            type="tel"
            placeholder="Телефон"
            value={formData.phone}
            onChange={(e) => handleChange("phone", e.target.value)}
          />
          <div className={styles.row}>
            <Input
              placeholder="Фамилия *"
              value={formData.lastName}
              onChange={(e) => handleChange("lastName", e.target.value)}
              required
            />
            <Input
              placeholder="Имя *"
              value={formData.firstName}
              onChange={(e) => handleChange("firstName", e.target.value)}
              required
            />
          </div>
          <Input
            placeholder="Отчество"
            value={formData.middleName}
            onChange={(e) => handleChange("middleName", e.target.value)}
          />
          <Input
            type="date"
            placeholder="Дата рождения"
            value={formData.birthDate}
            onChange={(e) => handleChange("birthDate", e.target.value)}
          />
          <Input
            type="password"
            placeholder="Пароль (минимум 8 символов) *"
            value={formData.password}
            onChange={(e) => handleChange("password", e.target.value)}
            required
          />
          <Input
            type="password"
            placeholder="Подтвердите пароль *"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
          {error && <Text className={styles.error}>{error}</Text>}
          <Button type="submit" appearance="primary" disabled={loading}>
            {loading ? "Регистрация..." : "Зарегистрироваться"}
          </Button>
          <div className={styles.link}>
            <Text>
              Уже есть аккаунт? <Link to="/login">Войти</Link>
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
