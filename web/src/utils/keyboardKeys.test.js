import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import {
  filterKeyboardKeys,
  formatKeyCode,
  normalizeKeyCodes,
  parseKeyboardKeyInput,
  serializeKeyCodes,
} from './keyboardKeys.js'

describe('keyboard key helpers', () => {
  it('parses common key names, aliases, labels, and raw key codes', () => {
    assert.equal(parseKeyboardKeyInput('esc'), 27)
    assert.equal(parseKeyboardKeyInput('Escape (27)'), 27)
    assert.equal(parseKeyboardKeyInput('回车'), 13)
    assert.equal(parseKeyboardKeyInput('F5'), 116)
    assert.equal(parseKeyboardKeyInput('A'), 65)
    assert.equal(parseKeyboardKeyInput('32'), 32)
    assert.equal(parseKeyboardKeyInput('not-a-key'), undefined)
  })

  it('searches by name, alias, and key code', () => {
    assert.equal(filterKeyboardKeys('空格')[0].code, 32)
    assert.equal(filterKeyboardKeys('enter')[0].code, 13)
    assert.equal(filterKeyboardKeys('116')[0].code, 116)
  })

  it('formats and serializes key selections for Maa fields', () => {
    assert.equal(formatKeyCode(27), 'Esc (27)')
    assert.deepEqual(normalizeKeyCodes([27, '13', 'bad', 27]), [27, 13])
    assert.equal(serializeKeyCodes([], true), undefined)
    assert.equal(serializeKeyCodes([27], true), 27)
    assert.deepEqual(serializeKeyCodes([27, 13], true), [27, 13])
    assert.equal(serializeKeyCodes([27, 13], false), 27)
  })
})
