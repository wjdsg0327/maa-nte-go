<template>
  <div class="flow-board">
    <div class="flow-board-toolbar">
      <button class="btn btn-primary btn-sm" @click="$emit('add')">+ 新建节点</button>
      <input v-model="searchText" class="search-input" placeholder="搜索节点...">
      <div class="flow-counts">
        <span>{{ visibleNames.length }} 节点</span>
        <span>{{ visibleLinks.length }} 连线</span>
      </div>
    </div>

    <div class="flow-stage">
      <svg class="flow-stage-svg" :width="layout.width" :height="layout.height">
        <defs>
          <marker id="board-arrow-next" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#252a3d"/>
          </marker>
          <marker id="board-arrow-interrupt" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#ff5f74"/>
          </marker>
          <marker id="board-arrow-wait" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#d89a00"/>
          </marker>
          <marker id="board-arrow-reverse" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
            <path d="M0,0 L0,6 L9,3 z" fill="#a79bff"/>
          </marker>
        </defs>

        <path
          v-for="link in visibleLinks"
          :key="`${link.from}-${link.type}-${link.to}`"
          :d="getPath(link)"
          class="board-flow-line"
          :class="`board-flow-line-${link.type}`"
          :marker-end="`url(#board-arrow-${link.type})`"
        />

        <g
          v-for="name in visibleNames"
          :key="name"
          class="board-flow-node"
          :class="{ active: activeNode === name, root: rootSet.has(name) }"
          :transform="`translate(${layout.positions[name].x}, ${layout.positions[name].y})`"
          @click="$emit('select', name)"
        >
          <rect class="board-node-rect" :width="nodeWidth" :height="nodeHeight" rx="8" ry="8" />
          <text class="board-node-kicker" x="14" y="20">
            {{ rootSet.has(name) ? 'ROOT' : 'NODE' }}
          </text>
          <text class="board-node-name" x="14" y="43">
            {{ truncate(name, 18) }}
          </text>
          <text class="board-node-meta" x="14" y="66">
            {{ getTypeLabel(nodes[name]) }}
          </text>
          <circle
            class="board-run-btn"
            :cx="nodeWidth - 18"
            cy="18"
            r="11"
            @click.stop="$emit('execute', name)"
          />
          <text
            class="board-run-icon"
            :x="nodeWidth - 18"
            y="22"
            text-anchor="middle"
            @click.stop="$emit('execute', name)"
          >▶</text>
        </g>
      </svg>

      <div class="flow-legend flow-board-legend">
        <span class="legend-item"><span class="legend-line next"></span> Next</span>
        <span class="legend-item"><span class="legend-line interrupt"></span> Interrupt</span>
        <span class="legend-item"><span class="legend-line wait"></span> Wait</span>
        <span class="legend-item"><span class="legend-line reverse"></span> Reverse</span>
      </div>

      <div v-if="visibleNames.length === 0" class="flow-empty empty-state">
        <div class="empty-icon">+</div>
        <p>暂无匹配节点，清空搜索或新建节点。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { getRootNodes, layoutGraph } from '../utils/pipelineGraph.js'

const props = defineProps({
  nodes: { type: Object, required: true },
  activeNode: { type: String, default: '' },
})

defineEmits(['select', 'add', 'execute'])

const nodeWidth = 190
const nodeHeight = 84
const searchText = ref('')

const filteredNodes = computed(() => {
  const search = searchText.value.trim().toLowerCase()
  if (!search) return props.nodes

  const result = {}
  for (const [name, node] of Object.entries(props.nodes)) {
    const text = `${name} ${node?.recognition || ''} ${node?.action || ''}`.toLowerCase()
    if (text.includes(search)) result[name] = node
  }
  return result
})

const layout = computed(() => layoutGraph(filteredNodes.value, {
  nodeWidth,
  nodeHeight,
  gapX: 96,
  gapY: 110,
}))

const visibleNames = computed(() => Object.keys(layout.value.positions))
const visibleLinks = computed(() => layout.value.links)
const rootSet = computed(() => new Set(getRootNodes(filteredNodes.value)))

const getPath = (link) => {
  const from = layout.value.positions[link.from]
  const to = layout.value.positions[link.to]
  if (!from || !to) return ''

  const x1 = from.x + nodeWidth / 2
  const y1 = from.y + nodeHeight
  const x2 = to.x + nodeWidth / 2
  const y2 = to.y
  const midY = (y1 + y2) / 2

  if (Math.abs(x1 - x2) < 12) {
    return `M${x1},${y1} C${x1},${midY} ${x2},${midY} ${x2},${y2}`
  }

  const offset = link.type === 'interrupt' ? 18 : link.type === 'wait' ? -18 : link.type === 'reverse' ? 30 : 0
  return `M${x1 + offset},${y1} C${x1 + offset},${midY} ${x2 + offset},${midY} ${x2},${y2}`
}

const truncate = (value, max) => {
  return value.length > max ? `${value.slice(0, max - 1)}...` : value
}

const getTypeLabel = (node) => {
  const recognition = node?.recognition?.type || node?.recognition || 'DirectHit'
  const action = node?.action?.type || node?.action || 'DoNothing'
  return `${recognition} -> ${action}`
}
</script>
