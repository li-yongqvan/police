<script setup>
import { onMounted, ref } from 'vue'
import { adminApi, userApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const users = ref([])

async function load() {
  users.value = await userApi.listUsers()
}

async function toggleStatus(user) {
  const next = user.status === 'active' ? 'banned' : 'active'
  await adminApi.setUserStatus(user.id, next)
  session.setFlash(`用户状态已切换为 ${next}`, 'success')
  await load()
}

onMounted(load)
</script>

<template>
  <section class="panel content-panel">
    <div class="section-title">
      <div>
        <p class="eyebrow">用户管理</p>
        <h3>演示封禁与恢复</h3>
      </div>
    </div>

    <div class="user-list">
      <article v-for="user in users" :key="user.id" class="user-row">
        <div>
          <strong>{{ user.name }}</strong>
          <p>{{ user.department }}</p>
        </div>
        <div class="user-actions">
          <span :class="['badge', user.status === 'active' ? 'good' : 'bad']">{{ user.status }}</span>
          <button class="secondary-button" @click="toggleStatus(user)">
            {{ user.status === 'active' ? '封禁账号' : '恢复账号' }}
          </button>
        </div>
      </article>
    </div>
  </section>
</template>
