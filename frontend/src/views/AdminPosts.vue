<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const posts = ref([])
const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '内容管理' },
]

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
  <div class="gx-page gx-admin-page">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="内容管理" title="全站帖子运营" :description="`第 ${page} / ${totalPages()} 页`" />

    <section class="gx-card">
      <div class="gx-section-head">
        <span class="gx-muted">{{ total }} 篇</span>
        <div class="gx-admin-actions">
          <button type="button" class="gx-btn gx-btn--secondary" :disabled="page <= 1" @click="changePage(page - 1)">
            上一页
          </button>
          <button
            type="button"
            class="gx-btn gx-btn--secondary"
            :disabled="page >= totalPages()"
            @click="changePage(page + 1)"
          >
            下一页
          </button>
        </div>
      </div>

      <div class="gx-admin-list">
        <article v-for="post in posts" :key="post.id" class="gx-admin-row">
          <div>
            <strong>{{ post.title }}</strong>
            <p class="gx-muted">{{ post.boardName }} · {{ post.authorName }} · {{ post.status }}</p>
          </div>
          <div class="gx-admin-actions">
            <button type="button" class="gx-btn gx-btn--secondary" @click="toggleFeatured(post)">
              {{ post.isFeatured ? '取消精华' : '设精华' }}
            </button>
            <button type="button" class="gx-btn gx-btn--secondary" @click="togglePinned(post)">
              {{ post.isPinned ? '取消置顶' : '置顶' }}
            </button>
            <button type="button" class="gx-btn gx-btn--danger" @click="remove(post)">删除</button>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>
