<script setup>
import { computed, onMounted, ref } from 'vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { adminApi, userApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const roles = ref([])
const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '角色权限' },
]

const users = ref([])
const selectedUserId = ref('')
const selectedRoleId = ref('')
const userRoles = ref([])
const selectedUser = computed(() => users.value.find((user) => user.id === selectedUserId.value))

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
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="角色权限" title="为用户分配中台角色" />

    <section class="gx-admin-summary">
      <article class="gx-admin-summary__item">
        <span>用户数量</span>
        <strong>{{ users.length }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>角色数量</span>
        <strong>{{ roles.length }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>当前用户角色</span>
        <strong>{{ userRoles.length }}</strong>
      </article>
    </section>

    <section class="gx-card gx-form">
      <h3 class="gx-panel__title">角色分配</h3>
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
      <div class="gx-section-head">
        <strong>已分配角色</strong>
        <span class="gx-muted">{{ selectedUser?.name || '未选择用户' }}</span>
      </div>
      <div v-if="!userRoles.length" class="gx-empty">当前用户还没有分配中台角色。</div>
      <div class="gx-admin-list">
        <article v-for="r in userRoles" :key="r.id || r.role_id" class="gx-admin-row">
          <strong>{{ r.name || r.role_name }}</strong>
          <button type="button" class="gx-btn gx-btn--secondary" @click="remove(r.id || r.role_id)">移除</button>
        </article>
      </div>
    </section>
  </div>
</template>
