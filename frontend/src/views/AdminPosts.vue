<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const posts = ref([])
const page = ref(1)
const total = ref(0)
const limit = 20

async function load() {
  const data = await adminApi.listPosts(page.value, limit)
  posts.value = data.posts
  total.value = data.total
}

async function changePage(next) {
  page.value = next
  await load()
}

async function toggleFeatured(post) {
  await adminApi.setPostFeatured(post.id, !post.isFeatured)
  session.setFlash(post.isFeatured ? '已取消精华' : '已设为精华', 'success')
  await load()
}

async function togglePinned(post) {
  await adminApi.setPostPinned(post.id, !post.isPinned)
  session.setFlash(post.isPinned ? '已取消置顶' : '已置顶', 'success')
  await load()
}

async function remove(post) {
  if (!window.confirm(`确定删除《${post.title}》？`)) return
  await adminApi.deletePost(post.id)
  session.setFlash('帖子已删除。', 'success')
  await load()
}

const totalPages = () => Math.max(1, Math.ceil(total.value / limit))

onMounted(load)
</script>

<template>
  <section class="panel content-panel">
    <div class="section-title">
      <div>
        <p class="eyebrow">内容管理</p>
        <h3>全站帖子运营 · 第 {{ page }} / {{ totalPages() }} 页</h3>
      </div>
      <div class="audit-actions mw-pagination">
          <button class="secondary-button" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
          <button class="secondary-button" :disabled="page >= totalPages()" @click="changePage(page + 1)">
            下一页
          </button>
        </div>
    </div>
    <div class="post-list compact">
      <article v-for="post in posts" :key="post.id" class="post-card compact mw-admin-post">
        <div>
          <strong>{{ post.title }}</strong>
          <p>{{ post.boardName }} · {{ post.authorName }} · {{ post.status }}</p>
        </div>
        <div class="audit-actions">
          <button class="secondary-button" @click="toggleFeatured(post)">
            {{ post.isFeatured ? '取消精华' : '设精华' }}
          </button>
          <button class="secondary-button" @click="togglePinned(post)">
            {{ post.isPinned ? '取消置顶' : '置顶' }}
          </button>
          <button class="danger-button" @click="remove(post)">删除</button>
        </div>
      </article>
    </div>
  </section>
</template>
