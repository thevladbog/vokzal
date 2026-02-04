-- Simple seed data for Vokzal.T??

-- Stations
INSERT INTO stations (name, code, address, timezone) VALUES
('?????? ???????????', 'MSK-BEL', '??????', 'Europe/Moscow'),
('?????-?????????', 'SPB-MAIN', '???', 'Europe/Moscow'),
('?????', 'TVR-MAIN', '?????', 'Europe/Moscow'),
('?????????', 'YAR-MAIN', '?????????', 'Europe/Moscow')
ON CONFLICT DO NOTHING;

-- Buses (reference stations by code for simplicity)
WITH s AS (
  SELECT id, code FROM stations WHERE code IN ('MSK-BEL', 'SPB-MAIN', 'TVR-MAIN', 'YAR-MAIN')
)
INSERT INTO buses (plate_number, model, capacity, status, station_id) 
SELECT * FROM (
  VALUES
    ('?001??99'::varchar, '???????? Sprinter'::varchar, 50::int, 'active'::varchar, (SELECT id FROM s WHERE code='MSK-BEL')),
    ('?002??99'::varchar, '???????? 5283'::varchar, 55::int, 'active'::varchar, (SELECT id FROM s WHERE code='MSK-BEL')),
    ('?001??78'::varchar, '???????? Sprinter'::varchar, 50::int, 'active'::varchar, (SELECT id FROM s WHERE code='SPB-MAIN')),
    ('?001??69'::varchar, '???? 5292'::varchar, 45::int, 'active'::varchar, (SELECT id FROM s WHERE code='TVR-MAIN')),
    ('?001??76'::varchar, '???????? Sprinter'::varchar, 50::int, 'active'::varchar, (SELECT id FROM s WHERE code='YAR-MAIN'))
) AS v(plate_number, model, capacity, status, station_id)
ON CONFLICT (plate_number) DO NOTHING;

-- Drivers
WITH s AS (
  SELECT id, code FROM stations WHERE code IN ('MSK-BEL', 'SPB-MAIN', 'TVR-MAIN', 'YAR-MAIN')
)
INSERT INTO drivers (full_name, license_number, experience_years, phone, station_id)
SELECT * FROM (
  VALUES
    ('?????? ?.?.'::varchar, '??0001'::varchar, 15::int, '+79001'::varchar, (SELECT id FROM s WHERE code='MSK-BEL')),
    ('?????? ?.?.'::varchar, '??0002'::varchar, 8::int, '+79002'::varchar, (SELECT id FROM s WHERE code='MSK-BEL')),
    ('??????? ?.?.'::varchar, '??0003'::varchar, 20::int, '+79003'::varchar, (SELECT id FROM s WHERE code='SPB-MAIN')),
    ('???????? ?.?.'::varchar, '??0004'::varchar, 12::int, '+79004'::varchar, (SELECT id FROM s WHERE code='TVR-MAIN')),
    ('??????? ?.?.'::varchar, '??0005'::varchar, 10::int, '+79005'::varchar, (SELECT id FROM s WHERE code='YAR-MAIN'))
) AS v(full_name, license_number, experience_years, phone, station_id)
ON CONFLICT (license_number) DO NOTHING;

-- Routes
INSERT INTO routes (name, stops, distance_km, duration_min, is_active) VALUES
('?-001: ??????  ???', '[]'::jsonb, 690, 480, true),
('?-002: ??????  ?????', '[]'::jsonb, 160, 120, true),
('?-003: ??????  ?????????', '[]'::jsonb, 250, 180, true)
ON CONFLICT DO NOTHING;

-- Schedules
WITH r AS (
  SELECT id, name FROM routes WHERE name LIKE '?-00%'
)
INSERT INTO schedules (route_id, departure_time, days_of_week, platform, is_active)
SELECT * FROM (
  VALUES
    ((SELECT id FROM r WHERE name = '?-001: ??????  ???'), '08:00:00'::time, '[1,2,3,4,5,6,7]'::jsonb, '1'::varchar, true),
    ((SELECT id FROM r WHERE name = '?-001: ??????  ???'), '18:00:00'::time, '[1,2,3,4,5,6,7]'::jsonb, '2'::varchar, true),
    ((SELECT id FROM r WHERE name = '?-002: ??????  ?????'), '07:00:00'::time, '[1,2,3,4,5,6,7]'::jsonb, '3'::varchar, true),
    ((SELECT id FROM r WHERE name = '?-002: ??????  ?????'), '12:00:00'::time, '[1,2,3,4,5,6,7]'::jsonb, '1'::varchar, true),
    ((SELECT id FROM r WHERE name = '?-003: ??????  ?????????'), '09:30:00'::time, '[1,3,5,6,7]'::jsonb, '4'::varchar, true)
) AS v(route_id, departure_time, days_of_week, platform, is_active)
ON CONFLICT DO NOTHING;

-- Users
WITH s AS (
  SELECT id FROM stations WHERE code = 'MSK-BEL' LIMIT 1
)
INSERT INTO users (username, password_hash, full_name, role, station_id, is_active)
SELECT * FROM (
  VALUES
    ('disp1'::varchar, '$2a$10$xNjK0Z5Z5Z5Z5Z5Z5Z5ZuO.n4VvvE1U/r5DvXm4UMZnIYzZb0QVnG'::varchar, 'Диспетчер 1'::varchar, 'dispatcher'::varchar, (SELECT id FROM s), true),
    ('cashier1'::varchar, '$2a$10$xNjK0Z5Z5Z5Z5Z5Z5Z5ZuO.n4VvvE1U/r5DvXm4UMZnIYzZb0QVnG'::varchar, 'Кассир 1'::varchar, 'cashier'::varchar, (SELECT id FROM s), true),
    ('ctrl1'::varchar, '$2a$10$xNjK0Z5Z5Z5Z5Z5Z5Z5ZuO.n4VvvE1U/r5DvXm4UMZnIYzZb0QVnG'::varchar, 'Контролер 1'::varchar, 'controller'::varchar, (SELECT id FROM s), true)
) AS v(username, password_hash, full_name, role, station_id, is_active)
ON CONFLICT (username) DO NOTHING;

-- Trips
WITH b AS (SELECT id FROM buses WHERE plate_number='?001??99' LIMIT 1),
     d AS (SELECT id FROM drivers WHERE license_number='??0001' LIMIT 1),
     sch AS (SELECT id FROM schedules LIMIT 3)
INSERT INTO trips (schedule_id, date, status, bus_id, driver_id)
SELECT sch.id, CURRENT_DATE + INTERVAL '1 day', 'scheduled'::varchar, b.id, d.id
FROM sch, b, d
ON CONFLICT (schedule_id, date) DO NOTHING;

-- Summary
SELECT 'Stations: ' || COUNT(*) FROM stations;
SELECT 'Buses: ' || COUNT(*) FROM buses;
SELECT 'Drivers: ' || COUNT(*) FROM drivers;
SELECT 'Routes: ' || COUNT(*) FROM routes;
SELECT 'Schedules: ' || COUNT(*) FROM schedules;
SELECT 'Users: ' || COUNT(*) FROM users;
SELECT 'Trips: ' || COUNT(*) FROM trips;
