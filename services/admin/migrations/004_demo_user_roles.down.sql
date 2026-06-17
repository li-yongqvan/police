DELETE FROM schema_admin.user_roles
WHERE user_id IN (
    SELECT id FROM schema_auth.users
    WHERE username IN ('demo_admin', 'demo_platform_admin')
);
