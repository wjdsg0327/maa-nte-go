const PREFIX_PATTERN = /^\[([A-Za-z_]+)\](.*)$/

function normalizeStringTarget(value) {
  const trimmed = typeof value === 'string' ? value.trim() : ''
  return trimmed || undefined
}

function normalizeNodeAttr(value) {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return undefined
  const name = normalizeStringTarget(value.name)
  if (!name) return undefined
  return { ...value, name }
}

function normalizeRelationItem(value) {
  if (typeof value === 'string') return normalizeStringTarget(value)
  return normalizeNodeAttr(value)
}

export function getRelationTargets(value) {
  const source = Array.isArray(value) ? value : [value]
  return source
    .map(normalizeRelationItem)
    .filter((item) => item !== undefined)
}

export function serializeRelationList(items) {
  const normalized = getRelationTargets(items)
  if (normalized.length === 0) return undefined
  return normalized.length === 1 ? normalized[0] : normalized
}

function splitRelationPrefix(value) {
  const text = normalizeStringTarget(value)
  if (!text) return { prefix: '', name: '' }

  const match = text.match(PREFIX_PATTERN)
  if (!match) return { prefix: '', name: text }
  return { prefix: `[${match[1]}]`, name: match[2].trim() }
}

function relationName(value) {
  if (typeof value === 'string') return splitRelationPrefix(value).name
  return normalizeStringTarget(value?.name) || ''
}

export function getRelationDisplayName(value) {
  return relationName(value)
}

export function getRelationGraphTargetName(value) {
  if (typeof value === 'string') {
    const { prefix, name } = splitRelationPrefix(value)
    return prefix.toLowerCase() === '[anchor]' ? '' : name
  }

  if (value?.anchor === true) return ''
  return relationName(value)
}

export function addRelationTarget(value, targetName) {
  const target = normalizeStringTarget(targetName)
  const current = getRelationTargets(value)
  if (!target) return serializeRelationList(current)

  const exists = current.some((item) => relationName(item) === target)
  return exists ? value : serializeRelationList([...current, target])
}

export function removeRelationTargetAt(value, index) {
  const current = getRelationTargets(value)
  if (index < 0 || index >= current.length) return serializeRelationList(current)
  current.splice(index, 1)
  return serializeRelationList(current)
}

export function removeRelationTargetName(value, targetName) {
  const target = normalizeStringTarget(targetName)
  if (!target) return serializeRelationList(value)
  return serializeRelationList(getRelationTargets(value).filter((item) => relationName(item) !== target))
}

export function renameRelationTarget(value, oldName, newName) {
  const oldTarget = normalizeStringTarget(oldName)
  const newTarget = normalizeStringTarget(newName)
  if (!oldTarget || !newTarget) return serializeRelationList(value)

  const updated = getRelationTargets(value).map((item) => {
    if (relationName(item) !== oldTarget) return item

    if (typeof item === 'string') {
      const { prefix } = splitRelationPrefix(item)
      return `${prefix}${newTarget}`
    }

    return { ...item, name: newTarget }
  })

  return serializeRelationList(updated)
}
