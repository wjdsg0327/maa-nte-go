import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import { dragToRoi, normalizeDragBox } from './roiPicker.js'

describe('ROI picker helpers', () => {
  it('scales a drag rectangle from rendered pixels to natural screenshot pixels', () => {
    const roi = dragToRoi({
      start: { x: 10, y: 20 },
      end: { x: 110, y: 80 },
      rendered: { width: 640, height: 360 },
      natural: { width: 1280, height: 720 },
    })

    assert.deepEqual(roi, [20, 40, 200, 120])
  })

  it('normalizes reverse dragging and clamps to image bounds', () => {
    const box = normalizeDragBox({
      start: { x: 300, y: 200 },
      end: { x: -20, y: 50 },
      bounds: { width: 320, height: 240 },
    })

    assert.deepEqual(box, { x: 0, y: 50, width: 300, height: 150 })
  })
})
