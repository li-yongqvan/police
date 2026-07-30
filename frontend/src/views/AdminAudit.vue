<script setup>
import { computed, onMounted, ref } from 'vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const items = ref([])
const preview = ref(null)
const selected = ref(new Set())
const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '内容审核' },
]

const rejectReason = ref('不符合社区规范')
const batchReason = ref('批量清理违规内容')
const selectedCount = computed(() => selected.value.size)
const previewTitle = computed(() => preview.value?.post?.title || '未选择')

async function load() {
  items.value = await adminApi.getPendingAudit()
  selected.value = new Set()
  preview.value = items.value[0] ? await adminApi.getPostDetail(items.value[0].postId) : null
}

function toggleSelect(id) {
  const next = new Set(selected.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  selected.value = next
}

async function act(item, action) {
  if (action === 'approve') {
    await adminApi.approveAudit(item.id)
  } else {
    const reason = window.prompt('请输入驳回理由', rejectReason.value)
    if (reason === null) return
    await adminApi.rejectAudit(item.id, reason || rejectReason.value)
  }
  session.setFlash(`审核操作已完成：${action === 'approve' ? '通过' : '驳回'}`, 'success')
  await load()
}

async function batchDelete() {
  if (!selected.value.size) {
    session.setFlash('请先勾选要删除的帖子', 'info')
    return
  }
  if (!window.confirm(`确定批量删除 ${selected.value.size} 篇待审帖子？`)) return
  await adminApi.batchDeleteAudit([...selected.value], batchReason.value)
  session.setFlash('批量删除已完成', 'success')
  await load()
}

async function showPreview(item) {
  preview.value = await adminApi.getPostDetail(item.postId)
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="内容审核" title="待审核帖子" />

    <section class="gx-admin-summary">
      <article class="gx-admin-summary__item">
        <span>待审核</span>
        <strong>{{ items.length }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>已选中</span>
        <strong>{{ selectedCount }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>当前预览</span>
        <strong>{{ previewTitle }}</strong>
      </article>
    </section>

    <div class="gx-audit-workbench">
      <section class="gx-card">
        <div class="gx-toolbar">
          <strong>审核队列</strong>
          <button type="button" class="gx-btn gx-btn--danger" :disabled="!selected.size" @click="batchDelete">
            批量删除 ({{ selected.size }})
          </button>
        </div>

        <div v-if="!items.length" class="gx-empty">当前没有待审核内容，社区流转正常。</div>

        <div v-else class="gx-audit-grid">
          <article v-for="item in items" :key="item.id" class="gx-admin-row gx-admin-row--stack">
            <label class="gx-muted gx-admin-check">
              <input type="checkbox" :checked="selected.has(item.postId)" @change="toggleSelect(item.postId)" />
              选中
            </label>
            <h4>{{ item.title }}</h4>
            <p class="gx-muted">{{ item.reason || '待人工审核' }}</p>
            <div class="gx-admin-actions">
              <button type="button" class="gx-btn gx-btn--secondary" @click="showPreview(item)">查看内容</button>
              <button type="button" class="gx-btn gx-btn--primary" @click="act(item, 'approve')">审核通过</button>
              <button type="button" class="gx-btn gx-btn--danger" @click="act(item, 'reject')">驳回</button>
            </div>
          </article>
        </div>
      </section>

      <section v-if="preview?.post" class="gx-card gx-audit-preview">
        <p class="gx-eyebrow">审核预览</p>
        <h2>{{ preview.post.title }}</h2>
        <p class="gx-post-body">{{ preview.post.content }}</p>
      </section>
      <section v-else class="gx-card gx-audit-preview gx-muted" style="min-height: 200px">
        <p class="gx-eyebrow">审核预览</p>
        <p>点击「查看内容」在此预览帖子正文。</p>
      </section>
    </div>
  </div>
</template>
