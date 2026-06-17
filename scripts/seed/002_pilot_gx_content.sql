-- Pilot seed: boards aligned with GX nav + sample posts (idempotent).

INSERT INTO schema_forum.boards (name, slug, description, sort_order, enabled)
VALUES
    ('学业研讨区', 'study', '课程学习、备考与学术讨论', 10, true),
    ('警务实训区', 'training', '技能训练、实习心得与实务交流', 11, true),
    ('校园公告区', 'notice', '学院与协会通知、活动信息', 12, true),
    ('社团风采区', 'club', '社团活动、文化建设与风采展示', 13, true)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, created_at)
SELECT
    '【置顶】校内试运行须知',
    '欢迎使用 AI智联平台。注册需邀请码；发帖请文明守法。遇到问题可联系辅导员或协会管理员。',
    u.id, b.id, 'published', true, true, 3, 0, NOW() - INTERVAL '1 day'
FROM schema_auth.users u
JOIN schema_forum.boards b ON b.slug = 'notice'
WHERE u.username = 'demo_admin'
  AND NOT EXISTS (
    SELECT 1 FROM schema_forum.posts p WHERE p.title = '【置顶】校内试运行须知'
  );

INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, like_count, created_at)
SELECT v.title, v.content, u.id, b.id, 'published', 2, NOW() - INTERVAL '2 days'
FROM schema_auth.users u
CROSS JOIN (VALUES
    ('刑法总论期末复习要点整理', 'study', '整理了课堂笔记与历年题型，欢迎补充。'),
    ('警务礼仪实训心得', 'training', '本周礼仪课要点：着装、敬礼、用语规范。'),
    ('人工智能社团纳新说明', 'club', '面向对 AI 应用感兴趣的同学，详情见本周公告。')
) AS v(title, slug, content)
JOIN schema_forum.boards b ON b.slug = v.slug
WHERE u.username = 'demo_student'
  AND NOT EXISTS (SELECT 1 FROM schema_forum.posts p WHERE p.title = v.title);
