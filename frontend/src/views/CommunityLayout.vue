<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import GxAdminFab from '../components/gx/GxAdminFab.vue'
import GxMobileTabBar from '../components/gx/GxMobileTabBar.vue'
import GxSiteFooter from '../components/gx/GxSiteFooter.vue'
import GxSiteHeader from '../components/gx/GxSiteHeader.vue'
import GxSiteSidebar from '../components/gx/GxSiteSidebar.vue'
import Button from '../components/ui/Button.vue'
import Card from '../components/ui/Card.vue'
import Input from '../components/ui/Input.vue'
import Label from '../components/ui/Label.vue'
import { useDrawerNav } from '../composables/useDrawerNav'
import { useSessionStore } from '../stores/session'
import { forumApi, userApi } from '../api'

const route = useRoute()
const session = useSessionStore()
const { drawerOpen, toggleDrawer, closeDrawer } = useDrawerNav()

const isFeedShell = computed(() =>
  ['community-home', 'board', 'rank', 'campus-circle', 'my-posts', 'my-favorites', 'my-history'].includes(
    route.name,
  ),
)
const showOnboarding = ref(false)
const onboard = ref({ department: '', squad: '', grade: '' })
const boards = ref([])
const unreadCount = ref(0)

onMounted(async () => {
  if (session.needsOnboarding) showOnboarding.value = true
  try {
    boards.value = await forumApi.getBoards()
  } catch {
    boards.value = []
  }
  if (session.currentUser) {
    try {
      const { items } = await forumApi.listNotifications(1, 50)
      unreadCount.value = items.filter((i) => !i.is_read).length
    } catch {
      unreadCount.value = 0
    }
  }
})

async function saveOnboarding() {
  await userApi.updateProfile(session.currentUser.id, {
    name: session.currentUser.name,
    department: onboard.value.department,
    squad: onboard.value.squad,
    grade: onboard.value.grade,
    profileCompleted: true,
  })
  await session.refreshMe()
  showOnboarding.value = false
  session.setFlash('资料已完善', 'success')
}
</script>

<template>
  <div class="gx-app gx-community-shell gx-app--with-tabbar">
    <GxSiteHeader
      :drawer-open="drawerOpen"
      :unread-count="unreadCount"
      @toggle-drawer="toggleDrawer"
    />
    <div class="gx-drawer-backdrop" :class="{ 'is-open': drawerOpen }" @click="closeDrawer" />
    <GxSiteSidebar :open="drawerOpen" :boards="boards" @navigate="closeDrawer" />

    <main class="gx-main" :class="{ 'gx-main--feed': isFeedShell }">
      <div class="gx-page-frame" :class="{ 'gx-page-frame--feed': isFeedShell }">
        <RouterView />
      </div>
      <GxSiteFooter v-if="!isFeedShell" />
    </main>
    <GxMobileTabBar />
    <GxAdminFab />

    <div v-if="showOnboarding" class="gx-modal-backdrop">
      <Card class="gx-modal mx-4 max-w-md p-6">
        <h2 class="text-title text-gx-primary">完善资料</h2>
        <p class="mt-1 text-meta text-gx-muted">请填写院系、区队、年级。</p>
        <div class="mt-4 space-y-4">
          <div class="space-y-2">
            <Label for="ob-dept">院系</Label>
            <Input id="ob-dept" v-model="onboard.department" />
          </div>
          <div class="space-y-2">
            <Label for="ob-squad">区队</Label>
            <Input id="ob-squad" v-model="onboard.squad" />
          </div>
          <div class="space-y-2">
            <Label for="ob-grade">年级</Label>
            <Input id="ob-grade" v-model="onboard.grade" />
          </div>
          <Button type="button" class="w-full" @click="saveOnboarding">保存并进入社区</Button>
        </div>
      </Card>
    </div>
  </div>
</template>
