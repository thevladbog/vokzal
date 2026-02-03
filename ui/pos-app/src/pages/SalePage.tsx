import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Card,
  Title2,
  Button,
  Input,
  Text,
  makeStyles,
  Spinner,
  tokens,
} from "@fluentui/react-components";
import {
  SearchRegular,
  TicketDiagonalRegular,
  MoneyRegular,
} from "@fluentui/react-icons";
import { tripService } from "@/services/api";
import { posService } from "@/services/pos";
import { useSaleStore } from "@/stores/saleStore";
import { format } from "date-fns";
import { ru } from "date-fns/locale";
import type { Trip, SaleRequest } from "@/types";

const useStyles = makeStyles({
  container: {
    padding: "24px",
    display: "grid",
    gridTemplateColumns: "2fr 1fr",
    gap: "24px",
    height: "100vh",
  },
  leftPanel: {
    display: "flex",
    flexDirection: "column",
    gap: "16px",
  },
  rightPanel: {
    display: "flex",
    flexDirection: "column",
    gap: "16px",
  },
  searchBar: {
    display: "flex",
    gap: "12px",
    alignItems: "center",
  },
  tripCard: {
    padding: "16px",
    cursor: "pointer",
    ":hover": {
      backgroundColor: tokens.colorNeutralBackground1Hover,
    },
  },
  selectedTrip: {
    backgroundColor: tokens.colorBrandBackground2,
  },
  saleForm: {
    display: "flex",
    flexDirection: "column",
    gap: "12px",
  },
  total: {
    fontSize: "24px",
    fontWeight: "bold",
    color: tokens.colorBrandForeground1,
  },
});

export const SalePage: React.FC = () => {
  const styles = useStyles();
  const { setSaleInProgress, setCurrentTicket } = useSaleStore();

  const [selectedTrip, setSelectedTrip] = useState<Trip | null>(null);
  const [passengerFio, setPassengerFio] = useState("");
  const [passengerPhone, setPassengerPhone] = useState("");
  const [selling, setSelling] = useState(false);

  const { data: trips, isLoading } = useQuery({
    queryKey: ["trips"],
    queryFn: () => tripService.getTrips({ status: "scheduled" }),
    refetchInterval: 30000,
  });

  const handleSell = async () => {
    if (!selectedTrip) return;

    setSelling(true);
    setSaleInProgress(true);

    try {
      const request: SaleRequest = {
        trip_id: selectedTrip.id,
        passenger_fio: passengerFio || undefined,
        passenger_phone: passengerPhone || undefined,
      };

      // Продать билет через Tauri backend
      const ticket = await posService.sellTicket(request);
      setCurrentTicket(ticket);

      // Печать билета
      await posService.printTicket({
        ticket_id: ticket.id,
        route: selectedTrip.route_name,
        date: format(new Date(selectedTrip.departure_datetime), "dd.MM.yyyy", {
          locale: ru,
        }),
        time: format(new Date(selectedTrip.departure_datetime), "HH:mm", {
          locale: ru,
        }),
        platform: selectedTrip.platform || "-",
        price: ticket.price,
        passenger_fio: passengerFio,
        qr_code: ticket.qr_code,
        bar_code: ticket.bar_code,
      });

      alert(`Билет продан успешно!\nID: ${ticket.id}`);

      // Сбросить форму
      setPassengerFio("");
      setPassengerPhone("");
      setSelectedTrip(null);
    } catch (err: unknown) {
      alert(
        `Ошибка продажи: ${err instanceof Error ? err.message : String(err)}`
      );
    } finally {
      setSelling(false);
      setSaleInProgress(false);
    }
  };

  if (isLoading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Spinner label="Загрузка рейсов..." />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Левая панель - список рейсов */}
      <div className={styles.leftPanel}>
        <Card>
          <div className={styles.searchBar}>
            <SearchRegular fontSize={24} />
            <Input placeholder="Поиск рейсов..." style={{ flex: 1 }} />
          </div>
        </Card>

        <Card style={{ flex: 1, overflow: "auto" }}>
          <Title2 style={{ marginBottom: "16px" }}>Доступные рейсы</Title2>

          {trips?.map((trip) => (
            <Card
              key={trip.id}
              className={`${styles.tripCard} ${
                selectedTrip?.id === trip.id ? styles.selectedTrip : ""
              }`}
              onClick={() => setSelectedTrip(trip)}
            >
              <Text weight="bold" block>
                {trip.route_name}
              </Text>
              <Text block>
                {format(new Date(trip.departure_datetime), "dd.MM.yyyy HH:mm", {
                  locale: ru,
                })}
              </Text>
              <Text block>
                Свободных мест: {trip.available_seats} | Цена: {trip.price} ₽
              </Text>
              {trip.platform && <Text block>Перрон: {trip.platform}</Text>}
            </Card>
          ))}

          {(!trips || trips.length === 0) && <Text>Нет доступных рейсов</Text>}
        </Card>
      </div>

      {/* Правая панель - продажа */}
      <div className={styles.rightPanel}>
        <Card style={{ padding: "24px" }}>
          <Title2 style={{ marginBottom: "16px" }}>Продажа билета</Title2>

          {selectedTrip ? (
            <div className={styles.saleForm}>
              <Card
                style={{
                  padding: "16px",
                  backgroundColor: tokens.colorNeutralBackground3,
                }}
              >
                <Text weight="bold" block>
                  {selectedTrip.route_name}
                </Text>
                <Text block>
                  {format(
                    new Date(selectedTrip.departure_datetime),
                    "dd.MM.yyyy HH:mm",
                    { locale: ru }
                  )}
                </Text>
                <Text block>Перрон: {selectedTrip.platform || "-"}</Text>
                <div className={styles.total} style={{ marginTop: "12px" }}>
                  {selectedTrip.price} ₽
                </div>
              </Card>

              <Input
                placeholder="ФИО пассажира (опционально)"
                value={passengerFio}
                onChange={(e) => setPassengerFio(e.target.value)}
              />

              <Input
                placeholder="Телефон пассажира (опционально)"
                value={passengerPhone}
                onChange={(e) => setPassengerPhone(e.target.value)}
              />

              <Button
                appearance="primary"
                size="large"
                icon={<TicketDiagonalRegular />}
                onClick={handleSell}
                disabled={selling}
              >
                {selling ? "Продажа..." : "Продать билет"}
              </Button>

              <Button
                appearance="secondary"
                onClick={() => setSelectedTrip(null)}
              >
                Отменить
              </Button>
            </div>
          ) : (
            <Text>Выберите рейс для продажи</Text>
          )}
        </Card>

        <Card style={{ padding: "24px" }}>
          <Button
            appearance="outline"
            size="large"
            icon={<MoneyRegular />}
            style={{ width: "100%" }}
          >
            Открыть кассу
          </Button>
        </Card>
      </div>
    </div>
  );
};
