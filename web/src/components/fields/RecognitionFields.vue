<template>
  <!-- Recognition Type -->
  <div class="field-group">
    <label>识别方式</label>
    <select v-model="node.recognition" class="field-input">
      <option value="DirectHit">DirectHit - 直接命中</option>
      <option value="TemplateMatch">TemplateMatch - 模板匹配</option>
      <option value="FeatureMatch">FeatureMatch - 特征匹配</option>
      <option value="ColorMatch">ColorMatch - 颜色匹配</option>
      <option value="OCR">OCR - 文字识别</option>
      <option value="NeuralNetworkClassify">NeuralNetworkClassify - NN分类</option>
      <option value="NeuralNetworkDetect">NeuralNetworkDetect - NN检测</option>
      <option value="And">And - 组合(与)</option>
      <option value="Or">Or - 组合(或)</option>
      <option value="Custom">Custom - 自定义</option>
    </select>
  </div>

  <!-- ROI Fields -->
  <template v-if="hasRoi">
    <div class="field-group">
      <label>ROI区域</label>
      <div class="roi-input-row">
        <input v-model="roiInput" :class="{ invalid: roiError }" placeholder="[x,y,w,h] 或节点名">
        <button class="btn btn-secondary btn-sm" type="button" @click="$emit('pick-roi')">框选</button>
      </div>
      <div v-if="roiError" class="field-error">{{ roiError }}</div>
      <div class="field-hint">识别区域 [x,y,w,h]，支持负数，0表示延伸到边缘</div>
    </div>
    <div class="field-group">
      <label>ROI偏移</label>
      <input v-model="roiOffsetInput" :class="{ invalid: roiOffsetError }" placeholder="[x,y,w,h] 可选">
      <div v-if="roiOffsetError" class="field-error">{{ roiOffsetError }}</div>
    </div>
  </template>

  <!-- TemplateMatch Fields -->
  <template v-if="node.recognition === 'TemplateMatch'">
    <div class="field-group">
      <label>模板图片</label>
      <div class="template-tags" v-if="templateItems.length > 0">
        <span v-for="(template, idx) in templateItems" :key="template" class="next-tag">
          {{ template }}
          <span class="remove" @click="removeTemplate(idx)">×</span>
        </span>
      </div>
      <div class="template-add-row">
        <select v-model="selectedTemplateImage" class="field-input">
          <option value="">选择图片...</option>
          <option v-for="img in availableTemplateImages" :key="img" :value="img">{{ img }}</option>
        </select>
        <button class="btn btn-primary btn-sm" @click="addSelectedTemplate" :disabled="!selectedTemplateImage">添加图片</button>
      </div>
      <div class="template-add-row">
        <input v-model="templatePathInput" placeholder="输入图片路径或文件夹路径" @keyup.enter="addTemplatePath">
        <button class="btn btn-secondary btn-sm" @click="addTemplatePath" :disabled="!templatePathInput.trim()">添加路径</button>
      </div>
      <div class="field-hint">支持多张图片或文件夹路径；单项保存为字符串，多项保存为数组。</div>
    </div>
    <div class="field-group">
      <label>阈值</label>
      <input v-model="thresholdInput" placeholder="0.7 或 [0.7,0.8]">
      <div class="field-hint">数组阈值长度需与模板数量一致；留空使用默认值 0.7。</div>
    </div>
    <div class="field-group">
      <label>排序方式</label>
      <select v-model="node.order_by" class="field-input">
        <option value="">默认(Horizontal)</option>
        <option value="Horizontal">Horizontal - 水平</option>
        <option value="Vertical">Vertical - 垂直</option>
        <option value="Score">Score - 分数</option>
        <option value="Random">Random - 随机</option>
      </select>
    </div>
    <div class="field-group">
      <label>结果索引</label>
      <input type="number" v-model.number="node.index" placeholder="0">
    </div>
    <div class="field-group">
      <label>匹配算法</label>
      <select v-model.number="node.method" class="field-input">
        <option :value="undefined">默认(5)</option>
        <option :value="5">5 - TM_CCOEFF_NORMED (推荐)</option>
        <option :value="3">3 - TM_CCORR_NORMED</option>
        <option :value="10001">10001 - TM_SQDIFF_NORMED(反转)</option>
      </select>
    </div>
    <div class="field-group">
      <div class="checkbox-group">
        <input type="checkbox" v-model="node.green_mask" id="greenMask">
        <label for="greenMask" style="margin-bottom: 0;">绿色掩码</label>
      </div>
      <div class="field-hint">图片中RGB(0,255,0)绿色区域不参与匹配</div>
    </div>
  </template>

  <!-- FeatureMatch Fields -->
  <template v-if="node.recognition === 'FeatureMatch'">
    <div class="field-group">
      <label>模板图片</label>
      <div class="template-tags" v-if="templateItems.length > 0">
        <span v-for="(template, idx) in templateItems" :key="template" class="next-tag">
          {{ template }}
          <span class="remove" @click="removeTemplate(idx)">×</span>
        </span>
      </div>
      <div class="template-add-row">
        <select v-model="selectedTemplateImage" class="field-input">
          <option value="">选择图片...</option>
          <option v-for="img in availableTemplateImages" :key="img" :value="img">{{ img }}</option>
        </select>
        <button class="btn btn-primary btn-sm" @click="addSelectedTemplate" :disabled="!selectedTemplateImage">添加图片</button>
      </div>
      <div class="template-add-row">
        <input v-model="templatePathInput" placeholder="输入图片路径或文件夹路径" @keyup.enter="addTemplatePath">
        <button class="btn btn-secondary btn-sm" @click="addTemplatePath" :disabled="!templatePathInput.trim()">添加路径</button>
      </div>
      <div class="field-hint">支持多张图片或文件夹路径；单项保存为字符串，多项保存为数组。</div>
    </div>
    <div class="field-group">
      <label>特征点数量</label>
      <input type="number" v-model.number="node.count" min="0" placeholder="4">
    </div>
    <div class="field-group">
      <label>检测器</label>
      <select v-model="node.detector" class="field-input">
        <option value="">默认(SIFT)</option>
        <option value="SIFT">SIFT - 精度最高(推荐)</option>
        <option value="KAZE">KAZE</option>
        <option value="AKAZE">AKAZE</option>
        <option value="BRISK">BRISK</option>
        <option value="ORB">ORB - 最快</option>
      </select>
    </div>
    <div class="field-group">
      <label>匹配比例</label>
      <input type="number" v-model.number="node.ratio" step="0.1" min="0" max="1" placeholder="0.6">
    </div>
    <div class="field-group">
      <label>排序方式</label>
      <select v-model="node.order_by" class="field-input">
        <option value="">默认(Horizontal)</option>
        <option value="Horizontal">Horizontal</option>
        <option value="Vertical">Vertical</option>
        <option value="Score">Score</option>
        <option value="Area">Area</option>
        <option value="Random">Random</option>
      </select>
    </div>
    <div class="field-group">
      <label>结果索引</label>
      <input type="number" v-model.number="node.index" placeholder="0">
    </div>
    <div class="field-group">
      <div class="checkbox-group">
        <input type="checkbox" v-model="node.green_mask" id="greenMaskFM">
        <label for="greenMaskFM" style="margin-bottom: 0;">绿色掩码</label>
      </div>
    </div>
  </template>

  <!-- ColorMatch Fields -->
  <template v-if="node.recognition === 'ColorMatch'">
    <div class="field-group">
      <label>颜色匹配方式</label>
      <select v-model.number="node.method" class="field-input">
        <option :value="undefined">默认(4-RGB)</option>
        <option :value="4">4 - RGB (3通道)</option>
        <option :value="40">40 - HSV (3通道)</option>
        <option :value="6">6 - GRAY (1通道)</option>
      </select>
    </div>
    <div class="field-group">
      <label>取色</label>
      <ColorMatchPicker v-model:lower="node.lower" v-model:upper="node.upper" />
    </div>
    <div class="field-group">
      <label>像素数量阈值</label>
      <input type="number" v-model.number="node.count" min="1" placeholder="1">
    </div>
    <div class="field-group">
      <div class="checkbox-group">
        <input type="checkbox" v-model="node.connected" id="connected">
        <label for="connected" style="margin-bottom: 0;">仅相连像素</label>
      </div>
    </div>
    <div class="field-group">
      <label>排序方式</label>
      <select v-model="node.order_by" class="field-input">
        <option value="">默认(Horizontal)</option>
        <option value="Horizontal">Horizontal</option>
        <option value="Vertical">Vertical</option>
        <option value="Score">Score</option>
        <option value="Area">Area</option>
        <option value="Random">Random</option>
      </select>
    </div>
    <div class="field-group">
      <label>结果索引</label>
      <input type="number" v-model.number="node.index" placeholder="0">
    </div>
  </template>

  <!-- OCR Fields -->
  <template v-if="node.recognition === 'OCR'">
    <div class="field-group">
      <label>期望文字</label>
      <input v-model="node.expected" placeholder="确定 或 正则表达式">
      <div class="field-hint">支持正则，多个用数组 ["OK","确定"]</div>
    </div>
    <div class="field-group">
      <label>置信度阈值</label>
      <input type="number" v-model.number="node.threshold" step="0.01" min="0" max="1" placeholder="0.3">
    </div>
    <div class="field-group">
      <label>文字替换</label>
      <input v-model="replaceInput" placeholder='[["旧1","新1"],["旧2","新2"]]'>
    </div>
    <div class="field-group">
      <label>排序方式</label>
      <select v-model="node.order_by" class="field-input">
        <option value="">默认(Horizontal)</option>
        <option value="Horizontal">Horizontal</option>
        <option value="Vertical">Vertical</option>
        <option value="Area">Area</option>
        <option value="Length">Length - 文本长度</option>
        <option value="Random">Random</option>
        <option value="Expected">Expected - 按期望顺序</option>
      </select>
    </div>
    <div class="field-group">
      <label>结果索引</label>
      <input type="number" v-model.number="node.index" placeholder="0">
    </div>
    <div class="field-group">
      <div class="checkbox-group">
        <input type="checkbox" v-model="node.only_rec" id="onlyRec">
        <label for="onlyRec" style="margin-bottom: 0;">仅识别(不检测)</label>
      </div>
    </div>
    <div class="field-group">
      <label>模型文件夹</label>
      <select v-model="node.model" class="field-input">
        <option value="">选择模型...</option>
        <option v-for="m in ocrModels" :key="m" :value="m">{{ m }}</option>
      </select>
    </div>
    <div class="field-group">
      <label>颜色过滤</label>
      <input v-model="node.color_filter" placeholder="ColorMatch节点名">
    </div>
  </template>

  <!-- NeuralNetworkClassify Fields -->
  <template v-if="node.recognition === 'NeuralNetworkClassify'">
    <div class="field-group">
      <label>模型文件</label>
      <select v-model="node.model" class="field-input">
        <option value="">选择模型...</option>
        <option v-for="m in detectModels" :key="m" :value="m">{{ m }}</option>
      </select>
      <input v-model="node.model" placeholder="或输入自定义路径" style="margin-top: 5px;">
    </div>
    <div class="field-group">
      <label>标签</label>
      <input v-model="labelsInput" placeholder='["Cat","Dog","Mouse"]'>
    </div>
    <div class="field-group">
      <label>期望分类</label>
      <input v-model="expectedIndexInput" placeholder="0 或 [0,2]">
    </div>
    <div class="field-group">
      <label>排序方式</label>
      <select v-model="node.order_by" class="field-input">
        <option value="">默认(Horizontal)</option>
        <option value="Horizontal">Horizontal</option>
        <option value="Vertical">Vertical</option>
        <option value="Score">Score</option>
        <option value="Random">Random</option>
        <option value="Expected">Expected</option>
      </select>
    </div>
    <div class="field-group">
      <label>结果索引</label>
      <input type="number" v-model.number="node.index" placeholder="0">
    </div>
  </template>

  <!-- NeuralNetworkDetect Fields -->
  <template v-if="node.recognition === 'NeuralNetworkDetect'">
    <div class="field-group">
      <label>模型文件</label>
      <select v-model="node.model" class="field-input">
        <option value="">选择模型...</option>
        <option v-for="m in detectModels" :key="m" :value="m">{{ m }}</option>
      </select>
    </div>
    <div class="field-group">
      <label>标签</label>
      <input v-model="labelsInput" placeholder='["Cat","Dog"]'>
    </div>
    <div class="field-group">
      <label>期望分类</label>
      <input v-model="expectedIndexInput" placeholder="0 或 [0,2]">
    </div>
    <div class="field-group">
      <label>置信度阈值</label>
      <input v-model="thresholdInput" placeholder="0.3 或 [0.3,0.5]">
    </div>
    <div class="field-group">
      <label>排序方式</label>
      <select v-model="node.order_by" class="field-input">
        <option value="">默认(Horizontal)</option>
        <option value="Horizontal">Horizontal</option>
        <option value="Vertical">Vertical</option>
        <option value="Score">Score</option>
        <option value="Area">Area</option>
        <option value="Random">Random</option>
        <option value="Expected">Expected</option>
      </select>
    </div>
    <div class="field-group">
      <label>结果索引</label>
      <input type="number" v-model.number="node.index" placeholder="0">
    </div>
  </template>

  <!-- And Fields -->
  <template v-if="node.recognition === 'And'">
    <div class="field-group">
      <label>子识别列表(与)</label>
      <textarea v-model="allOfInput" class="field-input" rows="4" placeholder='[{"recognition":"TemplateMatch","template":"a.png"}]'></textarea>
      <div class="field-hint">所有子识别都命中才算成功</div>
    </div>
    <div class="field-group">
      <label>输出框索引</label>
      <input type="number" v-model.number="node.box_index" min="0" placeholder="0">
    </div>
  </template>

  <!-- Or Fields -->
  <template v-if="node.recognition === 'Or'">
    <div class="field-group">
      <label>子识别列表(或)</label>
      <textarea v-model="anyOfInput" class="field-input" rows="4" placeholder='[{"recognition":"TemplateMatch","template":"a.png"}]'></textarea>
      <div class="field-hint">命中第一个即成功</div>
    </div>
  </template>

  <!-- Custom Recognition Fields -->
  <template v-if="node.recognition === 'Custom'">
    <div class="field-group">
      <label>识别器名称</label>
      <input v-model="node.custom_recognition" placeholder="注册的自定义识别器名">
    </div>
    <div class="field-group">
      <label>识别器参数</label>
      <input v-model="customRecognitionParamInput" placeholder='任意JSON参数'>
    </div>
  </template>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import ColorMatchPicker from '../ColorMatchPicker.vue'
import {
  appendListItem,
  formatNumberOrList,
  formatRoiValue,
  normalizeStringList,
  parseRectArrayValue,
  parseRoiValue,
  serializeNumberOrList,
  serializeStringOrList,
} from '../../utils/pipelineFields.js'

const props = defineProps({
  node: { type: Object, required: true },
  images: { type: Array, default: () => [] },
  ocrModels: { type: Array, default: () => [] },
  detectModels: { type: Array, default: () => [] }
})

defineEmits(['pick-roi'])

// Computed
const hasRoi = computed(() => {
  return ['DirectHit', 'TemplateMatch', 'FeatureMatch', 'ColorMatch', 'OCR', 'NeuralNetworkClassify', 'NeuralNetworkDetect', 'Custom'].includes(props.node.recognition)
})

// Input bindings
const roiInput = ref('')
const roiOffsetInput = ref('')
const roiError = ref('')
const roiOffsetError = ref('')
const replaceInput = ref('')
const labelsInput = ref('')
const expectedIndexInput = ref('')
const thresholdInput = ref('')
const templateItems = ref([])
const selectedTemplateImage = ref('')
const templatePathInput = ref('')
const allOfInput = ref('')
const anyOfInput = ref('')
const customRecognitionParamInput = ref('')

const availableTemplateImages = computed(() => {
  return props.images.filter((img) => !templateItems.value.includes(img))
})

const syncTemplate = () => {
  props.node.template = serializeStringOrList(templateItems.value)
}

const addTemplate = (value) => {
  const nextItems = appendListItem(templateItems.value, value)
  templateItems.value = nextItems
  syncTemplate()
}

const addSelectedTemplate = () => {
  addTemplate(selectedTemplateImage.value)
  selectedTemplateImage.value = ''
}

const addTemplatePath = () => {
  addTemplate(templatePathInput.value)
  templatePathInput.value = ''
}

const removeTemplate = (idx) => {
  templateItems.value = templateItems.value.filter((_, itemIdx) => itemIdx !== idx)
  syncTemplate()
}

// Initialize from node
watch(() => props.node, (node) => {
  templateItems.value = normalizeStringList(node.template)
  selectedTemplateImage.value = ''
  templatePathInput.value = ''
  roiError.value = ''
  roiOffsetError.value = ''
  roiInput.value = formatRoiValue(node.roi)
  roiOffsetInput.value = formatRoiValue(node.roi_offset)
  replaceInput.value = node.replace ? JSON.stringify(node.replace) : ''
  labelsInput.value = node.labels ? JSON.stringify(node.labels) : ''
  expectedIndexInput.value = node.expected ? JSON.stringify(node.expected) : ''
  thresholdInput.value = formatNumberOrList(node.threshold)
  allOfInput.value = node.allOf ? JSON.stringify(node.allOf) : ''
  anyOfInput.value = node.anyOf ? JSON.stringify(node.anyOf) : ''
  customRecognitionParamInput.value = node.custom_recognition_param ? JSON.stringify(node.custom_recognition_param) : ''
}, { immediate: true })

watch(() => props.node.roi, (value) => {
  const formatted = formatRoiValue(value)
  if (formatted !== roiInput.value) {
    roiError.value = ''
    roiInput.value = formatted
  }
})

watch(() => props.node.roi_offset, (value) => {
  const formatted = formatRoiValue(value)
  if (formatted !== roiOffsetInput.value) {
    roiOffsetError.value = ''
    roiOffsetInput.value = formatted
  }
})

// Sync inputs to node
watch(roiInput, (val) => {
  const result = parseRoiValue(val)
  roiError.value = result.error
  if (!result.error) props.node.roi = result.value
})
watch(roiOffsetInput, (val) => {
  const result = parseRectArrayValue(val, 'ROI offset')
  roiOffsetError.value = result.error
  if (!result.error) props.node.roi_offset = result.value
})
watch(replaceInput, (val) => {
  try { props.node.replace = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(labelsInput, (val) => {
  try { props.node.labels = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(expectedIndexInput, (val) => {
  try { props.node.expected = val ? JSON.parse(val) : undefined } catch (e) { props.node.expected = val || undefined }
})
watch(thresholdInput, (val) => {
  props.node.threshold = serializeNumberOrList(val)
})
watch(allOfInput, (val) => {
  try { props.node.allOf = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(anyOfInput, (val) => {
  try { props.node.anyOf = val ? JSON.parse(val) : undefined } catch (e) {}
})
watch(customRecognitionParamInput, (val) => {
  try { props.node.custom_recognition_param = val ? JSON.parse(val) : undefined } catch (e) {}
})
</script>
