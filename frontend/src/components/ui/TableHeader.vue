<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <thead :class="cn('[&_tr]:border-b', $attrs.class)">
    <slot />
  </thead>
</template>
