<template>
  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.enabled" id="enabled">
      <label for="enabled" style="margin-bottom: 0;">启用节点</label>
    </div>
  </div>

  <div class="field-group">
    <label>识别超时(ms)</label>
    <input type="number" v-model.number="node.timeout" min="-1" placeholder="20000">
    <div class="field-hint">-1 表示无限等待；留空使用默认值。</div>
  </div>

  <div class="field-group">
    <label>识别速率限制(ms)</label>
    <input type="number" v-model.number="node.rate_limit" min="0" placeholder="1000">
    <div class="field-hint">每轮识别至少消耗的时间。</div>
  </div>

  <div class="field-group">
    <label>最大命中次数</label>
    <input type="number" v-model.number="node.max_hit" min="1" placeholder="留空表示不限制">
  </div>

  <div class="field-group">
    <label>动作前静止等待(ms/JSON)</label>
    <input v-model="preWaitFreezesInput" :class="{ invalid: preWaitFreezesError }" placeholder="0 或 {&quot;time&quot;:500}">
    <div v-if="preWaitFreezesError" class="field-error">{{ preWaitFreezesError }}</div>
    <div class="field-hint">识别命中后，先等待画面静止，再执行前置延时和动作。</div>
  </div>

  <div class="field-group">
    <label>动作前延时(ms)</label>
    <input type="number" v-model.number="node.pre_delay" min="0" placeholder="200">
  </div>

  <div class="field-group">
    <label>动作后静止等待(ms/JSON)</label>
    <input v-model="postWaitFreezesInput" :class="{ invalid: postWaitFreezesError }" placeholder="0 或 {&quot;time&quot;:500}">
    <div v-if="postWaitFreezesError" class="field-error">{{ postWaitFreezesError }}</div>
  </div>

  <div class="field-group">
    <label>动作后延时(ms)</label>
    <input type="number" v-model.number="node.post_delay" min="0" placeholder="200">
  </div>

  <div class="field-group">
    <label>动作重复次数</label>
    <input type="number" v-model.number="node.repeat" min="1" placeholder="1">
    <div class="field-hint">同一个 Action 重复执行；多个不同动作建议拆节点用 next 串联。</div>
  </div>

  <div class="field-group">
    <label>重复间隔(ms)</label>
    <input type="number" v-model.number="node.repeat_delay" min="0" placeholder="0">
  </div>

  <div class="field-group">
    <label>重复前静止等待(ms/JSON)</label>
    <input v-model="repeatWaitFreezesInput" :class="{ invalid: repeatWaitFreezesError }" placeholder="0 或 {&quot;time&quot;:500}">
    <div v-if="repeatWaitFreezesError" class="field-error">{{ repeatWaitFreezesError }}</div>
    <div class="field-hint">仅在动作重复次数大于 1 时生效。</div>
  </div>

  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.inverse" id="inverse">
      <label for="inverse" style="margin-bottom: 0;">反转识别结果</label>
    </div>
    <div class="field-hint">识别成功当失败，失败当成功。</div>
  </div>

  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.save_result" id="saveResult">
      <label for="saveResult" style="margin-bottom: 0;">保存识别结果</label>
    </div>
  </div>

  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.save_draw" id="saveDraw">
      <label for="saveDraw" style="margin-bottom: 0;">保存标记图片</label>
    </div>
    <div class="field-hint">保存带识别框标记的截图到 log/vision 目录。</div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import {
  formatNumberOrObject,
  parseNumberOrObjectValue,
} from '../../utils/pipelineFields.js'

const props = defineProps({
  node: { type: Object, required: true },
})

const preWaitFreezesInput = ref('')
const postWaitFreezesInput = ref('')
const repeatWaitFreezesInput = ref('')
const preWaitFreezesError = ref('')
const postWaitFreezesError = ref('')
const repeatWaitFreezesError = ref('')

const migrateLegacyTimes = (node) => {
  if (!node) return
  if (node.repeat === undefined && node.times !== undefined) {
    node.repeat = node.times
  }
  if (node.times !== undefined) {
    delete node.times
  }
}

watch(() => props.node, (node) => {
  migrateLegacyTimes(node)
  preWaitFreezesInput.value = formatNumberOrObject(node.pre_wait_freezes)
  postWaitFreezesInput.value = formatNumberOrObject(node.post_wait_freezes)
  repeatWaitFreezesInput.value = formatNumberOrObject(node.repeat_wait_freezes)
  preWaitFreezesError.value = ''
  postWaitFreezesError.value = ''
  repeatWaitFreezesError.value = ''
}, { immediate: true })

const syncNumberOrObject = (field, value, errorRef, label) => {
  const result = parseNumberOrObjectValue(value, label)
  errorRef.value = result.error
  if (!result.error) {
    props.node[field] = result.value
  }
}

watch(preWaitFreezesInput, (val) => {
  syncNumberOrObject('pre_wait_freezes', val, preWaitFreezesError, '动作前静止等待')
})

watch(postWaitFreezesInput, (val) => {
  syncNumberOrObject('post_wait_freezes', val, postWaitFreezesError, '动作后静止等待')
})

watch(repeatWaitFreezesInput, (val) => {
  syncNumberOrObject('repeat_wait_freezes', val, repeatWaitFreezesError, '重复前静止等待')
})
</script>
