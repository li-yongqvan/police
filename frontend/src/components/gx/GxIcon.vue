<script setup>
import { computed } from 'vue'
import {
  Bell,
  BookOpen,
  Circle,
  Clock,
  Flag,
  Home,
  Info,
  LogOut,
  MessageSquare,
  Pencil,
  Search,
  Shield,
  ShieldOff,
  Star,
  User,
} from 'lucide-vue-next'

const props = defineProps({
  name: { type: String, required: true },
  size: { type: Number, default: 24 },
  class: { type: null, default: undefined },
})

const ICON_MAP = {
  book: BookOpen,
  shield: Shield,
  bell: Bell,
  flag: Flag,
  search: Search,
  message: MessageSquare,
  user: User,
  home: Home,
  info: Info,
  edit: Pencil,
  'shield-empty': ShieldOff,
  logout: LogOut,
  star: Star,
  clock: Clock,
}

const IconComponent = computed(() => ICON_MAP[props.name] || Circle)
</script>

<template>
  <component
    :is="IconComponent"
    class="gx-icon shrink-0"
    :size="size"
    :stroke-width="2"
    aria-hidden="true"
  />
</template>