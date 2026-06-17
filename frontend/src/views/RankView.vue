<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import GxEmptyState from '../components/gx/GxEmptyState.vue'
import GxFeedLayout from '../components/gx/GxFeedLayout.vue'
import GxFeedPostCard from '../components/gx/GxFeedPostCard.vue'
import GxIcon from '../components/gx/GxIcon.vue'
import { forumApi } from '../api'
import { loadPage } from '../composables/usePageLoad'
import { boardKeyFromName, CAMPUS_CIRCLE_SLUG } from '../composables/useGxNav'

function boardLink(b) {
  if (b.slug === CAMPUS_CIRCLE_SLUG || /校园圈/.test(b.name || '')) {
    return '/community/circle'
  }
  return `/community/boards/${boardKeyFromName(b.name)}`
}

const boards = ref([])
const posts = ref([])
const stats = ref(null)
const loading = ref(false)

const boardRank = computed(() =>
  [...boards.value].sort((a, b) => (b.postCount ?? 0) - (a.postCount ?? 0)).slice(0, 8),
)

const statCards = computed(() => {
  const s = stats.value
  if (!s) return []
  return [
    { label: '帖子总数', value: s.total_posts ?? '—' },
    { label: '今日发帖', value: s.posts_today ?? '—' },
    { label: '注册用户', value: s.total_users ?? '—' },
    { label: '在线同学', value: s.online_users ?? '—' },
  ]
})

async function load() {
  loading.value = true
  try {
    const data = await loadPage({
      boards: () => forumApi.getBoards(),
      stats: () => forumApi.getCommunityStats(),
      hot: () => forumApi.getPosts({ sort: 'hot', page: 1, limit: 20 }),
    })
    boards.value = data.boards
    stats.value = data.stats
    posts.value = data.hot.posts
  } finally {
    loading.value = false
  }
}

onMounted(() => load())
</script>

<template>
  <div class="gx-page gx-rank-page gx-feed-page">
    <GxFeedLayout>
      <template #header>
        <header class="gx-rank-hero">
          <span class="gx-rank-hero__icon" aria-hidden="true">
            <GxIcon name="star" :size="28" />
          </span>
          <div>
            <p class="gx-rank-hero__eyebrow">社区热度</p>
            <h1 class="gx-rank-hero__title">排行榜</h1>
            <p class="gx-rank-hero__desc">热门帖子与板块活跃度一览</p>
          </div>
        </header>
      </template>

      <div v-if="statCards.length" class="gx-rank-stats">
        <div v-for="card in statCards" :key="card.label" class="gx-rank-stats__item">
          <span class="gx-rank-stats__value">{{ card.value }}</span>
          <span class="gx-rank-stats__label">{{ card.label }}</span>
        </div>
      </div>

      <section v-if="boardRank.length" class="gx-rank-section">
        <h2 class="gx-rank-section__title">板块热度</h2>
        <ol class="gx-rank-board-list">
          <li v-for="(b, index) in boardRank" :key="b.id" class="gx-rank-board-list__item">
            <span class="gx-rank-board-list__rank" :class="{ 'is-top': index < 3 }">{{ index + 1 }}</span>
            <RouterLink :to="boardLink(b)" class="gx-rank-board-list__link">
              <span class="gx-rank-board-list__name">{{ b.name }}</span>
              <span class="gx-rank-board-list__count">{{ b.postCount ?? 0 }} 帖</span>
            </RouterLink>
          </li>
        </ol>
      </section>

      <section class="gx-rank-section">
        <h2 class="gx-rank-section__title">热门帖子</h2>
        <GxEmptyState
          v-if="!loading && !posts.length"
          title="暂无热门帖子"
          description="多发帖、多互动，热度会慢慢升上来"
        />
        <div v-else class="gx-feed-stream">
          <GxFeedPostCard
            v-for="post in posts"
            :key="post.id"
            :post="post"
            :pinned="post.isPinned"
            :announce="post.isFeatured"
          />
        </div>
      </section>

      <template #aside>
        <div class="gx-panel">
          <h3 class="gx-panel__title">关于榜单</h3>
          <p class="gx-rank-aside-note">
            热门帖子按互动与近期活跃度综合排序；板块热度按发帖数量统计。数据每次进入本页时刷新。
          </p>
          <RouterLink to="/community/circle" class="gx-rank-aside-link">去校园圈逛逛 →</RouterLink>
        </div>
      </template>
    </GxFeedLayout>
  </div>
</template>
