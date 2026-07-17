const fs = require('fs');
const path = require('path');
const {
  Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell,
  Header, Footer, AlignmentType, LevelFormat,
  HeadingLevel, BorderStyle, WidthType, ShadingType,
  PageNumber, PageBreak
} = require('docx');

const BLUE_DARK = "1A5276";
const BLUE_MID = "2980B9";
const BLUE_LIGHT = "D6EAF8";
const GREEN = "27AE60";
const ORANGE = "E67E22";
const WHITE = "FFFFFF";
const GRAY_BG = "F2F3F4";

const noBorder = { style: BorderStyle.NONE, size: 0 };
const noBorders = { top: noBorder, bottom: noBorder, left: noBorder, right: noBorder };
const cm = { top: 80, bottom: 80, left: 120, right: 120 };

const sec = (t) => new Paragraph({ heading: HeadingLevel.HEADING_1, spacing: { before: 360, after: 200 },
  children: [new TextRun({ text: t, font: "Microsoft YaHei", bold: true, size: 36, color: BLUE_DARK })] });
const sub = (t) => new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 280, after: 120 },
  children: [new TextRun({ text: t, font: "Microsoft YaHei", bold: true, size: 28, color: BLUE_MID })] });
const body = (t, o = {}) => new Paragraph({ spacing: { before: 60, after: 60 },
  children: [new TextRun({ text: t, font: "Microsoft YaHei", size: 21, color: o.color || "333333", bold: o.bold || false, italics: o.italics || false })] });
const check = (t) => new Paragraph({ spacing: { before: 40, after: 40 },
  children: [new TextRun({ text: "☐  ", font: "Microsoft YaHei", size: 22, color: BLUE_MID }), new TextRun({ text: t, font: "Microsoft YaHei", size: 21, color: "333333" })] });
const step = (n, t) => new Paragraph({ spacing: { before: 80, after: 40 },
  children: [new TextRun({ text: n + ". ", font: "Microsoft YaHei", size: 21, color: BLUE_MID, bold: true }), new TextRun({ text: t, font: "Microsoft YaHei", size: 21, color: "333333" })] });
const see = (t) => new Paragraph({ spacing: { before: 40, after: 20 }, indent: { left: 420 },
  children: [new TextRun({ text: "▸ 应该看到：", font: "Microsoft YaHei", size: 20, color: GREEN, bold: true }), new TextRun({ text: t, font: "Microsoft YaHei", size: 20, color: "555555" })] });
const tip = (t) => new Paragraph({ spacing: { before: 60, after: 60 }, indent: { left: 360 },
  children: [new TextRun({ text: "💡 ", font: "Microsoft YaHei", size: 20 }), new TextRun({ text: t, font: "Microsoft YaHei", size: 20, color: "888888", italics: true })] });
const warn = (t) => new Paragraph({ spacing: { before: 40, after: 40 }, indent: { left: 420 },
  children: [new TextRun({ text: "⚠ ", font: "Microsoft YaHei", size: 20, color: ORANGE, bold: true }), new TextRun({ text: t, font: "Microsoft YaHei", size: 20, color: ORANGE })] });
const div = () => new Paragraph({ spacing: { before: 120, after: 120 }, border: { bottom: { style: BorderStyle.SINGLE, size: 6, color: BLUE_LIGHT, space: 1 } }, children: [] });
const pb = () => new Paragraph({ children: [new PageBreak()] });

function labeled(k, v) {
  return new Paragraph({ spacing: { before: 40, after: 40 }, indent: { left: 420 },
    children: [new TextRun({ text: k + "：", font: "Microsoft YaHei", size: 20, bold: true, color: "555555" }), new TextRun({ text: "  " + v, font: "Microsoft YaHei", size: 20, color: "333333" })] });
}

function moduleCard(m) {
  return new Table({
    width: { size: 9026, type: WidthType.DXA },
    columnWidths: [600, 3200, 1200, 3626],
    rows: [new TableRow({ children: [
      new TableCell({ borders: noBorders, width: { size: 600, type: WidthType.DXA }, shading: { fill: m.col, type: ShadingType.CLEAR }, margins: { top: 120, bottom: 120, left: 100, right: 100 },
        children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: m.mod, font: "Microsoft YaHei", size: 24, bold: true, color: WHITE })] })] }),
      new TableCell({ borders: noBorders, width: { size: 3200, type: WidthType.DXA }, shading: { fill: WHITE, type: ShadingType.CLEAR }, margins: { top: 120, bottom: 120, left: 160, right: 80 },
        children: [
          new Paragraph({ spacing: { before: 0, after: 20 }, children: [new TextRun({ text: m.title, font: "Microsoft YaHei", size: 22, bold: true, color: "333333" })] }),
          new Paragraph({ spacing: { before: 0, after: 0 }, children: [new TextRun({ text: m.desc, font: "Microsoft YaHei", size: 18, color: "777777" })] })
        ] }),
      new TableCell({ borders: noBorders, width: { size: 1200, type: WidthType.DXA }, shading: { fill: WHITE, type: ShadingType.CLEAR }, margins: { top: 120, bottom: 120, left: 80, right: 80 },
        children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: m.time, font: "Microsoft YaHei", size: 18, bold: true, color: m.col })] })] }),
      new TableCell({ borders: noBorders, width: { size: 3626, type: WidthType.DXA }, shading: { fill: WHITE, type: ShadingType.CLEAR }, margins: { top: 120, bottom: 120, left: 80, right: 160 },
        children: [new Paragraph({ alignment: AlignmentType.RIGHT, children: [new TextRun({ text: m.who, font: "Microsoft YaHei", size: 18, color: "AAAAAA", italics: true })] })] }),
    ]})]
  });
}

function greenBox(lines) {
  return new Table({
    width: { size: 9026, type: WidthType.DXA }, columnWidths: [9026],
    rows: [new TableRow({ children: [new TableCell({ borders: noBorders, width: { size: 9026, type: WidthType.DXA }, shading: { fill: "E8F8F5", type: ShadingType.CLEAR }, margins: { top: 150, bottom: 150, left: 200, right: 200 },
      children: lines
    })] })]
  });
}

// Commented out — not currently used
// function blueBox(lines) {
//   return new Table({
//     width: { size: 9026, type: WidthType.DXA }, columnWidths: [9026],
//     rows: [new TableRow({ children: [new TableCell({ borders: noBorders, width: { size: 9026, type: WidthType.DXA }, shading: { fill: BLUE_LIGHT, type: ShadingType.CLEAR }, margins: { top: 150, bottom: 150, left: 200, right: 200 },
//       children: lines
//     })] })]
//   });
// }

const MODULES = [
  { col: BLUE_MID, mod: "A", title: "注册与登录", time: "15 分钟", desc: "注册新号、错误密码、退出登录", who: "想测账号相关" },
  { col: "1ABC9C", mod: "B", title: "浏览与搜索", time: "10 分钟", desc: "首页、板块、帖子详情、搜索", who: "只想随便点点" },
  { col: ORANGE, mod: "C", title: "发帖与互动", time: "25 分钟", desc: "发帖、编辑、评论、点赞、收藏、举报、删除", who: "想测核心功能" },
  { col: "8E44AD", mod: "D", title: "个人资料", time: "10 分钟", desc: "改昵称、院系、简介", who: "想测个人中心" },
  { col: "D35400", mod: "E", title: "消息通知", time: "15 分钟", desc: "两人互评看通知", who: "有搭子一起测" },
  { col: "7F8C8D", mod: "F", title: "权限检查", time: "2 分钟", desc: "学生能不能进后台", who: "只想交差" },
];

const doc = new Document({
  styles: {
    default: { document: { run: { font: "Microsoft YaHei", size: 21 } } },
    paragraphStyles: [
      { id: "Heading1", name: "Heading 1", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 36, bold: true, font: "Microsoft YaHei", color: BLUE_DARK },
        paragraph: { spacing: { before: 360, after: 200 }, outlineLevel: 0 } },
      { id: "Heading2", name: "Heading 2", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 28, bold: true, font: "Microsoft YaHei", color: BLUE_MID },
        paragraph: { spacing: { before: 280, after: 120 }, outlineLevel: 1 } },
    ]
  },
  sections: [
    // ===== COVER PAGE =====
    {
      properties: { page: { size: { width: 11906, height: 16838 }, margin: { top: 0, right: 0, bottom: 0, left: 0 } } },
      children: [
        new Table({ width: { size: 11906, type: WidthType.DXA }, columnWidths: [11906], rows: [new TableRow({ children: [
          new TableCell({ borders: noBorders, width: { size: 11906, type: WidthType.DXA }, shading: { fill: BLUE_DARK, type: ShadingType.CLEAR }, margins: { top: 600, bottom: 600, left: 800, right: 800 },
            children: [
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 300, after: 100 }, children: [new TextRun({ text: "AI 智联论坛", font: "Microsoft YaHei", size: 52, bold: true, color: WHITE })] }),
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 80, after: 300 }, children: [new TextRun({ text: "学生端测试指南", font: "Microsoft YaHei", size: 36, color: "D6EAF8" })] }),
            ] })
        ] })] }),
        new Table({ width: { size: 11906, type: WidthType.DXA }, columnWidths: [11906], rows: [new TableRow({ children: [
          new TableCell({ borders: noBorders, width: { size: 11906, type: WidthType.DXA }, shading: { fill: WHITE, type: ShadingType.CLEAR }, margins: { top: 400, bottom: 200, left: 1200, right: 1200 },
            children: [
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 200, after: 200 }, children: [new TextRun({ text: "不想写测试用例，但又得来测一下？", font: "Microsoft YaHei", size: 26, color: "555555" })] }),
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 80, after: 80 }, children: [new TextRun({ text: "快速开始 5 分钟 + 挑个感兴趣的模块", font: "Microsoft YaHei", size: 30, bold: true, color: BLUE_MID })] }),
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 80, after: 200 }, children: [new TextRun({ text: "做完一个就能交差", font: "Microsoft YaHei", size: 22, color: "888888" })] }),
            ] })
        ] })] }),
        new Table({ width: { size: 11906, type: WidthType.DXA }, columnWidths: [11906], rows: [new TableRow({ children: [
          new TableCell({ borders: noBorders, width: { size: 11906, type: WidthType.DXA }, shading: { fill: GRAY_BG, type: ShadingType.CLEAR }, margins: { top: 300, bottom: 300, left: 1200, right: 1200 },
            children: [
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 60, after: 60 }, children: [new TextRun({ text: "论坛地址：http://127.0.0.1:8091     |     密码统一：123456     |     邀请码：DEMO2026", font: "Microsoft YaHei", size: 20, color: "666666" })] }),
              new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 40, after: 40 }, children: [new TextRun({ text: "测试人：________     日期：________", font: "Microsoft YaHei", size: 20, color: "888888" })] }),
            ] })
        ] })] }),
      ]
    },

    // ===== MAIN CONTENT =====
    {
      properties: { page: { size: { width: 11906, height: 16838 }, margin: { top: 1200, right: 1200, bottom: 1200, left: 1200 } } },
      headers: { default: new Header({ children: [new Paragraph({ alignment: AlignmentType.RIGHT, border: { bottom: { style: BorderStyle.SINGLE, size: 4, color: BLUE_LIGHT, space: 4 } },
        children: [new TextRun({ text: "AI 智联论坛 · 学生端测试指南", font: "Microsoft YaHei", size: 16, color: "AAAAAA" })] })] }) },
      footers: { default: new Footer({ children: [new Paragraph({ alignment: AlignmentType.CENTER, border: { top: { style: BorderStyle.SINGLE, size: 4, color: BLUE_LIGHT, space: 4 } },
        children: [new TextRun({ text: "第 ", font: "Microsoft YaHei", size: 16, color: "AAAAAA" }), new TextRun({ children: [PageNumber.CURRENT], font: "Microsoft YaHei", size: 16, color: "AAAAAA" })] })] }) },
      children: [

        // QUICK START
        sec("快速开始（必做，5 分钟）"),

        sub("Q1：确认能打开论坛"),
        step("1", "打开 Chrome 或 Edge。"),
        step("2", "地址栏输入 http://127.0.0.1:8091，回车。"),
        see("登录页，标题「万千帖子，齐聚 AI智联平台。」，有 QQ 号/密码输入框和蓝色登录按钮。"),
        warn("如果打不开：找技术同学启动服务，不用继续往下。"),

        sub("Q2：登录"),
        step("1", "用户名输 demo_student，密码输 123456，点登录。"),
        see("进入社区首页 /community，顶栏有「警院论坛」、搜索框、铃铛、头像，中间有帖子列表。"),
        check("快速开始完成。下面挑模块！"),

        div(),

        // MODULE MENU
        sec("模块菜单（选你感兴趣的）"),
        body("每个模块独立可测，不必全做。做完一个就能交差。", { italics: true, color: "888888" }),
        ...MODULES.map(m => moduleCard(m)),
        greenBox([
          new Paragraph({ alignment: AlignmentType.CENTER, spacing: { before: 0, after: 0 },
            children: [
              new TextRun({ text: "⭐ 最省事方案：", font: "Microsoft YaHei", size: 20, bold: true, color: GREEN }),
              new TextRun({ text: "模块 B（10 分钟） + 模块 F（2 分钟），填个反馈单就走。", font: "Microsoft YaHei", size: 20, color: GREEN })
            ] }),
        ]),

        div(), pb(),

        // MODULE A
        sec("模块 A：注册与登录（15 分钟）"),
        tip("前提：已完成快速开始的 Q1（能打开论坛）。"),

        sub("A1：打开注册页"),
        step("1", "如果已登录，先退出：点右上角头像 → 退出登录。"),
        step("2", "在登录页点邀请码注册。"),
        see("地址变 /register，标题「加入校园论坛，开启交流之旅。」，按钮「注册并登录」。"),

        sub("A2：注册新账号"),
        labeled("QQ 号 / 用户名", "test_uat_01（别和现有的重复）"),
        labeled("密码", "Test1234"),
        labeled("邀请码", "DEMO2026"),
        labeled("院系、区队、年级", "随便填"),
        step("3", "填完点注册并登录。"),
        see("自动跳进社区 /community，有绿色成功提示。"),

        sub("A3：退出登录"),
        step("1", "点右上角头像 → 退出登录。"),
        see("回到登录页。"),

        sub("A4：错误密码"),
        step("1", "用户名输 demo_student，密码故意输错（如 wrong），点登录。"),
        see("登录失败，出现红色错误提示。"),

        sub("A5：重新登录"),
        step("1", "用 demo_student / 123456 正常登录回来。"),
        check("模块 A 完成。继续做别的模块？还是去收尾？"),

        pb(),

        // MODULE B
        sec("模块 B：浏览与搜索（10 分钟）"),
        tip("前提：已完成快速开始，当前已登录 demo_student。"),

        sub("B1：首页帖子列表"),
        step("1", "看 /community 首页，滚到最下面。"),
        see("帖子卡片（标题、作者、日期、点赞数、评论数）；底部可能有分页。"),

        sub("B2：翻页"),
        step("1", "点第 2 页或 ›。"),
        see("帖子换了一批。（帖子太少没分页就跳过。）"),

        sub("B3：侧栏菜单"),
        step("1", "点顶栏最左侧三条横线。"),
        see("左侧滑出菜单，有「首页」「我的帖子」「我的收藏」等导航，还有板块列表。点外面灰色区域关闭。"),

        sub("B4：板块"),
        step("1", "在侧栏点学业研讨（或顶栏点板块）。"),
        see("地址变 /community/boards/study，有板块标题和帖子列表。"),

        sub("B5：帖子详情"),
        step("1", "随便点一篇帖子的标题。"),
        see("帖子正文、点赞按钮、评论区、工具栏（点赞·收藏·举报）。"),

        sub("B6：搜索"),
        step("1", "点顶栏搜索框，输入关键词（如「刑法」），回车。"),
        see("提示「搜索\"刑法\"共 X 条结果」，列表是搜索结果。"),
        check("模块 B 完成。继续做别的模块？还是去收尾？"),

        pb(),

        // MODULE C
        sec("模块 C：发帖与互动（25 分钟）"),
        tip("前提：已完成快速开始，当前已登录 demo_student。"),

        sub("C1：打开发帖页"),
        step("1", "首页点发布帖子，或侧栏点发帖。"),
        see("左侧「发布到」列出板块，右侧有标题和正文输入框，右下角「发布帖子」按钮。"),

        sub("C2：正常发帖"),
        step("1", "左侧选一个板块（如「学业研讨区」）点一下让它高亮。"),
        step("2", "标题填：【UAT测试】正常发帖"),
        step("3", "正文填：测试内容，稍后删除。"),
        step("4", "点发布帖子。"),
        see("提示「帖子已成功发布」或「帖子已进入审核队列」。"),
        tip("记下地址栏里的帖子 ID 数字，后面编辑和删除用到。"),

        sub("C3：空标题拦截"),
        step("1", "再进发帖页，标题留空，正文随便写，点发布。"),
        see("按钮点不了，或提交失败。"),

        sub("C4：空正文拦截"),
        step("1", "标题填「测试空正文」，正文留空，点发布。"),
        see("同样无法发布。"),

        sub("C5：编辑帖子"),
        step("1", "打开你刚发的帖子（C2 那篇）。"),
        step("2", "点工具栏里的编辑（铅笔图标）。"),
        step("3", "标题改成：【UAT测试】已编辑"),
        step("4", "点保存修改。"),
        see("回到详情页，标题已更新。"),

        sub("C6：评论"),
        step("1", "在帖子详情页滚到评论区。"),
        step("2", "输入：测试评论内容"),
        step("3", "点发表评论。"),
        see("提示「评论已发布」，评论立刻出现在列表里。"),

        sub("C7：空评论拦截"),
        step("1", "不输入任何内容，直接点发表评论。"),
        see("按钮灰色点不了，或无法提交。"),

        sub("C8：点赞"),
        step("1", "打开一篇别人发的帖子（不是你自己的）。"),
        step("2", "点点赞（心形），看数字 +1。再点一次，数字减回去。"),

        sub("C9：收藏"),
        step("1", "在同一篇别人帖子，点收藏（书签）。"),
        step("2", "侧栏 → 我的收藏，确认列表里有这篇。"),
        step("3", "回到帖子再点一次收藏取消，回到我的收藏确认它消失了。"),

        sub("C10：举报"),
        step("1", "打开一篇别人帖子，点举报。理由填 UAT测试举报，提交。"),
        see("提示「举报已提交」。"),

        sub("C11：确认别人帖子没有编辑删除"),
        step("1", "打开任意别人帖子，看工具栏。"),
        see("只有点赞、收藏、举报，没有编辑和删除按钮。"),

        sub("C12：删除自己的帖子"),
        step("1", "打开你自己的帖子（C5 改过的那篇）。"),
        step("2", "点删除（垃圾桶），确认。"),
        see("帖子消失，跳回首页。"),
        check("模块 C 完成。继续做别的模块？还是去收尾？"),

        pb(),

        // MODULE D
        sec("模块 D：个人资料（10 分钟）"),
        tip("前提：已完成快速开始，当前已登录 demo_student。"),

        sub("D1：进入个人中心"),
        step("1", "点右上角头像 → 个人中心（或侧栏点「个人中心」）。"),
        see("地址 /community/profile，左侧显示昵称、账号，右侧有编辑表单。"),

        sub("D2：修改并保存"),
        step("1", "昵称改成「UAT测试同学」"),
        step("2", "院系、区队、年级、简介随便改改"),
        step("3", "点保存资料"),
        see("提示「资料已保存」。按 F5 刷新，填的内容还在。"),
        check("模块 D 完成。继续做别的模块？还是去收尾？"),

        pb(),

        // MODULE E
        sec("模块 E：消息通知（15 分钟，需两人）"),
        tip("需要两个账号。一个人用 demo_student，另一个人用 demo01 或其他学生号。可以一个人开两个浏览器窗口切着来。"),

        sub("E1：A 发帖"),
        step("1", "用 demo_student 发一篇帖，标题「【UAT】通知测试」。"),

        sub("E2：B 评论"),
        step("1", "demo_student 退出登录。"),
        step("2", "用另一个号（如 demo01）登录。"),
        step("3", "找到 demo_student 刚发的帖，评论：来自另一个号的评论。"),

        sub("E3：A 查看通知"),
        step("1", "退出 B，重新登录 demo_student。"),
        step("2", "看顶栏铃铛有没有红点。"),
        step("3", "点铃铛（或侧栏「消息」）。"),
        step("4", "点那条通知。"),
        see("地址 /community/messages，有评论相关通知，点击后详情显示，未读数减少。"),
        check("模块 E 完成。继续做别的模块？还是去收尾？"),

        pb(),

        // MODULE F
        sec("模块 F：权限检查（2 分钟）"),
        tip("前提：当前已登录 demo_student（学生号）。"),

        sub("F1：尝试访问管理端"),
        step("1", "在地址栏手动输入 http://127.0.0.1:8091/admin，回车。"),
        see("进不去，被送回社区首页 /community。"),
        check("模块 F 完成。去收尾填反馈单。"),

        pb(),

        // FEEDBACK
        sec("收尾：填反馈"),
        body("每发现一个不对劲就记一条。什么都没发现也写上「全通过」——这对开发很重要。", { bold: true }),
        div(),

        sub("反馈模板"),
        labeled("反馈编号", "01"),
        labeled("哪个模块/哪一步", "模块 __ ，操作：________"),
        labeled("类型", "□ BUG  □ 界面丑  □ 文案不对  □ 安全  □ 建议"),
        labeled("严重程度", "P0=彻底坏了  P1=主要功能坏了  P2=次要  P3=小毛病"),
        labeled("用的账号", "demo_student / 其他________"),
        labeled("我怎么操作的", "1.\n2.\n3."),
        labeled("我以为应该", "________"),
        labeled("实际", "________"),

        div(),

        sub("全通过模板（没问题时用）"),
        greenBox([
          new Paragraph({ spacing: { before: 0, after: 40 }, children: [new TextRun({ text: "反馈 — 通过", font: "Microsoft YaHei", size: 20, bold: true, color: GREEN })] }),
          new Paragraph({ spacing: { before: 20, after: 20 }, children: [new TextRun({ text: "环境：本地 8091", font: "Microsoft YaHei", size: 20, color: "333333" })] }),
          new Paragraph({ spacing: { before: 20, after: 20 }, children: [new TextRun({ text: "测试模块：__（你做了哪几个）", font: "Microsoft YaHei", size: 20, color: "333333" })] }),
          new Paragraph({ spacing: { before: 20, after: 0 }, children: [new TextRun({ text: "结论：以上模块均符合预期", font: "Microsoft YaHei", size: 20, color: "333333" })] }),
        ]),

        div(),

        body("把反馈内容发给开发或贴到 Cursor 对话里。", { bold: true }),

        new Paragraph({ spacing: { before: 200, after: 100 }, children: [] }),
        new Paragraph({ children: [
          new TextRun({ text: "测试人：________    ", font: "Microsoft YaHei", size: 22, color: "555555" }),
          new TextRun({ text: "日期：________", font: "Microsoft YaHei", size: 22, color: "555555" })
        ] }),
      ]
    }
  ]
});

Packer.toBuffer(doc).then(buffer => {
  const outPath = path.join(__dirname, "学生端测试指南.docx");
  fs.writeFileSync(outPath, buffer);
  console.log("Done: " + outPath);
}).catch(err => {
  console.error("Error:", err.message);
  process.exit(1);
});
