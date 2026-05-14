const LINK_FIELDS = [
  ['next', 'next'],
  ['interrupt', 'interrupt'],
  ['wait_freezes', 'wait'],
  ['reverse', 'reverse'],
]

export function buildGraphLinks(nodes) {
  const links = []

  for (const [name, node] of Object.entries(nodes || {})) {
    for (const [field, type] of LINK_FIELDS) {
      const targets = Array.isArray(node?.[field]) ? node[field] : []
      for (const target of targets) {
        if (nodes?.[target]) links.push({ from: name, to: target, type })
      }
    }
  }

  return links
}

export function getRootNodes(nodes) {
  const nextChildren = new Set()

  for (const node of Object.values(nodes || {})) {
    if (Array.isArray(node?.next)) {
      node.next.forEach((target) => nextChildren.add(target))
    }
  }

  return Object.keys(nodes || {}).filter((name) => !nextChildren.has(name))
}

export function layoutGraph(nodes, options = {}) {
  const nodeWidth = options.nodeWidth ?? 180
  const nodeHeight = options.nodeHeight ?? 84
  const gapX = options.gapX ?? 80
  const gapY = options.gapY ?? 92
  const padding = options.padding ?? 32
  const positions = {}
  const visited = new Set()
  let cursorX = padding

  const visit = (name, x, y) => {
    if (visited.has(name) || !nodes?.[name]) return x
    visited.add(name)
    positions[name] = { x, y }

    const children = Array.isArray(nodes[name].next)
      ? nodes[name].next.filter((child) => nodes[child])
      : []

    if (children.length === 0) return x + nodeWidth + gapX

    let childX = x
    for (const child of children) {
      childX = visit(child, childX, y + nodeHeight + gapY)
    }

    return Math.max(childX, x + nodeWidth + gapX)
  }

  for (const root of getRootNodes(nodes)) {
    cursorX = visit(root, cursorX, padding)
  }

  for (const name of Object.keys(nodes || {})) {
    if (!visited.has(name)) {
      positions[name] = { x: cursorX, y: padding }
      cursorX += nodeWidth + gapX
    }
  }

  const values = Object.values(positions)
  const width = values.length === 0
    ? 720
    : Math.max(720, Math.max(...values.map((pos) => pos.x)) + nodeWidth + padding)
  const height = values.length === 0
    ? 420
    : Math.max(420, Math.max(...values.map((pos) => pos.y)) + nodeHeight + padding)

  return {
    height,
    links: buildGraphLinks(nodes),
    positions,
    width,
  }
}
