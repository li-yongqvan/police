<template>
  <div class="board-list-page">
    <h1>核心板块</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="boards.length === 0" class="empty">暂无板块</div>
    <div v-else class="board-grid">
      <BoardCard v-for="board in boards" :key="board.id" :board="board" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useForumStore } from '../stores/forum'
import BoardCard from '../components/BoardCard.vue'

const forumStore = useForumStore()
const boards = ref([])
const loading = ref(true)

onMounted(async () => {
  await forumStore.loadBoards()
  boards.value = forumStore.boards
  loading.value = false
})
</script>

<style scoped>
.board-list-page {
  max-width: 1200px;
  margin: 0 auto;
}

.board-list-page h1 {
  margin-bottom: 1.5rem;
  color: #333;
}

.board-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.loading, .empty {
  text-align: center;
  padding: 3rem;
  color: #999;
  font-size: 1.125rem;
}

@media (max-width: 768px) {
  .board-grid {
    grid-template-columns: 1fr;
  }
}
</style>
