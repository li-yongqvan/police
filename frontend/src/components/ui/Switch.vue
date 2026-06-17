<script setup>
import { cn } from '@/lib/utils'
import { SwitchRoot, SwitchThumb } from 'radix-vue'

defineProps({
  class: { type: null, default: undefined },
  modelValue: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <SwitchRoot
    :checked="modelValue"
    :disabled="disabled"
    :class="cn('peer inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gx-primary/25 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-gx-primary data-[state=unchecked]:bg-gx-border', $attrs.class)"
    @update:checked="emit('update:modelValue', $event)"
  >
    <SwitchThumb class="pointer-events-none block h-4 w-4 rounded-full bg-white shadow-lg ring-0 transition-transform data-[state=checked]:translate-x-4 data-[state=unchecked]:translate-x-0" />
  </SwitchRoot>
</template>
