<template>
  <div class="profile-page">
    <div v-if="!authStore.isLoggedIn" class="login-prompt">
      <p>请先登录查看个人资料</p>
      <router-link to="/login" class="btn-primary">去登录</router-link>
    </div>
    <div v-else class="profile-content">
      <div class="profile-header">
        <div class="avatar-section">
          <UserAvatar :user="authStore.user" size="lg" @click="triggerAvatarUpload" />
          <input ref="avatarInput" type="file" accept="image/*" @change="handleAvatarUpload" hidden />
          <span class="change-avatar" @click="triggerAvatarUpload">更换头像</span>
        </div>
        <div class="user-info">
          <h2>{{ user.nickname || user.username }}</h2>
          <LevelBadge :level="user.level" />
          <p class="bio">{{ user.bio || '这个人很懒，什么都没写~' }}</p>
        </div>
      </div>
      <div class="edit-section">
        <h3>编辑资料</h3>
        <form @submit.prevent="handleUpdate">
          <div class="form-row">
            <label>用户名</label>
            <input v-model="editForm.username" />
          </div>
          <div class="form-row">
            <label>昵称</label>
            <input v-model="editForm.nickname" />
          </div>
          <div class="form-row">
            <label>签名</label>
            <textarea v-model="editForm.bio" rows="3"></textarea>
          </div>
          <button type="submit" class="btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </form>
      </div>
      <div class="tabs">
        <div class="tab active">我的帖子</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import UserAvatar from '../components/UserAvatar.vue'
import LevelBadge from '../components/LevelBadge.vue'

const authStore = useAuthStore()
const user = computed(() => authStore.user || {})
const saving = ref(false)
const avatarInput = ref(null)

const editForm = reactive({
  username: '',
  nickname: '',
  bio: '',
})

onMounted(() => {
  if (authStore.user) {
    editForm.username = authStore.user.username || ''
    editForm.nickname = authStore.user.nickname || ''
    editForm.bio = authStore.user.bio || ''
  }
})

async function handleUpdate() {
  saving.value = true
  try {
    await authStore.updateProfile(editForm)
  } catch (e) {
    alert(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

function triggerAvatarUpload() {
  avatarInput.value?.click()
}

async function handleAvatarUpload(e) {
  const file = e.target.files[0]
  if (!file) return
  try {
    await authStore.uploadAvatar(file)
  } catch (err) {
    alert('头像上传失败')
  }
}
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
}

.login-prompt {
  text-align: center;
  padding: 3rem;
}

.login-prompt p {
  color: #999;
  margin-bottom: 1rem;
}

.btn-primary {
  display: inline-block;
  padding: 0.5rem 1.5rem;
  background: #4a90d9;
  color: white;
  border-radius: 6px;
  text-decoration: none;
  border: none;
  cursor: pointer;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #eee;
}

.avatar-section {
  text-align: center;
}

.change-avatar {
  display: block;
  font-size: 0.75rem;
  color: #4a90d9;
  cursor: pointer;
  margin-top: 0.5rem;
}

.user-info h2 {
  color: #333;
  margin-bottom: 0.5rem;
}

.bio {
  color: #666;
  margin-top: 0.5rem;
}

.edit-section {
  margin-bottom: 2rem;
}

.edit-section h3 {
  margin-bottom: 1rem;
  color: #333;
}

.form-row {
  margin-bottom: 1rem;
}

.form-row label {
  display: block;
  margin-bottom: 0.25rem;
  color: #666;
  font-size: 0.875rem;
}

.form-row input,
.form-row textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

.tabs {
  display: flex;
  gap: 1rem;
  border-bottom: 1px solid #eee;
}

.tab {
  padding: 0.75rem 1rem;
  cursor: pointer;
  color: #666;
}

.tab.active {
  color: #4a90d9;
  border-bottom: 2px solid #4a90d9;
}

@media (max-width: 768px) {
  .profile-page { padding: 0 0.5rem; }
  .profile-header { flex-direction: column; align-items: flex-start; }
}
</style>
