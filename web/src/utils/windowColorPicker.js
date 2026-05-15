function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max)
}

function toPositiveSize(value) {
  const number = Number(value)
  return Number.isFinite(number) && number > 0 ? number : 0
}

export function pointToNaturalPixel({ point, rendered, natural }) {
  const renderedWidth = toPositiveSize(rendered?.width)
  const renderedHeight = toPositiveSize(rendered?.height)
  const naturalWidth = toPositiveSize(natural?.width)
  const naturalHeight = toPositiveSize(natural?.height)
  if (!renderedWidth || !renderedHeight || !naturalWidth || !naturalHeight) {
    return [0, 0]
  }

  const renderedX = clamp(Number(point?.x) || 0, 0, renderedWidth)
  const renderedY = clamp(Number(point?.y) || 0, 0, renderedHeight)
  const x = Math.floor((renderedX / renderedWidth) * naturalWidth)
  const y = Math.floor((renderedY / renderedHeight) * naturalHeight)

  return [
    clamp(x, 0, naturalWidth - 1),
    clamp(y, 0, naturalHeight - 1),
  ]
}
