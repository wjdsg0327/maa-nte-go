import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import {
  addRelationTarget,
  getRelationDisplayName,
  getRelationTargets,
  removeRelationTargetAt,
  renameRelationTarget,
  serializeRelationList,
} from './pipelineRelations.js'

describe('pipeline relation helpers', () => {
  it('normalizes Maa relation fields without losing NodeAttr data', () => {
    assert.deepEqual(getRelationTargets(undefined), [])
    assert.deepEqual(getRelationTargets('Confirm'), ['Confirm'])
    assert.deepEqual(getRelationTargets({ name: 'Confirm', jump_back: true }), [
      { name: 'Confirm', jump_back: true },
    ])
    assert.deepEqual(getRelationTargets(['A', { name: 'B', anchor: true }]), [
      'A',
      { name: 'B', anchor: true },
    ])
  })

  it('adds and removes selected next nodes from string-form fields', () => {
    assert.deepEqual(addRelationTarget('A', 'B'), ['A', 'B'])
    assert.deepEqual(addRelationTarget({ name: 'A', jump_back: true }, 'B'), [
      { name: 'A', jump_back: true },
      'B',
    ])
    assert.deepEqual(addRelationTarget(['A'], 'A'), ['A'])
    assert.equal(removeRelationTargetAt('A', 0), undefined)
    assert.equal(removeRelationTargetAt(['A', 'B'], 0), 'B')
  })

  it('renames and serializes targets while preserving single-value Maa syntax', () => {
    assert.equal(renameRelationTarget('Old', 'Old', 'New'), 'New')
    assert.deepEqual(renameRelationTarget({ name: 'Old', jump_back: true }, 'Old', 'New'), {
      name: 'New',
      jump_back: true,
    })
    assert.equal(serializeRelationList(['Only']), 'Only')
    assert.deepEqual(serializeRelationList(['A', 'B']), ['A', 'B'])
  })

  it('removes references to a deleted node regardless of relation shape', async () => {
    const { removeRelationTargetName } = await import('./pipelineRelations.js')

    assert.equal(removeRelationTargetName('Deleted', 'Deleted'), undefined)
    assert.deepEqual(removeRelationTargetName(['Keep', 'Deleted'], 'Deleted'), 'Keep')
    assert.deepEqual(
      removeRelationTargetName([{ name: 'Deleted', jump_back: true }, { name: 'Keep' }], 'Deleted'),
      { name: 'Keep' },
    )
  })

  it('extracts display names from plain, prefixed, anchor, and NodeAttr relations', () => {
    assert.equal(getRelationDisplayName('A'), 'A')
    assert.equal(getRelationDisplayName('[JumpBack]A'), 'A')
    assert.equal(getRelationDisplayName('[Anchor]Last'), 'Last')
    assert.equal(getRelationDisplayName({ name: 'A', jump_back: true }), 'A')
  })
})
