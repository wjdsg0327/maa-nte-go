import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import { cleanPipelineNodes } from './pipelineClean.js'

describe('pipeline cleaning', () => {
  it('drops empty values and legacy wait_freezes route fields before saving', () => {
    const cleaned = cleanPipelineNodes({
      Start: {
        recognition: 'DirectHit',
        action: 'DoNothing',
        next: [],
        wait_freezes: ['LegacyWait'],
        timeout: '',
      },
    })

    assert.deepEqual(cleaned, {
      Start: {
        recognition: 'DirectHit',
        action: 'DoNothing',
      },
    })
  })

  it('migrates legacy reverse route fields to Maa on_error', () => {
    const cleaned = cleanPipelineNodes({
      Start: {
        recognition: 'DirectHit',
        reverse: ['Retry'],
      },
      Retry: {
        recognition: 'DirectHit',
      },
    })

    assert.deepEqual(cleaned.Start, {
      recognition: 'DirectHit',
      on_error: 'Retry',
    })
  })

  it('combines legacy reverse with existing on_error without duplicating targets', () => {
    const cleaned = cleanPipelineNodes({
      Start: {
        recognition: 'DirectHit',
        on_error: ['Retry'],
        reverse: ['Retry', { name: 'Abort', jump_back: true }],
      },
    })

    assert.deepEqual(cleaned.Start, {
      recognition: 'DirectHit',
      on_error: ['Retry', { name: 'Abort', jump_back: true }],
    })
  })

  it('migrates legacy times to Maa repeat before saving', () => {
    const cleaned = cleanPipelineNodes({
      Start: {
        recognition: 'DirectHit',
        action: 'Click',
        times: 3,
      },
    })

    assert.deepEqual(cleaned.Start, {
      recognition: 'DirectHit',
      action: 'Click',
      repeat: 3,
    })
  })
})
