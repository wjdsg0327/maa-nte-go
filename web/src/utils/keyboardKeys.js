const BASE_KEYS = [
  [8, 'Backspace', ['退格']],
  [9, 'Tab', ['制表']],
  [13, 'Enter', ['回车', 'Return']],
  [16, 'Shift', ['Shift']],
  [17, 'Ctrl', ['Control', '控制']],
  [18, 'Alt', ['Menu']],
  [19, 'Pause', ['Break']],
  [20, 'Caps Lock', ['CapsLock', '大写锁定']],
  [27, 'Esc', ['Escape', '退出']],
  [32, 'Space', ['空格']],
  [33, 'Page Up', ['PageUp', '上一页']],
  [34, 'Page Down', ['PageDown', '下一页']],
  [35, 'End', ['末尾']],
  [36, 'Home', ['主页']],
  [37, 'Left', ['ArrowLeft', '左']],
  [38, 'Up', ['ArrowUp', '上']],
  [39, 'Right', ['ArrowRight', '右']],
  [40, 'Down', ['ArrowDown', '下']],
  [44, 'Print Screen', ['PrintScreen', '截图']],
  [45, 'Insert', ['插入']],
  [46, 'Delete', ['Del', '删除']],
  [91, 'Left Win', ['MetaLeft', 'Win', 'Windows']],
  [92, 'Right Win', ['MetaRight']],
  [93, 'Menu', ['ContextMenu', '菜单']],
  [106, 'Numpad *', ['Multiply', '小键盘乘']],
  [107, 'Numpad +', ['Add', '小键盘加']],
  [109, 'Numpad -', ['Subtract', '小键盘减']],
  [110, 'Numpad .', ['Decimal', '小键盘点']],
  [111, 'Numpad /', ['Divide', '小键盘除']],
  [144, 'Num Lock', ['NumLock', '数字锁定']],
  [145, 'Scroll Lock', ['ScrollLock']],
  [186, ';', ['Semicolon', '分号']],
  [187, '=', ['Equal', '等号']],
  [188, ',', ['Comma', '逗号']],
  [189, '-', ['Minus', '减号']],
  [190, '.', ['Period', '句号']],
  [191, '/', ['Slash', '斜杠']],
  [192, '`', ['Backquote', '反引号']],
  [219, '[', ['BracketLeft', '左中括号']],
  [220, '\\', ['Backslash', '反斜杠']],
  [221, ']', ['BracketRight', '右中括号']],
  [222, "'", ['Quote', '引号']],
]

const digitKeys = Array.from({ length: 10 }, (_, index) => {
  const label = String(index)
  return [48 + index, label, [`Digit${index}`, `数字${index}`]]
})

const letterKeys = Array.from({ length: 26 }, (_, index) => {
  const label = String.fromCharCode(65 + index)
  return [65 + index, label, [`Key${label}`]]
})

const numpadKeys = Array.from({ length: 10 }, (_, index) => [
  96 + index,
  `Numpad ${index}`,
  [`Numpad${index}`, `小键盘${index}`],
])

const functionKeys = Array.from({ length: 24 }, (_, index) => {
  const label = `F${index + 1}`
  return [112 + index, label, [label.toLowerCase()]]
})

function normalizeSearchText(value) {
  return String(value ?? '').trim().toLowerCase().replace(/\s+/g, '')
}

function makeOption([code, label, aliases = []]) {
  const searchText = normalizeSearchText([label, code, ...aliases].join(' '))
  return { code, label, aliases, searchText }
}

export const KEY_OPTIONS = [
  ...BASE_KEYS,
  ...digitKeys,
  ...letterKeys,
  ...numpadKeys,
  ...functionKeys,
].map(makeOption)

const KEY_BY_CODE = new Map(KEY_OPTIONS.map((option) => [option.code, option]))

function numberFromInput(value) {
  const trimmed = String(value ?? '').trim()
  if (!/^\d+$/.test(trimmed)) return undefined
  const number = Number(trimmed)
  return Number.isSafeInteger(number) ? number : undefined
}

export function parseKeyboardKeyInput(value) {
  const input = String(value ?? '').trim()
  if (!input) return undefined

  const labelMatch = input.match(/\((\d+)\)\s*$/)
  if (labelMatch) return Number(labelMatch[1])

  const numeric = numberFromInput(input)
  if (numeric !== undefined) return numeric

  const search = normalizeSearchText(input)
  const exact = KEY_OPTIONS.find((option) => {
    if (normalizeSearchText(option.label) === search) return true
    return option.aliases.some((alias) => normalizeSearchText(alias) === search)
  })
  return exact?.code
}

export function filterKeyboardKeys(query, limit = 24) {
  const search = normalizeSearchText(query)
  const options = search
    ? KEY_OPTIONS.filter((option) => option.searchText.includes(search))
    : KEY_OPTIONS
  return options.slice(0, limit)
}

export function formatKeyCode(code) {
  const number = Number(code)
  if (!Number.isSafeInteger(number)) return ''
  const option = KEY_BY_CODE.get(number)
  return `${option?.label || 'Key'} (${number})`
}

export function normalizeKeyCodes(value) {
  const source = Array.isArray(value) ? value : [value]
  const result = []

  for (const item of source) {
    const code = typeof item === 'number' ? item : parseKeyboardKeyInput(item)
    if (!Number.isSafeInteger(code) || result.includes(code)) continue
    result.push(code)
  }

  return result
}

export function serializeKeyCodes(codes, multiple) {
  const normalized = normalizeKeyCodes(codes)
  if (normalized.length === 0) return undefined
  if (!multiple || normalized.length === 1) return normalized[0]
  return normalized
}
