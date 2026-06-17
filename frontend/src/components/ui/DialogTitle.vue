<script setup>
import { cn } from '@/lib/utils'
import { DialogTitle } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <DialogTitle :class="cn('text-lg font-semibold text-gx-primary', $attrs.class)">
    <slot />
  </DialogTitle>
</template>
