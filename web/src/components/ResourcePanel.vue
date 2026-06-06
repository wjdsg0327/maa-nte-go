<template>
  <div class="card resource-card motion-card">
    <div class="card-titlebar">
      <div>
        <p class="eyebrow">RESOURCE KIT</p>
        <h2>开放资源包</h2>
      </div>
      <button class="btn btn-primary" :disabled="reloading" @click="$emit('reload')">
        {{ reloading ? '重载中' : '热重载' }}
      </button>
    </div>

    <div class="resource-path">
      <span>当前资源目录</span>
      <strong>{{ config.resourceDir || '-' }}</strong>
    </div>

    <div class="resource-metrics">
      <div>
        <span>Pipeline</span>
        <strong>{{ pipelineCount }}</strong>
      </div>
      <div>
        <span>图片</span>
        <strong>{{ imageCount }}</strong>
      </div>
      <div>
        <span>OCR</span>
        <strong>{{ ocrCount }}</strong>
      </div>
      <div>
        <span>Detect</span>
        <strong>{{ detectCount }}</strong>
      </div>
    </div>

    <div class="config-stack">
      <div>
        <span>监听</span>
        <code>{{ config.serverAddr || '-' }}</code>
      </div>
      <div>
        <span>自动连接</span>
        <code>{{ config.autoConnect ? config.targetWindowTitle || '已开启' : '关闭' }}</code>
      </div>
      <div>
        <span>日志目录</span>
        <code>{{ config.logDir || '-' }}</code>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  config: { type: Object, default: () => ({}) },
  pipelineCount: { type: Number, default: 0 },
  imageCount: { type: Number, default: 0 },
  ocrCount: { type: Number, default: 0 },
  detectCount: { type: Number, default: 0 },
  reloading: { type: Boolean, default: false }
})

defineEmits(['reload'])
</script>
