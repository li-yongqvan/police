DELETE FROM schema_auth.users WHERE username IN ('demo_student', 'demo_admin', 'demo_platform_admin');
DELETE FROM schema_auth.invite_codes WHERE code = 'DEMO2026';
