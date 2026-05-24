<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { forumApi } from '../api'

const route = useRoute()
const boards = ref([])
const posts = ref([])

const currentBoard = computed(() => boards.value.find((item) => item.slug === route.params.slug))
const currentPosts = computed(() =>
  posts.value.filter((item) => item.boardId === currentBoard.value?.id),
)

async function loadBoard() {
  boards.value = await forumApi.getBoards()
  posts.value = await forumApi.getPosts()
}

watch(() => route.params.slug, loadBoard, { immediate: true })
onMounted(loadBoard)
</script>

<template>
  <div class="page-stack page-stack--browse">
    <section class="panel content-panel board-header-compact">
      <div class="section-title section-title--compact">
        <div>
          <p class="eyebrow">板块</p>
          <h3>{{ currentBoard?.name }}</h3>
        </div>
        <RouterLink to="/community/posts/new" class="secondary-button">发帖</RouterLink>
      </div>
      <p class="board-desc-compact">{{ currentBoard?.description }}</p>
    </section>

    <section class="panel content-panel section-feed">
      <div class="section-title section-title--compact">
        <div>
          <p class="eyebrow">刷帖</p>
          <h3>本板块帖子</h3>
        </div>
        <span class="feed-count">{{ currentPosts.length }} 篇</span>
      </div>

      <div class="post-list post-list--feed">
        <RouterLink
          v-for="post in currentPosts"
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
    </section>
  </div>
</template>
