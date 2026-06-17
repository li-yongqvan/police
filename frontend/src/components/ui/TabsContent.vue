<script setup>
import { cn } from '@/lib/utils'
import { TabsContent } from 'radix-vue'

defineProps({
  value: { type: String, required: true },
  class: { type: null, default: undefined },
})
</script>

<template>
  <TabsContent :value="value" :class="cn('mt-2 focus-visible:outline-none', $attrs.class)">
    <slot />
  </TabsContent>
</template>
