<script setup>
import { Bookmark, Flag, Heart, Pencil, Trash2 } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'

defineProps({
  liked: { type: Boolean, default: false },
  collected: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  isAuthor: { type: Boolean, default: false },
  postId: { type: [String, Number], default: '' },
  layout: { type: String, default: 'horizontal' },
  likeCount: { type: Number, default: 0 },
})

const emit = defineEmits(['like', 'collect', 'report', 'delete'])
</script>

<template>
  <div class="gx-action-toolbar" :class="`gx-action-toolbar--${layout}`" role="toolbar" aria-label="帖子操作">
    <button
      type="button"
      class="gx-action-toolbar__btn"
      :class="{ 'is-active': liked }"
      :disabled="loading"
      aria-label="点赞"
      @click="emit('like')"
    >
      <Heart :size="18" :fill="liked ? 'currentColor' : 'none'" />
      <span>点赞<span v-if="likeCount > 0" class="gx-action-toolbar__count">{{ likeCount }}</span></span>
    </button>
    <button
      type="button"
      class="gx-action-toolbar__btn"
      :class="{ 'is-active': collected }"
      :disabled="loading"
      aria-label="收藏"
      @click="emit('collect')"
    >
      <Bookmark :size="18" :fill="collected ? 'currentColor' : 'none'" />
      <span>收藏</span>
    </button>
    <button
      type="button"
      class="gx-action-toolbar__btn gx-action-toolbar__btn--muted"
      :disabled="loading"
      aria-label="举报"
      @click="emit('report')"
    >
      <Flag :size="18" />
      <span>举报</span>
    </button>
    <RouterLink
      v-if="isAuthor && postId"
      :to="`/community/posts/${postId}/edit`"
      class="gx-action-toolbar__btn gx-action-toolbar__btn--ghost"
    >
      <Pencil :size="18" />
      <span>编辑</span>
    </RouterLink>
    <button
      v-if="isAuthor"
      type="button"
      class="gx-action-toolbar__btn gx-action-toolbar__btn--danger"
      :disabled="loading"
      aria-label="删除"
      @click="emit('delete')"
    >
      <Trash2 :size="18" />
      <span>删除</span>
    </button>
  </div>
</template>
