<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()
const boards = ref([])
const form = ref({
  boardId: '',
  title: '',
  content: '',
  tags: '',
  attachmentName: '',
  attachmentType: 'link',
  attachmentUrl: '',
})

onMounted(async () => {
  boards.value = await forumApi.getBoards()
  form.value.boardId = boards.value[0]?.id || ''
})

async function submit() {
  const attachments = form.value.attachmentName
    ? [
        {
          name: form.value.attachmentName,
          type: form.value.attachmentType,
          url: form.value.attachmentUrl,
        },
      ]
    : []
  const post = await forumApi.createPost({
    boardId: form.value.boardId,
    authorId: session.currentUser.id,
    title: form.value.title,
    content: form.value.content,
    tags: form.value.tags.split(',').map((item) => item.trim()).filter(Boolean),
    attachments,
  })
  session.setFlash(
    post.status === 'pending' ? '帖子已进入审核队列。' : '帖子已成功发布。',
    post.status === 'pending' ? 'warning' : 'success',
  )
  if (post.status === 'published') {
    router.push(`/community/posts/${post.id}`)
  } else {
    router.push('/community')
  }
}
</script>

<template>
  <section class="panel form-panel">
    <p class="eyebrow">发布新帖</p>
    <h2>用最小表单把社区内容流跑起来</h2>
    <div class="form-grid">
      <label>
        <span>选择板块</span>
        <select v-model="form.boardId">
          <option v-for="board in boards" :key="board.id" :value="board.id">{{ board.name }}</option>
        </select>
      </label>
      <label class="full-span">
        <span>标题</span>
        <input v-model="form.title" type="text" placeholder="例如：RAG 课程助手的 MVP 应该先做什么？" />
      </label>
      <label class="full-span">
        <span>正文</span>
        <textarea v-model="form.content" rows="8" placeholder="写下你的问题、经验或活动信息。" />
      </label>
      <label>
        <span>标签（逗号分隔）</span>
        <input v-model="form.tags" type="text" placeholder="RAG, 检索, 课程助手" />
      </label>
      <label>
        <span>附件类型</span>
        <select v-model="form.attachmentType">
          <option value="link">网盘链接</option>
          <option value="document">文档</option>
          <option value="image">图片</option>
        </select>
      </label>
      <label>
        <span>附件名称</span>
        <input v-model="form.attachmentName" type="text" placeholder="例如：项目草图.pdf" />
      </label>
      <label>
        <span>附件 URL</span>
        <input v-model="form.attachmentUrl" type="text" placeholder="https://example.com/resource" />
      </label>
    </div>
    <button
      class="primary-button"
      :disabled="!form.boardId || !form.title.trim() || !form.content.trim()"
      @click="submit"
    >
      提交帖子
    </button>
  </section>
</template>
