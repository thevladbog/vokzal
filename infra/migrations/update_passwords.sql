-- Dev-only: set a known bcrypt hash for three test accounts only (disp1, cashier1, ctrl1).
-- NOT a general password migration — does not update all users or production accounts.
-- Hash MUST match 002_seed_test_data.sql; see README.md section "Тестовый пароль (dev only)".
-- Does NOT run by default. Requires explicit opt-in in the same session:
--   SET app.allow_password_reset = 'yes';
-- Then run this file. Rolls back on error or if opt-in was not set.

DO $$
DECLARE
  updated_count integer;
BEGIN
  IF current_setting('app.allow_password_reset', true) IS DISTINCT FROM 'yes' THEN
    RAISE EXCEPTION 'Aborted: set app.allow_password_reset = ''yes'' in this session to run (dev only)';
  END IF;

  UPDATE users
  SET password_hash = '$2a$10$gvTooaZGTXfKUqHRGn1xeuaqvwFqlzd5Z3BH7WBEFvJ9bBOa9xAMu'
  WHERE username IN ('disp1', 'cashier1', 'ctrl1');

  GET DIAGNOSTICS updated_count = ROW_COUNT;
  RAISE NOTICE 'Updated % test account(s).', updated_count;
EXCEPTION
  WHEN OTHERS THEN
    RAISE;
END $$;
