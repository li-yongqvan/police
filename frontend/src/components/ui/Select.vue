<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
  modelValue: { type: [String, Number], default: '' },
  disabled: { type: Boolean, default: false },
  id: { type: String, default: undefined },
})

defineEmits(['update:modelValue'])
</script>

<template>
  <select
    :id="id"
    :value="modelValue"
    :disabled="disabled"
    :class="
      cn(
        'block min-h-11 w-full cursor-pointer rounded-gx-md border border-gx-border bg-gx-surface px-3 py-2.5 text-sm leading-normal text-gx-primary shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/25 disabled:cursor-not-allowed disabled:opacity-50',
        $props.class,
      )
    "
    @change="$emit('update:modelValue', $event.target.value)"
  >
    <slot />
  </select>
</template>
