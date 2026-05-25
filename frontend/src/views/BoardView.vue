<script setup>
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { forumApi } from '../api'

const route = useRoute()
const boards = ref([])
const posts = ref([])
const total = ref(0)
const page = ref(1)
const limit = 20
const loadingMore = ref(false)

const currentBoard = computed(() => boards.value.find((item) => item.slug === route.params.slug))

async function loadBoard(reset = true) {
  if (reset) page.value = 1
  boards.value = await forumApi.getBoards()
  const board = boards.value.find((item) => item.slug === route.params.slug)
  if (!board) {
    posts.value = []
    total.value = 0
    return
  }
  const feed = await forumApi.getPosts({ boardId: board.id, page: page.value, limit })
  if (reset) posts.value = feed.posts
  else posts.value = [...posts.value, ...feed.posts]
  total.value = feed.total
}

async function loadMore() {
  if (posts.value.length >= total.value) return
  loadingMore.value = true
  page.value += 1
  try {
    const board = currentBoard.value
    if (!board) return
    const feed = await forumApi.getPosts({ boardId: board.id, page: page.value, limit })
    posts.value = [...posts.value, ...feed.posts]
  } finally {
    loadingMore.value = false
  }
}

watch(() => route.params.slug, () => loadBoard(true), { immediate: true })
</script>

<template>
  <div class="page-stack">
    <section class="panel content-panel">
      <p class="eyebrow">板块详情</p>
      <h2>{{ currentBoard?.name }}</h2>
      <p>{{ currentBoard?.description }}</p>
    </section>

    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">帖子列表</p>
          <h3>{{ posts.length }} / {{ total }} 篇可展示内容</h3>
        </div>
        <RouterLink to="/community/posts/new" class="secondary-button">发布到本板块</RouterLink>
      </div>

      <div class="post-list">
        <RouterLink
          v-for="post in posts"
          :key="post.id"
          :to="`/community/posts/${post.id}`"
          class="post-card"
        >
          <div class="post-topline">
            <span class="badge subtle">{{ post.isFeatured ? '精华' : '讨论' }}</span>
            <span>{{ post.likeCount }} 赞 · {{ post.commentCount }} 评论</span>
          </div>
          <h4>{{ post.title }}</h4>
          <p>{{ post.content }}</p>
        </RouterLink>
      </div>
      <button
        v-if="posts.length < total"
        class="secondary-button load-more-btn"
        :disabled="loadingMore"
        @click="loadMore"
      >
        {{ loadingMore ? '加载中…' : '加载更多' }}
      </button>
    </section>
  </div>
</template>
