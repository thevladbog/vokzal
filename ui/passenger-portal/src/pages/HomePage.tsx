import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Button,
  Input,
  Text,
  Combobox,
  Option,
  makeStyles,
  tokens,
  Spinner,
} from "@fluentui/react-components";
import {
  ArrowSwap24Regular,
  Search24Regular,
  CalendarLtr24Regular,
} from "@fluentui/react-icons";
import { useSearchStore } from "@/stores/searchStore";
import { useBookingStore } from "@/stores/bookingStore";
import { stationService } from "@/services/station";
import { Station, Trip } from "@/types";
import { TripCard } from "@/components/TripCard";
import { formatDate } from "@/utils/format";

const useStyles = makeStyles({
  container: {
    minHeight: "100vh",
    backgroundColor: tokens.colorNeutralBackground2,
  },
  header: {
    backgroundColor: tokens.colorBrandBackground,
    color: tokens.colorNeutralForegroundInverted,
    padding: tokens.spacingVerticalXXL,
    textAlign: "center",
  },
  title: {
    fontSize: tokens.fontSizeHero900,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalS,
  },
  subtitle: {
    fontSize: tokens.fontSizeBase400,
    opacity: 0.9,
  },
  searchSection: {
    maxWidth: "900px",
    margin: "0 auto",
    padding: tokens.spacingVerticalXXL,
    transform: "translateY(-50%)",
  },
  searchCard: {
    backgroundColor: tokens.colorNeutralBackground1,
    borderRadius: tokens.borderRadiusLarge,
    padding: tokens.spacingVerticalXL,
    boxShadow: tokens.shadow16,
  },
  searchForm: {
    display: "grid",
    gridTemplateColumns: "1fr auto 1fr auto",
    gap: tokens.spacingHorizontalM,
    alignItems: "end",
    "@media (max-width: 768px)": {
      gridTemplateColumns: "1fr",
    },
  },
  swapButton: {
    alignSelf: "end",
  },
  dateInput: {
    width: "100%",
  },
  resultsSection: {
    maxWidth: "900px",
    margin: "0 auto",
    padding: `0 ${tokens.spacingVerticalXXL} ${tokens.spacingVerticalXXL}`,
  },
  resultsHeader: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalL,
  },
  emptyState: {
    textAlign: "center",
    padding: tokens.spacingVerticalXXL,
    color: tokens.colorNeutralForeground2,
  },
  loading: {
    display: "flex",
    justifyContent: "center",
    padding: tokens.spacingVerticalXXL,
  },
});

export const HomePage: React.FC = () => {
  const styles = useStyles();
  const navigate = useNavigate();
  const {
    fromStation,
    toStation,
    date,
    trips,
    isSearching,
    setFromStation,
    setToStation,
    setDate,
    searchTrips,
    swapStations,
  } = useSearchStore();
  const { selectTrip } = useBookingStore();

  const [stations, setStations] = useState<Station[]>([]);
  const [loadingStations, setLoadingStations] = useState(true);

  useEffect(() => {
    const loadStations = async () => {
      try {
        const data = await stationService.getAll();
        setStations(data.filter((s) => s.active));
      } catch (error) {
        console.error("Failed to load stations:", error);
      } finally {
        setLoadingStations(false);
      }
    };
    loadStations();
  }, []);

  const handleSearch = async () => {
    try {
      await searchTrips();
    } catch (error) {
      console.error("Search failed:", error);
    }
  };

  const handleTripSelect = (trip: Trip) => {
    selectTrip(trip);
    navigate("/booking");
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <Text className={styles.title}>Вокзал.ТЕХ</Text>
        <Text className={styles.subtitle}>
          Покупка автобусных билетов онлайн
        </Text>
      </header>

      <div className={styles.searchSection}>
        <div className={styles.searchCard}>
          <div className={styles.searchForm}>
            <div>
              <Text weight="semibold">Откуда</Text>
              <Combobox
                placeholder="Выберите станцию"
                value={fromStation?.name || ""}
                disabled={loadingStations}
                onOptionSelect={(_, data) => {
                  const station = stations.find(
                    (s) => s.id === data.optionValue
                  );
                  setFromStation(station || null);
                }}
              >
                {stations.map((station) => (
                  <Option
                    key={station.id}
                    value={station.id}
                    text={station.name}
                  >
                    {station.name} ({station.code})
                  </Option>
                ))}
              </Combobox>
            </div>

            <Button
              className={styles.swapButton}
              icon={<ArrowSwap24Regular />}
              appearance="subtle"
              onClick={swapStations}
            />

            <div>
              <Text weight="semibold">Куда</Text>
              <Combobox
                placeholder="Выберите станцию"
                value={toStation?.name || ""}
                disabled={loadingStations}
                onOptionSelect={(_, data) => {
                  const station = stations.find(
                    (s) => s.id === data.optionValue
                  );
                  setToStation(station || null);
                }}
              >
                {stations.map((station) => (
                  <Option
                    key={station.id}
                    value={station.id}
                    text={station.name}
                  >
                    {station.name} ({station.code})
                  </Option>
                ))}
              </Combobox>
            </div>

            <div>
              <Text weight="semibold">Дата</Text>
              <Input
                className={styles.dateInput}
                type="date"
                value={formatDate(date, "yyyy-MM-dd")}
                onChange={(e) => setDate(new Date(e.target.value))}
                contentBefore={<CalendarLtr24Regular />}
                min={formatDate(new Date(), "yyyy-MM-dd")}
              />
            </div>
          </div>

          <Button
            appearance="primary"
            icon={<Search24Regular />}
            onClick={handleSearch}
            disabled={!fromStation || !toStation || isSearching}
            style={{ marginTop: tokens.spacingVerticalL, width: "100%" }}
          >
            {isSearching ? "Поиск..." : "Найти билеты"}
          </Button>
        </div>
      </div>

      <div className={styles.resultsSection}>
        {isSearching && (
          <div className={styles.loading}>
            <Spinner label="Поиск рейсов..." />
          </div>
        )}

        {!isSearching && trips.length > 0 && (
          <>
            <Text className={styles.resultsHeader}>
              Найдено рейсов: {trips.length}
            </Text>
            {trips.map((trip) => (
              <TripCard key={trip.id} trip={trip} onSelect={handleTripSelect} />
            ))}
          </>
        )}

        {!isSearching && trips.length === 0 && fromStation && toStation && (
          <div className={styles.emptyState}>
            <Text>Рейсы не найдены. Попробуйте изменить параметры поиска.</Text>
          </div>
        )}
      </div>
    </div>
  );
};
