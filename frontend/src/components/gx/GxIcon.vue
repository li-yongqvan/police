<script setup>
import { computed, h } from 'vue'
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

const FlameUpIcon = (props, { attrs }) =>
  h(
    'svg',
    {
      ...attrs,
      width: props.size,
      height: props.size,
      viewBox: '0 0 24 24',
      fill: 'none',
      xmlns: 'http://www.w3.org/2000/svg',
      'aria-hidden': 'true',
    },
    [
      h('path', {
        d: 'M12 22c-4.7 0-8.2-3.1-8.2-7.4 0-2.9 1.6-5.3 3.3-7.2 1.7-1.8 3.1-3.9 3-5.7 0-.5.5-.8.9-.5 3.9 2.1 6.4 5.5 6.4 9 0 .5 0 .9-.1 1.3.8-.7 1.3-1.7 1.4-2.7 0-.5.6-.7 1-.3 1.5 1.6 2.4 3.7 2.4 6.1 0 4.3-3.5 7.4-8.1 7.4Z',
        fill: 'currentColor',
      }),
      h('path', {
        d: 'M8.2 15.1c.9 2.2 2.6 3.4 4.6 3.4s3.8-1.2 4.7-3.4',
        stroke: '#fff',
        'stroke-width': '1.8',
        'stroke-linecap': 'round',
      }),
      h('path', {
        d: 'M12.8 17.7V9.5',
        stroke: '#fff',
        'stroke-width': '2',
        'stroke-linecap': 'round',
      }),
      h('path', {
        d: 'M9.7 12.3 12.8 9l3.1 3.3',
        stroke: '#fff',
        'stroke-width': '2',
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
      }),
    ],
  )

FlameUpIcon.props = {
  size: { type: Number, default: 24 },
}

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
  'flame-up': FlameUpIcon,
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
