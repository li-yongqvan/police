<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import Badge from '../components/ui/Badge.vue'
import Button from '../components/ui/Button.vue'
import Card from '../components/ui/Card.vue'
import Input from '../components/ui/Input.vue'
import Label from '../components/ui/Label.vue'
import Textarea from '../components/ui/Textarea.vue'
import { userApi } from '../api'
import { useSessionStore } from '../stores/session'
import { formatAuthorLabel } from '../utils/displayName'

const session = useSessionStore()
const form = ref({
  name: session.currentUser?.name || '',
  bio: session.currentUser?.bio || '',
  username: session.currentUser?.username || '',
})
const avatarFile = ref(null)
const avatarFileName = ref('')
const avatarPreviewUrl = ref('')
const avatarInput = ref(null)
const saving = ref(false)
const uploadingAvatar = ref(false)

const breadcrumbItems = [
  { label: '首页', to: '/community' },
  { label: '个人中心' },
]

const displayName = computed(() => formatAuthorLabel(session.currentUser))
const visibleAvatar = computed(() => avatarPreviewUrl.value || session.currentUser?.avatar || '')
const avatarInitial = computed(() => {
  const raw = session.currentUser?.name || session.currentUser?.username || '?'
  return String(raw).trim()[0]?.toUpperCase() || '?'
})
function clearAvatarPreview() {
  if (avatarPreviewUrl.value) {
    URL.revokeObjectURL(avatarPreviewUrl.value)
    avatarPreviewUrl.value = ''
  }
}

const AVATAR_MAX_PX = 400

async function compressImage(file) {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => {
      let { width, height } = img
      if (width <= AVATAR_MAX_PX && height <= AVATAR_MAX_PX) {
        resolve(file)
        return
      }
      const ratio = Math.min(AVATAR_MAX_PX / width, AVATAR_MAX_PX / height)
      width = Math.round(width * ratio)
      height = Math.round(height * ratio)
      const canvas = document.createElement("canvas")
      canvas.width = width
      canvas.height = height
      const ctx = canvas.getContext("2d")
      ctx.drawImage(img, 0, 0, width, height)
      canvas.toBlob((blob) => {
        if (!blob) { resolve(file); return }
        const compressed = new File([blob], file.name, { type: "image/jpeg", lastModified: Date.now() })
        resolve(compressed)
      }, "image/jpeg", 0.85)
    }
    img.onerror = () => resolve(file)
    img.src = URL.createObjectURL(file)
  })
}

async function onAvatarPick(event) {
  const file = event.target.files?.[0]
  clearAvatarPreview()
  if (file) {
    const compressed = await compressImage(file)
    avatarFile.value = compressed || file
    avatarFileName.value = file.name
    avatarPreviewUrl.value = URL.createObjectURL(file)
  } else {
    avatarFile.value = null
    avatarFileName.value = ''
  }
}

function clearAvatarFile() {
  avatarFile.value = null
  avatarFileName.value = ''
  if (avatarInput.value) {
    avatarInput.value.value = ''
  }
}

function persistCurrentUser(user) {
  session.currentUser = user
  localStorage.setItem('ai-forum-user', JSON.stringify(user))
}

async function saveProfile() {
  saving.value = true
  uploadingAvatar.value = !!avatarFile.value
  try {
    let user = await userApi.updateProfile(session.currentUser.id, {
      name: form.value.name,
      bio: form.value.bio,
      username: form.value.username,
      profileCompleted: true,
    })
    if (avatarFile.value) {
      user = await userApi.uploadAvatar(user.id, avatarFile.value)
      clearAvatarFile()
      clearAvatarPreview()
    }
    persistCurrentUser(user)
    session.setFlash('资料已保存', 'success')
  } finally {
    saving.value = false
    uploadingAvatar.value = false
  }
}

async function saveAvatar() {
  if (!avatarFile.value) return
  uploadingAvatar.value = true
  try {
    const user = await userApi.uploadAvatar(session.currentUser.id, avatarFile.value)
    persistCurrentUser(user)
    clearAvatarFile()
    clearAvatarPreview()
    session.setFlash('头像已更新', 'success')
  } finally {
    uploadingAvatar.value = false
  }
}

onBeforeUnmount(clearAvatarPreview)

// Follow stats
const followCounts = ref({ following: 0, followers: 0 })

async function fetchFollowCounts() {
  try {
    const id = session.currentUser?.id
    if (!id) return
    const result = await userApi.getFollowCounts(id)
    followCounts.value = { following: result.following || 0, followers: result.followers || 0 }
  } catch (_) {}
}

onMounted(() => { fetchFollowCounts() })

</script>

<template>
  <div class="gx-page gx-profile-page">
    <GxBreadcrumb :items="breadcrumbItems" />


    <div class="gx-profile-layout">
      <aside class="gx-profile-sidebar">
        <Card class="gx-profile-card">
          <div class="gx-profile-card__banner" aria-hidden="true" />
          <div class="gx-profile-card__body">
            <div class="gx-profile-avatar gx-profile-avatar--page" aria-hidden="true">
              <img
                v-if="visibleAvatar"
                class="gx-profile-avatar__image"
                :src="visibleAvatar"
                :alt="`${displayName} 的头像`"
              />
              <span v-else>{{ avatarInitial }}</span>
            </div>
            <h1 class="gx-profile-card__name">{{ displayName }}</h1>
            <p class="gx-profile-card__meta">学号 {{ session.currentUser?.username }}</p>
            <div class="gx-profile-card__badges">
              <Badge variant="secondary">Lv.{{ session.currentUser?.level ?? 1 }}</Badge>
              <Badge variant="outline">{{ session.currentUser?.role || 'student' }}</Badge>
            </div>
            <div class="gx-profile-stats">
              <RouterLink to="/community/profile?tab=following" class="gx-profile-stat">
                <span class="gx-profile-stat__num">{{ followCounts.following }}</span>
                <span class="gx-profile-stat__label">关注</span>
              </RouterLink>
              <RouterLink to="/community/profile?tab=followers" class="gx-profile-stat">
                <span class="gx-profile-stat__num">{{ followCounts.followers }}</span>
                <span class="gx-profile-stat__label">粉丝</span>
              </RouterLink>
            </div>
            <p v-if="form.bio" class="gx-profile-card__bio">{{ form.bio }}</p>
            <RouterLink
              v-if="session.currentUser?.id"
              :to="`/community/users/${session.currentUser.id}`"
              class="gx-profile-card__link"
            >
              查看公开主页
            </RouterLink>
            <nav class="gx-profile-card__nav" aria-label="我的内容">
              <RouterLink to="/community/my/posts">我的帖子</RouterLink>
              <RouterLink to="/community/my/favorites">我的收藏</RouterLink>
              <RouterLink to="/community/my/history">浏览历史</RouterLink>
            </nav>
          </div>
        </Card>
      </aside>

      <main class="gx-profile-main">
        <Card class="gx-profile-form-card">
          <header class="gx-profile-form-head">
            <div>
              <p class="gx-profile-form-head__eyebrow">账户设置</p>
              <h2 class="gx-profile-form-head__title">编辑资料</h2>
              <p class="gx-profile-form-head__desc">设置头像、昵称和简介，展示你在论坛里的公开形象。</p>
            </div>
          </header>

          <form class="gx-profile-form" @submit.prevent="saveProfile">
            <section class="gx-form-section">
              <h3 class="gx-form-section__title">基本信息</h3>
              <div class="gx-form-section__grid gx-form-section__grid--2">
                <div class="gx-field">
                  <Label for="pf-name">昵称</Label>
                  <Input id="pf-name" v-model="form.name" placeholder="展示给其他同学的名称" />
                </div>
                <div class="gx-field">
                  <Label for="pf-user">用户名</Label>
                  <Input id="pf-user" v-model="form.username" placeholder="登录用用户名" />
                </div>
              </div>
            </section>

            <section class="gx-form-section">
              <h3 class="gx-form-section__title">个人简介</h3>
              <div class="gx-field">
                <Label for="pf-bio">简介</Label>
                <Textarea
                  id="pf-bio"
                  v-model="form.bio"
                  :rows="4"
                  placeholder="简单介绍自己，如兴趣方向、参与社团等"
                />
              </div>
            </section>

            <section class="gx-form-section gx-form-section--last">
              <h3 class="gx-form-section__title">头像</h3>
              <div class="gx-file-upload">
                <label class="gx-file-upload__btn" for="pf-avatar">选择图片</label>
                <input
                  id="pf-avatar"
                  ref="avatarInput"
                  type="file"
                  accept="image/*"
                  class="gx-file-upload__input"
                  @change="onAvatarPick"
                />
                <span class="gx-file-upload__hint">
                  {{ avatarFileName || '支持 JPG、PNG，建议正方形图片' }}
                </span>
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  :disabled="!avatarFile || uploadingAvatar || saving"
                  @click="saveAvatar"
                >
                  {{ uploadingAvatar ? '上传中…' : '上传头像' }}
                </Button>
              </div>
            </section>

            <footer class="gx-profile-form__footer">
              <Button type="submit" :disabled="saving">
                {{ saving ? '保存中…' : '保存资料' }}
              </Button>
            </footer>
          </form>
        </Card>
      </main>
    </div>
  </div>
</template>
