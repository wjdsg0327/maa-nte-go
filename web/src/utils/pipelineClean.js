import {
  addRelationTarget,
  getRelationDisplayName,
  getRelationTargets,
  serializeRelationList,
} from './pipelineRelations.js'

const LEGACY_ROUTE_FIELDS = new Set(['wait_freezes'])

function isEmptyValue(value) {
  return value === undefined
    || value === ''
    || value === null
    || (Array.isArray(value) && value.length === 0)
}

function mergeRelations(baseValue, extraValue) {
  let merged = serializeRelationList(baseValue)

  for (const item of getRelationTargets(extraValue)) {
    merged = addRelationTarget(merged, getRelationDisplayName(item))
    const current = getRelationTargets(merged)
    const lastIndex = current.length - 1
    if (typeof item !== 'string' && lastIndex >= 0 && getRelationDisplayName(current[lastIndex]) === item.name) {
      current[lastIndex] = item
      merged = serializeRelationList(current)
    }
  }

  return merged
}

export function cleanPipelineNodes(nodes) {
  const cleanNodes = {}

  for (const [name, node] of Object.entries(nodes || {})) {
    const cleanNode = {}

    for (const [key, value] of Object.entries(node || {})) {
      if (isEmptyValue(value) || LEGACY_ROUTE_FIELDS.has(key) || key === 'reverse' || key === 'times') continue
      cleanNode[key] = value
    }

    if (cleanNode.repeat === undefined && node?.times !== undefined && !isEmptyValue(node.times)) {
      cleanNode.repeat = node.times
    }

    if (node?.reverse !== undefined) {
      const onError = mergeRelations(cleanNode.on_error, node.reverse)
      if (onError === undefined) {
        delete cleanNode.on_error
      } else {
        cleanNode.on_error = onError
      }
    }

    cleanNodes[name] = cleanNode
  }

  return cleanNodes
}
