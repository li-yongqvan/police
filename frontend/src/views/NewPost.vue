<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxComposeShell from '../components/gx/GxComposeShell.vue'
import GxIcon from '../components/gx/GxIcon.vue'
import Alert from '../components/ui/Alert.vue'
import Dialog from '../components/ui/Dialog.vue'
import Input from '../components/ui/Input.vue'
import Textarea from '../components/ui/Textarea.vue'
import { boardKeyFromName } from '../composables/useGxNav'
import { formatApiError } from '../api/errors'
import { forumApi } from '../api'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const boards = ref([])
const uploading = ref(false)
const sensWordDialog = ref({ open: false, words: [] })
const needsLevelForAttachment = computed(() => (session.currentUser?.level ?? 1) < 1)

const breadcrumbItems = [
  { label: '首页', to: '/community' },
  { label: '发帖' },
]

const form = ref({
  boardId: '',
  title: '',
  content: '',
  attachmentType: 'link',
  attachmentUrl: '',
  attachmentFile: null,
})

const attachTypes = [
  { id: 'link', label: '网盘链接', icon: 'info', desc: '粘贴外链地址' },
  { id: 'document', label: '文档附件', icon: 'book', desc: 'PDF / Word 等' },
  { id: 'image', label: '图片附件', icon: 'edit', desc: 'JPG / PNG 等' },
]

const composeTips = [
  '标题简明扼要，正文文明理性。',
  '涉密、违规内容请勿发布。',
  '部分帖子需人工审核后展示。',
]

const selectedBoard = computed(() => boards.value.find((board) => board.id === form.value.boardId))
const selectedAttach = computed(() => attachTypes.find((type) => type.id === form.value.attachmentType))

const canSubmit = computed(
  () => form.value.boardId && form.value.title.trim() && form.value.content.trim() && !uploading.value,
)

function boardIcon(name = '') {
  if (/校园圈|生活|日常/.test(name)) return 'flag'
  const key = boardKeyFromName(name)
  const icons = { study: 'book', training: 'shield', notice: 'bell', club: 'flag' }
  return icons[key] || 'book'
}

function applyBoardFromQuery() {
  const slug = String(route.query.board || '').trim()
  if (!slug || !boards.value.length) return
  const match = boards.value.find((board) => board.slug === slug)
  if (match) form.value.boardId = match.id
}

onMounted(async () => {
  boards.value = await forumApi.getBoards()
  form.value.boardId = boards.value[0]?.id || ''
  applyBoardFromQuery()
})

async function uploadIfNeeded() {
  if (form.value.attachmentType === 'link' && form.value.attachmentUrl) {
    return [await forumApi.uploadAttachment({ type: 'link', linkUrl: form.value.attachmentUrl })]
  }
  if (form.value.attachmentFile) {
    return [
      await forumApi.uploadAttachment({
        type: form.value.attachmentType,
        file: form.value.attachmentFile,
      }),
    ]
  }
  return []
}

async function submit() {
  try {
    const check = await forumApi.checkSensitiveWords(`${form.value.title} ${form.value.content}`)
    if (!check.clean) {
      sensWordDialog.value = { open: true, words: check.matched_words || [] }
      return
    }
  } catch {
    // Backend still validates content if the pre-check is unavailable.
  }

  uploading.value = true
  try {
    const attachmentIds = await uploadIfNeeded()
    const post = await forumApi.createPost({
      boardId: form.value.boardId,
      title: form.value.title,
      content: form.value.content,
      attachmentIds,
    })
    session.setFlash(
      post.status === 'pending_review' ? '帖子已进入审核队列。' : '帖子已成功发布。',
      post.status === 'pending_review' ? 'warning' : 'success',
    )
    if (post.status === 'published') {
      router.push(`/community/posts/${post.id}`)
    } else {
      router.push('/community')
    }
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="gx-page">
    <GxBreadcrumb :items="breadcrumbItems" />

    <GxComposeShell
      :panel-title="selectedBoard?.name || '撰写新帖'"
      :panel-subtitle="selectedBoard?.description || '填写标题与正文后发布'"
    >
      <template #rail>
        <header class="gx-compose-rail__head">
          <div class="gx-compose-rail__brand">
            <span class="gx-compose-rail__logo" aria-hidden="true">
              <GxIcon name="edit" :size="20" />
            </span>
            <div>
              <h1 class="gx-compose-rail__title">发布新帖</h1>
              <p class="gx-compose-rail__meta">选择板块 · 撰写 · 发布</p>
            </div>
          </div>
        </header>

        <nav class="gx-compose-rail__block">
          <p class="gx-compose-rail__section">发布到</p>
          <button
            v-for="board in boards"
            :key="board.id"
            type="button"
            class="gx-compose-channel"
            :class="{ 'is-active': form.boardId === board.id }"
            @click="form.boardId = board.id"
          >
            <span class="gx-compose-channel__icon" aria-hidden="true">
              <GxIcon :name="boardIcon(board.name)" :size="18" />
            </span>
            <span class="gx-compose-channel__text">
              <span class="gx-compose-channel__label">{{ board.name }}</span>
              <span v-if="board.description" class="gx-compose-channel__desc">{{ board.description }}</span>
            </span>
          </button>
        </nav>

        <nav class="gx-compose-rail__block">
          <p class="gx-compose-rail__section">附件（可选）</p>
          <button
            v-for="type in attachTypes"
            :key="type.id"
            type="button"
            class="gx-compose-channel"
            :class="{ 'is-active': form.attachmentType === type.id, 'is-disabled': needsLevelForAttachment }"
            :disabled="needsLevelForAttachment"
            @click="form.attachmentType = type.id"
          >
            <span class="gx-compose-channel__icon" aria-hidden="true">
              <GxIcon :name="type.icon" :size="18" />
            </span>
            <span class="gx-compose-channel__text">
              <span class="gx-compose-channel__label">{{ type.label }}</span>
              <span class="gx-compose-channel__desc">{{ type.desc }}</span>
            </span>
          </button>
        </nav>

        <div class="gx-compose-rail__tips">
          <p class="gx-compose-rail__section">发帖须知</p>
          <ul class="gx-compose-tips">
            <li v-for="tip in composeTips" :key="tip">{{ tip }}</li>
          </ul>
        </div>
      </template>

      <template #panel-icon>
        <GxIcon :name="boardIcon(selectedBoard?.name)" :size="20" />
      </template>

      <form id="new-post-form" class="gx-compose-form" @submit.prevent="submit">
        <div class="gx-compose-field">
          <label class="gx-compose-field__label" for="title">帖子标题</label>
          <Input
            id="title"
            v-model="form.title"
            class="gx-compose-field__title"
            placeholder="输入清晰、具体的标题"
          />
        </div>

        <div class="gx-compose-field gx-compose-field--grow">
          <label class="gx-compose-field__label" for="content">正文内容</label>
          <Textarea
            id="content"
            v-model="form.content"
            class="gx-compose-field__body"
            :rows="12"
            placeholder="遵守警校规章制度，理性文明交流…"
          />
        </div>

        <div v-if="!needsLevelForAttachment && selectedAttach" class="gx-compose-attach">
          <p class="gx-compose-field__label">{{ selectedAttach.label }}</p>
          <Input
            v-if="form.attachmentType === 'link'"
            id="attach-url"
            v-model="form.attachmentUrl"
            type="url"
            placeholder="https://"
          />
          <label v-else class="gx-compose-file">
            <input
              id="attach-file"
              type="file"
              class="gx-compose-file__input"
              @change="(event) => (form.attachmentFile = event.target.files?.[0])"
            />
            <span class="gx-compose-file__btn">选择文件</span>
            <span class="gx-compose-file__name">{{ form.attachmentFile?.name || '未选择文件' }}</span>
          </label>
        </div>

        <Alert v-if="needsLevelForAttachment" class="gx-compose-alert">
          附件功能需等级 ≥ 2（当前 Lv.{{ session.currentUser?.level ?? 1 }}）
        </Alert>
      </form>

      <template #footer>
        <p class="gx-compose-footer__hint">
          {{ form.content.length }} 字 · 发布至「{{ selectedBoard?.name || '—' }}」
        </p>
        <button type="submit" form="new-post-form" class="gx-compose-submit" :disabled="!canSubmit">
          {{ uploading ? '提交中…' : '发布帖子' }}
        </button>
      </template>
    </GxComposeShell>

    <Dialog
      :open="sensWordDialog.open"
      title="内容违规提示"
      @update:open="(value) => { if (!value) sensWordDialog.open = false }"
    >
      <p class="text-body text-gx-muted">该内容含有违规词，请修改后重新发送。</p>
      <template #footer>
        <button type="button" class="gx-btn gx-btn--primary" @click="sensWordDialog.open = false">我知道了</button>
      </template>
    </Dialog>
  </div>
</template>
