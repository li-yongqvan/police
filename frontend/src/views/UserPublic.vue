<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { formatApiError } from '../api/errors'
import { userApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const session = useSessionStore()
const profile = ref(null)
const error = ref('')

onMounted(async () => {
  try {
    profile.value = await userApi.getProfile(route.params.id)
  } catch (e) {
    error.value = formatApiError(e, '无法加载用户资料')
  }
})
</script>

<template>
  <section class="panel content-panel">
    <p class="eyebrow">用户主页</p>
    <div v-if="error" class="empty-state">{{ error }}</div>
    <div v-else-if="profile">
      <h2>{{ profile.name }}</h2>
      <p>@{{ profile.username }} · Lv.{{ profile.level }}</p>
      <p v-if="profile.bio" class="detail-copy">{{ profile.bio }}</p>
      <p v-if="profile.id === session.currentUser?.id">
        <router-link to="/community/profile">编辑我的资料</router-link>
      </p>
    </div>
  </section>
</template>
