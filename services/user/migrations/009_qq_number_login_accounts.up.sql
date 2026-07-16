CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO schema_auth.users (
    username,
    password_hash,
    nickname,
    bio,
    avatar,
    level,
    status,
    department,
    squad,
    grade,
    profile_completed
)
SELECT DISTINCT
    qq_number,
    crypt('123456', gen_salt('bf', 10)),
    COALESCE(NULLIF(name, ''), qq_number),
    COALESCE('QQ普通登录账号 · 昵称：' || NULLIF(qq_nickname, ''), 'QQ普通登录账号'),
    '',
    1,
    'active',
    '',
    '',
    '',
    false
FROM schema_auth.member_qq_records
WHERE qq_number IS NOT NULL AND qq_number <> ''
ON CONFLICT (username) DO UPDATE
SET password_hash = EXCLUDED.password_hash,
    status = 'active',
    updated_at = NOW();
