<script setup>
import { cn } from '@/lib/utils'
import { SelectIcon, SelectTrigger } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <SelectTrigger :class="cn('flex h-10 w-full items-center justify-between rounded-gx-md border border-gx-border bg-gx-surface px-3 py-2 text-sm text-gx-primary shadow-sm placeholder:text-gx-muted focus:outline-none focus:ring-2 focus:ring-gx-primary/25 disabled:cursor-not-allowed disabled:opacity-50', $attrs.class)">
    <slot />
    <SelectIcon as-child>
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-down opacity-50"><path d="m6 9 6 6 6-6"/></svg>
    </SelectIcon>
  </SelectTrigger>
</template>
