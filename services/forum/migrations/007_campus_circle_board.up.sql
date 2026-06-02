INSERT INTO schema_forum.boards (name, slug, description, sort_order)
VALUES (
    '校园圈',
    'campus-circle',
    '校园日常、轻松闲聊，分享你在警院的小趣事',
    10
)
ON CONFLICT (slug) DO NOTHING;
