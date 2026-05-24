<template>
  <div class="board-detail-page">
    <div v-if="loading" class="loading">加载中...</div>
    <template v-else>
      <div class="board-header">
        <h1>{{ board?.name }}</h1>
        <p>{{ board?.description }}</p>
        <router-link v-if="isLoggedIn && authStore.userLevel >= 1"
          :to="{ name: 'PostCreate', query: { board_id: boardId } }"
          class="btn-primary">发帖</router-link>
      </div>
      <div v-if="posts.length === 0" class="empty">暂无帖子</div>
      <PostCard v-for="post in posts" :key="post.id" :post="post" />
      <div class="pagination">
        <button :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <span>第 {{ page }} 页</span>
        <button :disabled="posts.length < 20" @click="changePage(page + 1)">下一页</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useForumStore } from '../stores/forum'
import { useAuthStore } from '../stores/auth'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const forumStore = useForumStore()
const authStore = useAuthStore()
const boardId = computed(() => Number(route.params.id))
const page = ref(1)
const loading = ref(true)
const posts = computed(() => forumStore.posts)
const board = computed(() => forumStore.boards.find(b => b.id === boardId.value))
const isLoggedIn = computed(() => authStore.isLoggedIn)

onMounted(async () => {
  await forumStore.loadPosts(boardId.value, page.value)
  loading.value = false
})

async function changePage(p) {
  page.value = p
  loading.value = true
  await forumStore.loadPosts(boardId.value, p)
  loading.value = false
}
</script>

<style scoped>
.board-detail-page {
  max-width: 900px;
  margin: 0 auto;
}

.board-header {
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.board-header h1 {
  color: #333;
}

.board-header p {
  color: #666;
  margin-bottom: 0.5rem;
}

.btn-primary {
  display: inline-block;
  padding: 0.5rem 1.5rem;
  background-color: #4a90d9;
  color: white;
  border-radius: 6px;
  text-decoration: none;
}

.empty {
  text-align: center;
  padding: 3rem;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.pagination button {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #999;
}

@media (max-width: 768px) {
  .board-detail-page {
    padding: 0 0.5rem;
  }
}
</style>
