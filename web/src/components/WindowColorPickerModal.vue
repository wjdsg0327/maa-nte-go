<template>
  <div class="modal-overlay roi-overlay" @click.self="$emit('close')">
    <div class="modal roi-modal">
      <div class="modal-titlebar">
        <div>
          <p class="eyebrow">WINDOW COLOR PICKER</p>
          <h2>窗口取色</h2>
        </div>
        <button class="btn btn-secondary btn-sm" type="button" @click="$emit('close')">关闭</button>
      </div>

      <div class="roi-stage">
        <div
          ref="surfaceRef"
          class="roi-image-surface color-pick-surface"
          @click="pickColor"
        >
          <img
            ref="imageRef"
            class="roi-image"
            :src="screenshot.dataUrl"
            :width="screenshot.width"
            :height="screenshot.height"
            @load="handleImageLoad"
            alt=""
            draggable="false"
          >
          <div
            v-if="picked"
            class="color-pick-crosshair"
            :style="crosshairStyle"
          />
        </div>
      </div>

      <div class="color-pick-readout">
        <span
          class="color-pick-preview"
          :style="{ background: picked?.color || '#ffffff' }"
        />
        <div>
          <strong>{{ picked?.color || '点击截图取色' }}</strong>
          <span v-if="picked">RGB {{ picked.rgb.join(', ') }} · {{ picked.point.join(',') }}</span>
          <span v-else>{{ screenshot.width }} x {{ screenshot.height }}</span>
        </div>
      </div>

      <div v-if="error" class="field-error">{{ error }}</div>

      <div class="modal-actions">
        <button class="btn btn-secondary" type="button" @click="$emit('close')">取消</button>
        <button class="btn btn-primary" type="button" :disabled="!picked" @click="applyColor">应用颜色</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { rgbToHex } from '../utils/colorMatch.js'
import { pointToNaturalPixel } from '../utils/windowColorPicker.js'

const props = defineProps({
  screenshot: { type: Object, required: true },
})

const emit = defineEmits(['close', 'apply'])

const surfaceRef = ref(null)
const imageRef = ref(null)
const canvasRef = ref(null)
const surfaceSize = ref({ width: 0, height: 0 })
const picked = ref(null)
const pickedDisplayPoint = ref({ x: 0, y: 0 })
const error = ref('')

const syncSurface = () => {
  if (!surfaceRef.value) return
  const rect = surfaceRef.value.getBoundingClientRect()
  surfaceSize.value = { width: rect.width, height: rect.height }
}

const prepareCanvas = () => {
  if (!imageRef.value) return
  const canvas = document.createElement('canvas')
  canvas.width = props.screenshot.width
  canvas.height = props.screenshot.height
  const context = canvas.getContext('2d')
  context.drawImage(imageRef.value, 0, 0, canvas.width, canvas.height)
  canvasRef.value = canvas
}

const handleImageLoad = () => {
  error.value = ''
  syncSurface()
  prepareCanvas()
}

const pointFromEvent = (event) => {
  syncSurface()
  const rect = surfaceRef.value.getBoundingClientRect()
  return {
    x: Math.min(Math.max(event.clientX - rect.left, 0), rect.width),
    y: Math.min(Math.max(event.clientY - rect.top, 0), rect.height),
  }
}

const pickColor = (event) => {
  if (!surfaceRef.value) return
  if (!canvasRef.value) prepareCanvas()
  if (!canvasRef.value) {
    error.value = '截图还没有加载完成'
    return
  }

  const displayPoint = pointFromEvent(event)
  const point = pointToNaturalPixel({
    point: displayPoint,
    rendered: surfaceSize.value,
    natural: { width: props.screenshot.width, height: props.screenshot.height },
  })

  const context = canvasRef.value.getContext('2d')
  const [r, g, b] = Array.from(context.getImageData(point[0], point[1], 1, 1).data).slice(0, 3)
  pickedDisplayPoint.value = displayPoint
  picked.value = {
    color: rgbToHex([r, g, b]),
    rgb: [r, g, b],
    point,
  }
}

const crosshairStyle = computed(() => ({
  left: `${pickedDisplayPoint.value.x}px`,
  top: `${pickedDisplayPoint.value.y}px`,
  background: picked.value?.color || '#ffffff',
}))

const applyColor = () => {
  if (!picked.value) return
  emit('apply', picked.value)
}

onMounted(() => {
  window.addEventListener('resize', syncSurface)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncSurface)
})
</script>
