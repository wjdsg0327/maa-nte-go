<template>
  <div class="toast" :class="type">{{ message }}</div>
</template>

<script setup>
defineProps({
  message: { type: String, default: '' },
  type: { type: String, default: 'success' }
})
</script>
