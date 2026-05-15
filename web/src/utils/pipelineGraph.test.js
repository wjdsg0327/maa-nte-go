import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import {
  buildGraphLinks,
  getRootNodes,
  layoutGraph,
} from './pipelineGraph.js'

describe('pipeline graph helpers', () => {
  const nodes = {
    Start: { recognition: 'OCR', action: 'Click', next: ['Done'], interrupt: ['Abort'] },
    Done: { recognition: 'DirectHit', action: 'DoNothing', reverse: ['Retry'] },
    Abort: { recognition: 'TemplateMatch', action: 'StopTask' },
    Retry: { recognition: 'DirectHit', action: 'Click', on_error: ['Done'] },
  }

  it('builds all visible relation links', () => {
    assert.deepEqual(buildGraphLinks(nodes), [
      { from: 'Start', to: 'Done', type: 'next' },
      { from: 'Start', to: 'Abort', type: 'interrupt' },
      { from: 'Done', to: 'Retry', type: 'reverse' },
      { from: 'Retry', to: 'Done', type: 'error' },
    ])
  })

  it('detects roots from normal next relations only', () => {
    assert.deepEqual(getRootNodes(nodes), ['Start', 'Abort', 'Retry'])
  })

  it('lays graph nodes onto stable coordinates', () => {
    const layout = layoutGraph(nodes, { nodeWidth: 180, nodeHeight: 84, gapX: 80, gapY: 92 })

    assert.deepEqual(layout.positions.Start, { x: 32, y: 32 })
    assert.equal(layout.positions.Done.y, 208)
    assert.equal(layout.width >= 472, true)
    assert.equal(layout.height >= 324, true)
  })

  it('reads Maa string and NodeAttr relation targets', () => {
    const mixedNodes = {
      Start: { recognition: 'DirectHit', next: 'Middle' },
      Middle: { recognition: 'DirectHit', next: { name: 'End', jump_back: true } },
      End: { recognition: 'DirectHit', on_error: [{ name: 'Start' }] },
    }

    assert.deepEqual(buildGraphLinks(mixedNodes), [
      { from: 'Start', to: 'Middle', type: 'next' },
      { from: 'Middle', to: 'End', type: 'next' },
      { from: 'End', to: 'Start', type: 'error' },
    ])
    assert.deepEqual(getRootNodes(mixedNodes), ['Start'])
    assert.deepEqual(layoutGraph(mixedNodes).positions.Middle.y, 208)
  })
})
