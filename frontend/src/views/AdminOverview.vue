<script setup>
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'

const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '管理概览' },
]
</script>

<template>
  <div class="gx-page gx-admin-page gx-admin-overview">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="管理概览" title="AI智联平台 · 运营数据" description="可控、可管、可展示的管理闭环" />

    <section class="gx-card gx-admin-placeholder">
      <div class="gx-placeholder-content">
        <p class="gx-placeholder-icon">??</p>
        <h3>统计功能即将上线</h3>
        <p class="gx-muted">论坛已迁移至 Discourse，新的统计面板正在建设中，敬请期待。</p>
      </div>
    </section>
  </div>
</template>

<style scoped>
.gx-admin-overview {
  width: 100%;
  max-width: none;
}

.gx-admin-placeholder {
  text-align: center;
  padding: 48px 24px;
}

.gx-placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.gx-placeholder-icon {
  font-size: 48px;
  margin: 0;
  line-height: 1;
}

.gx-placeholder-content h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}
</style>