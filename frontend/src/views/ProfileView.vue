<script setup>
import { ref } from 'vue'
import { userApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const form = ref({
  name: session.currentUser?.name || '',
  bio: session.currentUser?.bio || '',
  username: session.currentUser?.username || '',
})
const avatarFile = ref(null)
const saving = ref(false)

async function saveProfile() {
  saving.value = true
  try {
    const user = await userApi.updateProfile(session.currentUser.id, form.value)
    session.currentUser = user
    localStorage.setItem('ai-forum-user', JSON.stringify(user))
    session.setFlash('资料已保存。', 'success')
  } finally {
    saving.value = false
  }
}

async function saveAvatar() {
  if (!avatarFile.value) return
  const user = await userApi.uploadAvatar(session.currentUser.id, avatarFile.value)
  session.currentUser = user
  localStorage.setItem('ai-forum-user', JSON.stringify(user))
  session.setFlash('头像已更新。', 'success')
}
</script>

<template>
  <section class="panel profile-card form-panel">
    <p class="eyebrow">个人主页</p>
    <h2>{{ session.currentUser?.name }}</h2>

    <div class="form-grid">
      <label>
        <span>昵称</span>
        <input v-model="form.name" type="text" />
      </label>
      <label>
        <span>用户名</span>
        <input v-model="form.username" type="text" />
      </label>
      <label class="full-span">
        <span>简介</span>
        <textarea v-model="form.bio" rows="4" />
      </label>
      <label class="full-span">
        <span>头像</span>
        <input type="file" accept="image/*" @change="(e) => (avatarFile = e.target.files?.[0])" />
      </label>
    </div>

    <div class="role-login-actions">
      <button class="primary-button" :disabled="saving" @click="saveProfile">保存资料</button>
      <button class="secondary-button" :disabled="!avatarFile" @click="saveAvatar">上传头像</button>
    </div>

    <div class="profile-grid">
      <div>
        <strong>{{ session.currentUser?.role }}</strong>
        <span>当前角色</span>
      </div>
      <div>
        <strong>{{ session.currentUser?.status }}</strong>
        <span>账号状态</span>
      </div>
      <div>
        <strong>Lv.{{ session.currentUser?.level ?? 1 }}</strong>
        <span>用户等级</span>
      </div>
    </div>
  </section>
</template>
