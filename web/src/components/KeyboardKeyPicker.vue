<template>
  <div class="key-picker">
    <div v-if="selectedCodes.length > 0" class="next-tags key-picker-tags">
      <span v-for="(code, idx) in selectedCodes" :key="code" class="next-tag">
        {{ formatKeyCode(code) }}
        <span class="remove" @click="removeKey(idx)">×</span>
      </span>
    </div>

    <div class="key-picker-search">
      <input
        v-model="searchText"
        :placeholder="placeholder"
        @keydown.enter.prevent="addFromSearch"
      >
      <button class="btn btn-primary btn-sm" @click="addFromSearch" :disabled="!searchText.trim()">添加</button>
    </div>

    <div class="key-picker-options">
      <button
        v-for="option in filteredOptions"
        :key="option.code"
        type="button"
        class="key-option"
        :class="{ selected: selectedCodes.includes(option.code) }"
        @click="selectKey(option.code)"
      >
        <span>{{ option.label }}</span>
        <small>{{ option.code }}</small>
      </button>
    </div>

    <div v-if="errorText" class="field-error">{{ errorText }}</div>
    <div class="field-hint">可搜索 Esc、Enter、F5、A、空格，也可以直接输入数字键码。</div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import {
  filterKeyboardKeys,
  formatKeyCode,
  normalizeKeyCodes,
  parseKeyboardKeyInput,
  serializeKeyCodes,
} from '../utils/keyboardKeys.js'

const props = defineProps({
  modelValue: { type: [Number, String, Array], default: undefined },
  multiple: { type: Boolean, default: false },
  placeholder: { type: String, default: '搜索按键，如 Esc / Enter / F5 / A' },
})

const emit = defineEmits(['update:modelValue'])

const searchText = ref('')
const errorText = ref('')

const selectedCodes = computed(() => normalizeKeyCodes(props.modelValue))
const filteredOptions = computed(() => filterKeyboardKeys(searchText.value))

const emitCodes = (codes) => {
  emit('update:modelValue', serializeKeyCodes(codes, props.multiple))
}

const selectKey = (code) => {
  errorText.value = ''
  if (!props.multiple) {
    emitCodes([code])
    searchText.value = ''
    return
  }

  emitCodes([...selectedCodes.value, code])
  searchText.value = ''
}

const addFromSearch = () => {
  const code = parseKeyboardKeyInput(searchText.value)
  if (code === undefined) {
    errorText.value = '没有找到这个按键，可以换个名字试试，比如 Esc、Enter、Space。'
    return
  }
  selectKey(code)
}

const removeKey = (idx) => {
  const nextCodes = [...selectedCodes.value]
  nextCodes.splice(idx, 1)
  emitCodes(nextCodes)
}
</script>
