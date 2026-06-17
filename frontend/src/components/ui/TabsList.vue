<script setup>
import { cn } from '@/lib/utils'
import { TabsList } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <TabsList :class="cn('inline-flex h-10 items-center justify-center rounded-gx-md bg-gx-bg p-1 text-gx-muted', $attrs.class)">
    <slot />
  </TabsList>
</template>
