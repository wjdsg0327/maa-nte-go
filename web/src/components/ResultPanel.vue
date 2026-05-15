<template>
  <div class="card">
    <div class="card-titlebar">
      <div>
        <p class="eyebrow">RESULT LOG</p>
        <h2>执行结果</h2>
      </div>
      <span class="status-badge" :class="result.status === 'success' ? 'connected' : 'disconnected'">
        {{ result.status || 'unknown' }}
      </span>
    </div>
    <div class="result-summary">
      <div>
        <span>任务</span>
        <strong>{{ result.task || '-' }}</strong>
      </div>
      <div>
        <span>入口</span>
        <strong>{{ result.entry || result.task || '-' }}</strong>
      </div>
      <div>
        <span>节点数</span>
        <strong>{{ nodes.length || result.node_count || 0 }}</strong>
      </div>
    </div>

    <div v-if="nodes.length" class="debug-node-list">
      <article v-for="node in nodes" :key="node.id || node.name" class="debug-node-card">
        <div class="debug-node-title">
          <strong>{{ node.name || `#${node.id}` }}</strong>
          <span class="mini-chip" :class="node.runCompleted || node.run_completed ? 'ok' : 'warn'">
            {{ node.runCompleted || node.run_completed ? '完成' : '未完成' }}
          </span>
        </div>
        <div v-if="node.recognition" class="debug-grid">
          <span>识别</span>
          <strong>{{ node.recognition.algorithm || '-' }}</strong>
          <span>命中</span>
          <strong>{{ node.recognition.hit ? '是' : '否' }}</strong>
          <span>区域</span>
          <code>{{ formatBox(node.recognition.box) }}</code>
        </div>
        <div v-if="node.action" class="debug-grid">
          <span>动作</span>
          <strong>{{ node.action.action || '-' }}</strong>
          <span>成功</span>
          <strong>{{ node.action.success ? '是' : '否' }}</strong>
          <span>区域</span>
          <code>{{ formatBox(node.action.box) }}</code>
        </div>
      </article>
    </div>

    <details class="raw-json-details">
      <summary>原始 JSON</summary>
      <pre class="json-block">{{ JSON.stringify(result, null, 2) }}</pre>
    </details>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  result: { type: Object, default: null }
})

const nodes = computed(() => props.result?.task_detail?.nodes || [])

const formatBox = (box) => Array.isArray(box) ? `[${box.join(',')}]` : '-'
</script>
