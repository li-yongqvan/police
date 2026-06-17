-- Batch demo accounts for parallel UAT (password equals username)
INSERT INTO schema_auth.users (username, password_hash, nickname, bio, avatar, level, status)
VALUES
    ('demo01', '$2a$10$YT1IlEw.HcIcbuFJZ6pE1uq8jYcRK60L2.SEcb69mCqzdOLCpsWv2', '演示学生01', '并行测试学生账号', '', 1, 'active'),
    ('demo02', '$2a$10$msIkTlaliuQdfHGqx/AeneKVHGP0CBAuuCkpVEJigfGVAyLqfzkJq', '演示学生02', '并行测试学生账号', '', 1, 'active'),
    ('demo03', '$2a$10$SPNb0bR.S7YfNDwBZ1OOhuPiaS65tPSGoJfQPL5ChypigusX11NY2', '演示学生03', '并行测试学生账号', '', 1, 'active'),
    ('demo04', '$2a$10$0zNpt3dD0NFPzluizh94oOBNv6R4CALuAGznoZ2xQtYj0wc/fd9Ui', '演示学生04', '并行测试学生账号', '', 1, 'active'),
    ('demo05', '$2a$10$IZHJb86G9rmAUapkqNFiF.sh3FI8zWsUaD4JLPAdFPqjZbkCrChkG', '演示学生05', '并行测试学生账号', '', 1, 'active'),
    ('demo06', '$2a$10$/XXaocudLGTfwgG/puTXretYobb6TPvYo0j74zsnQFe/J22kAfFfm', '演示学生06', '并行测试学生账号', '', 1, 'active'),
    ('admin01', '$2a$10$ardgKXHj2MXpGH/VJUzy7e16N4RbhfWWxF6lr0O1P1ztP.tVDTeK.', '协会管理员01', '并行测试协会管理账号', '', 2, 'active'),
    ('admin02', '$2a$10$pjPcOYUgi90wA.nLgkTpl.eaQr2Uyzwy7MqsKuruePoUdszdoAeqC', '协会管理员02', '并行测试协会管理账号', '', 2, 'active'),
    ('admin03', '$2a$10$s0JN/mwMVHpGGbZu6neet.DxpegrxK8aaZv1UbobSBlz.6OmgEDgO', '协会管理员03', '并行测试协会管理账号', '', 2, 'active'),
    ('admin04', '$2a$10$a30S7bXmFB2iX8bmG7PDEejemr7l8e302CAr8EEf5HxIEYlhJywpC', '协会管理员04', '并行测试协会管理账号', '', 2, 'active'),
    ('admin05', '$2a$10$jvp1/ZyW0mgdZtL6vYQGp.KX/wYiGu5o6CkyrWBd3PE4aaa.l8OSq', '协会管理员05', '并行测试协会管理账号', '', 2, 'active'),
    ('admin06', '$2a$10$H2aO9UKpnumFskwf/rr93OxvD4JhokT/yXnkRjtL.lwl7z7CGYk7G', '协会管理员06', '并行测试协会管理账号', '', 2, 'active'),
    ('plat01', '$2a$10$xdL9cu9ZZ9nWX0PRvMBze.qfVs3Ehrx/1GJg7tKv7/RplpuKLkZpS', '中台管理员01', '并行测试中台管理账号', '', 2, 'active'),
    ('plat02', '$2a$10$n48ejyKF4N2yDfCwqIPaZu5Qy7UXYBXfhG0YN9B.EVxEhDyCDZYoS', '中台管理员02', '并行测试中台管理账号', '', 2, 'active'),
    ('plat03', '$2a$10$XJ0aagNdGCHIXRCXVxTloODi5xnzj2YSPAXgLa/p6KnwaPZ0Zb8N6', '中台管理员03', '并行测试中台管理账号', '', 2, 'active'),
    ('plat04', '$2a$10$yrSePh9G7FR.Bm2Mb308sukmJia6pkEkIIi0vMWwsz9o1BKaz75uW', '中台管理员04', '并行测试中台管理账号', '', 2, 'active'),
    ('plat05', '$2a$10$vpw1dRzhAVqvbL.DG/uo8uOz4yN841RmTF.08excCg01k.KJDOrBK', '中台管理员05', '并行测试中台管理账号', '', 2, 'active'),
    ('plat06', '$2a$10$K7YPksssUzM1qzX14irPdekbltwRlVNDXyUz0At9u0eN1/UAXiqlu', '中台管理员06', '并行测试中台管理账号', '', 2, 'active')
ON CONFLICT (username) DO NOTHING;
