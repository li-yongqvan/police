<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'

const overview = ref(null)

onMounted(async () => {
  overview.value = await adminApi.getOverview()
})
</script>

<template>
  <div v-if="overview" class="page-stack">
    <section class="hero-strip panel">
      <div class="hero-copy">
        <p class="eyebrow">最小中台概览</p>
        <h2>先把可控、可管、可展示的管理闭环做出来。</h2>
      </div>
      <div class="hero-metrics admin">
        <div class="metric-card">
          <strong>{{ overview.userCount }}</strong>
          <span>注册用户数</span>
        </div>
        <div class="metric-card">
          <strong>{{ overview.todayPostCount }}</strong>
          <span>今日发帖量</span>
        </div>
        <div class="metric-card">
          <strong>{{ overview.pendingAuditCount }}</strong>
          <span>待审核数</span>
        </div>
        <div class="metric-card">
          <strong>{{ overview.postCount }}</strong>
          <span>公开帖子量</span>
        </div>
      </div>
    </section>

    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">板块活跃度</p>
          <h3>当前社区热度分布</h3>
        </div>
      </div>
      <div class="activity-grid">
        <article v-for="item in overview.boardActivity" :key="item.boardId" class="activity-card">
          <strong>{{ item.name }}</strong>
          <span>{{ item.count }} 篇公开帖</span>
        </article>
      </div>
    </section>
  </div>
</template>
