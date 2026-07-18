UPDATE schema_admin.system_config
SET value = '0', updated_at = NOW()
WHERE key IN ('post_requires_level', 'comment_requires_level')
  AND value = '1';
