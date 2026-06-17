<script setup>
import { cn } from '@/lib/utils'
import { DialogDescription } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <DialogDescription :class="cn('text-sm text-gx-muted', $attrs.class)">
    <slot />
  </DialogDescription>
</template>
