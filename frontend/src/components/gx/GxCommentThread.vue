<script setup>
import GxAuthorChip from './GxAuthorChip.vue'
import Button from '../ui/Button.vue'
import Textarea from '../ui/Textarea.vue'

defineProps({
  comments: { type: Array, default: () => [] },
  commentCount: { type: Number, default: 0 },
  modelValue: { type: String, default: '' },
  submitting: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'submit'])
</script>

<template>
  <section class="gx-comment-thread">
    <header class="gx-comment-thread__head">
      <h2 class="gx-comment-thread__title">评论区</h2>
      <span class="gx-stat-chip">{{ commentCount }}</span>
    </header>

    <p v-if="!comments.length" class="gx-comment-thread__empty">暂无评论，来抢沙发吧</p>

    <ul v-else class="gx-comment-thread__list">
      <li v-for="(comment, index) in comments" :key="comment.id" class="gx-comment">
        <div class="gx-comment__head">
          <span class="gx-comment__floor text-caption text-gx-muted">#{{ index + 1 }}</span>
          <GxAuthorChip
            :author-id="comment.authorId"
            :author-name="comment.authorName"
            :created-at="comment.createdAt"
            size="sm"
          />
        </div>
        <p class="gx-comment__body text-meta text-gx-muted">{{ comment.content }}</p>
      </li>
    </ul>

    <div class="gx-comment-composer">
      <label class="gx-comment-composer__label" for="gx-comment-input">发表评论</label>
      <Textarea
        id="gx-comment-input"
        :model-value="modelValue"
        :rows="4"
        placeholder="遵守警校规章制度，理性文明交流"
        @update:model-value="emit('update:modelValue', $event)"
      />
      <div class="gx-comment-composer__actions">
        <Button :disabled="!modelValue.trim() || submitting" @click="emit('submit')">
          {{ submitting ? '发布中…' : '发表评论' }}
        </Button>
      </div>
    </div>
  </section>
</template>
