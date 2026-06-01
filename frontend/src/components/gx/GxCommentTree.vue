<script setup>
import { computed, ref } from 'vue'
import GxCommentNode from './GxCommentNode.vue'
import Button from '../ui/Button.vue'
import Textarea from '../ui/Textarea.vue'
import { buildCommentTree } from '../../composables/buildCommentTree'

const props = defineProps({
  comments: { type: Array, default: () => [] },
  commentCount: { type: Number, default: 0 },
  modelValue: { type: String, default: '' },
  submitting: { type: Boolean, default: false },
  replyToId: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'submit', 'reply', 'cancel-reply'])

const collapsed = ref(new Set())

const tree = computed(() => buildCommentTree(props.comments))

function toggleCollapse(id) {
  const next = new Set(collapsed.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsed.value = next
}
</script>

<template>
  <section id="comments" class="gx-comment-tree">
    <header class="gx-comment-thread__head">
      <h2 class="gx-comment-thread__title">评论区</h2>
      <span class="gx-stat-chip">{{ commentCount }}</span>
    </header>

    <p v-if="!comments.length" class="gx-comment-thread__empty">暂无评论，来抢沙发吧</p>

    <ul v-else class="gx-comment-tree__list">
      <GxCommentNode
        v-for="node in tree"
        :key="node.id"
        :node="node"
        :collapsed-ids="collapsed"
        @toggle="toggleCollapse"
        @reply="emit('reply', $event)"
      />
    </ul>

    <div class="gx-comment-composer">
      <p v-if="replyToId" class="gx-comment-composer__reply-hint text-meta text-gx-muted">
        正在回复评论
        <button type="button" class="ml-2 text-gx-accent" @click="emit('cancel-reply')">取消</button>
      </p>
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
          {{ submitting ? '发布中…' : replyToId ? '回复' : '发表评论' }}
        </Button>
      </div>
    </div>
  </section>
</template>
