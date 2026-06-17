<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  rows: { type: [Number, String], default: 4 },
  disabled: { type: Boolean, default: false },
})

defineEmits(['update:modelValue'])
</script>

<template>
  <textarea
    :value="modelValue"
    :placeholder="placeholder"
    :rows="rows"
    :disabled="disabled"
    :class="
      cn(
        'flex min-h-[100px] w-full rounded-gx-md border border-gx-border bg-gx-surface px-3 py-2 text-sm text-[#222] shadow-sm transition-colors placeholder:text-gx-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/25 disabled:cursor-not-allowed disabled:opacity-50',
        $attrs.class,
      )
    "
    @input="$emit('update:modelValue', $event.target.value)"
  />
</template>
