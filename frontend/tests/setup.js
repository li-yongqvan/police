import { vi } from 'vitest'

window.matchMedia ??= vi.fn().mockImplementation((query) => ({
  matches: false,
  media: query,
  onchange: null,
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
  dispatchEvent: vi.fn(),
}))

window.ResizeObserver ??= class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

window.IntersectionObserver ??= class IntersectionObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}
