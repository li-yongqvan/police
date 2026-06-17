<script setup>
import { cn } from '@/lib/utils'

defineProps({
  class: { type: null, default: undefined },
  modelValue: { type: [String, Number], default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  autocomplete: { type: String, default: undefined },
  disabled: { type: Boolean, default: false },
})

defineEmits(['update:modelValue'])
</script>

<template>
  <input
    :type="type"
    :value="modelValue"
    :placeholder="placeholder"
    :autocomplete="autocomplete"
    :disabled="disabled"
    :class="
      cn(
        'block min-h-11 w-full rounded-gx-md border border-gx-border bg-gx-surface px-3 py-2.5 text-sm leading-normal text-[#222] shadow-sm transition-colors placeholder:text-gx-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/25 disabled:cursor-not-allowed disabled:opacity-50',
        $attrs.class,
      )
    "
    @input="$emit('update:modelValue', $event.target.value)"
  />
</template>
