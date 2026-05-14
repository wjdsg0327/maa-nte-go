<template>
  <div class="card">
    <div class="card-titlebar">
      <div>
        <p class="eyebrow">TARGET</p>
        <h2>窗口连接</h2>
      </div>
      <button class="btn btn-primary" @click="$emit('refresh')">刷新窗口</button>
    </div>
    <div class="window-list">
      <div
        v-for="w in windows"
        :key="w.handle"
        class="window-item"
        :class="{ connected: connectedHandle === w.handle }"
        @click="$emit('connect', w)"
      >
        <div class="window-main">
          <strong>{{ w.title }}</strong>
          <div class="text-muted">{{ w.class }}</div>
        </div>
        <button class="btn btn-success btn-sm window-action" v-if="connectedHandle === w.handle">ONLINE</button>
        <button class="btn btn-primary btn-sm window-action" v-else>连接</button>
      </div>
      <div v-if="windows.length === 0" class="empty-state">
        <div class="empty-icon">!</div>
        <p>点击刷新窗口，锁定目标程序。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  windows: { type: Array, default: () => [] },
  connectedHandle: { type: [Number, String, null], default: null }
})

defineEmits(['refresh', 'connect'])
</script>
