import { describe, expect, it } from 'vitest'
import { buildCommentTree } from '../../src/composables/buildCommentTree'

describe('buildCommentTree', () => {
  it('nests replies under their parent comments', () => {
    const tree = buildCommentTree([
      { id: '1', parentId: '', content: 'root' },
      { id: '2', parentId: '1', content: 'reply' },
      { id: '3', parentId: '2', content: 'nested reply' },
    ])

    expect(tree).toHaveLength(1)
    expect(tree[0].children[0].children[0]).toMatchObject({
      id: '3',
      content: 'nested reply',
    })
  })

  it('keeps orphaned replies visible as root comments', () => {
    const tree = buildCommentTree([{ id: '2', parentId: 'missing', content: 'orphan' }])

    expect(tree).toHaveLength(1)
    expect(tree[0]).toMatchObject({ id: '2', content: 'orphan' })
  })
})
