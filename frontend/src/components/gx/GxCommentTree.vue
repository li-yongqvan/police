<script setup>
import { computed, nextTick, ref, watch } from 'vue'
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
const composerOpen = ref(false)
const inputRef = ref(null)

const tree = computed(() => buildCommentTree(props.comments))
const showComposer = computed(() => composerOpen.value || !!props.replyToId)

function toggleCollapse(id) {
  const next = new Set(collapsed.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsed.value = next
}

async function openComposer() {
  composerOpen.value = true
  await nextTick()
  const target = inputRef.value?.$el || inputRef.value
  target?.focus?.()
}

function closeComposer() {
  composerOpen.value = false
  emit('cancel-reply')
  emit('update:modelValue', '')
}

function submitComposer() {
  if (!props.modelValue.trim() || props.submitting) return
  emit('submit')
}

watch(
  () => props.replyToId,
  (id) => {
    if (id) openComposer()
  },
)

watch(
  () => props.submitting,
  (submitting, wasSubmitting) => {
    if (wasSubmitting && !submitting && !props.modelValue.trim() && !props.replyToId) {
      composerOpen.value = false
    }
  },
)
</script>

<template>
  <section id="comments" class="gx-comment-tree">
    <header class="gx-comment-thread__head">
      <h2 class="gx-comment-thread__title">评论区</h2>
      <span class="gx-stat-chip">{{ commentCount }}</span>
    </header>

    <button
      v-if="!showComposer"
      type="button"
      class="gx-comment-prompt"
      @click="openComposer"
    >
      <span class="gx-comment-prompt__text">尊重是评论打动人心的入场券</span>
      <span class="gx-comment-prompt__action">评论</span>
    </button>

    <div v-else class="gx-comment-composer">
      <p v-if="replyToId" class="gx-comment-composer__reply-hint text-meta text-gx-muted">
        正在回复评论
        <button type="button" class="ml-2 text-gx-accent" @click="closeComposer">取消</button>
      </p>
      <label class="gx-comment-composer__label" for="gx-comment-input">发表评论</label>
      <Textarea
        id="gx-comment-input"
        ref="inputRef"
        :model-value="modelValue"
        :rows="4"
        placeholder="遵守警校规章制度，理性文明交流"
        @update:model-value="emit('update:modelValue', $event)"
      />
      <div class="gx-comment-composer__actions">
        <Button variant="secondary" :disabled="submitting" @click="closeComposer">取消</Button>
        <Button :disabled="!modelValue.trim() || submitting" @click="submitComposer">
          {{ submitting ? '发布中…' : replyToId ? '回复' : '发表评论' }}
        </Button>
      </div>
    </div>

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
  </section>
</template>
