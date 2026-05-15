<template>
  <div class="color-match-picker">
    <div
      v-for="(row, idx) in localRows"
      :key="idx"
      class="color-row"
    >
      <label class="color-swatch">
        <input type="color" v-model="row.color" @input="syncRows">
        <span :style="{ background: row.color }"></span>
      </label>
      <div class="color-row-fields">
        <input v-model="row.color" class="field-input" @input="syncRows" placeholder="#ff3300">
        <label>
          容差
          <input type="number" v-model.number="row.tolerance" min="0" max="255" @input="syncRows">
        </label>
      </div>
      <div class="color-row-actions">
        <button
          class="btn btn-secondary btn-sm"
          type="button"
          :disabled="pickerLoading"
          @click="openWindowColorPicker(idx)"
        >
          窗口取色
        </button>
        <button class="btn btn-danger btn-sm" type="button" @click="removeRow(idx)">删除</button>
      </div>
    </div>

    <div class="color-picker-actions">
      <button class="btn btn-secondary btn-sm" type="button" @click="addRow">添加颜色</button>
      <button
        class="btn btn-primary btn-sm"
        type="button"
        :disabled="pickerLoading"
        @click="addAndPickFromWindow"
      >
        {{ pickerLoading ? '截图中...' : '从窗口新增' }}
      </button>
    </div>
    <div v-if="pickerError" class="field-error">{{ pickerError }}</div>
    <div class="color-bounds-preview">
      <span>lower</span>
      <code>{{ JSON.stringify(bounds.lower || []) }}</code>
      <span>upper</span>
      <code>{{ JSON.stringify(bounds.upper || []) }}</code>
    </div>
    <div class="field-hint">RGB 模式下使用取色器和容差自动生成 Maa 的 lower/upper。</div>
  </div>
  <WindowColorPickerModal
    v-if="showWindowColorPicker && pickerScreenshot"
    :screenshot="pickerScreenshot"
    @close="showWindowColorPicker = false"
    @apply="applyPickedWindowColor"
  />
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { resourceApi } from '../api/index.js'
import WindowColorPickerModal from './WindowColorPickerModal.vue'
import {
  colorRowsFromBounds,
  serializeColorBounds,
} from '../utils/colorMatch.js'

const props = defineProps({
  lower: { type: [Array, Number, String], default: undefined },
  upper: { type: [Array, Number, String], default: undefined },
})

const emit = defineEmits(['update:lower', 'update:upper'])

const localRows = ref([])
const showWindowColorPicker = ref(false)
const pickerScreenshot = ref(null)
const pickerLoading = ref(false)
const pickerError = ref('')
const activePickerRow = ref(-1)

const bounds = computed(() => serializeColorBounds(localRows.value))

const resetRows = () => {
  const rows = colorRowsFromBounds(props.lower, props.upper)
  localRows.value = rows.length > 0 ? rows : [{ color: '#ff0000', tolerance: 10 }]
}

watch(() => [props.lower, props.upper], resetRows, { immediate: true, deep: true })

const syncRows = () => {
  const nextBounds = serializeColorBounds(localRows.value)
  emit('update:lower', nextBounds.lower)
  emit('update:upper', nextBounds.upper)
}

const addRow = () => {
  localRows.value = [...localRows.value, { color: '#ffffff', tolerance: 10 }]
  syncRows()
}

const addAndPickFromWindow = () => {
  const nextIndex = localRows.value.length
  localRows.value = [...localRows.value, { color: '#ffffff', tolerance: 10 }]
  syncRows()
  openWindowColorPicker(nextIndex)
}

const removeRow = (idx) => {
  localRows.value = localRows.value.filter((_, rowIdx) => rowIdx !== idx)
  syncRows()
}

const openWindowColorPicker = async (idx) => {
  if (pickerLoading.value) return
  pickerError.value = ''
  pickerLoading.value = true
  activePickerRow.value = idx
  try {
    pickerScreenshot.value = await resourceApi.screenshot()
    showWindowColorPicker.value = true
  } catch (e) {
    pickerError.value = e.message || '截图失败，请先连接窗口'
  } finally {
    pickerLoading.value = false
  }
}

const applyPickedWindowColor = ({ color }) => {
  if (activePickerRow.value < 0 || !localRows.value[activePickerRow.value]) return
  localRows.value[activePickerRow.value].color = color
  syncRows()
  showWindowColorPicker.value = false
}
</script>
