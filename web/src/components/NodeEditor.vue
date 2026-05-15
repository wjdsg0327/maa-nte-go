<template>
  <div class="node-editor">
    <div class="node-card active">
      <div class="node-header">
        <input
          :value="name"
          @change="$emit('update:name', $event.target.value)"
          placeholder="节点名称"
        >
        <button class="btn btn-danger btn-sm" @click="$emit('delete')">删除节点</button>
      </div>
      <div class="node-body">
        <section class="editor-section">
          <div class="editor-section-title">RECOGNITION</div>
          <RecognitionFields
            :node="node"
            :images="images"
            :ocr-models="ocrModels"
            :detect-models="detectModels"
            @pick-roi="$emit('pick-roi')"
          />
        </section>

        <section class="editor-section">
          <div class="editor-section-title">ACTION</div>
          <ActionFields :node="node" />
        </section>

        <section class="editor-section">
          <div class="editor-section-title">COMMON</div>
          <CommonFields :node="node" />
        </section>

        <section class="editor-section">
          <div class="editor-section-title">ROUTE</div>
          <TransitionFields :node="node" :all-nodes="allNodes" />
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import RecognitionFields from './fields/RecognitionFields.vue'
import ActionFields from './fields/ActionFields.vue'
import CommonFields from './fields/CommonFields.vue'
import TransitionFields from './fields/TransitionFields.vue'

defineProps({
  node: { type: Object, required: true },
  name: { type: String, required: true },
  images: { type: Array, default: () => [] },
  ocrModels: { type: Array, default: () => [] },
  detectModels: { type: Array, default: () => [] },
  allNodes: { type: Array, default: () => [] }
})

defineEmits(['update:name', 'delete', 'pick-roi'])
</script>
