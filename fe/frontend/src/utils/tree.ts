// 树形数据工具

interface TreeItem {
  id: number
  parentId: number
  children?: TreeItem[]
  [key: string]: any
}

// 列表转树形
export function listToTree<T extends TreeItem>(
  list: T[],
  rootId = 0
): T[] {
  const map = new Map<number, T>()
  const result: T[] = []

  // 创建映射
  list.forEach(item => {
    map.set(item.id, { ...item, children: [] })
  })

  // 构建树
  list.forEach(item => {
    const node = map.get(item.id)!
    if (item.parentId === rootId) {
      result.push(node)
    } else {
      const parent = map.get(item.parentId)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push(node)
      }
    }
  })

  return result
}

// 树形转列表
export function treeToList<T extends TreeItem>(tree: T[]): T[] {
  const result: T[] = []

  function traverse(nodes: T[]) {
    nodes.forEach(node => {
      const { children, ...rest } = node
      result.push(rest as T)
      if (children && children.length > 0) {
        traverse(children as T[])
      }
    })
  }

  traverse(tree)
  return result
}

// 查找节点
export function findNode<T extends TreeItem>(
  tree: T[],
  id: number
): T | undefined {
  for (const node of tree) {
    if (node.id === id) {
      return node
    }
    if (node.children && node.children.length > 0) {
      const found = findNode(node.children as T[], id)
      if (found) return found as T
    }
  }
  return undefined
}

// 获取节点路径
export function getNodePath<T extends TreeItem>(
  tree: T[],
  id: number
): T[] {
  const path: T[] = []

  function traverse(nodes: T[], targetId: number): boolean {
    for (const node of nodes) {
      path.push(node)
      if (node.id === targetId) {
        return true
      }
      if (node.children && node.children.length > 0) {
        if (traverse(node.children as T[], targetId)) {
          return true
        }
      }
      path.pop()
    }
    return false
  }

  traverse(tree, id)
  return path
}

// 获取所有叶子节点
export function getLeafNodes<T extends TreeItem>(tree: T[]): T[] {
  const leaves: T[] = []

  function traverse(nodes: T[]) {
    nodes.forEach(node => {
      if (!node.children || node.children.length === 0) {
        leaves.push(node)
      } else {
        traverse(node.children as T[])
      }
    })
  }

  traverse(tree)
  return leaves
}

// 过滤树形数据
export function filterTree<T extends TreeItem>(
  tree: T[],
  predicate: (node: T) => boolean
): T[] {
  return tree
    .filter(node => {
      if (node.children && node.children.length > 0) {
        node.children = filterTree(node.children as T[], predicate)
      }
      return predicate(node) || (node.children && node.children.length > 0)
    })
    .map(node => ({ ...node }))
}
