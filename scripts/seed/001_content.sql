-- Sample forum content for first launch (idempotent).
-- Requires: services migrated + demo users from user migration 004.

UPDATE schema_auth.users SET nickname = '林夏', bio = '关注大模型应用、AI 项目共创和技术答疑。'
WHERE username = 'demo_student';
UPDATE schema_auth.users SET nickname = '周沐', bio = '负责活动公告、内容审核和社区秩序维护。'
WHERE username = 'demo_admin';
UPDATE schema_auth.users SET nickname = '陈拓', bio = '负责系统配置、监管规则和核心数据看板。'
WHERE username = 'demo_platform_admin';

INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, is_pinned, is_featured, like_count, comment_count, created_at)
SELECT
    '本学期 AI 学习路线怎么排最稳？',
    '如果一学期只抓三个方向，我建议按 Python 基础、机器学习核心概念、一个能展示的 AI 应用项目来安排。',
    u.id,
    b.id,
    'published',
    false,
    true,
    14,
    2,
    NOW() - INTERVAL '4 days'
FROM schema_auth.users u
JOIN schema_forum.boards b ON b.slug = 'ai-learning'
WHERE u.username = 'demo_student'
  AND NOT EXISTS (SELECT 1 FROM schema_forum.posts LIMIT 1);

INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, is_featured, like_count, comment_count, created_at)
SELECT
    '本周五 AI Mini Hack Night 报名开启',
    '周五晚 19:00 在实验楼 302 举办 Mini Hack Night，欢迎带上自己的 idea 来现场组队。',
    u.id,
    b.id,
    'published',
    false,
    9,
    1,
    NOW() - INTERVAL '3 days'
FROM schema_auth.users u
JOIN schema_forum.boards b ON b.slug = 'announcements'
WHERE u.username = 'demo_admin'
  AND (SELECT COUNT(*) FROM schema_forum.posts) < 2;

INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, like_count, comment_count, created_at)
SELECT
    'RAG 项目里向量库先选 FAISS 还是 PGVector？',
    '数据量还不大，但后面可能会扩展到多个课程知识库。MVP 阶段先用哪个更合适？',
    u.id,
    b.id,
    'published',
    17,
    2,
    NOW() - INTERVAL '2 days'
FROM schema_auth.users u
JOIN schema_forum.boards b ON b.slug = 'tech-help'
WHERE u.username = 'demo_student'
  AND (SELECT COUNT(*) FROM schema_forum.posts) < 3;

INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, like_count, created_at)
SELECT
    '求一个校外技术交流群，最好能直接引流进群',
    '想找一个适合继续交流的群，如果有现成的技术引流入口也欢迎推荐。',
    u.id,
    b.id,
    'pending_review',
    0,
    NOW() - INTERVAL '1 day'
FROM schema_auth.users u
JOIN schema_forum.boards b ON b.slug = 'announcements'
WHERE u.username = 'demo_student'
  AND (SELECT COUNT(*) FROM schema_forum.posts) < 4;

INSERT INTO schema_forum.comments (post_id, author_id, content, created_at)
SELECT p.id, u.id, '这条路线很适合协会内部推广，后面可以顺手补一份工具安装清单。', NOW() - INTERVAL '3 days'
FROM schema_forum.posts p
JOIN schema_auth.users u ON u.username = 'demo_admin'
WHERE p.title LIKE '本学期 AI 学习路线%'
  AND NOT EXISTS (SELECT 1 FROM schema_forum.comments LIMIT 1);

INSERT INTO schema_forum.comments (post_id, author_id, content, created_at)
SELECT p.id, u.id, '如果你要做展示，最好在第三周就开始搭第一个 demo。', NOW() - INTERVAL '3 days'
FROM schema_forum.posts p
JOIN schema_auth.users u ON u.username = 'demo_platform_admin'
WHERE p.title LIKE '本学期 AI 学习路线%'
  AND (SELECT COUNT(*) FROM schema_forum.comments) < 2;

INSERT INTO schema_forum.comments (post_id, author_id, content, created_at)
SELECT p.id, u.id, '报名之后会提前公布组队主题吗？', NOW() - INTERVAL '2 days'
FROM schema_forum.posts p
JOIN schema_auth.users u ON u.username = 'demo_student'
WHERE p.title LIKE '%Hack Night%'
  AND (SELECT COUNT(*) FROM schema_forum.comments) < 3;

INSERT INTO schema_forum.comments (post_id, author_id, content, created_at)
SELECT p.id, u.id, '如果你现在主要目标是验证流程，先用 FAISS 本地跑通会更省力。', NOW() - INTERVAL '2 days'
FROM schema_forum.posts p
JOIN schema_auth.users u ON u.username = 'demo_platform_admin'
WHERE p.title LIKE '%FAISS%'
  AND (SELECT COUNT(*) FROM schema_forum.comments) < 4;

INSERT INTO schema_forum.comments (post_id, author_id, content, created_at)
SELECT p.id, u.id, '后续如果考虑多人协作和数据持久化，再切到 PGVector 会更顺一些。', NOW() - INTERVAL '2 days'
FROM schema_forum.posts p
JOIN schema_auth.users u ON u.username = 'demo_admin'
WHERE p.title LIKE '%FAISS%'
  AND (SELECT COUNT(*) FROM schema_forum.comments) < 5;
