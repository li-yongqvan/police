<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxFeedAside from '../components/gx/GxFeedAside.vue'
import GxFeedLayout from '../components/gx/GxFeedLayout.vue'
import GxFeedPagination from '../components/gx/GxFeedPagination.vue'
import GxFeedPostCard from '../components/gx/GxFeedPostCard.vue'
import GxHomeHero from '../components/gx/GxHomeHero.vue'
import GxCarousel from '../components/gx/GxCarousel.vue'
import GxFeedSortBar from '../components/gx/GxFeedSortBar.vue'
import { forumApi } from '../api'
import { loadPage } from '../composables/usePageLoad'
import { useFeedSort, useFeedSortWatch } from '../composables/useFeedSort'

const route = useRoute()
const { sort } = useFeedSort()
const limit = 10
const boards = ref([])
const posts = ref([])
const total = ref(0)
const stats = ref(null)
const page = ref(Number(route.query.page) || 1)
const loading = ref(false)
const searchQ = ref(route.query.q || '')

async function load(reset = true) {
  if (reset) page.value = 1
  loading.value = true
  try {
    const loaders = {
      boards: () => forumApi.getBoards(),
      feed: () =>
        forumApi.getPosts({
          page: page.value,
          limit,
          q: searchQ.value || '',
          sort: sort.value,
        }),
      stats: () => forumApi.getCommunityStats(),
    }
    const data = await loadPage(loaders)
    boards.value = data.boards
    stats.value = data.stats
    posts.value = data.feed.posts
    total.value = data.feed.total
  } finally {
    loading.value = false
  }
}

watch(page, () => {
  if (!searchQ.value) load(false)
})

watch(sort, () => {
  if (!searchQ.value) load(true)
})

watch(
  () => route.query.q,
  (q) => {
    searchQ.value = q || ''
    load(true)
  },
)

useFeedSortWatch(() => load(true))

onMounted(() => load(true))

function onPageChange(n) {
  page.value = n
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="gx-page gx-home-page gx-feed-page">
    <div v-if="searchQ" class="mb-4 rounded-gx-md border border-gx-primary/15 bg-gx-primary/5 px-4 py-3 text-body text-gx-primary">
      搜索「{{ searchQ }}」共 {{ total }} 条结果
    </div>

    <GxFeedLayout v-else :show-aside="true">
      <template #header>
        <GxCarousel class="mb-4" />
        <GxHomeHero :stats="stats" />
      </template>
      <template #sort>
        <GxFeedSortBar v-model="sort" />
      </template>
      <GxEmptyState
        v-if="!loading && !posts.length"
        title="暂无帖子"
        description="成为第一个发帖的同学吧"
      />
      <div v-else class="gx-feed-stream">
        <GxFeedPostCard
          v-for="post in posts"
          :key="post.id"
          :post="post"
          :pinned="post.isPinned"
          :announce="post.isFeatured"
        />
        <GxFeedPagination
          v-if="total > limit"
          :page="page"
          :total="total"
          :limit="limit"
          @update:page="onPageChange"
        />
      </div>
      <template #aside>
        <GxFeedAside variant="home" :stats="stats" :boards="boards" />
      </template>
    </GxFeedLayout>

    <template v-if="searchQ">
      <GxEmptyState
        v-if="!posts.length"
        title="未找到相关帖子"
        description="换个关键词试试"
      />
      <div v-else class="gx-feed-stream gx-feed-page">
        <GxFeedPostCard v-for="post in posts" :key="post.id" :post="post" />
        <GxFeedPagination
          v-if="total > limit"
          :page="page"
          :total="total"
          :limit="limit"
          @update:page="onPageChange"
        />
      </div>
    </template>
  </div>
</template>
