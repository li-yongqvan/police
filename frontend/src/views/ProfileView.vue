<script setup>
import { computed, ref } from 'vue'
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
  department: session.currentUser?.department || '',
  squad: session.currentUser?.squad || '',
  grade: session.currentUser?.grade || '',
})
const avatarFile = ref(null)
const avatarFileName = ref('')
const saving = ref(false)
const uploadingAvatar = ref(false)

const breadcrumbItems = [
  { label: '首页', to: '/community' },
  { label: '个人中心' },
]

const displayName = computed(() => formatAuthorLabel(session.currentUser))
const avatarInitial = computed(() => {
  const raw = session.currentUser?.name || session.currentUser?.username || '?'
  return String(raw).trim()[0]?.toUpperCase() || '?'
})
const profileComplete = computed(() => {
  const fields = [form.value.name, form.value.department, form.value.squad, form.value.grade, form.value.bio]
  const filled = fields.filter((v) => String(v || '').trim()).length
  return Math.round((filled / fields.length) * 100)
})

function onAvatarPick(event) {
  const file = event.target.files?.[0]
  avatarFile.value = file || null
  avatarFileName.value = file?.name || ''
}

async function saveProfile() {
  saving.value = true
  try {
    const user = await userApi.updateProfile(session.currentUser.id, {
      ...form.value,
      profileCompleted: !!(form.value.department && form.value.squad && form.value.grade),
    })
    session.currentUser = user
    localStorage.setItem('ai-forum-user', JSON.stringify(user))
    session.setFlash('资料已保存', 'success')
  } finally {
    saving.value = false
  }
}

async function saveAvatar() {
  if (!avatarFile.value) return
  uploadingAvatar.value = true
  try {
    const user = await userApi.uploadAvatar(session.currentUser.id, avatarFile.value)
    session.currentUser = user
    localStorage.setItem('ai-forum-user', JSON.stringify(user))
    avatarFile.value = null
    avatarFileName.value = ''
    session.setFlash('头像已更新', 'success')
  } finally {
    uploadingAvatar.value = false
  }
}
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
              {{ avatarInitial }}
            </div>
            <h1 class="gx-profile-card__name">{{ displayName }}</h1>
            <p class="gx-profile-card__meta">学号 {{ session.currentUser?.username }}</p>
            <div class="gx-profile-card__badges">
              <Badge variant="secondary">Lv.{{ session.currentUser?.level ?? 1 }}</Badge>
              <Badge variant="outline">{{ session.currentUser?.role || 'student' }}</Badge>
            </div>
            <div v-if="form.department" class="gx-profile-card__chips">
              <span class="gx-stat-chip">{{ form.department }}</span>
              <span class="gx-stat-chip">{{ form.squad }}</span>
              <span class="gx-stat-chip">{{ form.grade }}</span>
            </div>
            <p v-else-if="form.bio" class="gx-profile-card__bio">{{ form.bio }}</p>
            <div class="gx-profile-card__progress">
              <div class="gx-profile-card__progress-head">
                <span>资料完整度</span>
                <strong>{{ profileComplete }}%</strong>
              </div>
              <div class="gx-profile-card__progress-bar" role="progressbar" :aria-valuenow="profileComplete">
                <span :style="{ width: `${profileComplete}%` }" />
              </div>
            </div>
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
              <p class="gx-profile-form-head__desc">完善院系与简介，便于同学识别与交流。</p>
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
              <h3 class="gx-form-section__title">院系信息</h3>
              <div class="gx-form-section__grid gx-form-section__grid--2">
                <div class="gx-field">
                  <Label for="pf-dept">院系</Label>
                  <Input id="pf-dept" v-model="form.department" placeholder="例如：刑事科学技术学院" />
                </div>
                <div class="gx-field">
                  <Label for="pf-squad">区队</Label>
                  <Input id="pf-squad" v-model="form.squad" placeholder="例如：一区队" />
                </div>
                <div class="gx-field gx-field--full">
                  <Label for="pf-grade">年级</Label>
                  <Input id="pf-grade" v-model="form.grade" placeholder="例如：2024" />
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
                  :disabled="!avatarFile || uploadingAvatar"
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
