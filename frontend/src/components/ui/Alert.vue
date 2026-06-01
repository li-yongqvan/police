<script setup>
import { cn } from '@/lib/utils'

defineProps({
  variant: { type: String, default: 'default' },
  class: { type: null, default: undefined },
})
</script>

<template>
  <div
    role="alert"
    :class="
      cn(
        'relative w-full rounded-gx-md border px-4 py-3 text-sm',
        variant === 'destructive'
          ? 'border-gx-accent/30 bg-gx-accent/5 text-gx-accent'
          : 'border-gx-border bg-gx-bg text-gx-muted',
        $attrs.class,
      )
    "
  >
    <slot />
  </div>
</template>
