<template>
  <div class="card">
    <div class="card-titlebar">
      <div>
        <p class="eyebrow">RESULT LOG</p>
        <h2>执行结果</h2>
      </div>
    </div>
    <pre class="json-block">{{ JSON.stringify(result, null, 2) }}</pre>
  </div>
</template>

<script setup>
defineProps({
  result: { type: Object, default: null }
})
</script>
