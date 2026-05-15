export function normalizeStringList(value) {
  if (Array.isArray(value)) {
    return value
      .map((item) => String(item).trim())
      .filter(Boolean)
  }
  if (typeof value === 'string') {
    const trimmed = value.trim()
    return trimmed ? [trimmed] : []
  }
  return []
}

export function serializeStringOrList(items) {
  const normalized = normalizeStringList(items)
  if (normalized.length === 0) return undefined
  if (normalized.length === 1) return normalized[0]
  return normalized
}

export function appendListItem(items, item) {
  const normalized = normalizeStringList(items)
  const next = typeof item === 'string' ? item.trim() : ''
  if (!next || normalized.includes(next)) return normalized
  return [...normalized, next]
}

export function formatNumberOrList(value) {
  if (Array.isArray(value)) return JSON.stringify(value)
  if (typeof value === 'number' && Number.isFinite(value)) return String(value)
  if (typeof value === 'string') return value
  return ''
}

export function serializeNumberOrList(value) {
  const input = typeof value === 'string' ? value.trim() : String(value ?? '').trim()
  if (!input) return undefined

  if (input.startsWith('[')) {
    try {
      const parsed = JSON.parse(input)
      if (!Array.isArray(parsed)) return undefined
      const numbers = parsed.map((item) => Number(item))
      return numbers.every(Number.isFinite) ? numbers : undefined
    } catch (e) {
      return undefined
    }
  }

  const number = Number(input)
  return Number.isFinite(number) ? number : undefined
}

export function formatNumberOrObject(value) {
  if (typeof value === 'number' && Number.isFinite(value)) return String(value)
  if (value && typeof value === 'object' && !Array.isArray(value)) return JSON.stringify(value)
  if (typeof value === 'string') return value
  return ''
}

export function parseNumberOrObjectValue(value, label = '字段') {
  const input = textValue(value)
  if (!input) return { value: undefined, error: '' }

  if (input.startsWith('{')) {
    try {
      const parsed = JSON.parse(input)
      if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
        return { value: undefined, error: `${label} 需要填写非负整数或 JSON 对象` }
      }
      return { value: parsed, error: '' }
    } catch (e) {
      return { value: undefined, error: `${label} JSON 格式不正确` }
    }
  }

  const number = Number(input)
  if (!Number.isInteger(number) || number < 0) {
    return { value: undefined, error: `${label} 需要填写非负整数或 JSON 对象` }
  }
  return { value: number, error: '' }
}

function rectError(label) {
  return `${label} 需要填写 [x,y,w,h]，4 个值都必须是整数`
}

function textValue(value) {
  return typeof value === 'string' ? value.trim() : String(value ?? '').trim()
}

function looksLikeUnbracketedRect(input) {
  return input.includes(',') && /^[\d\s+\-.,]+$/.test(input)
}

export function formatRoiValue(value) {
  if (Array.isArray(value)) return JSON.stringify(value)
  if (typeof value === 'string') return value
  return ''
}

export function parseRectArrayValue(value, label = 'ROI') {
  const input = textValue(value)
  if (!input) return { value: undefined, error: '' }

  if (!input.startsWith('[')) {
    return { value: undefined, error: rectError(label) }
  }

  let parsed
  try {
    parsed = JSON.parse(input)
  } catch (e) {
    return { value: undefined, error: rectError(label) }
  }

  if (!Array.isArray(parsed) || parsed.length !== 4) {
    return { value: undefined, error: rectError(label) }
  }

  const numbers = parsed.map((item) => Number(item))
  if (!numbers.every((item) => Number.isFinite(item) && Number.isInteger(item))) {
    return { value: undefined, error: rectError(label) }
  }

  return { value: numbers, error: '' }
}

export function parseRoiValue(value) {
  const input = textValue(value)
  if (!input) return { value: undefined, error: '' }

  if (input.startsWith('[Anchor]')) {
    return { value: input, error: '' }
  }

  if (input.startsWith('[')) {
    return parseRectArrayValue(input, 'ROI')
  }

  if (looksLikeUnbracketedRect(input)) {
    return { value: undefined, error: 'ROI 需要使用 [x,y,w,h] 格式，不能写成逗号分隔文本' }
  }

  return { value: input, error: '' }
}
