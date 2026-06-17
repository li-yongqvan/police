<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <div :class="cn('shrink-0 bg-gx-border h-[1px] w-full', $attrs.class)" role="separator" />
</template>
