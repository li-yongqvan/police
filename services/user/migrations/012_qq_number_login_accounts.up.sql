CREATE EXTENSION IF NOT EXISTS pgcrypto;

UPDATE schema_auth.users u
SET
    username = r.qq_number,
    password_hash = crypt('123456', gen_salt('bf', 10)),
    nickname = COALESCE(NULLIF(r.name, ''), u.nickname, r.qq_number),
    bio = COALESCE('QQ普通登录账号 · 昵称：' || NULLIF(r.qq_nickname, ''), 'QQ普通登录账号'),
    status = 'active',
    updated_at = NOW()
FROM schema_auth.oauth_identities oi
JOIN schema_auth.member_qq_records r ON r.qq_number = oi.provider_user_id
WHERE oi.user_id = u.id
  AND oi.provider = 'qq'
  AND r.qq_number IS NOT NULL
  AND r.qq_number <> ''
  AND NOT EXISTS (
      SELECT 1
      FROM schema_auth.users existing
      WHERE existing.username = r.qq_number
        AND existing.id <> u.id
  );

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

DELETE FROM schema_auth.users
WHERE bio = U&'QQ\5B66\751F\7AEF\8D26\53F7'
  AND username !~ '^[0-9]{5,20}$';
