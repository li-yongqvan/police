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
          <h3>{{ currentPosts.length }} 篇可展示内容</h3>
        </div>
        <RouterLink to="/community/posts/new" class="secondary-button">发布到本板块</RouterLink>
      </div>

      <div class="post-list">
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
