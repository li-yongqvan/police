import { beforeEach, describe, expect, it, vi } from 'vitest'

const replace = vi.fn()
let route = { query: {} }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => ({ replace }),
}))

describe('useFeedSort', () => {
  beforeEach(() => {
    replace.mockClear()
    route = { query: {} }
  })

  it('falls back to hot when the query value is unknown', async () => {
    route = { query: { sort: 'unknown' } }
    const { useFeedSort } = await import('../../src/composables/useFeedSort')

    expect(useFeedSort().sort.value).toBe('hot')
  })

  it('updates the route query while omitting the default sort', async () => {
    route = { query: { page: '2', sort: 'new' } }
    const { useFeedSort } = await import('../../src/composables/useFeedSort')
    const { sort } = useFeedSort()

    sort.value = 'featured'
    expect(replace).toHaveBeenLastCalledWith({ query: { page: '2', sort: 'featured' } })

    sort.value = 'hot'
    expect(replace).toHaveBeenLastCalledWith({ query: { page: '2', sort: undefined } })
  })
})
