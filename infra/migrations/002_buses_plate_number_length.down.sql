-- Migration: 002_buses_plate_number_length (rollback)
-- Description: Возврат plate_number к VARCHAR(12)

-- Только если все значения укладываются в 12 символов
ALTER TABLE buses
    ALTER COLUMN plate_number TYPE VARCHAR(12);
