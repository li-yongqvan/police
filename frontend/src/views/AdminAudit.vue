<script setup>
import { onMounted, ref } from 'vue'
import { adminApi, forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const items = ref([])
const preview = ref(null)

async function load() {
  items.value = await adminApi.getPendingAudit()
  preview.value = items.value[0] ? await forumApi.getPost(items.value[0].postId, true) : null
}

async function act(item, action) {
  if (action === 'approve') {
    await adminApi.approveAudit(item.id, session.currentUser.id)
  } else {
    await adminApi.rejectAudit(item.id, session.currentUser.id)
  }
  session.setFlash(`审核操作已完成：${action === 'approve' ? '通过' : '驳回'}`, 'success')
  await load()
}

async function showPreview(item) {
  preview.value = await forumApi.getPost(item.postId, true)
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

      <div v-if="!items.length" class="empty-state">当前没有待审核内容，社区流转正常。</div>

      <div v-else class="audit-grid">
        <article v-for="item in items" :key="item.id" class="audit-card">
          <h4>{{ item.title }}</h4>
          <p>{{ item.reason }}</p>
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
