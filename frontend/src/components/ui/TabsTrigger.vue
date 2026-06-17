<script setup>
import { cn } from '@/lib/utils'
import { TabsTrigger } from 'radix-vue'

defineProps({
  value: { type: String, required: true },
  class: { type: null, default: undefined },
  disabled: { type: Boolean, default: false },
})
</script>

<template>
  <TabsTrigger :value="value" :disabled="disabled" :class="cn('inline-flex items-center justify-center whitespace-nowrap rounded-gx-sm px-3 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/25 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:bg-gx-surface data-[state=active]:text-gx-primary data-[state=active]:shadow-sm', $attrs.class)">
    <slot />
  </TabsTrigger>
</template>
