-- Migration: 002_buses_plate_number_length (rollback)
-- Description: Возврат plate_number к VARCHAR(12) после проверки, что нет значений длиннее 12 символов

DO $$
DECLARE
  cnt INTEGER;
  example_plate TEXT;
BEGIN
  SELECT COUNT(*) INTO cnt FROM buses WHERE char_length(plate_number) > 12;
  IF cnt > 0 THEN
    SELECT plate_number INTO example_plate FROM buses WHERE char_length(plate_number) > 12 LIMIT 1;
    RAISE EXCEPTION 'Cannot rollback buses.plate_number to VARCHAR(12): % row(s) exceed 12 characters. Example: %', cnt, COALESCE(example_plate, '?');
  END IF;
END $$;

ALTER TABLE buses
  ALTER COLUMN plate_number TYPE VARCHAR(12);
