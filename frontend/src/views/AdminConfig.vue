<script setup>
import { onMounted, ref } from 'vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { adminApi, forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const config = ref(null)
const boards = ref([])

const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '系统配置' },
]

onMounted(async () => {
  config.value = await adminApi.getConfig()
  boards.value = await forumApi.getBoards(true)
})

async function save() {
  config.value = await adminApi.updateConfig(config.value)
  window.dispatchEvent(new CustomEvent('forum-config-updated'))
  session.setFlash('系统配置已保存。', 'success')
}
</script>

<template>
  <div v-if="config" class="gx-page gx-admin-page">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="系统配置" title="基础控制能力" />

    <section class="gx-card gx-form">
      <label class="gx-admin-row">
        <span>开启发帖</span>
        <input v-model="config.postingEnabled" type="checkbox" />
      </label>

      <label>
        <span>审核模式</span>
        <select v-model="config.moderationMode">
          <option value="auto">自动审核</option>
          <option value="manual">人工审核</option>
        </select>
      </label>

      <div>
        <span class="gx-panel__title">板块开关</span>
        <label v-for="board in boards" :key="board.id" class="gx-admin-row">
          <span>{{ board.name }}</span>
          <input v-model="config.boardSwitches[board.id]" type="checkbox" />
        </label>
      </div>

      <button type="button" class="gx-btn gx-btn--primary" @click="save">保存配置</button>
    </section>
  </div>
</template>
