<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { formatApiError } from '../api/errors'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const payload = ref({ post: null, comments: [] })
const commentText = ref('')

const attachments = computed(() => payload.value.post?.attachments ?? [])
const isAuthor = computed(
  () => payload.value.post?.authorId === session.currentUser?.id,
)

async function loadPost() {
  payload.value = await forumApi.getPost(route.params.id)
}

async function likePost() {
  try {
    const resp = await forumApi.likePost(route.params.id)
    await loadPost()
    const count = resp.likeCount ?? payload.value.post?.likeCount
    session.setFlash(`点赞已记录${count != null ? `（${count} 赞）` : ''}。`, 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  }
}

async function collectPost() {
  try {
    const resp = await forumApi.collectPost(route.params.id)
    const msg = resp?.message || resp?.collected ? '已收藏。' : '收藏操作已完成。'
    session.setFlash(msg, 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  }
}

async function submitComment() {
  try {
    await forumApi.createComment(route.params.id, { content: commentText.value })
    commentText.value = ''
    await loadPost()
    session.setFlash('评论已发布。', 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  }
}

async function deletePost() {
  if (!window.confirm('确定删除这篇帖子？')) return
  await forumApi.deletePost(route.params.id)
  session.setFlash('帖子已删除。', 'success')
  router.push('/community')
}

onMounted(loadPost)
</script>

<template>
  <div v-if="payload.post" class="page-stack mw-post-detail">
    <article class="panel detail-card">
      <div class="post-topline">
        <span class="badge">{{ payload.post.isFeatured ? '精华帖' : '公开帖' }}</span>
        <span>{{ payload.post.likeCount }} 赞 · {{ payload.post.commentCount }} 评论</span>
      </div>
      <h2>{{ payload.post.title }}</h2>
      <p class="detail-copy">{{ payload.post.content }}</p>
      <p class="author-line">
        作者：
        <RouterLink :to="`/community/users/${payload.post.authorId}`">
          {{ payload.post.authorName || payload.post.authorId }}
        </RouterLink>
      </p>

      <section v-if="attachments.length" class="attachment-section">
        <h3>资源附件</h3>
        <div class="attachment-list">
          <a
            v-for="item in attachments"
            :key="item.id"
            :href="item.url"
            target="_blank"
            rel="noopener"
            class="attachment-card"
          >
            <strong>{{ item.title }}</strong>
            <span>{{ item.type }}</span>
          </a>
        </div>
      </section>

      <div class="action-row">
        <button class="primary-button" @click="likePost">点赞</button>
        <button class="secondary-button" @click="collectPost">收藏</button>
        <RouterLink
          v-if="isAuthor"
          :to="`/community/posts/${payload.post.id}/edit`"
          class="secondary-button"
        >
          编辑
        </RouterLink>
        <button v-if="isAuthor" class="danger-button" @click="deletePost">删除</button>
      </div>
    </article>

    <section class="panel content-panel">
      <div class="section-title">
        <div>
          <p class="eyebrow">评论区</p>
          <h3>继续把讨论往下走</h3>
        </div>
      </div>
      <div class="comment-list">
        <article v-for="comment in payload.comments" :key="comment.id" class="comment-card">
          <strong>
            <RouterLink :to="`/community/users/${comment.authorId}`">{{ comment.authorName }}</RouterLink>
          </strong>
          <p>{{ comment.content }}</p>
        </article>
      </div>

      <div class="comment-form">
        <textarea v-model="commentText" rows="4" placeholder="补充你的看法、解决方案或追问..." />
        <button class="primary-button" :disabled="!commentText.trim()" @click="submitComment">
          发表评论
        </button>
      </div>
    </section>
  </div>
</template>
