import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Card,
  CardHeader,
  Button,
  Input,
  Text,
  Radio,
  RadioGroup,
  Select,
  makeStyles,
  tokens,
  Spinner,
} from "@fluentui/react-components";
import { Add24Regular, Delete24Regular } from "@fluentui/react-icons";
import { useBookingStore } from "@/stores/bookingStore";
import { ticketService } from "@/services/ticket";
import { Passenger } from "@/types";
import { formatTime, formatPrice, formatDate } from "@/utils/format";

const useStyles = makeStyles({
  container: {
    minHeight: "100vh",
    backgroundColor: tokens.colorNeutralBackground2,
    padding: tokens.spacingVerticalXXL,
  },
  content: {
    maxWidth: "900px",
    margin: "0 auto",
  },
  header: {
    marginBottom: tokens.spacingVerticalXL,
  },
  title: {
    fontSize: tokens.fontSizeBase600,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalM,
  },
  tripInfo: {
    backgroundColor: tokens.colorNeutralBackground1,
    padding: tokens.spacingVerticalL,
    borderRadius: tokens.borderRadiusLarge,
    marginBottom: tokens.spacingVerticalXL,
  },
  section: {
    marginBottom: tokens.spacingVerticalXL,
  },
  sectionTitle: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalM,
  },
  passengerCard: {
    marginBottom: tokens.spacingVerticalM,
    padding: tokens.spacingVerticalL,
  },
  form: {
    display: "grid",
    gridTemplateColumns: "1fr 1fr",
    gap: tokens.spacingHorizontalM,
    "@media (max-width: 768px)": {
      gridTemplateColumns: "1fr",
    },
  },
  fullWidth: {
    gridColumn: "1 / -1",
  },
  buttonGroup: {
    display: "flex",
    gap: tokens.spacingHorizontalM,
    justifyContent: "space-between",
    marginTop: tokens.spacingVerticalXL,
  },
  summary: {
    backgroundColor: tokens.colorNeutralBackground1,
    padding: tokens.spacingVerticalL,
    borderRadius: tokens.borderRadiusLarge,
    marginTop: tokens.spacingVerticalXL,
  },
  summaryRow: {
    display: "flex",
    justifyContent: "space-between",
    marginBottom: tokens.spacingVerticalS,
  },
  total: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
    paddingTop: tokens.spacingVerticalM,
    borderTop: `1px solid ${tokens.colorNeutralStroke1}`,
  },
});

export const BookingPage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const {
    selectedTrip,
    passengers,
    paymentMethod,
    contactPhone,
    contactEmail,
    addPassenger,
    removePassenger,
    updatePassenger,
    setPaymentMethod,
    setContactPhone,
    setContactEmail,
    reset,
  } = useBookingStore();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  if (!selectedTrip) {
    navigate("/");
    return null;
  }

  const handleAddPassenger = () => {
    addPassenger({
      lastName: "",
      firstName: "",
      middleName: "",
      birthDate: "",
      documentType: "passport",
      documentNumber: "",
      phone: "",
      email: "",
      benefitType: "none",
      discount: 0,
    });
  };

  const handlePassengerChange = (
    index: number,
    field: keyof Passenger,
    value: Passenger[keyof Passenger]
  ) => {
    const updated = { ...passengers[index], [field]: value };
    updatePassenger(index, updated);
  };

  const calculateTotal = () => {
    return passengers.reduce((sum, passenger) => {
      const discount = passenger.discount || 0;
      return sum + selectedTrip.price * (1 - discount / 100);
    }, 0);
  };

  const handleSubmit = async () => {
    setError("");

    if (passengers.length === 0) {
      setError("Добавьте хотя бы одного пассажира");
      return;
    }

    if (!contactPhone) {
      setError("Укажите контактный телефон");
      return;
    }

    for (const passenger of passengers) {
      if (
        !passenger.lastName ||
        !passenger.firstName ||
        !passenger.birthDate ||
        !passenger.documentNumber
      ) {
        setError("Заполните все обязательные поля для каждого пассажира");
        return;
      }
    }

    setLoading(true);
    try {
      const result = await ticketService.sell({
        tripId: selectedTrip.id,
        passengers,
        paymentMethod,
        contactPhone,
        contactEmail,
      });

      // Navigate to payment or confirmation page
      if (paymentMethod === "cash") {
        navigate(
          `/confirmation?ticketIds=${result.tickets.map((t) => t.id).join(",")}`
        );
      } else {
        navigate(`/payment?paymentId=${result.payment.id}`);
      }
      reset();
    } catch (err: unknown) {
      const message =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { error?: string } } }).response?.data
              ?.error
          : undefined;
      setError(message || "Ошибка при оформлении билета");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <div className={styles.header}>
          <Text className={styles.title}>Оформление билета</Text>
        </div>

        <div className={styles.tripInfo}>
          <Text weight="semibold">Выбранный рейс</Text>
          <Text>
            {selectedTrip.route?.fromStation?.name} →{" "}
            {selectedTrip.route?.toStation?.name}
          </Text>
          <Text>
            {formatDate(selectedTrip.departureTime)} в{" "}
            {formatTime(selectedTrip.departureTime)}
          </Text>
          <Text>Цена: {formatPrice(selectedTrip.price)} за место</Text>
        </div>

        <div className={styles.section}>
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <Text className={styles.sectionTitle}>Пассажиры</Text>
            <Button
              icon={<Add24Regular />}
              onClick={handleAddPassenger}
              disabled={passengers.length >= selectedTrip.availableSeats}
            >
              Добавить пассажира
            </Button>
          </div>

          {passengers.map((passenger, index) => (
            <Card key={index} className={styles.passengerCard}>
              <CardHeader
                header={<Text weight="semibold">Пассажир {index + 1}</Text>}
                action={
                  <Button
                    icon={<Delete24Regular />}
                    appearance="subtle"
                    onClick={() => removePassenger(index)}
                  />
                }
              />
              <div className={styles.form}>
                <Input
                  placeholder="Фамилия *"
                  value={passenger.lastName}
                  onChange={(e) =>
                    handlePassengerChange(index, "lastName", e.target.value)
                  }
                  required
                />
                <Input
                  placeholder="Имя *"
                  value={passenger.firstName}
                  onChange={(e) =>
                    handlePassengerChange(index, "firstName", e.target.value)
                  }
                  required
                />
                <Input
                  placeholder="Отчество"
                  value={passenger.middleName}
                  onChange={(e) =>
                    handlePassengerChange(index, "middleName", e.target.value)
                  }
                />
                <Input
                  type="date"
                  placeholder="Дата рождения *"
                  value={passenger.birthDate}
                  onChange={(e) =>
                    handlePassengerChange(index, "birthDate", e.target.value)
                  }
                  required
                />
                <Select
                  value={passenger.documentType}
                  onChange={(e) =>
                    handlePassengerChange(index, "documentType", e.target.value)
                  }
                >
                  <option value="passport">Паспорт РФ</option>
                  <option value="birth_certificate">
                    Свидетельство о рождении
                  </option>
                  <option value="foreign_passport">Заграничный паспорт</option>
                </Select>
                {passenger.documentType === "passport" && (
                  <Input
                    placeholder="Серия паспорта"
                    value={passenger.documentSeries}
                    onChange={(e) =>
                      handlePassengerChange(
                        index,
                        "documentSeries",
                        e.target.value
                      )
                    }
                  />
                )}
                <Input
                  placeholder="Номер документа *"
                  value={passenger.documentNumber}
                  onChange={(e) =>
                    handlePassengerChange(
                      index,
                      "documentNumber",
                      e.target.value
                    )
                  }
                  required
                />
                <Select
                  value={passenger.benefitType}
                  onChange={(e) => {
                    const type = e.target.value as Passenger["benefitType"];
                    const discounts: Record<string, number> = {
                      none: 0,
                      child: 20,
                      student: 15,
                      pensioner: 10,
                      disabled: 50,
                    };
                    handlePassengerChange(index, "benefitType", type);
                    handlePassengerChange(
                      index,
                      "discount",
                      discounts[type || "none"] || 0
                    );
                  }}
                >
                  <option value="none">Без льгот</option>
                  <option value="child">Детский (скидка 20%)</option>
                  <option value="student">Студенческий (скидка 15%)</option>
                  <option value="pensioner">Пенсионный (скидка 10%)</option>
                  <option value="disabled">Инвалид (скидка 50%)</option>
                </Select>
              </div>
            </Card>
          ))}
        </div>

        <div className={styles.section}>
          <Text className={styles.sectionTitle}>Контактные данные</Text>
          <div className={styles.form}>
            <Input
              type="tel"
              placeholder="Телефон *"
              value={contactPhone}
              onChange={(e) => setContactPhone(e.target.value)}
              required
            />
            <Input
              type="email"
              placeholder="Email"
              value={contactEmail}
              onChange={(e) => setContactEmail(e.target.value)}
            />
          </div>
        </div>

        <div className={styles.section}>
          <Text className={styles.sectionTitle}>Способ оплаты</Text>
          <RadioGroup
            value={paymentMethod}
            onChange={(_, data) =>
              setPaymentMethod(
                (data.value as "card" | "sbp" | "cash") ?? "card"
              )
            }
          >
            <Radio value="card" label="Банковская карта" />
            <Radio value="sbp" label="СБП (Система быстрых платежей)" />
            <Radio value="cash" label="Наличные (оплата при посадке)" />
          </RadioGroup>
        </div>

        <div className={styles.summary}>
          <Text weight="semibold">Итого</Text>
          {passengers.map((passenger, index) => (
            <div key={index} className={styles.summaryRow}>
              <Text>
                Пассажир {index + 1}: {passenger.lastName} {passenger.firstName}
              </Text>
              <Text>
                {formatPrice(
                  selectedTrip.price * (1 - (passenger.discount || 0) / 100)
                )}
              </Text>
            </div>
          ))}
          <div className={`${styles.summaryRow} ${styles.total}`}>
            <Text>Всего к оплате:</Text>
            <Text>{formatPrice(calculateTotal())}</Text>
          </div>
        </div>

        {error && (
          <Text
            style={{
              color: tokens.colorPaletteRedForeground1,
              marginTop: tokens.spacingVerticalM,
            }}
          >
            {error}
          </Text>
        )}

        <div className={styles.buttonGroup}>
          <Button onClick={() => navigate("/")}>Отмена</Button>
          <Button
            appearance="primary"
            onClick={handleSubmit}
            disabled={loading || passengers.length === 0}
          >
            {loading ? <Spinner size="tiny" /> : "Оформить билет"}
          </Button>
        </div>
      </div>
    </div>
  );
};
