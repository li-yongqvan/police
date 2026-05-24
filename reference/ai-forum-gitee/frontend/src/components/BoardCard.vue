<template>
  <router-link :to="`/boards/${board.id}`" class="board-card">
    <div class="board-info">
      <h3>{{ board.name }}</h3>
      <p>{{ board.description }}</p>
    </div>
    <div class="board-stats">
      <span class="post-count">{{ board.post_count || 0 }} 帖子</span>
    </div>
  </router-link>
</template>

<script setup>
defineProps({
  board: { type: Object, required: true },
})
</script>

<style scoped>
.board-card {
  display: block;
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s, box-shadow 0.2s;
}

.board-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.board-info h3 {
  color: #4a90d9;
  margin-bottom: 0.5rem;
}

.board-info p {
  color: #666;
  font-size: 0.875rem;
}

.board-stats {
  margin-top: 0.75rem;
  color: #999;
  font-size: 0.75rem;
}
</style>
