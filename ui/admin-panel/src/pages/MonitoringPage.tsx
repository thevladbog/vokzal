import React from "react";
import { useTranslation } from "react-i18next";
import { Title2, Text, Card, Link, makeStyles } from "@fluentui/react-components";
import { AppLayout } from "@/components/layout/AppLayout";

const useStyles = makeStyles({
  card: { padding: "24px", marginBottom: "16px" },
});

export const MonitoringPage: React.FC = () => {
  const { t } = useTranslation();
  const styles = useStyles();
  const grafanaUrl =
    import.meta.env.VITE_GRAFANA_URL || "http://localhost:30001";

  return (
    <AppLayout>
      <Title2 style={{ marginBottom: "24px" }}>{t("monitoring.title")}</Title2>
      <Card className={styles.card}>
        <Text block style={{ marginBottom: "16px" }}>
          {t("monitoring.description")}
        </Text>
        <Link
          href={`${grafanaUrl}/d/services-overview`}
          target="_blank"
        >
          {t("monitoring.servicesOverview")}
        </Link>
      </Card>
      <Card className={styles.card}>
        <Link
          href={`${grafanaUrl}/d/database-monitoring`}
          target="_blank"
        >
          {t("monitoring.databaseMonitoring")}
        </Link>
      </Card>
    </AppLayout>
  );
};
