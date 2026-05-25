<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const codes = ref([])
const batchCount = ref(5)
const statusCode = ref('')
const statusResult = ref(null)

async function load() {
  codes.value = await adminApi.listInviteCodes()
}

async function generateOne() {
  const code = await adminApi.generateInviteCode()
  session.setFlash(`已生成邀请码：${code}`, 'success')
  await load()
}

async function generateBatch() {
  const list = await adminApi.generateInviteBatch(batchCount.value)
  session.setFlash(`已批量生成 ${list.length} 个邀请码`, 'success')
  await load()
}

async function voidCode(code) {
  const value = typeof code === 'string' ? code : code.code
  await adminApi.voidInviteCode(value)
  session.setFlash('邀请码已作废', 'success')
  await load()
}

async function queryStatus() {
  const code = statusCode.value.trim()
  if (!code) return
  statusResult.value = await adminApi.getInviteCodeStatus(code)
}

onMounted(load)
</script>

<template>
  <section class="panel content-panel">
    <div class="section-title">
      <div>
        <p class="eyebrow">邀请码</p>
        <h3>注册准入控制</h3>
      </div>
      <div class="audit-actions">
        <button class="primary-button" @click="generateOne">生成 1 个</button>
        <label class="stacked-row inline-batch">
          <input v-model.number="batchCount" type="number" min="1" max="50" />
          <button class="secondary-button" @click="generateBatch">批量生成</button>
        </label>
      </div>
    </div>

    <div class="form-grid invite-status-row">
      <label class="full-span">
        <span>查询邀请码状态</span>
        <input v-model="statusCode" type="text" placeholder="输入邀请码" />
      </label>
      <button class="secondary-button" @click="queryStatus">查询</button>
    </div>
    <p v-if="statusResult" class="status-hint">
      {{ statusResult.code || statusCode }}：
      {{ statusResult.status || statusResult.message || JSON.stringify(statusResult) }}
    </p>

    <div class="user-list">
      <article v-for="item in codes" :key="item.code || item" class="user-row">
        <div>
          <strong>{{ item.code || item }}</strong>
          <p>{{ item.status || '可用' }} · {{ item.used_by ? `已用：${item.used_by}` : '未使用' }}</p>
        </div>
        <button
          v-if="(item.status || 'active') !== 'void'"
          class="secondary-button"
          @click="voidCode(item)"
        >
          作废
        </button>
      </article>
    </div>
  </section>
</template>
