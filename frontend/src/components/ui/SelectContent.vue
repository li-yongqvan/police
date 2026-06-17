<script setup>
import { cn } from '@/lib/utils'
import { SelectContent, SelectPortal, SelectScrollDownButton, SelectScrollUpButton, SelectViewport } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <SelectPortal>
    <SelectContent :class="cn('relative z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-gx-md border border-gx-border bg-gx-surface text-gx-primary shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95', $attrs.class)" :side-offset="4">
      <SelectScrollUpButton class="flex h-6 items-center justify-center bg-gx-surface">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-up"><path d="m18 15-6-6-6 6"/></svg>
      </SelectScrollUpButton>
      <SelectViewport class="p-1">
        <slot />
      </SelectViewport>
      <SelectScrollDownButton class="flex h-6 items-center justify-center bg-gx-surface">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-down"><path d="m6 9 6 6 6-6"/></svg>
      </SelectScrollDownButton>
    </SelectContent>
  </SelectPortal>
</template>
