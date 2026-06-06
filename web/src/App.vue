<template>
  <main class="app-shell">
    <nav class="top-nav">
      <div class="brand-block">
        <div class="brand-mark">M</div>
        <div>
          <strong>Maa Console</strong>
          <span>{{ appConfig.resourceDir || '通用自动化资源包' }}</span>
        </div>
      </div>
      <div class="nav-actions">
        <span class="status-badge" :class="status.connected ? 'connected' : 'disconnected'">
          {{ status.connected ? '已连接' : '未连接' }}
        </span>
        <span class="status-badge running" v-if="status.running">执行中</span>
        <span class="status-badge idle" v-else>待命</span>
        <button class="btn btn-secondary" @click="refreshAll">同步状态</button>
      </div>
    </nav>

    <section class="hero-panel motion-card">
      <div class="hero-copy">
        <p class="eyebrow">MISSION CONTROL</p>
        <h1>把 Maa 资源包变成任何软件的自动化工作台</h1>
        <p class="hero-lede">
          当前资源、目标窗口、Pipeline 与调试结果集中在同一个控制台里。
        </p>
        <div class="hero-actions">
          <button class="btn btn-primary" @click="showCreateModal = true">新建 Pipeline</button>
          <button class="btn btn-secondary" @click="reloadResources" :disabled="reloadingResources">
            {{ reloadingResources ? '资源重载中' : '重载资源包' }}
          </button>
        </div>
      </div>
      <div class="hero-visual motion-image">
        <div class="signal-grid">
          <span v-for="item in heroNodes" :key="item" />
        </div>
        <div class="hero-terminal">
          <span>resource</span>
          <strong>{{ appConfig.resourceDir || './resource' }}</strong>
          <span>pipelines</span>
          <strong>{{ pipelines.length }}</strong>
          <span>windows</span>
          <strong>{{ windows.length }}</strong>
        </div>
      </div>
    </section>

    <section class="command-grid">
      <WindowPanel
        class="motion-card"
        :windows="windows"
        :connected-handle="connectedHandle"
        @refresh="refreshWindows"
        @connect="connectWindow"
      />

      <PipelinePanel
        class="motion-card"
        :pipelines="pipelines"
        :selected-pipeline="selectedPipeline"
        @refresh="refreshPipelines"
        @select="selectPipeline"
        @create="showCreateModal = true"
        @edit="editPipeline"
        @execute="executeTask"
        @delete="deletePipeline"
      />

      <ResourcePanel
        class="motion-card"
        :config="appConfig"
        :pipeline-count="pipelines.length"
        :image-count="images.length"
        :ocr-count="ocrModels.length"
        :detect-count="detectModels.length"
        :reloading="reloadingResources"
        @reload="reloadResources"
      />

      <ResultPanel v-if="lastResult" class="motion-card" :result="lastResult" />

      <section v-else class="card action-card motion-card">
        <div>
          <p class="eyebrow">OPEN PROJECT</p>
          <h2>下一步：为新软件建立资源包</h2>
        </div>
        <p>
          新资源包准备就绪后，会在这里显示窗口连接、任务执行和识别回放。
        </p>
        <div class="action-strip">
          <code>pipeline/</code>
          <code>image/</code>
          <code>ocr/</code>
          <code>detect/</code>
        </div>
      </section>
    </section>
  </main>

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
import { ref, reactive, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { pipelineApi, taskApi, windowApi, resourceApi, appApi } from './api/index.js'
import WindowPanel from './components/WindowPanel.vue'
import PipelinePanel from './components/PipelinePanel.vue'
import ResultPanel from './components/ResultPanel.vue'
import ResourcePanel from './components/ResourcePanel.vue'
import CreateModal from './components/CreateModal.vue'
import EditorModal from './components/EditorModal.vue'
import JsonPreviewModal from './components/JsonPreviewModal.vue'
import Toast from './components/Toast.vue'
import { cleanPipelineNodes } from './utils/pipelineClean.js'

gsap.registerPlugin(ScrollTrigger)

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
const appConfig = reactive({})
const toast = reactive({ show: false, message: '', type: 'success' })
const reloadingResources = ref(false)
const heroNodes = Array.from({ length: 24 }, (_, idx) => idx)
let statusTimer = null

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

const refreshConfig = async () => {
  try {
    const data = await appApi.config()
    Object.assign(appConfig, data)
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

const refreshAssets = async () => {
  await Promise.all([fetchImages(), fetchOcrModels(), fetchDetectModels()])
}

const refreshAll = async () => {
  await Promise.all([
    refreshConfig(),
    refreshPipelines(),
    refreshWindows(),
    refreshStatus(),
    refreshAssets()
  ])
  await nextTick()
  animateWorkbench()
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
    const data = await taskApi.run(name)
    lastResult.value = data
    showToast('任务执行完成: ' + name)
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

const reloadResources = async () => {
  reloadingResources.value = true
  try {
    await resourceApi.reload()
    await Promise.all([refreshPipelines(), refreshAssets()])
    showToast('资源包已热重载')
  } catch (e) {
    showToast('重载失败: ' + e.message, 'error')
  } finally {
    reloadingResources.value = false
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
    const cleanNodes = cleanPipelineNodes(nodes)
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
    const cleanNodes = cleanPipelineNodes(nodes)
    await pipelineApi.update(editingName.value, cleanNodes)

    // 再执行指定节点
    const data = await taskApi.run(editingName.value, nodeName)
    lastResult.value = data
    showToast('节点执行完成: ' + nodeName)
    refreshStatus()
  } catch (e) {
    showToast('执行失败: ' + e.message, 'error')
  }
}

const animateWorkbench = () => {
  gsap.fromTo(
    '.motion-card',
    { opacity: 0, y: 22, scale: 0.985 },
    { opacity: 1, y: 0, scale: 1, duration: 0.7, ease: 'power3.out', stagger: 0.06, overwrite: true }
  )
  gsap.fromTo(
    '.motion-image',
    { opacity: 0.45, scale: 0.92 },
    {
      opacity: 1,
      scale: 1,
      duration: 1,
      ease: 'power3.out',
      scrollTrigger: {
        trigger: '.hero-panel',
        start: 'top 80%',
        end: 'bottom 20%',
        scrub: true
      }
    }
  )
}

// Lifecycle
onMounted(async () => {
  await refreshAll()
  statusTimer = setInterval(refreshStatus, 2500)
})

onBeforeUnmount(() => {
  if (statusTimer) clearInterval(statusTimer)
  ScrollTrigger.getAll().forEach(trigger => trigger.kill())
})
</script>
