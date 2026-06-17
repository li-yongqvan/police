-- Idempotent: assign roles if users were created after 005 ran (startup race).
INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by)
SELECT u.id, 1, 0
FROM schema_auth.users u
WHERE u.username IN ('admin01', 'admin02', 'admin03', 'admin04', 'admin05', 'admin06')
  AND NOT EXISTS (
    SELECT 1 FROM schema_admin.user_roles ur
    WHERE ur.user_id = u.id AND ur.role_id = 1
  );

INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by)
SELECT u.id, 2, 0
FROM schema_auth.users u
WHERE u.username IN ('plat01', 'plat02', 'plat03', 'plat04', 'plat05', 'plat06')
  AND NOT EXISTS (
    SELECT 1 FROM schema_admin.user_roles ur
    WHERE ur.user_id = u.id AND ur.role_id = 2
  );
