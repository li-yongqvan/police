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
  <div class="page-stack">
    <section class="hero-strip panel">
      <div class="hero-copy">
        <p class="eyebrow">技术社区氛围优先</p>
        <h2>把 AI 学习、协会活动和技术问答都汇在同一块前台里。</h2>
        <p>
          这版 MVP 重点展示学院级论坛的内容组织能力，同时保留最小审核和配置闭环，方便做展示与答辩。
        </p>
      </div>
      <div class="hero-metrics">
        <div class="metric-card">
          <strong>{{ boards.length }}</strong>
          <span>开放板块</span>
        </div>
        <div class="metric-card">
          <strong>{{ posts.length }}</strong>
          <span>公开帖子</span>
        </div>
        <div class="metric-card">
          <strong>{{ overview?.pendingAuditCount ?? 0 }}</strong>
          <span>待审核</span>
        </div>
      </div>
    </section>

    <section class="two-column">
      <section class="panel content-panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">三大核心板块</p>
            <h3>用板块组织内容，用帖子带动讨论</h3>
          </div>
        </div>
        <div class="board-grid">
          <RouterLink
            v-for="board in boards"
            :key="board.id"
            :to="`/community/boards/${board.slug}`"
            class="board-card"
          >
            <h4>{{ board.name }}</h4>
            <p>{{ board.description }}</p>
          </RouterLink>
        </div>
      </section>

      <section class="panel content-panel">
        <div class="section-title">
          <div>
            <p class="eyebrow">精选内容</p>
            <h3>适合展示时先点开的高信号帖子</h3>
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
    </section>

    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">最新帖子流</p>
          <h3>展示社区活跃度和技术讨论感</h3>
        </div>
      </div>
      <div class="post-list">
        <RouterLink v-for="post in posts" :key="post.id" :to="`/community/posts/${post.id}`" class="post-card">
          <div class="post-topline">
            <span class="badge subtle">{{ post.status }}</span>
            <span>{{ post.commentCount }} 评论</span>
          </div>
          <h4>{{ post.title }}</h4>
          <p>{{ post.content }}</p>
          <div class="tag-row">
            <span v-for="tag in post.tags" :key="tag" class="tag">{{ tag }}</span>
          </div>
        </RouterLink>
      </div>
    </section>
  </div>
</template>
