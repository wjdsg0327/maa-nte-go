<template>
  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.save_result" id="saveResult">
      <label for="saveResult" style="margin-bottom: 0;">保存识别结果</label>
    </div>
  </div>
  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.save_draw" id="saveDraw">
      <label for="saveDraw" style="margin-bottom: 0;">保存标记图片</label>
    </div>
    <div class="field-hint">保存带识别框标记的截图到 log/vision 目录</div>
  </div>
  <div class="field-group">
    <label>超时时间(ms)</label>
    <input type="number" v-model.number="node.timeout" min="0" placeholder="留空使用默认值">
  </div>
  <div class="field-group">
    <label>前置延时(ms)</label>
    <input type="number" v-model.number="node.pre_delay" min="0" placeholder="0">
  </div>
  <div class="field-group">
    <label>后置延时(ms)</label>
    <input type="number" v-model.number="node.post_delay" min="0" placeholder="200">
  </div>
  <div class="field-group">
    <label>执行次数</label>
    <input type="number" v-model.number="node.times" min="1" placeholder="1">
  </div>
  <div class="field-group">
    <label>速率限制</label>
    <input type="number" v-model.number="node.rate_limit" min="0" placeholder="1000">
    <div class="field-hint">最短截图间隔(ms)</div>
  </div>
  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.inverse" id="inverse">
      <label for="inverse" style="margin-bottom: 0;">反转结果</label>
    </div>
    <div class="field-hint">识别成功当失败，失败当成功</div>
  </div>
  <div class="field-group">
    <div class="checkbox-group">
      <input type="checkbox" v-model="node.enabled" id="enabled">
      <label for="enabled" style="margin-bottom: 0;">启用节点</label>
    </div>
  </div>
</template>

<script setup>
defineProps({
  node: { type: Object, required: true }
})
</script>
