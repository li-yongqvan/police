# 用 Cursor 生成这份 PPT（操作说明）

## 已生成产物

| 文件 | 说明 |
|------|------|
| `origin_image/slide_01.png` | 第 1 页：智联 |
| `origin_image/slide_02.png` | 第 2 页：论坛 + 中台 |
| `origin_image/slide_03.png` | 第 3 页：一端（PC + 手机截图） |
| `origin_image/slide_04.png` | 第 4 页：80% |
| `ai-forum-deck.pptx` | 已合成的 4 页 PPT（含演讲者备注） |

直接双击打开 **`ai-forum-deck.pptx`** 即可预览。

---

## 在 Cursor 里重做 / 改某一页

在 **Agent 对话**里发送（可复制）：

```text
请阅读 docs/ppt-ai-forum-deck/AI_GENERATE.md，
按许岑式黑底白字风格重新生成第 N 页幻灯片图片（16:9），
保存到 docs/ppt-ai-forum-deck/origin_image/slide_0N.png，
然后重新合成 ai-forum-deck.pptx（演讲备注用 speech.md）。
```

把 `N` 换成 1、2 或 3。不必三页一起重做。

---

## 推荐顺序（以后自己做）

1. **先看** `outline.md` — 确认要表达什么  
2. **复制** `AI_GENERATE.md` 里「全局 + 某一页」— 让 Cursor **生成图片**  
3. **保存** 到 `origin_image/slide_01~03.png`  
4. **合成** — 对 Agent 说「用 origin_image 合成 pptx」，或手动插入 PowerPoint  
5. **彩排** — 对照 `speech.md` 口述，画面上的字不要念太长  

---

## 注意

- Cursor 出图是**整张幻灯片位图**，不是可编辑文字框；改字需要重新出图。  
- 若某页中文糊了或不对，只重做该页即可。  
- `deck_spec.json` 可忽略，除非你用 codex-ppt 自动化工具链。
