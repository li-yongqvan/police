<script setup>
import { computed, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import GxAvatarInitial from './GxAvatarInitial.vue'
import GxVoteRail from './GxVoteRail.vue'
import GxIcon from './GxIcon.vue'
import { boardKeyFromName } from '../../composables/useGxNav'
import { formatAuthorLabel, formatRelativeTime } from '../../utils/displayName'
import { formatApiError } from '../../api/errors'
import { forumApi } from '../../api'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  post: { type: Object, required: true },
  pinned: { type: Boolean, default: false },
  announce: { type: Boolean, default: false },
})

const session = useSessionStore()
const liked = ref(!!props.post.liked)
const likeCount = ref(props.post.likeCount ?? 0)
const voting = ref(false)

watch(
  () => [props.post.liked, props.post.likeCount],
  ([nextLiked, nextCount]) => {
    liked.value = !!nextLiked
    likeCount.value = nextCount ?? 0
  },
)

const boardPath = computed(() => {
  const slug = props.post.boardSlug || boardKeyFromName(props.post.boardName)
  return `/community/boards/${slug || 'study'}`
})

const excerpt = computed(() => {
  const text = String(props.post.content || '').trim()
  if (!text) return ''
  return text.length > 200 ? `${text.slice(0, 200)}…` : text
})

const authorLevel = computed(() => {
  const id = Number(props.post.authorId) || 1
  return Math.min(5, (id % 4) + 1)
})

const authorLabel = computed(() => formatAuthorLabel(null, props.post))
const authorAvatar = computed(() => {
  if (props.post.authorAvatar) return props.post.authorAvatar
  if (String(props.post.authorId) === String(session.currentUser?.id)) return session.currentUser?.avatar || ''
  return ''
})

function numericOrNull(value) {
  if (value === null || value === undefined || value === '') return null
  const n = Number(value)
  return Number.isFinite(n) ? n : null
}

const trendBadge = computed(() => {
  const explicitDirection = String(props.post.trendDirection || '').toLowerCase()
  const delta = numericOrNull(props.post.trendDelta)
  const heatScore =
    numericOrNull(props.post.hotScore) ??
    (Number(likeCount.value || 0) * 2 + Number(props.post.commentCount || 0) * 3)

  if (explicitDirection.includes('up') || explicitDirection.includes('rise')) {
    return { icon: 'trending-up', tone: 'rise', label: '这帖子的排名上升啦' }
  }
  if (explicitDirection.includes('down') || explicitDirection.includes('fall')) {
    return { icon: 'trending-down', tone: 'fall', label: '这帖子的讨论热度降了一些' }
  }
  if (delta !== null) {
    if (delta >= 12) return { icon: 'flame-up', tone: 'hot', label: '这个帖子最近很火爆' }
    if (delta >= 4) return { icon: 'trending-up', tone: 'rise', label: '这帖子的排名上升啦' }
    if (delta > 0) return { icon: 'arrow-up', tone: 'warm', label: '这帖子的关注度在升温' }
    if (delta <= -12) return { icon: 'trending-down', tone: 'fall-strong', label: '这帖子最近没那么热了' }
    if (delta <= -4) return { icon: 'trending-down', tone: 'fall', label: '这帖子的讨论热度降了一些' }
    if (delta < 0) return { icon: 'arrow-down', tone: 'cool', label: '这帖子的热度稍有回落' }
  }
  if (heatScore >= 30) return { icon: 'flame-up', tone: 'hot', label: '这个帖子最近很火爆' }
  if (heatScore >= 12) return { icon: 'trending-up', tone: 'rise', label: '这帖子的关注度在升温' }
  if (heatScore <= 0) return { icon: 'sparkles', tone: 'new', label: '新帖子，还在观察热度' }
  return { icon: 'minus', tone: 'stable', label: '这帖子的热度比较稳定' }
})

async function onVote() {
  if (!session.currentUser) {
    session.setFlash('请先登录后再点赞', 'info')
    return
  }
  voting.value = true
  try {
    const resp = await forumApi.likePost(props.post.id)
    liked.value = resp.liked
    likeCount.value = resp.likeCount
  } catch (error) {
    session.setFlash(formatApiError(error), 'info')
  } finally {
    voting.value = false
  }
}
</script>

<template>
  <article
    class="gx-feed-card gx-feed-card--rich"
    :class="{
      'gx-feed-card--announce': announce,
      'gx-feed-card--pinned': pinned && !announce,
    }"
  >
    <GxVoteRail
      :score="likeCount"
      :liked="liked"
      :loading="voting"
      :trend-icon="trendBadge.icon"
      :trend-tone="trendBadge.tone"
      :trend-label="trendBadge.label"
      compact
      trend-only
      @vote="onVote"
    />
    <div class="gx-feed-card__body">
      <div class="gx-feed-card__meta">
        <span v-if="pinned" class="gx-feed-card__pin-tag">[置顶]</span>
        <RouterLink v-else :to="boardPath" class="gx-feed-card__board" @click.stop>
          {{ post.boardName || '讨论' }}
        </RouterLink>
      </div>
      <RouterLink :to="`/community/posts/${post.id}`" class="gx-feed-card__title">
        {{ post.title }}
      </RouterLink>
      <p v-if="excerpt" class="gx-feed-card__excerpt">{{ excerpt }}</p>
      <footer class="gx-feed-card__foot">
        <RouterLink :to="`/community/users/${post.authorId}`" class="gx-feed-card__author" @click.stop>
          <GxAvatarInitial :name="authorLabel" :src="authorAvatar" :size="28" />
          <span class="gx-feed-card__author-name">{{ authorLabel }}</span>
          <span class="gx-feed-card__level">Lv.{{ authorLevel }}</span>
        </RouterLink>
        <span class="gx-feed-card__foot-meta">
          <time :datetime="post.createdAtIso">{{ formatRelativeTime(post.createdAtIso) }}</time>
          <RouterLink :to="`/community/posts/${post.id}#comments`" class="gx-feed-card__foot-stat">
            <GxIcon name="message" :size="14" />
            {{ post.commentCount }}
          </RouterLink>
          <button type="button" class="gx-feed-card__foot-stat" @click.prevent="onVote">
            <GxIcon name="star" :size="14" />
            {{ likeCount }}
          </button>
        </span>
      </footer>
    </div>
  </article>
</template>
