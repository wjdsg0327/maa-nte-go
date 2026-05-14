<template>
  <div class="card">
    <div class="card-titlebar">
      <div>
        <p class="eyebrow">MISSION DECK</p>
        <h2>Pipeline 管理</h2>
      </div>
      <div class="toolbar-actions">
        <button class="btn btn-primary" @click="$emit('create')">新建</button>
        <button class="btn btn-secondary" @click="$emit('refresh')">刷新列表</button>
      </div>
    </div>
    <div class="pipeline-list">
      <div
        v-for="p in pipelines"
        :key="p.name"
        class="pipeline-item"
        :class="{ selected: selectedPipeline === p.name }"
        @click="$emit('select', p.name)"
      >
        <div class="pipeline-header">
          <span class="pipeline-chip">RUN</span>
          <h3>{{ p.name }}</h3>
        </div>
        <div class="actions">
          <button class="btn btn-success" @click.stop="$emit('execute', p.name)">执行</button>
          <button class="btn btn-primary" @click.stop="$emit('edit', p.name)">编辑</button>
          <button class="btn btn-danger" @click.stop="$emit('delete', p.name)">删除</button>
        </div>
      </div>
      <div v-if="pipelines.length === 0" class="empty-state">
        <div class="empty-icon">+</div>
        <p>暂无 Pipeline，创建一个新任务开始编排。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  pipelines: { type: Array, default: () => [] },
  selectedPipeline: { type: String, default: null }
})

defineEmits(['refresh', 'select', 'create', 'edit', 'execute', 'delete'])
</script>
