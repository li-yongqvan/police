INSERT INTO schema_forum.boards (name, slug, description, sort_order, enabled)
VALUES
    ('学业研讨区', 'study', '课程学习、备考与学术讨论', 10, true),
    ('警务实训区', 'training', '技能训练、实习心得与实务交流', 11, true),
    ('校园公告区', 'notice', '学院与协会通知、活动信息', 12, true),
    ('社团风采区', 'club', '社团活动、文化建设与风采展示', 13, true)
ON CONFLICT (slug) DO UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    sort_order = EXCLUDED.sort_order,
    enabled = EXCLUDED.enabled;
