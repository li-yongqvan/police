<script setup>
import { onMounted, ref } from 'vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { RouterLink, useRouter } from 'vue-router'
import { adminApi } from '../api'
import { formatApiError } from '../api/errors'
import { useSessionStore } from '../stores/session'

const router = useRouter()

const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '管理概览' },
]
const session = useSessionStore()

const overview = ref(null)
const loading = ref(true)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    overview.value = await adminApi.getOverview()
  } catch (err) {
    error.value = formatApiError(
      err,
      '无法加载管理数据：请确认后端服务与 Cloudflare 隧道已启动，并使用正式域名或网关 8888 访问',
    )
    overview.value = null
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page gx-admin-overview">
      <GxBreadcrumb :items="breadcrumbItems" />
    <p v-if="loading" class="gx-muted">正在加载管理数据…</p>
    <div v-else-if="error" class="gx-card gx-admin-error">
      <p class="gx-error">{{ error }}</p>
      <p class="gx-muted">
        学生/管理入口：<strong>https://api.shgenren.dpdns.org</strong>（勿用失效的 trycloudflare 链接）。
      </p>
      <p class="gx-muted">服务器上可执行：<code>bash /opt/ai-forum/scripts/check-tunnel-health.sh</code></p>
      <div class="gx-action-row">
        <button type="button" class="gx-btn gx-btn--secondary" @click="load">重试</button>
        <button
          type="button"
          class="gx-btn gx-btn--ghost"
          @click="
            session.logout();
            router.push('/');
          "
        >
          退出并重新登录
        </button>
        <RouterLink to="/community" class="gx-btn gx-btn--secondary">返回社区</RouterLink>
      </div>
    </div>

    <template v-else-if="overview">
      <GxAdminPageHeader eyebrow="管理概览" title="AI智联平台 · 运营数据" description="可控、可管、可展示的管理闭环" />

      <section class="gx-card gx-admin-hero">
        <div class="gx-stat-grid gx-stat-grid--metrics">
          <div class="gx-stat-card">
            <strong>{{ overview.userCount }}</strong>
            <span>注册用户数</span>
          </div>
          <div class="gx-stat-card">
            <strong>{{ overview.todayPostCount }}</strong>
            <span>今日发帖量</span>
          </div>
          <div class="gx-stat-card">
            <strong>{{ overview.pendingAuditCount }}</strong>
            <span>待审核数</span>
          </div>
          <div class="gx-stat-card">
            <strong>{{ overview.postCount }}</strong>
            <span>公开帖子量</span>
          </div>
        </div>
      </section>

      <section class="gx-card">
        <h3 class="gx-panel__title">板块活跃度</h3>
        <div class="gx-stat-grid gx-stat-grid--boards">
          <article v-for="item in overview.boardActivity" :key="item.boardId" class="gx-stat-card">
            <strong>{{ item.count }}</strong>
            <span>{{ item.name }} · 公开帖</span>
          </article>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.gx-admin-overview {
  width: 100%;
  max-width: none;
}

.gx-admin-hero {
  width: 100%;
}

.gx-admin-hero .gx-stat-grid {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  gap: 16px;
  width: 100%;
}

.gx-admin-hero .gx-stat-card {
  flex: 1 1 0;
  min-width: 0;
}

.gx-stat-grid--boards {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 16px;
  width: 100%;
}

.gx-stat-grid--boards .gx-stat-card {
  flex: 1 1 220px;
}

@media (max-width: 599px) {
  .gx-admin-hero .gx-stat-grid {
    flex-wrap: wrap;
  }

  .gx-admin-hero .gx-stat-card {
    flex: 1 1 calc(50% - 8px);
  }
}
</style>
