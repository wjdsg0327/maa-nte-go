import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import {
  appendListItem,
  formatNumberOrObject,
  formatNumberOrList,
  formatRoiValue,
  normalizeStringList,
  parseNumberOrObjectValue,
  parseRectArrayValue,
  parseRoiValue,
  serializeNumberOrList,
  serializeStringOrList,
} from './pipelineFields.js'

describe('pipeline field helpers', () => {
  it('preserves Maa string-or-list template fields', () => {
    assert.deepEqual(normalizeStringList(undefined), [])
    assert.deepEqual(normalizeStringList('confirm.png'), ['confirm.png'])
    assert.deepEqual(normalizeStringList(['a.png', 'b.png']), ['a.png', 'b.png'])

    assert.equal(serializeStringOrList([]), undefined)
    assert.equal(serializeStringOrList(['confirm.png']), 'confirm.png')
    assert.deepEqual(serializeStringOrList(['a.png', 'b.png']), ['a.png', 'b.png'])
  })

  it('appends unique trimmed template entries', () => {
    assert.deepEqual(appendListItem(['a.png'], ' b.png '), ['a.png', 'b.png'])
    assert.deepEqual(appendListItem(['a.png'], 'a.png'), ['a.png'])
    assert.deepEqual(appendListItem(['a.png'], ''), ['a.png'])
  })

  it('formats and serializes threshold as a number or number list', () => {
    assert.equal(formatNumberOrList(undefined), '')
    assert.equal(formatNumberOrList(0.7), '0.7')
    assert.equal(formatNumberOrList([0.7, 0.8]), '[0.7,0.8]')

    assert.equal(serializeNumberOrList(''), undefined)
    assert.equal(serializeNumberOrList('0.72'), 0.72)
    assert.deepEqual(serializeNumberOrList('[0.7, 0.8]'), [0.7, 0.8])
    assert.equal(serializeNumberOrList('bad'), undefined)
  })

  it('formats and parses Maa uint-or-object wait fields', () => {
    assert.equal(formatNumberOrObject(undefined), '')
    assert.equal(formatNumberOrObject(500), '500')
    assert.equal(formatNumberOrObject({ time: 500, threshold: 0.95 }), '{"time":500,"threshold":0.95}')

    assert.deepEqual(parseNumberOrObjectValue('500', '静止等待'), {
      value: 500,
      error: '',
    })
    assert.deepEqual(parseNumberOrObjectValue('{"time":500,"threshold":0.95}', '静止等待'), {
      value: { time: 500, threshold: 0.95 },
      error: '',
    })
    assert.match(parseNumberOrObjectValue('-1', '静止等待').error, /静止等待/)
    assert.match(parseNumberOrObjectValue('[500]', '静止等待').error, /静止等待/)
  })

  it('formats and parses ROI values as arrays or string references', () => {
    assert.equal(formatRoiValue(undefined), '')
    assert.equal(formatRoiValue([983, 260, 169, 388]), '[983,260,169,388]')
    assert.equal(formatRoiValue('CreditIcon'), 'CreditIcon')

    assert.deepEqual(parseRoiValue('[983, 260, 169, 388]'), {
      value: [983, 260, 169, 388],
      error: '',
    })
    assert.deepEqual(parseRoiValue('CreditIcon'), {
      value: 'CreditIcon',
      error: '',
    })
    assert.deepEqual(parseRoiValue('[Anchor]CreditIcon'), {
      value: '[Anchor]CreditIcon',
      error: '',
    })
    assert.match(parseRoiValue('983,260.169,388').error, /ROI/)
  })

  it('rejects ROI arrays that are not four integers', () => {
    assert.match(parseRectArrayValue('[983, 260.169, 388]').error, /4/)
    assert.deepEqual(parseRectArrayValue('[983, 260, 169, 388]'), {
      value: [983, 260, 169, 388],
      error: '',
    })
  })
})
