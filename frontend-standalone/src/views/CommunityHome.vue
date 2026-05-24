<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { adminApi, forumApi } from '../api'

const boards = ref([])
const posts = ref([])
const overview = ref(null)

const featuredPosts = computed(() => posts.value.filter((item) => item.isFeatured))

onMounted(async () => {
  boards.value = await forumApi.getBoards()
  posts.value = await forumApi.getPosts()
  try {
    overview.value = await adminApi.getOverview()
  } catch {
    overview.value = null
  }
})
</script>

<template>
  <div class="page-stack page-stack--browse">
    <section class="section-feed panel content-panel">
      <div class="section-title section-title--compact">
        <div>
          <p class="eyebrow">刷帖</p>
          <h3>最新帖子</h3>
        </div>
        <span class="feed-count">{{ posts.length }} 篇</span>
      </div>
      <div class="post-list post-list--feed">
        <RouterLink v-for="post in posts" :key="post.id" :to="`/community/posts/${post.id}`" class="post-card">
          <div class="post-topline">
            <span class="badge subtle">{{ post.status }}</span>
            <span>{{ post.commentCount }} 评论 · {{ post.likeCount }} 赞</span>
          </div>
          <h4>{{ post.title }}</h4>
          <p>{{ post.content }}</p>
          <div v-if="post.tags?.length" class="tag-row">
            <span v-for="tag in post.tags" :key="tag" class="tag">{{ tag }}</span>
          </div>
        </RouterLink>
      </div>
    </section>

    <section class="section-boards panel content-panel">
      <div class="section-title section-title--compact">
        <div>
          <p class="eyebrow">板块</p>
          <h3>快速进入</h3>
        </div>
      </div>
      <div class="board-strip">
        <RouterLink
          v-for="board in boards"
          :key="board.id"
          :to="`/community/boards/${board.slug}`"
          class="board-chip"
        >
          {{ board.name }}
        </RouterLink>
      </div>
    </section>

    <section v-if="featuredPosts.length || posts.length" class="section-featured panel content-panel">
      <div class="section-title section-title--compact">
        <div>
          <p class="eyebrow">精选</p>
          <h3>高信号帖子</h3>
        </div>
      </div>
      <div class="post-list compact">
        <RouterLink
          v-for="post in featuredPosts.length ? featuredPosts : posts.slice(0, 2)"
          :key="post.id"
          :to="`/community/posts/${post.id}`"
          class="post-card"
        >
          <div class="post-topline">
            <span class="badge">精选</span>
            <span>{{ post.likeCount }} 赞</span>
          </div>
          <h4>{{ post.title }}</h4>
          <p>{{ post.content }}</p>
        </RouterLink>
      </div>
    </section>

    <section class="hero-strip panel hero-strip--secondary">
      <div class="hero-copy">
        <p class="eyebrow">AI 智联论坛</p>
        <h2 class="hero-title-short">学院 AI 社区</h2>
        <p class="hero-desc-long">
          这版 MVP 重点展示学院级论坛的内容组织能力，同时保留最小审核和配置闭环，方便做展示与答辩。
        </p>
      </div>
      <div class="hero-metrics stats-strip" aria-label="社区概览">
        <div class="metric-card metric-card--compact">
          <strong>{{ boards.length }}</strong>
          <span>板块</span>
        </div>
        <div class="metric-card metric-card--compact">
          <strong>{{ posts.length }}</strong>
          <span>帖子</span>
        </div>
        <div class="metric-card metric-card--compact">
          <strong>{{ overview?.pendingAuditCount ?? 0 }}</strong>
          <span>待审</span>
        </div>
      </div>
    </section>
  </div>
</template>
