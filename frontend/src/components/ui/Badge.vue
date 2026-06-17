<script setup>
import { computed } from 'vue'
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const props = defineProps({
  variant: { type: String, default: 'default' },
  class: { type: null, default: undefined },
})

const badgeVariants = cva(
  'inline-flex items-center rounded-gx-sm border px-2 py-0.5 text-xs font-medium transition-colors',
  {
    variants: {
      variant: {
        default: 'border-transparent bg-gx-primary text-white',
        secondary: 'border-transparent bg-gx-bg text-gx-primary',
        accent: 'border-transparent bg-gx-accent/10 text-gx-accent',
        gold: 'border-transparent bg-gx-gold/20 text-gx-primary-dark',
        danger: 'border-transparent bg-gx-danger/10 text-gx-danger',
        warning: 'border-transparent bg-amber-500/15 text-amber-800',
        outline: 'border-gx-border text-gx-muted',
      },
    },
    defaultVariants: { variant: 'default' },
  },
)

const classes = computed(() => cn(badgeVariants({ variant: props.variant }), props.class))
</script>

<template>
  <span :class="classes"><slot /></span>
</template>
