<template>
  <div class="app-shell">
    <div class="hero-panel">
      <div class="brand-block">
        <div class="brand-mark">M</div>
        <div>
          <p class="eyebrow">MISSION CONTROL</p>
          <h1>maa-nte-go</h1>
        </div>
      </div>
      <div class="hero-status">
        <span class="status-badge" :class="status.connected ? 'connected' : 'disconnected'">
          {{ status.connected ? '已连接' : '未连接' }}
        </span>
        <span class="status-badge running" v-if="status.running">执行中</span>
        <span class="status-badge idle" v-else>待命</span>
      </div>
    </div>

    <div class="command-grid">
      <WindowPanel
        :windows="windows"
        :connected-handle="connectedHandle"
        @refresh="refreshWindows"
        @connect="connectWindow"
      />

      <PipelinePanel
        :pipelines="pipelines"
        :selected-pipeline="selectedPipeline"
        @refresh="refreshPipelines"
        @select="selectPipeline"
        @create="showCreateModal = true"
        @edit="editPipeline"
        @execute="executeTask"
        @delete="deletePipeline"
      />

      <ResultPanel v-if="lastResult" :result="lastResult" />
    </div>
  </div>

  <CreateModal
    v-if="showCreateModal"
    v-model:name="newName"
    @close="showCreateModal = false"
    @create="createPipeline"
  />

  <EditorModal
    v-if="showEditModal"
    :name="editingName"
    :nodes="editNodes"
    :images="images"
    :ocr-models="ocrModels"
    :detect-models="detectModels"
    @close="showEditModal = false"
    @save="saveVisualPipeline"
    @execute="executeNodeFromEditor"
    @preview="showJsonPreview = true"
  />

  <JsonPreviewModal
    v-if="showJsonPreview"
    :data="editNodes"
    @close="showJsonPreview = false"
  />

  <Toast v-if="toast.show" :message="toast.message" :type="toast.type" />
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { pipelineApi, taskApi, windowApi, resourceApi } from './api/index.js'
import WindowPanel from './components/WindowPanel.vue'
import PipelinePanel from './components/PipelinePanel.vue'
import ResultPanel from './components/ResultPanel.vue'
import CreateModal from './components/CreateModal.vue'
import EditorModal from './components/EditorModal.vue'
import JsonPreviewModal from './components/JsonPreviewModal.vue'
import Toast from './components/Toast.vue'

// State
const pipelines = ref([])
const windows = ref([])
const selectedPipeline = ref(null)
const connectedHandle = ref(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showJsonPreview = ref(false)
const newName = ref('')
const editingName = ref('')
const lastResult = ref(null)
const status = reactive({ connected: false, running: false })
const toast = reactive({ show: false, message: '', type: 'success' })

// Editor state
const editNodes = ref({})
const images = ref([])
const ocrModels = ref([])
const detectModels = ref([])

// Toast helper
const showToast = (message, type = 'success') => {
  toast.show = true
  toast.message = message
  toast.type = type
  setTimeout(() => toast.show = false, 3000)
}

// API calls
const refreshPipelines = async () => {
  try {
    const data = await pipelineApi.list()
    pipelines.value = data.pipelines || []
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const refreshWindows = async () => {
  try {
    const data = await windowApi.list()
    windows.value = data.windows || []
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const refreshStatus = async () => {
  try {
    const data = await taskApi.status()
    Object.assign(status, data)
  } catch (e) {}
}

const selectPipeline = (name) => {
  selectedPipeline.value = name
}

const connectWindow = async (w) => {
  try {
    await windowApi.connect(w.handle)
    connectedHandle.value = w.handle
    showToast('已连接到: ' + w.title)
    refreshStatus()
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const executeTask = async (name) => {
  try {
    await taskApi.execute(name)
    showToast('任务已启动: ' + name)
    refreshStatus()
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const editPipeline = async (name) => {
  try {
    const data = await pipelineApi.get(name)
    editingName.value = name
    editNodes.value = JSON.parse(JSON.stringify(data.content || {}))
    await Promise.all([
      fetchImages(),
      fetchOcrModels(),
      fetchDetectModels()
    ])
    showEditModal.value = true
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const fetchImages = async () => {
  try {
    const data = await resourceApi.images()
    images.value = data.images || []
  } catch (e) {
    images.value = []
  }
}

const fetchOcrModels = async () => {
  try {
    const data = await resourceApi.ocrModels()
    ocrModels.value = data.models || []
  } catch (e) {
    ocrModels.value = []
  }
}

const fetchDetectModels = async () => {
  try {
    const data = await resourceApi.detectModels()
    detectModels.value = data.models || []
  } catch (e) {
    detectModels.value = []
  }
}

const createPipeline = async () => {
  if (!newName.value) {
    showToast('请输入名称', 'error')
    return
  }
  try {
    await pipelineApi.create(newName.value, {})
    showToast('Pipeline 已创建')
    editingName.value = newName.value
    editNodes.value = {}
    await Promise.all([fetchImages(), fetchOcrModels(), fetchDetectModels()])
    showCreateModal.value = false
    showEditModal.value = true
    newName.value = ''
    refreshPipelines()
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const deletePipeline = async (name) => {
  if (!confirm('确定删除 ' + name + ' ?')) return
  try {
    await pipelineApi.delete(name)
    showToast('Pipeline 已删除')
    refreshPipelines()
  } catch (e) {
    showToast(e.message, 'error')
  }
}

const saveVisualPipeline = async (nodes) => {
  try {
    const cleanNodes = {}
    for (const [name, node] of Object.entries(nodes)) {
      const cleanNode = {}
      for (const [key, value] of Object.entries(node)) {
        if (value !== undefined && value !== '' && value !== null) {
          if (Array.isArray(value) && value.length === 0) continue
          cleanNode[key] = value
        }
      }
      cleanNodes[name] = cleanNode
    }
    await pipelineApi.update(editingName.value, cleanNodes)
    showToast('Pipeline 已保存')
    showEditModal.value = false
    refreshPipelines()
  } catch (e) {
    showToast('保存失败: ' + e.message, 'error')
  }
}

// 从编辑器执行单个节点（先保存再执行）
const executeNodeFromEditor = async ({ nodes, nodeName }) => {
  try {
    // 先保存
    const cleanNodes = {}
    for (const [name, node] of Object.entries(nodes)) {
      const cleanNode = {}
      for (const [key, value] of Object.entries(node)) {
        if (value !== undefined && value !== '' && value !== null) {
          if (Array.isArray(value) && value.length === 0) continue
          cleanNode[key] = value
        }
      }
      cleanNodes[name] = cleanNode
    }
    await pipelineApi.update(editingName.value, cleanNodes)

    // 再执行指定节点
    await taskApi.execute(editingName.value, nodeName)
    showToast('执行节点: ' + nodeName)
    refreshStatus()
  } catch (e) {
    showToast('执行失败: ' + e.message, 'error')
  }
}

// Lifecycle
onMounted(() => {
  refreshPipelines()
  refreshStatus()
})
</script>
