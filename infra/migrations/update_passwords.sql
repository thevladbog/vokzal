-- Update all users with password: admin123
-- bcrypt hash: $2a$10$gvTooaZGTXfKUqHRGn1xeuaqvwFqlzd5Z3BH7WBEFvJ9bBOa9xAMu

UPDATE users SET password_hash = '$2a$10$gvTooaZGTXfKUqHRGn1xeuaqvwFqlzd5Z3BH7WBEFvJ9bBOa9xAMu';

SELECT username, role, LEFT(password_hash, 20) as hash_check FROM users;
