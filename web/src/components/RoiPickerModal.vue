<template>
  <div class="modal-overlay roi-overlay" @click.self="$emit('close')">
    <div class="modal roi-modal">
      <div class="modal-titlebar">
        <div>
          <p class="eyebrow">ROI PICKER</p>
          <h2>框选 ROI</h2>
        </div>
        <button class="btn btn-secondary btn-sm" type="button" @click="$emit('close')">关闭</button>
      </div>

      <div class="roi-stage">
        <div
          ref="surfaceRef"
          class="roi-image-surface"
          @pointerdown="startDrag"
          @pointermove="moveDrag"
          @pointerup="endDrag"
          @pointercancel="endDrag"
        >
          <img
            ref="imageRef"
            class="roi-image"
            :src="screenshot.dataUrl"
            :width="screenshot.width"
            :height="screenshot.height"
            @load="syncSurface"
            alt=""
            draggable="false"
          >
          <div
            v-if="selectionBox.width > 0 && selectionBox.height > 0"
            class="roi-selection"
            :style="selectionStyle"
          />
        </div>
      </div>

      <div class="roi-readout">
        <span>{{ roiText }}</span>
        <span>{{ screenshot.width }} x {{ screenshot.height }}</span>
      </div>

      <div class="modal-actions">
        <button class="btn btn-secondary" type="button" @click="$emit('close')">取消</button>
        <button class="btn btn-primary" type="button" :disabled="!canApply" @click="applyRoi">应用 ROI</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { dragToRoi, normalizeDragBox } from '../utils/roiPicker.js'

const props = defineProps({
  screenshot: { type: Object, required: true },
  initialRoi: { type: Array, default: () => [] },
})

const emit = defineEmits(['close', 'apply'])

const surfaceRef = ref(null)
const imageRef = ref(null)
const dragging = ref(false)
const start = ref({ x: 0, y: 0 })
const end = ref({ x: 0, y: 0 })
const surfaceSize = ref({ width: 0, height: 0 })

const syncSurface = () => {
  if (!surfaceRef.value) return
  const rect = surfaceRef.value.getBoundingClientRect()
  surfaceSize.value = { width: rect.width, height: rect.height }
}

const pointFromEvent = (event) => {
  syncSurface()
  const rect = surfaceRef.value.getBoundingClientRect()
  return {
    x: Math.min(Math.max(event.clientX - rect.left, 0), rect.width),
    y: Math.min(Math.max(event.clientY - rect.top, 0), rect.height),
  }
}

const startDrag = (event) => {
  if (!surfaceRef.value) return
  event.preventDefault()
  const point = pointFromEvent(event)
  start.value = point
  end.value = point
  dragging.value = true
  surfaceRef.value.setPointerCapture?.(event.pointerId)
}

const moveDrag = (event) => {
  if (!dragging.value) return
  event.preventDefault()
  end.value = pointFromEvent(event)
}

const endDrag = (event) => {
  if (!dragging.value) return
  end.value = pointFromEvent(event)
  dragging.value = false
  surfaceRef.value?.releasePointerCapture?.(event.pointerId)
}

const selectionBox = computed(() => normalizeDragBox({
  start: start.value,
  end: end.value,
  bounds: surfaceSize.value,
}))

const roi = computed(() => dragToRoi({
  start: start.value,
  end: end.value,
  rendered: surfaceSize.value,
  natural: { width: props.screenshot.width, height: props.screenshot.height },
}))

const selectionStyle = computed(() => ({
  left: `${selectionBox.value.x}px`,
  top: `${selectionBox.value.y}px`,
  width: `${selectionBox.value.width}px`,
  height: `${selectionBox.value.height}px`,
}))

const canApply = computed(() => roi.value[2] > 0 && roi.value[3] > 0)
const roiText = computed(() => `[${roi.value.join(',')}]`)

const applyRoi = () => {
  if (!canApply.value) return
  emit('apply', roi.value)
}

const seedInitialRoi = () => {
  syncSurface()
  if (!Array.isArray(props.initialRoi) || props.initialRoi.length !== 4) return
  if (!surfaceSize.value.width || !surfaceSize.value.height) return
  const scaleX = surfaceSize.value.width / props.screenshot.width
  const scaleY = surfaceSize.value.height / props.screenshot.height
  const [x, y, width, height] = props.initialRoi
  start.value = { x: x * scaleX, y: y * scaleY }
  end.value = { x: (x + width) * scaleX, y: (y + height) * scaleY }
}

watch(() => props.initialRoi, seedInitialRoi, { immediate: true })

onMounted(() => {
  window.addEventListener('resize', syncSurface)
  requestAnimationFrame(seedInitialRoi)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncSurface)
})
</script>
