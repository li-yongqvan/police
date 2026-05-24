<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'

const overview = ref(null)

onMounted(async () => {
  overview.value = await adminApi.getOverview()
})
</script>

<template>
  <div v-if="overview" class="page-stack page-stack--admin">
    <section class="panel content-panel section-activity">
      <div class="section-title section-title--compact">
        <div>
          <p class="eyebrow">板块热度</p>
          <h3>活跃度分布</h3>
        </div>
      </div>
      <div class="activity-grid">
        <article v-for="item in overview.boardActivity" :key="item.boardId" class="activity-card">
          <strong>{{ item.name }}</strong>
          <span>{{ item.count }} 篇</span>
        </article>
      </div>
    </section>

    <section class="hero-strip panel hero-strip--admin">
      <div class="hero-copy">
        <p class="eyebrow">中台概览</p>
        <h2 class="hero-title-short">管理数据</h2>
        <p class="hero-desc-long">先把可控、可管、可展示的管理闭环做出来。</p>
      </div>
      <div class="hero-metrics admin stats-strip" aria-label="运营指标">
        <div class="metric-card metric-card--compact">
          <strong>{{ overview.userCount }}</strong>
          <span>用户</span>
        </div>
        <div class="metric-card metric-card--compact">
          <strong>{{ overview.todayPostCount }}</strong>
          <span>今日帖</span>
        </div>
        <div class="metric-card metric-card--compact">
          <strong>{{ overview.pendingAuditCount }}</strong>
          <span>待审</span>
        </div>
        <div class="metric-card metric-card--compact">
          <strong>{{ overview.postCount }}</strong>
          <span>公开帖</span>
        </div>
      </div>
    </section>
  </div>
</template>
