<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxFeedLayout from '../components/gx/GxFeedLayout.vue'
import GxFeedPagination from '../components/gx/GxFeedPagination.vue'
import GxFeedPostCard from '../components/gx/GxFeedPostCard.vue'
import Button from '../components/ui/Button.vue'
import { clearBrowseHistory, getBrowseHistory } from '../composables/useBrowseHistory'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const session = useSessionStore()
const limit = 10
const posts = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const historyItems = ref([])

const mode = computed(() => route.meta.libraryMode || 'posts')

const pageConfig = computed(() => {
  const map = {
    posts: {
      title: '我的帖子',
      description: '你发布过的主题帖',
      emptyTitle: '还没有发过帖',
      emptyDesc: '去板块或首页发布第一条讨论吧',
      cta: { label: '去发帖', to: '/community/posts/new' },
    },
    favorites: {
      title: '我的收藏',
      description: '你收藏过的帖子',
      emptyTitle: '暂无收藏',
      emptyDesc: '浏览帖子时点击收藏，会显示在这里',
      cta: { label: '去逛逛', to: '/community' },
    },
    history: {
      title: '浏览历史',
      description: '最近打开过的帖子（仅保存在本浏览器）',
      emptyTitle: '暂无浏览记录',
      emptyDesc: '打开帖子详情后会自动记录',
      cta: { label: '去首页', to: '/community' },
    },
  }
  return map[mode.value] || map.posts
})

const breadcrumbItems = computed(() => [
  { label: '首页', to: '/community' },
  { label: pageConfig.value.title },
])

const historyPosts = computed(() =>
  historyItems.value.map((item) => ({
    id: item.id,
    title: item.title,
    boardName: item.boardName,
    boardSlug: item.boardSlug,
    authorName: item.authorName,
    content: '',
    likeCount: 0,
    commentCount: 0,
    createdAtIso: item.visitedAt,
    createdAt: '',
  })),
)

const displayPosts = computed(() => (mode.value === 'history' ? historyPosts.value : posts.value))

const displayTotal = computed(() => (mode.value === 'history' ? historyPosts.value.length : total.value))

async function loadFeed(reset = true) {
  if (mode.value === 'history') {
    historyItems.value = getBrowseHistory()
    return
  }
  if (reset) page.value = 1
  loading.value = true
  try {
    if (mode.value === 'favorites') {
      const feed = await forumApi.getMyCollections({ page: page.value, limit })
      posts.value = feed.posts
      total.value = feed.total
    } else {
      const uid = session.currentUser?.id
      if (!uid) {
        posts.value = []
        total.value = 0
        return
      }
      const feed = await forumApi.getPosts({
        authorId: uid,
        page: page.value,
        limit,
        sort: 'new',
      })
      posts.value = feed.posts
      total.value = feed.total
    }
  } finally {
    loading.value = false
  }
}

function onPageChange(n) {
  page.value = n
  loadFeed(false)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function clearHistory() {
  clearBrowseHistory()
  historyItems.value = []
  session.setFlash('浏览历史已清空', 'success')
}

onMounted(() => loadFeed(true))
watch(() => route.name, () => loadFeed(true))
</script>

<template>
  <div class="gx-page gx-feed-page">
    <GxFeedLayout :show-aside="false">
      <template #header>
        <GxBreadcrumb :items="breadcrumbItems" />
        <header class="gx-library-head">
          <div>
            <h1 class="text-display text-gx-primary">{{ pageConfig.title }}</h1>
            <p class="mt-1 text-body text-gx-muted">{{ pageConfig.description }}</p>
          </div>
          <div class="gx-library-head__actions">
            <Button v-if="mode === 'history' && historyItems.length" variant="ghost" @click="clearHistory">
              清空历史
            </Button>
            <RouterLink v-else-if="pageConfig.cta" :to="pageConfig.cta.to">
              <Button>{{ pageConfig.cta.label }}</Button>
            </RouterLink>
            <RouterLink to="/community/profile">
              <Button variant="outline">个人中心</Button>
            </RouterLink>
          </div>
        </header>
      </template>

      <GxEmptyState
        v-if="!loading && !displayPosts.length"
        :title="pageConfig.emptyTitle"
        :description="pageConfig.emptyDesc"
      />
      <div v-else class="gx-feed-stream">
        <GxFeedPostCard
          v-for="post in displayPosts"
          :key="post.id"
          :post="post"
          :pinned="post.isPinned"
          :announce="post.isFeatured"
        />
        <GxFeedPagination
          v-if="mode !== 'history' && displayTotal > limit"
          :page="page"
          :total="displayTotal"
          :limit="limit"
          @update:page="onPageChange"
        />
      </div>
    </GxFeedLayout>
  </div>
</template>
