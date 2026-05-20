<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const session = useSessionStore()
const payload = ref({ post: null, comments: [] })
const commentText = ref('')

const attachments = computed(() => payload.value.post?.attachments ?? [])

async function loadPost() {
  payload.value = await forumApi.getPost(route.params.id)
}

async function likePost() {
  await forumApi.likePost(route.params.id)
  await loadPost()
  session.setFlash('点赞已记录。', 'success')
}

async function submitComment() {
  await forumApi.createComment(route.params.id, {
    authorId: session.currentUser.id,
    content: commentText.value,
  })
  commentText.value = ''
  await loadPost()
  session.setFlash('评论已发布。', 'success')
}

onMounted(loadPost)
</script>

<template>
  <div v-if="payload.post" class="page-stack">
    <article class="panel detail-card">
      <div class="post-topline">
        <span class="badge">{{ payload.post.isFeatured ? '精华帖' : '公开帖' }}</span>
        <span>{{ payload.post.likeCount }} 赞 · {{ payload.post.commentCount }} 评论</span>
      </div>
      <h2>{{ payload.post.title }}</h2>
      <p class="detail-copy">{{ payload.post.content }}</p>
      <div class="tag-row">
        <span v-for="tag in payload.post.tags" :key="tag" class="tag">{{ tag }}</span>
      </div>

      <section v-if="attachments.length" class="attachment-section">
        <h3>资源附件</h3>
        <div class="attachment-list">
          <a v-for="item in attachments" :key="item.id" :href="item.url" target="_blank" class="attachment-card">
            <strong>{{ item.name }}</strong>
            <span>{{ item.type }}</span>
          </a>
        </div>
      </section>

      <div class="action-row">
        <button class="primary-button" @click="likePost">点赞这篇帖子</button>
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
          <strong>{{ comment.authorId }}</strong>
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
