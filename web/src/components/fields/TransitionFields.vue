<template>
  <!-- Next Nodes -->
  <div class="field-group">
    <label>下一个节点</label>
    <div class="next-tags">
      <span v-for="(next, idx) in node.next || []" :key="idx" class="next-tag">
        {{ next }}
        <span class="remove" @click="removeNext(idx)">×</span>
      </span>
    </div>
    <div class="add-next-input">
      <select v-model="newNextName" class="field-input" style="flex: 1;">
        <option value="">选择节点...</option>
        <option v-for="n in availableNodes" :key="n" :value="n">{{ n }}</option>
      </select>
      <button class="btn btn-primary btn-sm" @click="addNext" :disabled="!newNextName">添加</button>
    </div>
    <div class="field-hint">节点执行完成后跳转到的下一个节点</div>
  </div>

  <!-- Interrupt Nodes -->
  <div class="field-group">
    <label>中断节点</label>
    <div class="next-tags">
      <span v-for="(item, idx) in node.interrupt || []" :key="idx" class="next-tag">
        {{ item }}
        <span class="remove" @click="removeInterrupt(idx)">×</span>
      </span>
    </div>
    <div class="add-next-input">
      <select v-model="newInterruptName" class="field-input" style="flex: 1;">
        <option value="">选择节点...</option>
        <option v-for="n in availableNodes" :key="n" :value="n">{{ n }}</option>
      </select>
      <button class="btn btn-primary btn-sm" @click="addInterrupt" :disabled="!newInterruptName">添加</button>
    </div>
    <div class="field-hint">每次执行前检查的中断节点列表</div>
  </div>

  <!-- Wait Freezes Nodes -->
  <div class="field-group">
    <label>等待节点</label>
    <div class="next-tags">
      <span v-for="(item, idx) in node.wait_freezes || []" :key="idx" class="next-tag">
        {{ item }}
        <span class="remove" @click="removeWaitFreeze(idx)">×</span>
      </span>
    </div>
    <div class="add-next-input">
      <select v-model="newWaitFreezeName" class="field-input" style="flex: 1;">
        <option value="">选择节点...</option>
        <option v-for="n in availableNodes" :key="n" :value="n">{{ n }}</option>
      </select>
      <button class="btn btn-primary btn-sm" @click="addWaitFreeze" :disabled="!newWaitFreezeName">添加</button>
    </div>
    <div class="field-hint">等待画面稳定后再执行</div>
  </div>

  <!-- Reverse Nodes -->
  <div class="field-group">
    <label>反向节点</label>
    <div class="next-tags">
      <span v-for="(item, idx) in node.reverse || []" :key="idx" class="next-tag">
        {{ item }}
        <span class="remove" @click="removeReverse(idx)">×</span>
      </span>
    </div>
    <div class="add-next-input">
      <select v-model="newReverseName" class="field-input" style="flex: 1;">
        <option value="">选择节点...</option>
        <option v-for="n in availableNodes" :key="n" :value="n">{{ n }}</option>
      </select>
      <button class="btn btn-primary btn-sm" @click="addReverse" :disabled="!newReverseName">添加</button>
    </div>
    <div class="field-hint">当识别失败时跳转的节点列表</div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  node: { type: Object, required: true },
  allNodes: { type: Array, default: () => [] }
})

const newNextName = ref('')
const newInterruptName = ref('')
const newWaitFreezeName = ref('')
const newReverseName = ref('')

const availableNodes = computed(() => props.allNodes)

const addNext = () => {
  if (!newNextName.value) return
  if (!props.node.next) props.node.next = []
  if (!props.node.next.includes(newNextName.value)) {
    props.node.next.push(newNextName.value)
  }
  newNextName.value = ''
}

const removeNext = (idx) => {
  if (props.node.next) props.node.next.splice(idx, 1)
}

const addInterrupt = () => {
  if (!newInterruptName.value) return
  if (!props.node.interrupt) props.node.interrupt = []
  if (!props.node.interrupt.includes(newInterruptName.value)) {
    props.node.interrupt.push(newInterruptName.value)
  }
  newInterruptName.value = ''
}

const removeInterrupt = (idx) => {
  if (props.node.interrupt) props.node.interrupt.splice(idx, 1)
}

const addWaitFreeze = () => {
  if (!newWaitFreezeName.value) return
  if (!props.node.wait_freezes) props.node.wait_freezes = []
  if (!props.node.wait_freezes.includes(newWaitFreezeName.value)) {
    props.node.wait_freezes.push(newWaitFreezeName.value)
  }
  newWaitFreezeName.value = ''
}

const removeWaitFreeze = (idx) => {
  if (props.node.wait_freezes) props.node.wait_freezes.splice(idx, 1)
}

const addReverse = () => {
  if (!newReverseName.value) return
  if (!props.node.reverse) props.node.reverse = []
  if (!props.node.reverse.includes(newReverseName.value)) {
    props.node.reverse.push(newReverseName.value)
  }
  newReverseName.value = ''
}

const removeReverse = (idx) => {
  if (props.node.reverse) props.node.reverse.splice(idx, 1)
}
</script>
