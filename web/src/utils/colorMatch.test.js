import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import {
  colorRowsFromBounds,
  hexToRgb,
  rgbToHex,
  serializeColorBounds,
  toleranceBounds,
} from './colorMatch.js'

describe('color match helpers', () => {
  it('converts between hex and RGB arrays', () => {
    assert.deepEqual(hexToRgb('#ff3300'), [255, 51, 0])
    assert.deepEqual(hexToRgb('0a141e'), [10, 20, 30])
    assert.equal(rgbToHex([255, 51, 0]), '#ff3300')
  })

  it('builds clamped lower and upper bounds from color and tolerance', () => {
    assert.deepEqual(toleranceBounds('#ff3300', 10), {
      lower: [245, 41, 0],
      upper: [255, 61, 10],
    })
  })

  it('serializes one color as flat bounds and multiple colors as nested bounds', () => {
    assert.deepEqual(serializeColorBounds([{ color: '#ff3300', tolerance: 10 }]), {
      lower: [245, 41, 0],
      upper: [255, 61, 10],
    })
    assert.deepEqual(serializeColorBounds([
      { color: '#ff0000', tolerance: 0 },
      { color: '#00ff00', tolerance: 5 },
    ]), {
      lower: [[255, 0, 0], [0, 250, 0]],
      upper: [[255, 0, 0], [5, 255, 5]],
    })
  })

  it('restores editable rows from existing lower and upper values', () => {
    assert.deepEqual(colorRowsFromBounds([245, 41, 0], [255, 61, 10]), [
      { color: '#fa3305', tolerance: 10 },
    ])
    assert.deepEqual(colorRowsFromBounds([[255, 0, 0], [0, 250, 0]], [[255, 0, 0], [5, 255, 5]]), [
      { color: '#ff0000', tolerance: 0 },
      { color: '#03fd03', tolerance: 3 },
    ])
  })
})
