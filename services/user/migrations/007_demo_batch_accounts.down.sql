DELETE FROM schema_auth.users
WHERE username IN (
    'demo01', 'demo02', 'demo03', 'demo04', 'demo05', 'demo06',
    'admin01', 'admin02', 'admin03', 'admin04', 'admin05', 'admin06',
    'plat01', 'plat02', 'plat03', 'plat04', 'plat05', 'plat06'
);
