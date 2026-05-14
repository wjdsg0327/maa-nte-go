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
    Retry: { recognition: 'DirectHit', action: 'Click', wait_freezes: ['Done'] },
  }

  it('builds all visible relation links', () => {
    assert.deepEqual(buildGraphLinks(nodes), [
      { from: 'Start', to: 'Done', type: 'next' },
      { from: 'Start', to: 'Abort', type: 'interrupt' },
      { from: 'Done', to: 'Retry', type: 'reverse' },
      { from: 'Retry', to: 'Done', type: 'wait' },
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
})
