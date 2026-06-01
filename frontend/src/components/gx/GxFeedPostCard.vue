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
      compact
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
          <GxAvatarInitial :name="authorLabel" :size="28" />
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
