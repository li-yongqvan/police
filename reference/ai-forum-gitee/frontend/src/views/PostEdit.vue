<template>
  <div class="post-edit-page">
    <h1>编辑帖子</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <form v-else-if="post" @submit.prevent="handleSubmit">
      <div class="form-group">
        <label>标题</label>
        <input v-model="form.title" required />
      </div>
      <div class="form-group">
        <label>内容</label>
        <textarea v-model="form.content" rows="10" required minlength="10"></textarea>
      </div>
      <div v-if="error" class="error">{{ error }}</div>
      <div class="actions">
        <button type="submit" :disabled="loading">保存</button>
        <router-link :to="`/posts/${postId}`" class="btn-cancel">取消</router-link>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useForumStore } from '../stores/forum'

const route = useRoute()
const router = useRouter()
const forumStore = useForumStore()
const postId = computed(() => Number(route.params.id))
const post = computed(() => forumStore.currentPost)
const error = ref('')
const loading = ref(true)

const form = reactive({
  title: '',
  content: '',
})

onMounted(async () => {
  await forumStore.loadPost(postId.value)
  if (post.value) {
    form.title = post.value.title
    form.content = post.value.content
  }
  loading.value = false
})

async function handleSubmit() {
  error.value = ''
  try {
    await forumStore.updatePost(postId.value, { title: form.title, content: form.content })
    router.push(`/posts/${postId.value}`)
  } catch (e) {
    error.value = e.response?.data?.error || '保存失败，请重试'
  }
}
</script>

<style scoped>
.post-edit-page {
  max-width: 800px;
  margin: 0 auto;
}

.post-edit-page h1 {
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
.form-group textarea {
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

.error {
  color: #e74c3c;
  background: #fdf0ef;
  padding: 0.5rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.actions button {
  padding: 0.75rem 2rem;
  background: #4a90d9;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
}

.btn-cancel {
  padding: 0.75rem 1.5rem;
  color: #666;
  text-decoration: none;
  border: 1px solid #ddd;
  border-radius: 6px;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #999;
}

@media (max-width: 768px) {
  .post-edit-page { padding: 0 0.5rem; }
}
</style>
