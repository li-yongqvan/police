<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const canManageInvites = computed(() => session.currentUser?.role === 'platform_admin')
const discourseAdminUrl = 'http://122.51.233.225:8080/admin'

const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '管理首页' },
]
</script>

<template>
  <div class="gx-page gx-admin-page gx-admin-overview">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader
      eyebrow="管理首页"
      title="AI 智联论坛管理"
      description="当前仅开放试运行必需能力：用户管理和 Discourse 内容治理。"
    />

    <section class="gx-admin-actions-panel">
      <RouterLink class="gx-admin-entry" to="/admin/users">
        <strong>用户管理</strong>
        <span>查看用户、禁用异常账号、处理登录问题。</span>
      </RouterLink>

      <RouterLink v-if="canManageInvites" class="gx-admin-entry" to="/admin/invites">
        <strong>邀请码</strong>
        <span>仅在需要开放新用户注册时生成或作废邀请码。</span>
      </RouterLink>

      <a
        class="gx-admin-entry"
        :href="discourseAdminUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        <strong>Discourse 管理</strong>
        <span>处理帖子、回复、用户举报和论坛内容治理。</span>
      </a>
    </section>
  </div>
</template>

<style scoped>
.gx-admin-overview {
  width: 100%;
  max-width: none;
}

.gx-admin-actions-panel {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.gx-admin-entry {
  display: flex;
  min-height: 112px;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  border: 1px solid rgba(15, 23, 42, 0.1);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.86);
  padding: 18px;
  color: inherit;
  text-decoration: none;
}

.gx-admin-entry:hover {
  border-color: rgba(37, 99, 235, 0.28);
  background: rgba(239, 246, 255, 0.96);
}

.gx-admin-entry strong {
  font-size: 1rem;
}

.gx-admin-entry span {
  color: #475569;
  font-size: 0.875rem;
  line-height: 1.6;
}
</style>
