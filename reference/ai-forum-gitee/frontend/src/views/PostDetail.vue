<template>
  <div class="post-detail-page">
    <div v-if="loading" class="loading">加载中...</div>
    <template v-else-if="post">
      <div class="post-header">
        <h1>{{ post.title }}</h1>
        <div class="post-meta">
          <span class="author">{{ post.author_name }}</span>
          <span class="board">{{ post.board_name }}</span>
          <span class="time">{{ formatDate(post.created_at) }}</span>
        </div>
        <div v-if="post.is_pinned" class="badge pinned">置顶</div>
        <div v-if="post.is_featured" class="badge featured">精华</div>
      </div>
      <div class="post-content">{{ post.content }}</div>
      <div class="post-attachments" v-if="post.attachments?.length">
        <h3>附件</h3>
        <div class="attachment-list">
          <img v-for="att in post.attachments.filter(a => a.file_type === 'image')"
            :key="att.id" :src="getFileUrl(att.file_path)" :alt="att.filename" class="attachment-image" />
          <a v-for="att in post.attachments.filter(a => a.file_type === 'document')"
            :key="att.id" :href="getFileUrl(att.file_path)" download class="attachment-link">
            {{ att.filename }}
          </a>
          <a v-for="att in post.attachments.filter(a => a.file_type === 'link')"
            :key="att.id" :href="att.file_path" target="_blank" rel="noopener" class="attachment-link">
            网盘链接
          </a>
        </div>
      </div>
      <div class="post-actions">
        <button class="action-btn" @click="handleLike">
          {{ post.user_liked ? '已点赞' : '点赞' }} ({{ post.like_count }})
        </button>
        <button class="action-btn" @click="handleCollect">
          {{ post.user_collected ? '已收藏' : '收藏' }}
        </button>
        <router-link v-if="isAuthor" :to="`/post/${postId}/edit`" class="action-btn">编辑</router-link>
        <button v-if="isAuthor" class="action-btn danger" @click="handleDelete">删除</button>
      </div>
      <div class="comments-section">
        <h3>评论 ({{ post.comment_count }})</h3>
        <div v-if="isLoggedIn && authStore.userLevel >= 1" class="comment-form">
          <textarea v-model="commentContent" placeholder="写评论..." rows="3"></textarea>
          <button @click="handleComment" :disabled="!commentContent">发表评论</button>
        </div>
        <div v-else-if="isLoggedIn" class="level-notice">
          需要等级1才能评论
        </div>
        <div class="comment-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-meta">{{ formatDate(comment.created_at) }}</div>
          </div>
        </div>
      </div>
    </template>
    <div v-else class="empty">帖子不存在</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useForumStore } from '../stores/forum'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const forumStore = useForumStore()
const authStore = useAuthStore()
const postId = computed(() => Number(route.params.id))
const post = computed(() => forumStore.currentPost)
const comments = computed(() => forumStore.comments)
const loading = computed(() => forumStore.loading)
const isLoggedIn = computed(() => authStore.isLoggedIn)
const isAuthor = computed(() => authStore.user?.id === post.value?.author_id)
const commentContent = ref('')

onMounted(async () => {
  await forumStore.loadPost(postId.value)
  await forumStore.loadComments(postId.value)
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN') + ' ' + d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function getFileUrl(path) {
  if (path.startsWith('http')) return path
  return path
}

async function handleLike() {
  if (!isLoggedIn.value) { router.push('/login'); return }
  await forumStore.likePost(postId.value)
}

async function handleCollect() {
  if (!isLoggedIn.value) { router.push('/login'); return }
  await forumStore.collectPost(postId.value)
}

async function handleComment() {
  if (!commentContent.value.trim()) return
  await forumStore.createComment(postId.value, commentContent.value)
  commentContent.value = ''
}

async function handleDelete() {
  if (!confirm('确定要删除这篇帖子吗？')) return
  await forumStore.deletePost(postId.value)
  router.push(`/boards/${post.value?.board_id}`)
}
</script>

<style scoped>
.post-detail-page {
  max-width: 900px;
  margin: 0 auto;
}

.post-header {
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.post-header h1 {
  color: #333;
  margin-bottom: 0.5rem;
}

.post-meta {
  display: flex;
  gap: 1rem;
  color: #999;
  font-size: 0.875rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  margin-right: 0.5rem;
}

.badge.pinned { background: #e3f2fd; color: #1976d2; }
.badge.featured { background: #fff3e0; color: #f57c00; }

.post-content {
  white-space: pre-wrap;
  line-height: 1.8;
  margin-bottom: 1.5rem;
  color: #333;
}

.post-attachments {
  margin-bottom: 1.5rem;
}

.attachment-image {
  max-width: 100%;
  max-height: 400px;
  border-radius: 8px;
  margin: 0.5rem 0;
  display: block;
}

.attachment-link {
  display: inline-block;
  padding: 0.5rem 1rem;
  background: #f5f5f5;
  border-radius: 4px;
  color: #4a90d9;
  text-decoration: none;
  margin: 0.25rem;
}

.post-actions {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.action-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  color: #333;
  text-decoration: none;
}

.action-btn:hover { background: #f5f5f5; }
.action-btn.danger { color: #e74c3c; border-color: #e74c3c; }

.comments-section {
  border-top: 1px solid #eee;
  padding-top: 1.5rem;
}

.comment-form {
  margin-bottom: 1rem;
}

.comment-form textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  resize: vertical;
  font-size: 1rem;
}

.comment-form button {
  margin-top: 0.5rem;
  padding: 0.5rem 1.5rem;
  background: #4a90d9;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.comment-form button:disabled { opacity: 0.5; cursor: not-allowed; }

.level-notice {
  padding: 1rem;
  background: #fff3cd;
  border-radius: 6px;
  color: #856404;
  margin-bottom: 1rem;
}

.comment-item {
  padding: 1rem 0;
  border-bottom: 1px solid #f5f5f5;
}

.comment-content { color: #333; margin-bottom: 0.25rem; }
.comment-meta { color: #999; font-size: 0.75rem; }

.loading, .empty {
  text-align: center;
  padding: 3rem;
  color: #999;
}

@media (max-width: 768px) {
  .post-detail-page { padding: 0 0.5rem; }
  .post-meta { flex-direction: column; gap: 0.25rem; }
  .post-actions { flex-direction: column; }
}
</style>
