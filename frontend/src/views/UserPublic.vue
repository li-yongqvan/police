<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import GxAuthorChip from '../components/gx/GxAuthorChip.vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import Card from '../components/ui/Card.vue'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxReadingColumn from '../components/gx/GxReadingColumn.vue'
import { formatApiError } from '../api/errors'
import { userApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const session = useSessionStore()
const profile = ref(null)
const error = ref('')

const breadcrumbItems = [
  { label: '首页', to: '/community' },
  { label: '用户主页' },
]

onMounted(async () => {
  try {
    profile.value = await userApi.getProfile(route.params.id)
  } catch (e) {
    error.value = formatApiError(e, '无法加载用户资料')
  }
})
</script>

<template>
  <section class="gx-page">
    <GxReadingColumn center>
      <GxBreadcrumb :items="breadcrumbItems" />
      <div v-if="error" class="text-center text-body text-gx-muted">{{ error }}</div>
      <Card v-else-if="profile" class="p-6">
        <div class="gx-profile-hero flex flex-col gap-4 sm:flex-row sm:items-start">
          <GxAuthorChip :user="profile" :author-id="profile.id" size="lg" :linkable="false" />
          <div class="min-w-0 flex-1">
            <p class="text-body text-gx-muted">@{{ profile.username }} · Lv.{{ profile.level }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <span v-if="profile.department" class="gx-stat-chip">{{ profile.department }}</span>
              <span v-if="profile.squad" class="gx-stat-chip">{{ profile.squad }}</span>
              <span v-if="profile.grade" class="gx-stat-chip">{{ profile.grade }}</span>
            </div>
            <p v-if="profile.bio" class="mt-4 text-body leading-relaxed">{{ profile.bio }}</p>
            <p v-if="profile.id === session.currentUser?.id" class="mt-4 text-body">
              <router-link class="text-gx-primary hover:underline" to="/community/profile">编辑我的资料</router-link>
            </p>
          </div>
        </div>
        <section class="mt-6 border-t border-gx-border pt-6">
          <h2 class="text-title text-gx-primary">TA 的帖子</h2>
          <GxEmptyState
            class="!py-8"
            title="帖子列表即将上线"
            description="当前版本仅展示公开资料，后续将在此展示该同学发布的帖子。"
          />
        </section>
      </Card>
    </GxReadingColumn>
  </section>
</template>
