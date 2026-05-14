<template>
  <div class="nodes-panel">
    <div class="nodes-toolbar">
      <button class="btn btn-primary btn-sm" @click="$emit('add')">+ 新建节点</button>
      <input v-model="searchText" class="search-input" placeholder="搜索...">
      <button class="btn btn-sm btn-icon" :class="viewMode === 'tree' ? 'btn-primary' : 'btn-secondary'" @click="viewMode = 'tree'" title="列表视图" aria-label="列表视图">☰</button>
      <button class="btn btn-sm btn-icon" :class="viewMode === 'flow' ? 'btn-primary' : 'btn-secondary'" @click="viewMode = 'flow'" title="流程图" aria-label="流程图">⊞</button>
    </div>

    <!-- 流程图模式 -->
    <div v-if="viewMode === 'flow'" class="flow-container">
      <svg class="flow-svg" :width="svgWidth" :height="svgHeight">
        <!-- 连线 -->
        <defs>
          <marker id="arrow" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#2c3145"/>
          </marker>
          <marker id="arrow-interrupt" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#ff5f74"/>
          </marker>
          <marker id="arrow-wait" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#ffd84d"/>
          </marker>
        </defs>

        <!-- 绘制连线 -->
        <g v-for="(node, name) in filteredNodes" :key="'line-'+name">
          <!-- Next 连线 (绿色实线) -->
          <path
            v-for="target in (node.next || [])"
            :key="name+'-next-'+target"
            :d="getPath(name, target, 'next')"
            class="flow-line flow-line-next"
            marker-end="url(#arrow)"
          />
          <!-- Interrupt 连线 (红色虚线) -->
          <path
            v-for="target in (node.interrupt || [])"
            :key="name+'-int-'+target"
            :d="getPath(name, target, 'interrupt')"
            class="flow-line flow-line-interrupt"
            marker-end="url(#arrow-interrupt)"
          />
          <!-- Wait 连线 (橙色点线) -->
          <path
            v-for="target in (node.wait_freezes || [])"
            :key="name+'-wait-'+target"
            :d="getPath(name, target, 'wait')"
            class="flow-line flow-line-wait"
            marker-end="url(#arrow-wait)"
          />
        </g>

        <!-- 绘制节点 -->
        <g
          v-for="(pos, name) in nodePositions"
          :key="'node-'+name"
          :transform="`translate(${pos.x}, ${pos.y})`"
          class="flow-node"
          :class="{ active: activeNode === name, root: isRoot(name) }"
          @click="$emit('select', name)"
        >
          <rect
            class="node-rect"
            :width="nodeWidth"
            :height="nodeHeight"
            rx="6"
            ry="6"
          />
          <text class="node-label" :x="nodeWidth/2" y="20" text-anchor="middle">
            {{ truncateName(name) }}
          </text>
          <text class="node-type" :x="nodeWidth/2" y="38" text-anchor="middle">
            {{ filteredNodes[name]?.recognition?.[0] || 'D' }} → {{ filteredNodes[name]?.action?.[0] || 'N' }}
          </text>
          <circle
            class="run-btn-svg"
            :cx="nodeWidth - 12"
            cy="12"
            r="8"
            @click.stop="$emit('execute', name)"
          />
          <text :x="nodeWidth - 12" y="16" text-anchor="middle" class="run-icon" @click.stop="$emit('execute', name)">▶</text>
        </g>
      </svg>

      <!-- 图例 -->
      <div class="flow-legend">
        <span class="legend-item"><span class="legend-line next"></span> Next</span>
        <span class="legend-item"><span class="legend-line interrupt"></span> Interrupt</span>
        <span class="legend-item"><span class="legend-line wait"></span> Wait</span>
      </div>
      <div v-if="Object.keys(filteredNodes).length === 0" class="flow-empty empty-state">
        <div class="empty-icon">+</div>
        <p>暂无节点，点击“新建节点”开始编排。</p>
      </div>
    </div>

    <!-- 列表模式 -->
    <div v-else class="nodes-tree">
      <div v-for="group in filteredGroups" :key="group.root" class="node-group">
        <div
          class="node-item root-node"
          :class="{ active: activeNode === group.root, expanded: expandedGroups[group.root] }"
          @click="selectNode(group.root)"
        >
          <span class="expand-icon" @click.stop="toggleGroup(group.root)">
            {{ group.children.length > 0 ? (expandedGroups[group.root] ? '▼' : '▶') : '　' }}
          </span>
          <span class="node-name">{{ group.root }}</span>
          <span class="node-badge">{{ filteredNodes[group.root]?.recognition?.[0] || 'D' }}</span>
          <button class="run-btn" @click.stop="$emit('execute', group.root)" title="执行">▶</button>
        </div>
        <template v-if="expandedGroups[group.root] && group.children.length > 0">
          <div
            v-for="child in group.children"
            :key="child"
            class="node-item child-node"
            :class="{ active: activeNode === child }"
            @click="$emit('select', child)"
          >
            <span class="child-line"></span>
            <span class="node-name">{{ child }}</span>
            <span class="node-badge">{{ filteredNodes[child]?.recognition?.[0] || 'D' }}</span>
            <button class="run-btn" @click.stop="$emit('execute', child)" title="执行">▶</button>
          </div>
        </template>
      </div>

      <div v-if="standaloneNodes.length > 0" class="standalone-section">
        <div class="section-title" @click="showStandalone = !showStandalone">
          {{ showStandalone ? '▼' : '▶' }} 独立节点 ({{ standaloneNodes.length }})
        </div>
        <template v-if="showStandalone">
          <div
            v-for="name in standaloneNodes"
            :key="name"
            class="node-item standalone-node"
            :class="{ active: activeNode === name }"
            @click="$emit('select', name)"
          >
            <span class="node-name">{{ name }}</span>
            <span class="node-badge">{{ filteredNodes[name]?.recognition?.[0] || 'D' }}</span>
            <button class="run-btn" @click.stop="$emit('execute', name)" title="执行">▶</button>
          </div>
        </template>
      </div>
      <div v-if="Object.keys(filteredNodes).length === 0" class="nodes-empty">
        暂无匹配节点
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, reactive, watch } from 'vue'

const props = defineProps({
  nodes: { type: Object, required: true },
  activeNode: { type: String, default: '' }
})

const emit = defineEmits(['select', 'add', 'execute'])

const viewMode = ref('flow')
const searchText = ref('')
const showStandalone = ref(true)
const expandedGroups = reactive({})

// 节点尺寸
const nodeWidth = 100
const nodeHeight = 48
const nodeGapX = 40
const nodeGapY = 60

// 过滤后的节点
const filteredNodes = computed(() => {
  if (!searchText.value) return props.nodes
  const search = searchText.value.toLowerCase()
  const result = {}
  for (const [name, node] of Object.entries(props.nodes)) {
    if (name.toLowerCase().includes(search)) {
      result[name] = node
    }
  }
  return result
})

// 检查是否是根节点
const isRoot = (name) => {
  for (const node of Object.values(props.nodes)) {
    if (node.next && node.next.includes(name)) return false
  }
  return true
}

// 获取根节点列表
const rootNodes = computed(() => {
  const children = new Set()
  for (const node of Object.values(props.nodes)) {
    if (node.next) node.next.forEach(n => children.add(n))
  }
  return Object.keys(props.nodes).filter(n => !children.has(n))
})

// 计算节点位置（树状布局）
const nodePositions = computed(() => {
  const positions = {}
  const visited = new Set()
  let currentX = 20
  let currentY = 20

  const processNode = (name, x, y) => {
    if (visited.has(name) || !filteredNodes.value[name]) return x
    visited.add(name)
    positions[name] = { x, y }

    const node = filteredNodes.value[name]
    const children = node.next || []

    if (children.length === 0) {
      return x + nodeWidth + nodeGapX
    }

    let childX = x
    const childY = y + nodeHeight + nodeGapY

    for (const child of children) {
      if (!visited.has(child) && filteredNodes.value[child]) {
        childX = processNode(child, childX, childY)
      }
    }

    return Math.max(childX, x + nodeWidth + nodeGapX)
  }

  // 处理根节点
  for (const root of rootNodes.value) {
    if (filteredNodes.value[root]) {
      currentX = processNode(root, currentX, currentY)
    }
  }

  // 处理未访问的节点
  for (const name of Object.keys(filteredNodes.value)) {
    if (!visited.has(name)) {
      positions[name] = { x: currentX, y: currentY }
      currentX += nodeWidth + nodeGapX
    }
  }

  return positions
})

// SVG 尺寸
const svgWidth = computed(() => {
  const positions = Object.values(nodePositions.value)
  if (positions.length === 0) return 400
  return Math.max(400, Math.max(...positions.map(p => p.x)) + nodeWidth + 40)
})

const svgHeight = computed(() => {
  const positions = Object.values(nodePositions.value)
  if (positions.length === 0) return 300
  return Math.max(300, Math.max(...positions.map(p => p.y)) + nodeHeight + 40)
})

// 获取连线路径
const getPath = (from, to, type) => {
  const fromPos = nodePositions.value[from]
  const toPos = nodePositions.value[to]
  if (!fromPos || !toPos) return ''

  const x1 = fromPos.x + nodeWidth / 2
  const y1 = fromPos.y + nodeHeight
  const x2 = toPos.x + nodeWidth / 2
  const y2 = toPos.y

  // 贝塞尔曲线
  const midY = (y1 + y2) / 2
  return `M${x1},${y1} C${x1},${midY} ${x2},${midY} ${x2},${y2}`
}

// 截断名称
const truncateName = (name) => {
  return name.length > 8 ? name.substring(0, 7) + '...' : name
}

// 列表模式的分组
const nodeGroups = computed(() => {
  const children = new Set()
  for (const node of Object.values(props.nodes)) {
    if (node.next) node.next.forEach(n => children.add(n))
  }
  const roots = Object.keys(props.nodes).filter(n => !children.has(n))
  return roots.map(root => ({
    root,
    children: (props.nodes[root]?.next || []).filter(n => props.nodes[n])
  }))
})

const standaloneNodes = computed(() => {
  const rootSet = new Set(nodeGroups.value.map(g => g.root))
  const childSet = new Set()
  nodeGroups.value.forEach(g => g.children.forEach(c => childSet.add(c)))
  return Object.keys(props.nodes).filter(n => !rootSet.has(n) && !childSet.has(n))
})

const filteredGroups = computed(() => {
  if (!searchText.value) return nodeGroups.value
  const search = searchText.value.toLowerCase()
  return nodeGroups.value.filter(g =>
    g.root.toLowerCase().includes(search) || g.children.some(c => c.toLowerCase().includes(search))
  )
})

const toggleGroup = (rootName) => {
  expandedGroups[rootName] = !expandedGroups[rootName]
}

const selectNode = (name) => {
  emit('select', name)
  if (props.nodes[name]?.next?.length > 0) {
    toggleGroup(name)
  }
}

// 自动展开包含活跃节点的分组
watch(() => props.activeNode, (name) => {
  for (const group of nodeGroups.value) {
    if (group.root === name || group.children.includes(name)) {
      expandedGroups[group.root] = true
    }
  }
})
</script>

<style scoped>
.nodes-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-right: 2px dashed rgba(37, 42, 61, 0.22);
  padding-right: 12px;
}

.nodes-toolbar {
  display: flex;
  gap: 6px;
  margin-bottom: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 120px;
  font-size: 12px;
}

/* 流程图模式 */
.flow-container {
  flex: 1;
  overflow: auto;
  border: 2px solid #252a3d;
  border-radius: 8px;
  background:
    linear-gradient(90deg, rgba(37, 42, 61, 0.05) 1px, transparent 1px),
    linear-gradient(rgba(37, 42, 61, 0.05) 1px, transparent 1px),
    #fffef8;
  background-size: 24px 24px;
  position: relative;
}

.flow-svg {
  min-width: 100%;
  min-height: 100%;
}

.flow-line {
  fill: none;
  stroke-width: 3;
}

.flow-line-next {
  stroke: #2c3145;
}

.flow-line-interrupt {
  stroke: #ff5f74;
  stroke-dasharray: 5,5;
}

.flow-line-wait {
  stroke: #d89a00;
  stroke-dasharray: 2,3;
}

.flow-node {
  cursor: pointer;
  transition: transform 0.15s;
}

.flow-node:hover .node-rect {
  stroke: #252a3d;
  filter: drop-shadow(4px 4px 0 rgba(255, 127, 166, 0.85));
}

.flow-node.active .node-rect {
  fill: #fff3b8;
  stroke: #252a3d;
  stroke-width: 2;
  filter: drop-shadow(4px 4px 0 rgba(97, 215, 255, 0.9));
}

.flow-node.root .node-rect {
  fill: #eafff0;
}

.node-rect {
  fill: white;
  stroke: #252a3d;
  stroke-width: 2;
  transition: all 0.15s;
}

.node-label {
  font-size: 11px;
  font-weight: 800;
  fill: #252a3d;
}

.node-type {
  font-size: 9px;
  fill: #697086;
}

.run-btn-svg {
  fill: #6ee79b;
  stroke: #252a3d;
  stroke-width: 2;
  opacity: 0;
  cursor: pointer;
  transition: opacity 0.15s;
}

.flow-node:hover .run-btn-svg {
  opacity: 1;
}

.run-btn-svg:hover {
  fill: #ffd84d;
}

.run-icon {
  font-size: 8px;
  fill: #252a3d;
  pointer-events: none;
  opacity: 0;
}

.flow-node:hover .run-icon {
  opacity: 1;
}

.flow-legend {
  position: absolute;
  bottom: 8px;
  left: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  border: 2px solid #252a3d;
  background: rgba(255,253,246,0.94);
  padding: 6px 8px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 800;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.legend-line {
  width: 20px;
  height: 2px;
}

.legend-line.next {
  background: #252a3d;
}

.legend-line.interrupt {
  background: repeating-linear-gradient(90deg, #ff5f74 0, #ff5f74 4px, transparent 4px, transparent 8px);
}

.legend-line.wait {
  background: repeating-linear-gradient(90deg, #ffd84d 0, #ffd84d 2px, transparent 2px, transparent 5px);
}

/* 列表模式 */
.nodes-tree {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;
}

.nodes-tree::-webkit-scrollbar { width: 5px; }
.nodes-tree::-webkit-scrollbar-track { background: #f5f5f5; }
.nodes-tree::-webkit-scrollbar-thumb { background: #ccc; border-radius: 3px; }

.node-group {
  margin-bottom: 4px;
}

.node-item {
  display: flex;
  align-items: center;
  min-height: 36px;
  padding: 7px 9px;
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  gap: 6px;
  font-size: 13px;
}

.node-item:hover {
  border-color: #252a3d;
  background: #fff3b8;
}

.node-item.active {
  border-color: #252a3d;
  background: #e9fbff;
  box-shadow: 3px 3px 0 #61d7ff;
}

.root-node {
  font-weight: 500;
}

.child-node {
  padding-left: 28px;
  position: relative;
  color: #4c5368;
}

.child-line {
  position: absolute;
  left: 16px;
  top: 0;
  bottom: 50%;
  width: 1px;
  background: #252a3d;
}

.child-node::before {
  content: '';
  position: absolute;
  left: 16px;
  top: 50%;
  width: 8px;
  height: 1px;
  background: #252a3d;
}

.standalone-node {
  color: #666;
  padding-left: 12px;
}

.expand-icon {
  width: 16px;
  font-size: 10px;
  color: #697086;
  flex-shrink: 0;
}

.node-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-badge {
  width: 20px;
  height: 20px;
  border: 2px solid #252a3d;
  border-radius: 6px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: bold;
  color: #252a3d;
  flex-shrink: 0;
}

.root-node .node-badge {
  background: #61d7ff;
  color: #252a3d;
}

.run-btn {
  padding: 3px 7px;
  border: 2px solid #252a3d;
  border-radius: 6px;
  background: #6ee79b;
  color: #252a3d;
  cursor: pointer;
  font-size: 10px;
  font-weight: 900;
  opacity: 0;
  transition: opacity 0.15s;
}

.node-item:hover .run-btn {
  opacity: 1;
}

.standalone-section {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 2px dashed rgba(37, 42, 61, 0.18);
}

.section-title {
  padding: 6px 8px;
  font-size: 12px;
  color: #697086;
  font-weight: 800;
  cursor: pointer;
}

.section-title:hover {
  color: #252a3d;
}

@media (max-width: 620px) {
  .nodes-toolbar .search-input {
    order: 2;
    flex: 1 0 100%;
    min-width: 0;
  }

  .nodes-toolbar .btn-icon {
    order: 1;
  }
}
</style>
