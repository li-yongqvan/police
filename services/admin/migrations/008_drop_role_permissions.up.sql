-- Drop the never-consumed permission matrix: authorization is expressed by the role name alone (CONTEXT.md: Role).
ALTER TABLE schema_admin.roles DROP COLUMN IF EXISTS permissions;
