-- Assign admin roles to batch demo accounts (password equals username)
INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by)
SELECT u.id, 1, 0
FROM schema_auth.users u
WHERE u.username IN ('admin01', 'admin02', 'admin03', 'admin04', 'admin05', 'admin06')
ON CONFLICT DO NOTHING;

INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by)
SELECT u.id, 2, 0
FROM schema_auth.users u
WHERE u.username IN ('plat01', 'plat02', 'plat03', 'plat04', 'plat05', 'plat06')
ON CONFLICT DO NOTHING;
