import { describe, expect, it } from 'vitest'
import { cn } from '../../src/lib/utils'

describe('cn', () => {
  it('merges conditional classes and resolves Tailwind conflicts', () => {
    expect(cn('px-2 text-sm', ['px-4', false && 'hidden'], { 'font-bold': true })).toBe(
      'text-sm px-4 font-bold',
    )
  })
})
