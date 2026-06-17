<template>
  <div class="post-create-page">
    <h1>发布新帖子</h1>
    <div v-if="authStore.userLevel < 1" class="level-notice">
      需要等级1才能发帖
    </div>
    <form v-else @submit.prevent="handleSubmit">
      <div class="form-group">
        <label>板块</label>
        <select v-model="form.board_id" required>
          <option value="" disabled>选择板块</option>
          <option v-for="board in forumStore.boards" :key="board.id" :value="board.id">
            {{ board.name }}
          </option>
        </select>
      </div>
      <div class="form-group">
        <label>标题</label>
        <input v-model="form.title" placeholder="请输入标题" required minlength="2" />
      </div>
      <div class="form-group">
        <label>内容</label>
        <textarea v-model="form.content" placeholder="请输入内容（至少10个字符）" rows="10" required minlength="10"></textarea>
      </div>
      <div class="form-group">
        <label>附件上传</label>
        <input type="file" multiple accept="image/*,.pdf,.doc,.docx,.txt,.md" @change="handleFiles" />
        <div v-if="fileNames.length" class="file-list">
          <span v-for="(name, i) in fileNames" :key="i" class="file-tag">{{ name }}</span>
        </div>
      </div>
      <div class="form-group">
        <label>网盘链接</label>
        <input v-model="linkUrl" placeholder="可选，输入网盘分享链接" />
      </div>
      <div v-if="error" class="error">{{ error }}</div>
      <button type="submit" :disabled="loading" class="btn-primary">
        {{ loading ? '发布中...' : '发布' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useForumStore } from '../stores/forum'
import { useAuthStore } from '../stores/auth'
import { uploadAttachment } from '../api/forum'

const router = useRouter()
const route = useRoute()
const forumStore = useForumStore()
const authStore = useAuthStore()
const error = ref('')
const loading = ref(false)
const files = ref([])
const fileNames = ref([])
const linkUrl = ref('')

const form = reactive({
  title: '',
  content: '',
  board_id: Number(route.query.board_id) || 0,
  attachment_ids: [],
})

onMounted(async () => {
  if (forumStore.boards.length === 0) {
    await forumStore.loadBoards()
  }
})

function handleFiles(e) {
  files.value = Array.from(e.target.files)
  fileNames.value = files.value.map(f => f.name)
}

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    // Upload attachments first
    const attachmentIds = []
    for (const file of files.value) {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('type', file.type.startsWith('image') ? 'image' : 'document')
      const { data } = await uploadAttachment(formData)
      attachmentIds.push(data.id)
    }

    // Add link as attachment
    if (linkUrl.value) {
      const linkFormData = new FormData()
      linkFormData.append('type', 'link')
      linkFormData.append('link_url', linkUrl.value)
      const { data } = await uploadAttachment(linkFormData)
      attachmentIds.push(data.id)
    }

    const postData = {
      title: form.title,
      content: form.content,
      board_id: form.board_id,
      attachment_ids: attachmentIds,
    }

    const post = await forumStore.createPost(postData)
    router.push(`/posts/${post.id}`)
  } catch (e) {
    error.value = e.response?.data?.error || '发布失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.post-create-page {
  max-width: 800px;
  margin: 0 auto;
}

.post-create-page h1 {
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.25rem;
  color: #666;
  font-weight: 500;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

.form-group textarea {
  resize: vertical;
  min-height: 200px;
}

.file-list {
  margin-top: 0.5rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.file-tag {
  background: #e3f2fd;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.error {
  color: #e74c3c;
  background: #fdf0ef;
  padding: 0.5rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.btn-primary {
  padding: 0.75rem 2rem;
  background: #4a90d9;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
}

.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }

.level-notice {
  padding: 1rem;
  background: #fff3cd;
  border-radius: 6px;
  color: #856404;
}

@media (max-width: 768px) {
  .post-create-page { padding: 0 0.5rem; }
}
</style>
