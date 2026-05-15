import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import { pointToNaturalPixel } from './windowColorPicker.js'

describe('window color picker helpers', () => {
  it('maps a displayed screenshot point to the natural pixel coordinate', () => {
    const point = pointToNaturalPixel({
      point: { x: 50, y: 25 },
      rendered: { width: 100, height: 50 },
      natural: { width: 200, height: 100 },
    })

    assert.deepEqual(point, [100, 50])
  })

  it('clamps click coordinates to the screenshot bounds', () => {
    const point = pointToNaturalPixel({
      point: { x: 120, y: -10 },
      rendered: { width: 100, height: 50 },
      natural: { width: 200, height: 100 },
    })

    assert.deepEqual(point, [199, 0])
  })
})
