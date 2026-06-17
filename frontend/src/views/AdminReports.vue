<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const items = ref([])
const statusFilter = ref('pending')
const adminNote = ref('')
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    items.value = await adminApi.getReports(statusFilter.value)
  } finally {
    loading.value = false
  }
}

async function resolve(item, action) {
  const deletePost =
    action === 'resolved' &&
    window.confirm('是否同时删除被举报帖子？点「确定」删除，点「取消」仅标记已处理。')
  await adminApi.resolveReport(item.id, {
    action,
    delete_post: deletePost,
    admin_note: adminNote.value,
  })
  session.setFlash(action === 'resolved' ? '已处理举报' : '已驳回举报', 'success')
  await load()
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page">
    <header class="gx-page-head">
      <p class="gx-eyebrow">举报处理</p>
      <h1>用户举报</h1>
    </header>

    <section class="gx-card">
      <div class="gx-toolbar">
        <label class="gx-field gx-field--inline">
          <span>状态</span>
          <select v-model="statusFilter" @change="load">
            <option value="pending">待处理</option>
            <option value="resolved">已处理</option>
            <option value="dismissed">已驳回</option>
          </select>
        </label>
        <label class="gx-field gx-field--inline" style="flex: 1">
          <span>处理备注（可选）</span>
          <input v-model="adminNote" type="text" placeholder="内部备注" />
        </label>
        <button type="button" class="gx-btn gx-btn--secondary" :disabled="loading" @click="load">刷新</button>
      </div>

      <div v-if="!items.length" class="gx-empty">当前没有{{ statusFilter === 'pending' ? '待处理' : '' }}举报。</div>

      <article v-for="item in items" :key="item.id" class="gx-card" style="margin-top: 1rem">
        <h4>{{ item.postTitle }}</h4>
        <p class="gx-muted">
          举报人：{{ item.reporterName }} · {{ item.createdAt }}
        </p>
        <p>{{ item.reason }}</p>
        <p v-if="item.adminNote" class="gx-muted">备注：{{ item.adminNote }}</p>
        <div v-if="statusFilter === 'pending'" class="gx-admin-actions">
          <RouterLink class="gx-btn gx-btn--secondary" :to="`/community/posts/${item.postId}`" target="_blank">
            查看帖子
          </RouterLink>
          <button type="button" class="gx-btn gx-btn--primary" @click="resolve(item, 'resolved')">确认违规</button>
          <button type="button" class="gx-btn gx-btn--ghost" @click="resolve(item, 'dismissed')">驳回举报</button>
        </div>
      </article>
    </section>
  </div>
</template>
