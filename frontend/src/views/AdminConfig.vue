<script setup>
import { onMounted, ref } from 'vue'
import { adminApi, forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const config = ref(null)
const boards = ref([])

onMounted(async () => {
  config.value = await adminApi.getConfig(session.token)
  boards.value = await forumApi.getBoards(true)
})

async function save() {
  config.value = await adminApi.updateConfig(config.value, session.token)
  session.setFlash('系统配置已保存。', 'success')
}
</script>

<template>
  <section v-if="config" class="panel form-panel">
    <p class="eyebrow">系统配置</p>
    <h2>基础控制能力</h2>
    <div class="config-list">
      <label class="switch-row">
        <span>开启发帖</span>
        <input v-model="config.postingEnabled" type="checkbox" />
      </label>

      <label class="stacked-row">
        <span>审核模式</span>
        <select v-model="config.moderationMode">
          <option value="auto">自动审核</option>
          <option value="manual">人工审核</option>
        </select>
      </label>

      <div class="stacked-row">
        <span>板块开关</span>
        <div class="toggle-grid">
          <label v-for="board in boards" :key="board.id" class="switch-row">
            <span>{{ board.name }}</span>
            <input v-model="config.boardSwitches[board.id]" type="checkbox" />
          </label>
        </div>
      </div>
    </div>
    <button class="primary-button" @click="save">保存配置</button>
  </section>
</template>
