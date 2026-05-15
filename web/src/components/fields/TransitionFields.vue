<template>
  <div
    v-for="config in relationConfigs"
    :key="config.field"
    class="field-group"
  >
    <label>{{ config.label }}</label>
    <div class="next-tags">
      <span
        v-for="(item, idx) in relationItems(config.field)"
        :key="`${config.field}-${idx}-${relationLabel(item)}`"
        class="next-tag"
      >
        {{ relationLabel(item) }}
        <span class="remove" @click="removeRelation(config.field, idx)">×</span>
      </span>
    </div>
    <div class="add-next-input">
      <select v-model="newRelationNames[config.field]" class="field-input" style="flex: 1;">
        <option value="">选择节点...</option>
        <option v-for="n in availableNodes" :key="n" :value="n">{{ n }}</option>
      </select>
      <button
        class="btn btn-primary btn-sm"
        @click="addRelation(config.field)"
        :disabled="!newRelationNames[config.field]"
      >
        添加
      </button>
    </div>
    <div class="field-hint">{{ config.hint }}</div>
  </div>
</template>

<script setup>
import { computed, reactive } from 'vue'
import {
  addRelationTarget,
  getRelationDisplayName,
  getRelationTargets,
  removeRelationTargetAt,
} from '../../utils/pipelineRelations.js'

const props = defineProps({
  node: { type: Object, required: true },
  allNodes: { type: Array, default: () => [] }
})

const relationConfigs = [
  {
    field: 'next',
    label: '下一节点',
    hint: '节点动作完成后，按顺序识别并进入的候选节点。',
  },
  {
    field: 'on_error',
    label: '失败节点',
    hint: '当前节点动作失败，或 next 候选超时未命中时进入的节点。',
  },
  {
    field: 'interrupt',
    label: '中断节点（旧版）',
    hint: 'Maa 5.1 已废弃 interrupt；旧项目可保留，新流程推荐使用 next 的 JumpBack 节点属性。',
  },
]

const newRelationNames = reactive({
  next: '',
  on_error: '',
  interrupt: '',
})

const availableNodes = computed(() => props.allNodes)

const setRelationField = (field, value) => {
  if (value === undefined) {
    delete props.node[field]
  } else {
    props.node[field] = value
  }
}

const relationItems = (field) => getRelationTargets(props.node[field])

const relationLabel = (item) => {
  const name = getRelationDisplayName(item)
  const attrs = []

  if (typeof item === 'string') {
    if (item.startsWith('[JumpBack]')) attrs.push('JumpBack')
    if (item.startsWith('[Anchor]')) attrs.push('Anchor')
  } else {
    if (item?.jump_back) attrs.push('JumpBack')
    if (item?.anchor) attrs.push('Anchor')
  }

  return attrs.length > 0 ? `${name} · ${attrs.join('/')}` : name
}

const addRelation = (field) => {
  const nextValue = addRelationTarget(props.node[field], newRelationNames[field])
  setRelationField(field, nextValue)
  newRelationNames[field] = ''
}

const removeRelation = (field, idx) => {
  setRelationField(field, removeRelationTargetAt(props.node[field], idx))
}
</script>
