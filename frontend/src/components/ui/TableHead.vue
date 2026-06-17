<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <th :class="cn('h-12 px-4 text-left align-middle font-medium text-gx-muted [&:has([role=checkbox])]:pr-0', $attrs.class)">
    <slot />
  </th>
</template>
