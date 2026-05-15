function clampChannel(value) {
  const number = Number(value)
  if (!Number.isFinite(number)) return 0
  return Math.min(255, Math.max(0, Math.round(number)))
}

function normalizeRgb(value) {
  if (!Array.isArray(value) || value.length < 3) return [0, 0, 0]
  return value.slice(0, 3).map(clampChannel)
}

function normalizeColorList(value) {
  if (!Array.isArray(value)) return []
  if (value.length === 0) return []
  if (Array.isArray(value[0])) return value.map(normalizeRgb)
  return [normalizeRgb(value)]
}

export function hexToRgb(hex) {
  const normalized = String(hex || '').trim().replace(/^#/, '')
  if (!/^[0-9a-fA-F]{6}$/.test(normalized)) return [0, 0, 0]
  return [
    parseInt(normalized.slice(0, 2), 16),
    parseInt(normalized.slice(2, 4), 16),
    parseInt(normalized.slice(4, 6), 16),
  ]
}

export function rgbToHex(rgb) {
  return `#${normalizeRgb(rgb).map((channel) => channel.toString(16).padStart(2, '0')).join('')}`
}

export function toleranceBounds(color, tolerance) {
  const rgb = Array.isArray(color) ? normalizeRgb(color) : hexToRgb(color)
  const tol = Math.max(0, Math.round(Number(tolerance) || 0))
  return {
    lower: rgb.map((channel) => clampChannel(channel - tol)),
    upper: rgb.map((channel) => clampChannel(channel + tol)),
  }
}

export function serializeColorBounds(rows) {
  const bounds = (rows || [])
    .map((row) => toleranceBounds(row.color, row.tolerance))
    .filter((row) => Array.isArray(row.lower) && Array.isArray(row.upper))

  if (bounds.length === 0) return { lower: undefined, upper: undefined }
  if (bounds.length === 1) return bounds[0]
  return {
    lower: bounds.map((row) => row.lower),
    upper: bounds.map((row) => row.upper),
  }
}

export function colorRowsFromBounds(lowerValue, upperValue) {
  const lowers = normalizeColorList(lowerValue)
  const uppers = normalizeColorList(upperValue)
  const count = Math.min(lowers.length, uppers.length)
  const rows = []

  for (let index = 0; index < count; index += 1) {
    const lower = lowers[index]
    const upper = uppers[index]
    const color = lower.map((channel, channelIndex) => clampChannel((channel + upper[channelIndex]) / 2))
    const tolerance = Math.max(...upper.map((channel, channelIndex) => Math.abs(channel - lower[channelIndex]) / 2))
    rows.push({
      color: rgbToHex(color),
      tolerance: Math.round(tolerance),
    })
  }

  return rows
}
