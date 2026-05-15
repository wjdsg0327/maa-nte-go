function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max)
}

function toPositiveSize(value) {
  const number = Number(value)
  return Number.isFinite(number) && number > 0 ? number : 0
}

export function normalizeDragBox({ start, end, bounds }) {
  const maxWidth = toPositiveSize(bounds?.width)
  const maxHeight = toPositiveSize(bounds?.height)
  const startX = clamp(Number(start?.x) || 0, 0, maxWidth)
  const startY = clamp(Number(start?.y) || 0, 0, maxHeight)
  const endX = clamp(Number(end?.x) || 0, 0, maxWidth)
  const endY = clamp(Number(end?.y) || 0, 0, maxHeight)
  const x = Math.min(startX, endX)
  const y = Math.min(startY, endY)

  return {
    x,
    y,
    width: Math.abs(endX - startX),
    height: Math.abs(endY - startY),
  }
}

export function dragToRoi({ start, end, rendered, natural }) {
  const renderedWidth = toPositiveSize(rendered?.width)
  const renderedHeight = toPositiveSize(rendered?.height)
  const naturalWidth = toPositiveSize(natural?.width)
  const naturalHeight = toPositiveSize(natural?.height)
  if (!renderedWidth || !renderedHeight || !naturalWidth || !naturalHeight) {
    return [0, 0, 0, 0]
  }

  const box = normalizeDragBox({
    start,
    end,
    bounds: { width: renderedWidth, height: renderedHeight },
  })
  const scaleX = naturalWidth / renderedWidth
  const scaleY = naturalHeight / renderedHeight

  return [
    Math.round(box.x * scaleX),
    Math.round(box.y * scaleY),
    Math.round(box.width * scaleX),
    Math.round(box.height * scaleY),
  ]
}
