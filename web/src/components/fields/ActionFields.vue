<template>
  <!-- Action Type -->
  <div class="field-group">
    <label>动作类型</label>
    <select v-model="node.action" class="field-input">
      <option value="DoNothing">DoNothing - 无操作</option>
      <option value="Click">Click - 点击</option>
      <option value="LongPress">LongPress - 长按</option>
      <option value="Swipe">Swipe - 滑动</option>
      <option value="MultiSwipe">MultiSwipe - 多指滑动</option>
      <option value="Scroll">Scroll - 滚轮</option>
      <option value="ClickKey">ClickKey - 单击按键</option>
      <option value="LongPressKey">LongPressKey - 长按按键</option>
      <option value="KeyDown">KeyDown - 按下按键</option>
      <option value="KeyUp">KeyUp - 释放按键</option>
      <option value="InputText">InputText - 输入文字</option>
      <option value="StartApp">StartApp - 启动应用</option>
      <option value="StopApp">StopApp - 关闭应用</option>
      <option value="StopTask">StopTask - 停止任务</option>
      <option value="Command">Command - 执行命令</option>
      <option value="Shell">Shell - ADB命令</option>
      <option value="Screencap">Screencap - 保存截图</option>
      <option value="TouchDown">TouchDown - 触点按下</option>
      <option value="TouchMove">TouchMove - 触点移动</option>
      <option value="TouchUp">TouchUp - 触点抬起</option>
      <option value="Custom">Custom - 自定义</option>
    </select>
  </div>

  <!-- Click Fields -->
  <template v-if="node.action === 'Click'">
    <div class="field-group">
      <label>点击目标</label>
      <input v-model="targetInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
      <div class="field-hint">true=识别结果, [x,y]固定点, [x,y,w,h]区域内随机</div>
    </div>
    <div class="field-group">
      <label>目标偏移</label>
      <input v-model="targetOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
    <div class="field-group">
      <label>触点编号</label>
      <input type="number" v-model.number="node.contact" min="0" placeholder="0">
    </div>
    <div class="field-group">
      <label>触点力度</label>
      <input type="number" v-model.number="node.pressure" placeholder="1">
    </div>
  </template>

  <!-- LongPress Fields -->
  <template v-if="node.action === 'LongPress'">
    <div class="field-group">
      <label>长按目标</label>
      <input v-model="targetInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
    </div>
    <div class="field-group">
      <label>目标偏移</label>
      <input v-model="targetOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
    <div class="field-group">
      <label>长按时长(ms)</label>
      <input type="number" v-model.number="node.duration" min="0" placeholder="1000">
    </div>
    <div class="field-group">
      <label>触点编号</label>
      <input type="number" v-model.number="node.contact" min="0" placeholder="0">
    </div>
  </template>

  <!-- Swipe Fields -->
  <template v-if="node.action === 'Swipe'">
    <div class="field-group">
      <label>滑动起点</label>
      <input v-model="swipeBeginInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
    </div>
    <div class="field-group">
      <label>起点偏移</label>
      <input v-model="swipeBeginOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
    <div class="field-group">
      <label>滑动终点</label>
      <input v-model="swipeEndInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
    </div>
    <div class="field-group">
      <label>终点偏移</label>
      <input v-model="swipeEndOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
    <div class="field-group">
      <label>滑动时长(ms)</label>
      <input type="number" v-model.number="node.duration" min="0" placeholder="200">
    </div>
    <div class="field-group">
      <label>终点停留(ms)</label>
      <input type="number" v-model.number="node.end_hold" min="0" placeholder="0">
    </div>
    <div class="field-group">
      <div class="checkbox-group">
        <input type="checkbox" v-model="node.only_hover" id="onlyHover">
        <label for="onlyHover" style="margin-bottom: 0;">仅悬停移动</label>
      </div>
    </div>
  </template>

  <!-- ClickKey/LongPressKey Fields -->
  <template v-if="['ClickKey', 'LongPressKey'].includes(node.action)">
    <div class="field-group">
      <label>按键</label>
      <input v-model="keyInput" placeholder="27 或 [27,13]">
    </div>
    <div class="field-group" v-if="node.action === 'LongPressKey'">
      <label>长按时长(ms)</label>
      <input type="number" v-model.number="node.duration" min="0" placeholder="1000">
    </div>
  </template>

  <!-- KeyDown/KeyUp Fields -->
  <template v-if="['KeyDown', 'KeyUp'].includes(node.action)">
    <div class="field-group">
      <label>按键</label>
      <input type="number" v-model.number="node.key" placeholder="虚拟按键码">
    </div>
  </template>

  <!-- Scroll Fields -->
  <template v-if="node.action === 'Scroll'">
    <div class="field-group">
      <label>滚动目标</label>
      <input v-model="targetInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
    </div>
    <div class="field-group">
      <label>目标偏移</label>
      <input v-model="targetOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
    <div class="field-group">
      <label>水平滚动</label>
      <input type="number" v-model.number="node.dx" placeholder="0">
      <div class="field-hint">正值向右，负值向左</div>
    </div>
    <div class="field-group">
      <label>垂直滚动</label>
      <input type="number" v-model.number="node.dy" placeholder="0">
      <div class="field-hint">正值向上，负值向下</div>
    </div>
  </template>

  <!-- TouchDown/TouchMove/TouchUp Fields -->
  <template v-if="['TouchDown', 'TouchMove', 'TouchUp'].includes(node.action)">
    <div class="field-group">
      <label>触点编号</label>
      <input type="number" v-model.number="node.contact" min="0" placeholder="0">
    </div>
    <div class="field-group" v-if="node.action !== 'TouchUp'">
      <label>目标位置</label>
      <input v-model="targetInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
    </div>
    <div class="field-group" v-if="node.action !== 'TouchUp'">
      <label>目标偏移</label>
      <input v-model="targetOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
    <div class="field-group" v-if="node.action !== 'TouchUp'">
      <label>触控压力</label>
      <input type="number" v-model.number="node.pressure" placeholder="0">
    </div>
  </template>

  <!-- InputText Fields -->
  <template v-if="node.action === 'InputText'">
    <div class="field-group">
      <label>输入文字</label>
      <input v-model="node.input_text" placeholder="要输入的文字">
    </div>
  </template>

  <!-- StartApp/StopApp Fields -->
  <template v-if="node.action === 'StartApp' || node.action === 'StopApp'">
    <div class="field-group">
      <label>{{ node.action === 'StartApp' ? '启动入口' : '包名' }}</label>
      <input v-model="node.package" placeholder="com.example.app">
    </div>
  </template>

  <!-- Command Fields -->
  <template v-if="node.action === 'Command'">
    <div class="field-group">
      <label>执行程序</label>
      <input v-model="node.exec" placeholder="Python 或脚本路径">
    </div>
    <div class="field-group">
      <label>命令参数</label>
      <input v-model="commandArgsInput" placeholder='["arg1","{IMAGE}"]'>
    </div>
    <div class="field-group">
      <div class="checkbox-group">
        <input type="checkbox" v-model="node.detach" id="detach">
        <label for="detach" style="margin-bottom: 0;">分离模式</label>
      </div>
    </div>
  </template>

  <!-- Shell Fields -->
  <template v-if="node.action === 'Shell'">
    <div class="field-group">
      <label>Shell命令</label>
      <input v-model="node.cmd" placeholder="getprop ro.build.version.sdk">
    </div>
    <div class="field-group">
      <label>超时时间(ms)</label>
      <input type="number" v-model.number="node.shell_timeout" min="-1" placeholder="20000">
    </div>
  </template>

  <!-- Screencap Fields -->
  <template v-if="node.action === 'Screencap'">
    <div class="field-group">
      <label>文件名</label>
      <input v-model="node.filename" placeholder="可选，默认时间戳">
    </div>
    <div class="field-group">
      <label>图片格式</label>
      <select v-model="node.format" class="field-input">
        <option value="">默认(png)</option>
        <option value="png">png</option>
        <option value="jpg">jpg</option>
      </select>
    </div>
    <div class="field-group" v-if="node.format === 'jpg'">
      <label>图片质量</label>
      <input type="number" v-model.number="node.quality" min="0" max="100" placeholder="100">
    </div>
  </template>

  <!-- Custom Action Fields -->
  <template v-if="node.action === 'Custom'">
    <div class="field-group">
      <label>动作名称</label>
      <input v-model="node.custom_action" placeholder="注册的自定义动作名">
    </div>
    <div class="field-group">
      <label>动作参数</label>
      <input v-model="customActionParamInput" placeholder='任意JSON参数'>
    </div>
    <div class="field-group">
      <label>目标位置</label>
      <input v-model="targetInput" placeholder="true/节点名/[x,y]/[x,y,w,h]">
    </div>
    <div class="field-group">
      <label>目标偏移</label>
      <input v-model="targetOffsetInput" placeholder="[x,y,w,h] 可选">
    </div>
  </template>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  node: { type: Object, required: true }
})

// Input bindings
const targetInput = ref('')
const targetOffsetInput = ref('')
const swipeBeginInput = ref('')
const swipeBeginOffsetInput = ref('')
const swipeEndInput = ref('')
const swipeEndOffsetInput = ref('')
const keyInput = ref('')
const commandArgsInput = ref('')
const customActionParamInput = ref('')

// Initialize from node
watch(() => props.node, (node) => {
  targetInput.value = node.target !== undefined ? JSON.stringify(node.target) : ''
  targetOffsetInput.value = node.target_offset ? JSON.stringify(node.target_offset) : ''
  swipeBeginInput.value = node.begin ? JSON.stringify(node.begin) : ''
  swipeBeginOffsetInput.value = node.begin_offset ? JSON.stringify(node.begin_offset) : ''
  swipeEndInput.value = node.end ? JSON.stringify(node.end) : ''
  swipeEndOffsetInput.value = node.end_offset ? JSON.stringify(node.end_offset) : ''
  keyInput.value = node.key !== undefined ? JSON.stringify(node.key) : ''
  commandArgsInput.value = node.args ? JSON.stringify(node.args) : ''
  customActionParamInput.value = node.custom_action_param ? JSON.stringify(node.custom_action_param) : ''
}, { immediate: true })

// Sync inputs to node
watch(targetInput, (val) => {
  try {
    const parsed = JSON.parse(val)
    props.node.target = parsed
  } catch (e) {
    if (val === 'true') props.node.target = true
    else if (val === 'false') props.node.target = false
    else props.node.target = val || undefined
  }
})
watch(targetOffsetInput, (val) => {
  try { props.node.target_offset = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(swipeBeginInput, (val) => {
  try { props.node.begin = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(swipeBeginOffsetInput, (val) => {
  try { props.node.begin_offset = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(swipeEndInput, (val) => {
  try { props.node.end = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(swipeEndOffsetInput, (val) => {
  try { props.node.end_offset = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(keyInput, (val) => {
  try { props.node.key = val ? JSON.parse(val) : undefined } catch (e) { props.node.key = val || undefined }
})
watch(commandArgsInput, (val) => {
  try { props.node.args = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(customActionParamInput, (val) => {
  try { props.node.custom_action_param = val ? JSON.parse(val) : undefined } catch (e) {}
})
</script>
