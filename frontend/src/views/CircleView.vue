<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxFeedLayout from '../components/gx/GxFeedLayout.vue'
import GxFeedPagination from '../components/gx/GxFeedPagination.vue'
import GxFeedPostCard from '../components/gx/GxFeedPostCard.vue'
import GxFeedSortBar from '../components/gx/GxFeedSortBar.vue'
import GxIcon from '../components/gx/GxIcon.vue'
import { CAMPUS_CIRCLE_SLUG, resolveCircleBoard } from '../composables/useGxNav'
import { useFeedSort, useFeedSortWatch } from '../composables/useFeedSort'
import { forumApi } from '../api'

const route = useRoute()
const { sort } = useFeedSort()
const limit = 10
const boards = ref([])
const posts = ref([])
const total = ref(0)
const page = ref(Number(route.query.page) || 1)
const loading = ref(false)

const circleBoard = computed(() => resolveCircleBoard(boards.value))

const newPostTo = computed(() => ({
  path: '/community/posts/new',
  query: { board: circleBoard.value?.slug || CAMPUS_CIRCLE_SLUG },
}))

async function load(resetPage = false) {
  if (resetPage) page.value = 1
  loading.value = true
  try {
    boards.value = await forumApi.getBoards()
    const b = circleBoard.value
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
watch(sort, () => load(true))
useFeedSortWatch(() => load(true))

onMounted(() => load(true))

function onPageChange(n) {
  page.value = n
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="gx-page gx-circle-page gx-feed-page">
    <GxFeedLayout>
      <template #header>
        <header class="gx-circle-hero">
          <span class="gx-circle-hero__bubble" aria-hidden="true">☀️</span>
          <div class="gx-circle-hero__body">
            <p class="gx-circle-hero__eyebrow">轻松一刻</p>
            <h1 class="gx-circle-hero__title">校园圈</h1>
            <p class="gx-circle-hero__desc">
              {{ circleBoard?.description || '晒晒校园小事，和技术讨论之余的放松空间' }}
            </p>
          </div>
          <RouterLink :to="newPostTo" class="gx-circle-hero__cta">
            <GxIcon name="edit" :size="18" />
            发动态
          </RouterLink>
        </header>
      </template>
      <template #sort>
        <GxFeedSortBar v-model="sort" />
      </template>

      <GxEmptyState
        v-if="!loading && !circleBoard"
        title="校园圈板块未就绪"
        description="请管理员执行数据库迁移后刷新，或于后台创建 slug 为 campus-circle 的板块"
      />

      <GxEmptyState
        v-else-if="!loading && !posts.length"
        title="还没有校园动态"
        description="成为第一个分享校园小事的同学吧"
      >
        <RouterLink :to="newPostTo" class="gx-circle-empty-btn">发第一条校园动态</RouterLink>
      </GxEmptyState>

      <div v-else class="gx-feed-stream gx-circle-feed">
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
        <div class="gx-panel gx-circle-aside-panel">
          <h3 class="gx-panel__title">发帖小贴士</h3>
          <ul class="gx-circle-aside-tips">
            <li>分享宿舍趣事、食堂探店、运动打卡都可以</li>
            <li>请勿泄露他人隐私或涉密信息</li>
            <li>友善互动，让圈子保持轻松氛围</li>
          </ul>
          <RouterLink :to="newPostTo" class="gx-circle-aside-cta">去发动态 →</RouterLink>
        </div>
      </template>
    </GxFeedLayout>
  </div>
</template>
