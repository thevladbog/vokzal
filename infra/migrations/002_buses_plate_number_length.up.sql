-- Migration: 002_buses_plate_number_length
-- Description: Увеличение длины plate_number до VARCHAR(20) для международных форматов номеров

ALTER TABLE buses
    ALTER COLUMN plate_number TYPE VARCHAR(20);
