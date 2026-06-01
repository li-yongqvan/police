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
const collected = ref(false)
const dialog = ref({ open: false, mode: 'report', reason: '' })

const attachments = computed(() => payload.value.post?.attachments ?? [])
const isAuthor = computed(() => payload.value.post?.authorId === session.currentUser?.id)
const post = computed(() => payload.value.post)

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
    collected.value = !!p.collected
    recordBrowseHistory(p)
  }
}

async function likePost() {
  actionLoading.value = true
  try {
    const resp = await forumApi.likePost(route.params.id)
    liked.value = resp.liked
    if (payload.value.post) payload.value.post.likeCount = resp.likeCount
    session.setFlash(`点赞已记录（${resp.likeCount} 赞）`, 'success')
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

async function submitComment() {
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
            :loading="actionLoading"
            @vote="likePost"
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
                :created-at="post.createdAt"
              />
              <div class="gx-post-hero__stats">
                <span class="gx-stat-chip">{{ post.likeCount }} 赞</span>
                <span class="gx-stat-chip">{{ post.commentCount }} 评论</span>
              </div>
            </div>
          </header>

          <div class="gx-post-body">{{ post.content }}</div>

          <div v-if="attachments.length" class="gx-attachments mt-4 flex flex-wrap gap-2">
            <a
              v-for="item in attachments"
              :key="item.id"
              :href="item.url"
              target="_blank"
              rel="noopener"
              class="rounded-gx-sm border border-gx-border px-3 py-1 text-body text-gx-primary hover:bg-gx-bg"
            >
              {{ item.title }}
            </a>
          </div>

          <GxActionToolbar
            class="gx-action-row--in-article"
            :liked="liked"
            :collected="collected"
            :loading="actionLoading"
            :is-author="isAuthor"
            :post-id="post.id"
            :like-count="post.likeCount"
            @like="likePost"
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
          />
        </GxReadingColumn>
      </Card>
    </div>

    <aside class="gx-post-detail__aside">
      <Card class="gx-post-detail__aside-card p-5">
        <h3 class="gx-post-detail__aside-title">关于本帖</h3>
        <dl class="gx-meta-list gx-meta-list--compact">
          <div class="gx-meta-list__row">
            <dt>板块</dt>
            <dd>{{ post.boardName }}</dd>
          </div>
          <div class="gx-meta-list__row">
            <dt>作者</dt>
            <dd>
              <RouterLink class="text-gx-accent hover:underline" :to="`/community/users/${post.authorId}`">
                {{ post.authorName }}
              </RouterLink>
            </dd>
          </div>
          <div class="gx-meta-list__row">
            <dt>发布</dt>
            <dd>
              <time :datetime="post.createdAtIso">{{ post.createdAt }}</time>
            </dd>
          </div>
        </dl>

        <RouterLink :to="`/community/users/${post.authorId}`" class="gx-post-detail__author-link">
          查看作者主页
        </RouterLink>

        <GxActionToolbar
          layout="vertical"
          :liked="liked"
          :collected="collected"
          :loading="actionLoading"
          :is-author="isAuthor"
          :post-id="post.id"
          :like-count="post.likeCount"
          @like="likePost"
          @collect="collectPost"
          @report="openReport"
          @delete="openDelete"
        />
      </Card>
    </aside>

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
  </div>
</template>
