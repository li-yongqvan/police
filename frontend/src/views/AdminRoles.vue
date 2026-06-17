<script setup>
import { onMounted, ref } from 'vue'
import { adminApi, userApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const roles = ref([])
const users = ref([])
const selectedUserId = ref('')
const selectedRoleId = ref('')
const userRoles = ref([])

async function load() {
  roles.value = await adminApi.listRoles()
  const data = await userApi.listUsers(1, 100)
  users.value = data.users
  if (!selectedUserId.value && users.value[0]) {
    selectedUserId.value = users.value[0].id
    await loadUserRoles()
  }
}

async function loadUserRoles() {
  if (!selectedUserId.value) return
  userRoles.value = await adminApi.getUserRoles(selectedUserId.value)
}

async function assign() {
  await adminApi.assignRole(selectedUserId.value, selectedRoleId.value)
  session.setFlash('角色已分配', 'success')
  await loadUserRoles()
}

async function remove(roleId) {
  await adminApi.removeRole(selectedUserId.value, roleId)
  session.setFlash('角色已移除', 'success')
  await loadUserRoles()
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page">
    <header class="gx-page-head">
      <p class="gx-eyebrow">角色权限</p>
      <h1>为用户分配中台角色</h1>
    </header>

    <section class="gx-card gx-form">
      <label>
        <span>用户</span>
        <select v-model="selectedUserId" @change="loadUserRoles">
          <option v-for="u in users" :key="u.id" :value="u.id">{{ u.name }} ({{ u.username }})</option>
        </select>
      </label>
      <label>
        <span>角色</span>
        <select v-model="selectedRoleId">
          <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</option>
        </select>
      </label>
      <button type="button" class="gx-btn gx-btn--primary" @click="assign">分配角色</button>
    </section>

    <section class="gx-card">
      <div class="gx-admin-list">
        <article v-for="r in userRoles" :key="r.id || r.role_id" class="gx-admin-row">
          <strong>{{ r.name || r.role_name }}</strong>
          <button type="button" class="gx-btn gx-btn--secondary" @click="remove(r.id || r.role_id)">移除</button>
        </article>
      </div>
    </section>
  </div>
</template>
