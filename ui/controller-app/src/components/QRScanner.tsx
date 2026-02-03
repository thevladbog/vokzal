import React, { useCallback, useEffect, useRef, useState } from "react";
import { Html5Qrcode } from "html5-qrcode";
import { Button, Text, makeStyles, tokens } from "@fluentui/react-components";
import { Camera24Regular, Dismiss24Regular } from "@fluentui/react-icons";

const useStyles = makeStyles({
  container: {
    width: "100%",
    maxWidth: "600px",
    margin: "0 auto",
  },
  video: {
    width: "100%",
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground3,
  },
  controls: {
    display: "flex",
    gap: tokens.spacingHorizontalM,
    marginTop: tokens.spacingVerticalM,
    justifyContent: "center",
  },
  status: {
    textAlign: "center",
    marginTop: tokens.spacingVerticalS,
    color: tokens.colorNeutralForeground2,
  },
});

interface QRScannerProps {
  onScan: (qrCode: string) => void;
  onError?: (error: string) => void;
  isActive: boolean;
}

export const QRScanner: React.FC<QRScannerProps> = ({
  onScan,
  onError,
  isActive,
}) => {
  const styles = useStyles();
  const scannerRef = useRef<Html5Qrcode | null>(null);
  const [isScanning, setIsScanning] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const stopScanning = useCallback(async () => {
    if (scannerRef.current) {
      try {
        await scannerRef.current.stop();
        scannerRef.current.clear();
      } catch (err) {
        console.error("Error stopping scanner:", err);
      }
      scannerRef.current = null;
    }
    setIsScanning(false);
  }, []);

  const startScanning = useCallback(async () => {
    try {
      setError(null);
      setIsScanning(true);

      const scanner = new Html5Qrcode("qr-reader");
      scannerRef.current = scanner;

      await scanner.start(
        { facingMode: "environment" },
        {
          fps: 10,
          qrbox: { width: 250, height: 250 },
        },
        (decodedText) => {
          onScan(decodedText);
          stopScanning();
        },
        (errorMessage) => {
          // Ignore scanning errors (happens constantly during normal operation)
          console.debug("QR scan error:", errorMessage);
        }
      );
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : "Ошибка доступа к камере";
      setError(errorMsg);
      setIsScanning(false);
      onError?.(errorMsg);
    }
  }, [onScan, onError, stopScanning]);

  useEffect(() => {
    if (isActive && !isScanning) {
      startScanning();
    }

    return () => {
      stopScanning();
    };
  }, [isActive, isScanning, startScanning, stopScanning]);

  return (
    <div className={styles.container}>
      <div id="qr-reader" className={styles.video}></div>

      {error && (
        <Text
          className={styles.status}
          style={{ color: tokens.colorPaletteRedForeground1 }}
        >
          {error}
        </Text>
      )}

      {isScanning && !error && (
        <Text className={styles.status}>Наведите камеру на QR-код билета</Text>
      )}

      <div className={styles.controls}>
        {!isScanning ? (
          <Button
            appearance="primary"
            icon={<Camera24Regular />}
            onClick={startScanning}
            disabled={!isActive}
          >
            Начать сканирование
          </Button>
        ) : (
          <Button
            appearance="secondary"
            icon={<Dismiss24Regular />}
            onClick={stopScanning}
          >
            Остановить
          </Button>
        )}
      </div>
    </div>
  );
};
