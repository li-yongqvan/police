<script setup>
import { computed } from 'vue'
import GxAuthorChip from './GxAuthorChip.vue'
import { useSessionStore } from '../../stores/session'

const props = defineProps({
  node: { type: Object, required: true },
  collapsedIds: { type: Object, required: true },
})

const emit = defineEmits(['toggle', 'reply'])
const session = useSessionStore()

const collapsed = () => props.collapsedIds.has(props.node.id)
const hasChildren = () => (props.node.children?.length ?? 0) > 0
const authorAvatar = computed(() => {
  if (props.node.authorAvatar) return props.node.authorAvatar
  if (String(props.node.authorId) === String(session.currentUser?.id)) return session.currentUser?.avatar || ''
  return ''
})
</script>

<template>
  <li class="gx-comment-node" :style="{ marginLeft: `${Math.min(node.depth || 0, 6) * 16}px` }">
    <div class="gx-comment">
      <div class="gx-comment__head">
        <GxAuthorChip
          :author-id="node.authorId"
          :author-name="node.authorName"
          :author-avatar="authorAvatar"
          :created-at="node.createdAt"
          size="sm"
        />
        <button
          v-if="hasChildren()"
          type="button"
          class="gx-comment-node__collapse text-meta text-gx-muted"
          @click="emit('toggle', node.id)"
        >
          {{ collapsed() ? `[+] ${node.children.length} 条回复` : '[-] 收起回复' }}
        </button>
      </div>
      <template v-if="!collapsed()">
        <p class="gx-comment__body text-meta text-gx-muted">{{ node.content }}</p>
        <button type="button" class="gx-comment-node__reply text-meta text-gx-accent" @click="emit('reply', node)">
          回复
        </button>
      </template>
    </div>
    <ul v-if="!collapsed() && hasChildren()" class="gx-comment-tree__children">
      <GxCommentNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :collapsed-ids="collapsedIds"
        @toggle="emit('toggle', $event)"
        @reply="emit('reply', $event)"
      />
    </ul>
  </li>
</template>
