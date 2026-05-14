<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal large">
      <div class="modal-titlebar">
        <div>
          <p class="eyebrow">SOURCE VIEW</p>
          <h2>JSON 预览</h2>
        </div>
        <button class="btn btn-secondary btn-sm" @click="$emit('close')">关闭</button>
      </div>
      <pre class="json-block preview-block">{{ JSON.stringify(data, null, 2) }}</pre>
      <div class="modal-actions">
        <button class="btn btn-primary" @click="$emit('close')">完成</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  data: { type: Object, default: () => ({}) }
})

defineEmits(['close'])
</script>
