<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-titlebar">
        <div>
          <p class="eyebrow">NEW MISSION</p>
          <h2>新建 Pipeline</h2>
        </div>
        <button class="btn btn-secondary btn-sm" @click="$emit('close')">关闭</button>
      </div>
      <div class="form-group">
        <label>Pipeline 名称</label>
        <input
          :value="name"
          @input="$emit('update:name', $event.target.value)"
          placeholder="输入名称"
          @keyup.enter="$emit('create')"
        >
      </div>
      <div class="modal-actions">
        <button class="btn btn-primary" @click="$emit('create')">创建</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  name: { type: String, default: '' }
})

defineEmits(['update:name', 'close', 'create'])
</script>
