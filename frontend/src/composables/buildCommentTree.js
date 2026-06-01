/** 将扁平评论列表构建为树（用于 GxCommentTree） */
export function buildCommentTree(flat = []) {
  const byId = new Map()
  const roots = []

  flat.forEach((item) => {
    byId.set(item.id, { ...item, children: [] })
  })

  byId.forEach((node) => {
    const parentId = node.parentId
    if (parentId && byId.has(parentId)) {
      byId.get(parentId).children.push(node)
    } else {
      roots.push(node)
    }
  })

  return roots
}
