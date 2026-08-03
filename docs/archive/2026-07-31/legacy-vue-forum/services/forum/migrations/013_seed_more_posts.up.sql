-- Migration 013: Seed more posts for content richness
-- Adds 20 diverse AI-themed posts across boards

DO $\$
DECLARE
    v_student_id BIGINT;
    v_admin_id BIGINT;
    v_platform_id BIGINT;
    v_board_ai BIGINT;
    v_board_announce BIGINT;
    v_board_tech BIGINT;
    v_board_study BIGINT;
    v_board_training BIGINT;
    v_board_notice BIGINT;
    v_board_club BIGINT;
BEGIN
    -- Get user IDs
    SELECT id INTO v_student_id FROM schema_auth.users WHERE role = 'student' LIMIT 1;
    SELECT id INTO v_admin_id FROM schema_auth.users WHERE role = 'admin' LIMIT 1;
    SELECT id INTO v_platform_id FROM schema_auth.users WHERE role = 'platform_admin' LIMIT 1;

    -- Get board IDs
    SELECT id INTO v_board_ai FROM schema_forum.boards WHERE slug = 'ai-learning' LIMIT 1;
    SELECT id INTO v_board_announce FROM schema_forum.boards WHERE slug = 'announcements' LIMIT 1;
    SELECT id INTO v_board_tech FROM schema_forum.boards WHERE slug = 'tech-help' LIMIT 1;
    SELECT id INTO v_board_study FROM schema_forum.boards WHERE slug = 'study' LIMIT 1;
    SELECT id INTO v_board_training FROM schema_forum.boards WHERE slug = 'training' LIMIT 1;
    SELECT id INTO v_board_notice FROM schema_forum.boards WHERE slug = 'notice' LIMIT 1;
    SELECT id INTO v_board_club FROM schema_forum.boards WHERE slug = 'club' LIMIT 1;

    -- Skip if already seeded (check for existing posts beyond the initial 4)
    IF (SELECT COUNT(*) FROM schema_forum.posts) >= 10 THEN
        RAISE NOTICE 'Posts already seeded, skipping.';
        RETURN;
    END IF;

    -- AI Learning Zone posts
    INSERT INTO schema_forum.posts (title, content, author_id, board_id, status, like_count, comment_count, created_at, updated_at)
    VALUES
    ('用Stable Diffusion做校园文创设计，从Prompt工程到成品',
     '最近协会在筹备迎新文创，我尝试用SD来生成设计稿。分享几个心得：\n\n1. 提示词一定要包含材质描述（matte finish, holographic foil等），会极大影响最终质感\n2. 用ControlNet固定构图，否则每次生成的布局都不一样\n3. 批量出图后用CLIP score做初筛，效率提升明显\n\n附上我做的几版帆布袋和徽章设计，欢迎大家提意见！',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_ai, 1), 'published', 28, 6, NOW() - INTERVAL '1 day', NOW()),

    ('LLM微调避坑指南：我花了3天踩的5个坑',
     '最近在微调一个7B模型做校园知识问答，记录几个踩过的坑：\n\n1. LoRA rank设太大反而过拟合，建议从8开始试\n2. 学习率2e-4对大多数7B模型都偏大了，1e-4更稳\n3. 数据配比很重要——单轮对话太多会导致模型忘记多轮能力\n4. 验证集一定要和训练集同分布，否则eval loss没有参考意义\n5. 记得保存checkpoint，不然训了8小时崩了会想哭\n\n大家还有遇到什么坑？',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_ai, 1), 'published', 35, 8, NOW() - INTERVAL '2 days', NOW()),

    ('推荐一个超好用的论文阅读工具：用RAG搭建个人知识库',
     '期末赶论文的时候搭建了一个论文知识库，核心思路：\n\n- 用PyMuPDF提取论文全文\n- 按章节切分chunk，保留标题上下文\n- 用BGE-M3做embedding\n- Milvus做向量存储\n- 前端用Streamlit搭了个简单的问答界面\n\n效果意外地好，跨论文提问时能自动关联相关段落。代码已开源在GitHub，欢迎star！',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_ai, 1), 'published', 42, 12, NOW() - INTERVAL '3 days', NOW()),

    ('大模型时代，还有必要学传统NLP吗？',
     '最近和老师讨论了这个话题。我的看法是：\n\n传统NLP（分词、NER、句法分析）的价值不在于直接用，而在于理解语言的结构性。做RAG评估时需要计算准确率和召回率，做数据清洗时需要分词和去停用词——这些基本功绕不开。\n\n但确实不需要花太多时间在传统方法上，重点是理解"为什么"而不是"怎么做"。大家怎么看？',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_ai, 1), 'published', 19, 5, NOW() - INTERVAL '4 days', NOW()),

    -- Tech Q&A zone
    ('Ubuntu 24.04装CUDA死活不成功，求救！',
     '环境：RTX 4060 Laptop + Ubuntu 24.04\n\n装了nvidia-driver-550，nvidia-smi正常，但装cuda-toolkit 12.4后nvcc -V报错说找不到。\n\n试过的方法：\n- 用runfile装\n- 用deb(network)装\n- 手动加PATH\n\n都不行。有没有24.04上装成功的大佬？求个详细步骤😭',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_tech, 3), 'published', 8, 15, NOW() - INTERVAL '1 day', NOW()),

    ('求推荐适合学生党的高性价比GPU云平台',
     '显卡太贵买不起，一直在用云端GPU。目前用过：\n- AutoDL：便宜但高峰期要排队\n- Colab Pro：方便但显存小\n- 矩池云：没用过，有人说说吗？\n\n需求是跑7B模型的微调，24G显存够用。预算每月200以内。',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_tech, 3), 'published', 22, 10, NOW() - INTERVAL '2 days', NOW()),

    -- Study zone
    ('期末复习资料共享：模式识别重点整理',
     '整理了这学期模式识别的重点内容：\n\n1. 贝叶斯决策理论（最小错误率、最小风险）\n2. 线性判别函数（Fisher、感知机）\n3. SVM（硬间隔、软间隔、核方法）\n4. 集成学习（Bagging、Boosting）\n5. 聚类（K-means、GMM、DBSCAN）\n\n附上我的复习笔记和往年试题。祝大家考试顺利！',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_study, 4), 'published', 56, 18, NOW() - INTERVAL '3 days', NOW()),

    ('大三下学期才开始学AI，来得及吗？',
     '目前基础：Python会写，高数线代概统都学过，但机器学习零基础。\n\n计划下学期选课+自学，目标是毕业前能做出一个有模有样的AI项目。\n\n想问问过来人，这条路规划可行吗？重点是应该先学理论还是先动手做项目？',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_study, 4), 'published', 31, 9, NOW() - INTERVAL '4 days', NOW()),

    -- Training zone
    ('警务大数据分析实战：用Python做犯罪热点地图',
     '实训课上的一个小项目，分享一下：\n\n- 数据来源：公开的犯罪统计数据\n- 技术栈：Pandas清洗 + Folium可视化 + Scikit-learn做聚类\n- 结果：成功识别出了几个犯罪高发区域，和实际警务部署高度吻合\n\n代码和报告都在附件里。老师说可以进一步结合时间序列做预测，有兴趣的同学一起做吗？',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_training, 5), 'published', 15, 4, NOW() - INTERVAL '1 day', NOW()),

    ('无人机在警务巡逻中的应用前景',
     '最近参加了一个讲座，关于无人机在警务中的应用。分享几个有意思的点：\n\n1. 热成像在夜间搜救中的应用\n2. AI目标检测用于交通违章自动识别\n3. 5G+无人机实时回传方案\n\n感觉这个方向未来几年会有很大发展，有没有同好想组个学习小组？',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_training, 5), 'published', 18, 3, NOW() - INTERVAL '3 days', NOW()),

    -- Club zone
    ('AI协会招新了！欢迎所有对AI感兴趣的同学',
     '无论你是AI大佬还是零基础小白，AI协会都欢迎你！\n\n我们提供：\n- 每周技术分享会\n- 项目实战机会（有大佬带飞）\n- 参加各类AI竞赛\n- 与企业和实验室合作的机会\n\n本周五晚7点在教学楼201有招新宣讲，来就有小礼品！\n\n报名链接在公告区~',
     COALESCE(v_admin_id, v_student_id), COALESCE(v_board_club, 7), 'published', 67, 22, NOW() - INTERVAL '5 days', NOW()),

    ('协会项目组招募：校园AI助手二期开发',
     '一期已经完成了基础问答功能，现在准备做二期迭代：\n\n计划新增功能：\n- 课表查询\n- 空教室查询\n- 校园导航\n- 失物招领\n\n需要前端（Vue）、后端（Go/Python）、算法（NLP/RAG）的同学。\n\n技术栈和一期保持一致，有完整文档和代码可参考。报名私信我！',
     COALESCE(v_admin_id, v_student_id), COALESCE(v_board_club, 7), 'published', 45, 14, NOW() - INTERVAL '6 days', NOW()),

    -- Notice zone
    ('关于举办2026年AI创新大赛的通知',
     '各学院、各位同学：\n\n为激发同学们对人工智能的兴趣，学校决定举办2026年AI创新大赛。\n\n比赛主题：AI赋能校园生活\n参赛形式：1-4人组队\n作品形式：Demo + 答辩PPT\n\n时间安排：\n- 报名截止：8月15日\n- 初赛（线上评审）：8月30日\n- 决赛（现场答辩）：9月15日\n\n奖项设置：\n- 一等奖1名：5000元\n- 二等奖2名：3000元\n- 三等奖5名：1000元\n\n详细规则见附件。',
     COALESCE(v_admin_id, v_student_id), COALESCE(v_board_notice, 6), 'published', 89, 25, NOW() - INTERVAL '2 days', NOW()),

    -- Announcements zone
    ('本周活动预告：大模型应用开发Workshop',
     '本周六下午2点开始，带大家从零搭建一个大模型应用！\n\n内容：\n1. LangChain入门（14:00-15:00）\n2. 搭建RAG问答系统（15:00-16:30）\n3. 部署上线（16:30-17:00）\n4. 自由交流 + 答疑（17:00-18:00）\n\n地点：实验楼302\n准备：自带笔记本电脑，提前装好Python 3.10+\n\n无需报名，直接来就好！',
     COALESCE(v_admin_id, v_student_id), COALESCE(v_board_announce, 2), 'published', 33, 7, NOW() - INTERVAL '1 day', NOW()),

    -- More AI Learning posts
    ('万字长文：从零理解Attention机制',
     '写这篇文章的初衷是因为我发现很多同学对Attention的理解停留在"加权求和"这个层面，但面试和实际应用中需要更深的理解。\n\n本文从这几个角度展开：\n1. 为什么需要Attention——RNN的瓶颈\n2. Scaled Dot-Product Attention的数学直觉\n3. 多头注意力的并行化优势\n4. 复杂度分析和优化（Flash Attention简介）\n5. 实际代码实现（PyTorch）\n\n写了快一周，希望帮到大家！',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_ai, 1), 'published', 78, 20, NOW() - INTERVAL '5 days', NOW()),

    ('开源了一个AI配音工具，支持情感控制',
     '基于GPT-SoVITS做的一个配音工具，特点是：\n\n- 支持上传参考音频来克隆音色\n- 可以在文本中插入情感标签来控制语气\n- 支持批量生成\n- Web界面，部署简单\n\n本来是给协会宣传片做配音用的，后来觉得通用性很强就开源了。GitHub链接见评论区置顶。',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_ai, 1), 'published', 53, 16, NOW() - INTERVAL '6 days', NOW()),

    -- More Q&A posts
    ('请教：LangChain和LlamaIndex到底怎么选？',
     '两个框架都看了下文档，感觉功能重叠很多。\n\n我的场景：\n- 需要做文档问答（课程资料）\n- 数据量不大，几千页PDF\n- 需要支持流式输出\n- 后续可能会加Agent功能\n\n目前倾向LlamaIndex因为它的索引结构更清晰，但LangChain的生态好像更大。大家实际用过的给点建议？',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_tech, 3), 'published', 16, 11, NOW() - INTERVAL '3 days', NOW()),

    -- More Study posts
    ('分享我的AI学习路线（从入门到能找实习）',
     '从大二开始学AI，到现在大三拿到了AI实习offer，分享一下我的路线：\n\n第一阶段（1-2月）：Python + 数学基础\n- 重点：NumPy、Pandas、线代、概率论\n\n第二阶段（2-3月）：机器学习基础\n- 吴恩达课程 + 《统计学习方法》\n- 动手实现几个经典算法\n\n第三阶段（3-4月）：深度学习\n- 李沐《动手学深度学习》\n- 用PyTorch复现论文\n\n第四阶段（2-3月）：项目实战\n- Kaggle打比赛\n- GitHub上找开源项目贡献\n- 做自己的项目（简历上要有2-3个像样的）\n\n关键：不要只看不写，每个阶段都要有代码产出！',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_study, 4), 'published', 94, 28, NOW() - INTERVAL '7 days', NOW()),

    -- More Training posts
    ('智能安防系统设计：人脸识别+异常行为检测',
     '实训项目记录：设计一个校园智能安防系统。\n\n技术方案：\n- 人脸识别：ArcFace + 本地人脸库\n- 异常行为检测：基于YOLOv8的行为分类\n- 后端：Go + PostgreSQL\n- 前端：Vue + 视频流播放\n\n遇到的问题和解决方案：\n1. 夜间识别率低→加红外补光+图像增强预处理\n2. 多人同时检测延迟高→模型量化+TensorRT加速\n\n目前准确率白天98%，夜间91%，还在优化中。',
     COALESCE(v_student_id, v_admin_id), COALESCE(v_board_training, 5), 'published', 21, 6, NOW() - INTERVAL '4 days', NOW()),

    -- More Club posts
    ('AI协会年度回顾：这一年我们做了什么',
     '转眼一年过去了，回顾一下协会的成绩：\n\n📊 数据：\n- 举办了18场技术分享\n- 完成了5个开源项目\n- 成员增长了300%\n- 在全国AI竞赛中获得二等奖\n\n🎯 下一年目标：\n- 建立校内AI学习资源库\n- 组织跨校AI交流活动\n- 发布协会自己的AI产品\n\n感谢每一位成员的付出！新的一年继续加油💪',
     COALESCE(v_admin_id, v_student_id), COALESCE(v_board_club, 7), 'published', 112, 30, NOW() - INTERVAL '8 days', NOW());

    RAISE NOTICE 'Seed posts inserted successfully.';
END
\$\$;
