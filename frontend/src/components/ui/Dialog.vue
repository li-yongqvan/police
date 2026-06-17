<script setup>
import { cn } from '@/lib/utils'

defineProps({
  open: { type: Boolean, default: false },
  title: { type: String, default: '' },
  class: { type: null, default: undefined },
})

defineEmits(['update:open', 'confirm', 'cancel'])
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
      :aria-label="title"
    >
      <div class="absolute inset-0 bg-black/40" @click="$emit('update:open', false); $emit('cancel')" />
      <div
        :class="
          cn(
            'relative z-10 w-full max-w-md rounded-gx-lg border border-gx-border bg-gx-surface p-6 shadow-lg',
            $props.class,
          )
        "
      >
        <h2 v-if="title" class="text-title text-gx-primary">{{ title }}</h2>
        <div class="mt-4">
          <slot />
        </div>
        <div v-if="$slots.footer" class="mt-6 flex justify-end gap-2">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </Teleport>
</template>
