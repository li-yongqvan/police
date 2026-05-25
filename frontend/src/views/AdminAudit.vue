<script setup>
import { onMounted, ref } from 'vue'
import { adminApi, forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const items = ref([])
const preview = ref(null)
const selected = ref(new Set())
const rejectReason = ref('不符合社区规范')
const batchReason = ref('批量清理违规内容')

async function load() {
  items.value = await adminApi.getPendingAudit()
  selected.value = new Set()
  preview.value = items.value[0] ? await forumApi.getPost(items.value[0].postId) : null
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
  preview.value = await forumApi.getPost(item.postId)
}

onMounted(load)
</script>

<template>
  <div class="page-stack">
    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">内容审核</p>
          <h3>待审核帖子列表</h3>
        </div>
      </div>
      <div class="mw-sticky-toolbar">
        <button class="danger-button" :disabled="!selected.size" @click="batchDelete">
          批量删除 ({{ selected.size }})
        </button>
      </div>

      <div v-if="!items.length" class="empty-state">当前没有待审核内容，社区流转正常。</div>

      <div v-else class="audit-grid">
        <article v-for="item in items" :key="item.id" class="audit-card">
          <label class="audit-select">
            <input type="checkbox" :checked="selected.has(item.postId)" @change="toggleSelect(item.postId)" />
            选中
          </label>
          <h4>{{ item.title }}</h4>
          <p>{{ item.reason || '待人工审核' }}</p>
          <div class="audit-actions">
            <button class="secondary-button" @click="showPreview(item)">查看内容</button>
            <button class="primary-button" @click="act(item, 'approve')">审核通过</button>
            <button class="danger-button" @click="act(item, 'reject')">驳回</button>
          </div>
        </article>
      </div>
    </section>

    <section v-if="preview?.post" class="panel detail-card">
      <p class="eyebrow">审核预览</p>
      <h2>{{ preview.post.title }}</h2>
      <p class="detail-copy">{{ preview.post.content }}</p>
    </section>
  </div>
</template>
