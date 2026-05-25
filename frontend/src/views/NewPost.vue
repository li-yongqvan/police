<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { formatApiError } from '../api/errors'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()
const boards = ref([])
const uploading = ref(false)
const needsLevelForAttachment = computed(
  () => (session.currentUser?.level ?? 1) < 2,
)
const form = ref({
  boardId: '',
  title: '',
  content: '',
  tags: '',
  attachmentType: 'link',
  attachmentName: '',
  attachmentUrl: '',
  attachmentFile: null,
})

onMounted(async () => {
  boards.value = await forumApi.getBoards()
  form.value.boardId = boards.value[0]?.id || ''
})

async function uploadIfNeeded() {
  if (form.value.attachmentType === 'link' && form.value.attachmentUrl) {
    return [
      await forumApi.uploadAttachment({ type: 'link', linkUrl: form.value.attachmentUrl }),
    ]
  }
  if (form.value.attachmentFile) {
    return [
      await forumApi.uploadAttachment({
        type: form.value.attachmentType,
        file: form.value.attachmentFile,
      }),
    ]
  }
  return []
}

async function submit() {
  uploading.value = true
  try {
    const attachmentIds = await uploadIfNeeded()
    const post = await forumApi.createPost({
      boardId: form.value.boardId,
      title: form.value.title,
      content: form.value.content,
      attachmentIds,
    })
    session.setFlash(
      post.status === 'pending_review' ? '帖子已进入审核队列。' : '帖子已成功发布。',
      post.status === 'pending_review' ? 'warning' : 'success',
    )
    if (post.status === 'published') {
      router.push(`/community/posts/${post.id}`)
    } else {
      router.push('/community')
    }
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    uploading.value = false
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
      <p v-if="needsLevelForAttachment" class="status-hint">
        上传文件或链接附件需要账号等级 ≥ 2（当前 Lv.{{ session.currentUser?.level ?? 1 }}）。
      </p>
      <label>
        <span>附件类型</span>
        <select v-model="form.attachmentType">
          <option value="link">网盘链接</option>
          <option value="document">文档</option>
          <option value="image">图片</option>
        </select>
      </label>
      <label v-if="form.attachmentType === 'link'" class="full-span">
        <span>链接地址</span>
        <input v-model="form.attachmentUrl" type="url" placeholder="https://..." />
      </label>
      <label v-else class="full-span">
        <span>上传文件</span>
        <input type="file" @change="(e) => (form.attachmentFile = e.target.files?.[0])" />
      </label>
    </div>
    <button class="primary-button" :disabled="uploading" @click="submit">
      {{ uploading ? '提交中…' : '发布帖子' }}
    </button>
  </section>
</template>
