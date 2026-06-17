DELETE FROM schema_admin.user_roles
WHERE user_id IN (
    SELECT id FROM schema_auth.users
    WHERE username IN (
        'admin01', 'admin02', 'admin03', 'admin04', 'admin05', 'admin06',
        'plat01', 'plat02', 'plat03', 'plat04', 'plat05', 'plat06'
    )
);
