-- 002_admin_tables.up.sql

-- 1. System configuration table
CREATE TABLE schema_admin.system_config (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT DEFAULT '',
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT DEFAULT 0
);

-- Seed initial configuration
INSERT INTO schema_admin.system_config (key, value, description) VALUES
    ('post_max_title_length', '200', '帖子标题最大长度'),
    ('post_max_content_length', '50000', '帖子内容最大长度'),
    ('post_requires_level', '0', '发帖所需最低等级'),
    ('comment_requires_level', '0', '评论所需最低等级'),
    ('upload_requires_level', '2', '上传附件所需最低等级'),
    ('max_attachment_size_mb', '20', '单个附件最大大小（MB）'),
    ('board_ai-learning_enabled', 'true', 'AI学习交流区开关'),
    ('board_announcements_enabled', 'true', '协会公告区开关'),
    ('board_tech-help_enabled', 'true', '技术问答区开关'),
    ('sensitive_word_action', 'pending_review', '敏感词处理方式: pending_review / reject'),
    ('daily_post_limit_per_user', '50', '每用户每日发帖上限')
ON CONFLICT (key) DO NOTHING;

-- 2. Operation logs table
CREATE TABLE schema_admin.operation_logs (
    id SERIAL PRIMARY KEY,
    operator_id BIGINT NOT NULL,
    operator_username VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id BIGINT,
    detail JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_operation_logs_operator ON schema_admin.operation_logs(operator_id);
CREATE INDEX idx_operation_logs_action ON schema_admin.operation_logs(action);
CREATE INDEX idx_operation_logs_created ON schema_admin.operation_logs(created_at DESC);

-- 3. Roles table
CREATE TABLE schema_admin.roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    permissions JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed initial roles
INSERT INTO schema_admin.roles (id, name, description, permissions) VALUES
    (1, 'admin', '普通管理员', '["audit_posts","manage_posts","manage_boards","manage_sensitive_words","view_users","view_invite_codes"]'),
    (2, 'platform_admin', '中台管理员', '["audit_posts","manage_posts","manage_boards","manage_sensitive_words","manage_users","manage_levels","manage_invite_codes","manage_config","manage_roles","view_stats"]')
ON CONFLICT (id) DO NOTHING;

-- Reset sequence after seeding
SELECT setval('schema_admin.roles_id_seq', (SELECT MAX(id) FROM schema_admin.roles));

-- 4. User roles mapping table
CREATE TABLE schema_admin.user_roles (
    user_id BIGINT NOT NULL,
    role_id INT NOT NULL REFERENCES schema_admin.roles(id),
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    assigned_by BIGINT DEFAULT 0,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_role ON schema_admin.user_roles(role_id);
