<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import GxBoardHero from '../components/gx/GxBoardHero.vue'
import GxFeedAside from '../components/gx/GxFeedAside.vue'
import GxFeedLayout from '../components/gx/GxFeedLayout.vue'
import GxFeedPagination from '../components/gx/GxFeedPagination.vue'
import GxFeedPostCard from '../components/gx/GxFeedPostCard.vue'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxFeedSortBar from '../components/gx/GxFeedSortBar.vue'
import { GX_NAV_ITEMS, rememberBoardSlug, resolveBoardByKey } from '../composables/useGxNav'
import { useFeedSort, useFeedSortWatch } from '../composables/useFeedSort'
import { forumApi } from '../api'

const route = useRoute()
const { sort } = useFeedSort()
const limit = 10

watch(
  () => route.params.slug,
  (slug) => rememberBoardSlug(slug),
  { immediate: true },
)

const boards = ref([])
const posts = ref([])
const total = ref(0)
const stats = ref(null)
const page = ref(Number(route.query.page) || 1)
const loading = ref(false)

const navItem = computed(() => GX_NAV_ITEMS.find((n) => n.key === route.params.slug))
const board = computed(() => resolveBoardByKey(boards.value, route.params.slug))
const hasBoard = computed(() => !!board.value)

async function load(resetPage = false) {
  if (resetPage) page.value = 1
  loading.value = true
  try {
    boards.value = await forumApi.getBoards()
    stats.value = await forumApi.getCommunityStats()
    const b = board.value
    if (!b) {
      posts.value = []
      total.value = 0
      return
    }
    const feed = await forumApi.getPosts({
      boardId: b.id,
      page: page.value,
      limit,
      sort: sort.value,
    })
    posts.value = feed.posts
    total.value = feed.total
  } finally {
    loading.value = false
  }
}

watch(page, () => load(false))
watch(() => route.params.slug, () => load(true), { immediate: true })
watch(sort, () => load(true))
useFeedSortWatch(() => load(true))

function onPageChange(n) {
  page.value = n
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="gx-page gx-feed-page">
    <GxFeedLayout>
      <template #header>
        <GxBoardHero :board="board" :nav-label="navItem?.label" />
      </template>
      <template #sort>
        <GxFeedSortBar v-model="sort" />
      </template>
      <GxEmptyState
        v-if="!loading && !posts.length && !hasBoard"
        title="板块不存在或已下线"
        description="换个板块看看，或回到首页浏览最新动态"
      />
      <GxEmptyState
        v-else-if="!loading && !posts.length"
        title="该板块暂无帖子"
        description="成为第一个在本板块发帖的同学"
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
        <GxFeedAside variant="board" :stats="stats" :board="board" />
      </template>
    </GxFeedLayout>
  </div>
</template>
