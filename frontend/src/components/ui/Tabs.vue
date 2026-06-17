<script setup>
import { cn } from '@/lib/utils'
import { TabsList, TabsRoot, TabsTrigger } from 'radix-vue'

defineProps({
  defaultValue: { type: String, default: undefined },
  class: { type: null, default: undefined },
})
</script>

<template>
  <TabsRoot :default-value="defaultValue" :class="cn('', $attrs.class)">
    <slot />
  </TabsRoot>
</template>
