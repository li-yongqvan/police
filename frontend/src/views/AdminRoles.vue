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
  <section class="panel form-panel">
    <p class="eyebrow">角色权限</p>
    <h2>为用户分配中台角色</h2>
    <div class="form-grid">
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
    </div>
    <button class="primary-button" @click="assign">分配角色</button>

    <div class="user-list" style="margin-top: 22px">
      <article v-for="r in userRoles" :key="r.id || r.role_id" class="user-row">
        <strong>{{ r.name || r.role_name }}</strong>
        <button class="secondary-button" @click="remove(r.id || r.role_id)">移除</button>
      </article>
    </div>
  </section>
</template>
