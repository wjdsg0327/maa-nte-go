<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal large">
      <div class="modal-titlebar">
        <div>
          <p class="eyebrow">PIPELINE EDITOR</p>
          <h2>编辑 Pipeline - {{ name }}</h2>
        </div>
        <button class="btn btn-secondary btn-sm" @click="$emit('close')">关闭</button>
      </div>
      <div class="help-text">
        在流程画布里选择节点，右侧参数面板会同步编辑同一份 Pipeline。
      </div>
      <div class="editor-workbench">
        <PipelineFlowBoard
          :nodes="localNodes"
          :active-node="activeNodeName"
          @select="selectNode"
          @add="addNode"
          @execute="executeNode"
        />
        <div class="node-inspector">
          <div class="inspector-titlebar">
            <div>
              <p class="eyebrow">NODE INSPECTOR</p>
              <h3>{{ activeNodeName || '未选择节点' }}</h3>
            </div>
          </div>
          <NodeEditor
            v-if="activeNode"
            :node="activeNode"
            :name="activeNodeName"
            :images="images"
            :ocr-models="ocrModels"
            :detect-models="detectModels"
            :all-nodes="Object.keys(localNodes)"
            @update:name="renameNode"
            @delete="deleteActiveNode"
            @pick-roi="openRoiPicker"
          />
          <div v-else class="editor-empty">
            <div class="empty-icon">+</div>
            <strong>等待部署节点</strong>
            <p>从流程画布选择节点，或新建第一个节点开始编排。</p>
            <button class="btn btn-primary" @click="addNode">新建节点</button>
          </div>
        </div>
      </div>
      <div class="modal-actions">
        <span v-if="roiPickerError" class="inline-error">{{ roiPickerError }}</span>
        <button class="btn btn-secondary" @click="$emit('preview')">JSON预览</button>
        <button class="btn btn-success" v-if="activeNodeName" @click="executeNode(activeNodeName)">执行当前节点</button>
        <button class="btn btn-secondary" v-if="roiPickerLoading" disabled>截图中...</button>
        <button class="btn btn-primary" @click="save">保存</button>
      </div>
    </div>
  </div>
  <RoiPickerModal
    v-if="showRoiPicker && roiScreenshot"
    :screenshot="roiScreenshot"
    :initial-roi="Array.isArray(activeNode?.roi) ? activeNode.roi : []"
    @close="showRoiPicker = false"
    @apply="applyPickedRoi"
  />
</template>

<script setup>
import { ref, computed, watch, reactive } from 'vue'
import PipelineFlowBoard from './PipelineFlowBoard.vue'
import NodeEditor from './NodeEditor.vue'
import RoiPickerModal from './RoiPickerModal.vue'
import { resourceApi } from '../api/index.js'
import {
  removeRelationTargetName,
  renameRelationTarget,
} from '../utils/pipelineRelations.js'

const props = defineProps({
  name: { type: String, required: true },
  nodes: { type: Object, required: true },
  images: { type: Array, default: () => [] },
  ocrModels: { type: Array, default: () => [] },
  detectModels: { type: Array, default: () => [] }
})

const emit = defineEmits(['close', 'save', 'preview', 'execute'])
const RELATION_FIELDS = ['next', 'on_error', 'interrupt']

// Local copy of nodes for editing
const localNodes = ref(JSON.parse(JSON.stringify(props.nodes)))
const activeNodeName = ref('')
const showRoiPicker = ref(false)
const roiScreenshot = ref(null)
const roiPickerLoading = ref(false)
const roiPickerError = ref('')

// Watch for external nodes changes
watch(() => props.nodes, (newNodes) => {
  localNodes.value = JSON.parse(JSON.stringify(newNodes))
  const names = Object.keys(localNodes.value)
  if (names.length > 0 && !activeNodeName.value) {
    activeNodeName.value = names[0]
  }
}, { immediate: true })

const activeNode = computed(() => {
  if (!activeNodeName.value || !localNodes.value[activeNodeName.value]) return null
  return localNodes.value[activeNodeName.value]
})

const selectNode = (name) => {
  activeNodeName.value = name
}

const addNode = () => {
  const baseName = '新节点'
  let counter = 1
  let newName = baseName
  while (localNodes.value[newName]) {
    newName = `${baseName}_${counter}`
    counter++
  }
  localNodes.value[newName] = {
    recognition: 'DirectHit',
    action: 'DoNothing'
  }
  activeNodeName.value = newName
}

const renameNode = (newName) => {
  if (!activeNodeName.value || !newName || newName === activeNodeName.value) return

  const oldName = activeNodeName.value
  const node = localNodes.value[oldName]

  // Create new nodes object with renamed key
  const newNodes = {}
  for (const [key, value] of Object.entries(localNodes.value)) {
    if (key === oldName) {
      newNodes[newName] = value
    } else {
      newNodes[key] = value
    }
  }

  // Update all references
  for (const n of Object.values(newNodes)) {
    for (const field of RELATION_FIELDS) {
      if (n[field] === undefined) continue
      const updated = renameRelationTarget(n[field], oldName, newName)
      if (updated === undefined) {
        delete n[field]
      } else {
        n[field] = updated
      }
    }
  }

  localNodes.value = newNodes
  activeNodeName.value = newName
}

const deleteActiveNode = () => {
  if (!activeNodeName.value) return
  const nodeName = activeNodeName.value
  delete localNodes.value[nodeName]

  // Remove all references
  for (const node of Object.values(localNodes.value)) {
    for (const field of RELATION_FIELDS) {
      if (node[field] === undefined) continue
      const updated = removeRelationTargetName(node[field], nodeName)
      if (updated === undefined) {
        delete node[field]
      } else {
        node[field] = updated
      }
    }
  }

  // Select another node
  const names = Object.keys(localNodes.value)
  activeNodeName.value = names.length > 0 ? names[0] : ''
}

const openRoiPicker = async () => {
  if (!activeNode.value || roiPickerLoading.value) return
  roiPickerError.value = ''
  roiPickerLoading.value = true
  try {
    roiScreenshot.value = await resourceApi.screenshot()
    showRoiPicker.value = true
  } catch (e) {
    roiPickerError.value = e.message
  } finally {
    roiPickerLoading.value = false
  }
}

const applyPickedRoi = (roi) => {
  if (activeNode.value) {
    activeNode.value.roi = roi
  }
  showRoiPicker.value = false
}

const executeNode = (nodeName) => {
  // Emit execute with nodes data and node name
  emit('execute', { nodes: localNodes.value, nodeName })
}

const save = () => {
  emit('save', localNodes.value)
}
</script>
