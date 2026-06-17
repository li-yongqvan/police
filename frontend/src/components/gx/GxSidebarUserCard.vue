<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import GxAvatarInitial from './GxAvatarInitial.vue'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  postCount: { type: Number, default: 0 },
  replyCount: { type: Number, default: 0 },
  favoriteCount: { type: Number, default: 0 },
})

const session = useSessionStore()
const user = computed(() => session.currentUser)
const level = computed(() => user.value?.level ?? 1)
const xpCurrent = computed(() => Math.min(2000, 320 + level.value * 240))
const xpMax = 2000
const xpPct = computed(() => Math.round((xpCurrent.value / xpMax) * 100))
</script>

<template>
  <div v-if="user" class="gx-sidebar-user">
    <RouterLink to="/community/profile" class="gx-sidebar-user__head">
      <GxAvatarInitial :name="user.name" :size="40" />
      <div class="gx-sidebar-user__who">
        <strong>{{ user.name || user.username }}</strong>
        <span class="gx-sidebar-user__level">Lv.{{ level }}</span>
      </div>
    </RouterLink>
    <div class="gx-sidebar-user__xp" aria-hidden="true">
      <div class="gx-sidebar-user__xp-bar" :style="{ width: `${xpPct}%` }" />
    </div>
    <p class="gx-sidebar-user__xp-label text-caption text-gx-muted">
      {{ xpCurrent }} / {{ xpMax }} 经验
    </p>
    <dl class="gx-sidebar-user__stats">
      <div>
        <dt>发帖</dt>
        <dd>{{ postCount }}</dd>
      </div>
      <div>
        <dt>回帖</dt>
        <dd>{{ replyCount }}</dd>
      </div>
      <div>
        <dt>收藏</dt>
        <dd>{{ favoriteCount }}</dd>
      </div>
    </dl>
  </div>
</template>
