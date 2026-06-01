ALTER TABLE schema_auth.users
  DROP COLUMN IF EXISTS profile_completed,
  DROP COLUMN IF EXISTS grade,
  DROP COLUMN IF EXISTS squad,
  DROP COLUMN IF EXISTS department;
