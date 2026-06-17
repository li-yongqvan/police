<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import GxComposeShell from '../components/gx/GxComposeShell.vue'
import GxIcon from '../components/gx/GxIcon.vue'
import Input from '../components/ui/Input.vue'
import Textarea from '../components/ui/Textarea.vue'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const saving = ref(false)
const form = ref({ title: '', content: '' })

const editTips = [
  '仅修改必要内容，避免重复发帖。',
  '保存后将在原帖位置更新展示。',
  '重大变更可能再次进入审核。',
]

onMounted(async () => {
  const data = await forumApi.getPost(route.params.id)
  form.value.title = data.post.title
  form.value.content = data.post.content
})

async function save() {
  saving.value = true
  try {
    await forumApi.updatePost(route.params.id, form.value)
    session.setFlash('帖子已更新。', 'success')
    router.push(`/community/posts/${route.params.id}`)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <GxComposeShell panel-title="编辑帖子" panel-subtitle="更新标题与正文内容">
    <template #rail>
      <header class="gx-compose-rail__head">
        <div class="gx-compose-rail__brand">
          <span class="gx-compose-rail__logo gx-compose-rail__logo--edit" aria-hidden="true">
            <GxIcon name="edit" :size="20" />
          </span>
          <div>
            <h1 class="gx-compose-rail__title">编辑帖子</h1>
            <p class="gx-compose-rail__meta">修改后保存即可生效</p>
          </div>
        </div>
      </header>

      <div class="gx-compose-rail__tips">
        <p class="gx-compose-rail__section">编辑提示</p>
        <ul class="gx-compose-tips">
          <li v-for="tip in editTips" :key="tip">{{ tip }}</li>
        </ul>
      </div>

      <RouterLink :to="`/community/posts/${route.params.id}`" class="gx-compose-rail__back">
        <GxIcon name="home" :size="16" />
        返回帖子详情
      </RouterLink>
    </template>

    <template #panel-icon>
      <GxIcon name="edit" :size="20" />
    </template>

    <form id="edit-post-form" class="gx-compose-form" @submit.prevent="save">
      <div class="gx-compose-field">
        <label class="gx-compose-field__label" for="edit-title">帖子标题</label>
        <Input id="edit-title" v-model="form.title" class="gx-compose-field__title" />
      </div>
      <div class="gx-compose-field gx-compose-field--grow">
        <label class="gx-compose-field__label" for="edit-content">正文内容</label>
        <Textarea id="edit-content" v-model="form.content" class="gx-compose-field__body" :rows="12" />
      </div>
    </form>

    <template #footer>
      <p class="gx-compose-footer__hint">{{ form.content.length }} 字</p>
      <button
        type="submit"
        form="edit-post-form"
        class="gx-compose-submit"
        :disabled="saving || !form.title.trim() || !form.content.trim()"
      >
        {{ saving ? '保存中…' : '保存修改' }}
      </button>
    </template>
  </GxComposeShell>
</template>
