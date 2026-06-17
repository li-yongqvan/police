<script setup>
import { cn } from '@/lib/utils'
import { SelectItem, SelectItemIndicator, SelectItemText } from 'radix-vue'

defineProps({
  value: { type: String, required: true },
  class: { type: null, default: undefined },
  disabled: { type: Boolean, default: false },
})
</script>

<template>
  <SelectItem :value="value" :disabled="disabled" :class="cn('relative flex w-full cursor-default select-none items-center rounded-gx-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-gx-bg focus:text-gx-primary data-[disabled]:pointer-events-none data-[disabled]:opacity-50', $attrs.class)">
    <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
      <SelectItemIndicator>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-check"><path d="M20 6 9 17l-5-5"/></svg>
      </SelectItemIndicator>
    </span>
    <SelectItemText>
      <slot />
    </SelectItemText>
  </SelectItem>
</template>
