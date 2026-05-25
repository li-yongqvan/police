<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const form = ref({ title: '', content: '' })

onMounted(async () => {
  const data = await forumApi.getPost(route.params.id)
  form.value.title = data.post.title
  form.value.content = data.post.content
})

async function save() {
  await forumApi.updatePost(route.params.id, form.value)
  session.setFlash('帖子已更新。', 'success')
  router.push(`/community/posts/${route.params.id}`)
}
</script>

<template>
  <section class="panel form-panel">
    <p class="eyebrow">编辑帖子</p>
    <h2>更新标题与正文</h2>
    <div class="form-grid">
      <label class="full-span">
        <span>标题</span>
        <input v-model="form.title" type="text" />
      </label>
      <label class="full-span">
        <span>正文</span>
        <textarea v-model="form.content" rows="8" />
      </label>
    </div>
    <button class="primary-button" @click="save">保存修改</button>
  </section>
</template>
