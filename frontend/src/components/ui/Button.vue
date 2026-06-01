<script setup>
import { computed } from 'vue'
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const props = defineProps({
  variant: { type: String, default: 'default' },
  size: { type: String, default: 'default' },
  class: { type: null, default: undefined },
  type: { type: String, default: 'button' },
  disabled: { type: Boolean, default: false },
})

const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-gx-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/30 disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'bg-gx-primary text-white hover:bg-gx-primary-dark',
        secondary: 'border border-gx-border bg-gx-surface text-gx-primary hover:bg-gx-bg',
        ghost: 'text-gx-primary hover:bg-gx-bg',
        destructive: 'bg-gx-accent text-white hover:bg-gx-accent/90',
        outline: 'border border-gx-border bg-transparent hover:bg-gx-bg',
      },
      size: {
        default: 'h-10 px-4 py-2',
        sm: 'h-8 rounded-gx-sm px-3 text-xs',
        lg: 'h-11 rounded-gx-md px-8',
        icon: 'h-10 w-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

const classes = computed(() => cn(buttonVariants({ variant: props.variant, size: props.size }), props.class))
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled">
    <slot />
  </button>
</template>
