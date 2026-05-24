-- Demo accounts for MVP presentation (password: demo123456)
INSERT INTO schema_auth.users (username, password_hash, nickname, bio, avatar, level, status)
VALUES
    ('demo_student', '$2a$10$le.KP1QbQN9SkS5VajlApeUojQ/i7VGFDB0p1/fPcpE3qlANEj.VW', '演示学生', '学院 AI 社团成员', '', 1, 'active'),
    ('demo_admin', '$2a$10$le.KP1QbQN9SkS5VajlApeUojQ/i7VGFDB0p1/fPcpE3qlANEj.VW', '协会管理员', '负责内容审核与运营', '', 2, 'active'),
    ('demo_platform_admin', '$2a$10$le.KP1QbQN9SkS5VajlApeUojQ/i7VGFDB0p1/fPcpE3qlANEj.VW', '中台管理员', '系统配置与数据监管', '', 2, 'active')
ON CONFLICT (username) DO NOTHING;

INSERT INTO schema_auth.invite_codes (code, status, created_by)
VALUES ('DEMO2026', 'unused', 0)
ON CONFLICT (code) DO NOTHING;
