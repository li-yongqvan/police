<script setup>
import { cn } from '@/lib/utils'
import { CheckboxIndicator, CheckboxRoot } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
  modelValue: { type: [Boolean, String], default: false },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <CheckboxRoot
    :checked="modelValue"
    :disabled="disabled"
    :class="cn('peer h-4 w-4 shrink-0 rounded-gx-sm border border-gx-border shadow focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/25 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-gx-primary data-[state=checked]:text-white', $attrs.class)"
    @update:checked="emit('update:modelValue', $event)"
  >
    <CheckboxIndicator class="flex items-center justify-center text-current">
      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
    </CheckboxIndicator>
  </CheckboxRoot>
</template>
