<template>
  <div class="home-page">
    <section class="hero">
      <h1>AI智联论坛</h1>
      <p class="subtitle">面向高校学院的人工智能技术交流社区</p>
    </section>

    <section class="boards">
      <h2>核心板块</h2>
      <div class="board-list">
        <div class="board-card" v-for="board in boards" :key="board.id">
          <h3>{{ board.name }}</h3>
          <p>{{ board.description }}</p>
        </div>
      </div>
    </section>

    <section class="recent-posts">
      <h2>最新帖子</h2>
      <p class="placeholder">暂无内容，快去发布第一篇帖子吧！</p>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const boards = ref([
  { id: 1, name: 'AI学习交流区', description: '分享学习心得、讨论AI技术、交流学习经验' },
  { id: 2, name: '协会公告&活动区', description: '查看协会最新公告、活动通知、招新信息' },
  { id: 3, name: '技术问答求助区', description: '提出技术问题、寻求帮助、解答疑惑' },
])
</script>

<style scoped>
.hero {
  text-align: center;
  padding: 3rem 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  margin-bottom: 2rem;
}

.hero h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  font-size: 1.125rem;
  opacity: 0.9;
}

.boards h2,
.recent-posts h2 {
  margin-bottom: 1rem;
  color: #333;
}

.board-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.board-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.board-card h3 {
  color: #4a90d9;
  margin-bottom: 0.5rem;
}

.board-card p {
  color: #666;
  font-size: 0.875rem;
}

.placeholder {
  text-align: center;
  color: #999;
  padding: 2rem;
  background: white;
  border-radius: 8px;
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 1.75rem;
  }
  .subtitle {
    font-size: 1rem;
  }
}
</style>
