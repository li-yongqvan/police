<script setup>
import { computed } from 'vue'
import {
  ArrowDown,
  ArrowLeft,
  ArrowUp,
  Award,
  BarChart3,
  Bell,
  BookOpen,
  Check,
  ChevronDown,
  Circle,
  Clock,
  FileText,
  Flag,
  Flame,
  Home,
  Info,
  KeyRound,
  Lock,
  LogOut,
  MessageSquare,
  Minus,
  Pencil,
  Search,
  Settings,
  Shield,
  ShieldOff,
  Sparkles,
  Star,
  TrendingDown,
  TrendingUp,
  User,
  Users,
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
  award: Award,
  'bar-chart': BarChart3,
  settings: Settings,
  'arrow-left': ArrowLeft,
  key: KeyRound,
  lock: Lock,
  info: Info,
  edit: Pencil,
  'shield-empty': ShieldOff,
  logout: LogOut,
  star: Star,
  clock: Clock,
  file: FileText,
  users: Users,
  check: Check,
  'chevron-down': ChevronDown,
  flame: Flame,
  sparkles: Sparkles,
  'trending-up': TrendingUp,
  'trending-down': TrendingDown,
  'arrow-up': ArrowUp,
  'arrow-down': ArrowDown,
  minus: Minus,
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
