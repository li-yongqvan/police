DELETE FROM schema_auth.users
WHERE username IN (
    SELECT qq_number
    FROM schema_auth.member_qq_records
    WHERE qq_number IS NOT NULL AND qq_number <> ''
)
AND bio LIKE 'QQ普通登录账号%';
