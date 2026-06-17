-- Assign admin roles to demo accounts (user ids resolved by username)
INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by)
SELECT u.id, 1, 0
FROM schema_auth.users u
WHERE u.username = 'demo_admin'
ON CONFLICT DO NOTHING;

INSERT INTO schema_admin.user_roles (user_id, role_id, assigned_by)
SELECT u.id, 2, 0
FROM schema_auth.users u
WHERE u.username = 'demo_platform_admin'
ON CONFLICT DO NOTHING;
