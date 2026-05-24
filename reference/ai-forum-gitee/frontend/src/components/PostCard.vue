<template>
  <router-link :to="`/posts/${post.id}`" class="post-card">
    <div class="post-header">
      <h3 class="post-title">{{ post.title }}</h3>
      <div v-if="post.is_pinned" class="badge pinned">置顶</div>
      <div v-if="post.is_featured" class="badge featured">精华</div>
    </div>
    <div class="post-meta">
      <span class="author">{{ post.author_name }}</span>
      <span class="board">{{ post.board_name }}</span>
      <span class="time">{{ formatDate(post.created_at) }}</span>
    </div>
    <div class="post-stats">
      <span class="stat"><i class="icon-like"></i> {{ post.like_count }}</span>
      <span class="stat"><i class="icon-comment"></i> {{ post.comment_count }}</span>
    </div>
  </router-link>
</template>

<script setup>
defineProps({
  post: { type: Object, required: true },
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now - d
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.post-card {
  display: block;
  background: white;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  margin-bottom: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  text-decoration: none;
  color: inherit;
  transition: box-shadow 0.2s;
}

.post-card:hover {
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.post-title {
  font-size: 1rem;
  color: #333;
  margin: 0;
  flex: 1;
}

.badge {
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-size: 0.7rem;
}

.badge.pinned { background: #e3f2fd; color: #1976d2; }
.badge.featured { background: #fff3e0; color: #f57c00; }

.post-meta {
  display: flex;
  gap: 1rem;
  color: #999;
  font-size: 0.8rem;
  margin-bottom: 0.5rem;
}

.post-stats {
  display: flex;
  gap: 1rem;
  color: #999;
  font-size: 0.8rem;
}

.stat { display: flex; align-items: center; gap: 0.25rem; }

@media (max-width: 768px) {
  .post-card { padding: 0.75rem 1rem; }
  .post-meta { flex-wrap: wrap; gap: 0.5rem; }
}
</style>
