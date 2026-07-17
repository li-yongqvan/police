<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { formatApiError } from '../api/errors'
import { forumApi } from '../api'
import GxActionToolbar from '../components/gx/GxActionToolbar.vue'
import GxAuthorChip from '../components/gx/GxAuthorChip.vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxCommentTree from '../components/gx/GxCommentTree.vue'
import GxVoteRail from '../components/gx/GxVoteRail.vue'
import GxReadingColumn from '../components/gx/GxReadingColumn.vue'
import Badge from '../components/ui/Badge.vue'
import Button from '../components/ui/Button.vue'
import Card from '../components/ui/Card.vue'
import Dialog from '../components/ui/Dialog.vue'
import Input from '../components/ui/Input.vue'
import Label from '../components/ui/Label.vue'
import { boardKeyFromName, boardTagClass } from '../composables/useGxNav'
import { recordBrowseHistory } from '../composables/useBrowseHistory'
import { useSessionStore } from '../stores/session'

const route = useRoute()
const router = useRouter()
const session = useSessionStore()
const payload = ref({ post: null, comments: [] })
const commentText = ref('')
const replyToId = ref('')
const commentSubmitting = ref(false)
const actionLoading = ref(false)
const liked = ref(false)
const disliked = ref(false)
const collected = ref(false)
const dialog = ref({ open: false, mode: 'report', reason: '' })
const sensWordDialog = ref({ open: false, words: [] })

const attachments = computed(() => payload.value.post?.attachments ?? [])
const imageAttachments = computed(() => attachments.value.filter((item) => item.type === 'image'))
const fileAttachments = computed(() => attachments.value.filter((item) => item.type !== 'image'))
const isAuthor = computed(() => payload.value.post?.authorId === session.currentUser?.id)
const post = computed(() => payload.value.post)
const postAuthorAvatar = computed(() => {
  if (post.value?.authorAvatar) return post.value.authorAvatar
  if (String(post.value?.authorId) === String(session.currentUser?.id)) return session.currentUser?.avatar || ''
  return ''
})

const boardKey = computed(() => boardKeyFromName(post.value?.boardName || ''))

const breadcrumbItems = computed(() => [
  { label: '首页', to: '/community' },
  { label: post.value?.boardName || '板块', to: `/community/boards/${boardKey.value}` },
  { label: '帖子' },
])

function badgeVariant(boardName = '') {
  const cls = boardTagClass(boardName)
  if (cls.includes('club')) return 'gold'
  if (cls.includes('notice')) return 'accent'
  return 'secondary'
}

async function loadPost() {
  payload.value = await forumApi.getPost(route.params.id)
  const p = payload.value.post
  if (p) {
    liked.value = !!p.liked
    disliked.value = !!p.disliked
    collected.value = !!p.collected
    recordBrowseHistory(p)
  }
}

async function likePost() {
  actionLoading.value = true
  try {
    const resp = await forumApi.likePost(route.params.id)
    liked.value = resp.liked
    disliked.value = resp.disliked
    if (payload.value.post) {
      payload.value.post.likeCount = resp.likeCount
      payload.value.post.dislikeCount = resp.dislikeCount
    }
    session.setFlash(`点赞已记录（${resp.likeCount} 赞）`, 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    actionLoading.value = false
  }
}


async function dislikePost() {
  actionLoading.value = true
  try {
    const resp = await forumApi.dislikePost(route.params.id)
    liked.value = resp.liked
    disliked.value = resp.disliked
    if (payload.value.post) {
      payload.value.post.likeCount = resp.likeCount
      payload.value.post.dislikeCount = resp.dislikeCount
    }
    session.setFlash(`点踩已记录（${resp.dislikeCount} 踩）`, 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    actionLoading.value = false
  }
}

async function collectPost() {
  actionLoading.value = true
  try {
    const resp = await forumApi.collectPost(route.params.id)
    collected.value = resp.collected ?? !collected.value
    session.setFlash(collected.value ? '已收藏' : '已取消收藏', 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    actionLoading.value = false
  }
}

function openReport() {
  dialog.value = { open: true, mode: 'report', reason: '' }
}

function openDelete() {
  dialog.value = { open: true, mode: 'delete', reason: '' }
}

async function confirmDialog() {
  if (dialog.value.mode === 'report') {
    if (!dialog.value.reason.trim()) return
    try {
      await forumApi.reportPost(route.params.id, dialog.value.reason.trim())
      session.setFlash('举报已提交', 'success')
    } catch (error) {
      session.setFlash(formatApiError(error), 'info')
    }
  } else if (dialog.value.mode === 'delete') {
    await forumApi.deletePost(route.params.id)
    session.setFlash('帖子已删除', 'success')
    router.push('/community')
  }
  dialog.value.open = false
}

function onReply(comment) {
  replyToId.value = comment.id
}

function cancelReply() {
  replyToId.value = ''
}

function updateCommentInTree(commentId, updater) {
  const walk = (items) => {
    for (const item of items) {
      if (item.id === commentId) {
        updater(item)
        return true
      }
      if (item.children?.length && walk(item.children)) return true
    }
    return false
  }
  walk(payload.value.comments || [])
}

async function likeComment(comment) {
  try {
    const resp = await forumApi.likeComment(comment.id)
    updateCommentInTree(comment.id, (item) => {
      item.likeCount = resp.likeCount
      item.liked = resp.liked
      if (resp.liked) {
        item.disliked = false
      }
    })
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  }
}

async function dislikeComment(comment) {
  try {
    const resp = await forumApi.dislikeComment(comment.id)
    updateCommentInTree(comment.id, (item) => {
      item.dislikeCount = resp.dislikeCount
      item.disliked = resp.disliked
      if (resp.disliked) {
        item.liked = false
      }
    })
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  }
}

async function submitComment() {
  // Pre-check sensitive words
  try {
    const check = await forumApi.checkSensitiveWords(commentText.value)
    if (!check.clean) {
      sensWordDialog.value = { open: true, words: check.matched_words || [] }
      return
    }
  } catch (e) {
    // If check fails, proceed to submit (backend will handle)
  }

  commentSubmitting.value = true
  try {
    await forumApi.createComment(route.params.id, {
      content: commentText.value,
      parentId: replyToId.value || undefined,
    })
    commentText.value = ''
    replyToId.value = ''
    await loadPost()
    session.setFlash('评论已发布', 'success')
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    commentSubmitting.value = false
  }
}

onMounted(loadPost)
</script>

<template>
  <div v-if="post" class="gx-page gx-post-detail-layout">
    <div class="gx-post-detail__main">
      <Card class="gx-post-detail__article gx-post-detail__article--feed p-5 md:p-6">
        <div class="gx-post-detail__with-vote">
          <GxVoteRail
            :score="post.likeCount"
            :liked="liked"
            :disliked="disliked"
            :loading="actionLoading"
            @vote="likePost"
            @dislike="dislikePost"
          />
          <GxReadingColumn class="gx-post-detail__vote-content">
          <GxBreadcrumb :items="breadcrumbItems" />

          <header class="gx-post-hero">
            <div class="gx-post-hero__tags">
              <Badge :variant="badgeVariant(post.boardName)">{{ post.boardName }}</Badge>
            </div>
            <h1 class="gx-post-hero__title text-display text-gx-primary">{{ post.title }}</h1>
            <div class="gx-post-hero__meta">
              <GxAuthorChip
                :author-id="post.authorId"
                :author-name="post.authorName"
                :author-avatar="postAuthorAvatar"
                :created-at="post.createdAt"
              />
              <div class="gx-post-hero__stats">
                <span class="gx-stat-chip">{{ post.likeCount }} 赞</span>
                <span class="gx-stat-chip">{{ post.commentCount }} 评论</span>
              </div>
            </div>
          </header>

          <div class="gx-post-body">{{ post.content }}</div>

          <div v-if="attachments.length" class="gx-attachments">
            <div v-if="imageAttachments.length" class="gx-attachments__images">
              <a
                v-for="item in imageAttachments"
                :key="item.id"
                :href="item.url"
                target="_blank"
                rel="noopener"
                class="gx-attachment-image"
              >
                <img :src="item.url" :alt="item.title" loading="lazy" />
                <span>{{ item.title }}</span>
              </a>
            </div>

            <div v-if="fileAttachments.length" class="gx-attachments__files">
              <a
                v-for="item in fileAttachments"
                :key="item.id"
                :href="item.url"
                target="_blank"
                rel="noopener"
                class="gx-attachment-file"
              >
                {{ item.title }}
              </a>
            </div>
          </div>

          <GxActionToolbar
            class="gx-action-row--in-article"
            :liked="liked"
            :disliked="disliked"
            :collected="collected"
            :loading="actionLoading"
            :is-author="isAuthor"
            :post-id="post.id"
            :like-count="post.likeCount"
            :dislike-count="post.dislikeCount"
            @like="likePost"
            @dislike="dislikePost"
            @collect="collectPost"
            @report="openReport"
            @delete="openDelete"
          />
          </GxReadingColumn>
        </div>

        <div class="gx-post-detail__divider" />

        <GxReadingColumn>
          <GxCommentTree
            v-model="commentText"
            :comments="payload.comments"
            :comment-count="post.commentCount"
            :submitting="commentSubmitting"
            :reply-to-id="replyToId"
            @submit="submitComment"
            @reply="onReply"
            @cancel-reply="cancelReply"
            @like="likeComment"
            @dislike="dislikeComment"
          />
        </GxReadingColumn>
      </Card>
    </div>


    <Dialog
      :open="dialog.open"
      :title="dialog.mode === 'report' ? '举报帖子' : '删除帖子'"
      @update:open="(v) => { if (!v) dialog.open = false }"
    >
      <p v-if="dialog.mode === 'delete'" class="text-body text-gx-muted">确定删除这篇帖子？此操作不可撤销。</p>
      <div v-else class="space-y-2">
        <Label for="report-reason">举报理由</Label>
        <Input id="report-reason" v-model="dialog.reason" placeholder="涉嫌违规或不当内容" />
      </div>
      <template #footer>
        <Button variant="secondary" @click="dialog.open = false">取消</Button>
        <Button
          :variant="dialog.mode === 'delete' ? 'destructive' : 'default'"
          @click="confirmDialog"
        >
          {{ dialog.mode === 'delete' ? '确认删除' : '提交举报' }}
        </Button>
      </template>
    </Dialog>
    <Dialog
      :open="sensWordDialog.open"
      title="内容违规提示"
      @update:open="(v) => { if (!v) sensWordDialog.open = false }"
    >
      <p class="text-body text-gx-muted">该评论含有违规字体，请修改后重新发送。</p>
      <template #footer>
        <button type="button" class="gx-btn gx-btn--primary" @click="sensWordDialog.open = false">我知道了</button>
      </template>
    </Dialog>
  </div>
</template>
